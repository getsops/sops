package consul

import (
	"fmt"
	"os"
	"strings"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/ory/dockertest"
)

func PrepareTestContainer(t *testing.T, version string) (cleanup func(), retAddress string, consulToken string) {
	t.Logf("preparing test container")
	consulToken = os.Getenv("CONSUL_HTTP_TOKEN")
	retAddress = os.Getenv("CONSUL_HTTP_ADDR")
	if retAddress != "" {
		return func() {}, retAddress, consulToken
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	config := `acl { enabled = true default_policy = "deny" }`
	if strings.HasPrefix(version, "1.3") {
		config = `datacenter = "test" acl_default_policy = "deny" acl_datacenter = "test" acl_master_token = "test"`
	}

	dockerOptions := &dockertest.RunOptions{
		Repository: "consul",
		Tag:        version,
		Cmd:        []string{"agent", "-dev", "-client", "0.0.0.0", "-hcl", config},
	}
	resource, err := pool.RunWithOptions(dockerOptions)
	if err != nil {
		t.Fatalf("Could not start local Consul %s docker container: %s", version, err)
	}

	cleanup = func() {
		docker.CleanupResource(t, pool, resource)
	}

	retAddress = fmt.Sprintf("localhost:%s", resource.GetPort("8500/tcp"))

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		consulConfig := consulapi.DefaultNonPooledConfig()
		consulConfig.Address = retAddress
		consul, err := consulapi.NewClient(consulConfig)
		if err != nil {
			return err
		}

		// For version of Consul < 1.4
		if strings.HasPrefix(version, "1.3") {
			consulToken = "test"
			_, err = consul.KV().Put(&consulapi.KVPair{
				Key:   "setuptest",
				Value: []byte("setuptest"),
			}, &consulapi.WriteOptions{
				Token: consulToken,
			})
			if err != nil {
				return err
			}
			return nil
		}

		// New default behavior
		aclbootstrap, _, err := consul.ACL().Bootstrap()
		if err != nil {
			return err
		}
		consulToken = aclbootstrap.SecretID
		t.Logf("Generated Master token: %s", consulToken)
		policy := &consulapi.ACLPolicy{
			Name:        "test",
			Description: "test",
			Rules: `node_prefix "" {
                policy = "write"
              }

              service_prefix "" {
                policy = "read"
              }
      `,
		}
		q := &consulapi.WriteOptions{
			Token: consulToken,
		}
		_, _, err = consul.ACL().PolicyCreate(policy, q)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to docker: %s", err)
	}
	return cleanup, retAddress, consulToken
}
