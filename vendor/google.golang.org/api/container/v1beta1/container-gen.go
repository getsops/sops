// Package container provides access to the Google Container Engine API.
//
// See https://cloud.google.com/container-engine/
//
// Usage example:
//
//   import "google.golang.org/api/container/v1beta1"
//   ...
//   containerService, err := container.New(oauthHttpClient)
package container // import "google.golang.org/api/container/v1beta1"

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

const apiId = "container:v1beta1"
const apiName = "container"
const apiVersion = "v1beta1"
const basePath = "https://container.googleapis.com/"

// OAuth2 scopes used by this API.
const (
	// View and manage your data across Google Cloud Platform services
	CloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"
)

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
	rs.Locations = NewProjectsLocationsService(s)
	rs.Zones = NewProjectsZonesService(s)
	return rs
}

type ProjectsService struct {
	s *Service

	Locations *ProjectsLocationsService

	Zones *ProjectsZonesService
}

func NewProjectsLocationsService(s *Service) *ProjectsLocationsService {
	rs := &ProjectsLocationsService{s: s}
	rs.Clusters = NewProjectsLocationsClustersService(s)
	rs.Operations = NewProjectsLocationsOperationsService(s)
	return rs
}

type ProjectsLocationsService struct {
	s *Service

	Clusters *ProjectsLocationsClustersService

	Operations *ProjectsLocationsOperationsService
}

func NewProjectsLocationsClustersService(s *Service) *ProjectsLocationsClustersService {
	rs := &ProjectsLocationsClustersService{s: s}
	rs.NodePools = NewProjectsLocationsClustersNodePoolsService(s)
	return rs
}

type ProjectsLocationsClustersService struct {
	s *Service

	NodePools *ProjectsLocationsClustersNodePoolsService
}

func NewProjectsLocationsClustersNodePoolsService(s *Service) *ProjectsLocationsClustersNodePoolsService {
	rs := &ProjectsLocationsClustersNodePoolsService{s: s}
	return rs
}

type ProjectsLocationsClustersNodePoolsService struct {
	s *Service
}

func NewProjectsLocationsOperationsService(s *Service) *ProjectsLocationsOperationsService {
	rs := &ProjectsLocationsOperationsService{s: s}
	return rs
}

type ProjectsLocationsOperationsService struct {
	s *Service
}

func NewProjectsZonesService(s *Service) *ProjectsZonesService {
	rs := &ProjectsZonesService{s: s}
	rs.Clusters = NewProjectsZonesClustersService(s)
	rs.Operations = NewProjectsZonesOperationsService(s)
	return rs
}

type ProjectsZonesService struct {
	s *Service

	Clusters *ProjectsZonesClustersService

	Operations *ProjectsZonesOperationsService
}

func NewProjectsZonesClustersService(s *Service) *ProjectsZonesClustersService {
	rs := &ProjectsZonesClustersService{s: s}
	rs.NodePools = NewProjectsZonesClustersNodePoolsService(s)
	return rs
}

type ProjectsZonesClustersService struct {
	s *Service

	NodePools *ProjectsZonesClustersNodePoolsService
}

func NewProjectsZonesClustersNodePoolsService(s *Service) *ProjectsZonesClustersNodePoolsService {
	rs := &ProjectsZonesClustersNodePoolsService{s: s}
	return rs
}

type ProjectsZonesClustersNodePoolsService struct {
	s *Service
}

func NewProjectsZonesOperationsService(s *Service) *ProjectsZonesOperationsService {
	rs := &ProjectsZonesOperationsService{s: s}
	return rs
}

type ProjectsZonesOperationsService struct {
	s *Service
}

// AcceleratorConfig: AcceleratorConfig represents a Hardware
// Accelerator request.
type AcceleratorConfig struct {
	// AcceleratorCount: The number of the accelerator cards exposed to an
	// instance.
	AcceleratorCount int64 `json:"acceleratorCount,omitempty,string"`

	// AcceleratorType: The accelerator type resource name. List of
	// supported accelerators
	// [here](/compute/docs/gpus/#Introduction)
	AcceleratorType string `json:"acceleratorType,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AcceleratorCount") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AcceleratorCount") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *AcceleratorConfig) MarshalJSON() ([]byte, error) {
	type noMethod AcceleratorConfig
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// AddonsConfig: Configuration for the addons that can be automatically
// spun up in the
// cluster, enabling additional functionality.
type AddonsConfig struct {
	// HorizontalPodAutoscaling: Configuration for the horizontal pod
	// autoscaling feature, which
	// increases or decreases the number of replica pods a replication
	// controller
	// has based on the resource usage of the existing pods.
	HorizontalPodAutoscaling *HorizontalPodAutoscaling `json:"horizontalPodAutoscaling,omitempty"`

	// HttpLoadBalancing: Configuration for the HTTP (L7) load balancing
	// controller addon, which
	// makes it easy to set up HTTP load balancers for services in a
	// cluster.
	HttpLoadBalancing *HttpLoadBalancing `json:"httpLoadBalancing,omitempty"`

	// KubernetesDashboard: Configuration for the Kubernetes Dashboard.
	KubernetesDashboard *KubernetesDashboard `json:"kubernetesDashboard,omitempty"`

	// NetworkPolicyConfig: Configuration for NetworkPolicy. This only
	// tracks whether the addon
	// is enabled or not on the Master, it does not track whether network
	// policy
	// is enabled for the nodes.
	NetworkPolicyConfig *NetworkPolicyConfig `json:"networkPolicyConfig,omitempty"`

	// ForceSendFields is a list of field names (e.g.
	// "HorizontalPodAutoscaling") to unconditionally include in API
	// requests. By default, fields with empty values are omitted from API
	// requests. However, any non-pointer, non-interface field appearing in
	// ForceSendFields will be sent to the server regardless of whether the
	// field is empty or not. This may be used to include empty fields in
	// Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "HorizontalPodAutoscaling")
	// to include in API requests with the JSON null value. By default,
	// fields with empty values are omitted from API requests. However, any
	// field with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *AddonsConfig) MarshalJSON() ([]byte, error) {
	type noMethod AddonsConfig
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// AutoUpgradeOptions: AutoUpgradeOptions defines the set of options for
// the user to control how
// the Auto Upgrades will proceed.
type AutoUpgradeOptions struct {
	// AutoUpgradeStartTime: [Output only] This field is set when upgrades
	// are about to commence
	// with the approximate start time for the upgrades,
	// in
	// [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) text format.
	AutoUpgradeStartTime string `json:"autoUpgradeStartTime,omitempty"`

	// Description: [Output only] This field is set when upgrades are about
	// to commence
	// with the description of the upgrade.
	Description string `json:"description,omitempty"`

	// ForceSendFields is a list of field names (e.g.
	// "AutoUpgradeStartTime") to unconditionally include in API requests.
	// By default, fields with empty values are omitted from API requests.
	// However, any non-pointer, non-interface field appearing in
	// ForceSendFields will be sent to the server regardless of whether the
	// field is empty or not. This may be used to include empty fields in
	// Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AutoUpgradeStartTime") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *AutoUpgradeOptions) MarshalJSON() ([]byte, error) {
	type noMethod AutoUpgradeOptions
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CancelOperationRequest: CancelOperationRequest cancels a single
// operation.
type CancelOperationRequest struct {
	// Name: The name (project, location, operation id) of the operation to
	// cancel.
	// Specified in the format 'projects/*/locations/*/operations/*'.
	Name string `json:"name,omitempty"`

	// OperationId: The server-assigned `name` of the operation.
	// This field is deprecated, use name instead.
	OperationId string `json:"operationId,omitempty"`

	// ProjectId: The Google Developers Console [project ID or
	// project
	// number](https://support.google.com/cloud/answer/6158840).
	// This
	//  field is deprecated, use name instead.
	ProjectId string `json:"projectId,omitempty"`

	// Zone: The name of the Google Compute
	// Engine
	// [zone](/compute/docs/zones#available) in which the operation
	// resides.
	// This field is deprecated, use name instead.
	Zone string `json:"zone,omitempty"`

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

func (s *CancelOperationRequest) MarshalJSON() ([]byte, error) {
	type noMethod CancelOperationRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CidrBlock: CidrBlock contains an optional name and one CIDR block.
type CidrBlock struct {
	// CidrBlock: cidr_block must be specified in CIDR notation.
	CidrBlock string `json:"cidrBlock,omitempty"`

	// DisplayName: display_name is an optional field for users to identify
	// CIDR blocks.
	DisplayName string `json:"displayName,omitempty"`

	// ForceSendFields is a list of field names (e.g. "CidrBlock") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CidrBlock") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *CidrBlock) MarshalJSON() ([]byte, error) {
	type noMethod CidrBlock
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ClientCertificateConfig: Configuration for client certificates on the
// cluster.
type ClientCertificateConfig struct {
	// IssueClientCertificate: Issue a client certificate.
	IssueClientCertificate bool `json:"issueClientCertificate,omitempty"`

	// ForceSendFields is a list of field names (e.g.
	// "IssueClientCertificate") to unconditionally include in API requests.
	// By default, fields with empty values are omitted from API requests.
	// However, any non-pointer, non-interface field appearing in
	// ForceSendFields will be sent to the server regardless of whether the
	// field is empty or not. This may be used to include empty fields in
	// Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "IssueClientCertificate")
	// to include in API requests with the JSON null value. By default,
	// fields with empty values are omitted from API requests. However, any
	// field with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *ClientCertificateConfig) MarshalJSON() ([]byte, error) {
	type noMethod ClientCertificateConfig
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Cluster: A Google Container Engine cluster.
type Cluster struct {
	// AddonsConfig: Configurations for the various addons available to run
	// in the cluster.
	AddonsConfig *AddonsConfig `json:"addonsConfig,omitempty"`

	// ClusterIpv4Cidr: The IP address range of the container pods in this
	// cluster,
	// in
	// [CIDR](http://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing)
	//
	// notation (e.g. `10.96.0.0/14`). Leave blank to have
	// one automatically chosen or specify a `/14` block in `10.0.0.0/8`.
	ClusterIpv4Cidr string `json:"clusterIpv4Cidr,omitempty"`

	// CreateTime: [Output only] The time the cluster was created,
	// in
	// [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) text format.
	CreateTime string `json:"createTime,omitempty"`

	// CurrentMasterVersion: [Output only] The current software version of
	// the master endpoint.
	CurrentMasterVersion string `json:"currentMasterVersion,omitempty"`

	// CurrentNodeCount: [Output only] The number of nodes currently in the
	// cluster.
	CurrentNodeCount int64 `json:"currentNodeCount,omitempty"`

	// CurrentNodeVersion: [Output only] The current version of the node
	// software components.
	// If they are currently at multiple versions because they're in the
	// process
	// of being upgraded, this reflects the minimum version of all nodes.
	CurrentNodeVersion string `json:"currentNodeVersion,omitempty"`

	// Description: An optional description of this cluster.
	Description string `json:"description,omitempty"`

	// EnableKubernetesAlpha: Kubernetes alpha features are enabled on this
	// cluster. This includes alpha
	// API groups (e.g. v1beta1) and features that may not be production
	// ready in
	// the kubernetes version of the master and nodes.
	// The cluster has no SLA for uptime and master/node upgrades are
	// disabled.
	// Alpha enabled clusters are automatically deleted thirty days
	// after
	// creation.
	EnableKubernetesAlpha bool `json:"enableKubernetesAlpha,omitempty"`

	// Endpoint: [Output only] The IP address of this cluster's master
	// endpoint.
	// The endpoint can be accessed from the internet
	// at
	// `https://username:password@endpoint/`.
	//
	// See the `masterAuth` property of this resource for username
	// and
	// password information.
	Endpoint string `json:"endpoint,omitempty"`

	// ExpireTime: [Output only] The time the cluster will be
	// automatically
	// deleted in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) text
	// format.
	ExpireTime string `json:"expireTime,omitempty"`

	// InitialClusterVersion: The initial Kubernetes version for this
	// cluster.  Valid versions are those
	// found in validMasterVersions returned by getServerConfig.  The
	// version can
	// be upgraded over time; such upgrades are reflected
	// in
	// currentMasterVersion and currentNodeVersion.
	InitialClusterVersion string `json:"initialClusterVersion,omitempty"`

	// InitialNodeCount: The number of nodes to create in this cluster. You
	// must ensure that your
	// Compute Engine <a href="/compute/docs/resource-quotas">resource
	// quota</a>
	// is sufficient for this number of instances. You must also have
	// available
	// firewall and routes quota.
	// For requests, this field should only be used in lieu of a
	// "node_pool" object, since this configuration (along with
	// the
	// "node_config") will be used to create a "NodePool" object with
	// an
	// auto-generated name. Do not use this and a node_pool at the same
	// time.
	InitialNodeCount int64 `json:"initialNodeCount,omitempty"`

	// InstanceGroupUrls: [Output only] The resource URLs of
	// [instance
	// groups](/compute/docs/instance-groups/) associated with this
	// cluster.
	InstanceGroupUrls []string `json:"instanceGroupUrls,omitempty"`

	// IpAllocationPolicy: Configuration for cluster IP allocation.
	IpAllocationPolicy *IPAllocationPolicy `json:"ipAllocationPolicy,omitempty"`

	// Location: [Output only] The name of the Google Compute
	// Engine
	// [zone](/compute/docs/regions-zones/regions-zones#available)
	// or
	// [region](/compute/docs/regions-zones/regions-zones#available) in
	// which
	// the cluster resides.
	Location string `json:"location,omitempty"`

	// Locations: The list of Google Compute
	// Engine
	// [locations](/compute/docs/zones#available) in which the cluster's
	// nodes
	// should be located.
	Locations []string `json:"locations,omitempty"`

	// LoggingService: The logging service the cluster should use to write
	// logs.
	// Currently available options:
	//
	// * `logging.googleapis.com` - the Google Cloud Logging service.
	// * `none` - no logs will be exported from the cluster.
	// * if left as an empty string,`logging.googleapis.com` will be used.
	LoggingService string `json:"loggingService,omitempty"`

	// MaintenancePolicy: Configure the maintenance policy for this cluster.
	MaintenancePolicy *MaintenancePolicy `json:"maintenancePolicy,omitempty"`

	// MasterAuth: The authentication information for accessing the master
	// endpoint.
	MasterAuth *MasterAuth `json:"masterAuth,omitempty"`

	// MasterAuthorizedNetworksConfig: The configuration options for master
	// authorized networks feature.
	MasterAuthorizedNetworksConfig *MasterAuthorizedNetworksConfig `json:"masterAuthorizedNetworksConfig,omitempty"`

	// MonitoringService: The monitoring service the cluster should use to
	// write metrics.
	// Currently available options:
	//
	// * `monitoring.googleapis.com` - the Google Cloud Monitoring
	// service.
	// * `none` - no metrics will be exported from the cluster.
	// * if left as an empty string, `monitoring.googleapis.com` will be
	// used.
	MonitoringService string `json:"monitoringService,omitempty"`

	// Name: The name of this cluster. The name must be unique within this
	// project
	// and zone, and can be up to 40 characters with the following
	// restrictions:
	//
	// * Lowercase letters, numbers, and hyphens only.
	// * Must start with a letter.
	// * Must end with a number or a letter.
	Name string `json:"name,omitempty"`

	// Network: The name of the Google Compute
	// Engine
	// [network](/compute/docs/networks-and-firewalls#networks) to which
	// the
	// cluster is connected. If left unspecified, the `default` network
	// will be used.
	Network string `json:"network,omitempty"`

	// NetworkPolicy: Configuration options for the NetworkPolicy feature.
	NetworkPolicy *NetworkPolicy `json:"networkPolicy,omitempty"`

	// NodeConfig: Parameters used in creating the cluster's nodes.
	// See `nodeConfig` for the description of its properties.
	// For requests, this field should only be used in lieu of a
	// "node_pool" object, since this configuration (along with
	// the
	// "initial_node_count") will be used to create a "NodePool" object with
	// an
	// auto-generated name. Do not use this and a node_pool at the same
	// time.
	// For responses, this field will be populated with the node
	// configuration of
	// the first node pool.
	//
	// If unspecified, the defaults are used.
	NodeConfig *NodeConfig `json:"nodeConfig,omitempty"`

	// NodeIpv4CidrSize: [Output only] The size of the address space on each
	// node for hosting
	// containers. This is provisioned from within the
	// `container_ipv4_cidr`
	// range.
	NodeIpv4CidrSize int64 `json:"nodeIpv4CidrSize,omitempty"`

	// NodePools: The node pools associated with this cluster.
	// This field should not be set if "node_config" or "initial_node_count"
	// are
	// specified.
	NodePools []*NodePool `json:"nodePools,omitempty"`

	// PodSecurityPolicyConfig: Configuration for the PodSecurityPolicy
	// feature.
	PodSecurityPolicyConfig *PodSecurityPolicyConfig `json:"podSecurityPolicyConfig,omitempty"`

	// SelfLink: [Output only] Server-defined URL for the resource.
	SelfLink string `json:"selfLink,omitempty"`

	// ServicesIpv4Cidr: [Output only] The IP address range of the
	// Kubernetes services in
	// this cluster,
	// in
	// [CIDR](http://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing)
	//
	// notation (e.g. `1.2.3.4/29`). Service addresses are
	// typically put in the last `/16` from the container CIDR.
	ServicesIpv4Cidr string `json:"servicesIpv4Cidr,omitempty"`

	// Status: [Output only] The current status of this cluster.
	//
	// Possible values:
	//   "STATUS_UNSPECIFIED" - Not set.
	//   "PROVISIONING" - The PROVISIONING state indicates the cluster is
	// being created.
	//   "RUNNING" - The RUNNING state indicates the cluster has been
	// created and is fully
	// usable.
	//   "RECONCILING" - The RECONCILING state indicates that some work is
	// actively being done on
	// the cluster, such as upgrading the master or node software. Details
	// can
	// be found in the `statusMessage` field.
	//   "STOPPING" - The STOPPING state indicates the cluster is being
	// deleted.
	//   "ERROR" - The ERROR state indicates the cluster may be unusable.
	// Details
	// can be found in the `statusMessage` field.
	Status string `json:"status,omitempty"`

	// StatusMessage: [Output only] Additional information about the current
	// status of this
	// cluster, if available.
	StatusMessage string `json:"statusMessage,omitempty"`

	// Subnetwork: The name of the Google Compute
	// Engine
	// [subnetwork](/compute/docs/subnetworks) to which the
	// cluster is connected.
	Subnetwork string `json:"subnetwork,omitempty"`

	// Zone: [Output only] The name of the Google Compute
	// Engine
	// [zone](/compute/docs/zones#available) in which the
	// cluster
	// resides.
	// This field is deprecated, use location instead.
	Zone string `json:"zone,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "AddonsConfig") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AddonsConfig") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Cluster) MarshalJSON() ([]byte, error) {
	type noMethod Cluster
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ClusterUpdate: ClusterUpdate describes an update to the cluster.
// Exactly one update can
// be applied to a cluster with each request, so at most one field can
// be
// provided.
type ClusterUpdate struct {
	// DesiredAddonsConfig: Configurations for the various addons available
	// to run in the cluster.
	DesiredAddonsConfig *AddonsConfig `json:"desiredAddonsConfig,omitempty"`

	// DesiredImageType: The desired image type for the node pool.
	// NOTE: Set the "desired_node_pool" field as well.
	DesiredImageType string `json:"desiredImageType,omitempty"`

	// DesiredLocations: The desired list of Google Compute
	// Engine
	// [locations](/compute/docs/zones#available) in which the cluster's
	// nodes
	// should be located. Changing the locations a cluster is in will
	// result
	// in nodes being either created or removed from the cluster, depending
	// on
	// whether locations are being added or removed.
	//
	// This list must always include the cluster's primary zone.
	DesiredLocations []string `json:"desiredLocations,omitempty"`

	// DesiredMasterAuthorizedNetworksConfig: The desired configuration
	// options for master authorized networks feature.
	DesiredMasterAuthorizedNetworksConfig *MasterAuthorizedNetworksConfig `json:"desiredMasterAuthorizedNetworksConfig,omitempty"`

	// DesiredMasterVersion: The Kubernetes version to change the master to.
	// The only valid value is the
	// latest supported version. Use "-" to have the server automatically
	// select
	// the latest version.
	DesiredMasterVersion string `json:"desiredMasterVersion,omitempty"`

	// DesiredMonitoringService: The monitoring service the cluster should
	// use to write metrics.
	// Currently available options:
	//
	// * "monitoring.googleapis.com" - the Google Cloud Monitoring service
	// * "none" - no metrics will be exported from the cluster
	DesiredMonitoringService string `json:"desiredMonitoringService,omitempty"`

	// DesiredNodePoolAutoscaling: Autoscaler configuration for the node
	// pool specified in
	// desired_node_pool_id. If there is only one pool in the
	// cluster and desired_node_pool_id is not provided then
	// the change applies to that single node pool.
	DesiredNodePoolAutoscaling *NodePoolAutoscaling `json:"desiredNodePoolAutoscaling,omitempty"`

	// DesiredNodePoolId: The node pool to be upgraded. This field is
	// mandatory if
	// "desired_node_version", "desired_image_family"
	// or
	// "desired_node_pool_autoscaling" is specified and there is more than
	// one
	// node pool on the cluster.
	DesiredNodePoolId string `json:"desiredNodePoolId,omitempty"`

	// DesiredNodeVersion: The Kubernetes version to change the nodes to
	// (typically an
	// upgrade). Use `-` to upgrade to the latest version supported by
	// the server.
	DesiredNodeVersion string `json:"desiredNodeVersion,omitempty"`

	// DesiredPodSecurityPolicyConfig: The desired configuration options for
	// the PodSecurityPolicy feature.
	DesiredPodSecurityPolicyConfig *PodSecurityPolicyConfig `json:"desiredPodSecurityPolicyConfig,omitempty"`

	// ForceSendFields is a list of field names (e.g. "DesiredAddonsConfig")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "DesiredAddonsConfig") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *ClusterUpdate) MarshalJSON() ([]byte, error) {
	type noMethod ClusterUpdate
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CompleteIPRotationRequest: CompleteIPRotationRequest moves the
// cluster master back into single-IP mode.
type CompleteIPRotationRequest struct {
	// ClusterId: The name of the cluster.
	// This field is deprecated, use name instead.
	ClusterId string `json:"clusterId,omitempty"`

	// Name: The name (project, location, cluster id) of the cluster to
	// complete IP rotation.
	// Specified in the format 'projects/*/locations/*/clusters/*'.
	Name string `json:"name,omitempty"`

	// ProjectId: The Google Developers Console [project ID or
	// project
	// number](https://developers.google.com/console/help/new/#projec
	// tnumber).
	// This field is deprecated, use name instead.
	ProjectId string `json:"projectId,omitempty"`

	// Zone: The name of the Google Compute
	// Engine
	// [zone](/compute/docs/zones#available) in which the
	// cluster
	// resides.
	// This field is deprecated, use name instead.
	Zone string `json:"zone,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ClusterId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ClusterId") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *CompleteIPRotationRequest) MarshalJSON() ([]byte, error) {
	type noMethod CompleteIPRotationRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CreateClusterRequest: CreateClusterRequest creates a cluster.
type CreateClusterRequest struct {
	// Cluster: A
	// [cluster
	// resource](/container-engine/reference/rest/v1beta1/projects.z
	// ones.clusters)
	Cluster *Cluster `json:"cluster,omitempty"`

	// Parent: The parent (project and location) where the cluster will be
	// created.
	// Specified in the format 'projects/*/locations/*'.
	Parent string `json:"parent,omitempty"`

	// ProjectId: The Google Developers Console [project ID or
	// project
	// number](https://support.google.com/cloud/answer/6158840).
	// This
	//  field is deprecated, use parent instead.
	ProjectId string `json:"projectId,omitempty"`

	// Zone: The name of the Google Compute
	// Engine
	// [zone](/compute/docs/zones#available) in which the
	// cluster
	// resides.
	// This field is deprecated, use parent instead.
	Zone string `json:"zone,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Cluster") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Cluster") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *CreateClusterRequest) MarshalJSON() ([]byte, error) {
	type noMethod CreateClusterRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// CreateNodePoolRequest: CreateNodePoolRequest creates a node pool for
// a cluster.
type CreateNodePoolRequest struct {
	// ClusterId: The name of the cluster.
	// This field is deprecated, use parent instead.
	ClusterId string `json:"clusterId,omitempty"`

	// NodePool: The node pool to create.
	NodePool *NodePool `json:"nodePool,omitempty"`

	// Parent: The parent (project, location, cluster id) where the node
	// pool will be created.
	// Specified in the format
	// 'projects/*/locations/*/clusters/*/nodePools/*'.
	Parent string `json:"parent,omitempty"`

	// ProjectId: The Google Developers Console [project ID or
	// project
	// number](https://developers.google.com/console/help/new/#projec
	// tnumber).
	// This field is deprecated, use parent instead.
	ProjectId string `json:"projectId,omitempty"`

	// Zone: The name of the Google Compute
	// Engine
	// [zone](/compute/docs/zones#available) in which the
	// cluster
	// resides.
	// This field is deprecated, use parent instead.
	Zone string `json:"zone,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ClusterId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ClusterId") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *CreateNodePoolRequest) MarshalJSON() ([]byte, error) {
	type noMethod CreateNodePoolRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// DailyMaintenanceWindow: Time window specified for daily maintenance
// operations.
type DailyMaintenanceWindow struct {
	// Duration: [Output only] Duration of the time window, automatically
	// chosen to be
	// smallest possible in the given scenario.
	Duration string `json:"duration,omitempty"`

	// StartTime: Time within the maintenance window to start the
	// maintenance operations.
	// It must be in format "HH:MM”, where HH : [00-23] and MM : [00-59]
	// GMT.
	StartTime string `json:"startTime,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Duration") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Duration") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *DailyMaintenanceWindow) MarshalJSON() ([]byte, error) {
	type noMethod DailyMaintenanceWindow
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
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

// HorizontalPodAutoscaling: Configuration options for the horizontal
// pod autoscaling feature, which
// increases or decreases the number of replica pods a replication
// controller
// has based on the resource usage of the existing pods.
type HorizontalPodAutoscaling struct {
	// Disabled: Whether the Horizontal Pod Autoscaling feature is enabled
	// in the cluster.
	// When enabled, it ensures that a Heapster pod is running in the
	// cluster,
	// which is also used by the Cloud Monitoring service.
	Disabled bool `json:"disabled,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Disabled") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Disabled") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *HorizontalPodAutoscaling) MarshalJSON() ([]byte, error) {
	type noMethod HorizontalPodAutoscaling
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// HttpLoadBalancing: Configuration options for the HTTP (L7) load
// balancing controller addon,
// which makes it easy to set up HTTP load balancers for services in a
// cluster.
type HttpLoadBalancing struct {
	// Disabled: Whether the HTTP Load Balancing controller is enabled in
	// the cluster.
	// When enabled, it runs a small pod in the cluster that manages the
	// load
	// balancers.
	Disabled bool `json:"disabled,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Disabled") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Disabled") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *HttpLoadBalancing) MarshalJSON() ([]byte, error) {
	type noMethod HttpLoadBalancing
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// IPAllocationPolicy: Configuration for controlling how IPs are
// allocated in the cluster.
type IPAllocationPolicy struct {
	// ClusterIpv4Cidr: This field is deprecated, use
	// cluster_ipv4_cidr_block.
	ClusterIpv4Cidr string `json:"clusterIpv4Cidr,omitempty"`

	// ClusterIpv4CidrBlock: The IP address range for the cluster pod IPs.
	// If this field is set, then
	// `cluster.cluster_ipv4_cidr` must be left blank.
	//
	// This field is only applicable when `use_ip_aliases` is true.
	//
	// Set to blank to have a range chosen with the default size.
	//
	// Set to /netmask (e.g. `/14`) to have a range chosen with a
	// specific
	// netmask.
	//
	// Set to
	// a
	// [CIDR](http://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing)
	//
	// notation (e.g. `10.96.0.0/14`) from the RFC-1918 private networks
	// (e.g.
	// `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`) to pick a specific
	// range
	// to use.
	ClusterIpv4CidrBlock string `json:"clusterIpv4CidrBlock,omitempty"`

	// ClusterSecondaryRangeName: The name of the secondary range to be used
	// for the cluster CIDR
	// block.  The secondary range will be used for pod IP
	// addresses. This must be an existing secondary range associated
	// with the cluster subnetwork.
	//
	// This field is only applicable with use_ip_aliases
	// and
	// create_subnetwork is false.
	ClusterSecondaryRangeName string `json:"clusterSecondaryRangeName,omitempty"`

	// CreateSubnetwork: Whether a new subnetwork will be created
	// automatically for the cluster.
	//
	// This field is only applicable when `use_ip_aliases` is true.
	CreateSubnetwork bool `json:"createSubnetwork,omitempty"`

	// NodeIpv4Cidr: This field is deprecated, use node_ipv4_cidr_block.
	NodeIpv4Cidr string `json:"nodeIpv4Cidr,omitempty"`

	// NodeIpv4CidrBlock: The IP address range of the instance IPs in this
	// cluster.
	//
	// This is applicable only if `create_subnetwork` is true.
	//
	// Set to blank to have a range chosen with the default size.
	//
	// Set to /netmask (e.g. `/14`) to have a range chosen with a
	// specific
	// netmask.
	//
	// Set to
	// a
	// [CIDR](http://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing)
	//
	// notation (e.g. `10.96.0.0/14`) from the RFC-1918 private networks
	// (e.g.
	// `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`) to pick a specific
	// range
	// to use.
	NodeIpv4CidrBlock string `json:"nodeIpv4CidrBlock,omitempty"`

	// ServicesIpv4Cidr: This field is deprecated, use
	// services_ipv4_cidr_block.
	ServicesIpv4Cidr string `json:"servicesIpv4Cidr,omitempty"`

	// ServicesIpv4CidrBlock: The IP address range of the services IPs in
	// this cluster. If blank, a range
	// will be automatically chosen with the default size.
	//
	// This field is only applicable when `use_ip_aliases` is true.
	//
	// Set to blank to have a range chosen with the default size.
	//
	// Set to /netmask (e.g. `/14`) to have a range chosen with a
	// specific
	// netmask.
	//
	// Set to
	// a
	// [CIDR](http://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing)
	//
	// notation (e.g. `10.96.0.0/14`) from the RFC-1918 private networks
	// (e.g.
	// `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`) to pick a specific
	// range
	// to use.
	ServicesIpv4CidrBlock string `json:"servicesIpv4CidrBlock,omitempty"`

	// ServicesSecondaryRangeName: The name of the secondary range to be
	// used as for the services
	// CIDR block.  The secondary range will be used for service
	// ClusterIPs. This must be an existing secondary range associated
	// with the cluster subnetwork.
	//
	// This field is only applicable with use_ip_aliases
	// and
	// create_subnetwork is false.
	ServicesSecondaryRangeName string `json:"servicesSecondaryRangeName,omitempty"`

	// SubnetworkName: A custom subnetwork name to be used if
	// `create_subnetwork` is true.  If
	// this field is empty, then an automatic name will be chosen for the
	// new
	// subnetwork.
	SubnetworkName string `json:"subnetworkName,omitempty"`

	// UseIpAliases: Whether alias IPs will be used for pod IPs in the
	// cluster.
	UseIpAliases bool `json:"useIpAliases,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ClusterIpv4Cidr") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ClusterIpv4Cidr") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *IPAllocationPolicy) MarshalJSON() ([]byte, error) {
	type noMethod IPAllocationPolicy
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// KubernetesDashboard: Configuration for the Kubernetes Dashboard.
type KubernetesDashboard struct {
	// Disabled: Whether the Kubernetes Dashboard is enabled for this
	// cluster.
	Disabled bool `json:"disabled,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Disabled") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Disabled") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *KubernetesDashboard) MarshalJSON() ([]byte, error) {
	type noMethod KubernetesDashboard
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListClustersResponse: ListClustersResponse is the result of
// ListClustersRequest.
type ListClustersResponse struct {
	// Clusters: A list of clusters in the project in the specified zone,
	// or
	// across all ones.
	Clusters []*Cluster `json:"clusters,omitempty"`

	// MissingZones: If any zones are listed here, the list of clusters
	// returned
	// may be missing those zones.
	MissingZones []string `json:"missingZones,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Clusters") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Clusters") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListClustersResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListClustersResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListNodePoolsResponse: ListNodePoolsResponse is the result of
// ListNodePoolsRequest.
type ListNodePoolsResponse struct {
	// NodePools: A list of node pools for a cluster.
	NodePools []*NodePool `json:"nodePools,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "NodePools") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "NodePools") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListNodePoolsResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListNodePoolsResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListOperationsResponse: ListOperationsResponse is the result of
// ListOperationsRequest.
type ListOperationsResponse struct {
	// MissingZones: If any zones are listed here, the list of operations
	// returned
	// may be missing the operations from those zones.
	MissingZones []string `json:"missingZones,omitempty"`

	// Operations: A list of operations in the project in the specified
	// zone.
	Operations []*Operation `json:"operations,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "MissingZones") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "MissingZones") to include
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

// MaintenancePolicy: MaintenancePolicy defines the maintenance policy
// to be used for the cluster.
type MaintenancePolicy struct {
	// Window: Specifies the maintenance window in which maintenance may be
	// performed.
	Window *MaintenanceWindow `json:"window,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Window") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Window") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *MaintenancePolicy) MarshalJSON() ([]byte, error) {
	type noMethod MaintenancePolicy
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// MaintenanceWindow: MaintenanceWindow defines the maintenance window
// to be used for the cluster.
type MaintenanceWindow struct {
	// DailyMaintenanceWindow: DailyMaintenanceWindow specifies a daily
	// maintenance operation window.
	DailyMaintenanceWindow *DailyMaintenanceWindow `json:"dailyMaintenanceWindow,omitempty"`

	// ForceSendFields is a list of field names (e.g.
	// "DailyMaintenanceWindow") to unconditionally include in API requests.
	// By default, fields with empty values are omitted from API requests.
	// However, any non-pointer, non-interface field appearing in
	// ForceSendFields will be sent to the server regardless of whether the
	// field is empty or not. This may be used to include empty fields in
	// Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "DailyMaintenanceWindow")
	// to include in API requests with the JSON null value. By default,
	// fields with empty values are omitted from API requests. However, any
	// field with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *MaintenanceWindow) MarshalJSON() ([]byte, error) {
	type noMethod MaintenanceWindow
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// MasterAuth: The authentication information for accessing the master
// endpoint.
// Authentication can be done using HTTP basic auth or using
// client
// certificates.
type MasterAuth struct {
	// ClientCertificate: [Output only] Base64-encoded public certificate
	// used by clients to
	// authenticate to the cluster endpoint.
	ClientCertificate string `json:"clientCertificate,omitempty"`

	// ClientCertificateConfig: Configuration for client certificate
	// authentication on the cluster.  If no
	// configuration is specified, a client certificate is issued.
	ClientCertificateConfig *ClientCertificateConfig `json:"clientCertificateConfig,omitempty"`

	// ClientKey: [Output only] Base64-encoded private key used by clients
	// to authenticate
	// to the cluster endpoint.
	ClientKey string `json:"clientKey,omitempty"`

	// ClusterCaCertificate: [Output only] Base64-encoded public certificate
	// that is the root of
	// trust for the cluster.
	ClusterCaCertificate string `json:"clusterCaCertificate,omitempty"`

	// Password: The password to use for HTTP basic authentication to the
	// master endpoint.
	// Because the master endpoint is open to the Internet, you should
	// create a
	// strong password.  If a password is provided for cluster creation,
	// username
	// must be non-empty.
	Password string `json:"password,omitempty"`

	// Username: The username to use for HTTP basic authentication to the
	// master endpoint.
	// For clusters v1.6.0 and later, you can disable basic authentication
	// by
	// providing an empty username.
	Username string `json:"username,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ClientCertificate")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ClientCertificate") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *MasterAuth) MarshalJSON() ([]byte, error) {
	type noMethod MasterAuth
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// MasterAuthorizedNetworksConfig: Configuration options for the master
// authorized networks feature. Enabled
// master authorized networks will disallow all external traffic to
// access
// Kubernetes master through HTTPS except traffic from the given CIDR
// blocks,
// Google Compute Engine Public IPs and Google Prod IPs.
type MasterAuthorizedNetworksConfig struct {
	// CidrBlocks: cidr_blocks define up to 10 external networks that could
	// access
	// Kubernetes master through HTTPS.
	CidrBlocks []*CidrBlock `json:"cidrBlocks,omitempty"`

	// Enabled: Whether or not master authorized networks is enabled.
	Enabled bool `json:"enabled,omitempty"`

	// ForceSendFields is a list of field names (e.g. "CidrBlocks") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CidrBlocks") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *MasterAuthorizedNetworksConfig) MarshalJSON() ([]byte, error) {
	type noMethod MasterAuthorizedNetworksConfig
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// NetworkPolicy: Configuration options for the NetworkPolicy
// feature.
// https://kubernetes.io/docs/concepts/services-networking/netwo
// rkpolicies/
type NetworkPolicy struct {
	// Enabled: Whether network policy is enabled on the cluster.
	Enabled bool `json:"enabled,omitempty"`

	// Provider: The selected network policy provider.
	//
	// Possible values:
	//   "PROVIDER_UNSPECIFIED" - Not set
	//   "CALICO" - Tigera (Calico Felix).
	Provider string `json:"provider,omitempty"`

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

func (s *NetworkPolicy) MarshalJSON() ([]byte, error) {
	type noMethod NetworkPolicy
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// NetworkPolicyConfig: Configuration for NetworkPolicy. This only
// tracks whether the addon
// is enabled or not on the Master, it does not track whether network
// policy
// is enabled for the nodes.
type NetworkPolicyConfig struct {
	// Disabled: Whether NetworkPolicy is enabled for this cluster.
	Disabled bool `json:"disabled,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Disabled") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Disabled") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *NetworkPolicyConfig) MarshalJSON() ([]byte, error) {
	type noMethod NetworkPolicyConfig
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// NodeConfig: Parameters that describe the nodes in a cluster.
type NodeConfig struct {
	// Accelerators: A list of hardware accelerators to be attached to each
	// node.
	// See https://cloud.google.com/compute/docs/gpus for more information
	// about
	// support for GPUs.
	Accelerators []*AcceleratorConfig `json:"accelerators,omitempty"`

	// DiskSizeGb: Size of the disk attached to each node, specified in
	// GB.
	// The smallest allowed disk size is 10GB.
	//
	// If unspecified, the default disk size is 100GB.
	DiskSizeGb int64 `json:"diskSizeGb,omitempty"`

	// ImageType: The image type to use for this node. Note that for a given
	// image type,
	// the latest version of it will be used.
	ImageType string `json:"imageType,omitempty"`

	// Labels: The map of Kubernetes labels (key/value pairs) to be applied
	// to each node.
	// These will added in addition to any default label(s) that
	// Kubernetes may apply to the node.
	// In case of conflict in label keys, the applied set may differ
	// depending on
	// the Kubernetes version -- it's best to assume the behavior is
	// undefined
	// and conflicts should be avoided.
	// For more information, including usage and the valid values,
	// see:
	// https://kubernetes.io/docs/concepts/overview/working-with-objects
	// /labels/
	Labels map[string]string `json:"labels,omitempty"`

	// LocalSsdCount: The number of local SSD disks to be attached to the
	// node.
	//
	// The limit for this value is dependant upon the maximum number
	// of
	// disks available on a machine per zone.
	// See:
	// https://cloud.google.com/compute/docs/disks/local-ssd#local_ssd_l
	// imits
	// for more information.
	LocalSsdCount int64 `json:"localSsdCount,omitempty"`

	// MachineType: The name of a Google Compute Engine
	// [machine
	// type](/compute/docs/machine-types) (e.g.
	// `n1-standard-1`).
	//
	// If unspecified, the default machine type is
	// `n1-standard-1`.
	MachineType string `json:"machineType,omitempty"`

	// Metadata: The metadata key/value pairs assigned to instances in the
	// cluster.
	//
	// Keys must conform to the regexp [a-zA-Z0-9-_]+ and be less than 128
	// bytes
	// in length. These are reflected as part of a URL in the metadata
	// server.
	// Additionally, to avoid ambiguity, keys must not conflict with any
	// other
	// metadata keys for the project or be one of the four reserved
	// keys:
	// "instance-template", "kube-env", "startup-script", and
	// "user-data"
	//
	// Values are free-form strings, and only have meaning as interpreted
	// by
	// the image running in the instance. The only restriction placed on
	// them is
	// that each value's size must be less than or equal to 32 KB.
	//
	// The total size of all keys and values must be less than 512 KB.
	Metadata map[string]string `json:"metadata,omitempty"`

	// MinCpuPlatform: Minimum CPU platform to be used by this instance. The
	// instance may be
	// scheduled on the specified or newer CPU platform. Applicable values
	// are the
	// friendly names of CPU platforms, such as
	// <code>minCpuPlatform: &quot;Intel Haswell&quot;</code>
	// or
	// <code>minCpuPlatform: &quot;Intel Sandy Bridge&quot;</code>. For
	// more
	// information, read [how to specify min CPU
	// platform](https://cloud.google.com/compute/docs/instances/specify-min-
	// cpu-platform)
	MinCpuPlatform string `json:"minCpuPlatform,omitempty"`

	// OauthScopes: The set of Google API scopes to be made available on all
	// of the
	// node VMs under the "default" service account.
	//
	// The following scopes are recommended, but not required, and by
	// default are
	// not included:
	//
	// * `https://www.googleapis.com/auth/compute` is required for
	// mounting
	// persistent storage on your nodes.
	// * `https://www.googleapis.com/auth/devstorage.read_only` is required
	// for
	// communicating with **gcr.io**
	// (the [Google Container Registry](/container-registry/)).
	//
	// If unspecified, no scopes are added, unless Cloud Logging or
	// Cloud
	// Monitoring are enabled, in which case their required scopes will be
	// added.
	OauthScopes []string `json:"oauthScopes,omitempty"`

	// Preemptible: Whether the nodes are created as preemptible VM
	// instances.
	// See:
	// https://cloud.google.com/compute/docs/instances/preemptible for
	// more
	// inforamtion about preemptible VM instances.
	Preemptible bool `json:"preemptible,omitempty"`

	// ServiceAccount: The Google Cloud Platform Service Account to be used
	// by the node VMs. If
	// no Service Account is specified, the "default" service account is
	// used.
	ServiceAccount string `json:"serviceAccount,omitempty"`

	// Tags: The list of instance tags applied to all nodes. Tags are used
	// to identify
	// valid sources or targets for network firewalls and are specified
	// by
	// the client during cluster or node pool creation. Each tag within the
	// list
	// must comply with RFC1035.
	Tags []string `json:"tags,omitempty"`

	// Taints: List of kubernetes taints to be applied to each node.
	//
	// For more information, including usage and the valid values,
	// see:
	// https://kubernetes.io/docs/concepts/configuration/taint-and-toler
	// ation/
	Taints []*NodeTaint `json:"taints,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Accelerators") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Accelerators") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *NodeConfig) MarshalJSON() ([]byte, error) {
	type noMethod NodeConfig
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// NodeManagement: NodeManagement defines the set of node management
// services turned on for the
// node pool.
type NodeManagement struct {
	// AutoRepair: Whether the nodes will be automatically repaired.
	AutoRepair bool `json:"autoRepair,omitempty"`

	// AutoUpgrade: Whether the nodes will be automatically upgraded.
	AutoUpgrade bool `json:"autoUpgrade,omitempty"`

	// UpgradeOptions: Specifies the Auto Upgrade knobs for the node pool.
	UpgradeOptions *AutoUpgradeOptions `json:"upgradeOptions,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AutoRepair") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AutoRepair") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *NodeManagement) MarshalJSON() ([]byte, error) {
	type noMethod NodeManagement
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// NodePool: NodePool contains the name and configuration for a
// cluster's node pool.
// Node pools are a set of nodes (i.e. VM's), with a common
// configuration and
// specification, under the control of the cluster master. They may have
// a set
// of Kubernetes labels applied to them, which may be used to reference
// them
// during pod scheduling. They may also be resized up or down, to
// accommodate
// the workload.
type NodePool struct {
	// Autoscaling: Autoscaler configuration for this NodePool. Autoscaler
	// is enabled
	// only if a valid configuration is present.
	Autoscaling *NodePoolAutoscaling `json:"autoscaling,omitempty"`

	// Config: The node configuration of the pool.
	Config *NodeConfig `json:"config,omitempty"`

	// InitialNodeCount: The initial node count for the pool. You must
	// ensure that your
	// Compute Engine <a href="/compute/docs/resource-quotas">resource
	// quota</a>
	// is sufficient for this number of instances. You must also have
	// available
	// firewall and routes quota.
	InitialNodeCount int64 `json:"initialNodeCount,omitempty"`

	// InstanceGroupUrls: [Output only] The resource URLs of
	// [instance
	// groups](/compute/docs/instance-groups/) associated with this
	// node pool.
	InstanceGroupUrls []string `json:"instanceGroupUrls,omitempty"`

	// Management: NodeManagement configuration for this NodePool.
	Management *NodeManagement `json:"management,omitempty"`

	// Name: The name of the node pool.
	Name string `json:"name,omitempty"`

	// SelfLink: [Output only] Server-defined URL for the resource.
	SelfLink string `json:"selfLink,omitempty"`

	// Status: [Output only] The status of the nodes in this pool instance.
	//
	// Possible values:
	//   "STATUS_UNSPECIFIED" - Not set.
	//   "PROVISIONING" - The PROVISIONING state indicates the node pool is
	// being created.
	//   "RUNNING" - The RUNNING state indicates the node pool has been
	// created
	// and is fully usable.
	//   "RUNNING_WITH_ERROR" - The RUNNING_WITH_ERROR state indicates the
	// node pool has been created
	// and is partially usable. Some error state has occurred and
	// some
	// functionality may be impaired. Customer may need to reissue a
	// request
	// or trigger a new update.
	//   "RECONCILING" - The RECONCILING state indicates that some work is
	// actively being done on
	// the node pool, such as upgrading node software. Details can
	// be found in the `statusMessage` field.
	//   "STOPPING" - The STOPPING state indicates the node pool is being
	// deleted.
	//   "ERROR" - The ERROR state indicates the node pool may be unusable.
	// Details
	// can be found in the `statusMessage` field.
	Status string `json:"status,omitempty"`

	// StatusMessage: [Output only] Additional information about the current
	// status of this
	// node pool instance, if available.
	StatusMessage string `json:"statusMessage,omitempty"`

	// Version: [Output only] The version of the Kubernetes of this node.
	Version string `json:"version,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Autoscaling") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Autoscaling") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *NodePool) MarshalJSON() ([]byte, error) {
	type noMethod NodePool
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// NodePoolAutoscaling: NodePoolAutoscaling contains information
// required by cluster autoscaler to
// adjust the size of the node pool to the current cluster usage.
type NodePoolAutoscaling struct {
	// Enabled: Is autoscaling enabled for this node pool.
	Enabled bool `json:"enabled,omitempty"`

	// MaxNodeCount: Maximum number of nodes in the NodePool. Must be >=
	// min_node_count. There
	// has to enough quota to scale up the cluster.
	MaxNodeCount int64 `json:"maxNodeCount,omitempty"`

	// MinNodeCount: Minimum number of nodes in the NodePool. Must be >= 1
	// and <=
	// max_node_count.
	MinNodeCount int64 `json:"minNodeCount,omitempty"`

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

func (s *NodePoolAutoscaling) MarshalJSON() ([]byte, error) {
	type noMethod NodePoolAutoscaling
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// NodeTaint: Kubernetes taint is comprised of three fields: key, value,
// and effect. Effect
// can only be one of three types:  NoSchedule, PreferNoSchedule or
// NoExecute.
//
// For more information, including usage and the valid values,
// see:
// https://kubernetes.io/docs/concepts/configuration/taint-and-toler
// ation/
type NodeTaint struct {
	// Effect: Effect for taint.
	//
	// Possible values:
	//   "NO_SCHEDULE" - NoSchedule
	//   "PREFER_NO_SCHEDULE" - PreferNoSchedule
	//   "NO_EXECUTE" - NoExecute
	Effect string `json:"effect,omitempty"`

	// Key: Key for taint.
	Key string `json:"key,omitempty"`

	// Value: Value for taint.
	Value string `json:"value,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Effect") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Effect") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *NodeTaint) MarshalJSON() ([]byte, error) {
	type noMethod NodeTaint
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Operation: This operation resource represents operations that may
// have happened or are
// happening on the cluster. All fields are output only.
type Operation struct {
	// Detail: Detailed operation progress, if available.
	Detail string `json:"detail,omitempty"`

	// EndTime: [Output only] The time the operation completed,
	// in
	// [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) text format.
	EndTime string `json:"endTime,omitempty"`

	// Location: [Output only] The name of the Google Compute
	// Engine
	// [zone](/compute/docs/regions-zones/regions-zones#available)
	// or
	// [region](/compute/docs/regions-zones/regions-zones#available) in
	// which
	// the cluster resides.
	Location string `json:"location,omitempty"`

	// Name: The server-assigned ID for the operation.
	Name string `json:"name,omitempty"`

	// OperationType: The operation type.
	//
	// Possible values:
	//   "TYPE_UNSPECIFIED" - Not set.
	//   "CREATE_CLUSTER" - Cluster create.
	//   "DELETE_CLUSTER" - Cluster delete.
	//   "UPGRADE_MASTER" - A master upgrade.
	//   "UPGRADE_NODES" - A node upgrade.
	//   "REPAIR_CLUSTER" - Cluster repair.
	//   "UPDATE_CLUSTER" - Cluster update.
	//   "CREATE_NODE_POOL" - Node pool create.
	//   "DELETE_NODE_POOL" - Node pool delete.
	//   "SET_NODE_POOL_MANAGEMENT" - Set node pool management.
	//   "AUTO_REPAIR_NODES" - Automatic node pool repair.
	//   "AUTO_UPGRADE_NODES" - Automatic node upgrade.
	//   "SET_LABELS" - Set labels.
	//   "SET_MASTER_AUTH" - Set/generate master auth materials
	//   "SET_NODE_POOL_SIZE" - Set node pool size.
	//   "SET_NETWORK_POLICY" - Updates network policy for a cluster.
	//   "SET_MAINTENANCE_POLICY" - Set the maintenance policy.
	OperationType string `json:"operationType,omitempty"`

	// SelfLink: Server-defined URL for the resource.
	SelfLink string `json:"selfLink,omitempty"`

	// StartTime: [Output only] The time the operation started,
	// in
	// [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) text format.
	StartTime string `json:"startTime,omitempty"`

	// Status: The current status of the operation.
	//
	// Possible values:
	//   "STATUS_UNSPECIFIED" - Not set.
	//   "PENDING" - The operation has been created.
	//   "RUNNING" - The operation is currently running.
	//   "DONE" - The operation is done, either cancelled or completed.
	//   "ABORTING" - The operation is aborting.
	Status string `json:"status,omitempty"`

	// StatusMessage: If an error has occurred, a textual description of the
	// error.
	StatusMessage string `json:"statusMessage,omitempty"`

	// TargetLink: Server-defined URL for the target of the operation.
	TargetLink string `json:"targetLink,omitempty"`

	// Zone: The name of the Google Compute
	// Engine
	// [zone](/compute/docs/zones#available) in which the operation
	// is taking place.
	// This field is deprecated, use location instead.
	Zone string `json:"zone,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Detail") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Detail") to include in API
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

// PodSecurityPolicyConfig: Configuration for the PodSecurityPolicy
// feature.
type PodSecurityPolicyConfig struct {
	// Enabled: Enable the PodSecurityPolicy controller for this cluster. If
	// enabled, pods
	// must be valid under a PodSecurityPolicy to be created.
	Enabled bool `json:"enabled,omitempty"`

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

func (s *PodSecurityPolicyConfig) MarshalJSON() ([]byte, error) {
	type noMethod PodSecurityPolicyConfig
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// RollbackNodePoolUpgradeRequest: RollbackNodePoolUpgradeRequest
// rollbacks the previously Aborted or Failed
// NodePool upgrade. This will be an no-op if the last upgrade
// successfully
// completed.
type RollbackNodePoolUpgradeRequest struct {
	// ClusterId: The name of the cluster to rollback.
	// This field is deprecated, use name instead.
	ClusterId string `json:"clusterId,omitempty"`

	// Name: The name (project, location, cluster, node pool id) of the node
	// poll to
	// rollback upgrade.
	// Specified in the format
	// 'projects/*/locations/*/clusters/*/nodePools/*'.
	Name string `json:"name,omitempty"`

	// NodePoolId: The name of the node pool to rollback.
	// This field is deprecated, use name instead.
	NodePoolId string `json:"nodePoolId,omitempty"`

	// ProjectId: The Google Developers Console [project ID or
	// project
	// number](https://support.google.com/cloud/answer/6158840).
	// This
	//  field is deprecated, use name instead.
	ProjectId string `json:"projectId,omitempty"`

	// Zone: The name of the Google Compute
	// Engine
	// [zone](/compute/docs/zones#available) in which the
	// cluster
	// resides.
	// This field is deprecated, use name instead.
	Zone string `json:"zone,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ClusterId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ClusterId") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *RollbackNodePoolUpgradeRequest) MarshalJSON() ([]byte, error) {
	type noMethod RollbackNodePoolUpgradeRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ServerConfig: Container Engine service configuration.
type ServerConfig struct {
	// DefaultClusterVersion: Version of Kubernetes the service deploys by
	// default.
	DefaultClusterVersion string `json:"defaultClusterVersion,omitempty"`

	// DefaultImageType: Default image type.
	DefaultImageType string `json:"defaultImageType,omitempty"`

	// ValidImageTypes: List of valid image types.
	ValidImageTypes []string `json:"validImageTypes,omitempty"`

	// ValidMasterVersions: List of valid master versions.
	ValidMasterVersions []string `json:"validMasterVersions,omitempty"`

	// ValidNodeVersions: List of valid node upgrade target versions.
	ValidNodeVersions []string `json:"validNodeVersions,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g.
	// "DefaultClusterVersion") to unconditionally include in API requests.
	// By default, fields with empty values are omitted from API requests.
	// However, any non-pointer, non-interface field appearing in
	// ForceSendFields will be sent to the server regardless of whether the
	// field is empty or not. This may be used to include empty fields in
	// Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "DefaultClusterVersion") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *ServerConfig) MarshalJSON() ([]byte, error) {
	type noMethod ServerConfig
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// SetLabelsRequest: SetLabelsRequest sets the Google Cloud Platform
// labels on a Google Container
// Engine cluster, which will in turn set them for Google Compute
// Engine
// resources used by that cluster
type SetLabelsRequest struct {
	// ClusterId: The name of the cluster.
	// This field is deprecated, use name instead.
	ClusterId string `json:"clusterId,omitempty"`

	// LabelFingerprint: The fingerprint of the previous set of labels for
	// this resource,
	// used to detect conflicts. The fingerprint is initially generated
	// by
	// Container Engine and changes after every request to modify or
	// update
	// labels. You must always provide an up-to-date fingerprint hash
	// when
	// updating or changing labels. Make a <code>get()</code> request to
	// the
	// resource to get the latest fingerprint.
	LabelFingerprint string `json:"labelFingerprint,omitempty"`

	// Name: The name (project, location, cluster id) of the cluster to set
	// labels.
	// Specified in the format 'projects/*/locations/*/clusters/*'.
	Name string `json:"name,omitempty"`

	// ProjectId: The Google Developers Console [project ID or
	// project
	// number](https://developers.google.com/console/help/new/#projec
	// tnumber).
	// This field is deprecated, use name instead.
	ProjectId string `json:"projectId,omitempty"`

	// ResourceLabels: The labels to set for that cluster.
	ResourceLabels map[string]string `json:"resourceLabels,omitempty"`

	// Zone: The name of the Google Compute
	// Engine
	// [zone](/compute/docs/zones#available) in which the
	// cluster
	// resides.
	// This field is deprecated, use name instead.
	Zone string `json:"zone,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ClusterId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ClusterId") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *SetLabelsRequest) MarshalJSON() ([]byte, error) {
	type noMethod SetLabelsRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// SetLegacyAbacRequest: SetLegacyAbacRequest enables or disables the
// ABAC authorization mechanism for
// a cluster.
type SetLegacyAbacRequest struct {
	// ClusterId: The name of the cluster to update.
	// This field is deprecated, use name instead.
	ClusterId string `json:"clusterId,omitempty"`

	// Enabled: Whether ABAC authorization will be enabled in the cluster.
	Enabled bool `json:"enabled,omitempty"`

	// Name: The name (project, location, cluster id) of the cluster to set
	// legacy abac.
	// Specified in the format 'projects/*/locations/*/clusters/*'.
	Name string `json:"name,omitempty"`

	// ProjectId: The Google Developers Console [project ID or
	// project
	// number](https://support.google.com/cloud/answer/6158840).
	// This
	//  field is deprecated, use name instead.
	ProjectId string `json:"projectId,omitempty"`

	// Zone: The name of the Google Compute
	// Engine
	// [zone](/compute/docs/zones#available) in which the
	// cluster
	// resides.
	// This field is deprecated, use name instead.
	Zone string `json:"zone,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ClusterId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ClusterId") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *SetLegacyAbacRequest) MarshalJSON() ([]byte, error) {
	type noMethod SetLegacyAbacRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// SetMaintenancePolicyRequest: SetMaintenancePolicyRequest sets the
// maintenance policy for a cluster.
type SetMaintenancePolicyRequest struct {
	// ClusterId: The name of the cluster to update.
	ClusterId string `json:"clusterId,omitempty"`

	// MaintenancePolicy: The maintenance policy to be set for the cluster.
	// An empty field
	// clears the existing maintenance policy.
	MaintenancePolicy *MaintenancePolicy `json:"maintenancePolicy,omitempty"`

	// Name: The name (project, location, cluster id) of the cluster to set
	// maintenance
	// policy.
	// Specified in the format 'projects/*/locations/*/clusters/*'.
	Name string `json:"name,omitempty"`

	// ProjectId: The Google Developers Console [project ID or
	// project
	// number](https://support.google.com/cloud/answer/6158840).
	ProjectId string `json:"projectId,omitempty"`

	// Zone: The name of the Google Compute
	// Engine
	// [zone](/compute/docs/zones#available) in which the cluster
	// resides.
	Zone string `json:"zone,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ClusterId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ClusterId") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *SetMaintenancePolicyRequest) MarshalJSON() ([]byte, error) {
	type noMethod SetMaintenancePolicyRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// SetMasterAuthRequest: SetMasterAuthRequest updates the admin password
// of a cluster.
type SetMasterAuthRequest struct {
	// Action: The exact form of action to be taken on the master auth.
	//
	// Possible values:
	//   "UNKNOWN" - Operation is unknown and will error out.
	//   "SET_PASSWORD" - Set the password to a user generated value.
	//   "GENERATE_PASSWORD" - Generate a new password and set it to that.
	//   "SET_USERNAME" - Set the username.  If an empty username is
	// provided, basic authentication
	// is disabled for the cluster.  If a non-empty username is provided,
	// basic
	// authentication is enabled, with either a provided password or a
	// generated
	// one.
	Action string `json:"action,omitempty"`

	// ClusterId: The name of the cluster to upgrade.
	// This field is deprecated, use name instead.
	ClusterId string `json:"clusterId,omitempty"`

	// Name: The name (project, location, cluster) of the cluster to set
	// auth.
	// Specified in the format 'projects/*/locations/*/clusters/*'.
	Name string `json:"name,omitempty"`

	// ProjectId: The Google Developers Console [project ID or
	// project
	// number](https://support.google.com/cloud/answer/6158840).
	// This
	//  field is deprecated, use name instead.
	ProjectId string `json:"projectId,omitempty"`

	// Update: A description of the update.
	Update *MasterAuth `json:"update,omitempty"`

	// Zone: The name of the Google Compute
	// Engine
	// [zone](/compute/docs/zones#available) in which the
	// cluster
	// resides.
	// This field is deprecated, use name instead.
	Zone string `json:"zone,omitempty"`

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

func (s *SetMasterAuthRequest) MarshalJSON() ([]byte, error) {
	type noMethod SetMasterAuthRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// SetNetworkPolicyRequest: SetNetworkPolicyRequest enables/disables
// network policy for a cluster.
type SetNetworkPolicyRequest struct {
	// ClusterId: The name of the cluster.
	// This field is deprecated, use name instead.
	ClusterId string `json:"clusterId,omitempty"`

	// Name: The name (project, location, cluster id) of the cluster to set
	// networking policy.
	// Specified in the format 'projects/*/locations/*/clusters/*'.
	Name string `json:"name,omitempty"`

	// NetworkPolicy: Configuration options for the NetworkPolicy feature.
	NetworkPolicy *NetworkPolicy `json:"networkPolicy,omitempty"`

	// ProjectId: The Google Developers Console [project ID or
	// project
	// number](https://developers.google.com/console/help/new/#projec
	// tnumber).
	// This field is deprecated, use name instead.
	ProjectId string `json:"projectId,omitempty"`

	// Zone: The name of the Google Compute
	// Engine
	// [zone](/compute/docs/zones#available) in which the
	// cluster
	// resides.
	// This field is deprecated, use name instead.
	Zone string `json:"zone,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ClusterId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ClusterId") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *SetNetworkPolicyRequest) MarshalJSON() ([]byte, error) {
	type noMethod SetNetworkPolicyRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// SetNodePoolManagementRequest: SetNodePoolManagementRequest sets the
// node management properties of a node
// pool.
type SetNodePoolManagementRequest struct {
	// ClusterId: The name of the cluster to update.
	// This field is deprecated, use name instead.
	ClusterId string `json:"clusterId,omitempty"`

	// Management: NodeManagement configuration for the node pool.
	Management *NodeManagement `json:"management,omitempty"`

	// Name: The name (project, location, cluster, node pool id) of the node
	// pool to set
	// management properties. Specified in the
	// format
	// 'projects/*/locations/*/clusters/*/nodePools/*'.
	Name string `json:"name,omitempty"`

	// NodePoolId: The name of the node pool to update.
	// This field is deprecated, use name instead.
	NodePoolId string `json:"nodePoolId,omitempty"`

	// ProjectId: The Google Developers Console [project ID or
	// project
	// number](https://support.google.com/cloud/answer/6158840).
	// This
	//  field is deprecated, use name instead.
	ProjectId string `json:"projectId,omitempty"`

	// Zone: The name of the Google Compute
	// Engine
	// [zone](/compute/docs/zones#available) in which the
	// cluster
	// resides.
	// This field is deprecated, use name instead.
	Zone string `json:"zone,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ClusterId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ClusterId") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *SetNodePoolManagementRequest) MarshalJSON() ([]byte, error) {
	type noMethod SetNodePoolManagementRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// StartIPRotationRequest: StartIPRotationRequest creates a new IP for
// the cluster and then performs
// a node upgrade on each node pool to point to the new IP.
type StartIPRotationRequest struct {
	// ClusterId: The name of the cluster.
	// This field is deprecated, use name instead.
	ClusterId string `json:"clusterId,omitempty"`

	// Name: The name (project, location, cluster id) of the cluster to
	// start IP rotation.
	// Specified in the format 'projects/*/locations/*/clusters/*'.
	Name string `json:"name,omitempty"`

	// ProjectId: The Google Developers Console [project ID or
	// project
	// number](https://developers.google.com/console/help/new/#projec
	// tnumber).
	// This field is deprecated, use name instead.
	ProjectId string `json:"projectId,omitempty"`

	// Zone: The name of the Google Compute
	// Engine
	// [zone](/compute/docs/zones#available) in which the
	// cluster
	// resides.
	// This field is deprecated, use name instead.
	Zone string `json:"zone,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ClusterId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ClusterId") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *StartIPRotationRequest) MarshalJSON() ([]byte, error) {
	type noMethod StartIPRotationRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// UpdateClusterRequest: UpdateClusterRequest updates the settings of a
// cluster.
type UpdateClusterRequest struct {
	// ClusterId: The name of the cluster to upgrade.
	// This field is deprecated, use name instead.
	ClusterId string `json:"clusterId,omitempty"`

	// Name: The name (project, location, cluster) of the cluster to
	// update.
	// Specified in the format 'projects/*/locations/*/clusters/*'.
	Name string `json:"name,omitempty"`

	// ProjectId: The Google Developers Console [project ID or
	// project
	// number](https://support.google.com/cloud/answer/6158840).
	// This
	//  field is deprecated, use name instead.
	ProjectId string `json:"projectId,omitempty"`

	// Update: A description of the update.
	Update *ClusterUpdate `json:"update,omitempty"`

	// Zone: The name of the Google Compute
	// Engine
	// [zone](/compute/docs/zones#available) in which the
	// cluster
	// resides.
	// This field is deprecated, use name instead.
	Zone string `json:"zone,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ClusterId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ClusterId") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *UpdateClusterRequest) MarshalJSON() ([]byte, error) {
	type noMethod UpdateClusterRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// method id "container.projects.locations.getServerConfig":

type ProjectsLocationsGetServerConfigCall struct {
	s            *Service
	name         string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// GetServerConfig: Returns configuration info about the Container
// Engine service.
func (r *ProjectsLocationsService) GetServerConfig(name string) *ProjectsLocationsGetServerConfigCall {
	c := &ProjectsLocationsGetServerConfigCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	return c
}

// ProjectId sets the optional parameter "projectId": The Google
// Developers Console [project ID or
// project
// number](https://support.google.com/cloud/answer/6158840).
// This
//  field is deprecated, use name instead.
func (c *ProjectsLocationsGetServerConfigCall) ProjectId(projectId string) *ProjectsLocationsGetServerConfigCall {
	c.urlParams_.Set("projectId", projectId)
	return c
}

// Zone sets the optional parameter "zone": The name of the Google
// Compute Engine [zone](/compute/docs/zones#available)
// to return operations for.
// This field is deprecated, use name instead.
func (c *ProjectsLocationsGetServerConfigCall) Zone(zone string) *ProjectsLocationsGetServerConfigCall {
	c.urlParams_.Set("zone", zone)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsGetServerConfigCall) Fields(s ...googleapi.Field) *ProjectsLocationsGetServerConfigCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsLocationsGetServerConfigCall) IfNoneMatch(entityTag string) *ProjectsLocationsGetServerConfigCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsGetServerConfigCall) Context(ctx context.Context) *ProjectsLocationsGetServerConfigCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsGetServerConfigCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsGetServerConfigCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}/serverConfig")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.getServerConfig" call.
// Exactly one of *ServerConfig or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *ServerConfig.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsGetServerConfigCall) Do(opts ...googleapi.CallOption) (*ServerConfig, error) {
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
	ret := &ServerConfig{
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
	//   "description": "Returns configuration info about the Container Engine service.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/serverConfig",
	//   "httpMethod": "GET",
	//   "id": "container.projects.locations.getServerConfig",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The name (project and location) of the server config to get\nSpecified in the format 'projects/*/locations/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine [zone](/compute/docs/zones#available)\nto return operations for.\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}/serverConfig",
	//   "response": {
	//     "$ref": "ServerConfig"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.completeIpRotation":

type ProjectsLocationsClustersCompleteIpRotationCall struct {
	s                         *Service
	name                      string
	completeiprotationrequest *CompleteIPRotationRequest
	urlParams_                gensupport.URLParams
	ctx_                      context.Context
	header_                   http.Header
}

// CompleteIpRotation: Completes master IP rotation.
func (r *ProjectsLocationsClustersService) CompleteIpRotation(name string, completeiprotationrequest *CompleteIPRotationRequest) *ProjectsLocationsClustersCompleteIpRotationCall {
	c := &ProjectsLocationsClustersCompleteIpRotationCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.completeiprotationrequest = completeiprotationrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersCompleteIpRotationCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersCompleteIpRotationCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersCompleteIpRotationCall) Context(ctx context.Context) *ProjectsLocationsClustersCompleteIpRotationCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersCompleteIpRotationCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersCompleteIpRotationCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.completeiprotationrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:completeIpRotation")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.completeIpRotation" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersCompleteIpRotationCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Completes master IP rotation.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters/{clustersId}:completeIpRotation",
	//   "httpMethod": "POST",
	//   "id": "container.projects.locations.clusters.completeIpRotation",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The name (project, location, cluster id) of the cluster to complete IP rotation.\nSpecified in the format 'projects/*/locations/*/clusters/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/clusters/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}:completeIpRotation",
	//   "request": {
	//     "$ref": "CompleteIPRotationRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.create":

type ProjectsLocationsClustersCreateCall struct {
	s                    *Service
	parent               string
	createclusterrequest *CreateClusterRequest
	urlParams_           gensupport.URLParams
	ctx_                 context.Context
	header_              http.Header
}

// Create: Creates a cluster, consisting of the specified number and
// type of Google
// Compute Engine instances.
//
// By default, the cluster is created in the project's
// [default
// network](/compute/docs/networks-and-firewalls#networks).
//
// One firewall is added for the cluster. After cluster creation,
// the cluster creates routes for each node to allow the containers
// on that node to communicate with all other instances in
// the
// cluster.
//
// Finally, an entry is added to the project's global metadata
// indicating
// which CIDR range is being used by the cluster.
func (r *ProjectsLocationsClustersService) Create(parent string, createclusterrequest *CreateClusterRequest) *ProjectsLocationsClustersCreateCall {
	c := &ProjectsLocationsClustersCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	c.createclusterrequest = createclusterrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersCreateCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersCreateCall) Context(ctx context.Context) *ProjectsLocationsClustersCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.createclusterrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+parent}/clusters")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"parent": c.parent,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.create" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersCreateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Creates a cluster, consisting of the specified number and type of Google\nCompute Engine instances.\n\nBy default, the cluster is created in the project's\n[default network](/compute/docs/networks-and-firewalls#networks).\n\nOne firewall is added for the cluster. After cluster creation,\nthe cluster creates routes for each node to allow the containers\non that node to communicate with all other instances in the\ncluster.\n\nFinally, an entry is added to the project's global metadata indicating\nwhich CIDR range is being used by the cluster.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters",
	//   "httpMethod": "POST",
	//   "id": "container.projects.locations.clusters.create",
	//   "parameterOrder": [
	//     "parent"
	//   ],
	//   "parameters": {
	//     "parent": {
	//       "description": "The parent (project and location) where the cluster will be created.\nSpecified in the format 'projects/*/locations/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+parent}/clusters",
	//   "request": {
	//     "$ref": "CreateClusterRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.delete":

type ProjectsLocationsClustersDeleteCall struct {
	s          *Service
	name       string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Delete: Deletes the cluster, including the Kubernetes endpoint and
// all worker
// nodes.
//
// Firewalls and routes that were configured during cluster creation
// are also deleted.
//
// Other Google Compute Engine resources that might be in use by the
// cluster
// (e.g. load balancer resources) will not be deleted if they weren't
// present
// at the initial create time.
func (r *ProjectsLocationsClustersService) Delete(name string) *ProjectsLocationsClustersDeleteCall {
	c := &ProjectsLocationsClustersDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	return c
}

// ClusterId sets the optional parameter "clusterId": The name of the
// cluster to delete.
// This field is deprecated, use name instead.
func (c *ProjectsLocationsClustersDeleteCall) ClusterId(clusterId string) *ProjectsLocationsClustersDeleteCall {
	c.urlParams_.Set("clusterId", clusterId)
	return c
}

// ProjectId sets the optional parameter "projectId": The Google
// Developers Console [project ID or
// project
// number](https://support.google.com/cloud/answer/6158840).
// This
//  field is deprecated, use name instead.
func (c *ProjectsLocationsClustersDeleteCall) ProjectId(projectId string) *ProjectsLocationsClustersDeleteCall {
	c.urlParams_.Set("projectId", projectId)
	return c
}

// Zone sets the optional parameter "zone": The name of the Google
// Compute Engine
// [zone](/compute/docs/zones#available) in which the
// cluster
// resides.
// This field is deprecated, use name instead.
func (c *ProjectsLocationsClustersDeleteCall) Zone(zone string) *ProjectsLocationsClustersDeleteCall {
	c.urlParams_.Set("zone", zone)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersDeleteCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersDeleteCall) Context(ctx context.Context) *ProjectsLocationsClustersDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.delete" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersDeleteCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Deletes the cluster, including the Kubernetes endpoint and all worker\nnodes.\n\nFirewalls and routes that were configured during cluster creation\nare also deleted.\n\nOther Google Compute Engine resources that might be in use by the cluster\n(e.g. load balancer resources) will not be deleted if they weren't present\nat the initial create time.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters/{clustersId}",
	//   "httpMethod": "DELETE",
	//   "id": "container.projects.locations.clusters.delete",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster to delete.\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "name": {
	//       "description": "The name (project, location, cluster) of the cluster to delete.\nSpecified in the format 'projects/*/locations/*/clusters/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/clusters/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}",
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.get":

type ProjectsLocationsClustersGetCall struct {
	s            *Service
	name         string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Gets the details of a specific cluster.
func (r *ProjectsLocationsClustersService) Get(name string) *ProjectsLocationsClustersGetCall {
	c := &ProjectsLocationsClustersGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	return c
}

// ClusterId sets the optional parameter "clusterId": The name of the
// cluster to retrieve.
// This field is deprecated, use name instead.
func (c *ProjectsLocationsClustersGetCall) ClusterId(clusterId string) *ProjectsLocationsClustersGetCall {
	c.urlParams_.Set("clusterId", clusterId)
	return c
}

// ProjectId sets the optional parameter "projectId": The Google
// Developers Console [project ID or
// project
// number](https://support.google.com/cloud/answer/6158840).
// This
//  field is deprecated, use name instead.
func (c *ProjectsLocationsClustersGetCall) ProjectId(projectId string) *ProjectsLocationsClustersGetCall {
	c.urlParams_.Set("projectId", projectId)
	return c
}

// Zone sets the optional parameter "zone": The name of the Google
// Compute Engine
// [zone](/compute/docs/zones#available) in which the
// cluster
// resides.
// This field is deprecated, use name instead.
func (c *ProjectsLocationsClustersGetCall) Zone(zone string) *ProjectsLocationsClustersGetCall {
	c.urlParams_.Set("zone", zone)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersGetCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsLocationsClustersGetCall) IfNoneMatch(entityTag string) *ProjectsLocationsClustersGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersGetCall) Context(ctx context.Context) *ProjectsLocationsClustersGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.get" call.
// Exactly one of *Cluster or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Cluster.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *ProjectsLocationsClustersGetCall) Do(opts ...googleapi.CallOption) (*Cluster, error) {
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
	ret := &Cluster{
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
	//   "description": "Gets the details of a specific cluster.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters/{clustersId}",
	//   "httpMethod": "GET",
	//   "id": "container.projects.locations.clusters.get",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster to retrieve.\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "name": {
	//       "description": "The name (project, location, cluster) of the cluster to retrieve.\nSpecified in the format 'projects/*/locations/*/clusters/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/clusters/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}",
	//   "response": {
	//     "$ref": "Cluster"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.list":

type ProjectsLocationsClustersListCall struct {
	s            *Service
	parent       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists all clusters owned by a project in either the specified
// zone or all
// zones.
func (r *ProjectsLocationsClustersService) List(parent string) *ProjectsLocationsClustersListCall {
	c := &ProjectsLocationsClustersListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	return c
}

// ProjectId sets the optional parameter "projectId": The Google
// Developers Console [project ID or
// project
// number](https://support.google.com/cloud/answer/6158840).
// This
//  field is deprecated, use parent instead.
func (c *ProjectsLocationsClustersListCall) ProjectId(projectId string) *ProjectsLocationsClustersListCall {
	c.urlParams_.Set("projectId", projectId)
	return c
}

// Zone sets the optional parameter "zone": The name of the Google
// Compute Engine
// [zone](/compute/docs/zones#available) in which the cluster
// resides, or "-" for all zones.
// This field is deprecated, use parent instead.
func (c *ProjectsLocationsClustersListCall) Zone(zone string) *ProjectsLocationsClustersListCall {
	c.urlParams_.Set("zone", zone)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersListCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsLocationsClustersListCall) IfNoneMatch(entityTag string) *ProjectsLocationsClustersListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersListCall) Context(ctx context.Context) *ProjectsLocationsClustersListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+parent}/clusters")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"parent": c.parent,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.list" call.
// Exactly one of *ListClustersResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListClustersResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersListCall) Do(opts ...googleapi.CallOption) (*ListClustersResponse, error) {
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
	ret := &ListClustersResponse{
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
	//   "description": "Lists all clusters owned by a project in either the specified zone or all\nzones.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters",
	//   "httpMethod": "GET",
	//   "id": "container.projects.locations.clusters.list",
	//   "parameterOrder": [
	//     "parent"
	//   ],
	//   "parameters": {
	//     "parent": {
	//       "description": "The parent (project and location) where the clusters will be listed.\nSpecified in the format 'projects/*/locations/*'.\nLocation \"-\" matches all zones and all regions.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use parent instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides, or \"-\" for all zones.\nThis field is deprecated, use parent instead.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+parent}/clusters",
	//   "response": {
	//     "$ref": "ListClustersResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.setLegacyAbac":

type ProjectsLocationsClustersSetLegacyAbacCall struct {
	s                    *Service
	name                 string
	setlegacyabacrequest *SetLegacyAbacRequest
	urlParams_           gensupport.URLParams
	ctx_                 context.Context
	header_              http.Header
}

// SetLegacyAbac: Enables or disables the ABAC authorization mechanism
// on a cluster.
func (r *ProjectsLocationsClustersService) SetLegacyAbac(name string, setlegacyabacrequest *SetLegacyAbacRequest) *ProjectsLocationsClustersSetLegacyAbacCall {
	c := &ProjectsLocationsClustersSetLegacyAbacCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setlegacyabacrequest = setlegacyabacrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersSetLegacyAbacCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersSetLegacyAbacCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersSetLegacyAbacCall) Context(ctx context.Context) *ProjectsLocationsClustersSetLegacyAbacCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersSetLegacyAbacCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersSetLegacyAbacCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setlegacyabacrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setLegacyAbac")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.setLegacyAbac" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersSetLegacyAbacCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Enables or disables the ABAC authorization mechanism on a cluster.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters/{clustersId}:setLegacyAbac",
	//   "httpMethod": "POST",
	//   "id": "container.projects.locations.clusters.setLegacyAbac",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The name (project, location, cluster id) of the cluster to set legacy abac.\nSpecified in the format 'projects/*/locations/*/clusters/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/clusters/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}:setLegacyAbac",
	//   "request": {
	//     "$ref": "SetLegacyAbacRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.setMaintenancePolicy":

type ProjectsLocationsClustersSetMaintenancePolicyCall struct {
	s                           *Service
	name                        string
	setmaintenancepolicyrequest *SetMaintenancePolicyRequest
	urlParams_                  gensupport.URLParams
	ctx_                        context.Context
	header_                     http.Header
}

// SetMaintenancePolicy: Sets the maintenance policy for a cluster.
func (r *ProjectsLocationsClustersService) SetMaintenancePolicy(name string, setmaintenancepolicyrequest *SetMaintenancePolicyRequest) *ProjectsLocationsClustersSetMaintenancePolicyCall {
	c := &ProjectsLocationsClustersSetMaintenancePolicyCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setmaintenancepolicyrequest = setmaintenancepolicyrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersSetMaintenancePolicyCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersSetMaintenancePolicyCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersSetMaintenancePolicyCall) Context(ctx context.Context) *ProjectsLocationsClustersSetMaintenancePolicyCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersSetMaintenancePolicyCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersSetMaintenancePolicyCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setmaintenancepolicyrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setMaintenancePolicy")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.setMaintenancePolicy" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersSetMaintenancePolicyCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Sets the maintenance policy for a cluster.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters/{clustersId}:setMaintenancePolicy",
	//   "httpMethod": "POST",
	//   "id": "container.projects.locations.clusters.setMaintenancePolicy",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The name (project, location, cluster id) of the cluster to set maintenance\npolicy.\nSpecified in the format 'projects/*/locations/*/clusters/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/clusters/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}:setMaintenancePolicy",
	//   "request": {
	//     "$ref": "SetMaintenancePolicyRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.setMasterAuth":

type ProjectsLocationsClustersSetMasterAuthCall struct {
	s                    *Service
	name                 string
	setmasterauthrequest *SetMasterAuthRequest
	urlParams_           gensupport.URLParams
	ctx_                 context.Context
	header_              http.Header
}

// SetMasterAuth: Used to set master auth materials. Currently supports
// :-
// Changing the admin password of a specific cluster.
// This can be either via password generation or explicitly set.
// Modify basic_auth.csv and reset the K8S API server.
func (r *ProjectsLocationsClustersService) SetMasterAuth(name string, setmasterauthrequest *SetMasterAuthRequest) *ProjectsLocationsClustersSetMasterAuthCall {
	c := &ProjectsLocationsClustersSetMasterAuthCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setmasterauthrequest = setmasterauthrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersSetMasterAuthCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersSetMasterAuthCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersSetMasterAuthCall) Context(ctx context.Context) *ProjectsLocationsClustersSetMasterAuthCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersSetMasterAuthCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersSetMasterAuthCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setmasterauthrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setMasterAuth")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.setMasterAuth" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersSetMasterAuthCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Used to set master auth materials. Currently supports :-\nChanging the admin password of a specific cluster.\nThis can be either via password generation or explicitly set.\nModify basic_auth.csv and reset the K8S API server.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters/{clustersId}:setMasterAuth",
	//   "httpMethod": "POST",
	//   "id": "container.projects.locations.clusters.setMasterAuth",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The name (project, location, cluster) of the cluster to set auth.\nSpecified in the format 'projects/*/locations/*/clusters/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/clusters/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}:setMasterAuth",
	//   "request": {
	//     "$ref": "SetMasterAuthRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.setNetworkPolicy":

type ProjectsLocationsClustersSetNetworkPolicyCall struct {
	s                       *Service
	name                    string
	setnetworkpolicyrequest *SetNetworkPolicyRequest
	urlParams_              gensupport.URLParams
	ctx_                    context.Context
	header_                 http.Header
}

// SetNetworkPolicy: Enables/Disables Network Policy for a cluster.
func (r *ProjectsLocationsClustersService) SetNetworkPolicy(name string, setnetworkpolicyrequest *SetNetworkPolicyRequest) *ProjectsLocationsClustersSetNetworkPolicyCall {
	c := &ProjectsLocationsClustersSetNetworkPolicyCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setnetworkpolicyrequest = setnetworkpolicyrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersSetNetworkPolicyCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersSetNetworkPolicyCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersSetNetworkPolicyCall) Context(ctx context.Context) *ProjectsLocationsClustersSetNetworkPolicyCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersSetNetworkPolicyCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersSetNetworkPolicyCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setnetworkpolicyrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setNetworkPolicy")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.setNetworkPolicy" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersSetNetworkPolicyCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Enables/Disables Network Policy for a cluster.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters/{clustersId}:setNetworkPolicy",
	//   "httpMethod": "POST",
	//   "id": "container.projects.locations.clusters.setNetworkPolicy",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The name (project, location, cluster id) of the cluster to set networking policy.\nSpecified in the format 'projects/*/locations/*/clusters/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/clusters/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}:setNetworkPolicy",
	//   "request": {
	//     "$ref": "SetNetworkPolicyRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.setResourceLabels":

type ProjectsLocationsClustersSetResourceLabelsCall struct {
	s                *Service
	name             string
	setlabelsrequest *SetLabelsRequest
	urlParams_       gensupport.URLParams
	ctx_             context.Context
	header_          http.Header
}

// SetResourceLabels: Sets labels on a cluster.
func (r *ProjectsLocationsClustersService) SetResourceLabels(name string, setlabelsrequest *SetLabelsRequest) *ProjectsLocationsClustersSetResourceLabelsCall {
	c := &ProjectsLocationsClustersSetResourceLabelsCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setlabelsrequest = setlabelsrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersSetResourceLabelsCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersSetResourceLabelsCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersSetResourceLabelsCall) Context(ctx context.Context) *ProjectsLocationsClustersSetResourceLabelsCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersSetResourceLabelsCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersSetResourceLabelsCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setlabelsrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setResourceLabels")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.setResourceLabels" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersSetResourceLabelsCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Sets labels on a cluster.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters/{clustersId}:setResourceLabels",
	//   "httpMethod": "POST",
	//   "id": "container.projects.locations.clusters.setResourceLabels",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The name (project, location, cluster id) of the cluster to set labels.\nSpecified in the format 'projects/*/locations/*/clusters/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/clusters/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}:setResourceLabels",
	//   "request": {
	//     "$ref": "SetLabelsRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.startIpRotation":

type ProjectsLocationsClustersStartIpRotationCall struct {
	s                      *Service
	name                   string
	startiprotationrequest *StartIPRotationRequest
	urlParams_             gensupport.URLParams
	ctx_                   context.Context
	header_                http.Header
}

// StartIpRotation: Start master IP rotation.
func (r *ProjectsLocationsClustersService) StartIpRotation(name string, startiprotationrequest *StartIPRotationRequest) *ProjectsLocationsClustersStartIpRotationCall {
	c := &ProjectsLocationsClustersStartIpRotationCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.startiprotationrequest = startiprotationrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersStartIpRotationCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersStartIpRotationCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersStartIpRotationCall) Context(ctx context.Context) *ProjectsLocationsClustersStartIpRotationCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersStartIpRotationCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersStartIpRotationCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.startiprotationrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:startIpRotation")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.startIpRotation" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersStartIpRotationCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Start master IP rotation.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters/{clustersId}:startIpRotation",
	//   "httpMethod": "POST",
	//   "id": "container.projects.locations.clusters.startIpRotation",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The name (project, location, cluster id) of the cluster to start IP rotation.\nSpecified in the format 'projects/*/locations/*/clusters/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/clusters/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}:startIpRotation",
	//   "request": {
	//     "$ref": "StartIPRotationRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.update":

type ProjectsLocationsClustersUpdateCall struct {
	s                    *Service
	name                 string
	updateclusterrequest *UpdateClusterRequest
	urlParams_           gensupport.URLParams
	ctx_                 context.Context
	header_              http.Header
}

// Update: Updates the settings of a specific cluster.
func (r *ProjectsLocationsClustersService) Update(name string, updateclusterrequest *UpdateClusterRequest) *ProjectsLocationsClustersUpdateCall {
	c := &ProjectsLocationsClustersUpdateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.updateclusterrequest = updateclusterrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersUpdateCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersUpdateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersUpdateCall) Context(ctx context.Context) *ProjectsLocationsClustersUpdateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersUpdateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersUpdateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.updateclusterrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PUT", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.update" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersUpdateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Updates the settings of a specific cluster.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters/{clustersId}",
	//   "httpMethod": "PUT",
	//   "id": "container.projects.locations.clusters.update",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The name (project, location, cluster) of the cluster to update.\nSpecified in the format 'projects/*/locations/*/clusters/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/clusters/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}",
	//   "request": {
	//     "$ref": "UpdateClusterRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.nodePools.create":

type ProjectsLocationsClustersNodePoolsCreateCall struct {
	s                     *Service
	parent                string
	createnodepoolrequest *CreateNodePoolRequest
	urlParams_            gensupport.URLParams
	ctx_                  context.Context
	header_               http.Header
}

// Create: Creates a node pool for a cluster.
func (r *ProjectsLocationsClustersNodePoolsService) Create(parent string, createnodepoolrequest *CreateNodePoolRequest) *ProjectsLocationsClustersNodePoolsCreateCall {
	c := &ProjectsLocationsClustersNodePoolsCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	c.createnodepoolrequest = createnodepoolrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersNodePoolsCreateCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersNodePoolsCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersNodePoolsCreateCall) Context(ctx context.Context) *ProjectsLocationsClustersNodePoolsCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersNodePoolsCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersNodePoolsCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.createnodepoolrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+parent}/nodePools")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"parent": c.parent,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.nodePools.create" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersNodePoolsCreateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Creates a node pool for a cluster.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters/{clustersId}/nodePools",
	//   "httpMethod": "POST",
	//   "id": "container.projects.locations.clusters.nodePools.create",
	//   "parameterOrder": [
	//     "parent"
	//   ],
	//   "parameters": {
	//     "parent": {
	//       "description": "The parent (project, location, cluster id) where the node pool will be created.\nSpecified in the format 'projects/*/locations/*/clusters/*/nodePools/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/clusters/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+parent}/nodePools",
	//   "request": {
	//     "$ref": "CreateNodePoolRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.nodePools.delete":

type ProjectsLocationsClustersNodePoolsDeleteCall struct {
	s          *Service
	name       string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Delete: Deletes a node pool from a cluster.
func (r *ProjectsLocationsClustersNodePoolsService) Delete(name string) *ProjectsLocationsClustersNodePoolsDeleteCall {
	c := &ProjectsLocationsClustersNodePoolsDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	return c
}

// ClusterId sets the optional parameter "clusterId": The name of the
// cluster.
// This field is deprecated, use name instead.
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) ClusterId(clusterId string) *ProjectsLocationsClustersNodePoolsDeleteCall {
	c.urlParams_.Set("clusterId", clusterId)
	return c
}

// NodePoolId sets the optional parameter "nodePoolId": The name of the
// node pool to delete.
// This field is deprecated, use name instead.
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) NodePoolId(nodePoolId string) *ProjectsLocationsClustersNodePoolsDeleteCall {
	c.urlParams_.Set("nodePoolId", nodePoolId)
	return c
}

// ProjectId sets the optional parameter "projectId": The Google
// Developers Console [project ID or
// project
// number](https://developers.google.com/console/help/new/#projec
// tnumber).
// This field is deprecated, use name instead.
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) ProjectId(projectId string) *ProjectsLocationsClustersNodePoolsDeleteCall {
	c.urlParams_.Set("projectId", projectId)
	return c
}

// Zone sets the optional parameter "zone": The name of the Google
// Compute Engine
// [zone](/compute/docs/zones#available) in which the
// cluster
// resides.
// This field is deprecated, use name instead.
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) Zone(zone string) *ProjectsLocationsClustersNodePoolsDeleteCall {
	c.urlParams_.Set("zone", zone)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersNodePoolsDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) Context(ctx context.Context) *ProjectsLocationsClustersNodePoolsDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersNodePoolsDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.nodePools.delete" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersNodePoolsDeleteCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Deletes a node pool from a cluster.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters/{clustersId}/nodePools/{nodePoolsId}",
	//   "httpMethod": "DELETE",
	//   "id": "container.projects.locations.clusters.nodePools.delete",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster.\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "name": {
	//       "description": "The name (project, location, cluster, node pool id) of the node pool to delete.\nSpecified in the format 'projects/*/locations/*/clusters/*/nodePools/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/clusters/[^/]+/nodePools/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "nodePoolId": {
	//       "description": "The name of the node pool to delete.\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://developers.google.com/console/help/new/#projectnumber).\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}",
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.nodePools.get":

type ProjectsLocationsClustersNodePoolsGetCall struct {
	s            *Service
	name         string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Retrieves the node pool requested.
func (r *ProjectsLocationsClustersNodePoolsService) Get(name string) *ProjectsLocationsClustersNodePoolsGetCall {
	c := &ProjectsLocationsClustersNodePoolsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	return c
}

// ClusterId sets the optional parameter "clusterId": The name of the
// cluster.
// This field is deprecated, use name instead.
func (c *ProjectsLocationsClustersNodePoolsGetCall) ClusterId(clusterId string) *ProjectsLocationsClustersNodePoolsGetCall {
	c.urlParams_.Set("clusterId", clusterId)
	return c
}

// NodePoolId sets the optional parameter "nodePoolId": The name of the
// node pool.
// This field is deprecated, use name instead.
func (c *ProjectsLocationsClustersNodePoolsGetCall) NodePoolId(nodePoolId string) *ProjectsLocationsClustersNodePoolsGetCall {
	c.urlParams_.Set("nodePoolId", nodePoolId)
	return c
}

// ProjectId sets the optional parameter "projectId": The Google
// Developers Console [project ID or
// project
// number](https://developers.google.com/console/help/new/#projec
// tnumber).
// This field is deprecated, use name instead.
func (c *ProjectsLocationsClustersNodePoolsGetCall) ProjectId(projectId string) *ProjectsLocationsClustersNodePoolsGetCall {
	c.urlParams_.Set("projectId", projectId)
	return c
}

// Zone sets the optional parameter "zone": The name of the Google
// Compute Engine
// [zone](/compute/docs/zones#available) in which the
// cluster
// resides.
// This field is deprecated, use name instead.
func (c *ProjectsLocationsClustersNodePoolsGetCall) Zone(zone string) *ProjectsLocationsClustersNodePoolsGetCall {
	c.urlParams_.Set("zone", zone)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersNodePoolsGetCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersNodePoolsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsLocationsClustersNodePoolsGetCall) IfNoneMatch(entityTag string) *ProjectsLocationsClustersNodePoolsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersNodePoolsGetCall) Context(ctx context.Context) *ProjectsLocationsClustersNodePoolsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersNodePoolsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersNodePoolsGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.nodePools.get" call.
// Exactly one of *NodePool or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *NodePool.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersNodePoolsGetCall) Do(opts ...googleapi.CallOption) (*NodePool, error) {
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
	ret := &NodePool{
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
	//   "description": "Retrieves the node pool requested.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters/{clustersId}/nodePools/{nodePoolsId}",
	//   "httpMethod": "GET",
	//   "id": "container.projects.locations.clusters.nodePools.get",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster.\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "name": {
	//       "description": "The name (project, location, cluster, node pool id) of the node pool to get.\nSpecified in the format 'projects/*/locations/*/clusters/*/nodePools/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/clusters/[^/]+/nodePools/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "nodePoolId": {
	//       "description": "The name of the node pool.\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://developers.google.com/console/help/new/#projectnumber).\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}",
	//   "response": {
	//     "$ref": "NodePool"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.nodePools.list":

type ProjectsLocationsClustersNodePoolsListCall struct {
	s            *Service
	parent       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists the node pools for a cluster.
func (r *ProjectsLocationsClustersNodePoolsService) List(parent string) *ProjectsLocationsClustersNodePoolsListCall {
	c := &ProjectsLocationsClustersNodePoolsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	return c
}

// ClusterId sets the optional parameter "clusterId": The name of the
// cluster.
// This field is deprecated, use parent instead.
func (c *ProjectsLocationsClustersNodePoolsListCall) ClusterId(clusterId string) *ProjectsLocationsClustersNodePoolsListCall {
	c.urlParams_.Set("clusterId", clusterId)
	return c
}

// ProjectId sets the optional parameter "projectId": The Google
// Developers Console [project ID or
// project
// number](https://developers.google.com/console/help/new/#projec
// tnumber).
// This field is deprecated, use parent instead.
func (c *ProjectsLocationsClustersNodePoolsListCall) ProjectId(projectId string) *ProjectsLocationsClustersNodePoolsListCall {
	c.urlParams_.Set("projectId", projectId)
	return c
}

// Zone sets the optional parameter "zone": The name of the Google
// Compute Engine
// [zone](/compute/docs/zones#available) in which the
// cluster
// resides.
// This field is deprecated, use parent instead.
func (c *ProjectsLocationsClustersNodePoolsListCall) Zone(zone string) *ProjectsLocationsClustersNodePoolsListCall {
	c.urlParams_.Set("zone", zone)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersNodePoolsListCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersNodePoolsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsLocationsClustersNodePoolsListCall) IfNoneMatch(entityTag string) *ProjectsLocationsClustersNodePoolsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersNodePoolsListCall) Context(ctx context.Context) *ProjectsLocationsClustersNodePoolsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersNodePoolsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersNodePoolsListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+parent}/nodePools")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"parent": c.parent,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.nodePools.list" call.
// Exactly one of *ListNodePoolsResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListNodePoolsResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersNodePoolsListCall) Do(opts ...googleapi.CallOption) (*ListNodePoolsResponse, error) {
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
	ret := &ListNodePoolsResponse{
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
	//   "description": "Lists the node pools for a cluster.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters/{clustersId}/nodePools",
	//   "httpMethod": "GET",
	//   "id": "container.projects.locations.clusters.nodePools.list",
	//   "parameterOrder": [
	//     "parent"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster.\nThis field is deprecated, use parent instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "parent": {
	//       "description": "The parent (project, location, cluster id) where the node pools will be listed.\nSpecified in the format 'projects/*/locations/*/clusters/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/clusters/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://developers.google.com/console/help/new/#projectnumber).\nThis field is deprecated, use parent instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use parent instead.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+parent}/nodePools",
	//   "response": {
	//     "$ref": "ListNodePoolsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.nodePools.rollback":

type ProjectsLocationsClustersNodePoolsRollbackCall struct {
	s                              *Service
	name                           string
	rollbacknodepoolupgraderequest *RollbackNodePoolUpgradeRequest
	urlParams_                     gensupport.URLParams
	ctx_                           context.Context
	header_                        http.Header
}

// Rollback: Roll back the previously Aborted or Failed NodePool
// upgrade.
// This will be an no-op if the last upgrade successfully completed.
func (r *ProjectsLocationsClustersNodePoolsService) Rollback(name string, rollbacknodepoolupgraderequest *RollbackNodePoolUpgradeRequest) *ProjectsLocationsClustersNodePoolsRollbackCall {
	c := &ProjectsLocationsClustersNodePoolsRollbackCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.rollbacknodepoolupgraderequest = rollbacknodepoolupgraderequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersNodePoolsRollbackCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersNodePoolsRollbackCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersNodePoolsRollbackCall) Context(ctx context.Context) *ProjectsLocationsClustersNodePoolsRollbackCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersNodePoolsRollbackCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersNodePoolsRollbackCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.rollbacknodepoolupgraderequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:rollback")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.nodePools.rollback" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersNodePoolsRollbackCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Roll back the previously Aborted or Failed NodePool upgrade.\nThis will be an no-op if the last upgrade successfully completed.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters/{clustersId}/nodePools/{nodePoolsId}:rollback",
	//   "httpMethod": "POST",
	//   "id": "container.projects.locations.clusters.nodePools.rollback",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The name (project, location, cluster, node pool id) of the node poll to\nrollback upgrade.\nSpecified in the format 'projects/*/locations/*/clusters/*/nodePools/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/clusters/[^/]+/nodePools/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}:rollback",
	//   "request": {
	//     "$ref": "RollbackNodePoolUpgradeRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.clusters.nodePools.setManagement":

type ProjectsLocationsClustersNodePoolsSetManagementCall struct {
	s                            *Service
	name                         string
	setnodepoolmanagementrequest *SetNodePoolManagementRequest
	urlParams_                   gensupport.URLParams
	ctx_                         context.Context
	header_                      http.Header
}

// SetManagement: Sets the NodeManagement options for a node pool.
func (r *ProjectsLocationsClustersNodePoolsService) SetManagement(name string, setnodepoolmanagementrequest *SetNodePoolManagementRequest) *ProjectsLocationsClustersNodePoolsSetManagementCall {
	c := &ProjectsLocationsClustersNodePoolsSetManagementCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.setnodepoolmanagementrequest = setnodepoolmanagementrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsClustersNodePoolsSetManagementCall) Fields(s ...googleapi.Field) *ProjectsLocationsClustersNodePoolsSetManagementCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsClustersNodePoolsSetManagementCall) Context(ctx context.Context) *ProjectsLocationsClustersNodePoolsSetManagementCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsClustersNodePoolsSetManagementCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsClustersNodePoolsSetManagementCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setnodepoolmanagementrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:setManagement")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.clusters.nodePools.setManagement" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsClustersNodePoolsSetManagementCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Sets the NodeManagement options for a node pool.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/clusters/{clustersId}/nodePools/{nodePoolsId}:setManagement",
	//   "httpMethod": "POST",
	//   "id": "container.projects.locations.clusters.nodePools.setManagement",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The name (project, location, cluster, node pool id) of the node pool to set\nmanagement properties. Specified in the format\n'projects/*/locations/*/clusters/*/nodePools/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/clusters/[^/]+/nodePools/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}:setManagement",
	//   "request": {
	//     "$ref": "SetNodePoolManagementRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.operations.cancel":

type ProjectsLocationsOperationsCancelCall struct {
	s                      *Service
	name                   string
	canceloperationrequest *CancelOperationRequest
	urlParams_             gensupport.URLParams
	ctx_                   context.Context
	header_                http.Header
}

// Cancel: Cancels the specified operation.
func (r *ProjectsLocationsOperationsService) Cancel(name string, canceloperationrequest *CancelOperationRequest) *ProjectsLocationsOperationsCancelCall {
	c := &ProjectsLocationsOperationsCancelCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	c.canceloperationrequest = canceloperationrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsOperationsCancelCall) Fields(s ...googleapi.Field) *ProjectsLocationsOperationsCancelCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsOperationsCancelCall) Context(ctx context.Context) *ProjectsLocationsOperationsCancelCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsOperationsCancelCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsOperationsCancelCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.canceloperationrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}:cancel")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.operations.cancel" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *ProjectsLocationsOperationsCancelCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
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
	//   "description": "Cancels the specified operation.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/operations/{operationsId}:cancel",
	//   "httpMethod": "POST",
	//   "id": "container.projects.locations.operations.cancel",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The name (project, location, operation id) of the operation to cancel.\nSpecified in the format 'projects/*/locations/*/operations/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/operations/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}:cancel",
	//   "request": {
	//     "$ref": "CancelOperationRequest"
	//   },
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.operations.get":

type ProjectsLocationsOperationsGetCall struct {
	s            *Service
	name         string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Gets the specified operation.
func (r *ProjectsLocationsOperationsService) Get(name string) *ProjectsLocationsOperationsGetCall {
	c := &ProjectsLocationsOperationsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.name = name
	return c
}

// OperationId sets the optional parameter "operationId": The
// server-assigned `name` of the operation.
// This field is deprecated, use name instead.
func (c *ProjectsLocationsOperationsGetCall) OperationId(operationId string) *ProjectsLocationsOperationsGetCall {
	c.urlParams_.Set("operationId", operationId)
	return c
}

// ProjectId sets the optional parameter "projectId": The Google
// Developers Console [project ID or
// project
// number](https://support.google.com/cloud/answer/6158840).
// This
//  field is deprecated, use name instead.
func (c *ProjectsLocationsOperationsGetCall) ProjectId(projectId string) *ProjectsLocationsOperationsGetCall {
	c.urlParams_.Set("projectId", projectId)
	return c
}

// Zone sets the optional parameter "zone": The name of the Google
// Compute Engine
// [zone](/compute/docs/zones#available) in which the
// cluster
// resides.
// This field is deprecated, use name instead.
func (c *ProjectsLocationsOperationsGetCall) Zone(zone string) *ProjectsLocationsOperationsGetCall {
	c.urlParams_.Set("zone", zone)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsOperationsGetCall) Fields(s ...googleapi.Field) *ProjectsLocationsOperationsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsLocationsOperationsGetCall) IfNoneMatch(entityTag string) *ProjectsLocationsOperationsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsOperationsGetCall) Context(ctx context.Context) *ProjectsLocationsOperationsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsOperationsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsOperationsGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"name": c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.operations.get" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsLocationsOperationsGetCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Gets the specified operation.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/operations/{operationsId}",
	//   "httpMethod": "GET",
	//   "id": "container.projects.locations.operations.get",
	//   "parameterOrder": [
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The name (project, location, operation id) of the operation to get.\nSpecified in the format 'projects/*/locations/*/operations/*'.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+/operations/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "operationId": {
	//       "description": "The server-assigned `name` of the operation.\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+name}",
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.locations.operations.list":

type ProjectsLocationsOperationsListCall struct {
	s            *Service
	parent       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists all operations in a project in a specific zone or all
// zones.
func (r *ProjectsLocationsOperationsService) List(parent string) *ProjectsLocationsOperationsListCall {
	c := &ProjectsLocationsOperationsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	return c
}

// ProjectId sets the optional parameter "projectId": The Google
// Developers Console [project ID or
// project
// number](https://support.google.com/cloud/answer/6158840).
// This
//  field is deprecated, use parent instead.
func (c *ProjectsLocationsOperationsListCall) ProjectId(projectId string) *ProjectsLocationsOperationsListCall {
	c.urlParams_.Set("projectId", projectId)
	return c
}

// Zone sets the optional parameter "zone": The name of the Google
// Compute Engine [zone](/compute/docs/zones#available)
// to return operations for, or `-` for all zones.
// This field is deprecated, use parent instead.
func (c *ProjectsLocationsOperationsListCall) Zone(zone string) *ProjectsLocationsOperationsListCall {
	c.urlParams_.Set("zone", zone)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsLocationsOperationsListCall) Fields(s ...googleapi.Field) *ProjectsLocationsOperationsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsLocationsOperationsListCall) IfNoneMatch(entityTag string) *ProjectsLocationsOperationsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsLocationsOperationsListCall) Context(ctx context.Context) *ProjectsLocationsOperationsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsLocationsOperationsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsLocationsOperationsListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/{+parent}/operations")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"parent": c.parent,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.locations.operations.list" call.
// Exactly one of *ListOperationsResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListOperationsResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsLocationsOperationsListCall) Do(opts ...googleapi.CallOption) (*ListOperationsResponse, error) {
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
	//   "description": "Lists all operations in a project in a specific zone or all zones.",
	//   "flatPath": "v1beta1/projects/{projectsId}/locations/{locationsId}/operations",
	//   "httpMethod": "GET",
	//   "id": "container.projects.locations.operations.list",
	//   "parameterOrder": [
	//     "parent"
	//   ],
	//   "parameters": {
	//     "parent": {
	//       "description": "The parent (project and location) where the operations will be listed.\nSpecified in the format 'projects/*/locations/*'.\nLocation \"-\" matches all zones and all regions.",
	//       "location": "path",
	//       "pattern": "^projects/[^/]+/locations/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use parent instead.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine [zone](/compute/docs/zones#available)\nto return operations for, or `-` for all zones.\nThis field is deprecated, use parent instead.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/{+parent}/operations",
	//   "response": {
	//     "$ref": "ListOperationsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.getServerconfig":

type ProjectsZonesGetServerconfigCall struct {
	s            *Service
	projectId    string
	zone         string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// GetServerconfig: Returns configuration info about the Container
// Engine service.
func (r *ProjectsZonesService) GetServerconfig(projectId string, zone string) *ProjectsZonesGetServerconfigCall {
	c := &ProjectsZonesGetServerconfigCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	return c
}

// Name sets the optional parameter "name": The name (project and
// location) of the server config to get
// Specified in the format 'projects/*/locations/*'.
func (c *ProjectsZonesGetServerconfigCall) Name(name string) *ProjectsZonesGetServerconfigCall {
	c.urlParams_.Set("name", name)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesGetServerconfigCall) Fields(s ...googleapi.Field) *ProjectsZonesGetServerconfigCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsZonesGetServerconfigCall) IfNoneMatch(entityTag string) *ProjectsZonesGetServerconfigCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesGetServerconfigCall) Context(ctx context.Context) *ProjectsZonesGetServerconfigCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesGetServerconfigCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesGetServerconfigCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/serverconfig")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"zone":      c.zone,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.getServerconfig" call.
// Exactly one of *ServerConfig or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *ServerConfig.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesGetServerconfigCall) Do(opts ...googleapi.CallOption) (*ServerConfig, error) {
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
	ret := &ServerConfig{
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
	//   "description": "Returns configuration info about the Container Engine service.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/serverconfig",
	//   "httpMethod": "GET",
	//   "id": "container.projects.zones.getServerconfig",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The name (project and location) of the server config to get\nSpecified in the format 'projects/*/locations/*'.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine [zone](/compute/docs/zones#available)\nto return operations for.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/serverconfig",
	//   "response": {
	//     "$ref": "ServerConfig"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.completeIpRotation":

type ProjectsZonesClustersCompleteIpRotationCall struct {
	s                         *Service
	projectId                 string
	zone                      string
	clusterId                 string
	completeiprotationrequest *CompleteIPRotationRequest
	urlParams_                gensupport.URLParams
	ctx_                      context.Context
	header_                   http.Header
}

// CompleteIpRotation: Completes master IP rotation.
func (r *ProjectsZonesClustersService) CompleteIpRotation(projectId string, zone string, clusterId string, completeiprotationrequest *CompleteIPRotationRequest) *ProjectsZonesClustersCompleteIpRotationCall {
	c := &ProjectsZonesClustersCompleteIpRotationCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.completeiprotationrequest = completeiprotationrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersCompleteIpRotationCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersCompleteIpRotationCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersCompleteIpRotationCall) Context(ctx context.Context) *ProjectsZonesClustersCompleteIpRotationCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersCompleteIpRotationCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersCompleteIpRotationCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.completeiprotationrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:completeIpRotation")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"zone":      c.zone,
		"clusterId": c.clusterId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.completeIpRotation" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesClustersCompleteIpRotationCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Completes master IP rotation.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:completeIpRotation",
	//   "httpMethod": "POST",
	//   "id": "container.projects.zones.clusters.completeIpRotation",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "clusterId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://developers.google.com/console/help/new/#projectnumber).\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:completeIpRotation",
	//   "request": {
	//     "$ref": "CompleteIPRotationRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.create":

type ProjectsZonesClustersCreateCall struct {
	s                    *Service
	projectId            string
	zone                 string
	createclusterrequest *CreateClusterRequest
	urlParams_           gensupport.URLParams
	ctx_                 context.Context
	header_              http.Header
}

// Create: Creates a cluster, consisting of the specified number and
// type of Google
// Compute Engine instances.
//
// By default, the cluster is created in the project's
// [default
// network](/compute/docs/networks-and-firewalls#networks).
//
// One firewall is added for the cluster. After cluster creation,
// the cluster creates routes for each node to allow the containers
// on that node to communicate with all other instances in
// the
// cluster.
//
// Finally, an entry is added to the project's global metadata
// indicating
// which CIDR range is being used by the cluster.
func (r *ProjectsZonesClustersService) Create(projectId string, zone string, createclusterrequest *CreateClusterRequest) *ProjectsZonesClustersCreateCall {
	c := &ProjectsZonesClustersCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.createclusterrequest = createclusterrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersCreateCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersCreateCall) Context(ctx context.Context) *ProjectsZonesClustersCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.createclusterrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"zone":      c.zone,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.create" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesClustersCreateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Creates a cluster, consisting of the specified number and type of Google\nCompute Engine instances.\n\nBy default, the cluster is created in the project's\n[default network](/compute/docs/networks-and-firewalls#networks).\n\nOne firewall is added for the cluster. After cluster creation,\nthe cluster creates routes for each node to allow the containers\non that node to communicate with all other instances in the\ncluster.\n\nFinally, an entry is added to the project's global metadata indicating\nwhich CIDR range is being used by the cluster.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters",
	//   "httpMethod": "POST",
	//   "id": "container.projects.zones.clusters.create",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone"
	//   ],
	//   "parameters": {
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use parent instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use parent instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters",
	//   "request": {
	//     "$ref": "CreateClusterRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.delete":

type ProjectsZonesClustersDeleteCall struct {
	s          *Service
	projectId  string
	zone       string
	clusterId  string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Delete: Deletes the cluster, including the Kubernetes endpoint and
// all worker
// nodes.
//
// Firewalls and routes that were configured during cluster creation
// are also deleted.
//
// Other Google Compute Engine resources that might be in use by the
// cluster
// (e.g. load balancer resources) will not be deleted if they weren't
// present
// at the initial create time.
func (r *ProjectsZonesClustersService) Delete(projectId string, zone string, clusterId string) *ProjectsZonesClustersDeleteCall {
	c := &ProjectsZonesClustersDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	return c
}

// Name sets the optional parameter "name": The name (project, location,
// cluster) of the cluster to delete.
// Specified in the format 'projects/*/locations/*/clusters/*'.
func (c *ProjectsZonesClustersDeleteCall) Name(name string) *ProjectsZonesClustersDeleteCall {
	c.urlParams_.Set("name", name)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersDeleteCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersDeleteCall) Context(ctx context.Context) *ProjectsZonesClustersDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"zone":      c.zone,
		"clusterId": c.clusterId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.delete" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesClustersDeleteCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Deletes the cluster, including the Kubernetes endpoint and all worker\nnodes.\n\nFirewalls and routes that were configured during cluster creation\nare also deleted.\n\nOther Google Compute Engine resources that might be in use by the cluster\n(e.g. load balancer resources) will not be deleted if they weren't present\nat the initial create time.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}",
	//   "httpMethod": "DELETE",
	//   "id": "container.projects.zones.clusters.delete",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "clusterId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster to delete.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "name": {
	//       "description": "The name (project, location, cluster) of the cluster to delete.\nSpecified in the format 'projects/*/locations/*/clusters/*'.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}",
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.get":

type ProjectsZonesClustersGetCall struct {
	s            *Service
	projectId    string
	zone         string
	clusterId    string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Gets the details of a specific cluster.
func (r *ProjectsZonesClustersService) Get(projectId string, zone string, clusterId string) *ProjectsZonesClustersGetCall {
	c := &ProjectsZonesClustersGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	return c
}

// Name sets the optional parameter "name": The name (project, location,
// cluster) of the cluster to retrieve.
// Specified in the format 'projects/*/locations/*/clusters/*'.
func (c *ProjectsZonesClustersGetCall) Name(name string) *ProjectsZonesClustersGetCall {
	c.urlParams_.Set("name", name)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersGetCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsZonesClustersGetCall) IfNoneMatch(entityTag string) *ProjectsZonesClustersGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersGetCall) Context(ctx context.Context) *ProjectsZonesClustersGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"zone":      c.zone,
		"clusterId": c.clusterId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.get" call.
// Exactly one of *Cluster or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Cluster.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *ProjectsZonesClustersGetCall) Do(opts ...googleapi.CallOption) (*Cluster, error) {
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
	ret := &Cluster{
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
	//   "description": "Gets the details of a specific cluster.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}",
	//   "httpMethod": "GET",
	//   "id": "container.projects.zones.clusters.get",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "clusterId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster to retrieve.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "name": {
	//       "description": "The name (project, location, cluster) of the cluster to retrieve.\nSpecified in the format 'projects/*/locations/*/clusters/*'.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}",
	//   "response": {
	//     "$ref": "Cluster"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.legacyAbac":

type ProjectsZonesClustersLegacyAbacCall struct {
	s                    *Service
	projectId            string
	zone                 string
	clusterId            string
	setlegacyabacrequest *SetLegacyAbacRequest
	urlParams_           gensupport.URLParams
	ctx_                 context.Context
	header_              http.Header
}

// LegacyAbac: Enables or disables the ABAC authorization mechanism on a
// cluster.
func (r *ProjectsZonesClustersService) LegacyAbac(projectId string, zone string, clusterId string, setlegacyabacrequest *SetLegacyAbacRequest) *ProjectsZonesClustersLegacyAbacCall {
	c := &ProjectsZonesClustersLegacyAbacCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.setlegacyabacrequest = setlegacyabacrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersLegacyAbacCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersLegacyAbacCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersLegacyAbacCall) Context(ctx context.Context) *ProjectsZonesClustersLegacyAbacCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersLegacyAbacCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersLegacyAbacCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setlegacyabacrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/legacyAbac")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"zone":      c.zone,
		"clusterId": c.clusterId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.legacyAbac" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesClustersLegacyAbacCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Enables or disables the ABAC authorization mechanism on a cluster.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/legacyAbac",
	//   "httpMethod": "POST",
	//   "id": "container.projects.zones.clusters.legacyAbac",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "clusterId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster to update.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/legacyAbac",
	//   "request": {
	//     "$ref": "SetLegacyAbacRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.list":

type ProjectsZonesClustersListCall struct {
	s            *Service
	projectId    string
	zone         string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists all clusters owned by a project in either the specified
// zone or all
// zones.
func (r *ProjectsZonesClustersService) List(projectId string, zone string) *ProjectsZonesClustersListCall {
	c := &ProjectsZonesClustersListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	return c
}

// Parent sets the optional parameter "parent": The parent (project and
// location) where the clusters will be listed.
// Specified in the format 'projects/*/locations/*'.
// Location "-" matches all zones and all regions.
func (c *ProjectsZonesClustersListCall) Parent(parent string) *ProjectsZonesClustersListCall {
	c.urlParams_.Set("parent", parent)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersListCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsZonesClustersListCall) IfNoneMatch(entityTag string) *ProjectsZonesClustersListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersListCall) Context(ctx context.Context) *ProjectsZonesClustersListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"zone":      c.zone,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.list" call.
// Exactly one of *ListClustersResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListClustersResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsZonesClustersListCall) Do(opts ...googleapi.CallOption) (*ListClustersResponse, error) {
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
	ret := &ListClustersResponse{
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
	//   "description": "Lists all clusters owned by a project in either the specified zone or all\nzones.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters",
	//   "httpMethod": "GET",
	//   "id": "container.projects.zones.clusters.list",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone"
	//   ],
	//   "parameters": {
	//     "parent": {
	//       "description": "The parent (project and location) where the clusters will be listed.\nSpecified in the format 'projects/*/locations/*'.\nLocation \"-\" matches all zones and all regions.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use parent instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides, or \"-\" for all zones.\nThis field is deprecated, use parent instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters",
	//   "response": {
	//     "$ref": "ListClustersResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.resourceLabels":

type ProjectsZonesClustersResourceLabelsCall struct {
	s                *Service
	projectId        string
	zone             string
	clusterId        string
	setlabelsrequest *SetLabelsRequest
	urlParams_       gensupport.URLParams
	ctx_             context.Context
	header_          http.Header
}

// ResourceLabels: Sets labels on a cluster.
func (r *ProjectsZonesClustersService) ResourceLabels(projectId string, zone string, clusterId string, setlabelsrequest *SetLabelsRequest) *ProjectsZonesClustersResourceLabelsCall {
	c := &ProjectsZonesClustersResourceLabelsCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.setlabelsrequest = setlabelsrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersResourceLabelsCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersResourceLabelsCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersResourceLabelsCall) Context(ctx context.Context) *ProjectsZonesClustersResourceLabelsCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersResourceLabelsCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersResourceLabelsCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setlabelsrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/resourceLabels")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"zone":      c.zone,
		"clusterId": c.clusterId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.resourceLabels" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesClustersResourceLabelsCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Sets labels on a cluster.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/resourceLabels",
	//   "httpMethod": "POST",
	//   "id": "container.projects.zones.clusters.resourceLabels",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "clusterId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://developers.google.com/console/help/new/#projectnumber).\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/resourceLabels",
	//   "request": {
	//     "$ref": "SetLabelsRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.setMaintenancePolicy":

type ProjectsZonesClustersSetMaintenancePolicyCall struct {
	s                           *Service
	projectId                   string
	zone                        string
	clusterId                   string
	setmaintenancepolicyrequest *SetMaintenancePolicyRequest
	urlParams_                  gensupport.URLParams
	ctx_                        context.Context
	header_                     http.Header
}

// SetMaintenancePolicy: Sets the maintenance policy for a cluster.
func (r *ProjectsZonesClustersService) SetMaintenancePolicy(projectId string, zone string, clusterId string, setmaintenancepolicyrequest *SetMaintenancePolicyRequest) *ProjectsZonesClustersSetMaintenancePolicyCall {
	c := &ProjectsZonesClustersSetMaintenancePolicyCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.setmaintenancepolicyrequest = setmaintenancepolicyrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersSetMaintenancePolicyCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersSetMaintenancePolicyCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersSetMaintenancePolicyCall) Context(ctx context.Context) *ProjectsZonesClustersSetMaintenancePolicyCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersSetMaintenancePolicyCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersSetMaintenancePolicyCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setmaintenancepolicyrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:setMaintenancePolicy")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"zone":      c.zone,
		"clusterId": c.clusterId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.setMaintenancePolicy" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesClustersSetMaintenancePolicyCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Sets the maintenance policy for a cluster.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:setMaintenancePolicy",
	//   "httpMethod": "POST",
	//   "id": "container.projects.zones.clusters.setMaintenancePolicy",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "clusterId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster to update.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:setMaintenancePolicy",
	//   "request": {
	//     "$ref": "SetMaintenancePolicyRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.setMasterAuth":

type ProjectsZonesClustersSetMasterAuthCall struct {
	s                    *Service
	projectId            string
	zone                 string
	clusterId            string
	setmasterauthrequest *SetMasterAuthRequest
	urlParams_           gensupport.URLParams
	ctx_                 context.Context
	header_              http.Header
}

// SetMasterAuth: Used to set master auth materials. Currently supports
// :-
// Changing the admin password of a specific cluster.
// This can be either via password generation or explicitly set.
// Modify basic_auth.csv and reset the K8S API server.
func (r *ProjectsZonesClustersService) SetMasterAuth(projectId string, zone string, clusterId string, setmasterauthrequest *SetMasterAuthRequest) *ProjectsZonesClustersSetMasterAuthCall {
	c := &ProjectsZonesClustersSetMasterAuthCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.setmasterauthrequest = setmasterauthrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersSetMasterAuthCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersSetMasterAuthCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersSetMasterAuthCall) Context(ctx context.Context) *ProjectsZonesClustersSetMasterAuthCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersSetMasterAuthCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersSetMasterAuthCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setmasterauthrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:setMasterAuth")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"zone":      c.zone,
		"clusterId": c.clusterId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.setMasterAuth" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesClustersSetMasterAuthCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Used to set master auth materials. Currently supports :-\nChanging the admin password of a specific cluster.\nThis can be either via password generation or explicitly set.\nModify basic_auth.csv and reset the K8S API server.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:setMasterAuth",
	//   "httpMethod": "POST",
	//   "id": "container.projects.zones.clusters.setMasterAuth",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "clusterId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster to upgrade.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:setMasterAuth",
	//   "request": {
	//     "$ref": "SetMasterAuthRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.setNetworkPolicy":

type ProjectsZonesClustersSetNetworkPolicyCall struct {
	s                       *Service
	projectId               string
	zone                    string
	clusterId               string
	setnetworkpolicyrequest *SetNetworkPolicyRequest
	urlParams_              gensupport.URLParams
	ctx_                    context.Context
	header_                 http.Header
}

// SetNetworkPolicy: Enables/Disables Network Policy for a cluster.
func (r *ProjectsZonesClustersService) SetNetworkPolicy(projectId string, zone string, clusterId string, setnetworkpolicyrequest *SetNetworkPolicyRequest) *ProjectsZonesClustersSetNetworkPolicyCall {
	c := &ProjectsZonesClustersSetNetworkPolicyCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.setnetworkpolicyrequest = setnetworkpolicyrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersSetNetworkPolicyCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersSetNetworkPolicyCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersSetNetworkPolicyCall) Context(ctx context.Context) *ProjectsZonesClustersSetNetworkPolicyCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersSetNetworkPolicyCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersSetNetworkPolicyCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setnetworkpolicyrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:setNetworkPolicy")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"zone":      c.zone,
		"clusterId": c.clusterId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.setNetworkPolicy" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesClustersSetNetworkPolicyCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Enables/Disables Network Policy for a cluster.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:setNetworkPolicy",
	//   "httpMethod": "POST",
	//   "id": "container.projects.zones.clusters.setNetworkPolicy",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "clusterId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://developers.google.com/console/help/new/#projectnumber).\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:setNetworkPolicy",
	//   "request": {
	//     "$ref": "SetNetworkPolicyRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.startIpRotation":

type ProjectsZonesClustersStartIpRotationCall struct {
	s                      *Service
	projectId              string
	zone                   string
	clusterId              string
	startiprotationrequest *StartIPRotationRequest
	urlParams_             gensupport.URLParams
	ctx_                   context.Context
	header_                http.Header
}

// StartIpRotation: Start master IP rotation.
func (r *ProjectsZonesClustersService) StartIpRotation(projectId string, zone string, clusterId string, startiprotationrequest *StartIPRotationRequest) *ProjectsZonesClustersStartIpRotationCall {
	c := &ProjectsZonesClustersStartIpRotationCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.startiprotationrequest = startiprotationrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersStartIpRotationCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersStartIpRotationCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersStartIpRotationCall) Context(ctx context.Context) *ProjectsZonesClustersStartIpRotationCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersStartIpRotationCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersStartIpRotationCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.startiprotationrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:startIpRotation")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"zone":      c.zone,
		"clusterId": c.clusterId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.startIpRotation" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesClustersStartIpRotationCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Start master IP rotation.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:startIpRotation",
	//   "httpMethod": "POST",
	//   "id": "container.projects.zones.clusters.startIpRotation",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "clusterId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://developers.google.com/console/help/new/#projectnumber).\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}:startIpRotation",
	//   "request": {
	//     "$ref": "StartIPRotationRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.update":

type ProjectsZonesClustersUpdateCall struct {
	s                    *Service
	projectId            string
	zone                 string
	clusterId            string
	updateclusterrequest *UpdateClusterRequest
	urlParams_           gensupport.URLParams
	ctx_                 context.Context
	header_              http.Header
}

// Update: Updates the settings of a specific cluster.
func (r *ProjectsZonesClustersService) Update(projectId string, zone string, clusterId string, updateclusterrequest *UpdateClusterRequest) *ProjectsZonesClustersUpdateCall {
	c := &ProjectsZonesClustersUpdateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.updateclusterrequest = updateclusterrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersUpdateCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersUpdateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersUpdateCall) Context(ctx context.Context) *ProjectsZonesClustersUpdateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersUpdateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersUpdateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.updateclusterrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PUT", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"zone":      c.zone,
		"clusterId": c.clusterId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.update" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesClustersUpdateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Updates the settings of a specific cluster.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}",
	//   "httpMethod": "PUT",
	//   "id": "container.projects.zones.clusters.update",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "clusterId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster to upgrade.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}",
	//   "request": {
	//     "$ref": "UpdateClusterRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.nodePools.create":

type ProjectsZonesClustersNodePoolsCreateCall struct {
	s                     *Service
	projectId             string
	zone                  string
	clusterId             string
	createnodepoolrequest *CreateNodePoolRequest
	urlParams_            gensupport.URLParams
	ctx_                  context.Context
	header_               http.Header
}

// Create: Creates a node pool for a cluster.
func (r *ProjectsZonesClustersNodePoolsService) Create(projectId string, zone string, clusterId string, createnodepoolrequest *CreateNodePoolRequest) *ProjectsZonesClustersNodePoolsCreateCall {
	c := &ProjectsZonesClustersNodePoolsCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.createnodepoolrequest = createnodepoolrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersNodePoolsCreateCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersNodePoolsCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersNodePoolsCreateCall) Context(ctx context.Context) *ProjectsZonesClustersNodePoolsCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersNodePoolsCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersNodePoolsCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.createnodepoolrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"zone":      c.zone,
		"clusterId": c.clusterId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.nodePools.create" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesClustersNodePoolsCreateCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Creates a node pool for a cluster.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools",
	//   "httpMethod": "POST",
	//   "id": "container.projects.zones.clusters.nodePools.create",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "clusterId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster.\nThis field is deprecated, use parent instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://developers.google.com/console/help/new/#projectnumber).\nThis field is deprecated, use parent instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use parent instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools",
	//   "request": {
	//     "$ref": "CreateNodePoolRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.nodePools.delete":

type ProjectsZonesClustersNodePoolsDeleteCall struct {
	s          *Service
	projectId  string
	zone       string
	clusterId  string
	nodePoolId string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Delete: Deletes a node pool from a cluster.
func (r *ProjectsZonesClustersNodePoolsService) Delete(projectId string, zone string, clusterId string, nodePoolId string) *ProjectsZonesClustersNodePoolsDeleteCall {
	c := &ProjectsZonesClustersNodePoolsDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.nodePoolId = nodePoolId
	return c
}

// Name sets the optional parameter "name": The name (project, location,
// cluster, node pool id) of the node pool to delete.
// Specified in the format
// 'projects/*/locations/*/clusters/*/nodePools/*'.
func (c *ProjectsZonesClustersNodePoolsDeleteCall) Name(name string) *ProjectsZonesClustersNodePoolsDeleteCall {
	c.urlParams_.Set("name", name)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersNodePoolsDeleteCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersNodePoolsDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersNodePoolsDeleteCall) Context(ctx context.Context) *ProjectsZonesClustersNodePoolsDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersNodePoolsDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersNodePoolsDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":  c.projectId,
		"zone":       c.zone,
		"clusterId":  c.clusterId,
		"nodePoolId": c.nodePoolId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.nodePools.delete" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesClustersNodePoolsDeleteCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Deletes a node pool from a cluster.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}",
	//   "httpMethod": "DELETE",
	//   "id": "container.projects.zones.clusters.nodePools.delete",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "clusterId",
	//     "nodePoolId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "name": {
	//       "description": "The name (project, location, cluster, node pool id) of the node pool to delete.\nSpecified in the format 'projects/*/locations/*/clusters/*/nodePools/*'.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "nodePoolId": {
	//       "description": "The name of the node pool to delete.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://developers.google.com/console/help/new/#projectnumber).\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}",
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.nodePools.get":

type ProjectsZonesClustersNodePoolsGetCall struct {
	s            *Service
	projectId    string
	zone         string
	clusterId    string
	nodePoolId   string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Retrieves the node pool requested.
func (r *ProjectsZonesClustersNodePoolsService) Get(projectId string, zone string, clusterId string, nodePoolId string) *ProjectsZonesClustersNodePoolsGetCall {
	c := &ProjectsZonesClustersNodePoolsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.nodePoolId = nodePoolId
	return c
}

// Name sets the optional parameter "name": The name (project, location,
// cluster, node pool id) of the node pool to get.
// Specified in the format
// 'projects/*/locations/*/clusters/*/nodePools/*'.
func (c *ProjectsZonesClustersNodePoolsGetCall) Name(name string) *ProjectsZonesClustersNodePoolsGetCall {
	c.urlParams_.Set("name", name)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersNodePoolsGetCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersNodePoolsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsZonesClustersNodePoolsGetCall) IfNoneMatch(entityTag string) *ProjectsZonesClustersNodePoolsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersNodePoolsGetCall) Context(ctx context.Context) *ProjectsZonesClustersNodePoolsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersNodePoolsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersNodePoolsGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":  c.projectId,
		"zone":       c.zone,
		"clusterId":  c.clusterId,
		"nodePoolId": c.nodePoolId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.nodePools.get" call.
// Exactly one of *NodePool or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *NodePool.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesClustersNodePoolsGetCall) Do(opts ...googleapi.CallOption) (*NodePool, error) {
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
	ret := &NodePool{
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
	//   "description": "Retrieves the node pool requested.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}",
	//   "httpMethod": "GET",
	//   "id": "container.projects.zones.clusters.nodePools.get",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "clusterId",
	//     "nodePoolId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "name": {
	//       "description": "The name (project, location, cluster, node pool id) of the node pool to get.\nSpecified in the format 'projects/*/locations/*/clusters/*/nodePools/*'.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "nodePoolId": {
	//       "description": "The name of the node pool.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://developers.google.com/console/help/new/#projectnumber).\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}",
	//   "response": {
	//     "$ref": "NodePool"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.nodePools.list":

type ProjectsZonesClustersNodePoolsListCall struct {
	s            *Service
	projectId    string
	zone         string
	clusterId    string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists the node pools for a cluster.
func (r *ProjectsZonesClustersNodePoolsService) List(projectId string, zone string, clusterId string) *ProjectsZonesClustersNodePoolsListCall {
	c := &ProjectsZonesClustersNodePoolsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	return c
}

// Parent sets the optional parameter "parent": The parent (project,
// location, cluster id) where the node pools will be listed.
// Specified in the format 'projects/*/locations/*/clusters/*'.
func (c *ProjectsZonesClustersNodePoolsListCall) Parent(parent string) *ProjectsZonesClustersNodePoolsListCall {
	c.urlParams_.Set("parent", parent)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersNodePoolsListCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersNodePoolsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsZonesClustersNodePoolsListCall) IfNoneMatch(entityTag string) *ProjectsZonesClustersNodePoolsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersNodePoolsListCall) Context(ctx context.Context) *ProjectsZonesClustersNodePoolsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersNodePoolsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersNodePoolsListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"zone":      c.zone,
		"clusterId": c.clusterId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.nodePools.list" call.
// Exactly one of *ListNodePoolsResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListNodePoolsResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsZonesClustersNodePoolsListCall) Do(opts ...googleapi.CallOption) (*ListNodePoolsResponse, error) {
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
	ret := &ListNodePoolsResponse{
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
	//   "description": "Lists the node pools for a cluster.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools",
	//   "httpMethod": "GET",
	//   "id": "container.projects.zones.clusters.nodePools.list",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "clusterId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster.\nThis field is deprecated, use parent instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "parent": {
	//       "description": "The parent (project, location, cluster id) where the node pools will be listed.\nSpecified in the format 'projects/*/locations/*/clusters/*'.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://developers.google.com/console/help/new/#projectnumber).\nThis field is deprecated, use parent instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use parent instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools",
	//   "response": {
	//     "$ref": "ListNodePoolsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.nodePools.rollback":

type ProjectsZonesClustersNodePoolsRollbackCall struct {
	s                              *Service
	projectId                      string
	zone                           string
	clusterId                      string
	nodePoolId                     string
	rollbacknodepoolupgraderequest *RollbackNodePoolUpgradeRequest
	urlParams_                     gensupport.URLParams
	ctx_                           context.Context
	header_                        http.Header
}

// Rollback: Roll back the previously Aborted or Failed NodePool
// upgrade.
// This will be an no-op if the last upgrade successfully completed.
func (r *ProjectsZonesClustersNodePoolsService) Rollback(projectId string, zone string, clusterId string, nodePoolId string, rollbacknodepoolupgraderequest *RollbackNodePoolUpgradeRequest) *ProjectsZonesClustersNodePoolsRollbackCall {
	c := &ProjectsZonesClustersNodePoolsRollbackCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.nodePoolId = nodePoolId
	c.rollbacknodepoolupgraderequest = rollbacknodepoolupgraderequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersNodePoolsRollbackCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersNodePoolsRollbackCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersNodePoolsRollbackCall) Context(ctx context.Context) *ProjectsZonesClustersNodePoolsRollbackCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersNodePoolsRollbackCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersNodePoolsRollbackCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.rollbacknodepoolupgraderequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}:rollback")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":  c.projectId,
		"zone":       c.zone,
		"clusterId":  c.clusterId,
		"nodePoolId": c.nodePoolId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.nodePools.rollback" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesClustersNodePoolsRollbackCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Roll back the previously Aborted or Failed NodePool upgrade.\nThis will be an no-op if the last upgrade successfully completed.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}:rollback",
	//   "httpMethod": "POST",
	//   "id": "container.projects.zones.clusters.nodePools.rollback",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "clusterId",
	//     "nodePoolId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster to rollback.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "nodePoolId": {
	//       "description": "The name of the node pool to rollback.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}:rollback",
	//   "request": {
	//     "$ref": "RollbackNodePoolUpgradeRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.clusters.nodePools.setManagement":

type ProjectsZonesClustersNodePoolsSetManagementCall struct {
	s                            *Service
	projectId                    string
	zone                         string
	clusterId                    string
	nodePoolId                   string
	setnodepoolmanagementrequest *SetNodePoolManagementRequest
	urlParams_                   gensupport.URLParams
	ctx_                         context.Context
	header_                      http.Header
}

// SetManagement: Sets the NodeManagement options for a node pool.
func (r *ProjectsZonesClustersNodePoolsService) SetManagement(projectId string, zone string, clusterId string, nodePoolId string, setnodepoolmanagementrequest *SetNodePoolManagementRequest) *ProjectsZonesClustersNodePoolsSetManagementCall {
	c := &ProjectsZonesClustersNodePoolsSetManagementCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.clusterId = clusterId
	c.nodePoolId = nodePoolId
	c.setnodepoolmanagementrequest = setnodepoolmanagementrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesClustersNodePoolsSetManagementCall) Fields(s ...googleapi.Field) *ProjectsZonesClustersNodePoolsSetManagementCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesClustersNodePoolsSetManagementCall) Context(ctx context.Context) *ProjectsZonesClustersNodePoolsSetManagementCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesClustersNodePoolsSetManagementCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesClustersNodePoolsSetManagementCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.setnodepoolmanagementrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}/setManagement")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":  c.projectId,
		"zone":       c.zone,
		"clusterId":  c.clusterId,
		"nodePoolId": c.nodePoolId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.clusters.nodePools.setManagement" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesClustersNodePoolsSetManagementCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Sets the NodeManagement options for a node pool.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}/setManagement",
	//   "httpMethod": "POST",
	//   "id": "container.projects.zones.clusters.nodePools.setManagement",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "clusterId",
	//     "nodePoolId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "The name of the cluster to update.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "nodePoolId": {
	//       "description": "The name of the node pool to update.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/clusters/{clusterId}/nodePools/{nodePoolId}/setManagement",
	//   "request": {
	//     "$ref": "SetNodePoolManagementRequest"
	//   },
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.operations.cancel":

type ProjectsZonesOperationsCancelCall struct {
	s                      *Service
	projectId              string
	zone                   string
	operationId            string
	canceloperationrequest *CancelOperationRequest
	urlParams_             gensupport.URLParams
	ctx_                   context.Context
	header_                http.Header
}

// Cancel: Cancels the specified operation.
func (r *ProjectsZonesOperationsService) Cancel(projectId string, zone string, operationId string, canceloperationrequest *CancelOperationRequest) *ProjectsZonesOperationsCancelCall {
	c := &ProjectsZonesOperationsCancelCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.operationId = operationId
	c.canceloperationrequest = canceloperationrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesOperationsCancelCall) Fields(s ...googleapi.Field) *ProjectsZonesOperationsCancelCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesOperationsCancelCall) Context(ctx context.Context) *ProjectsZonesOperationsCancelCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesOperationsCancelCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesOperationsCancelCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.canceloperationrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/operations/{operationId}:cancel")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":   c.projectId,
		"zone":        c.zone,
		"operationId": c.operationId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.operations.cancel" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *ProjectsZonesOperationsCancelCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
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
	//   "description": "Cancels the specified operation.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/operations/{operationId}:cancel",
	//   "httpMethod": "POST",
	//   "id": "container.projects.zones.operations.cancel",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "operationId"
	//   ],
	//   "parameters": {
	//     "operationId": {
	//       "description": "The server-assigned `name` of the operation.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the operation resides.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/operations/{operationId}:cancel",
	//   "request": {
	//     "$ref": "CancelOperationRequest"
	//   },
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.operations.get":

type ProjectsZonesOperationsGetCall struct {
	s            *Service
	projectId    string
	zone         string
	operationId  string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Gets the specified operation.
func (r *ProjectsZonesOperationsService) Get(projectId string, zone string, operationId string) *ProjectsZonesOperationsGetCall {
	c := &ProjectsZonesOperationsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	c.operationId = operationId
	return c
}

// Name sets the optional parameter "name": The name (project, location,
// operation id) of the operation to get.
// Specified in the format 'projects/*/locations/*/operations/*'.
func (c *ProjectsZonesOperationsGetCall) Name(name string) *ProjectsZonesOperationsGetCall {
	c.urlParams_.Set("name", name)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesOperationsGetCall) Fields(s ...googleapi.Field) *ProjectsZonesOperationsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsZonesOperationsGetCall) IfNoneMatch(entityTag string) *ProjectsZonesOperationsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesOperationsGetCall) Context(ctx context.Context) *ProjectsZonesOperationsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesOperationsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesOperationsGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/operations/{operationId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":   c.projectId,
		"zone":        c.zone,
		"operationId": c.operationId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.operations.get" call.
// Exactly one of *Operation or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Operation.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsZonesOperationsGetCall) Do(opts ...googleapi.CallOption) (*Operation, error) {
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
	//   "description": "Gets the specified operation.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/operations/{operationId}",
	//   "httpMethod": "GET",
	//   "id": "container.projects.zones.operations.get",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone",
	//     "operationId"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "The name (project, location, operation id) of the operation to get.\nSpecified in the format 'projects/*/locations/*/operations/*'.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "operationId": {
	//       "description": "The server-assigned `name` of the operation.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine\n[zone](/compute/docs/zones#available) in which the cluster\nresides.\nThis field is deprecated, use name instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/operations/{operationId}",
	//   "response": {
	//     "$ref": "Operation"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "container.projects.zones.operations.list":

type ProjectsZonesOperationsListCall struct {
	s            *Service
	projectId    string
	zone         string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists all operations in a project in a specific zone or all
// zones.
func (r *ProjectsZonesOperationsService) List(projectId string, zone string) *ProjectsZonesOperationsListCall {
	c := &ProjectsZonesOperationsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.zone = zone
	return c
}

// Parent sets the optional parameter "parent": The parent (project and
// location) where the operations will be listed.
// Specified in the format 'projects/*/locations/*'.
// Location "-" matches all zones and all regions.
func (c *ProjectsZonesOperationsListCall) Parent(parent string) *ProjectsZonesOperationsListCall {
	c.urlParams_.Set("parent", parent)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsZonesOperationsListCall) Fields(s ...googleapi.Field) *ProjectsZonesOperationsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsZonesOperationsListCall) IfNoneMatch(entityTag string) *ProjectsZonesOperationsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsZonesOperationsListCall) Context(ctx context.Context) *ProjectsZonesOperationsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsZonesOperationsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsZonesOperationsListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1beta1/projects/{projectId}/zones/{zone}/operations")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"zone":      c.zone,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "container.projects.zones.operations.list" call.
// Exactly one of *ListOperationsResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListOperationsResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsZonesOperationsListCall) Do(opts ...googleapi.CallOption) (*ListOperationsResponse, error) {
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
	//   "description": "Lists all operations in a project in a specific zone or all zones.",
	//   "flatPath": "v1beta1/projects/{projectId}/zones/{zone}/operations",
	//   "httpMethod": "GET",
	//   "id": "container.projects.zones.operations.list",
	//   "parameterOrder": [
	//     "projectId",
	//     "zone"
	//   ],
	//   "parameters": {
	//     "parent": {
	//       "description": "The parent (project and location) where the operations will be listed.\nSpecified in the format 'projects/*/locations/*'.\nLocation \"-\" matches all zones and all regions.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The Google Developers Console [project ID or project\nnumber](https://support.google.com/cloud/answer/6158840).\nThis field is deprecated, use parent instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "zone": {
	//       "description": "The name of the Google Compute Engine [zone](/compute/docs/zones#available)\nto return operations for, or `-` for all zones.\nThis field is deprecated, use parent instead.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1beta1/projects/{projectId}/zones/{zone}/operations",
	//   "response": {
	//     "$ref": "ListOperationsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}
