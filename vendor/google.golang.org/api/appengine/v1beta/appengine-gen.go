// Package appengine provides access to the Google App Engine Admin API.
//
// See https://cloud.google.com/appengine/docs/admin-api/
//
// Usage example:
//
//   import "google.golang.org/api/appengine/v1beta"
//   ...
//   appengineService, err := appengine.New(oauthHttpClient)
package appengine // import "google.golang.org/api/appengine/v1beta"

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

const apiId = "appengine:v1beta"
const apiName = "appengine"
const apiVersion = "v1beta"
const basePath = "https://appengine.googleapis.com/"

// OAuth2 scopes used by this API.
const (
	// View and manage your applications deployed on Google App Engine
	AppengineAdminScope = "https://www.googleapis.com/auth/appengine.admin"

	// View and manage your data across Google Cloud Platform services
	CloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"

	// View your data across Google Cloud Platform services
	CloudPlatformReadOnlyScope = "https://www.googleapis.com/auth/cloud-platform.read-only"
)

func New(client *http.Client) (*APIService, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	s := &APIService{client: client, BasePath: basePath}
	s.Apps = NewAppsService(s)
	return s, nil
}

type APIService struct {
	client    *http.Client
	BasePath  string // API endpoint base URL
	UserAgent string // optional additional User-Agent fragment

	Apps *AppsService
}

func (s *APIService) userAgent() string {
	if s.UserAgent == "" {
		return googleapi.UserAgent
	}
	return googleapi.UserAgent + " " + s.UserAgent
}

func NewAppsService(s *APIService) *AppsService {
	rs := &AppsService{s: s}
	rs.AuthorizedCertificates = NewAppsAuthorizedCertificatesService(s)
	rs.AuthorizedDomains = NewAppsAuthorizedDomainsService(s)
	rs.DomainMappings = NewAppsDomainMappingsService(s)
	rs.Firewall = NewAppsFirewallService(s)
	rs.Locations = NewAppsLocationsService(s)
	rs.Operations = NewAppsOperationsService(s)
	rs.Services = NewAppsServicesService(s)
	return rs
}

type AppsService struct {
	s *APIService

	AuthorizedCertificates *AppsAuthorizedCertificatesService

	AuthorizedDomains *AppsAuthorizedDomainsService

	DomainMappings *AppsDomainMappingsService

	Firewall *AppsFirewallService

	Locations *AppsLocationsService

	Operations *AppsOperationsService

	Services *AppsServicesService
}

func NewAppsAuthorizedCertificatesService(s *APIService) *AppsAuthorizedCertificatesService {
	rs := &AppsAuthorizedCertificatesService{s: s}
	return rs
}

type AppsAuthorizedCertificatesService struct {
	s *APIService
}

func NewAppsAuthorizedDomainsService(s *APIService) *AppsAuthorizedDomainsService {
	rs := &AppsAuthorizedDomainsService{s: s}
	return rs
}

type AppsAuthorizedDomainsService struct {
	s *APIService
}

func NewAppsDomainMappingsService(s *APIService) *AppsDomainMappingsService {
	rs := &AppsDomainMappingsService{s: s}
	return rs
}

type AppsDomainMappingsService struct {
	s *APIService
}

func NewAppsFirewallService(s *APIService) *AppsFirewallService {
	rs := &AppsFirewallService{s: s}
	rs.IngressRules = NewAppsFirewallIngressRulesService(s)
	return rs
}

type AppsFirewallService struct {
	s *APIService

	IngressRules *AppsFirewallIngressRulesService
}

func NewAppsFirewallIngressRulesService(s *APIService) *AppsFirewallIngressRulesService {
	rs := &AppsFirewallIngressRulesService{s: s}
	return rs
}

type AppsFirewallIngressRulesService struct {
	s *APIService
}

func NewAppsLocationsService(s *APIService) *AppsLocationsService {
	rs := &AppsLocationsService{s: s}
	return rs
}

type AppsLocationsService struct {
	s *APIService
}

func NewAppsOperationsService(s *APIService) *AppsOperationsService {
	rs := &AppsOperationsService{s: s}
	return rs
}

type AppsOperationsService struct {
	s *APIService
}

func NewAppsServicesService(s *APIService) *AppsServicesService {
	rs := &AppsServicesService{s: s}
	rs.Versions = NewAppsServicesVersionsService(s)
	return rs
}

type AppsServicesService struct {
	s *APIService

	Versions *AppsServicesVersionsService
}

func NewAppsServicesVersionsService(s *APIService) *AppsServicesVersionsService {
	rs := &AppsServicesVersionsService{s: s}
	rs.Instances = NewAppsServicesVersionsInstancesService(s)
	return rs
}

type AppsServicesVersionsService struct {
	s *APIService

	Instances *AppsServicesVersionsInstancesService
}

func NewAppsServicesVersionsInstancesService(s *APIService) *AppsServicesVersionsInstancesService {
	rs := &AppsServicesVersionsInstancesService{s: s}
	return rs
}

type AppsServicesVersionsInstancesService struct {
	s *APIService
}

// ApiConfigHandler: Google Cloud Endpoints
// (https://cloud.google.com/appengine/docs/python/endpoints/)
// configuration for API handlers.
type ApiConfigHandler struct {
	// AuthFailAction: Action to take when users access resources that
	// require authentication. Defaults to redirect.
	//
	// Possible values:
	//   "AUTH_FAIL_ACTION_UNSPECIFIED" - Not specified.
	// AUTH_FAIL_ACTION_REDIRECT is assumed.
	//   "AUTH_FAIL_ACTION_REDIRECT" - Redirects user to
	// "accounts.google.com". The user is redirected back to the application
	// URL after signing in or creating an account.
	//   "AUTH_FAIL_ACTION_UNAUTHORIZED" - Rejects request with a 401 HTTP
	// status code and an error message.
	AuthFailAction string `json:"authFailAction,omitempty"`

	// Login: Level of login required to access this resource. Defaults to
	// optional.
	//
	// Possible values:
	//   "LOGIN_UNSPECIFIED" - Not specified. LOGIN_OPTIONAL is assumed.
	//   "LOGIN_OPTIONAL" - Does not require that the user is signed in.
	//   "LOGIN_ADMIN" - If the user is not signed in, the auth_fail_action
	// is taken. In addition, if the user is not an administrator for the
	// application, they are given an error message regardless of
	// auth_fail_action. If the user is an administrator, the handler
	// proceeds.
	//   "LOGIN_REQUIRED" - If the user has signed in, the handler proceeds
	// normally. Otherwise, the auth_fail_action is taken.
	Login string `json:"login,omitempty"`

	// Script: Path to the script from the application root directory.
	Script string `json:"script,omitempty"`

	// SecurityLevel: Security (HTTPS) enforcement for this URL.
	//
	// Possible values:
	//   "SECURE_UNSPECIFIED" - Not specified.
	//   "SECURE_DEFAULT" - Both HTTP and HTTPS requests with URLs that
	// match the handler succeed without redirects. The application can
	// examine the request to determine which protocol was used, and respond
	// accordingly.
	//   "SECURE_NEVER" - Requests for a URL that match this handler that
	// use HTTPS are automatically redirected to the HTTP equivalent URL.
	//   "SECURE_OPTIONAL" - Both HTTP and HTTPS requests with URLs that
	// match the handler succeed without redirects. The application can
	// examine the request to determine which protocol was used and respond
	// accordingly.
	//   "SECURE_ALWAYS" - Requests for a URL that match this handler that
	// do not use HTTPS are automatically redirected to the HTTPS URL with
	// the same path. Query parameters are reserved for the redirect.
	SecurityLevel string `json:"securityLevel,omitempty"`

	// Url: URL to serve the endpoint at.
	Url string `json:"url,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AuthFailAction") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AuthFailAction") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *ApiConfigHandler) MarshalJSON() ([]byte, error) {
	type noMethod ApiConfigHandler
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ApiEndpointHandler: Uses Google Cloud Endpoints to handle requests.
type ApiEndpointHandler struct {
	// ScriptPath: Path to the script from the application root directory.
	ScriptPath string `json:"scriptPath,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ScriptPath") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ScriptPath") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ApiEndpointHandler) MarshalJSON() ([]byte, error) {
	type noMethod ApiEndpointHandler
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Application: An Application resource contains the top-level
// configuration of an App Engine application. Next tag: 20
type Application struct {
	// AuthDomain: Google Apps authentication domain that controls which
	// users can access this application.Defaults to open access for any
	// Google Account.
	AuthDomain string `json:"authDomain,omitempty"`

	// CodeBucket: Google Cloud Storage bucket that can be used for storing
	// files associated with this application. This bucket is associated
	// with the application and can be used by the gcloud deployment
	// commands.@OutputOnly
	CodeBucket string `json:"codeBucket,omitempty"`

	// DefaultBucket: Google Cloud Storage bucket that can be used by this
	// application to store content.@OutputOnly
	DefaultBucket string `json:"defaultBucket,omitempty"`

	// DefaultCookieExpiration: Cookie expiration policy for this
	// application.
	DefaultCookieExpiration string `json:"defaultCookieExpiration,omitempty"`

	// DefaultHostname: Hostname used to reach this application, as resolved
	// by App Engine.@OutputOnly
	DefaultHostname string `json:"defaultHostname,omitempty"`

	// DispatchRules: HTTP path dispatch rules for requests to the
	// application that do not explicitly target a service or version. Rules
	// are order-dependent. Up to 20 dispatch rules can be
	// supported.@OutputOnly
	DispatchRules []*UrlDispatchRule `json:"dispatchRules,omitempty"`

	// FeatureSettings: The feature specific settings to be used in the
	// application.
	FeatureSettings *FeatureSettings `json:"featureSettings,omitempty"`

	// GcrDomain: The Google Container Registry domain used for storing
	// managed build docker images for this application.
	GcrDomain string `json:"gcrDomain,omitempty"`

	Iap *IdentityAwareProxy `json:"iap,omitempty"`

	// Id: Identifier of the Application resource. This identifier is
	// equivalent to the project ID of the Google Cloud Platform project
	// where you want to deploy your application. Example: myapp.
	Id string `json:"id,omitempty"`

	// LocationId: Location from which this application will be run.
	// Application instances will run out of data centers in the chosen
	// location, which is also where all of the application's end user
	// content is stored.Defaults to us-central.Options are:us-central -
	// Central USeurope-west - Western Europeus-east1 - Eastern US
	LocationId string `json:"locationId,omitempty"`

	// Name: Full path to the Application resource in the API. Example:
	// apps/myapp.@OutputOnly
	Name string `json:"name,omitempty"`

	// ServingStatus: Serving status of this application.
	//
	// Possible values:
	//   "UNSPECIFIED" - Serving status is unspecified.
	//   "SERVING" - Application is serving.
	//   "USER_DISABLED" - Application has been disabled by the user.
	//   "SYSTEM_DISABLED" - Application has been disabled by the system.
	ServingStatus string `json:"servingStatus,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "AuthDomain") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AuthDomain") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Application) MarshalJSON() ([]byte, error) {
	type noMethod Application
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// AuthorizedCertificate: An SSL certificate that a user has been
// authorized to administer. A user is authorized to administer any
// certificate that applies to one of their authorized domains.
type AuthorizedCertificate struct {
	// CertificateRawData: The SSL certificate serving the
	// AuthorizedCertificate resource. This must be obtained independently
	// from a certificate authority.
	CertificateRawData *CertificateRawData `json:"certificateRawData,omitempty"`

	// DisplayName: The user-specified display name of the certificate. This
	// is not guaranteed to be unique. Example: My Certificate.
	DisplayName string `json:"displayName,omitempty"`

	// DomainMappingsCount: Aggregate count of the domain mappings with this
	// certificate mapped. This count includes domain mappings on
	// applications for which the user does not have VIEWER permissions.Only
	// returned by GET or LIST requests when specifically requested by the
	// view=FULL_CERTIFICATE option.@OutputOnly
	DomainMappingsCount int64 `json:"domainMappingsCount,omitempty"`

	// DomainNames: Topmost applicable domains of this certificate. This
	// certificate applies to these domains and their subdomains. Example:
	// example.com.@OutputOnly
	DomainNames []string `json:"domainNames,omitempty"`

	// ExpireTime: The time when this certificate expires. To update the
	// renewal time on this certificate, upload an SSL certificate with a
	// different expiration time using
	// AuthorizedCertificates.UpdateAuthorizedCertificate.@OutputOnly
	ExpireTime string `json:"expireTime,omitempty"`

	// Id: Relative name of the certificate. This is a unique value
	// autogenerated on AuthorizedCertificate resource creation. Example:
	// 12345.@OutputOnly
	Id string `json:"id,omitempty"`

	// ManagedCertificate: Only applicable if this certificate is managed by
	// App Engine. Managed certificates are tied to the lifecycle of a
	// DomainMapping and cannot be updated or deleted via the
	// AuthorizedCertificates API. If this certificate is manually
	// administered by the user, this field will be empty.@OutputOnly
	ManagedCertificate *ManagedCertificate `json:"managedCertificate,omitempty"`

	// Name: Full path to the AuthorizedCertificate resource in the API.
	// Example: apps/myapp/authorizedCertificates/12345.@OutputOnly
	Name string `json:"name,omitempty"`

	// VisibleDomainMappings: The full paths to user visible Domain Mapping
	// resources that have this certificate mapped. Example:
	// apps/myapp/domainMappings/example.com.This may not represent the full
	// list of mapped domain mappings if the user does not have VIEWER
	// permissions on all of the applications that have this certificate
	// mapped. See domain_mappings_count for a complete count.Only returned
	// by GET or LIST requests when specifically requested by the
	// view=FULL_CERTIFICATE option.@OutputOnly
	VisibleDomainMappings []string `json:"visibleDomainMappings,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "CertificateRawData")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CertificateRawData") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *AuthorizedCertificate) MarshalJSON() ([]byte, error) {
	type noMethod AuthorizedCertificate
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// AuthorizedDomain: A domain that a user has been authorized to
// administer. To authorize use of a domain, verify ownership via
// Webmaster Central
// (https://www.google.com/webmasters/verification/home).
type AuthorizedDomain struct {
	// Id: Fully qualified domain name of the domain authorized for use.
	// Example: example.com.
	Id string `json:"id,omitempty"`

	// Name: Full path to the AuthorizedDomain resource in the API. Example:
	// apps/myapp/authorizedDomains/example.com.@OutputOnly
	Name string `json:"name,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Id") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Id") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *AuthorizedDomain) MarshalJSON() ([]byte, error) {
	type noMethod AuthorizedDomain
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// AutomaticScaling: Automatic scaling is based on request rate,
// response latencies, and other application metrics.
type AutomaticScaling struct {
	// CoolDownPeriod: Amount of time that the Autoscaler
	// (https://cloud.google.com/compute/docs/autoscaler/) should wait
	// between changes to the number of virtual machines. Only applicable
	// for VM runtimes.
	CoolDownPeriod string `json:"coolDownPeriod,omitempty"`

	// CpuUtilization: Target scaling by CPU usage.
	CpuUtilization *CpuUtilization `json:"cpuUtilization,omitempty"`

	// DiskUtilization: Target scaling by disk usage.
	DiskUtilization *DiskUtilization `json:"diskUtilization,omitempty"`

	// MaxConcurrentRequests: Number of concurrent requests an automatic
	// scaling instance can accept before the scheduler spawns a new
	// instance.Defaults to a runtime-specific value.
	MaxConcurrentRequests int64 `json:"maxConcurrentRequests,omitempty"`

	// MaxIdleInstances: Maximum number of idle instances that should be
	// maintained for this version.
	MaxIdleInstances int64 `json:"maxIdleInstances,omitempty"`

	// MaxPendingLatency: Maximum amount of time that a request should wait
	// in the pending queue before starting a new instance to handle it.
	MaxPendingLatency string `json:"maxPendingLatency,omitempty"`

	// MaxTotalInstances: Maximum number of instances that should be started
	// to handle requests.
	MaxTotalInstances int64 `json:"maxTotalInstances,omitempty"`

	// MinIdleInstances: Minimum number of idle instances that should be
	// maintained for this version. Only applicable for the default version
	// of a service.
	MinIdleInstances int64 `json:"minIdleInstances,omitempty"`

	// MinPendingLatency: Minimum amount of time a request should wait in
	// the pending queue before starting a new instance to handle it.
	MinPendingLatency string `json:"minPendingLatency,omitempty"`

	// MinTotalInstances: Minimum number of instances that should be
	// maintained for this version.
	MinTotalInstances int64 `json:"minTotalInstances,omitempty"`

	// NetworkUtilization: Target scaling by network usage.
	NetworkUtilization *NetworkUtilization `json:"networkUtilization,omitempty"`

	// RequestUtilization: Target scaling by request utilization.
	RequestUtilization *RequestUtilization `json:"requestUtilization,omitempty"`

	// StandardSchedulerSettings: Scheduler settings for standard
	// environment.
	StandardSchedulerSettings *StandardSchedulerSettings `json:"standardSchedulerSettings,omitempty"`

	// ForceSendFields is a list of field names (e.g. "CoolDownPeriod") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CoolDownPeriod") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *AutomaticScaling) MarshalJSON() ([]byte, error) {
	type noMethod AutomaticScaling
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// BasicScaling: A service with basic scaling will create an instance
// when the application receives a request. The instance will be turned
// down when the app becomes idle. Basic scaling is ideal for work that
// is intermittent or driven by user activity.
type BasicScaling struct {
	// IdleTimeout: Duration of time after the last request that an instance
	// must wait before the instance is shut down.
	IdleTimeout string `json:"idleTimeout,omitempty"`

	// MaxInstances: Maximum number of instances to create for this version.
	MaxInstances int64 `json:"maxInstances,omitempty"`

	// ForceSendFields is a list of field names (e.g. "IdleTimeout") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "IdleTimeout") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *BasicScaling) MarshalJSON() ([]byte, error) {
	type noMethod BasicScaling
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// BatchUpdateIngressRulesRequest: Request message for
// Firewall.BatchUpdateIngressRules.
type BatchUpdateIngressRulesRequest struct {
	// IngressRules: A list of FirewallRules to replace the existing set.
	IngressRules []*FirewallRule `json:"ingressRules,omitempty"`

	// ForceSendFields is a list of field names (e.g. "IngressRules") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "IngressRules") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *BatchUpdateIngressRulesRequest) MarshalJSON() ([]byte, error) {
	type noMethod BatchUpdateIngressRulesRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// BatchUpdateIngressRulesResponse: Response message for
// Firewall.UpdateAllIngressRules.
type BatchUpdateIngressRulesResponse struct {
	// IngressRules: The full list of ingress FirewallRules for this
	// application.
	IngressRules []*FirewallRule `json:"ingressRules,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "IngressRules") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "IngressRules") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *BatchUpdateIngressRulesResponse) MarshalJSON() ([]byte, error) {
	type noMethod BatchUpdateIngressRulesResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// BuildInfo: Google Cloud Container Builder build information.
type BuildInfo struct {
	// CloudBuildId: The Google Cloud Container Builder build id. Example:
	// "f966068f-08b2-42c8-bdfe-74137dff2bf9"
	CloudBuildId string `json:"cloudBuildId,omitempty"`

	// ForceSendFields is a list of field names (e.g. "CloudBuildId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CloudBuildId") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *BuildInfo) MarshalJSON() ([]byte, error) {
	type noMethod BuildInfo
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CertificateRawData: An SSL certificate obtained from a certificate
// authority.
type CertificateRawData struct {
	// PrivateKey: Unencrypted PEM encoded RSA private key. This field is
	// set once on certificate creation and then encrypted. The key size
	// must be 2048 bits or fewer. Must include the header and footer.
	// Example: <pre> -----BEGIN RSA PRIVATE KEY-----
	// <unencrypted_key_value> -----END RSA PRIVATE KEY----- </pre>
	// @InputOnly
	PrivateKey string `json:"privateKey,omitempty"`

	// PublicCertificate: PEM encoded x.509 public key certificate. This
	// field is set once on certificate creation. Must include the header
	// and footer. Example: <pre> -----BEGIN CERTIFICATE-----
	// <certificate_value> -----END CERTIFICATE----- </pre>
	PublicCertificate string `json:"publicCertificate,omitempty"`

	// ForceSendFields is a list of field names (e.g. "PrivateKey") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "PrivateKey") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *CertificateRawData) MarshalJSON() ([]byte, error) {
	type noMethod CertificateRawData
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ContainerInfo: Docker image that is used to create a container and
// start a VM instance for the version that you deploy. Only applicable
// for instances running in the App Engine flexible environment.
type ContainerInfo struct {
	// Image: URI to the hosted container image in Google Container
	// Registry. The URI must be fully qualified and include a tag or
	// digest. Examples: "gcr.io/my-project/image:tag" or
	// "gcr.io/my-project/image@digest"
	Image string `json:"image,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Image") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Image") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ContainerInfo) MarshalJSON() ([]byte, error) {
	type noMethod ContainerInfo
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CpuUtilization: Target scaling by CPU usage.
type CpuUtilization struct {
	// AggregationWindowLength: Period of time over which CPU utilization is
	// calculated.
	AggregationWindowLength string `json:"aggregationWindowLength,omitempty"`

	// TargetUtilization: Target CPU utilization ratio to maintain when
	// scaling. Must be between 0 and 1.
	TargetUtilization float64 `json:"targetUtilization,omitempty"`

	// ForceSendFields is a list of field names (e.g.
	// "AggregationWindowLength") to unconditionally include in API
	// requests. By default, fields with empty values are omitted from API
	// requests. However, any non-pointer, non-interface field appearing in
	// ForceSendFields will be sent to the server regardless of whether the
	// field is empty or not. This may be used to include empty fields in
	// Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AggregationWindowLength")
	// to include in API requests with the JSON null value. By default,
	// fields with empty values are omitted from API requests. However, any
	// field with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *CpuUtilization) MarshalJSON() ([]byte, error) {
	type noMethod CpuUtilization
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

func (s *CpuUtilization) UnmarshalJSON(data []byte) error {
	type noMethod CpuUtilization
	var s1 struct {
		TargetUtilization gensupport.JSONFloat64 `json:"targetUtilization"`
		*noMethod
	}
	s1.noMethod = (*noMethod)(s)
	if err := json.Unmarshal(data, &s1); err != nil {
		return err
	}
	s.TargetUtilization = float64(s1.TargetUtilization)
	return nil
}

// DebugInstanceRequest: Request message for Instances.DebugInstance.
type DebugInstanceRequest struct {
	// SshKey: Public SSH key to add to the instance.
	// Examples:
	// [USERNAME]:ssh-rsa [KEY_VALUE] [USERNAME]
	// [USERNAME]:ssh-rsa [KEY_VALUE] google-ssh
	// {"userName":"[USERNAME]","expireOn":"[EXPIRE_TIME]"}For more
	// information, see Adding and Removing SSH Keys
	// (https://cloud.google.com/compute/docs/instances/adding-removing-ssh-k
	// eys).
	SshKey string `json:"sshKey,omitempty"`

	// ForceSendFields is a list of field names (e.g. "SshKey") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "SshKey") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *DebugInstanceRequest) MarshalJSON() ([]byte, error) {
	type noMethod DebugInstanceRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Deployment: Code and application artifacts used to deploy a version
// to App Engine.
type Deployment struct {
	// Build: Google Cloud Container Builder build information.
	Build *BuildInfo `json:"build,omitempty"`

	// Container: The Docker image for the container that runs the version.
	// Only applicable for instances running in the App Engine flexible
	// environment.
	Container *ContainerInfo `json:"container,omitempty"`

	// Files: Manifest of the files stored in Google Cloud Storage that are
	// included as part of this version. All files must be readable using
	// the credentials supplied with this call.
	Files map[string]FileInfo `json:"files,omitempty"`

	// Zip: The zip file for this deployment, if this is a zip deployment.
	Zip *ZipInfo `json:"zip,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Build") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Build") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Deployment) MarshalJSON() ([]byte, error) {
	type noMethod Deployment
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// DiskUtilization: Target scaling by disk usage. Only applicable for VM
// runtimes.
type DiskUtilization struct {
	// TargetReadBytesPerSecond: Target bytes read per second.
	TargetReadBytesPerSecond int64 `json:"targetReadBytesPerSecond,omitempty"`

	// TargetReadOpsPerSecond: Target ops read per seconds.
	TargetReadOpsPerSecond int64 `json:"targetReadOpsPerSecond,omitempty"`

	// TargetWriteBytesPerSecond: Target bytes written per second.
	TargetWriteBytesPerSecond int64 `json:"targetWriteBytesPerSecond,omitempty"`

	// TargetWriteOpsPerSecond: Target ops written per second.
	TargetWriteOpsPerSecond int64 `json:"targetWriteOpsPerSecond,omitempty"`

	// ForceSendFields is a list of field names (e.g.
	// "TargetReadBytesPerSecond") to unconditionally include in API
	// requests. By default, fields with empty values are omitted from API
	// requests. However, any non-pointer, non-interface field appearing in
	// ForceSendFields will be sent to the server regardless of whether the
	// field is empty or not. This may be used to include empty fields in
	// Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "TargetReadBytesPerSecond")
	// to include in API requests with the JSON null value. By default,
	// fields with empty values are omitted from API requests. However, any
	// field with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *DiskUtilization) MarshalJSON() ([]byte, error) {
	type noMethod DiskUtilization
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// DomainMapping: A domain serving an App Engine application.
type DomainMapping struct {
	// Id: Relative name of the domain serving the application. Example:
	// example.com.
	Id string `json:"id,omitempty"`

	// Name: Full path to the DomainMapping resource in the API. Example:
	// apps/myapp/domainMapping/example.com.@OutputOnly
	Name string `json:"name,omitempty"`

	// ResourceRecords: The resource records required to configure this
	// domain mapping. These records must be added to the domain's DNS
	// configuration in order to serve the application via this domain
	// mapping.@OutputOnly
	ResourceRecords []*ResourceRecord `json:"resourceRecords,omitempty"`

	// SslSettings: SSL configuration for this domain. If unconfigured, this
	// domain will not serve with SSL.
	SslSettings *SslSettings `json:"sslSettings,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Id") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Id") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *DomainMapping) MarshalJSON() ([]byte, error) {
	type noMethod DomainMapping
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Empty: A generic empty message that you can re-use to avoid defining
// duplicated empty messages in your APIs. A typical example is to use
// it as the request or the response type of an API method. For
// instance:
// service Foo {
//   rpc Bar(google.protobuf.Empty) returns
// (google.protobuf.Empty);
// }
// The JSON representation for Empty is empty JSON object {}.
type Empty struct {
	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`
}

// EndpointsApiService: Cloud Endpoints
// (https://cloud.google.com/endpoints) configuration. The Endpoints API
// Service provides tooling for serving Open API and gRPC endpoints via
// an NGINX proxy.The fields here refer to the name and configuration id
// of a "service" resource in the Service Management API
// (https://cloud.google.com/service-management/overview).
type EndpointsApiService struct {
	// ConfigId: Endpoints service configuration id as specified by the
	// Service Management API. For example "2016-09-19r1"
	ConfigId string `json:"configId,omitempty"`

	// Name: Endpoints service name which is the name of the "service"
	// resource in the Service Management API. For example
	// "myapi.endpoints.myproject.cloud.goog"
	Name string `json:"name,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ConfigId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ConfigId") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *EndpointsApiService) MarshalJSON() ([]byte, error) {
	type noMethod EndpointsApiService
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ErrorHandler: Custom static error page to be served when an error
// occurs.
type ErrorHandler struct {
	// ErrorCode: Error condition this handler applies to.
	//
	// Possible values:
	//   "ERROR_CODE_UNSPECIFIED" - Not specified. ERROR_CODE_DEFAULT is
	// assumed.
	//   "ERROR_CODE_DEFAULT" - All other error types.
	//   "ERROR_CODE_OVER_QUOTA" - Application has exceeded a resource
	// quota.
	//   "ERROR_CODE_DOS_API_DENIAL" - Client blocked by the application's
	// Denial of Service protection configuration.
	//   "ERROR_CODE_TIMEOUT" - Deadline reached before the application
	// responds.
	ErrorCode string `json:"errorCode,omitempty"`

	// MimeType: MIME type of file. Defaults to text/html.
	MimeType string `json:"mimeType,omitempty"`

	// StaticFile: Static file content to be served for this error.
	StaticFile string `json:"staticFile,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ErrorCode") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ErrorCode") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ErrorHandler) MarshalJSON() ([]byte, error) {
	type noMethod ErrorHandler
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// FeatureSettings: The feature specific settings to be used in the
// application. These define behaviors that are user configurable.
type FeatureSettings struct {
	// SplitHealthChecks: Boolean value indicating if split health checks
	// should be used instead of the legacy health checks. At an app.yaml
	// level, this means defaulting to 'readiness_check' and
	// 'liveness_check' values instead of 'health_check' ones. Once the
	// legacy 'health_check' behavior is deprecated, and this value is
	// always true, this setting can be removed.
	SplitHealthChecks bool `json:"splitHealthChecks,omitempty"`

	// ForceSendFields is a list of field names (e.g. "SplitHealthChecks")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "SplitHealthChecks") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *FeatureSettings) MarshalJSON() ([]byte, error) {
	type noMethod FeatureSettings
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// FileInfo: Single source file that is part of the version to be
// deployed. Each source file that is deployed must be specified
// separately.
type FileInfo struct {
	// MimeType: The MIME type of the file.Defaults to the value from Google
	// Cloud Storage.
	MimeType string `json:"mimeType,omitempty"`

	// Sha1Sum: The SHA1 hash of the file, in hex.
	Sha1Sum string `json:"sha1Sum,omitempty"`

	// SourceUrl: URL source to use to fetch this file. Must be a URL to a
	// resource in Google Cloud Storage in the form
	// 'http(s)://storage.googleapis.com/<bucket>/<object>'.
	SourceUrl string `json:"sourceUrl,omitempty"`

	// ForceSendFields is a list of field names (e.g. "MimeType") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "MimeType") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *FileInfo) MarshalJSON() ([]byte, error) {
	type noMethod FileInfo
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// FirewallRule: A single firewall rule that is evaluated against
// incoming traffic and provides an action to take on matched requests.
type FirewallRule struct {
	// Action: The action to take on matched requests.
	//
	// Possible values:
	//   "UNSPECIFIED_ACTION"
	//   "ALLOW" - Matching requests are allowed.
	//   "DENY" - Matching requests are denied.
	Action string `json:"action,omitempty"`

	// Description: An optional string description of this rule. This field
	// has a maximum length of 100 characters.
	Description string `json:"description,omitempty"`

	// Priority: A positive integer between 1, Int32.MaxValue-1 that defines
	// the order of rule evaluation. Rules with the lowest priority are
	// evaluated first.A default rule at priority Int32.MaxValue matches all
	// IPv4 and IPv6 traffic when no previous rule matches. Only the action
	// of this rule can be modified by the user.
	Priority int64 `json:"priority,omitempty"`

	// SourceRange: IP address or range, defined using CIDR notation, of
	// requests that this rule applies to. You can use the wildcard
	// character "*" to match all IPs equivalent to "0/0" and "::/0"
	// together. Examples: 192.168.1.1 or 192.168.0.0/16 or 2001:db8::/32
	// or 2001:0db8:0000:0042:0000:8a2e:0370:7334.<p>Truncation will be
	// silently performed on addresses which are not properly truncated. For
	// example, 1.2.3.4/24 is accepted as the same address as 1.2.3.0/24.
	// Similarly, for IPv6, 2001:db8::1/32 is accepted as the same address
	// as 2001:db8::/32.
	SourceRange string `json:"sourceRange,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Action") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Action") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *FirewallRule) MarshalJSON() ([]byte, error) {
	type noMethod FirewallRule
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// HealthCheck: Health checking configuration for VM instances.
// Unhealthy instances are killed and replaced with new instances. Only
// applicable for instances in App Engine flexible environment.
type HealthCheck struct {
	// CheckInterval: Interval between health checks.
	CheckInterval string `json:"checkInterval,omitempty"`

	// DisableHealthCheck: Whether to explicitly disable health checks for
	// this instance.
	DisableHealthCheck bool `json:"disableHealthCheck,omitempty"`

	// HealthyThreshold: Number of consecutive successful health checks
	// required before receiving traffic.
	HealthyThreshold int64 `json:"healthyThreshold,omitempty"`

	// Host: Host header to send when performing an HTTP health check.
	// Example: "myapp.appspot.com"
	Host string `json:"host,omitempty"`

	// RestartThreshold: Number of consecutive failed health checks required
	// before an instance is restarted.
	RestartThreshold int64 `json:"restartThreshold,omitempty"`

	// Timeout: Time before the health check is considered failed.
	Timeout string `json:"timeout,omitempty"`

	// UnhealthyThreshold: Number of consecutive failed health checks
	// required before removing traffic.
	UnhealthyThreshold int64 `json:"unhealthyThreshold,omitempty"`

	// ForceSendFields is a list of field names (e.g. "CheckInterval") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CheckInterval") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *HealthCheck) MarshalJSON() ([]byte, error) {
	type noMethod HealthCheck
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// IdentityAwareProxy: Identity-Aware Proxy
type IdentityAwareProxy struct {
	// Enabled: Whether the serving infrastructure will authenticate and
	// authorize all incoming requests.If true, the oauth2_client_id and
	// oauth2_client_secret fields must be non-empty.
	Enabled bool `json:"enabled,omitempty"`

	// Oauth2ClientId: OAuth2 client ID to use for the authentication flow.
	Oauth2ClientId string `json:"oauth2ClientId,omitempty"`

	// Oauth2ClientSecret: OAuth2 client secret to use for the
	// authentication flow.For security reasons, this value cannot be
	// retrieved via the API. Instead, the SHA-256 hash of the value is
	// returned in the oauth2_client_secret_sha256 field.@InputOnly
	Oauth2ClientSecret string `json:"oauth2ClientSecret,omitempty"`

	// Oauth2ClientSecretSha256: Hex-encoded SHA-256 hash of the client
	// secret.@OutputOnly
	Oauth2ClientSecretSha256 string `json:"oauth2ClientSecretSha256,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Enabled") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Enabled") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *IdentityAwareProxy) MarshalJSON() ([]byte, error) {
	type noMethod IdentityAwareProxy
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Instance: An Instance resource is the computing unit that App Engine
// uses to automatically scale an application.
type Instance struct {
	// AppEngineRelease: App Engine release this instance is running
	// on.@OutputOnly
	AppEngineRelease string `json:"appEngineRelease,omitempty"`

	// Availability: Availability of the instance.@OutputOnly
	//
	// Possible values:
	//   "UNSPECIFIED"
	//   "RESIDENT"
	//   "DYNAMIC"
	Availability string `json:"availability,omitempty"`

	// AverageLatency: Average latency (ms) over the last minute.@OutputOnly
	AverageLatency int64 `json:"averageLatency,omitempty"`

	// Errors: Number of errors since this instance was started.@OutputOnly
	Errors int64 `json:"errors,omitempty"`

	// Id: Relative name of the instance within the version. Example:
	// instance-1.@OutputOnly
	Id string `json:"id,omitempty"`

	// MemoryUsage: Total memory in use (bytes).@OutputOnly
	MemoryUsage int64 `json:"memoryUsage,omitempty,string"`

	// Name: Full path to the Instance resource in the API. Example:
	// apps/myapp/services/default/versions/v1/instances/instance-1.@OutputOn
	// ly
	Name string `json:"name,omitempty"`

	// Qps: Average queries per second (QPS) over the last
	// minute.@OutputOnly
	Qps float64 `json:"qps,omitempty"`

	// Requests: Number of requests since this instance was
	// started.@OutputOnly
	Requests int64 `json:"requests,omitempty"`

	// StartTime: Time that this instance was started.@OutputOnly
	StartTime string `json:"startTime,omitempty"`

	// VmDebugEnabled: Whether this instance is in debug mode. Only
	// applicable for instances in App Engine flexible
	// environment.@OutputOnly
	VmDebugEnabled bool `json:"vmDebugEnabled,omitempty"`

	// VmId: Virtual machine ID of this instance. Only applicable for
	// instances in App Engine flexible environment.@OutputOnly
	VmId string `json:"vmId,omitempty"`

	// VmIp: The IP address of this instance. Only applicable for instances
	// in App Engine flexible environment.@OutputOnly
	VmIp string `json:"vmIp,omitempty"`

	// VmName: Name of the virtual machine where this instance lives. Only
	// applicable for instances in App Engine flexible
	// environment.@OutputOnly
	VmName string `json:"vmName,omitempty"`

	// VmStatus: Status of the virtual machine where this instance lives.
	// Only applicable for instances in App Engine flexible
	// environment.@OutputOnly
	VmStatus string `json:"vmStatus,omitempty"`

	// VmZoneName: Zone where the virtual machine is located. Only
	// applicable for instances in App Engine flexible
	// environment.@OutputOnly
	VmZoneName string `json:"vmZoneName,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "AppEngineRelease") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AppEngineRelease") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *Instance) MarshalJSON() ([]byte, error) {
	type noMethod Instance
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

func (s *Instance) UnmarshalJSON(data []byte) error {
	type noMethod Instance
	var s1 struct {
		Qps gensupport.JSONFloat64 `json:"qps"`
		*noMethod
	}
	s1.noMethod = (*noMethod)(s)
	if err := json.Unmarshal(data, &s1); err != nil {
		return err
	}
	s.Qps = float64(s1.Qps)
	return nil
}

// Library: Third-party Python runtime library that is required by the
// application.
type Library struct {
	// Name: Name of the library. Example: "django".
	Name string `json:"name,omitempty"`

	// Version: Version of the library to select, or "latest".
	Version string `json:"version,omitempty"`

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

func (s *Library) MarshalJSON() ([]byte, error) {
	type noMethod Library
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListAuthorizedCertificatesResponse: Response message for
// AuthorizedCertificates.ListAuthorizedCertificates.
type ListAuthorizedCertificatesResponse struct {
	// Certificates: The SSL certificates the user is authorized to
	// administer.
	Certificates []*AuthorizedCertificate `json:"certificates,omitempty"`

	// NextPageToken: Continuation token for fetching the next page of
	// results.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Certificates") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Certificates") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListAuthorizedCertificatesResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListAuthorizedCertificatesResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListAuthorizedDomainsResponse: Response message for
// AuthorizedDomains.ListAuthorizedDomains.
type ListAuthorizedDomainsResponse struct {
	// Domains: The authorized domains belonging to the user.
	Domains []*AuthorizedDomain `json:"domains,omitempty"`

	// NextPageToken: Continuation token for fetching the next page of
	// results.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Domains") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Domains") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListAuthorizedDomainsResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListAuthorizedDomainsResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListDomainMappingsResponse: Response message for
// DomainMappings.ListDomainMappings.
type ListDomainMappingsResponse struct {
	// DomainMappings: The domain mappings for the application.
	DomainMappings []*DomainMapping `json:"domainMappings,omitempty"`

	// NextPageToken: Continuation token for fetching the next page of
	// results.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "DomainMappings") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "DomainMappings") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *ListDomainMappingsResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListDomainMappingsResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListIngressRulesResponse: Response message for
// Firewall.ListIngressRules.
type ListIngressRulesResponse struct {
	// IngressRules: The ingress FirewallRules for this application.
	IngressRules []*FirewallRule `json:"ingressRules,omitempty"`

	// NextPageToken: Continuation token for fetching the next page of
	// results.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "IngressRules") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "IngressRules") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListIngressRulesResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListIngressRulesResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListInstancesResponse: Response message for Instances.ListInstances.
type ListInstancesResponse struct {
	// Instances: The instances belonging to the requested version.
	Instances []*Instance `json:"instances,omitempty"`

	// NextPageToken: Continuation token for fetching the next page of
	// results.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Instances") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Instances") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListInstancesResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListInstancesResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListLocationsResponse: The response message for
// Locations.ListLocations.
type ListLocationsResponse struct {
	// Locations: A list of locations that matches the specified filter in
	// the request.
	Locations []*Location `json:"locations,omitempty"`

	// NextPageToken: The standard List next-page token.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Locations") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Locations") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListLocationsResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListLocationsResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListOperationsResponse: The response message for
// Operations.ListOperations.
type ListOperationsResponse struct {
	// NextPageToken: The standard List next-page token.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// Operations: A list of operations that matches the specified filter in
	// the request.
	Operations []*Operation `json:"operations,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "NextPageToken") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "NextPageToken") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListOperationsResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListOperationsResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListServicesResponse: Response message for Services.ListServices.
type ListServicesResponse struct {
	// NextPageToken: Continuation token for fetching the next page of
	// results.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// Services: The services belonging to the requested application.
	Services []*Service `json:"services,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "NextPageToken") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "NextPageToken") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListServicesResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListServicesResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListVersionsResponse: Response message for Versions.ListVersions.
type ListVersionsResponse struct {
	// NextPageToken: Continuation token for fetching the next page of
	// results.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// Versions: The versions belonging to the requested service.
	Versions []*Version `json:"versions,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "NextPageToken") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "NextPageToken") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListVersionsResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListVersionsResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// LivenessCheck: Health checking configuration for VM instances.
// Unhealthy instances are killed and replaced with new instances.
type LivenessCheck struct {
	// CheckInterval: Interval between health checks.
	CheckInterval string `json:"checkInterval,omitempty"`

	// FailureThreshold: Number of consecutive failed checks required before
	// considering the VM unhealthy.
	FailureThreshold int64 `json:"failureThreshold,omitempty"`

	// Host: Host header to send when performing a HTTP Liveness check.
	// Example: "myapp.appspot.com"
	Host string `json:"host,omitempty"`

	// InitialDelay: The initial delay before starting to execute the
	// checks.
	InitialDelay string `json:"initialDelay,omitempty"`

	// Path: The request path.
	Path string `json:"path,omitempty"`

	// SuccessThreshold: Number of consecutive successful checks required
	// before considering the VM healthy.
	SuccessThreshold int64 `json:"successThreshold,omitempty"`

	// Timeout: Time before the check is considered failed.
	Timeout string `json:"timeout,omitempty"`

	// ForceSendFields is a list of field names (e.g. "CheckInterval") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CheckInterval") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *LivenessCheck) MarshalJSON() ([]byte, error) {
	type noMethod LivenessCheck
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Location: A resource that represents Google Cloud Platform location.
type Location struct {
	// Labels: Cross-service attributes for the location. For
	// example
	// {"cloud.googleapis.com/region": "us-east1"}
	//
	Labels map[string]string `json:"labels,omitempty"`

	// LocationId: The canonical id for this location. For example:
	// "us-east1".
	LocationId string `json:"locationId,omitempty"`

	// Metadata: Service-specific metadata. For example the available
	// capacity at the given location.
	Metadata googleapi.RawMessage `json:"metadata,omitempty"`

	// Name: Resource name for the location, which may vary between
	// implementations. For example:
	// "projects/example-project/locations/us-east1"
	Name string `json:"name,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Labels") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Labels") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Location) MarshalJSON() ([]byte, error) {
	type noMethod Location
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// LocationMetadata: Metadata for the given
// google.cloud.location.Location.
type LocationMetadata struct {
	// FlexibleEnvironmentAvailable: App Engine Flexible Environment is
	// available in the given location.@OutputOnly
	FlexibleEnvironmentAvailable bool `json:"flexibleEnvironmentAvailable,omitempty"`

	// StandardEnvironmentAvailable: App Engine Standard Environment is
	// available in the given location.@OutputOnly
	StandardEnvironmentAvailable bool `json:"standardEnvironmentAvailable,omitempty"`

	// ForceSendFields is a list of field names (e.g.
	// "FlexibleEnvironmentAvailable") to unconditionally include in API
	// requests. By default, fields with empty values are omitted from API
	// requests. However, any non-pointer, non-interface field appearing in
	// ForceSendFields will be sent to the server regardless of whether the
	// field is empty or not. This may be used to include empty fields in
	// Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g.
	// "FlexibleEnvironmentAvailable") to include in API requests with the
	// JSON null value. By default, fields with empty values are omitted
	// from API requests. However, any field with an empty value appearing
	// in NullFields will be sent to the server as null. It is an error if a
	// field in this list has a non-empty value. This may be used to include
	// null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *LocationMetadata) MarshalJSON() ([]byte, error) {
	type noMethod LocationMetadata
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ManagedCertificate: A certificate managed by App Engine.
type ManagedCertificate struct {
	// LastRenewalTime: Time at which the certificate was last renewed. The
	// renewal process is fully managed. Certificate renewal will
	// automatically occur before the certificate expires. Renewal errors
	// can be tracked via ManagementStatus.@OutputOnly
	LastRenewalTime string `json:"lastRenewalTime,omitempty"`

	// Status: Status of certificate management. Refers to the most recent
	// certificate acquisition or renewal attempt.@OutputOnly
	//
	// Possible values:
	//   "MANAGEMENT_STATUS_UNSPECIFIED"
	//   "OK" - Certificate was successfully obtained and inserted into the
	// serving system.
	//   "PENDING" - Certificate is under active attempts to acquire or
	// renew.
	//   "FAILED_RETRYING_NOT_VISIBLE" - Most recent renewal failed due to
	// an invalid DNS setup and will be retried. Renewal attempts will
	// continue to fail until the certificate domain's DNS configuration is
	// fixed. The last successfully provisioned certificate may still be
	// serving.
	//   "FAILED_PERMANENT" - All renewal attempts have been exhausted,
	// likely due to an invalid DNS setup.
	Status string `json:"status,omitempty"`

	// ForceSendFields is a list of field names (e.g. "LastRenewalTime") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "LastRenewalTime") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *ManagedCertificate) MarshalJSON() ([]byte, error) {
	type noMethod ManagedCertificate
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ManualScaling: A service with manual scaling runs continuously,
// allowing you to perform complex initialization and rely on the state
// of its memory over time.
type ManualScaling struct {
	// Instances: Number of instances to assign to the service at the start.
	// This number can later be altered by using the Modules API
	// (https://cloud.google.com/appengine/docs/python/modules/functions)
	// set_num_instances() function.
	Instances int64 `json:"instances,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Instances") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Instances") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ManualScaling) MarshalJSON() ([]byte, error) {
	type noMethod ManualScaling
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Network: Extra network settings. Only applicable for App Engine
// flexible environment versions
type Network struct {
	// ForwardedPorts: List of ports, or port pairs, to forward from the
	// virtual machine to the application container. Only applicable for App
	// Engine flexible environment versions.
	ForwardedPorts []string `json:"forwardedPorts,omitempty"`

	// InstanceTag: Tag to apply to the VM instance during creation. Only
	// applicable for for App Engine flexible environment versions.
	InstanceTag string `json:"instanceTag,omitempty"`

	// Name: Google Compute Engine network where the virtual machines are
	// created. Specify the short name, not the resource path.Defaults to
	// default.
	Name string `json:"name,omitempty"`

	// SubnetworkName: Google Cloud Platform sub-network where the virtual
	// machines are created. Specify the short name, not the resource
	// path.If a subnetwork name is specified, a network name will also be
	// required unless it is for the default network.
	// If the network the VM instance is being created in is a Legacy
	// network, then the IP address is allocated from the IPv4Range.
	// If the network the VM instance is being created in is an auto Subnet
	// Mode Network, then only network name should be specified (not the
	// subnetwork_name) and the IP address is created from the IPCidrRange
	// of the subnetwork that exists in that zone for that network.
	// If the network the VM instance is being created in is a custom Subnet
	// Mode Network, then the subnetwork_name must be specified and the IP
	// address is created from the IPCidrRange of the subnetwork.If
	// specified, the subnetwork must exist in the same region as the App
	// Engine flexible environment application.
	SubnetworkName string `json:"subnetworkName,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ForwardedPorts") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ForwardedPorts") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *Network) MarshalJSON() ([]byte, error) {
	type noMethod Network
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// NetworkUtilization: Target scaling by network usage. Only applicable
// for VM runtimes.
type NetworkUtilization struct {
	// TargetReceivedBytesPerSecond: Target bytes received per second.
	TargetReceivedBytesPerSecond int64 `json:"targetReceivedBytesPerSecond,omitempty"`

	// TargetReceivedPacketsPerSecond: Target packets received per second.
	TargetReceivedPacketsPerSecond int64 `json:"targetReceivedPacketsPerSecond,omitempty"`

	// TargetSentBytesPerSecond: Target bytes sent per second.
	TargetSentBytesPerSecond int64 `json:"targetSentBytesPerSecond,omitempty"`

	// TargetSentPacketsPerSecond: Target packets sent per second.
	TargetSentPacketsPerSecond int64 `json:"targetSentPacketsPerSecond,omitempty"`

	// ForceSendFields is a list of field names (e.g.
	// "TargetReceivedBytesPerSecond") to unconditionally include in API
	// requests. By default, fields with empty values are omitted from API
	// requests. However, any non-pointer, non-interface field appearing in
	// ForceSendFields will be sent to the server regardless of whether the
	// field is empty or not. This may be used to include empty fields in
	// Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g.
	// "TargetReceivedBytesPerSecond") to include in API requests with the
	// JSON null value. By default, fields with empty values are omitted
	// from API requests. However, any field with an empty value appearing
	// in NullFields will be sent to the server as null. It is an error if a
	// field in this list has a non-empty value. This may be used to include
	// null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *NetworkUtilization) MarshalJSON() ([]byte, error) {
	type noMethod NetworkUtilization
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Operation: This resource represents a long-running operation that is
// the result of a network API call.
type Operation struct {
	// Done: If the value is false, it means the operation is still in
	// progress. If true, the operation is completed, and either error or
	// response is available.
	Done bool `json:"done,omitempty"`

	// Error: The error result of the operation in case of failure or
	// cancellation.
	Error *Status `json:"error,omitempty"`

	// Metadata: Service-specific metadata associated with the operation. It
	// typically contains progress information and common metadata such as
	// create time. Some services might not provide such metadata. Any
	// method that returns a long-running operation should document the
	// metadata type, if any.
	Metadata googleapi.RawMessage `json:"metadata,omitempty"`

	// Name: The server-assigned name, which is only unique within the same
	// service that originally returns it. If you use the default HTTP
	// mapping, the name should have the format of
	// operations/some/unique/name.
	Name string `json:"name,omitempty"`

	// Response: The normal response of the operation in case of success. If
	// the original method returns no data on success, such as Delete, the
	// response is google.protobuf.Empty. If the original method is standard
	// Get/Create/Update, the response should be the resource. For other
	// methods, the response should have the type XxxResponse, where Xxx is
	// the original method name. For example, if the original method name is
	// TakeSnapshot(), the inferred response type is TakeSnapshotResponse.
	Response googleapi.RawMessage `json:"response,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Done") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Done") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Operation) MarshalJSON() ([]byte, error) {
	type noMethod Operation
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// OperationMetadata: Metadata for the given
// google.longrunning.Operation.
type OperationMetadata struct {
	// EndTime: Timestamp that this operation completed.@OutputOnly
	EndTime string `json:"endTime,omitempty"`

	// InsertTime: Timestamp that this operation was created.@OutputOnly
	InsertTime string `json:"insertTime,omitempty"`

	// Method: API method that initiated this operation. Example:
	// google.appengine.v1beta4.Version.CreateVersion.@OutputOnly
	Method string `json:"method,omitempty"`

	// OperationType: Type of this operation. Deprecated, use method field
	// instead. Example: "create_version".@OutputOnly
	OperationType string `json:"operationType,omitempty"`

	// Target: Name of the resource that this operation is acting on.
	// Example: apps/myapp/modules/default.@OutputOnly
	Target string `json:"target,omitempty"`

	// User: User who requested this operation.@OutputOnly
	User string `json:"user,omitempty"`

	// ForceSendFields is a list of field names (e.g. "EndTime") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "EndTime") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *OperationMetadata) MarshalJSON() ([]byte, error) {
	type noMethod OperationMetadata
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// OperationMetadataExperimental: Metadata for the given
// google.longrunning.Operation.
type OperationMetadataExperimental struct {
	// EndTime: Time that this operation completed.@OutputOnly
	EndTime string `json:"endTime,omitempty"`

	// InsertTime: Time that this operation was created.@OutputOnly
	InsertTime string `json:"insertTime,omitempty"`

	// Method: API method that initiated this operation. Example:
	// google.appengine.experimental.CustomDomains.CreateCustomDomain.@Output
	// Only
	Method string `json:"method,omitempty"`

	// Target: Name of the resource that this operation is acting on.
	// Example: apps/myapp/customDomains/example.com.@OutputOnly
	Target string `json:"target,omitempty"`

	// User: User who requested this operation.@OutputOnly
	User string `json:"user,omitempty"`

	// ForceSendFields is a list of field names (e.g. "EndTime") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "EndTime") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *OperationMetadataExperimental) MarshalJSON() ([]byte, error) {
	type noMethod OperationMetadataExperimental
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// OperationMetadataV1: Metadata for the given
// google.longrunning.Operation.
type OperationMetadataV1 struct {
	// EndTime: Time that this operation completed.@OutputOnly
	EndTime string `json:"endTime,omitempty"`

	// EphemeralMessage: Ephemeral message that may change every time the
	// operation is polled. @OutputOnly
	EphemeralMessage string `json:"ephemeralMessage,omitempty"`

	// InsertTime: Time that this operation was created.@OutputOnly
	InsertTime string `json:"insertTime,omitempty"`

	// Method: API method that initiated this operation. Example:
	// google.appengine.v1.Versions.CreateVersion.@OutputOnly
	Method string `json:"method,omitempty"`

	// Target: Name of the resource that this operation is acting on.
	// Example: apps/myapp/services/default.@OutputOnly
	Target string `json:"target,omitempty"`

	// User: User who requested this operation.@OutputOnly
	User string `json:"user,omitempty"`

	// Warning: Durable messages that persist on every operation poll.
	// @OutputOnly
	Warning []string `json:"warning,omitempty"`

	// ForceSendFields is a list of field names (e.g. "EndTime") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "EndTime") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *OperationMetadataV1) MarshalJSON() ([]byte, error) {
	type noMethod OperationMetadataV1
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// OperationMetadataV1Alpha: Metadata for the given
// google.longrunning.Operation.
type OperationMetadataV1Alpha struct {
	// EndTime: Time that this operation completed.@OutputOnly
	EndTime string `json:"endTime,omitempty"`

	// EphemeralMessage: Ephemeral message that may change every time the
	// operation is polled. @OutputOnly
	EphemeralMessage string `json:"ephemeralMessage,omitempty"`

	// InsertTime: Time that this operation was created.@OutputOnly
	InsertTime string `json:"insertTime,omitempty"`

	// Method: API method that initiated this operation. Example:
	// google.appengine.v1alpha.Versions.CreateVersion.@OutputOnly
	Method string `json:"method,omitempty"`

	// Target: Name of the resource that this operation is acting on.
	// Example: apps/myapp/services/default.@OutputOnly
	Target string `json:"target,omitempty"`

	// User: User who requested this operation.@OutputOnly
	User string `json:"user,omitempty"`

	// Warning: Durable messages that persist on every operation poll.
	// @OutputOnly
	Warning []string `json:"warning,omitempty"`

	// ForceSendFields is a list of field names (e.g. "EndTime") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "EndTime") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *OperationMetadataV1Alpha) MarshalJSON() ([]byte, error) {
	type noMethod OperationMetadataV1Alpha
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// OperationMetadataV1Beta: Metadata for the given
// google.longrunning.Operation.
type OperationMetadataV1Beta struct {
	// EndTime: Time that this operation completed.@OutputOnly
	EndTime string `json:"endTime,omitempty"`

	// EphemeralMessage: Ephemeral message that may change every time the
	// operation is polled. @OutputOnly
	EphemeralMessage string `json:"ephemeralMessage,omitempty"`

	// InsertTime: Time that this operation was created.@OutputOnly
	InsertTime string `json:"insertTime,omitempty"`

	// Method: API method that initiated this operation. Example:
	// google.appengine.v1beta.Versions.CreateVersion.@OutputOnly
	Method string `json:"method,omitempty"`

	// Target: Name of the resource that this operation is acting on.
	// Example: apps/myapp/services/default.@OutputOnly
	Target string `json:"target,omitempty"`

	// User: User who requested this operation.@OutputOnly
	User string `json:"user,omitempty"`

	// Warning: Durable messages that persist on every operation poll.
	// @OutputOnly
	Warning []string `json:"warning,omitempty"`

	// ForceSendFields is a list of field names (e.g. "EndTime") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "EndTime") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *OperationMetadataV1Beta) MarshalJSON() ([]byte, error) {
	type noMethod OperationMetadataV1Beta
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// OperationMetadataV1Beta5: Metadata for the given
// google.longrunning.Operation.
type OperationMetadataV1Beta5 struct {
	// EndTime: Timestamp that this operation completed.@OutputOnly
	EndTime string `json:"endTime,omitempty"`

	// InsertTime: Timestamp that this operation was created.@OutputOnly
	InsertTime string `json:"insertTime,omitempty"`

	// Method: API method name that initiated this operation. Example:
	// google.appengine.v1beta5.Version.CreateVersion.@OutputOnly
	Method string `json:"method,omitempty"`

	// Target: Name of the resource that this operation is acting on.
	// Example: apps/myapp/services/default.@OutputOnly
	Target string `json:"target,omitempty"`

	// User: User who requested this operation.@OutputOnly
	User string `json:"user,omitempty"`

	// ForceSendFields is a list of field names (e.g. "EndTime") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "EndTime") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *OperationMetadataV1Beta5) MarshalJSON() ([]byte, error) {
	type noMethod OperationMetadataV1Beta5
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ReadinessCheck: Readiness checking configuration for VM instances.
// Unhealthy instances are removed from traffic rotation.
type ReadinessCheck struct {
	// AppStartTimeout: A maximum time limit on application initialization,
	// measured from moment the application successfully replies to a
	// healthcheck until it is ready to serve traffic.
	AppStartTimeout string `json:"appStartTimeout,omitempty"`

	// CheckInterval: Interval between health checks.
	CheckInterval string `json:"checkInterval,omitempty"`

	// FailureThreshold: Number of consecutive failed checks required before
	// removing traffic.
	FailureThreshold int64 `json:"failureThreshold,omitempty"`

	// Host: Host header to send when performing a HTTP Readiness check.
	// Example: "myapp.appspot.com"
	Host string `json:"host,omitempty"`

	// Path: The request path.
	Path string `json:"path,omitempty"`

	// SuccessThreshold: Number of consecutive successful checks required
	// before receiving traffic.
	SuccessThreshold int64 `json:"successThreshold,omitempty"`

	// Timeout: Time before the check is considered failed.
	Timeout string `json:"timeout,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AppStartTimeout") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AppStartTimeout") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *ReadinessCheck) MarshalJSON() ([]byte, error) {
	type noMethod ReadinessCheck
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// RepairApplicationRequest: Request message for
// 'Applications.RepairApplication'.
type RepairApplicationRequest struct {
}

// RequestUtilization: Target scaling by request utilization. Only
// applicable for VM runtimes.
type RequestUtilization struct {
	// TargetConcurrentRequests: Target number of concurrent requests.
	TargetConcurrentRequests int64 `json:"targetConcurrentRequests,omitempty"`

	// TargetRequestCountPerSecond: Target requests per second.
	TargetRequestCountPerSecond int64 `json:"targetRequestCountPerSecond,omitempty"`

	// ForceSendFields is a list of field names (e.g.
	// "TargetConcurrentRequests") to unconditionally include in API
	// requests. By default, fields with empty values are omitted from API
	// requests. However, any non-pointer, non-interface field appearing in
	// ForceSendFields will be sent to the server regardless of whether the
	// field is empty or not. This may be used to include empty fields in
	// Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "TargetConcurrentRequests")
	// to include in API requests with the JSON null value. By default,
	// fields with empty values are omitted from API requests. However, any
	// field with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *RequestUtilization) MarshalJSON() ([]byte, error) {
	type noMethod RequestUtilization
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ResourceRecord: A DNS resource record.
type ResourceRecord struct {
	// Name: Relative name of the object affected by this record. Only
	// applicable for CNAME records. Example: 'www'.
	Name string `json:"name,omitempty"`

	// Rrdata: Data for this record. Values vary by record type, as defined
	// in RFC 1035 (section 5) and RFC 1034 (section 3.6.1).
	Rrdata string `json:"rrdata,omitempty"`

	// Type: Resource record type. Example: AAAA.
	//
	// Possible values:
	//   "A" - An A resource record. Data is an IPv4 address.
	//   "AAAA" - An AAAA resource record. Data is an IPv6 address.
	//   "CNAME" - A CNAME resource record. Data is a domain name to be
	// aliased.
	Type string `json:"type,omitempty"`

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

func (s *ResourceRecord) MarshalJSON() ([]byte, error) {
	type noMethod ResourceRecord
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Resources: Machine resources for a version.
type Resources struct {
	// Cpu: Number of CPU cores needed.
	Cpu float64 `json:"cpu,omitempty"`

	// DiskGb: Disk size (GB) needed.
	DiskGb float64 `json:"diskGb,omitempty"`

	// MemoryGb: Memory (GB) needed.
	MemoryGb float64 `json:"memoryGb,omitempty"`

	// Volumes: User specified volumes.
	Volumes []*Volume `json:"volumes,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Cpu") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Cpu") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Resources) MarshalJSON() ([]byte, error) {
	type noMethod Resources
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

func (s *Resources) UnmarshalJSON(data []byte) error {
	type noMethod Resources
	var s1 struct {
		Cpu      gensupport.JSONFloat64 `json:"cpu"`
		DiskGb   gensupport.JSONFloat64 `json:"diskGb"`
		MemoryGb gensupport.JSONFloat64 `json:"memoryGb"`
		*noMethod
	}
	s1.noMethod = (*noMethod)(s)
	if err := json.Unmarshal(data, &s1); err != nil {
		return err
	}
	s.Cpu = float64(s1.Cpu)
	s.DiskGb = float64(s1.DiskGb)
	s.MemoryGb = float64(s1.MemoryGb)
	return nil
}

// ScriptHandler: Executes a script to handle the request that matches
// the URL pattern.
type ScriptHandler struct {
	// ScriptPath: Path to the script from the application root directory.
	ScriptPath string `json:"scriptPath,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ScriptPath") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ScriptPath") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ScriptHandler) MarshalJSON() ([]byte, error) {
	type noMethod ScriptHandler
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Service: A Service resource is a logical component of an application
// that can share state and communicate in a secure fashion with other
// services. For example, an application that handles customer requests
// might include separate services to handle tasks such as backend data
// analysis or API requests from mobile devices. Each service has a
// collection of versions that define a specific set of code used to
// implement the functionality of that service.
type Service struct {
	// Id: Relative name of the service within the application. Example:
	// default.@OutputOnly
	Id string `json:"id,omitempty"`

	// Name: Full path to the Service resource in the API. Example:
	// apps/myapp/services/default.@OutputOnly
	Name string `json:"name,omitempty"`

	// Split: Mapping that defines fractional HTTP traffic diversion to
	// different versions within the service.
	Split *TrafficSplit `json:"split,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Id") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Id") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Service) MarshalJSON() ([]byte, error) {
	type noMethod Service
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// SslSettings: SSL configuration for a DomainMapping resource.
type SslSettings struct {
	// CertificateId: ID of the AuthorizedCertificate resource configuring
	// SSL for the application. Clearing this field will remove SSL
	// support.By default, a managed certificate is automatically created
	// for every domain mapping. To omit SSL support or to configure SSL
	// manually, specify SslManagementType.MANUAL on a CREATE or UPDATE
	// request. You must be authorized to administer the
	// AuthorizedCertificate resource to manually map it to a DomainMapping
	// resource. Example: 12345.
	CertificateId string `json:"certificateId,omitempty"`

	// PendingManagedCertificateId: ID of the managed AuthorizedCertificate
	// resource currently being provisioned, if applicable. Until the new
	// managed certificate has been successfully provisioned, the previous
	// SSL state will be preserved. Once the provisioning process completes,
	// the certificate_id field will reflect the new managed certificate and
	// this field will be left empty. To remove SSL support while there is
	// still a pending managed certificate, clear the certificate_id field
	// with an UpdateDomainMappingRequest.@OutputOnly
	PendingManagedCertificateId string `json:"pendingManagedCertificateId,omitempty"`

	// SslManagementType: SSL management type for this domain. If AUTOMATIC,
	// a managed certificate is automatically provisioned. If MANUAL,
	// certificate_id must be manually specified in order to configure SSL
	// for this domain.
	//
	// Possible values:
	//   "AUTOMATIC" - SSL support for this domain is configured
	// automatically. The mapped SSL certificate will be automatically
	// renewed.
	//   "MANUAL" - SSL support for this domain is configured manually by
	// the user. Either the domain has no SSL support or a user-obtained SSL
	// certificate has been explictly mapped to this domain.
	SslManagementType string `json:"sslManagementType,omitempty"`

	// ForceSendFields is a list of field names (e.g. "CertificateId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CertificateId") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *SslSettings) MarshalJSON() ([]byte, error) {
	type noMethod SslSettings
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// StandardSchedulerSettings: Scheduler settings for standard
// environment.
type StandardSchedulerSettings struct {
	// MaxInstances: Maximum number of instances for an app version. Set to
	// a non-positive value (0 by convention) to disable max_instances
	// configuration.
	MaxInstances int64 `json:"maxInstances,omitempty"`

	// MinInstances: Minimum number of instances for an app version. Set to
	// a non-positive value (0 by convention) to disable min_instances
	// configuration.
	MinInstances int64 `json:"minInstances,omitempty"`

	// TargetCpuUtilization: Target CPU utilization ratio to maintain when
	// scaling.
	TargetCpuUtilization float64 `json:"targetCpuUtilization,omitempty"`

	// TargetThroughputUtilization: Target throughput utilization ratio to
	// maintain when scaling
	TargetThroughputUtilization float64 `json:"targetThroughputUtilization,omitempty"`

	// ForceSendFields is a list of field names (e.g. "MaxInstances") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "MaxInstances") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *StandardSchedulerSettings) MarshalJSON() ([]byte, error) {
	type noMethod StandardSchedulerSettings
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

func (s *StandardSchedulerSettings) UnmarshalJSON(data []byte) error {
	type noMethod StandardSchedulerSettings
	var s1 struct {
		TargetCpuUtilization        gensupport.JSONFloat64 `json:"targetCpuUtilization"`
		TargetThroughputUtilization gensupport.JSONFloat64 `json:"targetThroughputUtilization"`
		*noMethod
	}
	s1.noMethod = (*noMethod)(s)
	if err := json.Unmarshal(data, &s1); err != nil {
		return err
	}
	s.TargetCpuUtilization = float64(s1.TargetCpuUtilization)
	s.TargetThroughputUtilization = float64(s1.TargetThroughputUtilization)
	return nil
}

// StaticFilesHandler: Files served directly to the user for a given
// URL, such as images, CSS stylesheets, or JavaScript source files.
// Static file handlers describe which files in the application
// directory are static files, and which URLs serve them.
type StaticFilesHandler struct {
	// ApplicationReadable: Whether files should also be uploaded as code
	// data. By default, files declared in static file handlers are uploaded
	// as static data and are only served to end users; they cannot be read
	// by the application. If enabled, uploads are charged against both your
	// code and static data storage resource quotas.
	ApplicationReadable bool `json:"applicationReadable,omitempty"`

	// Expiration: Time a static file served by this handler should be
	// cached by web proxies and browsers.
	Expiration string `json:"expiration,omitempty"`

	// HttpHeaders: HTTP headers to use for all responses from these URLs.
	HttpHeaders map[string]string `json:"httpHeaders,omitempty"`

	// MimeType: MIME type used to serve all files served by this
	// handler.Defaults to file-specific MIME types, which are derived from
	// each file's filename extension.
	MimeType string `json:"mimeType,omitempty"`

	// Path: Path to the static files matched by the URL pattern, from the
	// application root directory. The path can refer to text matched in
	// groupings in the URL pattern.
	Path string `json:"path,omitempty"`

	// RequireMatchingFile: Whether this handler should match the request if
	// the file referenced by the handler does not exist.
	RequireMatchingFile bool `json:"requireMatchingFile,omitempty"`

	// UploadPathRegex: Regular expression that matches the file paths for
	// all files that should be referenced by this handler.
	UploadPathRegex string `json:"uploadPathRegex,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ApplicationReadable")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ApplicationReadable") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *StaticFilesHandler) MarshalJSON() ([]byte, error) {
	type noMethod StaticFilesHandler
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Status: The Status type defines a logical error model that is
// suitable for different programming environments, including REST APIs
// and RPC APIs. It is used by gRPC (https://github.com/grpc). The error
// model is designed to be:
// Simple to use and understand for most users
// Flexible enough to meet unexpected needsOverviewThe Status message
// contains three pieces of data: error code, error message, and error
// details. The error code should be an enum value of google.rpc.Code,
// but it may accept additional error codes if needed. The error message
// should be a developer-facing English message that helps developers
// understand and resolve the error. If a localized user-facing error
// message is needed, put the localized message in the error details or
// localize it in the client. The optional error details may contain
// arbitrary information about the error. There is a predefined set of
// error detail types in the package google.rpc that can be used for
// common error conditions.Language mappingThe Status message is the
// logical representation of the error model, but it is not necessarily
// the actual wire format. When the Status message is exposed in
// different client libraries and different wire protocols, it can be
// mapped differently. For example, it will likely be mapped to some
// exceptions in Java, but more likely mapped to some error codes in
// C.Other usesThe error model and the Status message can be used in a
// variety of environments, either with or without APIs, to provide a
// consistent developer experience across different environments.Example
// uses of this error model include:
// Partial errors. If a service needs to return partial errors to the
// client, it may embed the Status in the normal response to indicate
// the partial errors.
// Workflow errors. A typical workflow has multiple steps. Each step may
// have a Status message for error reporting.
// Batch operations. If a client uses batch request and batch response,
// the Status message should be used directly inside batch response, one
// for each error sub-response.
// Asynchronous operations. If an API call embeds asynchronous operation
// results in its response, the status of those operations should be
// represented directly using the Status message.
// Logging. If some API errors are stored in logs, the message Status
// could be used directly after any stripping needed for
// security/privacy reasons.
type Status struct {
	// Code: The status code, which should be an enum value of
	// google.rpc.Code.
	Code int64 `json:"code,omitempty"`

	// Details: A list of messages that carry the error details. There is a
	// common set of message types for APIs to use.
	Details []googleapi.RawMessage `json:"details,omitempty"`

	// Message: A developer-facing error message, which should be in
	// English. Any user-facing error message should be localized and sent
	// in the google.rpc.Status.details field, or localized by the client.
	Message string `json:"message,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Code") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Code") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Status) MarshalJSON() ([]byte, error) {
	type noMethod Status
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// TrafficSplit: Traffic routing configuration for versions within a
// single service. Traffic splits define how traffic directed to the
// service is assigned to versions.
type TrafficSplit struct {
	// Allocations: Mapping from version IDs within the service to
	// fractional (0.000, 1] allocations of traffic for that version. Each
	// version can be specified only once, but some versions in the service
	// may not have any traffic allocation. Services that have traffic
	// allocated cannot be deleted until either the service is deleted or
	// their traffic allocation is removed. Allocations must sum to 1. Up to
	// two decimal place precision is supported for IP-based splits and up
	// to three decimal places is supported for cookie-based splits.
	Allocations map[string]float64 `json:"allocations,omitempty"`

	// ShardBy: Mechanism used to determine which version a request is sent
	// to. The traffic selection algorithm will be stable for either type
	// until allocations are changed.
	//
	// Possible values:
	//   "UNSPECIFIED" - Diversion method unspecified.
	//   "COOKIE" - Diversion based on a specially named cookie,
	// "GOOGAPPUID." The cookie must be set by the application itself or no
	// diversion will occur.
	//   "IP" - Diversion based on applying the modulus operation to a
	// fingerprint of the IP address.
	//   "RANDOM" - Diversion based on weighted random assignment. An
	// incoming request is randomly routed to a version in the traffic
	// split, with probability proportional to the version's traffic share.
	ShardBy string `json:"shardBy,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Allocations") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Allocations") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *TrafficSplit) MarshalJSON() ([]byte, error) {
	type noMethod TrafficSplit
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// UrlDispatchRule: Rules to match an HTTP request and dispatch that
// request to a service.
type UrlDispatchRule struct {
	// Domain: Domain name to match against. The wildcard "*" is supported
	// if specified before a period: "*.".Defaults to matching all domains:
	// "*".
	Domain string `json:"domain,omitempty"`

	// Path: Pathname within the host. Must start with a "/". A single "*"
	// can be included at the end of the path.The sum of the lengths of the
	// domain and path may not exceed 100 characters.
	Path string `json:"path,omitempty"`

	// Service: Resource ID of a service in this application that should
	// serve the matched request. The service must already exist. Example:
	// default.
	Service string `json:"service,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Domain") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Domain") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *UrlDispatchRule) MarshalJSON() ([]byte, error) {
	type noMethod UrlDispatchRule
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// UrlMap: URL pattern and description of how the URL should be handled.
// App Engine can handle URLs by executing application code or by
// serving static files uploaded with the version, such as images, CSS,
// or JavaScript.
type UrlMap struct {
	// ApiEndpoint: Uses API Endpoints to handle requests.
	ApiEndpoint *ApiEndpointHandler `json:"apiEndpoint,omitempty"`

	// AuthFailAction: Action to take when users access resources that
	// require authentication. Defaults to redirect.
	//
	// Possible values:
	//   "AUTH_FAIL_ACTION_UNSPECIFIED" - Not specified.
	// AUTH_FAIL_ACTION_REDIRECT is assumed.
	//   "AUTH_FAIL_ACTION_REDIRECT" - Redirects user to
	// "accounts.google.com". The user is redirected back to the application
	// URL after signing in or creating an account.
	//   "AUTH_FAIL_ACTION_UNAUTHORIZED" - Rejects request with a 401 HTTP
	// status code and an error message.
	AuthFailAction string `json:"authFailAction,omitempty"`

	// Login: Level of login required to access this resource.
	//
	// Possible values:
	//   "LOGIN_UNSPECIFIED" - Not specified. LOGIN_OPTIONAL is assumed.
	//   "LOGIN_OPTIONAL" - Does not require that the user is signed in.
	//   "LOGIN_ADMIN" - If the user is not signed in, the auth_fail_action
	// is taken. In addition, if the user is not an administrator for the
	// application, they are given an error message regardless of
	// auth_fail_action. If the user is an administrator, the handler
	// proceeds.
	//   "LOGIN_REQUIRED" - If the user has signed in, the handler proceeds
	// normally. Otherwise, the auth_fail_action is taken.
	Login string `json:"login,omitempty"`

	// RedirectHttpResponseCode: 30x code to use when performing redirects
	// for the secure field. Defaults to 302.
	//
	// Possible values:
	//   "REDIRECT_HTTP_RESPONSE_CODE_UNSPECIFIED" - Not specified. 302 is
	// assumed.
	//   "REDIRECT_HTTP_RESPONSE_CODE_301" - 301 Moved Permanently code.
	//   "REDIRECT_HTTP_RESPONSE_CODE_302" - 302 Moved Temporarily code.
	//   "REDIRECT_HTTP_RESPONSE_CODE_303" - 303 See Other code.
	//   "REDIRECT_HTTP_RESPONSE_CODE_307" - 307 Temporary Redirect code.
	RedirectHttpResponseCode string `json:"redirectHttpResponseCode,omitempty"`

	// Script: Executes a script to handle the request that matches this URL
	// pattern.
	Script *ScriptHandler `json:"script,omitempty"`

	// SecurityLevel: Security (HTTPS) enforcement for this URL.
	//
	// Possible values:
	//   "SECURE_UNSPECIFIED" - Not specified.
	//   "SECURE_DEFAULT" - Both HTTP and HTTPS requests with URLs that
	// match the handler succeed without redirects. The application can
	// examine the request to determine which protocol was used, and respond
	// accordingly.
	//   "SECURE_NEVER" - Requests for a URL that match this handler that
	// use HTTPS are automatically redirected to the HTTP equivalent URL.
	//   "SECURE_OPTIONAL" - Both HTTP and HTTPS requests with URLs that
	// match the handler succeed without redirects. The application can
	// examine the request to determine which protocol was used and respond
	// accordingly.
	//   "SECURE_ALWAYS" - Requests for a URL that match this handler that
	// do not use HTTPS are automatically redirected to the HTTPS URL with
	// the same path. Query parameters are reserved for the redirect.
	SecurityLevel string `json:"securityLevel,omitempty"`

	// StaticFiles: Returns the contents of a file, such as an image, as the
	// response.
	StaticFiles *StaticFilesHandler `json:"staticFiles,omitempty"`

	// UrlRegex: URL prefix. Uses regular expression syntax, which means
	// regexp special characters must be escaped, but should not contain
	// groupings. All URLs that begin with this prefix are handled by this
	// handler, using the portion of the URL after the prefix as part of the
	// file path.
	UrlRegex string `json:"urlRegex,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ApiEndpoint") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ApiEndpoint") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *UrlMap) MarshalJSON() ([]byte, error) {
	type noMethod UrlMap
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Version: A Version resource is a specific set of source code and
// configuration files that are deployed into a service.
type Version struct {
	// ApiConfig: Serving configuration for Google Cloud Endpoints
	// (https://cloud.google.com/appengine/docs/python/endpoints/).Only
	// returned in GET requests if view=FULL is set.
	ApiConfig *ApiConfigHandler `json:"apiConfig,omitempty"`

	// AutomaticScaling: Automatic scaling is based on request rate,
	// response latencies, and other application metrics.
	AutomaticScaling *AutomaticScaling `json:"automaticScaling,omitempty"`

	// BasicScaling: A service with basic scaling will create an instance
	// when the application receives a request. The instance will be turned
	// down when the app becomes idle. Basic scaling is ideal for work that
	// is intermittent or driven by user activity.
	BasicScaling *BasicScaling `json:"basicScaling,omitempty"`

	// BetaSettings: Metadata settings that are supplied to this version to
	// enable beta runtime features.
	BetaSettings map[string]string `json:"betaSettings,omitempty"`

	// CreateTime: Time that this version was created.@OutputOnly
	CreateTime string `json:"createTime,omitempty"`

	// CreatedBy: Email address of the user who created this
	// version.@OutputOnly
	CreatedBy string `json:"createdBy,omitempty"`

	// DefaultExpiration: Duration that static files should be cached by web
	// proxies and browsers. Only applicable if the corresponding
	// StaticFilesHandler
	// (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1be
	// ta/apps.services.versions#staticfileshandler) does not specify its
	// own expiration time.Only returned in GET requests if view=FULL is
	// set.
	DefaultExpiration string `json:"defaultExpiration,omitempty"`

	// Deployment: Code and application artifacts that make up this
	// version.Only returned in GET requests if view=FULL is set.
	Deployment *Deployment `json:"deployment,omitempty"`

	// DiskUsageBytes: Total size in bytes of all the files that are
	// included in this version and curerntly hosted on the App Engine
	// disk.@OutputOnly
	DiskUsageBytes int64 `json:"diskUsageBytes,omitempty,string"`

	// EndpointsApiService: Cloud Endpoints configuration.If
	// endpoints_api_service is set, the Cloud Endpoints Extensible Service
	// Proxy will be provided to serve the API implemented by the app.
	EndpointsApiService *EndpointsApiService `json:"endpointsApiService,omitempty"`

	// Env: App Engine execution environment for this version.Defaults to
	// standard.
	Env string `json:"env,omitempty"`

	// EnvVariables: Environment variables available to the application.Only
	// returned in GET requests if view=FULL is set.
	EnvVariables map[string]string `json:"envVariables,omitempty"`

	// ErrorHandlers: Custom static error pages. Limited to 10KB per
	// page.Only returned in GET requests if view=FULL is set.
	ErrorHandlers []*ErrorHandler `json:"errorHandlers,omitempty"`

	// Handlers: An ordered list of URL-matching patterns that should be
	// applied to incoming requests. The first matching URL handles the
	// request and other request handlers are not attempted.Only returned in
	// GET requests if view=FULL is set.
	Handlers []*UrlMap `json:"handlers,omitempty"`

	// HealthCheck: Configures health checking for VM instances. Unhealthy
	// instances are stopped and replaced with new instances. Only
	// applicable for VM runtimes.Only returned in GET requests if view=FULL
	// is set.
	HealthCheck *HealthCheck `json:"healthCheck,omitempty"`

	// Id: Relative name of the version within the service. Example: v1.
	// Version names can contain only lowercase letters, numbers, or
	// hyphens. Reserved names: "default", "latest", and any name with the
	// prefix "ah-".
	Id string `json:"id,omitempty"`

	// InboundServices: Before an application can receive email or XMPP
	// messages, the application must be configured to enable the service.
	//
	// Possible values:
	//   "INBOUND_SERVICE_UNSPECIFIED" - Not specified.
	//   "INBOUND_SERVICE_MAIL" - Allows an application to receive mail.
	//   "INBOUND_SERVICE_MAIL_BOUNCE" - Allows an application to receive
	// email-bound notifications.
	//   "INBOUND_SERVICE_XMPP_ERROR" - Allows an application to receive
	// error stanzas.
	//   "INBOUND_SERVICE_XMPP_MESSAGE" - Allows an application to receive
	// instant messages.
	//   "INBOUND_SERVICE_XMPP_SUBSCRIBE" - Allows an application to receive
	// user subscription POSTs.
	//   "INBOUND_SERVICE_XMPP_PRESENCE" - Allows an application to receive
	// a user's chat presence.
	//   "INBOUND_SERVICE_CHANNEL_PRESENCE" - Registers an application for
	// notifications when a client connects or disconnects from a channel.
	//   "INBOUND_SERVICE_WARMUP" - Enables warmup requests.
	InboundServices []string `json:"inboundServices,omitempty"`

	// InstanceClass: Instance class that is used to run this version. Valid
	// values are:
	// AutomaticScaling: F1, F2, F4, F4_1G
	// ManualScaling or BasicScaling: B1, B2, B4, B8, B4_1GDefaults to F1
	// for AutomaticScaling and B1 for ManualScaling or BasicScaling.
	InstanceClass string `json:"instanceClass,omitempty"`

	// Libraries: Configuration for third-party Python runtime libraries
	// that are required by the application.Only returned in GET requests if
	// view=FULL is set.
	Libraries []*Library `json:"libraries,omitempty"`

	// LivenessCheck: Configures liveness health checking for VM instances.
	// Unhealthy instances are stopped and replaced with new instancesOnly
	// returned in GET requests if view=FULL is set.
	LivenessCheck *LivenessCheck `json:"livenessCheck,omitempty"`

	// ManualScaling: A service with manual scaling runs continuously,
	// allowing you to perform complex initialization and rely on the state
	// of its memory over time.
	ManualScaling *ManualScaling `json:"manualScaling,omitempty"`

	// Name: Full path to the Version resource in the API. Example:
	// apps/myapp/services/default/versions/v1.@OutputOnly
	Name string `json:"name,omitempty"`

	// Network: Extra network settings. Only applicable for App Engine
	// flexible environment versions.
	Network *Network `json:"network,omitempty"`

	// NobuildFilesRegex: Files that match this pattern will not be built
	// into this version. Only applicable for Go runtimes.Only returned in
	// GET requests if view=FULL is set.
	NobuildFilesRegex string `json:"nobuildFilesRegex,omitempty"`

	// ReadinessCheck: Configures readiness health checking for VM
	// instances. Unhealthy instances are not put into the backend traffic
	// rotation.Only returned in GET requests if view=FULL is set.
	ReadinessCheck *ReadinessCheck `json:"readinessCheck,omitempty"`

	// Resources: Machine resources for this version. Only applicable for VM
	// runtimes.
	Resources *Resources `json:"resources,omitempty"`

	// Runtime: Desired runtime. Example: python27.
	Runtime string `json:"runtime,omitempty"`

	// RuntimeApiVersion: The version of the API in the given runtime
	// environment. Please see the app.yaml reference for valid values at
	// https://cloud.google.com/appengine/docs/standard/<language>/config/appref
	RuntimeApiVersion string `json:"runtimeApiVersion,omitempty"`

	// ServingStatus: Current serving status of this version. Only the
	// versions with a SERVING status create instances and can be
	// billed.SERVING_STATUS_UNSPECIFIED is an invalid value. Defaults to
	// SERVING.
	//
	// Possible values:
	//   "SERVING_STATUS_UNSPECIFIED" - Not specified.
	//   "SERVING" - Currently serving. Instances are created according to
	// the scaling settings of the version.
	//   "STOPPED" - Disabled. No instances will be created and the scaling
	// settings are ignored until the state of the version changes to
	// SERVING.
	ServingStatus string `json:"servingStatus,omitempty"`

	// Threadsafe: Whether multiple requests can be dispatched to this
	// version at once.
	Threadsafe bool `json:"threadsafe,omitempty"`

	// VersionUrl: Serving URL for this version. Example:
	// "https://myversion-dot-myservice-dot-myapp.appspot.com"@OutputOnly
	VersionUrl string `json:"versionUrl,omitempty"`

	// Vm: Whether to deploy this version in a container on a virtual
	// machine.
	Vm bool `json:"vm,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "ApiConfig") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ApiConfig") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Version) MarshalJSON() ([]byte, error) {
	type noMethod Version
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Volume: Volumes mounted within the app container. Only applicable for
// VM runtimes.
type Volume struct {
	// Name: Unique name for the volume.
	Name string `json:"name,omitempty"`

	// SizeGb: Volume size in gigabytes.
	SizeGb float64 `json:"sizeGb,omitempty"`

	// VolumeType: Underlying volume type, e.g. 'tmpfs'.
	VolumeType string `json:"volumeType,omitempty"`

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

func (s *Volume) MarshalJSON() ([]byte, error) {
	type noMethod Volume
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

func (s *Volume) UnmarshalJSON(data []byte) error {
	type noMethod Volume
	var s1 struct {
		SizeGb gensupport.JSONFloat64 `json:"sizeGb"`
		*noMethod
	}
	s1.noMethod = (*noMethod)(s)
	if err := json.Unmarshal(data, &s1); err != nil {
		return err
	}
	s.SizeGb = float64(s1.SizeGb)
	return nil
}

// ZipInfo: The zip file information for a zip deployment.
type ZipInfo struct {
	// FilesCount: An estimate of the number of files in a zip for a zip
	// deployment. If set, must be greater than or equal to the actual
	// number of files. Used for optimizing performance; if not provided,
	// deployment may be slow.
	FilesCount int64 `json:"filesCount,omitempty"`

	// SourceUrl: URL of the zip file to deploy from. Must be a URL to a
	// resource in Google Cloud Storage in the form
	// 'http(s)://storage.googleapis.com/<bucket>/<object>'.
	SourceUrl string `json:"sourceUrl,omitempty"`

	// ForceSendFields is a list of field names (e.g. "FilesCount") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "FilesCount") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ZipInfo) MarshalJSON() ([]byte, error) {
	type noMethod ZipInfo
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// method id "appengine.apps.create":

type AppsCreateCall struct {
	s           *APIService
	application *Application
	urlParams_  gensupport.URLParams
	ctx_        context.Context
	header_     http.Header
}

// Create: Creates an App Engine application for a Google Cloud Platform
// project. Required fields:
// id - The ID of the target Cloud Platform project.
// location - The region
// (https://cloud.google.com/appengine/docs/locations) where you want
// the App Engine application located.For more information about App
// Engine applications, see Managing Projects, Applications, and Billing
// (https://cloud.google.com/appengine/docs/python/console/).
func (r *AppsService) Create(application *Application) *AppsCreateCall {
	c := &AppsCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.application = application
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsCreateCall) Fields(s ...googleapi.Field) *AppsCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsCreateCall) Context(ctx context.Context) *AppsCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.application)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.create" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsCreateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	ret := &Operation{
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
	//   "description": "Creates an App Engine application for a Google Cloud Platform project. Required fields:\nid - The ID of the target Cloud Platform project.\nlocation - The region (https://cloud.google.com/appengine/docs/locations) where you want the App Engine application located.For more information about App Engine applications, see Managing Projects, Applications, and Billing (https://cloud.google.com/appengine/docs/python/console/).",
	//   "flatPath": "v1beta/apps",
	//   "httpMethod": "POST",
	//   "id": "appengine.apps.create",
	//   "parameterOrder": [],
	//   "parameters": {},
	//   "path": "v1beta/apps",
	//   "request": {
	//     "$ref": "Application"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.get":

type AppsGetCall struct {
	s            *APIService
	appsId       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Gets information about an application.
func (r *AppsService) Get(appsId string) *AppsGetCall {
	c := &AppsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsGetCall) Fields(s ...googleapi.Field) *AppsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsGetCall) IfNoneMatch(entityTag string) *AppsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsGetCall) Context(ctx context.Context) *AppsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId": c.appsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.get" call.
// Exactly one of *Application or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Application.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsGetCall) Do(opts ...googleapi.CallOption) (*Application, error) {
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
	ret := &Application{
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
	//   "description": "Gets information about an application.",
	//   "flatPath": "v1beta/apps/{appsId}",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.get",
	//   "parameterOrder": [
	//     "appsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the Application resource to get. Example: apps/myapp.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}",
	//   "response": {
	//     "$ref": "Application"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// method id "appengine.apps.patch":

type AppsPatchCall struct {
	s           *APIService
	appsId      string
	application *Application
	urlParams_  gensupport.URLParams
	ctx_        context.Context
	header_     http.Header
}

// Patch: Updates the specified Application resource. You can update the
// following fields:
// auth_domain - Google authentication domain for controlling user
// access to the application.
// default_cookie_expiration - Cookie expiration policy for the
// application.
func (r *AppsService) Patch(appsId string, application *Application) *AppsPatchCall {
	c := &AppsPatchCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.application = application
	return c
}

// UpdateMask sets the optional parameter "updateMask": Standard field
// mask for the set of fields to be updated.
func (c *AppsPatchCall) UpdateMask(updateMask string) *AppsPatchCall {
	c.urlParams_.Set("updateMask", updateMask)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsPatchCall) Fields(s ...googleapi.Field) *AppsPatchCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsPatchCall) Context(ctx context.Context) *AppsPatchCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsPatchCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsPatchCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.application)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId": c.appsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.patch" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsPatchCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	ret := &Operation{
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
	//   "description": "Updates the specified Application resource. You can update the following fields:\nauth_domain - Google authentication domain for controlling user access to the application.\ndefault_cookie_expiration - Cookie expiration policy for the application.",
	//   "flatPath": "v1beta/apps/{appsId}",
	//   "httpMethod": "PATCH",
	//   "id": "appengine.apps.patch",
	//   "parameterOrder": [
	//     "appsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the Application resource to update. Example: apps/myapp.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "updateMask": {
	//       "description": "Standard field mask for the set of fields to be updated.",
	//       "format": "google-fieldmask",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}",
	//   "request": {
	//     "$ref": "Application"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.repair":

type AppsRepairCall struct {
	s                        *APIService
	appsId                   string
	repairapplicationrequest *RepairApplicationRequest
	urlParams_               gensupport.URLParams
	ctx_                     context.Context
	header_                  http.Header
}

// Repair: Recreates the required App Engine features for the specified
// App Engine application, for example a Cloud Storage bucket or App
// Engine service account. Use this method if you receive an error
// message about a missing feature, for example, Error retrieving the
// App Engine service account.
func (r *AppsService) Repair(appsId string, repairapplicationrequest *RepairApplicationRequest) *AppsRepairCall {
	c := &AppsRepairCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.repairapplicationrequest = repairapplicationrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsRepairCall) Fields(s ...googleapi.Field) *AppsRepairCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsRepairCall) Context(ctx context.Context) *AppsRepairCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsRepairCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsRepairCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.repairapplicationrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}:repair")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId": c.appsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.repair" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsRepairCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	ret := &Operation{
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
	//   "description": "Recreates the required App Engine features for the specified App Engine application, for example a Cloud Storage bucket or App Engine service account. Use this method if you receive an error message about a missing feature, for example, Error retrieving the App Engine service account.",
	//   "flatPath": "v1beta/apps/{appsId}:repair",
	//   "httpMethod": "POST",
	//   "id": "appengine.apps.repair",
	//   "parameterOrder": [
	//     "appsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the application to repair. Example: apps/myapp",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}:repair",
	//   "request": {
	//     "$ref": "RepairApplicationRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.authorizedCertificates.create":

type AppsAuthorizedCertificatesCreateCall struct {
	s                     *APIService
	appsId                string
	authorizedcertificate *AuthorizedCertificate
	urlParams_            gensupport.URLParams
	ctx_                  context.Context
	header_               http.Header
}

// Create: Uploads the specified SSL certificate.
func (r *AppsAuthorizedCertificatesService) Create(appsId string, authorizedcertificate *AuthorizedCertificate) *AppsAuthorizedCertificatesCreateCall {
	c := &AppsAuthorizedCertificatesCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.authorizedcertificate = authorizedcertificate
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsAuthorizedCertificatesCreateCall) Fields(s ...googleapi.Field) *AppsAuthorizedCertificatesCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsAuthorizedCertificatesCreateCall) Context(ctx context.Context) *AppsAuthorizedCertificatesCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsAuthorizedCertificatesCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsAuthorizedCertificatesCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.authorizedcertificate)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/authorizedCertificates")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId": c.appsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.authorizedCertificates.create" call.
// Exactly one of *AuthorizedCertificate or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *AuthorizedCertificate.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *AppsAuthorizedCertificatesCreateCall) Do(opts ...googleapi.CallOption) (*AuthorizedCertificate, error) {
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
	ret := &AuthorizedCertificate{
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
	//   "description": "Uploads the specified SSL certificate.",
	//   "flatPath": "v1beta/apps/{appsId}/authorizedCertificates",
	//   "httpMethod": "POST",
	//   "id": "appengine.apps.authorizedCertificates.create",
	//   "parameterOrder": [
	//     "appsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `parent`. Name of the parent Application resource. Example: apps/myapp.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/authorizedCertificates",
	//   "request": {
	//     "$ref": "AuthorizedCertificate"
	//   },
	//   "response": {
	//     "$ref": "AuthorizedCertificate"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.authorizedCertificates.delete":

type AppsAuthorizedCertificatesDeleteCall struct {
	s                        *APIService
	appsId                   string
	authorizedCertificatesId string
	urlParams_               gensupport.URLParams
	ctx_                     context.Context
	header_                  http.Header
}

// Delete: Deletes the specified SSL certificate.
func (r *AppsAuthorizedCertificatesService) Delete(appsId string, authorizedCertificatesId string) *AppsAuthorizedCertificatesDeleteCall {
	c := &AppsAuthorizedCertificatesDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.authorizedCertificatesId = authorizedCertificatesId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsAuthorizedCertificatesDeleteCall) Fields(s ...googleapi.Field) *AppsAuthorizedCertificatesDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsAuthorizedCertificatesDeleteCall) Context(ctx context.Context) *AppsAuthorizedCertificatesDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsAuthorizedCertificatesDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsAuthorizedCertificatesDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":                   c.appsId,
		"authorizedCertificatesId": c.authorizedCertificatesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.authorizedCertificates.delete" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *AppsAuthorizedCertificatesDeleteCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
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
	//   "description": "Deletes the specified SSL certificate.",
	//   "flatPath": "v1beta/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}",
	//   "httpMethod": "DELETE",
	//   "id": "appengine.apps.authorizedCertificates.delete",
	//   "parameterOrder": [
	//     "appsId",
	//     "authorizedCertificatesId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the resource to delete. Example: apps/myapp/authorizedCertificates/12345.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "authorizedCertificatesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}",
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.authorizedCertificates.get":

type AppsAuthorizedCertificatesGetCall struct {
	s                        *APIService
	appsId                   string
	authorizedCertificatesId string
	urlParams_               gensupport.URLParams
	ifNoneMatch_             string
	ctx_                     context.Context
	header_                  http.Header
}

// Get: Gets the specified SSL certificate.
func (r *AppsAuthorizedCertificatesService) Get(appsId string, authorizedCertificatesId string) *AppsAuthorizedCertificatesGetCall {
	c := &AppsAuthorizedCertificatesGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.authorizedCertificatesId = authorizedCertificatesId
	return c
}

// View sets the optional parameter "view": Controls the set of fields
// returned in the GET response.
//
// Possible values:
//   "BASIC_CERTIFICATE"
//   "FULL_CERTIFICATE"
func (c *AppsAuthorizedCertificatesGetCall) View(view string) *AppsAuthorizedCertificatesGetCall {
	c.urlParams_.Set("view", view)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsAuthorizedCertificatesGetCall) Fields(s ...googleapi.Field) *AppsAuthorizedCertificatesGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsAuthorizedCertificatesGetCall) IfNoneMatch(entityTag string) *AppsAuthorizedCertificatesGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsAuthorizedCertificatesGetCall) Context(ctx context.Context) *AppsAuthorizedCertificatesGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsAuthorizedCertificatesGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsAuthorizedCertificatesGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":                   c.appsId,
		"authorizedCertificatesId": c.authorizedCertificatesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.authorizedCertificates.get" call.
// Exactly one of *AuthorizedCertificate or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *AuthorizedCertificate.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *AppsAuthorizedCertificatesGetCall) Do(opts ...googleapi.CallOption) (*AuthorizedCertificate, error) {
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
	ret := &AuthorizedCertificate{
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
	//   "description": "Gets the specified SSL certificate.",
	//   "flatPath": "v1beta/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.authorizedCertificates.get",
	//   "parameterOrder": [
	//     "appsId",
	//     "authorizedCertificatesId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the resource requested. Example: apps/myapp/authorizedCertificates/12345.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "authorizedCertificatesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "view": {
	//       "description": "Controls the set of fields returned in the GET response.",
	//       "enum": [
	//         "BASIC_CERTIFICATE",
	//         "FULL_CERTIFICATE"
	//       ],
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}",
	//   "response": {
	//     "$ref": "AuthorizedCertificate"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// method id "appengine.apps.authorizedCertificates.list":

type AppsAuthorizedCertificatesListCall struct {
	s            *APIService
	appsId       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists all SSL certificates the user is authorized to
// administer.
func (r *AppsAuthorizedCertificatesService) List(appsId string) *AppsAuthorizedCertificatesListCall {
	c := &AppsAuthorizedCertificatesListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum results to
// return per page.
func (c *AppsAuthorizedCertificatesListCall) PageSize(pageSize int64) *AppsAuthorizedCertificatesListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken": Continuation token
// for fetching the next page of results.
func (c *AppsAuthorizedCertificatesListCall) PageToken(pageToken string) *AppsAuthorizedCertificatesListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// View sets the optional parameter "view": Controls the set of fields
// returned in the LIST response.
//
// Possible values:
//   "BASIC_CERTIFICATE"
//   "FULL_CERTIFICATE"
func (c *AppsAuthorizedCertificatesListCall) View(view string) *AppsAuthorizedCertificatesListCall {
	c.urlParams_.Set("view", view)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsAuthorizedCertificatesListCall) Fields(s ...googleapi.Field) *AppsAuthorizedCertificatesListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsAuthorizedCertificatesListCall) IfNoneMatch(entityTag string) *AppsAuthorizedCertificatesListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsAuthorizedCertificatesListCall) Context(ctx context.Context) *AppsAuthorizedCertificatesListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsAuthorizedCertificatesListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsAuthorizedCertificatesListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/authorizedCertificates")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId": c.appsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.authorizedCertificates.list" call.
// Exactly one of *ListAuthorizedCertificatesResponse or error will be
// non-nil. Any non-2xx status code is an error. Response headers are in
// either *ListAuthorizedCertificatesResponse.ServerResponse.Header or
// (if a response was returned at all) in
// error.(*googleapi.Error).Header. Use googleapi.IsNotModified to check
// whether the returned error was because http.StatusNotModified was
// returned.
func (c *AppsAuthorizedCertificatesListCall) Do(opts ...googleapi.CallOption) (*ListAuthorizedCertificatesResponse, error) {
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
	ret := &ListAuthorizedCertificatesResponse{
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
	//   "description": "Lists all SSL certificates the user is authorized to administer.",
	//   "flatPath": "v1beta/apps/{appsId}/authorizedCertificates",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.authorizedCertificates.list",
	//   "parameterOrder": [
	//     "appsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `parent`. Name of the parent Application resource. Example: apps/myapp.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum results to return per page.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "Continuation token for fetching the next page of results.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "view": {
	//       "description": "Controls the set of fields returned in the LIST response.",
	//       "enum": [
	//         "BASIC_CERTIFICATE",
	//         "FULL_CERTIFICATE"
	//       ],
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/authorizedCertificates",
	//   "response": {
	//     "$ref": "ListAuthorizedCertificatesResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *AppsAuthorizedCertificatesListCall) Pages(ctx context.Context, f func(*ListAuthorizedCertificatesResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "appengine.apps.authorizedCertificates.patch":

type AppsAuthorizedCertificatesPatchCall struct {
	s                        *APIService
	appsId                   string
	authorizedCertificatesId string
	authorizedcertificate    *AuthorizedCertificate
	urlParams_               gensupport.URLParams
	ctx_                     context.Context
	header_                  http.Header
}

// Patch: Updates the specified SSL certificate. To renew a certificate
// and maintain its existing domain mappings, update certificate_data
// with a new certificate. The new certificate must be applicable to the
// same domains as the original certificate. The certificate
// display_name may also be updated.
func (r *AppsAuthorizedCertificatesService) Patch(appsId string, authorizedCertificatesId string, authorizedcertificate *AuthorizedCertificate) *AppsAuthorizedCertificatesPatchCall {
	c := &AppsAuthorizedCertificatesPatchCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.authorizedCertificatesId = authorizedCertificatesId
	c.authorizedcertificate = authorizedcertificate
	return c
}

// UpdateMask sets the optional parameter "updateMask": Standard field
// mask for the set of fields to be updated. Updates are only supported
// on the certificate_raw_data and display_name fields.
func (c *AppsAuthorizedCertificatesPatchCall) UpdateMask(updateMask string) *AppsAuthorizedCertificatesPatchCall {
	c.urlParams_.Set("updateMask", updateMask)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsAuthorizedCertificatesPatchCall) Fields(s ...googleapi.Field) *AppsAuthorizedCertificatesPatchCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsAuthorizedCertificatesPatchCall) Context(ctx context.Context) *AppsAuthorizedCertificatesPatchCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsAuthorizedCertificatesPatchCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsAuthorizedCertificatesPatchCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.authorizedcertificate)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":                   c.appsId,
		"authorizedCertificatesId": c.authorizedCertificatesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.authorizedCertificates.patch" call.
// Exactly one of *AuthorizedCertificate or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *AuthorizedCertificate.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *AppsAuthorizedCertificatesPatchCall) Do(opts ...googleapi.CallOption) (*AuthorizedCertificate, error) {
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
	ret := &AuthorizedCertificate{
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
	//   "description": "Updates the specified SSL certificate. To renew a certificate and maintain its existing domain mappings, update certificate_data with a new certificate. The new certificate must be applicable to the same domains as the original certificate. The certificate display_name may also be updated.",
	//   "flatPath": "v1beta/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}",
	//   "httpMethod": "PATCH",
	//   "id": "appengine.apps.authorizedCertificates.patch",
	//   "parameterOrder": [
	//     "appsId",
	//     "authorizedCertificatesId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the resource to update. Example: apps/myapp/authorizedCertificates/12345.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "authorizedCertificatesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "updateMask": {
	//       "description": "Standard field mask for the set of fields to be updated. Updates are only supported on the certificate_raw_data and display_name fields.",
	//       "format": "google-fieldmask",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}",
	//   "request": {
	//     "$ref": "AuthorizedCertificate"
	//   },
	//   "response": {
	//     "$ref": "AuthorizedCertificate"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.authorizedDomains.list":

type AppsAuthorizedDomainsListCall struct {
	s            *APIService
	appsId       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists all domains the user is authorized to administer.
func (r *AppsAuthorizedDomainsService) List(appsId string) *AppsAuthorizedDomainsListCall {
	c := &AppsAuthorizedDomainsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum results to
// return per page.
func (c *AppsAuthorizedDomainsListCall) PageSize(pageSize int64) *AppsAuthorizedDomainsListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken": Continuation token
// for fetching the next page of results.
func (c *AppsAuthorizedDomainsListCall) PageToken(pageToken string) *AppsAuthorizedDomainsListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsAuthorizedDomainsListCall) Fields(s ...googleapi.Field) *AppsAuthorizedDomainsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsAuthorizedDomainsListCall) IfNoneMatch(entityTag string) *AppsAuthorizedDomainsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsAuthorizedDomainsListCall) Context(ctx context.Context) *AppsAuthorizedDomainsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsAuthorizedDomainsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsAuthorizedDomainsListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/authorizedDomains")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId": c.appsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.authorizedDomains.list" call.
// Exactly one of *ListAuthorizedDomainsResponse or error will be
// non-nil. Any non-2xx status code is an error. Response headers are in
// either *ListAuthorizedDomainsResponse.ServerResponse.Header or (if a
// response was returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *AppsAuthorizedDomainsListCall) Do(opts ...googleapi.CallOption) (*ListAuthorizedDomainsResponse, error) {
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
	ret := &ListAuthorizedDomainsResponse{
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
	//   "description": "Lists all domains the user is authorized to administer.",
	//   "flatPath": "v1beta/apps/{appsId}/authorizedDomains",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.authorizedDomains.list",
	//   "parameterOrder": [
	//     "appsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `parent`. Name of the parent Application resource. Example: apps/myapp.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum results to return per page.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "Continuation token for fetching the next page of results.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/authorizedDomains",
	//   "response": {
	//     "$ref": "ListAuthorizedDomainsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *AppsAuthorizedDomainsListCall) Pages(ctx context.Context, f func(*ListAuthorizedDomainsResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "appengine.apps.domainMappings.create":

type AppsDomainMappingsCreateCall struct {
	s             *APIService
	appsId        string
	domainmapping *DomainMapping
	urlParams_    gensupport.URLParams
	ctx_          context.Context
	header_       http.Header
}

// Create: Maps a domain to an application. A user must be authorized to
// administer a domain in order to map it to an application. For a list
// of available authorized domains, see
// AuthorizedDomains.ListAuthorizedDomains.
func (r *AppsDomainMappingsService) Create(appsId string, domainmapping *DomainMapping) *AppsDomainMappingsCreateCall {
	c := &AppsDomainMappingsCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.domainmapping = domainmapping
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsDomainMappingsCreateCall) Fields(s ...googleapi.Field) *AppsDomainMappingsCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsDomainMappingsCreateCall) Context(ctx context.Context) *AppsDomainMappingsCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsDomainMappingsCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsDomainMappingsCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.domainmapping)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/domainMappings")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId": c.appsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.domainMappings.create" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsDomainMappingsCreateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	ret := &Operation{
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
	//   "description": "Maps a domain to an application. A user must be authorized to administer a domain in order to map it to an application. For a list of available authorized domains, see AuthorizedDomains.ListAuthorizedDomains.",
	//   "flatPath": "v1beta/apps/{appsId}/domainMappings",
	//   "httpMethod": "POST",
	//   "id": "appengine.apps.domainMappings.create",
	//   "parameterOrder": [
	//     "appsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `parent`. Name of the parent Application resource. Example: apps/myapp.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/domainMappings",
	//   "request": {
	//     "$ref": "DomainMapping"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.domainMappings.delete":

type AppsDomainMappingsDeleteCall struct {
	s                *APIService
	appsId           string
	domainMappingsId string
	urlParams_       gensupport.URLParams
	ctx_             context.Context
	header_          http.Header
}

// Delete: Deletes the specified domain mapping. A user must be
// authorized to administer the associated domain in order to delete a
// DomainMapping resource.
func (r *AppsDomainMappingsService) Delete(appsId string, domainMappingsId string) *AppsDomainMappingsDeleteCall {
	c := &AppsDomainMappingsDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.domainMappingsId = domainMappingsId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsDomainMappingsDeleteCall) Fields(s ...googleapi.Field) *AppsDomainMappingsDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsDomainMappingsDeleteCall) Context(ctx context.Context) *AppsDomainMappingsDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsDomainMappingsDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsDomainMappingsDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/domainMappings/{domainMappingsId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":           c.appsId,
		"domainMappingsId": c.domainMappingsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.domainMappings.delete" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsDomainMappingsDeleteCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	ret := &Operation{
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
	//   "description": "Deletes the specified domain mapping. A user must be authorized to administer the associated domain in order to delete a DomainMapping resource.",
	//   "flatPath": "v1beta/apps/{appsId}/domainMappings/{domainMappingsId}",
	//   "httpMethod": "DELETE",
	//   "id": "appengine.apps.domainMappings.delete",
	//   "parameterOrder": [
	//     "appsId",
	//     "domainMappingsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the resource to delete. Example: apps/myapp/domainMappings/example.com.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "domainMappingsId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/domainMappings/{domainMappingsId}",
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.domainMappings.get":

type AppsDomainMappingsGetCall struct {
	s                *APIService
	appsId           string
	domainMappingsId string
	urlParams_       gensupport.URLParams
	ifNoneMatch_     string
	ctx_             context.Context
	header_          http.Header
}

// Get: Gets the specified domain mapping.
func (r *AppsDomainMappingsService) Get(appsId string, domainMappingsId string) *AppsDomainMappingsGetCall {
	c := &AppsDomainMappingsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.domainMappingsId = domainMappingsId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsDomainMappingsGetCall) Fields(s ...googleapi.Field) *AppsDomainMappingsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsDomainMappingsGetCall) IfNoneMatch(entityTag string) *AppsDomainMappingsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsDomainMappingsGetCall) Context(ctx context.Context) *AppsDomainMappingsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsDomainMappingsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsDomainMappingsGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/domainMappings/{domainMappingsId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":           c.appsId,
		"domainMappingsId": c.domainMappingsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.domainMappings.get" call.
// Exactly one of *DomainMapping or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *DomainMapping.ServerResponse.Header or (if a response was returned
// at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *AppsDomainMappingsGetCall) Do(opts ...googleapi.CallOption) (*DomainMapping, error) {
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
	ret := &DomainMapping{
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
	//   "description": "Gets the specified domain mapping.",
	//   "flatPath": "v1beta/apps/{appsId}/domainMappings/{domainMappingsId}",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.domainMappings.get",
	//   "parameterOrder": [
	//     "appsId",
	//     "domainMappingsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the resource requested. Example: apps/myapp/domainMappings/example.com.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "domainMappingsId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/domainMappings/{domainMappingsId}",
	//   "response": {
	//     "$ref": "DomainMapping"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// method id "appengine.apps.domainMappings.list":

type AppsDomainMappingsListCall struct {
	s            *APIService
	appsId       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists the domain mappings on an application.
func (r *AppsDomainMappingsService) List(appsId string) *AppsDomainMappingsListCall {
	c := &AppsDomainMappingsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum results to
// return per page.
func (c *AppsDomainMappingsListCall) PageSize(pageSize int64) *AppsDomainMappingsListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken": Continuation token
// for fetching the next page of results.
func (c *AppsDomainMappingsListCall) PageToken(pageToken string) *AppsDomainMappingsListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsDomainMappingsListCall) Fields(s ...googleapi.Field) *AppsDomainMappingsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsDomainMappingsListCall) IfNoneMatch(entityTag string) *AppsDomainMappingsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsDomainMappingsListCall) Context(ctx context.Context) *AppsDomainMappingsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsDomainMappingsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsDomainMappingsListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/domainMappings")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId": c.appsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.domainMappings.list" call.
// Exactly one of *ListDomainMappingsResponse or error will be non-nil.
// Any non-2xx status code is an error. Response headers are in either
// *ListDomainMappingsResponse.ServerResponse.Header or (if a response
// was returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *AppsDomainMappingsListCall) Do(opts ...googleapi.CallOption) (*ListDomainMappingsResponse, error) {
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
	ret := &ListDomainMappingsResponse{
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
	//   "description": "Lists the domain mappings on an application.",
	//   "flatPath": "v1beta/apps/{appsId}/domainMappings",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.domainMappings.list",
	//   "parameterOrder": [
	//     "appsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `parent`. Name of the parent Application resource. Example: apps/myapp.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum results to return per page.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "Continuation token for fetching the next page of results.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/domainMappings",
	//   "response": {
	//     "$ref": "ListDomainMappingsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *AppsDomainMappingsListCall) Pages(ctx context.Context, f func(*ListDomainMappingsResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "appengine.apps.domainMappings.patch":

type AppsDomainMappingsPatchCall struct {
	s                *APIService
	appsId           string
	domainMappingsId string
	domainmapping    *DomainMapping
	urlParams_       gensupport.URLParams
	ctx_             context.Context
	header_          http.Header
}

// Patch: Updates the specified domain mapping. To map an SSL
// certificate to a domain mapping, update certificate_id to point to an
// AuthorizedCertificate resource. A user must be authorized to
// administer the associated domain in order to update a DomainMapping
// resource.
func (r *AppsDomainMappingsService) Patch(appsId string, domainMappingsId string, domainmapping *DomainMapping) *AppsDomainMappingsPatchCall {
	c := &AppsDomainMappingsPatchCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.domainMappingsId = domainMappingsId
	c.domainmapping = domainmapping
	return c
}

// UpdateMask sets the optional parameter "updateMask": Standard field
// mask for the set of fields to be updated.
func (c *AppsDomainMappingsPatchCall) UpdateMask(updateMask string) *AppsDomainMappingsPatchCall {
	c.urlParams_.Set("updateMask", updateMask)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsDomainMappingsPatchCall) Fields(s ...googleapi.Field) *AppsDomainMappingsPatchCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsDomainMappingsPatchCall) Context(ctx context.Context) *AppsDomainMappingsPatchCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsDomainMappingsPatchCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsDomainMappingsPatchCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.domainmapping)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/domainMappings/{domainMappingsId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":           c.appsId,
		"domainMappingsId": c.domainMappingsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.domainMappings.patch" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsDomainMappingsPatchCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	ret := &Operation{
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
	//   "description": "Updates the specified domain mapping. To map an SSL certificate to a domain mapping, update certificate_id to point to an AuthorizedCertificate resource. A user must be authorized to administer the associated domain in order to update a DomainMapping resource.",
	//   "flatPath": "v1beta/apps/{appsId}/domainMappings/{domainMappingsId}",
	//   "httpMethod": "PATCH",
	//   "id": "appengine.apps.domainMappings.patch",
	//   "parameterOrder": [
	//     "appsId",
	//     "domainMappingsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the resource to update. Example: apps/myapp/domainMappings/example.com.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "domainMappingsId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "updateMask": {
	//       "description": "Standard field mask for the set of fields to be updated.",
	//       "format": "google-fieldmask",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/domainMappings/{domainMappingsId}",
	//   "request": {
	//     "$ref": "DomainMapping"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.firewall.ingressRules.batchUpdate":

type AppsFirewallIngressRulesBatchUpdateCall struct {
	s                              *APIService
	appsId                         string
	batchupdateingressrulesrequest *BatchUpdateIngressRulesRequest
	urlParams_                     gensupport.URLParams
	ctx_                           context.Context
	header_                        http.Header
}

// BatchUpdate: Replaces the entire firewall ruleset in one bulk
// operation. This overrides and replaces the rules of an existing
// firewall with the new rules.If the final rule does not match traffic
// with the '*' wildcard IP range, then an "allow all" rule is
// explicitly added to the end of the list.
func (r *AppsFirewallIngressRulesService) BatchUpdate(appsId string, batchupdateingressrulesrequest *BatchUpdateIngressRulesRequest) *AppsFirewallIngressRulesBatchUpdateCall {
	c := &AppsFirewallIngressRulesBatchUpdateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.batchupdateingressrulesrequest = batchupdateingressrulesrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsFirewallIngressRulesBatchUpdateCall) Fields(s ...googleapi.Field) *AppsFirewallIngressRulesBatchUpdateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsFirewallIngressRulesBatchUpdateCall) Context(ctx context.Context) *AppsFirewallIngressRulesBatchUpdateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsFirewallIngressRulesBatchUpdateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsFirewallIngressRulesBatchUpdateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.batchupdateingressrulesrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/firewall/ingressRules:batchUpdate")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId": c.appsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.firewall.ingressRules.batchUpdate" call.
// Exactly one of *BatchUpdateIngressRulesResponse or error will be
// non-nil. Any non-2xx status code is an error. Response headers are in
// either *BatchUpdateIngressRulesResponse.ServerResponse.Header or (if
// a response was returned at all) in error.(*googleapi.Error).Header.
// Use googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *AppsFirewallIngressRulesBatchUpdateCall) Do(opts ...googleapi.CallOption) (*BatchUpdateIngressRulesResponse, error) {
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
	ret := &BatchUpdateIngressRulesResponse{
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
	//   "description": "Replaces the entire firewall ruleset in one bulk operation. This overrides and replaces the rules of an existing firewall with the new rules.If the final rule does not match traffic with the '*' wildcard IP range, then an \"allow all\" rule is explicitly added to the end of the list.",
	//   "flatPath": "v1beta/apps/{appsId}/firewall/ingressRules:batchUpdate",
	//   "httpMethod": "POST",
	//   "id": "appengine.apps.firewall.ingressRules.batchUpdate",
	//   "parameterOrder": [
	//     "appsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the Firewall collection to set. Example: apps/myapp/firewall/ingressRules.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/firewall/ingressRules:batchUpdate",
	//   "request": {
	//     "$ref": "BatchUpdateIngressRulesRequest"
	//   },
	//   "response": {
	//     "$ref": "BatchUpdateIngressRulesResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.firewall.ingressRules.create":

type AppsFirewallIngressRulesCreateCall struct {
	s            *APIService
	appsId       string
	firewallrule *FirewallRule
	urlParams_   gensupport.URLParams
	ctx_         context.Context
	header_      http.Header
}

// Create: Creates a firewall rule for the application.
func (r *AppsFirewallIngressRulesService) Create(appsId string, firewallrule *FirewallRule) *AppsFirewallIngressRulesCreateCall {
	c := &AppsFirewallIngressRulesCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.firewallrule = firewallrule
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsFirewallIngressRulesCreateCall) Fields(s ...googleapi.Field) *AppsFirewallIngressRulesCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsFirewallIngressRulesCreateCall) Context(ctx context.Context) *AppsFirewallIngressRulesCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsFirewallIngressRulesCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsFirewallIngressRulesCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.firewallrule)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/firewall/ingressRules")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId": c.appsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.firewall.ingressRules.create" call.
// Exactly one of *FirewallRule or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *FirewallRule.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsFirewallIngressRulesCreateCall) Do(opts ...googleapi.CallOption) (*FirewallRule, error) {
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
	ret := &FirewallRule{
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
	//   "description": "Creates a firewall rule for the application.",
	//   "flatPath": "v1beta/apps/{appsId}/firewall/ingressRules",
	//   "httpMethod": "POST",
	//   "id": "appengine.apps.firewall.ingressRules.create",
	//   "parameterOrder": [
	//     "appsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `parent`. Name of the parent Firewall collection in which to create a new rule. Example: apps/myapp/firewall/ingressRules.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/firewall/ingressRules",
	//   "request": {
	//     "$ref": "FirewallRule"
	//   },
	//   "response": {
	//     "$ref": "FirewallRule"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.firewall.ingressRules.delete":

type AppsFirewallIngressRulesDeleteCall struct {
	s              *APIService
	appsId         string
	ingressRulesId string
	urlParams_     gensupport.URLParams
	ctx_           context.Context
	header_        http.Header
}

// Delete: Deletes the specified firewall rule.
func (r *AppsFirewallIngressRulesService) Delete(appsId string, ingressRulesId string) *AppsFirewallIngressRulesDeleteCall {
	c := &AppsFirewallIngressRulesDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.ingressRulesId = ingressRulesId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsFirewallIngressRulesDeleteCall) Fields(s ...googleapi.Field) *AppsFirewallIngressRulesDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsFirewallIngressRulesDeleteCall) Context(ctx context.Context) *AppsFirewallIngressRulesDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsFirewallIngressRulesDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsFirewallIngressRulesDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/firewall/ingressRules/{ingressRulesId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":         c.appsId,
		"ingressRulesId": c.ingressRulesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.firewall.ingressRules.delete" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *AppsFirewallIngressRulesDeleteCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
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
	//   "description": "Deletes the specified firewall rule.",
	//   "flatPath": "v1beta/apps/{appsId}/firewall/ingressRules/{ingressRulesId}",
	//   "httpMethod": "DELETE",
	//   "id": "appengine.apps.firewall.ingressRules.delete",
	//   "parameterOrder": [
	//     "appsId",
	//     "ingressRulesId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the Firewall resource to delete. Example: apps/myapp/firewall/ingressRules/100.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "ingressRulesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/firewall/ingressRules/{ingressRulesId}",
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.firewall.ingressRules.get":

type AppsFirewallIngressRulesGetCall struct {
	s              *APIService
	appsId         string
	ingressRulesId string
	urlParams_     gensupport.URLParams
	ifNoneMatch_   string
	ctx_           context.Context
	header_        http.Header
}

// Get: Gets the specified firewall rule.
func (r *AppsFirewallIngressRulesService) Get(appsId string, ingressRulesId string) *AppsFirewallIngressRulesGetCall {
	c := &AppsFirewallIngressRulesGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.ingressRulesId = ingressRulesId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsFirewallIngressRulesGetCall) Fields(s ...googleapi.Field) *AppsFirewallIngressRulesGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsFirewallIngressRulesGetCall) IfNoneMatch(entityTag string) *AppsFirewallIngressRulesGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsFirewallIngressRulesGetCall) Context(ctx context.Context) *AppsFirewallIngressRulesGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsFirewallIngressRulesGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsFirewallIngressRulesGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/firewall/ingressRules/{ingressRulesId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":         c.appsId,
		"ingressRulesId": c.ingressRulesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.firewall.ingressRules.get" call.
// Exactly one of *FirewallRule or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *FirewallRule.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsFirewallIngressRulesGetCall) Do(opts ...googleapi.CallOption) (*FirewallRule, error) {
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
	ret := &FirewallRule{
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
	//   "description": "Gets the specified firewall rule.",
	//   "flatPath": "v1beta/apps/{appsId}/firewall/ingressRules/{ingressRulesId}",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.firewall.ingressRules.get",
	//   "parameterOrder": [
	//     "appsId",
	//     "ingressRulesId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the Firewall resource to retrieve. Example: apps/myapp/firewall/ingressRules/100.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "ingressRulesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/firewall/ingressRules/{ingressRulesId}",
	//   "response": {
	//     "$ref": "FirewallRule"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// method id "appengine.apps.firewall.ingressRules.list":

type AppsFirewallIngressRulesListCall struct {
	s            *APIService
	appsId       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists the firewall rules of an application.
func (r *AppsFirewallIngressRulesService) List(appsId string) *AppsFirewallIngressRulesListCall {
	c := &AppsFirewallIngressRulesListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	return c
}

// MatchingAddress sets the optional parameter "matchingAddress": A
// valid IP Address. If set, only rules matching this address will be
// returned. The first returned rule will be the rule that fires on
// requests from this IP.
func (c *AppsFirewallIngressRulesListCall) MatchingAddress(matchingAddress string) *AppsFirewallIngressRulesListCall {
	c.urlParams_.Set("matchingAddress", matchingAddress)
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum results to
// return per page.
func (c *AppsFirewallIngressRulesListCall) PageSize(pageSize int64) *AppsFirewallIngressRulesListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken": Continuation token
// for fetching the next page of results.
func (c *AppsFirewallIngressRulesListCall) PageToken(pageToken string) *AppsFirewallIngressRulesListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsFirewallIngressRulesListCall) Fields(s ...googleapi.Field) *AppsFirewallIngressRulesListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsFirewallIngressRulesListCall) IfNoneMatch(entityTag string) *AppsFirewallIngressRulesListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsFirewallIngressRulesListCall) Context(ctx context.Context) *AppsFirewallIngressRulesListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsFirewallIngressRulesListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsFirewallIngressRulesListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/firewall/ingressRules")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId": c.appsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.firewall.ingressRules.list" call.
// Exactly one of *ListIngressRulesResponse or error will be non-nil.
// Any non-2xx status code is an error. Response headers are in either
// *ListIngressRulesResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *AppsFirewallIngressRulesListCall) Do(opts ...googleapi.CallOption) (*ListIngressRulesResponse, error) {
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
	ret := &ListIngressRulesResponse{
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
	//   "description": "Lists the firewall rules of an application.",
	//   "flatPath": "v1beta/apps/{appsId}/firewall/ingressRules",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.firewall.ingressRules.list",
	//   "parameterOrder": [
	//     "appsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `parent`. Name of the Firewall collection to retrieve. Example: apps/myapp/firewall/ingressRules.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "matchingAddress": {
	//       "description": "A valid IP Address. If set, only rules matching this address will be returned. The first returned rule will be the rule that fires on requests from this IP.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum results to return per page.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "Continuation token for fetching the next page of results.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/firewall/ingressRules",
	//   "response": {
	//     "$ref": "ListIngressRulesResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *AppsFirewallIngressRulesListCall) Pages(ctx context.Context, f func(*ListIngressRulesResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "appengine.apps.firewall.ingressRules.patch":

type AppsFirewallIngressRulesPatchCall struct {
	s              *APIService
	appsId         string
	ingressRulesId string
	firewallrule   *FirewallRule
	urlParams_     gensupport.URLParams
	ctx_           context.Context
	header_        http.Header
}

// Patch: Updates the specified firewall rule.
func (r *AppsFirewallIngressRulesService) Patch(appsId string, ingressRulesId string, firewallrule *FirewallRule) *AppsFirewallIngressRulesPatchCall {
	c := &AppsFirewallIngressRulesPatchCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.ingressRulesId = ingressRulesId
	c.firewallrule = firewallrule
	return c
}

// UpdateMask sets the optional parameter "updateMask": Standard field
// mask for the set of fields to be updated.
func (c *AppsFirewallIngressRulesPatchCall) UpdateMask(updateMask string) *AppsFirewallIngressRulesPatchCall {
	c.urlParams_.Set("updateMask", updateMask)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsFirewallIngressRulesPatchCall) Fields(s ...googleapi.Field) *AppsFirewallIngressRulesPatchCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsFirewallIngressRulesPatchCall) Context(ctx context.Context) *AppsFirewallIngressRulesPatchCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsFirewallIngressRulesPatchCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsFirewallIngressRulesPatchCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.firewallrule)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/firewall/ingressRules/{ingressRulesId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":         c.appsId,
		"ingressRulesId": c.ingressRulesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.firewall.ingressRules.patch" call.
// Exactly one of *FirewallRule or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *FirewallRule.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsFirewallIngressRulesPatchCall) Do(opts ...googleapi.CallOption) (*FirewallRule, error) {
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
	ret := &FirewallRule{
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
	//   "description": "Updates the specified firewall rule.",
	//   "flatPath": "v1beta/apps/{appsId}/firewall/ingressRules/{ingressRulesId}",
	//   "httpMethod": "PATCH",
	//   "id": "appengine.apps.firewall.ingressRules.patch",
	//   "parameterOrder": [
	//     "appsId",
	//     "ingressRulesId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the Firewall resource to update. Example: apps/myapp/firewall/ingressRules/100.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "ingressRulesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "updateMask": {
	//       "description": "Standard field mask for the set of fields to be updated.",
	//       "format": "google-fieldmask",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/firewall/ingressRules/{ingressRulesId}",
	//   "request": {
	//     "$ref": "FirewallRule"
	//   },
	//   "response": {
	//     "$ref": "FirewallRule"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.locations.get":

type AppsLocationsGetCall struct {
	s            *APIService
	appsId       string
	locationsId  string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Get information about a location.
func (r *AppsLocationsService) Get(appsId string, locationsId string) *AppsLocationsGetCall {
	c := &AppsLocationsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.locationsId = locationsId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsLocationsGetCall) Fields(s ...googleapi.Field) *AppsLocationsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsLocationsGetCall) IfNoneMatch(entityTag string) *AppsLocationsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsLocationsGetCall) Context(ctx context.Context) *AppsLocationsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsLocationsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsLocationsGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/locations/{locationsId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":      c.appsId,
		"locationsId": c.locationsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.locations.get" call.
// Exactly one of *Location or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Location.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsLocationsGetCall) Do(opts ...googleapi.CallOption) (*Location, error) {
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
	ret := &Location{
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
	//   "description": "Get information about a location.",
	//   "flatPath": "v1beta/apps/{appsId}/locations/{locationsId}",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.locations.get",
	//   "parameterOrder": [
	//     "appsId",
	//     "locationsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Resource name for the location.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "locationsId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/locations/{locationsId}",
	//   "response": {
	//     "$ref": "Location"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// method id "appengine.apps.locations.list":

type AppsLocationsListCall struct {
	s            *APIService
	appsId       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists information about the supported locations for this
// service.
func (r *AppsLocationsService) List(appsId string) *AppsLocationsListCall {
	c := &AppsLocationsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	return c
}

// Filter sets the optional parameter "filter": The standard list
// filter.
func (c *AppsLocationsListCall) Filter(filter string) *AppsLocationsListCall {
	c.urlParams_.Set("filter", filter)
	return c
}

// PageSize sets the optional parameter "pageSize": The standard list
// page size.
func (c *AppsLocationsListCall) PageSize(pageSize int64) *AppsLocationsListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken": The standard list
// page token.
func (c *AppsLocationsListCall) PageToken(pageToken string) *AppsLocationsListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsLocationsListCall) Fields(s ...googleapi.Field) *AppsLocationsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsLocationsListCall) IfNoneMatch(entityTag string) *AppsLocationsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsLocationsListCall) Context(ctx context.Context) *AppsLocationsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsLocationsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsLocationsListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/locations")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId": c.appsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.locations.list" call.
// Exactly one of *ListLocationsResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListLocationsResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *AppsLocationsListCall) Do(opts ...googleapi.CallOption) (*ListLocationsResponse, error) {
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
	ret := &ListLocationsResponse{
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
	//   "description": "Lists information about the supported locations for this service.",
	//   "flatPath": "v1beta/apps/{appsId}/locations",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.locations.list",
	//   "parameterOrder": [
	//     "appsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. The resource that owns the locations collection, if applicable.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "filter": {
	//       "description": "The standard list filter.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "The standard list page size.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "The standard list page token.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/locations",
	//   "response": {
	//     "$ref": "ListLocationsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *AppsLocationsListCall) Pages(ctx context.Context, f func(*ListLocationsResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "appengine.apps.operations.get":

type AppsOperationsGetCall struct {
	s            *APIService
	appsId       string
	operationsId string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Gets the latest state of a long-running operation. Clients can
// use this method to poll the operation result at intervals as
// recommended by the API service.
func (r *AppsOperationsService) Get(appsId string, operationsId string) *AppsOperationsGetCall {
	c := &AppsOperationsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.operationsId = operationsId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsOperationsGetCall) Fields(s ...googleapi.Field) *AppsOperationsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsOperationsGetCall) IfNoneMatch(entityTag string) *AppsOperationsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsOperationsGetCall) Context(ctx context.Context) *AppsOperationsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsOperationsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsOperationsGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/operations/{operationsId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":       c.appsId,
		"operationsId": c.operationsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.operations.get" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsOperationsGetCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	ret := &Operation{
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
	//   "description": "Gets the latest state of a long-running operation. Clients can use this method to poll the operation result at intervals as recommended by the API service.",
	//   "flatPath": "v1beta/apps/{appsId}/operations/{operationsId}",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.operations.get",
	//   "parameterOrder": [
	//     "appsId",
	//     "operationsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. The name of the operation resource.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "operationsId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/operations/{operationsId}",
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// method id "appengine.apps.operations.list":

type AppsOperationsListCall struct {
	s            *APIService
	appsId       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists operations that match the specified filter in the
// request. If the server doesn't support this method, it returns
// UNIMPLEMENTED.NOTE: the name binding allows API services to override
// the binding to use different resource name schemes, such as
// users/*/operations. To override the binding, API services can add a
// binding such as "/v1/{name=users/*}/operations" to their service
// configuration. For backwards compatibility, the default name includes
// the operations collection id, however overriding users must ensure
// the name binding is the parent resource, without the operations
// collection id.
func (r *AppsOperationsService) List(appsId string) *AppsOperationsListCall {
	c := &AppsOperationsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	return c
}

// Filter sets the optional parameter "filter": The standard list
// filter.
func (c *AppsOperationsListCall) Filter(filter string) *AppsOperationsListCall {
	c.urlParams_.Set("filter", filter)
	return c
}

// PageSize sets the optional parameter "pageSize": The standard list
// page size.
func (c *AppsOperationsListCall) PageSize(pageSize int64) *AppsOperationsListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken": The standard list
// page token.
func (c *AppsOperationsListCall) PageToken(pageToken string) *AppsOperationsListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsOperationsListCall) Fields(s ...googleapi.Field) *AppsOperationsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsOperationsListCall) IfNoneMatch(entityTag string) *AppsOperationsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsOperationsListCall) Context(ctx context.Context) *AppsOperationsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsOperationsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsOperationsListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/operations")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId": c.appsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.operations.list" call.
// Exactly one of *ListOperationsResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListOperationsResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *AppsOperationsListCall) Do(opts ...googleapi.CallOption) (*ListOperationsResponse, error) {
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
	ret := &ListOperationsResponse{
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
	//   "description": "Lists operations that match the specified filter in the request. If the server doesn't support this method, it returns UNIMPLEMENTED.NOTE: the name binding allows API services to override the binding to use different resource name schemes, such as users/*/operations. To override the binding, API services can add a binding such as \"/v1/{name=users/*}/operations\" to their service configuration. For backwards compatibility, the default name includes the operations collection id, however overriding users must ensure the name binding is the parent resource, without the operations collection id.",
	//   "flatPath": "v1beta/apps/{appsId}/operations",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.operations.list",
	//   "parameterOrder": [
	//     "appsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. The name of the operation's parent resource.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "filter": {
	//       "description": "The standard list filter.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "The standard list page size.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "The standard list page token.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/operations",
	//   "response": {
	//     "$ref": "ListOperationsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *AppsOperationsListCall) Pages(ctx context.Context, f func(*ListOperationsResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "appengine.apps.services.delete":

type AppsServicesDeleteCall struct {
	s          *APIService
	appsId     string
	servicesId string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Delete: Deletes the specified service and all enclosed versions.
func (r *AppsServicesService) Delete(appsId string, servicesId string) *AppsServicesDeleteCall {
	c := &AppsServicesDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.servicesId = servicesId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsServicesDeleteCall) Fields(s ...googleapi.Field) *AppsServicesDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsServicesDeleteCall) Context(ctx context.Context) *AppsServicesDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsServicesDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsServicesDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/services/{servicesId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":     c.appsId,
		"servicesId": c.servicesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.services.delete" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsServicesDeleteCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	ret := &Operation{
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
	//   "description": "Deletes the specified service and all enclosed versions.",
	//   "flatPath": "v1beta/apps/{appsId}/services/{servicesId}",
	//   "httpMethod": "DELETE",
	//   "id": "appengine.apps.services.delete",
	//   "parameterOrder": [
	//     "appsId",
	//     "servicesId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the resource requested. Example: apps/myapp/services/default.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "servicesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/services/{servicesId}",
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.services.get":

type AppsServicesGetCall struct {
	s            *APIService
	appsId       string
	servicesId   string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Gets the current configuration of the specified service.
func (r *AppsServicesService) Get(appsId string, servicesId string) *AppsServicesGetCall {
	c := &AppsServicesGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.servicesId = servicesId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsServicesGetCall) Fields(s ...googleapi.Field) *AppsServicesGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsServicesGetCall) IfNoneMatch(entityTag string) *AppsServicesGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsServicesGetCall) Context(ctx context.Context) *AppsServicesGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsServicesGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsServicesGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/services/{servicesId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":     c.appsId,
		"servicesId": c.servicesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.services.get" call.
// Exactly one of *Service or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Service.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *AppsServicesGetCall) Do(opts ...googleapi.CallOption) (*Service, error) {
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
	ret := &Service{
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
	//   "description": "Gets the current configuration of the specified service.",
	//   "flatPath": "v1beta/apps/{appsId}/services/{servicesId}",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.services.get",
	//   "parameterOrder": [
	//     "appsId",
	//     "servicesId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the resource requested. Example: apps/myapp/services/default.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "servicesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/services/{servicesId}",
	//   "response": {
	//     "$ref": "Service"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// method id "appengine.apps.services.list":

type AppsServicesListCall struct {
	s            *APIService
	appsId       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists all the services in the application.
func (r *AppsServicesService) List(appsId string) *AppsServicesListCall {
	c := &AppsServicesListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum results to
// return per page.
func (c *AppsServicesListCall) PageSize(pageSize int64) *AppsServicesListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken": Continuation token
// for fetching the next page of results.
func (c *AppsServicesListCall) PageToken(pageToken string) *AppsServicesListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsServicesListCall) Fields(s ...googleapi.Field) *AppsServicesListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsServicesListCall) IfNoneMatch(entityTag string) *AppsServicesListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsServicesListCall) Context(ctx context.Context) *AppsServicesListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsServicesListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsServicesListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/services")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId": c.appsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.services.list" call.
// Exactly one of *ListServicesResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListServicesResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *AppsServicesListCall) Do(opts ...googleapi.CallOption) (*ListServicesResponse, error) {
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
	ret := &ListServicesResponse{
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
	//   "description": "Lists all the services in the application.",
	//   "flatPath": "v1beta/apps/{appsId}/services",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.services.list",
	//   "parameterOrder": [
	//     "appsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `parent`. Name of the parent Application resource. Example: apps/myapp.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum results to return per page.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "Continuation token for fetching the next page of results.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/services",
	//   "response": {
	//     "$ref": "ListServicesResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *AppsServicesListCall) Pages(ctx context.Context, f func(*ListServicesResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "appengine.apps.services.patch":

type AppsServicesPatchCall struct {
	s          *APIService
	appsId     string
	servicesId string
	service    *Service
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Patch: Updates the configuration of the specified service.
func (r *AppsServicesService) Patch(appsId string, servicesId string, service *Service) *AppsServicesPatchCall {
	c := &AppsServicesPatchCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.servicesId = servicesId
	c.service = service
	return c
}

// MigrateTraffic sets the optional parameter "migrateTraffic": Set to
// true to gradually shift traffic to one or more versions that you
// specify. By default, traffic is shifted immediately. For gradual
// traffic migration, the target versions must be located within
// instances that are configured for both warmup requests
// (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1be
// ta/apps.services.versions#inboundservicetype) and automatic scaling
// (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1be
// ta/apps.services.versions#automaticscaling). You must specify the
// shardBy
// (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1be
// ta/apps.services#shardby) field in the Service resource. Gradual
// traffic migration is not supported in the App Engine flexible
// environment. For examples, see Migrating and Splitting Traffic
// (https://cloud.google.com/appengine/docs/admin-api/migrating-splitting
// -traffic).
func (c *AppsServicesPatchCall) MigrateTraffic(migrateTraffic bool) *AppsServicesPatchCall {
	c.urlParams_.Set("migrateTraffic", fmt.Sprint(migrateTraffic))
	return c
}

// UpdateMask sets the optional parameter "updateMask": Standard field
// mask for the set of fields to be updated.
func (c *AppsServicesPatchCall) UpdateMask(updateMask string) *AppsServicesPatchCall {
	c.urlParams_.Set("updateMask", updateMask)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsServicesPatchCall) Fields(s ...googleapi.Field) *AppsServicesPatchCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsServicesPatchCall) Context(ctx context.Context) *AppsServicesPatchCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsServicesPatchCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsServicesPatchCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.service)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/services/{servicesId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":     c.appsId,
		"servicesId": c.servicesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.services.patch" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsServicesPatchCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	ret := &Operation{
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
	//   "description": "Updates the configuration of the specified service.",
	//   "flatPath": "v1beta/apps/{appsId}/services/{servicesId}",
	//   "httpMethod": "PATCH",
	//   "id": "appengine.apps.services.patch",
	//   "parameterOrder": [
	//     "appsId",
	//     "servicesId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the resource to update. Example: apps/myapp/services/default.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "migrateTraffic": {
	//       "description": "Set to true to gradually shift traffic to one or more versions that you specify. By default, traffic is shifted immediately. For gradual traffic migration, the target versions must be located within instances that are configured for both warmup requests (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1beta/apps.services.versions#inboundservicetype) and automatic scaling (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1beta/apps.services.versions#automaticscaling). You must specify the shardBy (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1beta/apps.services#shardby) field in the Service resource. Gradual traffic migration is not supported in the App Engine flexible environment. For examples, see Migrating and Splitting Traffic (https://cloud.google.com/appengine/docs/admin-api/migrating-splitting-traffic).",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "servicesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "updateMask": {
	//       "description": "Standard field mask for the set of fields to be updated.",
	//       "format": "google-fieldmask",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/services/{servicesId}",
	//   "request": {
	//     "$ref": "Service"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.services.versions.create":

type AppsServicesVersionsCreateCall struct {
	s          *APIService
	appsId     string
	servicesId string
	version    *Version
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Create: Deploys code and resource files to a new version.
func (r *AppsServicesVersionsService) Create(appsId string, servicesId string, version *Version) *AppsServicesVersionsCreateCall {
	c := &AppsServicesVersionsCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.servicesId = servicesId
	c.version = version
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsServicesVersionsCreateCall) Fields(s ...googleapi.Field) *AppsServicesVersionsCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsServicesVersionsCreateCall) Context(ctx context.Context) *AppsServicesVersionsCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsServicesVersionsCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsServicesVersionsCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.version)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/services/{servicesId}/versions")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":     c.appsId,
		"servicesId": c.servicesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.services.versions.create" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsServicesVersionsCreateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	ret := &Operation{
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
	//   "description": "Deploys code and resource files to a new version.",
	//   "flatPath": "v1beta/apps/{appsId}/services/{servicesId}/versions",
	//   "httpMethod": "POST",
	//   "id": "appengine.apps.services.versions.create",
	//   "parameterOrder": [
	//     "appsId",
	//     "servicesId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `parent`. Name of the parent resource to create this version under. Example: apps/myapp/services/default.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "servicesId": {
	//       "description": "Part of `parent`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/services/{servicesId}/versions",
	//   "request": {
	//     "$ref": "Version"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.services.versions.delete":

type AppsServicesVersionsDeleteCall struct {
	s          *APIService
	appsId     string
	servicesId string
	versionsId string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Delete: Deletes an existing Version resource.
func (r *AppsServicesVersionsService) Delete(appsId string, servicesId string, versionsId string) *AppsServicesVersionsDeleteCall {
	c := &AppsServicesVersionsDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.servicesId = servicesId
	c.versionsId = versionsId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsServicesVersionsDeleteCall) Fields(s ...googleapi.Field) *AppsServicesVersionsDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsServicesVersionsDeleteCall) Context(ctx context.Context) *AppsServicesVersionsDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsServicesVersionsDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsServicesVersionsDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":     c.appsId,
		"servicesId": c.servicesId,
		"versionsId": c.versionsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.services.versions.delete" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsServicesVersionsDeleteCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	ret := &Operation{
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
	//   "description": "Deletes an existing Version resource.",
	//   "flatPath": "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}",
	//   "httpMethod": "DELETE",
	//   "id": "appengine.apps.services.versions.delete",
	//   "parameterOrder": [
	//     "appsId",
	//     "servicesId",
	//     "versionsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the resource requested. Example: apps/myapp/services/default/versions/v1.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "servicesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "versionsId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}",
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.services.versions.get":

type AppsServicesVersionsGetCall struct {
	s            *APIService
	appsId       string
	servicesId   string
	versionsId   string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Gets the specified Version resource. By default, only a
// BASIC_VIEW will be returned. Specify the FULL_VIEW parameter to get
// the full resource.
func (r *AppsServicesVersionsService) Get(appsId string, servicesId string, versionsId string) *AppsServicesVersionsGetCall {
	c := &AppsServicesVersionsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.servicesId = servicesId
	c.versionsId = versionsId
	return c
}

// View sets the optional parameter "view": Controls the set of fields
// returned in the Get response.
//
// Possible values:
//   "BASIC"
//   "FULL"
func (c *AppsServicesVersionsGetCall) View(view string) *AppsServicesVersionsGetCall {
	c.urlParams_.Set("view", view)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsServicesVersionsGetCall) Fields(s ...googleapi.Field) *AppsServicesVersionsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsServicesVersionsGetCall) IfNoneMatch(entityTag string) *AppsServicesVersionsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsServicesVersionsGetCall) Context(ctx context.Context) *AppsServicesVersionsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsServicesVersionsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsServicesVersionsGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":     c.appsId,
		"servicesId": c.servicesId,
		"versionsId": c.versionsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.services.versions.get" call.
// Exactly one of *Version or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Version.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *AppsServicesVersionsGetCall) Do(opts ...googleapi.CallOption) (*Version, error) {
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
	ret := &Version{
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
	//   "description": "Gets the specified Version resource. By default, only a BASIC_VIEW will be returned. Specify the FULL_VIEW parameter to get the full resource.",
	//   "flatPath": "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.services.versions.get",
	//   "parameterOrder": [
	//     "appsId",
	//     "servicesId",
	//     "versionsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the resource requested. Example: apps/myapp/services/default/versions/v1.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "servicesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "versionsId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "view": {
	//       "description": "Controls the set of fields returned in the Get response.",
	//       "enum": [
	//         "BASIC",
	//         "FULL"
	//       ],
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}",
	//   "response": {
	//     "$ref": "Version"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// method id "appengine.apps.services.versions.list":

type AppsServicesVersionsListCall struct {
	s            *APIService
	appsId       string
	servicesId   string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists the versions of a service.
func (r *AppsServicesVersionsService) List(appsId string, servicesId string) *AppsServicesVersionsListCall {
	c := &AppsServicesVersionsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.servicesId = servicesId
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum results to
// return per page.
func (c *AppsServicesVersionsListCall) PageSize(pageSize int64) *AppsServicesVersionsListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken": Continuation token
// for fetching the next page of results.
func (c *AppsServicesVersionsListCall) PageToken(pageToken string) *AppsServicesVersionsListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// View sets the optional parameter "view": Controls the set of fields
// returned in the List response.
//
// Possible values:
//   "BASIC"
//   "FULL"
func (c *AppsServicesVersionsListCall) View(view string) *AppsServicesVersionsListCall {
	c.urlParams_.Set("view", view)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsServicesVersionsListCall) Fields(s ...googleapi.Field) *AppsServicesVersionsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsServicesVersionsListCall) IfNoneMatch(entityTag string) *AppsServicesVersionsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsServicesVersionsListCall) Context(ctx context.Context) *AppsServicesVersionsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsServicesVersionsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsServicesVersionsListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/services/{servicesId}/versions")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":     c.appsId,
		"servicesId": c.servicesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.services.versions.list" call.
// Exactly one of *ListVersionsResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListVersionsResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *AppsServicesVersionsListCall) Do(opts ...googleapi.CallOption) (*ListVersionsResponse, error) {
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
	ret := &ListVersionsResponse{
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
	//   "description": "Lists the versions of a service.",
	//   "flatPath": "v1beta/apps/{appsId}/services/{servicesId}/versions",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.services.versions.list",
	//   "parameterOrder": [
	//     "appsId",
	//     "servicesId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `parent`. Name of the parent Service resource. Example: apps/myapp/services/default.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum results to return per page.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "Continuation token for fetching the next page of results.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "servicesId": {
	//       "description": "Part of `parent`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "view": {
	//       "description": "Controls the set of fields returned in the List response.",
	//       "enum": [
	//         "BASIC",
	//         "FULL"
	//       ],
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/services/{servicesId}/versions",
	//   "response": {
	//     "$ref": "ListVersionsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *AppsServicesVersionsListCall) Pages(ctx context.Context, f func(*ListVersionsResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "appengine.apps.services.versions.patch":

type AppsServicesVersionsPatchCall struct {
	s          *APIService
	appsId     string
	servicesId string
	versionsId string
	version    *Version
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Patch: Updates the specified Version resource. You can specify the
// following fields depending on the App Engine environment and type of
// scaling that the version resource uses:
// serving_status
// (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1be
// ta/apps.services.versions#Version.FIELDS.serving_status):  For
// Version resources that use basic scaling, manual scaling, or run in
// the App Engine flexible environment.
// instance_class
// (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1be
// ta/apps.services.versions#Version.FIELDS.instance_class):  For
// Version resources that run in the App Engine standard
// environment.
// automatic_scaling.min_idle_instances
// (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1be
// ta/apps.services.versions#Version.FIELDS.automatic_scaling):  For
// Version resources that use automatic scaling and run in the App
// Engine standard environment.
// automatic_scaling.max_idle_instances
// (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1be
// ta/apps.services.versions#Version.FIELDS.automatic_scaling):  For
// Version resources that use automatic scaling and run in the App
// Engine standard environment.
// automatic_scaling.min_total_instances
// (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1be
// ta/apps.services.versions#Version.FIELDS.automatic_scaling):  For
// Version resources that use automatic scaling and run in the App
// Engine Flexible environment.
// automatic_scaling.max_total_instances
// (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1be
// ta/apps.services.versions#Version.FIELDS.automatic_scaling):  For
// Version resources that use automatic scaling and run in the App
// Engine Flexible environment.
// automatic_scaling.cool_down_period_sec
// (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1be
// ta/apps.services.versions#Version.FIELDS.automatic_scaling):  For
// Version resources that use automatic scaling and run in the App
// Engine Flexible
// environment.
// automatic_scaling.cpu_utilization.target_utilization
// (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1be
// ta/apps.services.versions#Version.FIELDS.automatic_scaling):  For
// Version resources that use automatic scaling and run in the App
// Engine Flexible environment.
func (r *AppsServicesVersionsService) Patch(appsId string, servicesId string, versionsId string, version *Version) *AppsServicesVersionsPatchCall {
	c := &AppsServicesVersionsPatchCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.servicesId = servicesId
	c.versionsId = versionsId
	c.version = version
	return c
}

// UpdateMask sets the optional parameter "updateMask": Standard field
// mask for the set of fields to be updated.
func (c *AppsServicesVersionsPatchCall) UpdateMask(updateMask string) *AppsServicesVersionsPatchCall {
	c.urlParams_.Set("updateMask", updateMask)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsServicesVersionsPatchCall) Fields(s ...googleapi.Field) *AppsServicesVersionsPatchCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsServicesVersionsPatchCall) Context(ctx context.Context) *AppsServicesVersionsPatchCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsServicesVersionsPatchCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsServicesVersionsPatchCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.version)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":     c.appsId,
		"servicesId": c.servicesId,
		"versionsId": c.versionsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.services.versions.patch" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsServicesVersionsPatchCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	ret := &Operation{
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
	//   "description": "Updates the specified Version resource. You can specify the following fields depending on the App Engine environment and type of scaling that the version resource uses:\nserving_status (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1beta/apps.services.versions#Version.FIELDS.serving_status):  For Version resources that use basic scaling, manual scaling, or run in  the App Engine flexible environment.\ninstance_class (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1beta/apps.services.versions#Version.FIELDS.instance_class):  For Version resources that run in the App Engine standard environment.\nautomatic_scaling.min_idle_instances (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1beta/apps.services.versions#Version.FIELDS.automatic_scaling):  For Version resources that use automatic scaling and run in the App  Engine standard environment.\nautomatic_scaling.max_idle_instances (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1beta/apps.services.versions#Version.FIELDS.automatic_scaling):  For Version resources that use automatic scaling and run in the App  Engine standard environment.\nautomatic_scaling.min_total_instances (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1beta/apps.services.versions#Version.FIELDS.automatic_scaling):  For Version resources that use automatic scaling and run in the App  Engine Flexible environment.\nautomatic_scaling.max_total_instances (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1beta/apps.services.versions#Version.FIELDS.automatic_scaling):  For Version resources that use automatic scaling and run in the App  Engine Flexible environment.\nautomatic_scaling.cool_down_period_sec (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1beta/apps.services.versions#Version.FIELDS.automatic_scaling):  For Version resources that use automatic scaling and run in the App  Engine Flexible environment.\nautomatic_scaling.cpu_utilization.target_utilization (https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1beta/apps.services.versions#Version.FIELDS.automatic_scaling):  For Version resources that use automatic scaling and run in the App  Engine Flexible environment.",
	//   "flatPath": "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}",
	//   "httpMethod": "PATCH",
	//   "id": "appengine.apps.services.versions.patch",
	//   "parameterOrder": [
	//     "appsId",
	//     "servicesId",
	//     "versionsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the resource to update. Example: apps/myapp/services/default/versions/1.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "servicesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "updateMask": {
	//       "description": "Standard field mask for the set of fields to be updated.",
	//       "format": "google-fieldmask",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "versionsId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}",
	//   "request": {
	//     "$ref": "Version"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.services.versions.instances.debug":

type AppsServicesVersionsInstancesDebugCall struct {
	s                    *APIService
	appsId               string
	servicesId           string
	versionsId           string
	instancesId          string
	debuginstancerequest *DebugInstanceRequest
	urlParams_           gensupport.URLParams
	ctx_                 context.Context
	header_              http.Header
}

// Debug: Enables debugging on a VM instance. This allows you to use the
// SSH command to connect to the virtual machine where the instance
// lives. While in "debug mode", the instance continues to serve live
// traffic. You should delete the instance when you are done debugging
// and then allow the system to take over and determine if another
// instance should be started.Only applicable for instances in App
// Engine flexible environment.
func (r *AppsServicesVersionsInstancesService) Debug(appsId string, servicesId string, versionsId string, instancesId string, debuginstancerequest *DebugInstanceRequest) *AppsServicesVersionsInstancesDebugCall {
	c := &AppsServicesVersionsInstancesDebugCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.servicesId = servicesId
	c.versionsId = versionsId
	c.instancesId = instancesId
	c.debuginstancerequest = debuginstancerequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsServicesVersionsInstancesDebugCall) Fields(s ...googleapi.Field) *AppsServicesVersionsInstancesDebugCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsServicesVersionsInstancesDebugCall) Context(ctx context.Context) *AppsServicesVersionsInstancesDebugCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsServicesVersionsInstancesDebugCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsServicesVersionsInstancesDebugCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.debuginstancerequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}/instances/{instancesId}:debug")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":      c.appsId,
		"servicesId":  c.servicesId,
		"versionsId":  c.versionsId,
		"instancesId": c.instancesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.services.versions.instances.debug" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsServicesVersionsInstancesDebugCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	ret := &Operation{
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
	//   "description": "Enables debugging on a VM instance. This allows you to use the SSH command to connect to the virtual machine where the instance lives. While in \"debug mode\", the instance continues to serve live traffic. You should delete the instance when you are done debugging and then allow the system to take over and determine if another instance should be started.Only applicable for instances in App Engine flexible environment.",
	//   "flatPath": "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}/instances/{instancesId}:debug",
	//   "httpMethod": "POST",
	//   "id": "appengine.apps.services.versions.instances.debug",
	//   "parameterOrder": [
	//     "appsId",
	//     "servicesId",
	//     "versionsId",
	//     "instancesId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the resource requested. Example: apps/myapp/services/default/versions/v1/instances/instance-1.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "instancesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "servicesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "versionsId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}/instances/{instancesId}:debug",
	//   "request": {
	//     "$ref": "DebugInstanceRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.services.versions.instances.delete":

type AppsServicesVersionsInstancesDeleteCall struct {
	s           *APIService
	appsId      string
	servicesId  string
	versionsId  string
	instancesId string
	urlParams_  gensupport.URLParams
	ctx_        context.Context
	header_     http.Header
}

// Delete: Stops a running instance.
func (r *AppsServicesVersionsInstancesService) Delete(appsId string, servicesId string, versionsId string, instancesId string) *AppsServicesVersionsInstancesDeleteCall {
	c := &AppsServicesVersionsInstancesDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.servicesId = servicesId
	c.versionsId = versionsId
	c.instancesId = instancesId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsServicesVersionsInstancesDeleteCall) Fields(s ...googleapi.Field) *AppsServicesVersionsInstancesDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsServicesVersionsInstancesDeleteCall) Context(ctx context.Context) *AppsServicesVersionsInstancesDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsServicesVersionsInstancesDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsServicesVersionsInstancesDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}/instances/{instancesId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":      c.appsId,
		"servicesId":  c.servicesId,
		"versionsId":  c.versionsId,
		"instancesId": c.instancesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.services.versions.instances.delete" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsServicesVersionsInstancesDeleteCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	ret := &Operation{
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
	//   "description": "Stops a running instance.",
	//   "flatPath": "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}/instances/{instancesId}",
	//   "httpMethod": "DELETE",
	//   "id": "appengine.apps.services.versions.instances.delete",
	//   "parameterOrder": [
	//     "appsId",
	//     "servicesId",
	//     "versionsId",
	//     "instancesId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the resource requested. Example: apps/myapp/services/default/versions/v1/instances/instance-1.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "instancesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "servicesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "versionsId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}/instances/{instancesId}",
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "appengine.apps.services.versions.instances.get":

type AppsServicesVersionsInstancesGetCall struct {
	s            *APIService
	appsId       string
	servicesId   string
	versionsId   string
	instancesId  string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Gets instance information.
func (r *AppsServicesVersionsInstancesService) Get(appsId string, servicesId string, versionsId string, instancesId string) *AppsServicesVersionsInstancesGetCall {
	c := &AppsServicesVersionsInstancesGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.servicesId = servicesId
	c.versionsId = versionsId
	c.instancesId = instancesId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsServicesVersionsInstancesGetCall) Fields(s ...googleapi.Field) *AppsServicesVersionsInstancesGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsServicesVersionsInstancesGetCall) IfNoneMatch(entityTag string) *AppsServicesVersionsInstancesGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsServicesVersionsInstancesGetCall) Context(ctx context.Context) *AppsServicesVersionsInstancesGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsServicesVersionsInstancesGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsServicesVersionsInstancesGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}/instances/{instancesId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":      c.appsId,
		"servicesId":  c.servicesId,
		"versionsId":  c.versionsId,
		"instancesId": c.instancesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.services.versions.instances.get" call.
// Exactly one of *Instance or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Instance.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *AppsServicesVersionsInstancesGetCall) Do(opts ...googleapi.CallOption) (*Instance, error) {
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
	ret := &Instance{
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
	//   "description": "Gets instance information.",
	//   "flatPath": "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}/instances/{instancesId}",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.services.versions.instances.get",
	//   "parameterOrder": [
	//     "appsId",
	//     "servicesId",
	//     "versionsId",
	//     "instancesId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `name`. Name of the resource requested. Example: apps/myapp/services/default/versions/v1/instances/instance-1.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "instancesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "servicesId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "versionsId": {
	//       "description": "Part of `name`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}/instances/{instancesId}",
	//   "response": {
	//     "$ref": "Instance"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// method id "appengine.apps.services.versions.instances.list":

type AppsServicesVersionsInstancesListCall struct {
	s            *APIService
	appsId       string
	servicesId   string
	versionsId   string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists the instances of a version.Tip: To aggregate details
// about instances over time, see the Stackdriver Monitoring API
// (https://cloud.google.com/monitoring/api/ref_v3/rest/v3/projects.timeS
// eries/list).
func (r *AppsServicesVersionsInstancesService) List(appsId string, servicesId string, versionsId string) *AppsServicesVersionsInstancesListCall {
	c := &AppsServicesVersionsInstancesListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.appsId = appsId
	c.servicesId = servicesId
	c.versionsId = versionsId
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum results to
// return per page.
func (c *AppsServicesVersionsInstancesListCall) PageSize(pageSize int64) *AppsServicesVersionsInstancesListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken": Continuation token
// for fetching the next page of results.
func (c *AppsServicesVersionsInstancesListCall) PageToken(pageToken string) *AppsServicesVersionsInstancesListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AppsServicesVersionsInstancesListCall) Fields(s ...googleapi.Field) *AppsServicesVersionsInstancesListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AppsServicesVersionsInstancesListCall) IfNoneMatch(entityTag string) *AppsServicesVersionsInstancesListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AppsServicesVersionsInstancesListCall) Context(ctx context.Context) *AppsServicesVersionsInstancesListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AppsServicesVersionsInstancesListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AppsServicesVersionsInstancesListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}/instances")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"appsId":     c.appsId,
		"servicesId": c.servicesId,
		"versionsId": c.versionsId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "appengine.apps.services.versions.instances.list" call.
// Exactly one of *ListInstancesResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListInstancesResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *AppsServicesVersionsInstancesListCall) Do(opts ...googleapi.CallOption) (*ListInstancesResponse, error) {
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
	ret := &ListInstancesResponse{
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
	//   "description": "Lists the instances of a version.Tip: To aggregate details about instances over time, see the Stackdriver Monitoring API (https://cloud.google.com/monitoring/api/ref_v3/rest/v3/projects.timeSeries/list).",
	//   "flatPath": "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}/instances",
	//   "httpMethod": "GET",
	//   "id": "appengine.apps.services.versions.instances.list",
	//   "parameterOrder": [
	//     "appsId",
	//     "servicesId",
	//     "versionsId"
	//   ],
	//   "parameters": {
	//     "appsId": {
	//       "description": "Part of `parent`. Name of the parent Version resource. Example: apps/myapp/services/default/versions/v1.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum results to return per page.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "Continuation token for fetching the next page of results.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "servicesId": {
	//       "description": "Part of `parent`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "versionsId": {
	//       "description": "Part of `parent`. See documentation of `appsId`.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta/apps/{appsId}/services/{servicesId}/versions/{versionsId}/instances",
	//   "response": {
	//     "$ref": "ListInstancesResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/appengine.admin",
	//     "https://www.googleapis.com/auth/cloud-platform",
	//     "https://www.googleapis.com/auth/cloud-platform.read-only"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *AppsServicesVersionsInstancesListCall) Pages(ctx context.Context, f func(*ListInstancesResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}
