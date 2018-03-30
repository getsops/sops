// Package firebaseremoteconfig provides access to the Firebase Remote Config API.
//
// See https://firebase.google.com/docs/remote-config/
//
// Usage example:
//
//   import "google.golang.org/api/firebaseremoteconfig/v1"
//   ...
//   firebaseremoteconfigService, err := firebaseremoteconfig.New(oauthHttpClient)
package firebaseremoteconfig // import "google.golang.org/api/firebaseremoteconfig/v1"

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

const apiId = "firebaseremoteconfig:v1"
const apiName = "firebaseremoteconfig"
const apiVersion = "v1"
const basePath = "https://firebaseremoteconfig.googleapis.com/"

func New(client *http.Client) (*Service, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	s := &Service{client: client, BasePath: basePath}
	s.Projects = NewProjectsService(s)
	return s, nil
}

type Service struct {
	client    *http.Client
	BasePath  string // API endpoint base URL
	UserAgent string // optional additional User-Agent fragment

	Projects *ProjectsService
}

func (s *Service) userAgent() string {
	if s.UserAgent == "" {
		return googleapi.UserAgent
	}
	return googleapi.UserAgent + " " + s.UserAgent
}

func NewProjectsService(s *Service) *ProjectsService {
	rs := &ProjectsService{s: s}
	return rs
}

type ProjectsService struct {
	s *Service
}

// RemoteConfig: *
// The RemoteConfig consists of a list of conditions (which can
// be
// thought of as named "if" statements) and a map of parameters
// (parameter key
// to a stucture containing an optional default value, as well as a
// optional
// submap of (condition name to value when that condition is true).
type RemoteConfig struct {
	// Conditions: The list of named conditions. The order *does* affect the
	// semantics.
	// The condition_name values of these entries must be unique.
	//
	// The resolved value of a config parameter P is determined as follow:
	// * Let Y be the set of values from the submap of P that refer to
	// conditions
	//   that evaluate to <code>true</code>.
	// * If Y is non empty, the value is taken from the specific submap in Y
	// whose
	//   condition_name is the earliest in this condition list.
	// * Else, if P has a default value option (condition_name is empty)
	// then
	//   the value is taken from that option.
	// * Else, parameter P has no value and is omitted from the config
	// result.
	//
	// Example: parameter key "p1", default value "v1", submap specified
	// as
	// {"c1": v2, "c2": v3} where "c1" and "c2" are names of conditions in
	// the
	// condition list (where "c1" in this example appears before "c2").
	// The
	// value of p1 would be v2 as long as c1 is true.  Otherwise, if c2 is
	// true,
	// p1 would evaluate to v3, and if c1 and c2 are both false, p1 would
	// evaluate
	// to v1.  If no default value was specified, and c1 and c2 were both
	// false,
	// no value for p1 would be generated.
	Conditions []*RemoteConfigCondition `json:"conditions,omitempty"`

	// Parameters: Map of parameter keys to their optional default values
	// and optional submap
	// of (condition name : value). Order doesn't affect semantics, and so
	// is
	// sorted by the server. The 'key' values of the params must be unique.
	Parameters map[string]RemoteConfigParameter `json:"parameters,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Conditions") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Conditions") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *RemoteConfig) MarshalJSON() ([]byte, error) {
	type noMethod RemoteConfig
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// RemoteConfigCondition: A single RemoteConfig Condition.  A list of
// these (because order matters) are
// part of a single RemoteConfig template.
type RemoteConfigCondition struct {
	// Expression: Required.
	Expression string `json:"expression,omitempty"`

	// Name: Required.
	// A non empty and unique name of this condition.
	Name string `json:"name,omitempty"`

	// TagColor: Optional.
	// The display (tag) color of this condition. This serves as part of a
	// tag
	// (in the future, we may add tag text as well as tag color, but that is
	// not
	// yet implemented in the UI).
	// This value has no affect on the semantics of the delivered config and
	// it
	// is ignored by the backend, except for passing it through
	// write/read
	// requests.
	// Not having this value or having the
	// "CONDITION_DISPLAY_COLOR_UNSPECIFIED"
	// value (0) have the same meaning:  Let the UI choose any valid color
	// when
	// displaying the condition.
	//
	// Possible values:
	//   "CONDITION_DISPLAY_COLOR_UNSPECIFIED"
	//   "BLUE" - Blue
	//   "BROWN" - Brown
	//   "CYAN" - Cyan
	//   "DEEP_ORANGE" - aka "Red Orange"
	//   "GREEN" - Green
	//   "INDIGO" - Indigo
	// *
	//   "LIME" - Lime - Approved deviation from Material color palette
	//   "ORANGE" - Orange
	//   "PINK" - Pink
	//   "PURPLE" - Purple
	//   "TEAL" - Teal
	TagColor string `json:"tagColor,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Expression") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Expression") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *RemoteConfigCondition) MarshalJSON() ([]byte, error) {
	type noMethod RemoteConfigCondition
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// RemoteConfigParameter: While default_value and conditional_values are
// each optional, at least one of
// the two is required - otherwise, the parameter is meaningless (and
// an
// exception will be thrown by the validation logic).
type RemoteConfigParameter struct {
	// ConditionalValues: Optional - a map of (condition_name, value). The
	// condition_name of the
	// highest priority (the one listed first in the conditions array)
	// determines
	// the value of this parameter.
	ConditionalValues map[string]RemoteConfigParameterValue `json:"conditionalValues,omitempty"`

	// DefaultValue: Optional - value to set the parameter to, when none of
	// the named conditions
	// evaluate to <code>true</code>.
	DefaultValue *RemoteConfigParameterValue `json:"defaultValue,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ConditionalValues")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ConditionalValues") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *RemoteConfigParameter) MarshalJSON() ([]byte, error) {
	type noMethod RemoteConfigParameter
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// RemoteConfigParameterValue: A RemoteConfigParameter's "value" (either
// the default value, or the value
// associated with a condition name) is either a string, or
// the
// "use_in_app_default" indicator (which means to leave out the
// parameter from
// the returned <key, value> map that is the output of the parameter
// fetch).
// We represent the "use_in_app_default" as a bool, but (when using the
// boolean
// instead of the string) it should always be <code>true</code>.
type RemoteConfigParameterValue struct {
	// UseInAppDefault: if true, omit the parameter from the map of fetched
	// parameter values
	UseInAppDefault bool `json:"useInAppDefault,omitempty"`

	// Value: the string to set the parameter to
	Value string `json:"value,omitempty"`

	// ForceSendFields is a list of field names (e.g. "UseInAppDefault") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "UseInAppDefault") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *RemoteConfigParameterValue) MarshalJSON() ([]byte, error) {
	type noMethod RemoteConfigParameterValue
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// method id "firebaseremoteconfig.projects.getRemoteConfig":

type ProjectsGetRemoteConfigCall struct {
	s            *Service
	projectid    string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// GetRemoteConfig: Get the latest version Remote Configuration for a
// project.
// Returns the RemoteConfig as the payload, and also the eTag as
// a
// response header.
func (r *ProjectsService) GetRemoteConfig(projectid string) *ProjectsGetRemoteConfigCall {
	c := &ProjectsGetRemoteConfigCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectid = projectid
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsGetRemoteConfigCall) Fields(s ...googleapi.Field) *ProjectsGetRemoteConfigCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsGetRemoteConfigCall) IfNoneMatch(entityTag string) *ProjectsGetRemoteConfigCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsGetRemoteConfigCall) Context(ctx context.Context) *ProjectsGetRemoteConfigCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsGetRemoteConfigCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsGetRemoteConfigCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/{+project}/remoteConfig")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"project": c.projectid,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "firebaseremoteconfig.projects.getRemoteConfig" call.
// Exactly one of *RemoteConfig or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *RemoteConfig.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsGetRemoteConfigCall) Do(opts ...googleapi.CallOption) (*RemoteConfig, error) {
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
	ret := &RemoteConfig{
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
	//   "description": "Get the latest version Remote Configuration for a project.\nReturns the RemoteConfig as the payload, and also the eTag as a\nresponse header.",
	//   "flatPath": "v1/projects/{projectsId}/remoteConfig",
	//   "httpMethod": "GET",
	//   "id": "firebaseremoteconfig.projects.getRemoteConfig",
	//   "parameterOrder": [
	//     "project"
	//   ],
	//   "parameters": {
	//     "project": {
	//       "description": "The GMP project identifier. Required.\nSee note at the beginning of this file regarding project ids.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/{+project}/remoteConfig",
	//   "response": {
	//     "$ref": "RemoteConfig"
	//   }
	// }

}

// method id "firebaseremoteconfig.projects.updateRemoteConfig":

type ProjectsUpdateRemoteConfigCall struct {
	s            *Service
	projectid    string
	remoteconfig *RemoteConfig
	urlParams_   gensupport.URLParams
	ctx_         context.Context
	header_      http.Header
}

// UpdateRemoteConfig: Update a RemoteConfig. We treat this as an
// always-existing
// resource (when it is not found in our data store, we treat it as
// version
// 0, a template with zero conditions and zero parameters). Hence there
// are
// no Create or Delete operations. Returns the updated template
// when
// successful (and the updated eTag as a response header), or an error
// if
// things go wrong.
// Possible error messages:
// * VALIDATION_ERROR (HTTP status 400) with additional details if
// the
// template being passed in can not be validated.
// * AUTHENTICATION_ERROR (HTTP status 401) if the request can not
// be
// authenticate (e.g. no access token, or invalid access token).
// * AUTHORIZATION_ERROR (HTTP status 403) if the request can not
// be
// authorized (e.g. the user has no access to the specified project
// id).
// * VERSION_MISMATCH (HTTP status 412) when trying to update when
// the
// expected eTag (passed in via the "If-match" header) is not specified,
// or
// is specified but does does not match the current eTag.
// * Internal error (HTTP status 500) for Database problems or other
// internal
// errors.
func (r *ProjectsService) UpdateRemoteConfig(projectid string, remoteconfig *RemoteConfig) *ProjectsUpdateRemoteConfigCall {
	c := &ProjectsUpdateRemoteConfigCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectid = projectid
	c.remoteconfig = remoteconfig
	return c
}

// ValidateOnly sets the optional parameter "validateOnly": Defaults to
// <code>false</code> (UpdateRemoteConfig call should
// update the backend if there are no validation/interal errors). May be
// set
// to <code>true</code> to indicate that, should no validation errors
// occur,
// the call should return a "200 OK" instead of performing the update.
// Note
// that other error messages (500 Internal Error, 412 Version Mismatch,
// etc)
// may still result after flipping to <code>false</code>, even if
// getting a
// "200 OK" when calling with <code>true</code>.
func (c *ProjectsUpdateRemoteConfigCall) ValidateOnly(validateOnly bool) *ProjectsUpdateRemoteConfigCall {
	c.urlParams_.Set("validateOnly", fmt.Sprint(validateOnly))
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsUpdateRemoteConfigCall) Fields(s ...googleapi.Field) *ProjectsUpdateRemoteConfigCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsUpdateRemoteConfigCall) Context(ctx context.Context) *ProjectsUpdateRemoteConfigCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsUpdateRemoteConfigCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsUpdateRemoteConfigCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.remoteconfig)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/{+project}/remoteConfig")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PUT", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"project": c.projectid,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "firebaseremoteconfig.projects.updateRemoteConfig" call.
// Exactly one of *RemoteConfig or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *RemoteConfig.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsUpdateRemoteConfigCall) Do(opts ...googleapi.CallOption) (*RemoteConfig, error) {
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
	ret := &RemoteConfig{
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
	//   "description": "Update a RemoteConfig. We treat this as an always-existing\nresource (when it is not found in our data store, we treat it as version\n0, a template with zero conditions and zero parameters). Hence there are\nno Create or Delete operations. Returns the updated template when\nsuccessful (and the updated eTag as a response header), or an error if\nthings go wrong.\nPossible error messages:\n* VALIDATION_ERROR (HTTP status 400) with additional details if the\ntemplate being passed in can not be validated.\n* AUTHENTICATION_ERROR (HTTP status 401) if the request can not be\nauthenticate (e.g. no access token, or invalid access token).\n* AUTHORIZATION_ERROR (HTTP status 403) if the request can not be\nauthorized (e.g. the user has no access to the specified project id).\n* VERSION_MISMATCH (HTTP status 412) when trying to update when the\nexpected eTag (passed in via the \"If-match\" header) is not specified, or\nis specified but does does not match the current eTag.\n* Internal error (HTTP status 500) for Database problems or other internal\nerrors.",
	//   "flatPath": "v1/projects/{projectsId}/remoteConfig",
	//   "httpMethod": "PUT",
	//   "id": "firebaseremoteconfig.projects.updateRemoteConfig",
	//   "parameterOrder": [
	//     "project"
	//   ],
	//   "parameters": {
	//     "project": {
	//       "description": "The GMP project identifier. Required.\nSee note at the beginning of this file regarding project ids.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "validateOnly": {
	//       "description": "Optional. Defaults to \u003ccode\u003efalse\u003c/code\u003e (UpdateRemoteConfig call should\nupdate the backend if there are no validation/interal errors). May be set\nto \u003ccode\u003etrue\u003c/code\u003e to indicate that, should no validation errors occur,\nthe call should return a \"200 OK\" instead of performing the update. Note\nthat other error messages (500 Internal Error, 412 Version Mismatch, etc)\nmay still result after flipping to \u003ccode\u003efalse\u003c/code\u003e, even if getting a\n\"200 OK\" when calling with \u003ccode\u003etrue\u003c/code\u003e.",
	//       "location": "query",
	//       "type": "boolean"
	//     }
	//   },
	//   "path": "v1/{+project}/remoteConfig",
	//   "request": {
	//     "$ref": "RemoteConfig"
	//   },
	//   "response": {
	//     "$ref": "RemoteConfig"
	//   }
	// }

}
