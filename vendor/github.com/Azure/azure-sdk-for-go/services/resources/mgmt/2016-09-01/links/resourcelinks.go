package links

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// ResourceLinksClient is the azure resources can be linked together to form logical relationships. You can establish
// links between resources belonging to different resource groups. However, all the linked resources must belong to the
// same subscription. Each resource can be linked to 50 other resources. If any of the linked resources are deleted or
// moved, the link owner must clean up the remaining link.
type ResourceLinksClient struct {
	BaseClient
}

// NewResourceLinksClient creates an instance of the ResourceLinksClient client.
func NewResourceLinksClient(subscriptionID string) ResourceLinksClient {
	return NewResourceLinksClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewResourceLinksClientWithBaseURI creates an instance of the ResourceLinksClient client.
func NewResourceLinksClientWithBaseURI(baseURI string, subscriptionID string) ResourceLinksClient {
	return ResourceLinksClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate creates or updates a resource link between the specified resources.
// Parameters:
// linkID - the fully qualified ID of the resource link. Use the format,
// /subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/{provider-namespace}/{resource-type}/{resource-name}/Microsoft.Resources/links/{link-name}.
// For example,
// /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup/Microsoft.Web/sites/mySite/Microsoft.Resources/links/myLink
// parameters - parameters for creating or updating a resource link.
func (client ResourceLinksClient) CreateOrUpdate(ctx context.Context, linkID string, parameters ResourceLink) (result ResourceLink, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ResourceLinksClient.CreateOrUpdate")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.Properties", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "parameters.Properties.TargetID", Name: validation.Null, Rule: true, Chain: nil}}}}}}); err != nil {
		return result, validation.NewError("links.ResourceLinksClient", "CreateOrUpdate", err.Error())
	}

	req, err := client.CreateOrUpdatePreparer(ctx, linkID, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "CreateOrUpdate", resp, "Failure sending request")
		return
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "CreateOrUpdate", resp, "Failure responding to request")
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client ResourceLinksClient) CreateOrUpdatePreparer(ctx context.Context, linkID string, parameters ResourceLink) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"linkId": linkID,
	}

	const APIVersion = "2016-09-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	parameters.ID = nil
	parameters.Name = nil
	parameters.Type = ""
	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{linkId}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client ResourceLinksClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client ResourceLinksClient) CreateOrUpdateResponder(resp *http.Response) (result ResourceLink, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes a resource link with the specified ID.
// Parameters:
// linkID - the fully qualified ID of the resource link. Use the format,
// /subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/{provider-namespace}/{resource-type}/{resource-name}/Microsoft.Resources/links/{link-name}.
// For example,
// /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup/Microsoft.Web/sites/mySite/Microsoft.Resources/links/myLink
func (client ResourceLinksClient) Delete(ctx context.Context, linkID string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ResourceLinksClient.Delete")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeletePreparer(ctx, linkID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client ResourceLinksClient) DeletePreparer(ctx context.Context, linkID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"linkId": linkID,
	}

	const APIVersion = "2016-09-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{linkId}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client ResourceLinksClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client ResourceLinksClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets a resource link with the specified ID.
// Parameters:
// linkID - the fully qualified Id of the resource link. For example,
// /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup/Microsoft.Web/sites/mySite/Microsoft.Resources/links/myLink
func (client ResourceLinksClient) Get(ctx context.Context, linkID string) (result ResourceLink, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ResourceLinksClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, linkID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client ResourceLinksClient) GetPreparer(ctx context.Context, linkID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"linkId": linkID,
	}

	const APIVersion = "2016-09-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{linkId}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ResourceLinksClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ResourceLinksClient) GetResponder(resp *http.Response) (result ResourceLink, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListAtSourceScope gets a list of resource links at and below the specified source scope.
// Parameters:
// scope - the fully qualified ID of the scope for getting the resource links. For example, to list resource
// links at and under a resource group, set the scope to
// /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup.
// filter - the filter to apply when getting resource links. To get links only at the specified scope (not
// below the scope), use Filter.atScope().
func (client ResourceLinksClient) ListAtSourceScope(ctx context.Context, scope string, filter Filter) (result ResourceLinkResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ResourceLinksClient.ListAtSourceScope")
		defer func() {
			sc := -1
			if result.rlr.Response.Response != nil {
				sc = result.rlr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listAtSourceScopeNextResults
	req, err := client.ListAtSourceScopePreparer(ctx, scope, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "ListAtSourceScope", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListAtSourceScopeSender(req)
	if err != nil {
		result.rlr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "ListAtSourceScope", resp, "Failure sending request")
		return
	}

	result.rlr, err = client.ListAtSourceScopeResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "ListAtSourceScope", resp, "Failure responding to request")
	}

	return
}

// ListAtSourceScopePreparer prepares the ListAtSourceScope request.
func (client ResourceLinksClient) ListAtSourceScopePreparer(ctx context.Context, scope string, filter Filter) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"scope": scope,
	}

	const APIVersion = "2016-09-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(string(filter)) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{scope}/providers/Microsoft.Resources/links", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListAtSourceScopeSender sends the ListAtSourceScope request. The method will close the
// http.Response Body if it receives an error.
func (client ResourceLinksClient) ListAtSourceScopeSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListAtSourceScopeResponder handles the response to the ListAtSourceScope request. The method always
// closes the http.Response Body.
func (client ResourceLinksClient) ListAtSourceScopeResponder(resp *http.Response) (result ResourceLinkResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listAtSourceScopeNextResults retrieves the next set of results, if any.
func (client ResourceLinksClient) listAtSourceScopeNextResults(ctx context.Context, lastResults ResourceLinkResult) (result ResourceLinkResult, err error) {
	req, err := lastResults.resourceLinkResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "links.ResourceLinksClient", "listAtSourceScopeNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListAtSourceScopeSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "links.ResourceLinksClient", "listAtSourceScopeNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListAtSourceScopeResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "listAtSourceScopeNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListAtSourceScopeComplete enumerates all values, automatically crossing page boundaries as required.
func (client ResourceLinksClient) ListAtSourceScopeComplete(ctx context.Context, scope string, filter Filter) (result ResourceLinkResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ResourceLinksClient.ListAtSourceScope")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListAtSourceScope(ctx, scope, filter)
	return
}

// ListAtSubscription gets all the linked resources for the subscription.
// Parameters:
// filter - the filter to apply on the list resource links operation. The supported filter for list resource
// links is targetId. For example, $filter=targetId eq {value}
func (client ResourceLinksClient) ListAtSubscription(ctx context.Context, filter string) (result ResourceLinkResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ResourceLinksClient.ListAtSubscription")
		defer func() {
			sc := -1
			if result.rlr.Response.Response != nil {
				sc = result.rlr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listAtSubscriptionNextResults
	req, err := client.ListAtSubscriptionPreparer(ctx, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "ListAtSubscription", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListAtSubscriptionSender(req)
	if err != nil {
		result.rlr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "ListAtSubscription", resp, "Failure sending request")
		return
	}

	result.rlr, err = client.ListAtSubscriptionResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "ListAtSubscription", resp, "Failure responding to request")
	}

	return
}

// ListAtSubscriptionPreparer prepares the ListAtSubscription request.
func (client ResourceLinksClient) ListAtSubscriptionPreparer(ctx context.Context, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-09-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Resources/links", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListAtSubscriptionSender sends the ListAtSubscription request. The method will close the
// http.Response Body if it receives an error.
func (client ResourceLinksClient) ListAtSubscriptionSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListAtSubscriptionResponder handles the response to the ListAtSubscription request. The method always
// closes the http.Response Body.
func (client ResourceLinksClient) ListAtSubscriptionResponder(resp *http.Response) (result ResourceLinkResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listAtSubscriptionNextResults retrieves the next set of results, if any.
func (client ResourceLinksClient) listAtSubscriptionNextResults(ctx context.Context, lastResults ResourceLinkResult) (result ResourceLinkResult, err error) {
	req, err := lastResults.resourceLinkResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "links.ResourceLinksClient", "listAtSubscriptionNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListAtSubscriptionSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "links.ResourceLinksClient", "listAtSubscriptionNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListAtSubscriptionResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.ResourceLinksClient", "listAtSubscriptionNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListAtSubscriptionComplete enumerates all values, automatically crossing page boundaries as required.
func (client ResourceLinksClient) ListAtSubscriptionComplete(ctx context.Context, filter string) (result ResourceLinkResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ResourceLinksClient.ListAtSubscription")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListAtSubscription(ctx, filter)
	return
}
