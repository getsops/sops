// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pubsub

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	gax "github.com/googleapis/gax-go"

	"golang.org/x/net/context"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/internal"
	"cloud.google.com/go/internal/testutil"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	topicIDs = testutil.NewUIDSpace("topic")
	subIDs   = testutil.NewUIDSpace("sub")
)

// messageData is used to hold the contents of a message so that it can be compared against the contents
// of another message without regard to irrelevant fields.
type messageData struct {
	ID         string
	Data       []byte
	Attributes map[string]string
}

func extractMessageData(m *Message) *messageData {
	return &messageData{
		ID:         m.ID,
		Data:       m.Data,
		Attributes: m.Attributes,
	}
}

func integrationTestClient(t *testing.T, ctx context.Context) *Client {
	if testing.Short() {
		t.Skip("Integration tests skipped in short mode")
	}
	projID := testutil.ProjID()
	if projID == "" {
		t.Skip("Integration tests skipped. See CONTRIBUTING.md for details")
	}
	ts := testutil.TokenSource(ctx, ScopePubSub, ScopeCloudPlatform)
	if ts == nil {
		t.Skip("Integration tests skipped. See CONTRIBUTING.md for details")
	}
	client, err := NewClient(ctx, projID, option.WithTokenSource(ts))
	if err != nil {
		t.Fatalf("Creating client error: %v", err)
	}
	return client
}

func TestAll(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client := integrationTestClient(t, ctx)
	defer client.Close()

	topic, err := client.CreateTopic(ctx, topicIDs.New())
	if err != nil {
		t.Errorf("CreateTopic error: %v", err)
	}
	defer topic.Stop()

	var sub *Subscription
	if sub, err = client.CreateSubscription(ctx, subIDs.New(), SubscriptionConfig{Topic: topic}); err != nil {
		t.Errorf("CreateSub error: %v", err)
	}

	exists, err := topic.Exists(ctx)
	if err != nil {
		t.Fatalf("TopicExists error: %v", err)
	}
	if !exists {
		t.Errorf("topic %v should exist, but it doesn't", topic)
	}

	exists, err = sub.Exists(ctx)
	if err != nil {
		t.Fatalf("SubExists error: %v", err)
	}
	if !exists {
		t.Errorf("subscription %s should exist, but it doesn't", sub.ID())
	}

	var msgs []*Message
	for i := 0; i < 10; i++ {
		text := fmt.Sprintf("a message with an index %d", i)
		attrs := make(map[string]string)
		attrs["foo"] = "bar"
		msgs = append(msgs, &Message{
			Data:       []byte(text),
			Attributes: attrs,
		})
	}

	// Publish the messages.
	type pubResult struct {
		m *Message
		r *PublishResult
	}
	var rs []pubResult
	for _, m := range msgs {
		r := topic.Publish(ctx, m)
		rs = append(rs, pubResult{m, r})
	}
	want := make(map[string]*messageData)
	for _, res := range rs {
		id, err := res.r.Get(ctx)
		if err != nil {
			t.Fatal(err)
		}
		md := extractMessageData(res.m)
		md.ID = id
		want[md.ID] = md
	}

	// Use a timeout to ensure that Pull does not block indefinitely if there are unexpectedly few messages available.
	timeoutCtx, _ := context.WithTimeout(ctx, time.Minute)
	gotMsgs, err := pullN(timeoutCtx, sub, len(want), func(ctx context.Context, m *Message) {
		m.Ack()
	})
	if err != nil {
		t.Fatalf("Pull: %v", err)
	}
	got := make(map[string]*messageData)
	for _, m := range gotMsgs {
		md := extractMessageData(m)
		got[md.ID] = md
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("messages: got: %v ; want: %v", got, want)
	}

	if msg, ok := testIAM(ctx, topic.IAM(), "pubsub.topics.get"); !ok {
		t.Errorf("topic IAM: %s", msg)
	}
	if msg, ok := testIAM(ctx, sub.IAM(), "pubsub.subscriptions.get"); !ok {
		t.Errorf("sub IAM: %s", msg)
	}

	snap, err := sub.createSnapshot(ctx, "")
	if err != nil {
		t.Fatalf("CreateSnapshot error: %v", err)
	}

	timeoutCtx, _ = context.WithTimeout(ctx, time.Minute)
	err = internal.Retry(timeoutCtx, gax.Backoff{}, func() (bool, error) {
		snapIt := client.snapshots(timeoutCtx)
		for {
			s, err := snapIt.Next()
			if err == nil && s.name == snap.name {
				return true, nil
			}
			if err == iterator.Done {
				return false, fmt.Errorf("cannot find snapshot: %q", snap.name)
			}
			if err != nil {
				return false, err
			}
		}
	})
	if err != nil {
		t.Error(err)
	}

	err = internal.Retry(timeoutCtx, gax.Backoff{}, func() (bool, error) {
		err := sub.seekToSnapshot(timeoutCtx, snap.snapshot)
		return err == nil, err
	})
	if err != nil {
		t.Error(err)
	}

	err = internal.Retry(timeoutCtx, gax.Backoff{}, func() (bool, error) {
		err := sub.seekToTime(timeoutCtx, time.Now())
		return err == nil, err
	})
	if err != nil {
		t.Error(err)
	}

	err = internal.Retry(timeoutCtx, gax.Backoff{}, func() (bool, error) {
		snapHandle := client.snapshot(snap.ID())
		err := snapHandle.delete(timeoutCtx)
		return err == nil, err
	})
	if err != nil {
		t.Error(err)
	}

	if err := sub.Delete(ctx); err != nil {
		t.Errorf("DeleteSub error: %v", err)
	}

	if err := topic.Delete(ctx); err != nil {
		t.Errorf("DeleteTopic error: %v", err)
	}
}

// IAM tests.
// NOTE: for these to succeed, the test runner identity must have the Pub/Sub Admin or Owner roles.
// To set, visit https://console.developers.google.com, select "IAM & Admin" from the top-left
// menu, choose the account, click the Roles dropdown, and select "Pub/Sub > Pub/Sub Admin".
// TODO(jba): move this to a testing package within cloud.google.com/iam, so we can re-use it.
func testIAM(ctx context.Context, h *iam.Handle, permission string) (msg string, ok bool) {
	// Attempting to add an non-existent identity  (e.g. "alice@example.com") causes the service
	// to return an internal error, so use a real identity.
	const member = "domain:google.com"

	var policy *iam.Policy
	var err error

	if policy, err = h.Policy(ctx); err != nil {
		return fmt.Sprintf("Policy: %v", err), false
	}
	// The resource is new, so the policy should be empty.
	if got := policy.Roles(); len(got) > 0 {
		return fmt.Sprintf("initially: got roles %v, want none", got), false
	}
	// Add a member, set the policy, then check that the member is present.
	policy.Add(member, iam.Viewer)
	if err := h.SetPolicy(ctx, policy); err != nil {
		return fmt.Sprintf("SetPolicy: %v", err), false
	}
	if policy, err = h.Policy(ctx); err != nil {
		return fmt.Sprintf("Policy: %v", err), false
	}
	if got, want := policy.Members(iam.Viewer), []string{member}; !reflect.DeepEqual(got, want) {
		return fmt.Sprintf("after Add: got %v, want %v", got, want), false
	}
	// Now remove that member, set the policy, and check that it's empty again.
	policy.Remove(member, iam.Viewer)
	if err := h.SetPolicy(ctx, policy); err != nil {
		return fmt.Sprintf("SetPolicy: %v", err), false
	}
	if policy, err = h.Policy(ctx); err != nil {
		return fmt.Sprintf("Policy: %v", err), false
	}
	if got := policy.Roles(); len(got) > 0 {
		return fmt.Sprintf("after Remove: got roles %v, want none", got), false
	}
	// Call TestPermissions.
	// Because this user is an admin, it has all the permissions on the
	// resource type. Note: the service fails if we ask for inapplicable
	// permissions (e.g. a subscription permission on a topic, or a topic
	// create permission on a topic rather than its parent).
	wantPerms := []string{permission}
	gotPerms, err := h.TestPermissions(ctx, wantPerms)
	if err != nil {
		return fmt.Sprintf("TestPermissions: %v", err), false
	}
	if !reflect.DeepEqual(gotPerms, wantPerms) {
		return fmt.Sprintf("TestPermissions: got %v, want %v", gotPerms, wantPerms), false
	}
	return "", true
}

func TestSubscriptionUpdate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client := integrationTestClient(t, ctx)
	defer client.Close()

	topic, err := client.CreateTopic(ctx, topicIDs.New())
	if err != nil {
		t.Fatalf("CreateTopic error: %v", err)
	}
	defer topic.Stop()
	defer topic.Delete(ctx)

	var sub *Subscription
	if sub, err = client.CreateSubscription(ctx, subIDs.New(), SubscriptionConfig{Topic: topic}); err != nil {
		t.Fatalf("CreateSub error: %v", err)
	}
	defer sub.Delete(ctx)

	sc, err := sub.Config(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(sc.PushConfig, PushConfig{}) {
		t.Fatalf("got %+v, want empty PushConfig")
	}
	// Add a PushConfig.
	projID := testutil.ProjID()
	pc := PushConfig{
		Endpoint:   "https://" + projID + ".appspot.com/_ah/push-handlers/push",
		Attributes: map[string]string{"x-goog-version": "v1"},
	}
	sc, err = sub.Update(ctx, SubscriptionConfigToUpdate{PushConfig: &pc})
	if err != nil {
		t.Fatal(err)
	}
	// Despite the docs which say that Get always returns a valid "x-goog-version"
	// attribute, none is returned. See
	// https://cloud.google.com/pubsub/docs/reference/rpc/google.pubsub.v1#google.pubsub.v1.PushConfig
	pc.Attributes = nil
	if got, want := sc.PushConfig, pc; !reflect.DeepEqual(got, want) {
		t.Fatalf("setting push config: got\n%+v\nwant\n%+v", got, want)
	}
	// Remove the PushConfig, turning the subscription back into pull mode.
	pc = PushConfig{}
	sc, err = sub.Update(ctx, SubscriptionConfigToUpdate{PushConfig: &pc})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := sc.PushConfig, pc; !reflect.DeepEqual(got, want) {
		t.Fatalf("removing push config: got\n%+v\nwant %+v", got, want)
	}

	// If nothing changes, our client returns an error.
	_, err = sub.Update(ctx, SubscriptionConfigToUpdate{})
	if err == nil {
		t.Fatal("got nil, wanted error")
	}
}

func TestPublicTopic(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client := integrationTestClient(t, ctx)
	defer client.Close()

	sub, err := client.CreateSubscription(ctx, subIDs.New(), SubscriptionConfig{
		Topic: client.TopicInProject("taxirides-realtime", "pubsub-public-data"),
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Delete(ctx)
	// Confirm that Receive works. It doesn't matter if we actually get any
	// messages.
	ctxt, cancel := context.WithTimeout(ctx, 5*time.Second)
	err = sub.Receive(ctxt, func(_ context.Context, msg *Message) {
		msg.Ack()
		cancel()
	})
	if err != nil {
		t.Fatal(err)
	}
}
