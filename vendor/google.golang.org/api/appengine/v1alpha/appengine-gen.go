// Package appengine provides access to the Google App Engine Admin API.
//
// See https://cloud.google.com/appengine/docs/admin-api/
//
// Usage example:
//
//   import "google.golang.org/api/appengine/v1alpha"
//   ...
//   appengineService, err := appengine.New(oauthHttpClient)
package appengine // import "google.golang.org/api/appengine/v1alpha"

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

const apiId = "appengine:v1alpha"
const apiName = "appengine"
const apiVersion = "v1alpha"
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
	rs.Locations = NewAppsLocationsService(s)
	rs.Operations = NewAppsOperationsService(s)
	return rs
}

type AppsService struct {
	s *APIService

	AuthorizedCertificates *AppsAuthorizedCertificatesService

	AuthorizedDomains *AppsAuthorizedDomainsService

	DomainMappings *AppsDomainMappingsService

	Locations *AppsLocationsService

	Operations *AppsOperationsService
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
	//   "UNSPECIFIED_STATUS"
	//   "OK" - Certificate was successfully obtained and inserted into the
	// serving system.
	//   "PENDING" - Certificate is under active attempts to acquire or
	// renew.
	//   "FAILED_RETRYING_INTERNAL" - Most recent renewal failed due to a
	// system failure and will be retried. System failure is likely
	// transient, and subsequent renewal attempts may succeed. The last
	// successfully provisioned certificate may still be serving.
	//   "FAILED_RETRYING_NOT_VISIBLE" - Most recent renewal failed due to
	// an invalid DNS setup and will be retried. Renewal attempts will
	// continue to fail until the certificate domain's DNS configuration is
	// fixed. The last successfully provisioned certificate may still be
	// serving.
	//   "FAILED_PERMANENTLY_NOT_VISIBLE" - All renewal attempts have been
	// exhausted. Most recent renewal failed due to an invalid DNS setup and
	// will not be retried. The last successfully provisioned certificate
	// may still be serving.
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

// SslSettings: SSL configuration for a DomainMapping resource.
type SslSettings struct {
	// CertificateId: ID of the AuthorizedCertificate resource configuring
	// SSL for the application. Clearing this field will remove SSL
	// support.By default, a managed certificate is automatically created
	// for every domain mapping. To omit SSL support or to configure SSL
	// manually, specify no_managed_certificate on a CREATE or UPDATE
	// request. You must be authorized to administer the
	// AuthorizedCertificate resource to manually map it to a DomainMapping
	// resource. Example: 12345.
	CertificateId string `json:"certificateId,omitempty"`

	// IsManagedCertificate: Whether the mapped certificate is an App Engine
	// managed certificate. Managed certificates are created by default with
	// a domain mapping. To opt out, specify no_managed_certificate on a
	// CREATE or UPDATE request.@OutputOnly
	IsManagedCertificate bool `json:"isManagedCertificate,omitempty"`

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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1alpha/apps/{appsId}/authorizedCertificates")
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
	//   "flatPath": "v1alpha/apps/{appsId}/authorizedCertificates",
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
	//   "path": "v1alpha/apps/{appsId}/authorizedCertificates",
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1alpha/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}")
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
	//   "flatPath": "v1alpha/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}",
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
	//   "path": "v1alpha/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}",
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1alpha/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}")
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
	//   "flatPath": "v1alpha/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}",
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
	//   "path": "v1alpha/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}",
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1alpha/apps/{appsId}/authorizedCertificates")
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
	//   "flatPath": "v1alpha/apps/{appsId}/authorizedCertificates",
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
	//   "path": "v1alpha/apps/{appsId}/authorizedCertificates",
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1alpha/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}")
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
	//   "flatPath": "v1alpha/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}",
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
	//   "path": "v1alpha/apps/{appsId}/authorizedCertificates/{authorizedCertificatesId}",
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1alpha/apps/{appsId}/authorizedDomains")
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
	//   "flatPath": "v1alpha/apps/{appsId}/authorizedDomains",
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
	//   "path": "v1alpha/apps/{appsId}/authorizedDomains",
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

// NoManagedCertificate sets the optional parameter
// "noManagedCertificate": Whether a managed certificate should be
// provided by App Engine. If true, a certificate ID must be manaually
// set in the DomainMapping resource to configure SSL for this domain.
// If false, a managed certificate will be provisioned and a certificate
// ID will be automatically populated.
func (c *AppsDomainMappingsCreateCall) NoManagedCertificate(noManagedCertificate bool) *AppsDomainMappingsCreateCall {
	c.urlParams_.Set("noManagedCertificate", fmt.Sprint(noManagedCertificate))
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1alpha/apps/{appsId}/domainMappings")
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
	//   "flatPath": "v1alpha/apps/{appsId}/domainMappings",
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
	//     },
	//     "noManagedCertificate": {
	//       "description": "Whether a managed certificate should be provided by App Engine. If true, a certificate ID must be manaually set in the DomainMapping resource to configure SSL for this domain. If false, a managed certificate will be provisioned and a certificate ID will be automatically populated.",
	//       "location": "query",
	//       "type": "boolean"
	//     }
	//   },
	//   "path": "v1alpha/apps/{appsId}/domainMappings",
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1alpha/apps/{appsId}/domainMappings/{domainMappingsId}")
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
	//   "flatPath": "v1alpha/apps/{appsId}/domainMappings/{domainMappingsId}",
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
	//   "path": "v1alpha/apps/{appsId}/domainMappings/{domainMappingsId}",
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1alpha/apps/{appsId}/domainMappings/{domainMappingsId}")
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
	//   "flatPath": "v1alpha/apps/{appsId}/domainMappings/{domainMappingsId}",
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
	//   "path": "v1alpha/apps/{appsId}/domainMappings/{domainMappingsId}",
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1alpha/apps/{appsId}/domainMappings")
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
	//   "flatPath": "v1alpha/apps/{appsId}/domainMappings",
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
	//   "path": "v1alpha/apps/{appsId}/domainMappings",
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

// NoManagedCertificate sets the optional parameter
// "noManagedCertificate": Whether a managed certificate should be
// provided by App Engine. If true, a certificate ID must be manually
// set in the DomainMapping resource to configure SSL for this domain.
// If false, a managed certificate will be provisioned and a certificate
// ID will be automatically populated. Only applicable if
// ssl_settings.certificate_id is specified in the update mask.
func (c *AppsDomainMappingsPatchCall) NoManagedCertificate(noManagedCertificate bool) *AppsDomainMappingsPatchCall {
	c.urlParams_.Set("noManagedCertificate", fmt.Sprint(noManagedCertificate))
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1alpha/apps/{appsId}/domainMappings/{domainMappingsId}")
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
	//   "flatPath": "v1alpha/apps/{appsId}/domainMappings/{domainMappingsId}",
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
	//     "noManagedCertificate": {
	//       "description": "Whether a managed certificate should be provided by App Engine. If true, a certificate ID must be manually set in the DomainMapping resource to configure SSL for this domain. If false, a managed certificate will be provisioned and a certificate ID will be automatically populated. Only applicable if ssl_settings.certificate_id is specified in the update mask.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "updateMask": {
	//       "description": "Standard field mask for the set of fields to be updated.",
	//       "format": "google-fieldmask",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1alpha/apps/{appsId}/domainMappings/{domainMappingsId}",
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1alpha/apps/{appsId}/locations/{locationsId}")
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
	//   "flatPath": "v1alpha/apps/{appsId}/locations/{locationsId}",
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
	//   "path": "v1alpha/apps/{appsId}/locations/{locationsId}",
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1alpha/apps/{appsId}/locations")
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
	//   "flatPath": "v1alpha/apps/{appsId}/locations",
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
	//   "path": "v1alpha/apps/{appsId}/locations",
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1alpha/apps/{appsId}/operations/{operationsId}")
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
	//   "flatPath": "v1alpha/apps/{appsId}/operations/{operationsId}",
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
	//   "path": "v1alpha/apps/{appsId}/operations/{operationsId}",
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1alpha/apps/{appsId}/operations")
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
	//   "flatPath": "v1alpha/apps/{appsId}/operations",
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
	//   "path": "v1alpha/apps/{appsId}/operations",
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
