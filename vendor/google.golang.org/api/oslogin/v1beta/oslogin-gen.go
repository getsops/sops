// Package oslogin provides access to the Google Cloud OS Login API.
//
// See https://cloud.google.com/compute/docs/oslogin/rest/
//
// Usage example:
//
//   import "google.golang.org/api/oslogin/v1beta"
//   ...
//   osloginService, err := oslogin.New(oauthHttpClient)
package oslogin // import "google.golang.org/api/oslogin/v1beta"

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

const apiId = "oslogin:v1beta"
const apiName = "oslogin"
const apiVersion = "v1beta"
const basePath = "https://oslogin.googleapis.com/"

// OAuth2 scopes used by this API.
const (
	// View and manage your data across Google Cloud Platform services
	CloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"

	// View your data across Google Cloud Platform services
	CloudPlatformReadOnlyScope = "https://www.googleapis.com/auth/cloud-platform.read-only"

	// View and manage your Google Compute Engine resources
	ComputeScope = "https://www.googleapis.com/auth/compute"

	// View your Google Compute Engine resources
	ComputeReadonlyScope = "https://www.googleapis.com/auth/compute.readonly"
)

func New(client *http.Client) (*Service, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	s := &Service{client: client, BasePath: basePath}
	s.Users = NewUsersService(s)
	return s, nil
}

type Service struct {
	client    *http.Client
	BasePath  string // API endpoint base URL
	UserAgent string // optional additional User-Agent fragment

	Users *UsersService
}

func (s *Service) userAgent() string {
	if s.UserAgent == "" {
		return googleapi.UserAgent
	}
	return googleapi.UserAgent + " " + s.UserAgent
}

func NewUsersService(s *Service) *UsersService {
	rs := &UsersService{s: s}
	rs.SshPublicKeys = NewUsersSshPublicKeysService(s)
	return rs
}

type UsersService struct {
	s *Service

	SshPublicKeys *UsersSshPublicKeysService
}

func NewUsersSshPublicKeysService(s *Service) *UsersSshPublicKeysService {
	rs := &UsersSshPublicKeysService{s: s}
	return rs
}

type UsersSshPublicKeysService struct {
	s *Service
}

// Empty: A generic empty message that you can re-use to avoid defining
// duplicated
// empty messages in your APIs. A typical example is to use it as the
// request
// or the response type of an API method. For instance:
//
//     service Foo {
//       rpc Bar(google.protobuf.Empty) returns
// (google.protobuf.Empty);
//     }
//
// The JSON representation for `Empty` is empty JSON object `{}`.
type Empty struct {
	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`
}

// ImportSshPublicKeyResponse: A response message for importing an SSH
// public key.
type ImportSshPublicKeyResponse struct {
	// LoginProfile: The login profile information for the user.
	LoginProfile *LoginProfile `json:"loginProfile,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "LoginProfile") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "LoginProfile") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ImportSshPublicKeyResponse) MarshalJSON() ([]byte, error) {
	type noMethod ImportSshPublicKeyResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// LoginProfile: The user profile information used for logging in to a
// virtual machine on
// Google Compute Engine.
type LoginProfile struct {
	// Name: The primary email address that uniquely identifies the user.
	Name string `json:"name,omitempty"`

	// PosixAccounts: The list of POSIX accounts associated with the user.
	PosixAccounts []*PosixAccount `json:"posixAccounts,omitempty"`

	// SshPublicKeys: A map from SSH public key fingerprint to the
	// associated key object.
	SshPublicKeys map[string]SshPublicKey `json:"sshPublicKeys,omitempty"`

	// Suspended: Indicates if the user is suspended. A suspended user
	// cannot log in but
	// their profile information is retained.
	Suspended bool `json:"suspended,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Name") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Name") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *LoginProfile) MarshalJSON() ([]byte, error) {
	type noMethod LoginProfile
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// PosixAccount: The POSIX account information associated with a
// Directory API User.
type PosixAccount struct {
	// Gecos: The GECOS (user information) entry for this account.
	Gecos string `json:"gecos,omitempty"`

	// Gid: The default group ID.
	Gid int64 `json:"gid,omitempty,string"`

	// HomeDirectory: The path to the home directory for this account.
	HomeDirectory string `json:"homeDirectory,omitempty"`

	// Primary: Only one POSIX account can be marked as primary.
	Primary bool `json:"primary,omitempty"`

	// Shell: The path to the logic shell for this account.
	Shell string `json:"shell,omitempty"`

	// SystemId: System identifier for which account the username or uid
	// applies to.
	// By default, the empty value is used.
	SystemId string `json:"systemId,omitempty"`

	// Uid: The user ID.
	Uid int64 `json:"uid,omitempty,string"`

	// Username: The username of the POSIX account.
	Username string `json:"username,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Gecos") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Gecos") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *PosixAccount) MarshalJSON() ([]byte, error) {
	type noMethod PosixAccount
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// SshPublicKey: The SSH public key information associated with a
// Directory API User.
type SshPublicKey struct {
	// ExpirationTimeUsec: An expiration time in microseconds since epoch.
	ExpirationTimeUsec int64 `json:"expirationTimeUsec,omitempty,string"`

	// Fingerprint: The SHA-256 fingerprint of the SSH public key.
	// Output only.
	Fingerprint string `json:"fingerprint,omitempty"`

	// Key: Public key text in SSH format, defined by
	// <a href="https://www.ietf.org/rfc/rfc4253.txt"
	// target="_blank">RFC4253</a>
	// section 6.6.
	Key string `json:"key,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "ExpirationTimeUsec")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ExpirationTimeUsec") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *SshPublicKey) MarshalJSON() ([]byte, error) {
	type noMethod SshPublicKey
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// method id "oslogin.users.getLoginProfile":

type UsersGetLoginProfileCall struct {
	s            *Service
	name         string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// GetLoginProfile: Retrieves the profile information used for logging
// in to a virtual machine
// on Google Compute Engine.
func (r *UsersService) GetLoginProfile(name string) *UsersGetLoginProfileCall {
	c := &UsersGetLoginProfileCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UsersGetLoginProfileCall) Fields(s ...googleapi.Field) *UsersGetLoginProfileCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *UsersGetLoginProfileCall) IfNoneMatch(entityTag string) *UsersGetLoginProfileCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *UsersGetLoginProfileCall) Context(ctx context.Context) *UsersGetLoginProfileCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *UsersGetLoginProfileCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *UsersGetLoginProfileCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/{+name}/loginProfile")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "oslogin.users.getLoginProfile" call.
// Exactly one of *LoginProfile or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *LoginProfile.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *UsersGetLoginProfileCall) Do(opts ...googleapi.CallOption) (*LoginProfile, error) {
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
	ret := &LoginProfile{
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
	//   "description": "Retrieves the profile information used for logging in to a virtual machine\non Google Compute Engine.",
	//   "flatPath": "v1beta/users/{usersId}/loginProfile",
	//   "httpMethod": "GET",
	//   "id": "oslogin.users.getLoginProfile",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The unique ID for the user in format `users/{user}`.",
	//       "location": "path",
	//       "pattern": "^users/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/{+name}/loginProfile",
	//   "response": {
	//     "$ref": "LoginProfile"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only",
	//     "https://www.googleapis.com/auth/compute",
	//     "https://www.googleapis.com/auth/compute.readonly"
	//   ]
	// }

}

// method id "oslogin.users.importSshPublicKey":

type UsersImportSshPublicKeyCall struct {
	s            *Service
	parent       string
	sshpublickey *SshPublicKey
	urlParams_   gensupport.URLParams
	ctx_         context.Context
	header_      http.Header
}

// ImportSshPublicKey: Adds an SSH public key and returns the profile
// information. Default POSIX
// account information is set when no username and UID exist as part of
// the
// login profile.
func (r *UsersService) ImportSshPublicKey(parent string, sshpublickey *SshPublicKey) *UsersImportSshPublicKeyCall {
	c := &UsersImportSshPublicKeyCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	c.sshpublickey = sshpublickey
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UsersImportSshPublicKeyCall) Fields(s ...googleapi.Field) *UsersImportSshPublicKeyCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *UsersImportSshPublicKeyCall) Context(ctx context.Context) *UsersImportSshPublicKeyCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *UsersImportSshPublicKeyCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *UsersImportSshPublicKeyCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.sshpublickey)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/{+parent}:importSshPublicKey")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"parent": c.parent,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "oslogin.users.importSshPublicKey" call.
// Exactly one of *ImportSshPublicKeyResponse or error will be non-nil.
// Any non-2xx status code is an error. Response headers are in either
// *ImportSshPublicKeyResponse.ServerResponse.Header or (if a response
// was returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *UsersImportSshPublicKeyCall) Do(opts ...googleapi.CallOption) (*ImportSshPublicKeyResponse, error) {
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
	ret := &ImportSshPublicKeyResponse{
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
	//   "description": "Adds an SSH public key and returns the profile information. Default POSIX\naccount information is set when no username and UID exist as part of the\nlogin profile.",
	//   "flatPath": "v1beta/users/{usersId}:importSshPublicKey",
	//   "httpMethod": "POST",
	//   "id": "oslogin.users.importSshPublicKey",
	//   "parameterOrder": [
	//     "parent"
	//   ],
	//   "parameters": {
	//     "parent": {
	//       "description": "The unique ID for the user in format `users/{user}`.",
	//       "location": "path",
	//       "pattern": "^users/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/{+parent}:importSshPublicKey",
	//   "request": {
	//     "$ref": "SshPublicKey"
	//   },
	//   "response": {
	//     "$ref": "ImportSshPublicKeyResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/compute"
	//   ]
	// }

}

// method id "oslogin.users.sshPublicKeys.delete":

type UsersSshPublicKeysDeleteCall struct {
	s          *Service
	name       string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Delete: Deletes an SSH public key.
func (r *UsersSshPublicKeysService) Delete(name string) *UsersSshPublicKeysDeleteCall {
	c := &UsersSshPublicKeysDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UsersSshPublicKeysDeleteCall) Fields(s ...googleapi.Field) *UsersSshPublicKeysDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *UsersSshPublicKeysDeleteCall) Context(ctx context.Context) *UsersSshPublicKeysDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *UsersSshPublicKeysDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *UsersSshPublicKeysDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "oslogin.users.sshPublicKeys.delete" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *UsersSshPublicKeysDeleteCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
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
	ret := &Empty{
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
	//   "description": "Deletes an SSH public key.",
	//   "flatPath": "v1beta/users/{usersId}/sshPublicKeys/{sshPublicKeysId}",
	//   "httpMethod": "DELETE",
	//   "id": "oslogin.users.sshPublicKeys.delete",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The fingerprint of the public key to update. Public keys are identified by\ntheir SHA-256 fingerprint. The fingerprint of the public key is in format\n`users/{user}/sshPublicKeys/{fingerprint}`.",
	//       "location": "path",
	//       "pattern": "^users/[^/]+/sshPublicKeys/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/{+name}",
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/compute"
	//   ]
	// }

}

// method id "oslogin.users.sshPublicKeys.get":

type UsersSshPublicKeysGetCall struct {
	s            *Service
	name         string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Retrieves an SSH public key.
func (r *UsersSshPublicKeysService) Get(name string) *UsersSshPublicKeysGetCall {
	c := &UsersSshPublicKeysGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UsersSshPublicKeysGetCall) Fields(s ...googleapi.Field) *UsersSshPublicKeysGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *UsersSshPublicKeysGetCall) IfNoneMatch(entityTag string) *UsersSshPublicKeysGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *UsersSshPublicKeysGetCall) Context(ctx context.Context) *UsersSshPublicKeysGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *UsersSshPublicKeysGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *UsersSshPublicKeysGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "oslogin.users.sshPublicKeys.get" call.
// Exactly one of *SshPublicKey or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *SshPublicKey.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *UsersSshPublicKeysGetCall) Do(opts ...googleapi.CallOption) (*SshPublicKey, error) {
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
	ret := &SshPublicKey{
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
	//   "description": "Retrieves an SSH public key.",
	//   "flatPath": "v1beta/users/{usersId}/sshPublicKeys/{sshPublicKeysId}",
	//   "httpMethod": "GET",
	//   "id": "oslogin.users.sshPublicKeys.get",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The fingerprint of the public key to retrieve. Public keys are identified\nby their SHA-256 fingerprint. The fingerprint of the public key is in\nformat `users/{user}/sshPublicKeys/{fingerprint}`.",
	//       "location": "path",
	//       "pattern": "^users/[^/]+/sshPublicKeys/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/{+name}",
	//   "response": {
	//     "$ref": "SshPublicKey"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/compute"
	//   ]
	// }

}

// method id "oslogin.users.sshPublicKeys.patch":

type UsersSshPublicKeysPatchCall struct {
	s            *Service
	name         string
	sshpublickey *SshPublicKey
	urlParams_   gensupport.URLParams
	ctx_         context.Context
	header_      http.Header
}

// Patch: Updates an SSH public key and returns the profile information.
// This method
// supports patch semantics.
func (r *UsersSshPublicKeysService) Patch(name string, sshpublickey *SshPublicKey) *UsersSshPublicKeysPatchCall {
	c := &UsersSshPublicKeysPatchCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.sshpublickey = sshpublickey
	return c
}

// UpdateMask sets the optional parameter "updateMask": Mask to control
// which fields get updated. Updates all if not present.
func (c *UsersSshPublicKeysPatchCall) UpdateMask(updateMask string) *UsersSshPublicKeysPatchCall {
	c.urlParams_.Set("updateMask", updateMask)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *UsersSshPublicKeysPatchCall) Fields(s ...googleapi.Field) *UsersSshPublicKeysPatchCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *UsersSshPublicKeysPatchCall) Context(ctx context.Context) *UsersSshPublicKeysPatchCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *UsersSshPublicKeysPatchCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *UsersSshPublicKeysPatchCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.sshpublickey)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "oslogin.users.sshPublicKeys.patch" call.
// Exactly one of *SshPublicKey or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *SshPublicKey.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *UsersSshPublicKeysPatchCall) Do(opts ...googleapi.CallOption) (*SshPublicKey, error) {
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
	ret := &SshPublicKey{
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
	//   "description": "Updates an SSH public key and returns the profile information. This method\nsupports patch semantics.",
	//   "flatPath": "v1beta/users/{usersId}/sshPublicKeys/{sshPublicKeysId}",
	//   "httpMethod": "PATCH",
	//   "id": "oslogin.users.sshPublicKeys.patch",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The fingerprint of the public key to update. Public keys are identified by\ntheir SHA-256 fingerprint. The fingerprint of the public key is in format\n`users/{user}/sshPublicKeys/{fingerprint}`.",
	//       "location": "path",
	//       "pattern": "^users/[^/]+/sshPublicKeys/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "updateMask": {
	//       "description": "Mask to control which fields get updated. Updates all if not present.",
	//       "format": "google-fieldmask",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/{+name}",
	//   "request": {
	//     "$ref": "SshPublicKey"
	//   },
	//   "response": {
	//     "$ref": "SshPublicKey"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/compute"
	//   ]
	// }

}
