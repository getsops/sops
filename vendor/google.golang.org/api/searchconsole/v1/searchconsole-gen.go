// Package searchconsole provides access to the Google Search Console URL Testing Tools API.
//
// See https://developers.google.com/webmaster-tools/search-console-api/
//
// Usage example:
//
//   import "google.golang.org/api/searchconsole/v1"
//   ...
//   searchconsoleService, err := searchconsole.New(oauthHttpClient)
package searchconsole // import "google.golang.org/api/searchconsole/v1"

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	context "golang.org/x/net/context"
	ctxhttp "golang.org/x/net/context/ctxhttp"
	gensupport "google.golang.org/api/gensupport"
	googleapi "google.golang.org/api/googleapi"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Always reference these packages, just in case the auto-generated code
// below doesn't.
var _ = bytes.NewBuffer
var _ = strconv.Itoa
var _ = fmt.Sprintf
var _ = json.NewDecoder
var _ = io.Copy
var _ = url.Parse
var _ = gensupport.MarshalJSON
var _ = googleapi.Version
var _ = errors.New
var _ = strings.Replace
var _ = context.Canceled
var _ = ctxhttp.Do

const apiId = "searchconsole:v1"
const apiName = "searchconsole"
const apiVersion = "v1"
const basePath = "https://searchconsole.googleapis.com/"

func New(client *http.Client) (*Service, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	s := &Service{client: client, BasePath: basePath}
	s.UrlTestingTools = NewUrlTestingToolsService(s)
	return s, nil
}

type Service struct {
	client    *http.Client
	BasePath  string // API endpoint base URL
	UserAgent string // optional additional User-Agent fragment

	UrlTestingTools *UrlTestingToolsService
}

func (s *Service) userAgent() string {
	if s.UserAgent == "" {
		return googleapi.UserAgent
	}
	return googleapi.UserAgent + " " + s.UserAgent
}

func NewUrlTestingToolsService(s *Service) *UrlTestingToolsService {
	rs := &UrlTestingToolsService{s: s}
	rs.MobileFriendlyTest = NewUrlTestingToolsMobileFriendlyTestService(s)
	return rs
}

type UrlTestingToolsService struct {
	s *Service

	MobileFriendlyTest *UrlTestingToolsMobileFriendlyTestService
}

func NewUrlTestingToolsMobileFriendlyTestService(s *Service) *UrlTestingToolsMobileFriendlyTestService {
	rs := &UrlTestingToolsMobileFriendlyTestService{s: s}
	return rs
}

type UrlTestingToolsMobileFriendlyTestService struct {
	s *Service
}

// BlockedResource: Blocked resource.
type BlockedResource struct {
	// Url: URL of the blocked resource.
	Url string `json:"url,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Url") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Url") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *BlockedResource) MarshalJSON() ([]byte, error) {
	type noMethod BlockedResource
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Image: Describe image data.
type Image struct {
	// Data: Image data in format determined by the mime type. Currently,
	// the format
	// will always be "image/png", but this might change in the future.
	Data string `json:"data,omitempty"`

	// MimeType: The mime-type of the image data.
	MimeType string `json:"mimeType,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Data") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Data") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Image) MarshalJSON() ([]byte, error) {
	type noMethod Image
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// MobileFriendlyIssue: Mobile-friendly issue.
type MobileFriendlyIssue struct {
	// Rule: Rule violated.
	//
	// Possible values:
	//   "MOBILE_FRIENDLY_RULE_UNSPECIFIED" - Unknown rule. Sorry, we don't
	// have any description for the rule that was
	// broken.
	//   "USES_INCOMPATIBLE_PLUGINS" - Plugins incompatible with mobile
	// devices are being used. [Learn
	// more]
	// (https://support.google.com/webmasters/answer/6352293#flash_usag
	// e).
	//   "CONFIGURE_VIEWPORT" - Viewsport is not specified using the meta
	// viewport tag. [Learn
	// more]
	// (https://support.google.com/webmasters/answer/6352293#viewport_n
	// ot_configured).
	//   "FIXED_WIDTH_VIEWPORT" - Viewport defined to a fixed width. [Learn
	// more]
	// (https://support.google.com/webmasters/answer/6352293#fixed-widt
	// h_viewport).
	//   "SIZE_CONTENT_TO_VIEWPORT" - Content not sized to viewport. [Learn
	// more]
	// (https://support.google.com/webmasters/answer/6352293#content_no
	// t_sized_to_viewport).
	//   "USE_LEGIBLE_FONT_SIZES" - Font size is too small for easy reading
	// on a small screen. [Learn
	// More]
	// (https://support.google.com/webmasters/answer/6352293#small_font
	// _size).
	//   "TAP_TARGETS_TOO_CLOSE" - Touch elements are too close to each
	// other. [Learn
	// more]
	// (https://support.google.com/webmasters/answer/6352293#touch_elem
	// ents_too_close).
	Rule string `json:"rule,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Rule") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Rule") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *MobileFriendlyIssue) MarshalJSON() ([]byte, error) {
	type noMethod MobileFriendlyIssue
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ResourceIssue: Information about a resource with issue.
type ResourceIssue struct {
	// BlockedResource: Describes a blocked resource issue.
	BlockedResource *BlockedResource `json:"blockedResource,omitempty"`

	// ForceSendFields is a list of field names (e.g. "BlockedResource") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "BlockedResource") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *ResourceIssue) MarshalJSON() ([]byte, error) {
	type noMethod ResourceIssue
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// RunMobileFriendlyTestRequest: Mobile-friendly test request.
type RunMobileFriendlyTestRequest struct {
	// RequestScreenshot: Whether or not screenshot is requested. Default is
	// false.
	RequestScreenshot bool `json:"requestScreenshot,omitempty"`

	// Url: URL for inspection.
	Url string `json:"url,omitempty"`

	// ForceSendFields is a list of field names (e.g. "RequestScreenshot")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "RequestScreenshot") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *RunMobileFriendlyTestRequest) MarshalJSON() ([]byte, error) {
	type noMethod RunMobileFriendlyTestRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// RunMobileFriendlyTestResponse: Mobile-friendly test response,
// including mobile-friendly issues and resource
// issues.
type RunMobileFriendlyTestResponse struct {
	// MobileFriendliness: Test verdict, whether the page is mobile friendly
	// or not.
	//
	// Possible values:
	//   "MOBILE_FRIENDLY_TEST_RESULT_UNSPECIFIED" - Internal error when
	// running this test. Please try running the test again.
	//   "MOBILE_FRIENDLY" - The page is mobile friendly.
	//   "NOT_MOBILE_FRIENDLY" - The page is not mobile friendly.
	MobileFriendliness string `json:"mobileFriendliness,omitempty"`

	// MobileFriendlyIssues: List of mobile-usability issues.
	MobileFriendlyIssues []*MobileFriendlyIssue `json:"mobileFriendlyIssues,omitempty"`

	// ResourceIssues: Information about embedded resources issues.
	ResourceIssues []*ResourceIssue `json:"resourceIssues,omitempty"`

	// Screenshot: Screenshot of the requested URL.
	Screenshot *Image `json:"screenshot,omitempty"`

	// TestStatus: Final state of the test, can be either complete or an
	// error.
	TestStatus *TestStatus `json:"testStatus,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "MobileFriendliness")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "MobileFriendliness") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *RunMobileFriendlyTestResponse) MarshalJSON() ([]byte, error) {
	type noMethod RunMobileFriendlyTestResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// TestStatus: Final state of the test, including error details if
// necessary.
type TestStatus struct {
	// Details: Error details if applicable.
	Details string `json:"details,omitempty"`

	// Status: Status of the test.
	//
	// Possible values:
	//   "TEST_STATUS_UNSPECIFIED" - Internal error when running this test.
	// Please try running the test again.
	//   "COMPLETE" - Inspection has completed without errors.
	//   "INTERNAL_ERROR" - Inspection terminated in an error state. This
	// indicates a problem in
	// Google's infrastructure, not a user error. Please try again later.
	//   "PAGE_UNREACHABLE" - Google can not access the URL because of a
	// user error such as a robots.txt
	// blockage, a 403 or 500 code etc. Please make sure that the URL
	// provided is
	// accessible by Googlebot and is not password protected.
	Status string `json:"status,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Details") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Details") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *TestStatus) MarshalJSON() ([]byte, error) {
	type noMethod TestStatus
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// method id "searchconsole.urlTestingTools.mobileFriendlyTest.run":

type UrlTestingToolsMobileFriendlyTestRunCall struct {
	s                            *Service
	runmobilefriendlytestrequest *RunMobileFriendlyTestRequest
	urlParams_                   gensupport.URLParams
	ctx_                         context.Context
	header_                      http.Header
}

// Run: Runs Mobile-Friendly Test for a given URL.
func (r *UrlTestingToolsMobileFriendlyTestService) Run(runmobilefriendlytestrequest *RunMobileFriendlyTestRequest) *UrlTestingToolsMobileFriendlyTestRunCall {
	c := &UrlTestingToolsMobileFriendlyTestRunCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.runmobilefriendlytestrequest = runmobilefriendlytestrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UrlTestingToolsMobileFriendlyTestRunCall) Fields(s ...googleapi.Field) *UrlTestingToolsMobileFriendlyTestRunCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *UrlTestingToolsMobileFriendlyTestRunCall) Context(ctx context.Context) *UrlTestingToolsMobileFriendlyTestRunCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *UrlTestingToolsMobileFriendlyTestRunCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *UrlTestingToolsMobileFriendlyTestRunCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.runmobilefriendlytestrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/urlTestingTools/mobileFriendlyTest:run")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "searchconsole.urlTestingTools.mobileFriendlyTest.run" call.
// Exactly one of *RunMobileFriendlyTestResponse or error will be
// non-nil. Any non-2xx status code is an error. Response headers are in
// either *RunMobileFriendlyTestResponse.ServerResponse.Header or (if a
// response was returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *UrlTestingToolsMobileFriendlyTestRunCall) Do(opts ...googleapi.CallOption) (*RunMobileFriendlyTestResponse, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &RunMobileFriendlyTestResponse{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Runs Mobile-Friendly Test for a given URL.",
	//   "flatPath": "v1/urlTestingTools/mobileFriendlyTest:run",
	//   "httpMethod": "POST",
	//   "id": "searchconsole.urlTestingTools.mobileFriendlyTest.run",
	//   "parameterOrder": [],
	//   "parameters": {},
	//   "path": "v1/urlTestingTools/mobileFriendlyTest:run",
	//   "request": {
	//     "$ref": "RunMobileFriendlyTestRequest"
	//   },
	//   "response": {
	//     "$ref": "RunMobileFriendlyTestResponse"
	//   }
	// }

}
