// Package toolresults provides access to the Cloud Tool Results API.
//
// See https://firebase.google.com/docs/test-lab/
//
// Usage example:
//
//   import "google.golang.org/api/toolresults/v1beta3"
//   ...
//   toolresultsService, err := toolresults.New(oauthHttpClient)
package toolresults // import "google.golang.org/api/toolresults/v1beta3"

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

const apiId = "toolresults:v1beta3"
const apiName = "toolresults"
const apiVersion = "v1beta3"
const basePath = "https://www.googleapis.com/toolresults/v1beta3/projects/"

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
	rs.Histories = NewProjectsHistoriesService(s)
	return rs
}

type ProjectsService struct {
	s *Service

	Histories *ProjectsHistoriesService
}

func NewProjectsHistoriesService(s *Service) *ProjectsHistoriesService {
	rs := &ProjectsHistoriesService{s: s}
	rs.Executions = NewProjectsHistoriesExecutionsService(s)
	return rs
}

type ProjectsHistoriesService struct {
	s *Service

	Executions *ProjectsHistoriesExecutionsService
}

func NewProjectsHistoriesExecutionsService(s *Service) *ProjectsHistoriesExecutionsService {
	rs := &ProjectsHistoriesExecutionsService{s: s}
	rs.Clusters = NewProjectsHistoriesExecutionsClustersService(s)
	rs.Steps = NewProjectsHistoriesExecutionsStepsService(s)
	return rs
}

type ProjectsHistoriesExecutionsService struct {
	s *Service

	Clusters *ProjectsHistoriesExecutionsClustersService

	Steps *ProjectsHistoriesExecutionsStepsService
}

func NewProjectsHistoriesExecutionsClustersService(s *Service) *ProjectsHistoriesExecutionsClustersService {
	rs := &ProjectsHistoriesExecutionsClustersService{s: s}
	return rs
}

type ProjectsHistoriesExecutionsClustersService struct {
	s *Service
}

func NewProjectsHistoriesExecutionsStepsService(s *Service) *ProjectsHistoriesExecutionsStepsService {
	rs := &ProjectsHistoriesExecutionsStepsService{s: s}
	rs.PerfMetricsSummary = NewProjectsHistoriesExecutionsStepsPerfMetricsSummaryService(s)
	rs.PerfSampleSeries = NewProjectsHistoriesExecutionsStepsPerfSampleSeriesService(s)
	rs.Thumbnails = NewProjectsHistoriesExecutionsStepsThumbnailsService(s)
	return rs
}

type ProjectsHistoriesExecutionsStepsService struct {
	s *Service

	PerfMetricsSummary *ProjectsHistoriesExecutionsStepsPerfMetricsSummaryService

	PerfSampleSeries *ProjectsHistoriesExecutionsStepsPerfSampleSeriesService

	Thumbnails *ProjectsHistoriesExecutionsStepsThumbnailsService
}

func NewProjectsHistoriesExecutionsStepsPerfMetricsSummaryService(s *Service) *ProjectsHistoriesExecutionsStepsPerfMetricsSummaryService {
	rs := &ProjectsHistoriesExecutionsStepsPerfMetricsSummaryService{s: s}
	return rs
}

type ProjectsHistoriesExecutionsStepsPerfMetricsSummaryService struct {
	s *Service
}

func NewProjectsHistoriesExecutionsStepsPerfSampleSeriesService(s *Service) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesService {
	rs := &ProjectsHistoriesExecutionsStepsPerfSampleSeriesService{s: s}
	rs.Samples = NewProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesService(s)
	return rs
}

type ProjectsHistoriesExecutionsStepsPerfSampleSeriesService struct {
	s *Service

	Samples *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesService
}

func NewProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesService(s *Service) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesService {
	rs := &ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesService{s: s}
	return rs
}

type ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesService struct {
	s *Service
}

func NewProjectsHistoriesExecutionsStepsThumbnailsService(s *Service) *ProjectsHistoriesExecutionsStepsThumbnailsService {
	rs := &ProjectsHistoriesExecutionsStepsThumbnailsService{s: s}
	return rs
}

type ProjectsHistoriesExecutionsStepsThumbnailsService struct {
	s *Service
}

// AndroidAppInfo: Android app information.
type AndroidAppInfo struct {
	// Name: The name of the app. Optional
	Name string `json:"name,omitempty"`

	// PackageName: The package name of the app. Required.
	PackageName string `json:"packageName,omitempty"`

	// VersionCode: The internal version code of the app. Optional.
	VersionCode string `json:"versionCode,omitempty"`

	// VersionName: The version name of the app. Optional.
	VersionName string `json:"versionName,omitempty"`

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

func (s *AndroidAppInfo) MarshalJSON() ([]byte, error) {
	type noMethod AndroidAppInfo
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// AndroidInstrumentationTest: A test of an Android application that can
// control an Android component independently of its normal
// lifecycle.
//
// See  for more information on types of Android tests.
type AndroidInstrumentationTest struct {
	// TestPackageId: The java package for the test to be executed. Required
	TestPackageId string `json:"testPackageId,omitempty"`

	// TestRunnerClass: The InstrumentationTestRunner class. Required
	TestRunnerClass string `json:"testRunnerClass,omitempty"`

	// TestTargets: Each target must be fully qualified with the package
	// name or class name, in one of these formats: - "package package_name"
	// - "class package_name.class_name" - "class
	// package_name.class_name#method_name"
	//
	// If empty, all targets in the module will be run.
	TestTargets []string `json:"testTargets,omitempty"`

	// UseOrchestrator: The flag indicates whether Android Test Orchestrator
	// will be used to run test or not. Test orchestrator is used if either:
	// - orchestrator_option field is USE_ORCHESTRATOR, and test runner is
	// compatible with orchestrator. Or - orchestrator_option field is
	// unspecified or ORCHESTRATOR_OPTION_UNSPECIFIED, and test runner is
	// compatible with orchestrator.
	UseOrchestrator bool `json:"useOrchestrator,omitempty"`

	// ForceSendFields is a list of field names (e.g. "TestPackageId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "TestPackageId") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *AndroidInstrumentationTest) MarshalJSON() ([]byte, error) {
	type noMethod AndroidInstrumentationTest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// AndroidRoboTest: A test of an android application that explores the
// application on a virtual or physical Android device, finding culprits
// and crashes as it goes.
type AndroidRoboTest struct {
	// AppInitialActivity: The initial activity that should be used to start
	// the app. Optional
	AppInitialActivity string `json:"appInitialActivity,omitempty"`

	// BootstrapPackageId: The java package for the bootstrap. Optional
	BootstrapPackageId string `json:"bootstrapPackageId,omitempty"`

	// BootstrapRunnerClass: The runner class for the bootstrap. Optional
	BootstrapRunnerClass string `json:"bootstrapRunnerClass,omitempty"`

	// MaxDepth: The max depth of the traversal stack Robo can explore.
	// Optional
	MaxDepth int64 `json:"maxDepth,omitempty"`

	// MaxSteps: The max number of steps/actions Robo can execute. Default
	// is no limit (0). Optional
	MaxSteps int64 `json:"maxSteps,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AppInitialActivity")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AppInitialActivity") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *AndroidRoboTest) MarshalJSON() ([]byte, error) {
	type noMethod AndroidRoboTest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// AndroidTest: An Android mobile test specification.
type AndroidTest struct {
	// AndroidAppInfo: Infomation about the application under test.
	AndroidAppInfo *AndroidAppInfo `json:"androidAppInfo,omitempty"`

	// AndroidInstrumentationTest: An Android instrumentation test.
	AndroidInstrumentationTest *AndroidInstrumentationTest `json:"androidInstrumentationTest,omitempty"`

	// AndroidRoboTest: An Android robo test.
	AndroidRoboTest *AndroidRoboTest `json:"androidRoboTest,omitempty"`

	// TestTimeout: Max time a test is allowed to run before it is
	// automatically cancelled.
	TestTimeout *Duration `json:"testTimeout,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AndroidAppInfo") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AndroidAppInfo") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *AndroidTest) MarshalJSON() ([]byte, error) {
	type noMethod AndroidTest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Any: `Any` contains an arbitrary serialized protocol buffer message
// along with a URL that describes the type of the serialized
// message.
//
// Protobuf library provides support to pack/unpack Any values in the
// form of utility functions or additional generated methods of the Any
// type.
//
// Example 1: Pack and unpack a message in C++.
//
// Foo foo = ...; Any any; any.PackFrom(foo); ... if
// (any.UnpackTo(&foo)) { ... }
//
// Example 2: Pack and unpack a message in Java.
//
// Foo foo = ...; Any any = Any.pack(foo); ... if (any.is(Foo.class)) {
// foo = any.unpack(Foo.class); }
//
// Example 3: Pack and unpack a message in Python.
//
// foo = Foo(...) any = Any() any.Pack(foo) ... if
// any.Is(Foo.DESCRIPTOR): any.Unpack(foo) ...
//
// Example 4: Pack and unpack a message in Go
//
// foo := &pb.Foo{...} any, err := ptypes.MarshalAny(foo) ... foo :=
// &pb.Foo{} if err := ptypes.UnmarshalAny(any, foo); err != nil { ...
// }
//
// The pack methods provided by protobuf library will by default use
// 'type.googleapis.com/full.type.name' as the type URL and the unpack
// methods only use the fully qualified type name after the last '/' in
// the type URL, for example "foo.bar.com/x/y.z" will yield type name
// "y.z".
//
//
//
// JSON ==== The JSON representation of an `Any` value uses the regular
// representation of the deserialized, embedded message, with an
// additional field `@type` which contains the type URL.
// Example:
//
// package google.profile; message Person { string first_name = 1;
// string last_name = 2; }
//
// { "@type": "type.googleapis.com/google.profile.Person", "firstName":
// , "lastName":  }
//
// If the embedded message type is well-known and has a custom JSON
// representation, that representation will be embedded adding a field
// `value` which holds the custom JSON in addition to the `@type` field.
// Example (for message [google.protobuf.Duration][]):
//
// { "@type": "type.googleapis.com/google.protobuf.Duration", "value":
// "1.212s" }
type Any struct {
	// TypeUrl: A URL/resource name whose content describes the type of the
	// serialized protocol buffer message.
	//
	// For URLs which use the scheme `http`, `https`, or no scheme, the
	// following restrictions and interpretations apply:
	//
	// * If no scheme is provided, `https` is assumed. * The last segment of
	// the URL's path must represent the fully qualified name of the type
	// (as in `path/google.protobuf.Duration`). The name should be in a
	// canonical form (e.g., leading "." is not accepted). * An HTTP GET on
	// the URL must yield a [google.protobuf.Type][] value in binary format,
	// or produce an error. * Applications are allowed to cache lookup
	// results based on the URL, or have them precompiled into a binary to
	// avoid any lookup. Therefore, binary compatibility needs to be
	// preserved on changes to types. (Use versioned type names to manage
	// breaking changes.)
	//
	// Schemes other than `http`, `https` (or the empty scheme) might be
	// used with implementation specific semantics.
	TypeUrl string `json:"typeUrl,omitempty"`

	// Value: Must be a valid serialized protocol buffer of the above
	// specified type.
	Value string `json:"value,omitempty"`

	// ForceSendFields is a list of field names (e.g. "TypeUrl") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "TypeUrl") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Any) MarshalJSON() ([]byte, error) {
	type noMethod Any
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type AppStartTime struct {
	// FullyDrawnTime: Optional. The time from app start to reaching the
	// developer-reported "fully drawn" time. This is only stored if the app
	// includes a call to Activity.reportFullyDrawn(). See
	// https://developer.android.com/topic/performance/launch-time.html#time-full
	FullyDrawnTime *Duration `json:"fullyDrawnTime,omitempty"`

	// InitialDisplayTime: The time from app start to the first displayed
	// activity being drawn, as reported in Logcat. See
	// https://developer.android.com/topic/performance/launch-time.html#time-initial
	InitialDisplayTime *Duration `json:"initialDisplayTime,omitempty"`

	// ForceSendFields is a list of field names (e.g. "FullyDrawnTime") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "FullyDrawnTime") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *AppStartTime) MarshalJSON() ([]byte, error) {
	type noMethod AppStartTime
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// BasicPerfSampleSeries: Encapsulates the metadata for basic sample
// series represented by a line chart
type BasicPerfSampleSeries struct {
	// Possible values:
	//   "cpu"
	//   "graphics"
	//   "memory"
	//   "network"
	//   "perfMetricTypeUnspecified"
	PerfMetricType string `json:"perfMetricType,omitempty"`

	// Possible values:
	//   "byte"
	//   "bytesPerSecond"
	//   "framesPerSecond"
	//   "kibibyte"
	//   "percent"
	//   "perfUnitUnspecified"
	PerfUnit string `json:"perfUnit,omitempty"`

	// Possible values:
	//   "cpuKernel"
	//   "cpuTotal"
	//   "cpuUser"
	//   "graphicsFrameRate"
	//   "memoryRssPrivate"
	//   "memoryRssShared"
	//   "memoryRssTotal"
	//   "memoryTotal"
	//   "networkReceived"
	//   "networkSent"
	//   "ntBytesReceived"
	//   "ntBytesTransferred"
	//   "sampleSeriesTypeUnspecified"
	SampleSeriesLabel string `json:"sampleSeriesLabel,omitempty"`

	// ForceSendFields is a list of field names (e.g. "PerfMetricType") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "PerfMetricType") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *BasicPerfSampleSeries) MarshalJSON() ([]byte, error) {
	type noMethod BasicPerfSampleSeries
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// BatchCreatePerfSamplesRequest: The request must provide up to a
// maximum of 5000 samples to be created; a larger sample size will
// cause an INVALID_ARGUMENT error
type BatchCreatePerfSamplesRequest struct {
	// PerfSamples: The set of PerfSamples to create should not include
	// existing timestamps
	PerfSamples []*PerfSample `json:"perfSamples,omitempty"`

	// ForceSendFields is a list of field names (e.g. "PerfSamples") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "PerfSamples") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *BatchCreatePerfSamplesRequest) MarshalJSON() ([]byte, error) {
	type noMethod BatchCreatePerfSamplesRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type BatchCreatePerfSamplesResponse struct {
	PerfSamples []*PerfSample `json:"perfSamples,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "PerfSamples") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "PerfSamples") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *BatchCreatePerfSamplesResponse) MarshalJSON() ([]byte, error) {
	type noMethod BatchCreatePerfSamplesResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type CPUInfo struct {
	// CpuProcessor: description of the device processor ie '1.8 GHz hexa
	// core 64-bit ARMv8-A'
	CpuProcessor string `json:"cpuProcessor,omitempty"`

	// CpuSpeedInGhz: the CPU clock speed in GHz
	CpuSpeedInGhz float64 `json:"cpuSpeedInGhz,omitempty"`

	// NumberOfCores: the number of CPU cores
	NumberOfCores int64 `json:"numberOfCores,omitempty"`

	// ForceSendFields is a list of field names (e.g. "CpuProcessor") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CpuProcessor") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *CPUInfo) MarshalJSON() ([]byte, error) {
	type noMethod CPUInfo
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

func (s *CPUInfo) UnmarshalJSON(data []byte) error {
	type noMethod CPUInfo
	var s1 struct {
		CpuSpeedInGhz gensupport.JSONFloat64 `json:"cpuSpeedInGhz"`
		*noMethod
	}
	s1.noMethod = (*noMethod)(s)
	if err := json.Unmarshal(data, &s1); err != nil {
		return err
	}
	s.CpuSpeedInGhz = float64(s1.CpuSpeedInGhz)
	return nil
}

// Duration: A Duration represents a signed, fixed-length span of time
// represented as a count of seconds and fractions of seconds at
// nanosecond resolution. It is independent of any calendar and concepts
// like "day" or "month". It is related to Timestamp in that the
// difference between two Timestamp values is a Duration and it can be
// added or subtracted from a Timestamp. Range is approximately +-10,000
// years.
//
// # Examples
//
// Example 1: Compute Duration from two Timestamps in pseudo
// code.
//
// Timestamp start = ...; Timestamp end = ...; Duration duration =
// ...;
//
// duration.seconds = end.seconds - start.seconds; duration.nanos =
// end.nanos - start.nanos;
//
// if (duration.seconds  0) { duration.seconds += 1; duration.nanos -=
// 1000000000; } else if (durations.seconds > 0 && duration.nanos < 0) {
// duration.seconds -= 1; duration.nanos += 1000000000; }
//
// Example 2: Compute Timestamp from Timestamp + Duration in pseudo
// code.
//
// Timestamp start = ...; Duration duration = ...; Timestamp end =
// ...;
//
// end.seconds = start.seconds + duration.seconds; end.nanos =
// start.nanos + duration.nanos;
//
// if (end.nanos = 1000000000) { end.seconds += 1; end.nanos -=
// 1000000000; }
//
// Example 3: Compute Duration from datetime.timedelta in Python.
//
// td = datetime.timedelta(days=3, minutes=10) duration = Duration()
// duration.FromTimedelta(td)
//
// # JSON Mapping
//
// In JSON format, the Duration type is encoded as a string rather than
// an object, where the string ends in the suffix "s" (indicating
// seconds) and is preceded by the number of seconds, with nanoseconds
// expressed as fractional seconds. For example, 3 seconds with 0
// nanoseconds should be encoded in JSON format as "3s", while 3 seconds
// and 1 nanosecond should be expressed in JSON format as
// "3.000000001s", and 3 seconds and 1 microsecond should be expressed
// in JSON format as "3.000001s".
type Duration struct {
	// Nanos: Signed fractions of a second at nanosecond resolution of the
	// span of time. Durations less than one second are represented with a 0
	// `seconds` field and a positive or negative `nanos` field. For
	// durations of one second or more, a non-zero value for the `nanos`
	// field must be of the same sign as the `seconds` field. Must be from
	// -999,999,999 to +999,999,999 inclusive.
	Nanos int64 `json:"nanos,omitempty"`

	// Seconds: Signed seconds of the span of time. Must be from
	// -315,576,000,000 to +315,576,000,000 inclusive. Note: these bounds
	// are computed from: 60 sec/min * 60 min/hr * 24 hr/day * 365.25
	// days/year * 10000 years
	Seconds int64 `json:"seconds,omitempty,string"`

	// ForceSendFields is a list of field names (e.g. "Nanos") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Nanos") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Duration) MarshalJSON() ([]byte, error) {
	type noMethod Duration
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Execution: An Execution represents a collection of Steps. For
// instance, it could represent: - a mobile test executed across a range
// of device configurations - a jenkins job with a build step followed
// by a test step
//
// The maximum size of an execution message is 1 MiB.
//
// An Execution can be updated until its state is set to COMPLETE at
// which point it becomes immutable.
type Execution struct {
	// CompletionTime: The time when the Execution status transitioned to
	// COMPLETE.
	//
	// This value will be set automatically when state transitions to
	// COMPLETE.
	//
	// - In response: set if the execution state is COMPLETE. - In
	// create/update request: never set
	CompletionTime *Timestamp `json:"completionTime,omitempty"`

	// CreationTime: The time when the Execution was created.
	//
	// This value will be set automatically when CreateExecution is
	// called.
	//
	// - In response: always set - In create/update request: never set
	CreationTime *Timestamp `json:"creationTime,omitempty"`

	// ExecutionId: A unique identifier within a History for this
	// Execution.
	//
	// Returns INVALID_ARGUMENT if this field is set or overwritten by the
	// caller.
	//
	// - In response always set - In create/update request: never set
	ExecutionId string `json:"executionId,omitempty"`

	// Outcome: Classify the result, for example into SUCCESS or FAILURE
	//
	// - In response: present if set by create/update request - In
	// create/update request: optional
	Outcome *Outcome `json:"outcome,omitempty"`

	// Specification: Lightweight information about execution request.
	//
	// - In response: present if set by create - In create: optional - In
	// update: optional
	Specification *Specification `json:"specification,omitempty"`

	// State: The initial state is IN_PROGRESS.
	//
	// The only legal state transitions is from IN_PROGRESS to COMPLETE.
	//
	// A PRECONDITION_FAILED will be returned if an invalid transition is
	// requested.
	//
	// The state can only be set to COMPLETE once. A FAILED_PRECONDITION
	// will be returned if the state is set to COMPLETE multiple times.
	//
	// If the state is set to COMPLETE, all the in-progress steps within the
	// execution will be set as COMPLETE. If the outcome of the step is not
	// set, the outcome will be set to INCONCLUSIVE.
	//
	// - In response always set - In create/update request: optional
	//
	// Possible values:
	//   "complete"
	//   "inProgress"
	//   "pending"
	//   "unknownState"
	State string `json:"state,omitempty"`

	// TestExecutionMatrixId: TestExecution Matrix ID that the
	// TestExecutionService uses.
	//
	// - In response: present if set by create - In create: optional - In
	// update: never set
	TestExecutionMatrixId string `json:"testExecutionMatrixId,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "CompletionTime") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CompletionTime") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *Execution) MarshalJSON() ([]byte, error) {
	type noMethod Execution
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type FailureDetail struct {
	// Crashed: If the failure was severe because the system (app) under
	// test crashed.
	Crashed bool `json:"crashed,omitempty"`

	// NotInstalled: If an app is not installed and thus no test can be run
	// with the app. This might be caused by trying to run a test on an
	// unsupported platform.
	NotInstalled bool `json:"notInstalled,omitempty"`

	// OtherNativeCrash: If a native process (including any other than the
	// app) crashed.
	OtherNativeCrash bool `json:"otherNativeCrash,omitempty"`

	// TimedOut: If the test overran some time limit, and that is why it
	// failed.
	TimedOut bool `json:"timedOut,omitempty"`

	// UnableToCrawl: If the robo was unable to crawl the app; perhaps
	// because the app did not start.
	UnableToCrawl bool `json:"unableToCrawl,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Crashed") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Crashed") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *FailureDetail) MarshalJSON() ([]byte, error) {
	type noMethod FailureDetail
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// FileReference: A reference to a file.
type FileReference struct {
	// FileUri: The URI of a file stored in Google Cloud Storage.
	//
	// For example: http://storage.googleapis.com/mybucket/path/to/test.xml
	// or in gsutil format: gs://mybucket/path/to/test.xml with
	// version-specific info,
	// gs://mybucket/path/to/test.xml#1360383693690000
	//
	// An INVALID_ARGUMENT error will be returned if the URI format is not
	// supported.
	//
	// - In response: always set - In create/update request: always set
	FileUri string `json:"fileUri,omitempty"`

	// ForceSendFields is a list of field names (e.g. "FileUri") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "FileUri") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *FileReference) MarshalJSON() ([]byte, error) {
	type noMethod FileReference
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// GraphicsStats: Graphics statistics for the App. The information is
// collected from 'adb shell dumpsys graphicsstats'. For more info see:
// https://developer.android.com/training/testing/performance.html
// Statistics will only be present for API 23+.
type GraphicsStats struct {
	// Buckets: Histogram of frame render times. There should be 154 buckets
	// ranging from [5ms, 6ms) to [4950ms, infinity)
	Buckets []*GraphicsStatsBucket `json:"buckets,omitempty"`

	// HighInputLatencyCount: Total "high input latency" events.
	HighInputLatencyCount int64 `json:"highInputLatencyCount,omitempty,string"`

	// JankyFrames: Total frames with slow render time. Should be <=
	// total_frames.
	JankyFrames int64 `json:"jankyFrames,omitempty,string"`

	// MissedVsyncCount: Total "missed vsync" events.
	MissedVsyncCount int64 `json:"missedVsyncCount,omitempty,string"`

	// P50Millis: 50th percentile frame render time in milliseconds.
	P50Millis int64 `json:"p50Millis,omitempty,string"`

	// P90Millis: 90th percentile frame render time in milliseconds.
	P90Millis int64 `json:"p90Millis,omitempty,string"`

	// P95Millis: 95th percentile frame render time in milliseconds.
	P95Millis int64 `json:"p95Millis,omitempty,string"`

	// P99Millis: 99th percentile frame render time in milliseconds.
	P99Millis int64 `json:"p99Millis,omitempty,string"`

	// SlowBitmapUploadCount: Total "slow bitmap upload" events.
	SlowBitmapUploadCount int64 `json:"slowBitmapUploadCount,omitempty,string"`

	// SlowDrawCount: Total "slow draw" events.
	SlowDrawCount int64 `json:"slowDrawCount,omitempty,string"`

	// SlowUiThreadCount: Total "slow UI thread" events.
	SlowUiThreadCount int64 `json:"slowUiThreadCount,omitempty,string"`

	// TotalFrames: Total frames rendered by package.
	TotalFrames int64 `json:"totalFrames,omitempty,string"`

	// ForceSendFields is a list of field names (e.g. "Buckets") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Buckets") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *GraphicsStats) MarshalJSON() ([]byte, error) {
	type noMethod GraphicsStats
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type GraphicsStatsBucket struct {
	// FrameCount: Number of frames in the bucket.
	FrameCount int64 `json:"frameCount,omitempty,string"`

	// RenderMillis: Lower bound of render time in milliseconds.
	RenderMillis int64 `json:"renderMillis,omitempty,string"`

	// ForceSendFields is a list of field names (e.g. "FrameCount") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "FrameCount") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *GraphicsStatsBucket) MarshalJSON() ([]byte, error) {
	type noMethod GraphicsStatsBucket
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// History: A History represents a sorted list of Executions ordered by
// the start_timestamp_millis field (descending). It can be used to
// group all the Executions of a continuous build.
//
// Note that the ordering only operates on one-dimension. If a
// repository has multiple branches, it means that multiple histories
// will need to be used in order to order Executions per branch.
type History struct {
	// DisplayName: A short human-readable (plain text) name to display in
	// the UI. Maximum of 100 characters.
	//
	// - In response: present if set during create. - In create request:
	// optional
	DisplayName string `json:"displayName,omitempty"`

	// HistoryId: A unique identifier within a project for this
	// History.
	//
	// Returns INVALID_ARGUMENT if this field is set or overwritten by the
	// caller.
	//
	// - In response always set - In create request: never set
	HistoryId string `json:"historyId,omitempty"`

	// Name: A name to uniquely identify a history within a project. Maximum
	// of 100 characters.
	//
	// - In response always set - In create request: always set
	Name string `json:"name,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "DisplayName") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "DisplayName") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *History) MarshalJSON() ([]byte, error) {
	type noMethod History
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Image: An image, with a link to the main image and a thumbnail.
type Image struct {
	// Error: An error explaining why the thumbnail could not be rendered.
	Error *Status `json:"error,omitempty"`

	// SourceImage: A reference to the full-size, original image.
	//
	// This is the same as the tool_outputs entry for the image under its
	// Step.
	//
	// Always set.
	SourceImage *ToolOutputReference `json:"sourceImage,omitempty"`

	// StepId: The step to which the image is attached.
	//
	// Always set.
	StepId string `json:"stepId,omitempty"`

	// Thumbnail: The thumbnail.
	Thumbnail *Thumbnail `json:"thumbnail,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Error") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Error") to include in API
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

type InconclusiveDetail struct {
	// AbortedByUser: If the end user aborted the test execution before a
	// pass or fail could be determined. For example, the user pressed
	// ctrl-c which sent a kill signal to the test runner while the test was
	// running.
	AbortedByUser bool `json:"abortedByUser,omitempty"`

	// InfrastructureFailure: If the test runner could not determine success
	// or failure because the test depends on a component other than the
	// system under test which failed.
	//
	// For example, a mobile test requires provisioning a device where the
	// test executes, and that provisioning can fail.
	InfrastructureFailure bool `json:"infrastructureFailure,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AbortedByUser") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AbortedByUser") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *InconclusiveDetail) MarshalJSON() ([]byte, error) {
	type noMethod InconclusiveDetail
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ListExecutionsResponse struct {
	// Executions: Executions.
	//
	// Always set.
	Executions []*Execution `json:"executions,omitempty"`

	// NextPageToken: A continuation token to resume the query at the next
	// item.
	//
	// Will only be set if there are more Executions to fetch.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Executions") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Executions") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListExecutionsResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListExecutionsResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListHistoriesResponse: Response message for HistoryService.List
type ListHistoriesResponse struct {
	// Histories: Histories.
	Histories []*History `json:"histories,omitempty"`

	// NextPageToken: A continuation token to resume the query at the next
	// item.
	//
	// Will only be set if there are more histories to fetch.
	//
	// Tokens are valid for up to one hour from the time of the first list
	// request. For instance, if you make a list request at 1PM and use the
	// token from this first request 10 minutes later, the token from this
	// second response will only be valid for 50 minutes.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Histories") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Histories") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListHistoriesResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListHistoriesResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ListPerfSampleSeriesResponse struct {
	// PerfSampleSeries: The resulting PerfSampleSeries sorted by id
	PerfSampleSeries []*PerfSampleSeries `json:"perfSampleSeries,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "PerfSampleSeries") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "PerfSampleSeries") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *ListPerfSampleSeriesResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListPerfSampleSeriesResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ListPerfSamplesResponse struct {
	// NextPageToken: Optional, returned if result size exceeds the page
	// size specified in the request (or the default page size, 500, if
	// unspecified). It indicates the last sample timestamp to be used as
	// page_token in subsequent request
	NextPageToken string `json:"nextPageToken,omitempty"`

	PerfSamples []*PerfSample `json:"perfSamples,omitempty"`

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

func (s *ListPerfSamplesResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListPerfSamplesResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ListScreenshotClustersResponse struct {
	// Clusters: The set of clustres associated with an execution Always set
	Clusters []*ScreenshotCluster `json:"clusters,omitempty"`

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

func (s *ListScreenshotClustersResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListScreenshotClustersResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListStepThumbnailsResponse: A response containing the thumbnails in a
// step.
type ListStepThumbnailsResponse struct {
	// NextPageToken: A continuation token to resume the query at the next
	// item.
	//
	// If set, indicates that there are more thumbnails to read, by calling
	// list again with this value in the page_token field.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// Thumbnails: A list of image data.
	//
	// Images are returned in a deterministic order; they are ordered by
	// these factors, in order of importance: * First, by their associated
	// test case. Images without a test case are considered greater than
	// images with one. * Second, by their creation time. Images without a
	// creation time are greater than images with one. * Third, by the order
	// in which they were added to the step (by calls to CreateStep or
	// UpdateStep).
	Thumbnails []*Image `json:"thumbnails,omitempty"`

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

func (s *ListStepThumbnailsResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListStepThumbnailsResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ListStepsResponse: Response message for StepService.List.
type ListStepsResponse struct {
	// NextPageToken: A continuation token to resume the query at the next
	// item.
	//
	// If set, indicates that there are more steps to read, by calling list
	// again with this value in the page_token field.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// Steps: Steps.
	Steps []*Step `json:"steps,omitempty"`

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

func (s *ListStepsResponse) MarshalJSON() ([]byte, error) {
	type noMethod ListStepsResponse
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type MemoryInfo struct {
	// MemoryCapInKibibyte: Maximum memory that can be allocated to the
	// process in KiB
	MemoryCapInKibibyte int64 `json:"memoryCapInKibibyte,omitempty,string"`

	// MemoryTotalInKibibyte: Total memory available on the device in KiB
	MemoryTotalInKibibyte int64 `json:"memoryTotalInKibibyte,omitempty,string"`

	// ForceSendFields is a list of field names (e.g. "MemoryCapInKibibyte")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "MemoryCapInKibibyte") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *MemoryInfo) MarshalJSON() ([]byte, error) {
	type noMethod MemoryInfo
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Outcome: Interprets a result so that humans and machines can act on
// it.
type Outcome struct {
	// FailureDetail: More information about a FAILURE outcome.
	//
	// Returns INVALID_ARGUMENT if this field is set but the summary is not
	// FAILURE.
	//
	// Optional
	FailureDetail *FailureDetail `json:"failureDetail,omitempty"`

	// InconclusiveDetail: More information about an INCONCLUSIVE
	// outcome.
	//
	// Returns INVALID_ARGUMENT if this field is set but the summary is not
	// INCONCLUSIVE.
	//
	// Optional
	InconclusiveDetail *InconclusiveDetail `json:"inconclusiveDetail,omitempty"`

	// SkippedDetail: More information about a SKIPPED outcome.
	//
	// Returns INVALID_ARGUMENT if this field is set but the summary is not
	// SKIPPED.
	//
	// Optional
	SkippedDetail *SkippedDetail `json:"skippedDetail,omitempty"`

	// SuccessDetail: More information about a SUCCESS outcome.
	//
	// Returns INVALID_ARGUMENT if this field is set but the summary is not
	// SUCCESS.
	//
	// Optional
	SuccessDetail *SuccessDetail `json:"successDetail,omitempty"`

	// Summary: The simplest way to interpret a result.
	//
	// Required
	//
	// Possible values:
	//   "failure"
	//   "inconclusive"
	//   "skipped"
	//   "success"
	//   "unset"
	Summary string `json:"summary,omitempty"`

	// ForceSendFields is a list of field names (e.g. "FailureDetail") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "FailureDetail") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Outcome) MarshalJSON() ([]byte, error) {
	type noMethod Outcome
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// PerfEnvironment: Encapsulates performance environment info
type PerfEnvironment struct {
	// CpuInfo: CPU related environment info
	CpuInfo *CPUInfo `json:"cpuInfo,omitempty"`

	// MemoryInfo: Memory related environment info
	MemoryInfo *MemoryInfo `json:"memoryInfo,omitempty"`

	// ForceSendFields is a list of field names (e.g. "CpuInfo") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CpuInfo") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *PerfEnvironment) MarshalJSON() ([]byte, error) {
	type noMethod PerfEnvironment
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// PerfMetricsSummary: A summary of perf metrics collected and
// performance environment info
type PerfMetricsSummary struct {
	AppStartTime *AppStartTime `json:"appStartTime,omitempty"`

	// ExecutionId: A tool results execution ID.
	ExecutionId string `json:"executionId,omitempty"`

	// GraphicsStats: Graphics statistics for the entire run. Statistics are
	// reset at the beginning of the run and collected at the end of the
	// run.
	GraphicsStats *GraphicsStats `json:"graphicsStats,omitempty"`

	// HistoryId: A tool results history ID.
	HistoryId string `json:"historyId,omitempty"`

	// PerfEnvironment: Describes the environment in which the performance
	// metrics were collected
	PerfEnvironment *PerfEnvironment `json:"perfEnvironment,omitempty"`

	// PerfMetrics: Set of resource collected
	//
	// Possible values:
	//   "cpu"
	//   "graphics"
	//   "memory"
	//   "network"
	//   "perfMetricTypeUnspecified"
	PerfMetrics []string `json:"perfMetrics,omitempty"`

	// ProjectId: The cloud project
	ProjectId string `json:"projectId,omitempty"`

	// StepId: A tool results step ID.
	StepId string `json:"stepId,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "AppStartTime") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AppStartTime") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *PerfMetricsSummary) MarshalJSON() ([]byte, error) {
	type noMethod PerfMetricsSummary
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// PerfSample: Resource representing a single performance measure or
// data point
type PerfSample struct {
	// SampleTime: Timestamp of collection
	SampleTime *Timestamp `json:"sampleTime,omitempty"`

	// Value: Value observed
	Value float64 `json:"value,omitempty"`

	// ForceSendFields is a list of field names (e.g. "SampleTime") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "SampleTime") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *PerfSample) MarshalJSON() ([]byte, error) {
	type noMethod PerfSample
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

func (s *PerfSample) UnmarshalJSON(data []byte) error {
	type noMethod PerfSample
	var s1 struct {
		Value gensupport.JSONFloat64 `json:"value"`
		*noMethod
	}
	s1.noMethod = (*noMethod)(s)
	if err := json.Unmarshal(data, &s1); err != nil {
		return err
	}
	s.Value = float64(s1.Value)
	return nil
}

// PerfSampleSeries: Resource representing a collection of performance
// samples (or data points)
type PerfSampleSeries struct {
	// BasicPerfSampleSeries: Basic series represented by a line chart
	BasicPerfSampleSeries *BasicPerfSampleSeries `json:"basicPerfSampleSeries,omitempty"`

	// ExecutionId: A tool results execution ID.
	ExecutionId string `json:"executionId,omitempty"`

	// HistoryId: A tool results history ID.
	HistoryId string `json:"historyId,omitempty"`

	// ProjectId: The cloud project
	ProjectId string `json:"projectId,omitempty"`

	// SampleSeriesId: A sample series id
	SampleSeriesId string `json:"sampleSeriesId,omitempty"`

	// StepId: A tool results step ID.
	StepId string `json:"stepId,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g.
	// "BasicPerfSampleSeries") to unconditionally include in API requests.
	// By default, fields with empty values are omitted from API requests.
	// However, any non-pointer, non-interface field appearing in
	// ForceSendFields will be sent to the server regardless of whether the
	// field is empty or not. This may be used to include empty fields in
	// Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "BasicPerfSampleSeries") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *PerfSampleSeries) MarshalJSON() ([]byte, error) {
	type noMethod PerfSampleSeries
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ProjectSettings: Per-project settings for the Tool Results service.
type ProjectSettings struct {
	// DefaultBucket: The name of the Google Cloud Storage bucket to which
	// results are written.
	//
	// By default, this is unset.
	//
	// In update request: optional In response: optional
	DefaultBucket string `json:"defaultBucket,omitempty"`

	// Name: The name of the project's settings.
	//
	// Always of the form: projects/{project-id}/settings
	//
	// In update request: never set In response: always set
	Name string `json:"name,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "DefaultBucket") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "DefaultBucket") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ProjectSettings) MarshalJSON() ([]byte, error) {
	type noMethod ProjectSettings
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// PublishXunitXmlFilesRequest: Request message for
// StepService.PublishXunitXmlFiles.
type PublishXunitXmlFilesRequest struct {
	// XunitXmlFiles: URI of the Xunit XML files to publish.
	//
	// The maximum size of the file this reference is pointing to is
	// 50MB.
	//
	// Required.
	XunitXmlFiles []*FileReference `json:"xunitXmlFiles,omitempty"`

	// ForceSendFields is a list of field names (e.g. "XunitXmlFiles") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "XunitXmlFiles") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *PublishXunitXmlFilesRequest) MarshalJSON() ([]byte, error) {
	type noMethod PublishXunitXmlFilesRequest
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type Screen struct {
	// FileReference: File reference of the png file. Required.
	FileReference string `json:"fileReference,omitempty"`

	// Locale: Locale of the device that the screenshot was taken on.
	// Required.
	Locale string `json:"locale,omitempty"`

	// Model: Model of the device that the screenshot was taken on.
	// Required.
	Model string `json:"model,omitempty"`

	// Version: OS version of the device that the screenshot was taken on.
	// Required.
	Version string `json:"version,omitempty"`

	// ForceSendFields is a list of field names (e.g. "FileReference") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "FileReference") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Screen) MarshalJSON() ([]byte, error) {
	type noMethod Screen
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ScreenshotCluster struct {
	// Activity: A string that describes the activity of every screen in the
	// cluster.
	Activity string `json:"activity,omitempty"`

	// ClusterId: A unique identifier for the cluster.
	ClusterId string `json:"clusterId,omitempty"`

	// KeyScreen: A singular screen that represents the cluster as a whole.
	// This screen will act as the "cover" of the entire cluster. When users
	// look at the clusters, only the key screen from each cluster will be
	// shown. Which screen is the key screen is determined by the
	// ClusteringAlgorithm
	KeyScreen *Screen `json:"keyScreen,omitempty"`

	// Screens: Full list of screens.
	Screens []*Screen `json:"screens,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Activity") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Activity") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ScreenshotCluster) MarshalJSON() ([]byte, error) {
	type noMethod ScreenshotCluster
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type SkippedDetail struct {
	// IncompatibleAppVersion: If the App doesn't support the specific API
	// level.
	IncompatibleAppVersion bool `json:"incompatibleAppVersion,omitempty"`

	// IncompatibleArchitecture: If the App doesn't run on the specific
	// architecture, for example, x86.
	IncompatibleArchitecture bool `json:"incompatibleArchitecture,omitempty"`

	// IncompatibleDevice: If the requested OS version doesn't run on the
	// specific device model.
	IncompatibleDevice bool `json:"incompatibleDevice,omitempty"`

	// ForceSendFields is a list of field names (e.g.
	// "IncompatibleAppVersion") to unconditionally include in API requests.
	// By default, fields with empty values are omitted from API requests.
	// However, any non-pointer, non-interface field appearing in
	// ForceSendFields will be sent to the server regardless of whether the
	// field is empty or not. This may be used to include empty fields in
	// Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "IncompatibleAppVersion")
	// to include in API requests with the JSON null value. By default,
	// fields with empty values are omitted from API requests. However, any
	// field with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *SkippedDetail) MarshalJSON() ([]byte, error) {
	type noMethod SkippedDetail
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Specification: The details about how to run the execution.
type Specification struct {
	// AndroidTest: An Android mobile test execution specification.
	AndroidTest *AndroidTest `json:"androidTest,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AndroidTest") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AndroidTest") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Specification) MarshalJSON() ([]byte, error) {
	type noMethod Specification
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// StackTrace: A stacktrace.
type StackTrace struct {
	// ClusterId: Exception cluster ID
	ClusterId string `json:"clusterId,omitempty"`

	// Exception: The stack trace message.
	//
	// Required
	Exception string `json:"exception,omitempty"`

	// ReportId: Exception report ID
	ReportId string `json:"reportId,omitempty"`

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

func (s *StackTrace) MarshalJSON() ([]byte, error) {
	type noMethod StackTrace
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Status: The `Status` type defines a logical error model that is
// suitable for different programming environments, including REST APIs
// and RPC APIs. It is used by [gRPC](https://github.com/grpc). The
// error model is designed to be:
//
// - Simple to use and understand for most users - Flexible enough to
// meet unexpected needs
//
// # Overview
//
// The `Status` message contains three pieces of data: error code, error
// message, and error details. The error code should be an enum value of
// [google.rpc.Code][], but it may accept additional error codes if
// needed. The error message should be a developer-facing English
// message that helps developers *understand* and *resolve* the error.
// If a localized user-facing error message is needed, put the localized
// message in the error details or localize it in the client. The
// optional error details may contain arbitrary information about the
// error. There is a predefined set of error detail types in the package
// `google.rpc` that can be used for common error conditions.
//
// # Language mapping
//
// The `Status` message is the logical representation of the error
// model, but it is not necessarily the actual wire format. When the
// `Status` message is exposed in different client libraries and
// different wire protocols, it can be mapped differently. For example,
// it will likely be mapped to some exceptions in Java, but more likely
// mapped to some error codes in C.
//
// # Other uses
//
// The error model and the `Status` message can be used in a variety of
// environments, either with or without APIs, to provide a consistent
// developer experience across different environments.
//
// Example uses of this error model include:
//
// - Partial errors. If a service needs to return partial errors to the
// client, it may embed the `Status` in the normal response to indicate
// the partial errors.
//
// - Workflow errors. A typical workflow has multiple steps. Each step
// may have a `Status` message for error reporting.
//
// - Batch operations. If a client uses batch request and batch
// response, the `Status` message should be used directly inside batch
// response, one for each error sub-response.
//
// - Asynchronous operations. If an API call embeds asynchronous
// operation results in its response, the status of those operations
// should be represented directly using the `Status` message.
//
// - Logging. If some API errors are stored in logs, the message
// `Status` could be used directly after any stripping needed for
// security/privacy reasons.
type Status struct {
	// Code: The status code, which should be an enum value of
	// [google.rpc.Code][].
	Code int64 `json:"code,omitempty"`

	// Details: A list of messages that carry the error details. There is a
	// common set of message types for APIs to use.
	Details []*Any `json:"details,omitempty"`

	// Message: A developer-facing error message, which should be in
	// English. Any user-facing error message should be localized and sent
	// in the [google.rpc.Status.details][] field, or localized by the
	// client.
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

// Step: A Step represents a single operation performed as part of
// Execution. A step can be used to represent the execution of a tool (
// for example a test runner execution or an execution of a
// compiler).
//
// Steps can overlap (for instance two steps might have the same start
// time if some operations are done in parallel).
//
// Here is an example, let's consider that we have a continuous build is
// executing a test runner for each iteration. The workflow would look
// like: - user creates a Execution with id 1 - user creates an
// TestExecutionStep with id 100 for Execution 1 - user update
// TestExecutionStep with id 100 to add a raw xml log + the service
// parses the xml logs and returns a TestExecutionStep with updated
// TestResult(s). - user update the status of TestExecutionStep with id
// 100 to COMPLETE
//
// A Step can be updated until its state is set to COMPLETE at which
// points it becomes immutable.
type Step struct {
	// CompletionTime: The time when the step status was set to
	// complete.
	//
	// This value will be set automatically when state transitions to
	// COMPLETE.
	//
	// - In response: set if the execution state is COMPLETE. - In
	// create/update request: never set
	CompletionTime *Timestamp `json:"completionTime,omitempty"`

	// CreationTime: The time when the step was created.
	//
	// - In response: always set - In create/update request: never set
	CreationTime *Timestamp `json:"creationTime,omitempty"`

	// Description: A description of this tool For example: mvn clean
	// package -D skipTests=true
	//
	// - In response: present if set by create/update request - In
	// create/update request: optional
	Description string `json:"description,omitempty"`

	// DeviceUsageDuration: How much the device resource is used to perform
	// the test.
	//
	// This is the device usage used for billing purpose, which is different
	// from the run_duration, for example, infrastructure failure won't be
	// charged for device usage.
	//
	// PRECONDITION_FAILED will be returned if one attempts to set a
	// device_usage on a step which already has this field set.
	//
	// - In response: present if previously set. - In create request:
	// optional - In update request: optional
	DeviceUsageDuration *Duration `json:"deviceUsageDuration,omitempty"`

	// DimensionValue: If the execution containing this step has any
	// dimension_definition set, then this field allows the child to specify
	// the values of the dimensions.
	//
	// The keys must exactly match the dimension_definition of the
	// execution.
	//
	// For example, if the execution has `dimension_definition = ['attempt',
	// 'device']` then a step must define values for those dimensions, eg.
	// `dimension_value = ['attempt': '1', 'device': 'Nexus 6']`
	//
	// If a step does not participate in one dimension of the matrix, the
	// value for that dimension should be empty string. For example, if one
	// of the tests is executed by a runner which does not support retries,
	// the step could have `dimension_value = ['attempt': '', 'device':
	// 'Nexus 6']`
	//
	// If the step does not participate in any dimensions of the matrix, it
	// may leave dimension_value unset.
	//
	// A PRECONDITION_FAILED will be returned if any of the keys do not
	// exist in the dimension_definition of the execution.
	//
	// A PRECONDITION_FAILED will be returned if another step in this
	// execution already has the same name and dimension_value, but differs
	// on other data fields, for example, step field is different.
	//
	// A PRECONDITION_FAILED will be returned if dimension_value is set, and
	// there is a dimension_definition in the execution which is not
	// specified as one of the keys.
	//
	// - In response: present if set by create - In create request: optional
	// - In update request: never set
	DimensionValue []*StepDimensionValueEntry `json:"dimensionValue,omitempty"`

	// HasImages: Whether any of the outputs of this step are images whose
	// thumbnails can be fetched with ListThumbnails.
	//
	// - In response: always set - In create/update request: never set
	HasImages bool `json:"hasImages,omitempty"`

	// Labels: Arbitrary user-supplied key/value pairs that are associated
	// with the step.
	//
	// Users are responsible for managing the key namespace such that keys
	// don't accidentally collide.
	//
	// An INVALID_ARGUMENT will be returned if the number of labels exceeds
	// 100 or if the length of any of the keys or values exceeds 100
	// characters.
	//
	// - In response: always set - In create request: optional - In update
	// request: optional; any new key/value pair will be added to the map,
	// and any new value for an existing key will update that key's value
	Labels []*StepLabelsEntry `json:"labels,omitempty"`

	// Name: A short human-readable name to display in the UI. Maximum of
	// 100 characters. For example: Clean build
	//
	// A PRECONDITION_FAILED will be returned upon creating a new step if it
	// shares its name and dimension_value with an existing step. If two
	// steps represent a similar action, but have different dimension
	// values, they should share the same name. For instance, if the same
	// set of tests is run on two different platforms, the two steps should
	// have the same name.
	//
	// - In response: always set - In create request: always set - In update
	// request: never set
	Name string `json:"name,omitempty"`

	// Outcome: Classification of the result, for example into SUCCESS or
	// FAILURE
	//
	// - In response: present if set by create/update request - In
	// create/update request: optional
	Outcome *Outcome `json:"outcome,omitempty"`

	// RunDuration: How long it took for this step to run.
	//
	// If unset, this is set to the difference between creation_time and
	// completion_time when the step is set to the COMPLETE state. In some
	// cases, it is appropriate to set this value separately: For instance,
	// if a step is created, but the operation it represents is queued for a
	// few minutes before it executes, it would be appropriate not to
	// include the time spent queued in its
	// run_duration.
	//
	// PRECONDITION_FAILED will be returned if one attempts to set a
	// run_duration on a step which already has this field set.
	//
	// - In response: present if previously set; always present on COMPLETE
	// step - In create request: optional - In update request: optional
	RunDuration *Duration `json:"runDuration,omitempty"`

	// State: The initial state is IN_PROGRESS. The only legal state
	// transitions are * IN_PROGRESS -> COMPLETE
	//
	// A PRECONDITION_FAILED will be returned if an invalid transition is
	// requested.
	//
	// It is valid to create Step with a state set to COMPLETE. The state
	// can only be set to COMPLETE once. A PRECONDITION_FAILED will be
	// returned if the state is set to COMPLETE multiple times.
	//
	// - In response: always set - In create/update request: optional
	//
	// Possible values:
	//   "complete"
	//   "inProgress"
	//   "pending"
	//   "unknownState"
	State string `json:"state,omitempty"`

	// StepId: A unique identifier within a Execution for this
	// Step.
	//
	// Returns INVALID_ARGUMENT if this field is set or overwritten by the
	// caller.
	//
	// - In response: always set - In create/update request: never set
	StepId string `json:"stepId,omitempty"`

	// TestExecutionStep: An execution of a test runner.
	TestExecutionStep *TestExecutionStep `json:"testExecutionStep,omitempty"`

	// ToolExecutionStep: An execution of a tool (used for steps we don't
	// explicitly support).
	ToolExecutionStep *ToolExecutionStep `json:"toolExecutionStep,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "CompletionTime") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CompletionTime") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *Step) MarshalJSON() ([]byte, error) {
	type noMethod Step
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type StepDimensionValueEntry struct {
	Key string `json:"key,omitempty"`

	Value string `json:"value,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Key") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Key") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *StepDimensionValueEntry) MarshalJSON() ([]byte, error) {
	type noMethod StepDimensionValueEntry
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type StepLabelsEntry struct {
	Key string `json:"key,omitempty"`

	Value string `json:"value,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Key") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Key") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *StepLabelsEntry) MarshalJSON() ([]byte, error) {
	type noMethod StepLabelsEntry
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type SuccessDetail struct {
	// OtherNativeCrash: If a native process other than the app crashed.
	OtherNativeCrash bool `json:"otherNativeCrash,omitempty"`

	// ForceSendFields is a list of field names (e.g. "OtherNativeCrash") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "OtherNativeCrash") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *SuccessDetail) MarshalJSON() ([]byte, error) {
	type noMethod SuccessDetail
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// TestCaseReference: A reference to a test case.
//
// Test case references are canonically ordered lexicographically by
// these three factors: * First, by test_suite_name. * Second, by
// class_name. * Third, by name.
type TestCaseReference struct {
	// ClassName: The name of the class.
	ClassName string `json:"className,omitempty"`

	// Name: The name of the test case.
	//
	// Required.
	Name string `json:"name,omitempty"`

	// TestSuiteName: The name of the test suite to which this test case
	// belongs.
	TestSuiteName string `json:"testSuiteName,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ClassName") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ClassName") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *TestCaseReference) MarshalJSON() ([]byte, error) {
	type noMethod TestCaseReference
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// TestExecutionStep: A step that represents running tests.
//
// It accepts ant-junit xml files which will be parsed into structured
// test results by the service. Xml file paths are updated in order to
// append more files, however they can't be deleted.
//
// Users can also add test results manually by using the test_result
// field.
type TestExecutionStep struct {
	// TestIssues: Issues observed during the test execution.
	//
	// For example, if the mobile app under test crashed during the test,
	// the error message and the stack trace content can be recorded here to
	// assist debugging.
	//
	// - In response: present if set by create or update - In create/update
	// request: optional
	TestIssues []*TestIssue `json:"testIssues,omitempty"`

	// TestSuiteOverviews: List of test suite overview contents. This could
	// be parsed from xUnit XML log by server, or uploaded directly by user.
	// This references should only be called when test suites are fully
	// parsed or uploaded.
	//
	// The maximum allowed number of test suite overviews per step is
	// 1000.
	//
	// - In response: always set - In create request: optional - In update
	// request: never (use publishXunitXmlFiles custom method instead)
	TestSuiteOverviews []*TestSuiteOverview `json:"testSuiteOverviews,omitempty"`

	// TestTiming: The timing break down of the test execution.
	//
	// - In response: present if set by create or update - In create/update
	// request: optional
	TestTiming *TestTiming `json:"testTiming,omitempty"`

	// ToolExecution: Represents the execution of the test runner.
	//
	// The exit code of this tool will be used to determine if the test
	// passed.
	//
	// - In response: always set - In create/update request: optional
	ToolExecution *ToolExecution `json:"toolExecution,omitempty"`

	// ForceSendFields is a list of field names (e.g. "TestIssues") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "TestIssues") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *TestExecutionStep) MarshalJSON() ([]byte, error) {
	type noMethod TestExecutionStep
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// TestIssue: An issue detected occurring during a test execution.
type TestIssue struct {
	// ErrorMessage: A brief human-readable message describing the issue.
	// Required.
	ErrorMessage string `json:"errorMessage,omitempty"`

	// Severity: Severity of issue. Required.
	//
	// Possible values:
	//   "info"
	//   "severe"
	//   "unspecifiedSeverity"
	//   "warning"
	Severity string `json:"severity,omitempty"`

	// StackTrace: Deprecated in favor of stack trace fields inside specific
	// warnings.
	StackTrace *StackTrace `json:"stackTrace,omitempty"`

	// Type: Type of issue. Required.
	//
	// Possible values:
	//   "anr"
	//   "fatalException"
	//   "nativeCrash"
	//   "unspecifiedType"
	Type string `json:"type,omitempty"`

	// Warning: Warning message with additional details of the issue. Should
	// always be a message from com.google.devtools.toolresults.v1.warnings
	// Required.
	Warning *Any `json:"warning,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ErrorMessage") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ErrorMessage") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *TestIssue) MarshalJSON() ([]byte, error) {
	type noMethod TestIssue
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// TestSuiteOverview: A summary of a test suite result either parsed
// from XML or uploaded directly by a user.
//
// Note: the API related comments are for StepService only. This message
// is also being used in ExecutionService in a read only mode for the
// corresponding step.
type TestSuiteOverview struct {
	// ErrorCount: Number of test cases in error, typically set by the
	// service by parsing the xml_source.
	//
	// - In create/response: always set - In update request: never
	ErrorCount int64 `json:"errorCount,omitempty"`

	// FailureCount: Number of failed test cases, typically set by the
	// service by parsing the xml_source. May also be set by the user.
	//
	// - In create/response: always set - In update request: never
	FailureCount int64 `json:"failureCount,omitempty"`

	// Name: The name of the test suite.
	//
	// - In create/response: always set - In update request: never
	Name string `json:"name,omitempty"`

	// SkippedCount: Number of test cases not run, typically set by the
	// service by parsing the xml_source.
	//
	// - In create/response: always set - In update request: never
	SkippedCount int64 `json:"skippedCount,omitempty"`

	// TotalCount: Number of test cases, typically set by the service by
	// parsing the xml_source.
	//
	// - In create/response: always set - In update request: never
	TotalCount int64 `json:"totalCount,omitempty"`

	// XmlSource: If this test suite was parsed from XML, this is the URI
	// where the original XML file is stored.
	//
	// Note: Multiple test suites can share the same xml_source
	//
	// Returns INVALID_ARGUMENT if the uri format is not supported.
	//
	// - In create/response: optional - In update request: never
	XmlSource *FileReference `json:"xmlSource,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ErrorCount") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ErrorCount") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *TestSuiteOverview) MarshalJSON() ([]byte, error) {
	type noMethod TestSuiteOverview
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// TestTiming: Testing timing break down to know phases.
type TestTiming struct {
	// TestProcessDuration: How long it took to run the test process.
	//
	// - In response: present if previously set. - In create/update request:
	// optional
	TestProcessDuration *Duration `json:"testProcessDuration,omitempty"`

	// ForceSendFields is a list of field names (e.g. "TestProcessDuration")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "TestProcessDuration") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *TestTiming) MarshalJSON() ([]byte, error) {
	type noMethod TestTiming
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Thumbnail: A single thumbnail, with its size and format.
type Thumbnail struct {
	// ContentType: The thumbnail's content type, i.e. "image/png".
	//
	// Always set.
	ContentType string `json:"contentType,omitempty"`

	// Data: The thumbnail file itself.
	//
	// That is, the bytes here are precisely the bytes that make up the
	// thumbnail file; they can be served as an image as-is (with the
	// appropriate content type.)
	//
	// Always set.
	Data string `json:"data,omitempty"`

	// HeightPx: The height of the thumbnail, in pixels.
	//
	// Always set.
	HeightPx int64 `json:"heightPx,omitempty"`

	// WidthPx: The width of the thumbnail, in pixels.
	//
	// Always set.
	WidthPx int64 `json:"widthPx,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ContentType") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ContentType") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Thumbnail) MarshalJSON() ([]byte, error) {
	type noMethod Thumbnail
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Timestamp: A Timestamp represents a point in time independent of any
// time zone or calendar, represented as seconds and fractions of
// seconds at nanosecond resolution in UTC Epoch time. It is encoded
// using the Proleptic Gregorian Calendar which extends the Gregorian
// calendar backwards to year one. It is encoded assuming all minutes
// are 60 seconds long, i.e. leap seconds are "smeared" so that no leap
// second table is needed for interpretation. Range is from
// 0001-01-01T00:00:00Z to 9999-12-31T23:59:59.999999999Z. By
// restricting to that range, we ensure that we can convert to and from
// RFC 3339 date strings. See
// [https://www.ietf.org/rfc/rfc3339.txt](https://www.ietf.org/rfc/rfc333
// 9.txt).
//
// # Examples
//
// Example 1: Compute Timestamp from POSIX `time()`.
//
// Timestamp timestamp; timestamp.set_seconds(time(NULL));
// timestamp.set_nanos(0);
//
// Example 2: Compute Timestamp from POSIX `gettimeofday()`.
//
// struct timeval tv; gettimeofday(&tv, NULL);
//
// Timestamp timestamp; timestamp.set_seconds(tv.tv_sec);
// timestamp.set_nanos(tv.tv_usec * 1000);
//
// Example 3: Compute Timestamp from Win32
// `GetSystemTimeAsFileTime()`.
//
// FILETIME ft; GetSystemTimeAsFileTime(&ft); UINT64 ticks =
// (((UINT64)ft.dwHighDateTime) << 32) | ft.dwLowDateTime;
//
// // A Windows tick is 100 nanoseconds. Windows epoch
// 1601-01-01T00:00:00Z // is 11644473600 seconds before Unix epoch
// 1970-01-01T00:00:00Z. Timestamp timestamp;
// timestamp.set_seconds((INT64) ((ticks / 10000000) - 11644473600LL));
// timestamp.set_nanos((INT32) ((ticks % 10000000) * 100));
//
// Example 4: Compute Timestamp from Java
// `System.currentTimeMillis()`.
//
// long millis = System.currentTimeMillis();
//
// Timestamp timestamp = Timestamp.newBuilder().setSeconds(millis /
// 1000) .setNanos((int) ((millis % 1000) *
// 1000000)).build();
//
//
//
// Example 5: Compute Timestamp from current time in Python.
//
// timestamp = Timestamp() timestamp.GetCurrentTime()
//
// # JSON Mapping
//
// In JSON format, the Timestamp type is encoded as a string in the [RFC
// 3339](https://www.ietf.org/rfc/rfc3339.txt) format. That is, the
// format is "{year}-{month}-{day}T{hour}:{min}:{sec}[.{frac_sec}]Z"
// where {year} is always expressed using four digits while {month},
// {day}, {hour}, {min}, and {sec} are zero-padded to two digits each.
// The fractional seconds, which can go up to 9 digits (i.e. up to 1
// nanosecond resolution), are optional. The "Z" suffix indicates the
// timezone ("UTC"); the timezone is required, though only UTC (as
// indicated by "Z") is presently supported.
//
// For example, "2017-01-15T01:30:15.01Z" encodes 15.01 seconds past
// 01:30 UTC on January 15, 2017.
//
// In JavaScript, one can convert a Date object to this format using the
// standard
// [toISOString()](https://developer.mozilla.org/en-US/docs/Web/JavaScrip
// t/Reference/Global_Objects/Date/toISOString] method. In Python, a
// standard `datetime.datetime` object can be converted to this format
// using
// [`strftime`](https://docs.python.org/2/library/time.html#time.strftime
// ) with the time format spec '%Y-%m-%dT%H:%M:%S.%fZ'. Likewise, in
// Java, one can use the Joda Time's [`ISODateTimeFormat.dateTime()`](
// http://joda-time.sourceforge.net/apidocs/org/joda/time/format/ISODateTimeFormat.html#dateTime()) to obtain a formatter capable of generating timestamps in this
// format.
type Timestamp struct {
	// Nanos: Non-negative fractions of a second at nanosecond resolution.
	// Negative second values with fractions must still have non-negative
	// nanos values that count forward in time. Must be from 0 to
	// 999,999,999 inclusive.
	Nanos int64 `json:"nanos,omitempty"`

	// Seconds: Represents seconds of UTC time since Unix epoch
	// 1970-01-01T00:00:00Z. Must be from 0001-01-01T00:00:00Z to
	// 9999-12-31T23:59:59Z inclusive.
	Seconds int64 `json:"seconds,omitempty,string"`

	// ForceSendFields is a list of field names (e.g. "Nanos") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Nanos") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Timestamp) MarshalJSON() ([]byte, error) {
	type noMethod Timestamp
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ToolExecution: An execution of an arbitrary tool. It could be a test
// runner or a tool copying artifacts or deploying code.
type ToolExecution struct {
	// CommandLineArguments: The full tokenized command line including the
	// program name (equivalent to argv in a C program).
	//
	// - In response: present if set by create request - In create request:
	// optional - In update request: never set
	CommandLineArguments []string `json:"commandLineArguments,omitempty"`

	// ExitCode: Tool execution exit code. This field will be set once the
	// tool has exited.
	//
	// - In response: present if set by create/update request - In create
	// request: optional - In update request: optional, a
	// FAILED_PRECONDITION error will be returned if an exit_code is already
	// set.
	ExitCode *ToolExitCode `json:"exitCode,omitempty"`

	// ToolLogs: References to any plain text logs output the tool
	// execution.
	//
	// This field can be set before the tool has exited in order to be able
	// to have access to a live view of the logs while the tool is
	// running.
	//
	// The maximum allowed number of tool logs per step is 1000.
	//
	// - In response: present if set by create/update request - In create
	// request: optional - In update request: optional, any value provided
	// will be appended to the existing list
	ToolLogs []*FileReference `json:"toolLogs,omitempty"`

	// ToolOutputs: References to opaque files of any format output by the
	// tool execution.
	//
	// The maximum allowed number of tool outputs per step is 1000.
	//
	// - In response: present if set by create/update request - In create
	// request: optional - In update request: optional, any value provided
	// will be appended to the existing list
	ToolOutputs []*ToolOutputReference `json:"toolOutputs,omitempty"`

	// ForceSendFields is a list of field names (e.g.
	// "CommandLineArguments") to unconditionally include in API requests.
	// By default, fields with empty values are omitted from API requests.
	// However, any non-pointer, non-interface field appearing in
	// ForceSendFields will be sent to the server regardless of whether the
	// field is empty or not. This may be used to include empty fields in
	// Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CommandLineArguments") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *ToolExecution) MarshalJSON() ([]byte, error) {
	type noMethod ToolExecution
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ToolExecutionStep: Generic tool step to be used for binaries we do
// not explicitly support. For example: running cp to copy artifacts
// from one location to another.
type ToolExecutionStep struct {
	// ToolExecution: A Tool execution.
	//
	// - In response: present if set by create/update request - In
	// create/update request: optional
	ToolExecution *ToolExecution `json:"toolExecution,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ToolExecution") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ToolExecution") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ToolExecutionStep) MarshalJSON() ([]byte, error) {
	type noMethod ToolExecutionStep
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ToolExitCode: Exit code from a tool execution.
type ToolExitCode struct {
	// Number: Tool execution exit code. A value of 0 means that the
	// execution was successful.
	//
	// - In response: always set - In create/update request: always set
	Number int64 `json:"number,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Number") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Number") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ToolExitCode) MarshalJSON() ([]byte, error) {
	type noMethod ToolExitCode
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ToolOutputReference: A reference to a ToolExecution output file.
type ToolOutputReference struct {
	// CreationTime: The creation time of the file.
	//
	// - In response: present if set by create/update request - In
	// create/update request: optional
	CreationTime *Timestamp `json:"creationTime,omitempty"`

	// Output: A FileReference to an output file.
	//
	// - In response: always set - In create/update request: always set
	Output *FileReference `json:"output,omitempty"`

	// TestCase: The test case to which this output file belongs.
	//
	// - In response: present if set by create/update request - In
	// create/update request: optional
	TestCase *TestCaseReference `json:"testCase,omitempty"`

	// ForceSendFields is a list of field names (e.g. "CreationTime") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "CreationTime") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ToolOutputReference) MarshalJSON() ([]byte, error) {
	type noMethod ToolOutputReference
	raw := noMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// method id "toolresults.projects.getSettings":

type ProjectsGetSettingsCall struct {
	s            *Service
	projectId    string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// GetSettings: Gets the Tool Results settings for a project.
//
// May return any of the following canonical error codes:
//
// - PERMISSION_DENIED - if the user is not authorized to read from
// project
func (r *ProjectsService) GetSettings(projectId string) *ProjectsGetSettingsCall {
	c := &ProjectsGetSettingsCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsGetSettingsCall) Fields(s ...googleapi.Field) *ProjectsGetSettingsCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsGetSettingsCall) IfNoneMatch(entityTag string) *ProjectsGetSettingsCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsGetSettingsCall) Context(ctx context.Context) *ProjectsGetSettingsCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsGetSettingsCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsGetSettingsCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/settings")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.getSettings" call.
// Exactly one of *ProjectSettings or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *ProjectSettings.ServerResponse.Header or (if a response was returned
// at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsGetSettingsCall) Do(opts ...googleapi.CallOption) (*ProjectSettings, error) {
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
	ret := &ProjectSettings{
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
	//   "description": "Gets the Tool Results settings for a project.\n\nMay return any of the following canonical error codes:\n\n- PERMISSION_DENIED - if the user is not authorized to read from project",
	//   "httpMethod": "GET",
	//   "id": "toolresults.projects.getSettings",
	//   "parameterOrder": [
	//     "projectId"
	//   ],
	//   "parameters": {
	//     "projectId": {
	//       "description": "A Project id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/settings",
	//   "response": {
	//     "$ref": "ProjectSettings"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.initializeSettings":

type ProjectsInitializeSettingsCall struct {
	s          *Service
	projectId  string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// InitializeSettings: Creates resources for settings which have not yet
// been set.
//
// Currently, this creates a single resource: a Google Cloud Storage
// bucket, to be used as the default bucket for this project. The bucket
// is created in an FTL-own storage project. Except for in rare cases,
// calling this method in parallel from multiple clients will only
// create a single bucket. In order to avoid unnecessary storage
// charges, the bucket is configured to automatically delete objects
// older than 90 days.
//
// The bucket is created with the following permissions: - Owner access
// for owners of central storage project (FTL-owned) - Writer access for
// owners/editors of customer project - Reader access for viewers of
// customer project The default ACL on objects created in the bucket is:
// - Owner access for owners of central storage project - Reader access
// for owners/editors/viewers of customer project See Google Cloud
// Storage documentation for more details.
//
// If there is already a default bucket set and the project can access
// the bucket, this call does nothing. However, if the project doesn't
// have the permission to access the bucket or the bucket is deleted, a
// new bucket will be created.
//
// May return any canonical error codes, including the following:
//
// - PERMISSION_DENIED - if the user is not authorized to write to
// project - Any error code raised by Google Cloud Storage
func (r *ProjectsService) InitializeSettings(projectId string) *ProjectsInitializeSettingsCall {
	c := &ProjectsInitializeSettingsCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsInitializeSettingsCall) Fields(s ...googleapi.Field) *ProjectsInitializeSettingsCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsInitializeSettingsCall) Context(ctx context.Context) *ProjectsInitializeSettingsCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsInitializeSettingsCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsInitializeSettingsCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}:initializeSettings")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.initializeSettings" call.
// Exactly one of *ProjectSettings or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *ProjectSettings.ServerResponse.Header or (if a response was returned
// at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsInitializeSettingsCall) Do(opts ...googleapi.CallOption) (*ProjectSettings, error) {
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
	ret := &ProjectSettings{
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
	//   "description": "Creates resources for settings which have not yet been set.\n\nCurrently, this creates a single resource: a Google Cloud Storage bucket, to be used as the default bucket for this project. The bucket is created in an FTL-own storage project. Except for in rare cases, calling this method in parallel from multiple clients will only create a single bucket. In order to avoid unnecessary storage charges, the bucket is configured to automatically delete objects older than 90 days.\n\nThe bucket is created with the following permissions: - Owner access for owners of central storage project (FTL-owned) - Writer access for owners/editors of customer project - Reader access for viewers of customer project The default ACL on objects created in the bucket is: - Owner access for owners of central storage project - Reader access for owners/editors/viewers of customer project See Google Cloud Storage documentation for more details.\n\nIf there is already a default bucket set and the project can access the bucket, this call does nothing. However, if the project doesn't have the permission to access the bucket or the bucket is deleted, a new bucket will be created.\n\nMay return any canonical error codes, including the following:\n\n- PERMISSION_DENIED - if the user is not authorized to write to project - Any error code raised by Google Cloud Storage",
	//   "httpMethod": "POST",
	//   "id": "toolresults.projects.initializeSettings",
	//   "parameterOrder": [
	//     "projectId"
	//   ],
	//   "parameters": {
	//     "projectId": {
	//       "description": "A Project id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}:initializeSettings",
	//   "response": {
	//     "$ref": "ProjectSettings"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.histories.create":

type ProjectsHistoriesCreateCall struct {
	s          *Service
	projectId  string
	history    *History
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Create: Creates a History.
//
// The returned History will have the id set.
//
// May return any of the following canonical error codes:
//
// - PERMISSION_DENIED - if the user is not authorized to write to
// project - INVALID_ARGUMENT - if the request is malformed - NOT_FOUND
// - if the containing project does not exist
func (r *ProjectsHistoriesService) Create(projectId string, history *History) *ProjectsHistoriesCreateCall {
	c := &ProjectsHistoriesCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.history = history
	return c
}

// RequestId sets the optional parameter "requestId": A unique request
// ID for server to detect duplicated requests. For example, a
// UUID.
//
// Optional, but strongly recommended.
func (c *ProjectsHistoriesCreateCall) RequestId(requestId string) *ProjectsHistoriesCreateCall {
	c.urlParams_.Set("requestId", requestId)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesCreateCall) Fields(s ...googleapi.Field) *ProjectsHistoriesCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesCreateCall) Context(ctx context.Context) *ProjectsHistoriesCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.history)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.create" call.
// Exactly one of *History or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *History.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *ProjectsHistoriesCreateCall) Do(opts ...googleapi.CallOption) (*History, error) {
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
	ret := &History{
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
	//   "description": "Creates a History.\n\nThe returned History will have the id set.\n\nMay return any of the following canonical error codes:\n\n- PERMISSION_DENIED - if the user is not authorized to write to project - INVALID_ARGUMENT - if the request is malformed - NOT_FOUND - if the containing project does not exist",
	//   "httpMethod": "POST",
	//   "id": "toolresults.projects.histories.create",
	//   "parameterOrder": [
	//     "projectId"
	//   ],
	//   "parameters": {
	//     "projectId": {
	//       "description": "A Project id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "requestId": {
	//       "description": "A unique request ID for server to detect duplicated requests. For example, a UUID.\n\nOptional, but strongly recommended.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories",
	//   "request": {
	//     "$ref": "History"
	//   },
	//   "response": {
	//     "$ref": "History"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.histories.get":

type ProjectsHistoriesGetCall struct {
	s            *Service
	projectId    string
	historyId    string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Gets a History.
//
// May return any of the following canonical error codes:
//
// - PERMISSION_DENIED - if the user is not authorized to read project -
// INVALID_ARGUMENT - if the request is malformed - NOT_FOUND - if the
// History does not exist
func (r *ProjectsHistoriesService) Get(projectId string, historyId string) *ProjectsHistoriesGetCall {
	c := &ProjectsHistoriesGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesGetCall) Fields(s ...googleapi.Field) *ProjectsHistoriesGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsHistoriesGetCall) IfNoneMatch(entityTag string) *ProjectsHistoriesGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesGetCall) Context(ctx context.Context) *ProjectsHistoriesGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"historyId": c.historyId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.get" call.
// Exactly one of *History or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *History.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *ProjectsHistoriesGetCall) Do(opts ...googleapi.CallOption) (*History, error) {
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
	ret := &History{
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
	//   "description": "Gets a History.\n\nMay return any of the following canonical error codes:\n\n- PERMISSION_DENIED - if the user is not authorized to read project - INVALID_ARGUMENT - if the request is malformed - NOT_FOUND - if the History does not exist",
	//   "httpMethod": "GET",
	//   "id": "toolresults.projects.histories.get",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId"
	//   ],
	//   "parameters": {
	//     "historyId": {
	//       "description": "A History id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "A Project id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}",
	//   "response": {
	//     "$ref": "History"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.histories.list":

type ProjectsHistoriesListCall struct {
	s            *Service
	projectId    string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists Histories for a given Project.
//
// The histories are sorted by modification time in descending order.
// The history_id key will be used to order the history with the same
// modification time.
//
// May return any of the following canonical error codes:
//
// - PERMISSION_DENIED - if the user is not authorized to read project -
// INVALID_ARGUMENT - if the request is malformed - NOT_FOUND - if the
// History does not exist
func (r *ProjectsHistoriesService) List(projectId string) *ProjectsHistoriesListCall {
	c := &ProjectsHistoriesListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	return c
}

// FilterByName sets the optional parameter "filterByName": If set, only
// return histories with the given name.
func (c *ProjectsHistoriesListCall) FilterByName(filterByName string) *ProjectsHistoriesListCall {
	c.urlParams_.Set("filterByName", filterByName)
	return c
}

// PageSize sets the optional parameter "pageSize": The maximum number
// of Histories to fetch.
//
// Default value: 20. The server will use this default if the field is
// not set or has a value of 0. Any value greater than 100 will be
// treated as 100.
func (c *ProjectsHistoriesListCall) PageSize(pageSize int64) *ProjectsHistoriesListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken": A continuation
// token to resume the query at the next item.
func (c *ProjectsHistoriesListCall) PageToken(pageToken string) *ProjectsHistoriesListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesListCall) Fields(s ...googleapi.Field) *ProjectsHistoriesListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsHistoriesListCall) IfNoneMatch(entityTag string) *ProjectsHistoriesListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesListCall) Context(ctx context.Context) *ProjectsHistoriesListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.list" call.
// Exactly one of *ListHistoriesResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListHistoriesResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsHistoriesListCall) Do(opts ...googleapi.CallOption) (*ListHistoriesResponse, error) {
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
	ret := &ListHistoriesResponse{
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
	//   "description": "Lists Histories for a given Project.\n\nThe histories are sorted by modification time in descending order. The history_id key will be used to order the history with the same modification time.\n\nMay return any of the following canonical error codes:\n\n- PERMISSION_DENIED - if the user is not authorized to read project - INVALID_ARGUMENT - if the request is malformed - NOT_FOUND - if the History does not exist",
	//   "httpMethod": "GET",
	//   "id": "toolresults.projects.histories.list",
	//   "parameterOrder": [
	//     "projectId"
	//   ],
	//   "parameters": {
	//     "filterByName": {
	//       "description": "If set, only return histories with the given name.\n\nOptional.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "The maximum number of Histories to fetch.\n\nDefault value: 20. The server will use this default if the field is not set or has a value of 0. Any value greater than 100 will be treated as 100.\n\nOptional.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "A continuation token to resume the query at the next item.\n\nOptional.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "A Project id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories",
	//   "response": {
	//     "$ref": "ListHistoriesResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *ProjectsHistoriesListCall) Pages(ctx context.Context, f func(*ListHistoriesResponse) error) error {
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

// method id "toolresults.projects.histories.executions.create":

type ProjectsHistoriesExecutionsCreateCall struct {
	s          *Service
	projectId  string
	historyId  string
	execution  *Execution
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Create: Creates an Execution.
//
// The returned Execution will have the id set.
//
// May return any of the following canonical error codes:
//
// - PERMISSION_DENIED - if the user is not authorized to write to
// project - INVALID_ARGUMENT - if the request is malformed - NOT_FOUND
// - if the containing History does not exist
func (r *ProjectsHistoriesExecutionsService) Create(projectId string, historyId string, execution *Execution) *ProjectsHistoriesExecutionsCreateCall {
	c := &ProjectsHistoriesExecutionsCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.execution = execution
	return c
}

// RequestId sets the optional parameter "requestId": A unique request
// ID for server to detect duplicated requests. For example, a
// UUID.
//
// Optional, but strongly recommended.
func (c *ProjectsHistoriesExecutionsCreateCall) RequestId(requestId string) *ProjectsHistoriesExecutionsCreateCall {
	c.urlParams_.Set("requestId", requestId)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsCreateCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsCreateCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.execution)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"historyId": c.historyId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.create" call.
// Exactly one of *Execution or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Execution.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsHistoriesExecutionsCreateCall) Do(opts ...googleapi.CallOption) (*Execution, error) {
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
	ret := &Execution{
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
	//   "description": "Creates an Execution.\n\nThe returned Execution will have the id set.\n\nMay return any of the following canonical error codes:\n\n- PERMISSION_DENIED - if the user is not authorized to write to project - INVALID_ARGUMENT - if the request is malformed - NOT_FOUND - if the containing History does not exist",
	//   "httpMethod": "POST",
	//   "id": "toolresults.projects.histories.executions.create",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId"
	//   ],
	//   "parameters": {
	//     "historyId": {
	//       "description": "A History id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "A Project id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "requestId": {
	//       "description": "A unique request ID for server to detect duplicated requests. For example, a UUID.\n\nOptional, but strongly recommended.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions",
	//   "request": {
	//     "$ref": "Execution"
	//   },
	//   "response": {
	//     "$ref": "Execution"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.histories.executions.get":

type ProjectsHistoriesExecutionsGetCall struct {
	s            *Service
	projectId    string
	historyId    string
	executionId  string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Gets an Execution.
//
// May return any of the following canonical error codes:
//
// - PERMISSION_DENIED - if the user is not authorized to write to
// project - INVALID_ARGUMENT - if the request is malformed - NOT_FOUND
// - if the Execution does not exist
func (r *ProjectsHistoriesExecutionsService) Get(projectId string, historyId string, executionId string) *ProjectsHistoriesExecutionsGetCall {
	c := &ProjectsHistoriesExecutionsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsGetCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsHistoriesExecutionsGetCall) IfNoneMatch(entityTag string) *ProjectsHistoriesExecutionsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsGetCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":   c.projectId,
		"historyId":   c.historyId,
		"executionId": c.executionId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.get" call.
// Exactly one of *Execution or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Execution.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsHistoriesExecutionsGetCall) Do(opts ...googleapi.CallOption) (*Execution, error) {
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
	ret := &Execution{
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
	//   "description": "Gets an Execution.\n\nMay return any of the following canonical error codes:\n\n- PERMISSION_DENIED - if the user is not authorized to write to project - INVALID_ARGUMENT - if the request is malformed - NOT_FOUND - if the Execution does not exist",
	//   "httpMethod": "GET",
	//   "id": "toolresults.projects.histories.executions.get",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId"
	//   ],
	//   "parameters": {
	//     "executionId": {
	//       "description": "An Execution id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "A History id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "A Project id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}",
	//   "response": {
	//     "$ref": "Execution"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.histories.executions.list":

type ProjectsHistoriesExecutionsListCall struct {
	s            *Service
	projectId    string
	historyId    string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists Histories for a given Project.
//
// The executions are sorted by creation_time in descending order. The
// execution_id key will be used to order the executions with the same
// creation_time.
//
// May return any of the following canonical error codes:
//
// - PERMISSION_DENIED - if the user is not authorized to read project -
// INVALID_ARGUMENT - if the request is malformed - NOT_FOUND - if the
// containing History does not exist
func (r *ProjectsHistoriesExecutionsService) List(projectId string, historyId string) *ProjectsHistoriesExecutionsListCall {
	c := &ProjectsHistoriesExecutionsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	return c
}

// PageSize sets the optional parameter "pageSize": The maximum number
// of Executions to fetch.
//
// Default value: 25. The server will use this default if the field is
// not set or has a value of 0.
func (c *ProjectsHistoriesExecutionsListCall) PageSize(pageSize int64) *ProjectsHistoriesExecutionsListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken": A continuation
// token to resume the query at the next item.
func (c *ProjectsHistoriesExecutionsListCall) PageToken(pageToken string) *ProjectsHistoriesExecutionsListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsListCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsHistoriesExecutionsListCall) IfNoneMatch(entityTag string) *ProjectsHistoriesExecutionsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsListCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId": c.projectId,
		"historyId": c.historyId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.list" call.
// Exactly one of *ListExecutionsResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListExecutionsResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsHistoriesExecutionsListCall) Do(opts ...googleapi.CallOption) (*ListExecutionsResponse, error) {
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
	ret := &ListExecutionsResponse{
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
	//   "description": "Lists Histories for a given Project.\n\nThe executions are sorted by creation_time in descending order. The execution_id key will be used to order the executions with the same creation_time.\n\nMay return any of the following canonical error codes:\n\n- PERMISSION_DENIED - if the user is not authorized to read project - INVALID_ARGUMENT - if the request is malformed - NOT_FOUND - if the containing History does not exist",
	//   "httpMethod": "GET",
	//   "id": "toolresults.projects.histories.executions.list",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId"
	//   ],
	//   "parameters": {
	//     "historyId": {
	//       "description": "A History id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "The maximum number of Executions to fetch.\n\nDefault value: 25. The server will use this default if the field is not set or has a value of 0.\n\nOptional.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "A continuation token to resume the query at the next item.\n\nOptional.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "A Project id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions",
	//   "response": {
	//     "$ref": "ListExecutionsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *ProjectsHistoriesExecutionsListCall) Pages(ctx context.Context, f func(*ListExecutionsResponse) error) error {
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

// method id "toolresults.projects.histories.executions.patch":

type ProjectsHistoriesExecutionsPatchCall struct {
	s           *Service
	projectId   string
	historyId   string
	executionId string
	execution   *Execution
	urlParams_  gensupport.URLParams
	ctx_        context.Context
	header_     http.Header
}

// Patch: Updates an existing Execution with the supplied partial
// entity.
//
// May return any of the following canonical error codes:
//
// - PERMISSION_DENIED - if the user is not authorized to write to
// project - INVALID_ARGUMENT - if the request is malformed -
// FAILED_PRECONDITION - if the requested state transition is illegal -
// NOT_FOUND - if the containing History does not exist
func (r *ProjectsHistoriesExecutionsService) Patch(projectId string, historyId string, executionId string, execution *Execution) *ProjectsHistoriesExecutionsPatchCall {
	c := &ProjectsHistoriesExecutionsPatchCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	c.execution = execution
	return c
}

// RequestId sets the optional parameter "requestId": A unique request
// ID for server to detect duplicated requests. For example, a
// UUID.
//
// Optional, but strongly recommended.
func (c *ProjectsHistoriesExecutionsPatchCall) RequestId(requestId string) *ProjectsHistoriesExecutionsPatchCall {
	c.urlParams_.Set("requestId", requestId)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsPatchCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsPatchCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsPatchCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsPatchCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsPatchCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsPatchCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.execution)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":   c.projectId,
		"historyId":   c.historyId,
		"executionId": c.executionId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.patch" call.
// Exactly one of *Execution or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *Execution.ServerResponse.Header or (if a response was returned at
// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
// to check whether the returned error was because
// http.StatusNotModified was returned.
func (c *ProjectsHistoriesExecutionsPatchCall) Do(opts ...googleapi.CallOption) (*Execution, error) {
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
	ret := &Execution{
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
	//   "description": "Updates an existing Execution with the supplied partial entity.\n\nMay return any of the following canonical error codes:\n\n- PERMISSION_DENIED - if the user is not authorized to write to project - INVALID_ARGUMENT - if the request is malformed - FAILED_PRECONDITION - if the requested state transition is illegal - NOT_FOUND - if the containing History does not exist",
	//   "httpMethod": "PATCH",
	//   "id": "toolresults.projects.histories.executions.patch",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId"
	//   ],
	//   "parameters": {
	//     "executionId": {
	//       "description": "Required.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "Required.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "A Project id. Required.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "requestId": {
	//       "description": "A unique request ID for server to detect duplicated requests. For example, a UUID.\n\nOptional, but strongly recommended.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}",
	//   "request": {
	//     "$ref": "Execution"
	//   },
	//   "response": {
	//     "$ref": "Execution"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.histories.executions.clusters.get":

type ProjectsHistoriesExecutionsClustersGetCall struct {
	s            *Service
	projectId    string
	historyId    string
	executionId  string
	clusterId    string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Retrieves a single screenshot cluster by its ID
func (r *ProjectsHistoriesExecutionsClustersService) Get(projectId string, historyId string, executionId string, clusterId string) *ProjectsHistoriesExecutionsClustersGetCall {
	c := &ProjectsHistoriesExecutionsClustersGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	c.clusterId = clusterId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsClustersGetCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsClustersGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsHistoriesExecutionsClustersGetCall) IfNoneMatch(entityTag string) *ProjectsHistoriesExecutionsClustersGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsClustersGetCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsClustersGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsClustersGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsClustersGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}/clusters/{clusterId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":   c.projectId,
		"historyId":   c.historyId,
		"executionId": c.executionId,
		"clusterId":   c.clusterId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.clusters.get" call.
// Exactly one of *ScreenshotCluster or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ScreenshotCluster.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsHistoriesExecutionsClustersGetCall) Do(opts ...googleapi.CallOption) (*ScreenshotCluster, error) {
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
	ret := &ScreenshotCluster{
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
	//   "description": "Retrieves a single screenshot cluster by its ID",
	//   "httpMethod": "GET",
	//   "id": "toolresults.projects.histories.executions.clusters.get",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId",
	//     "clusterId"
	//   ],
	//   "parameters": {
	//     "clusterId": {
	//       "description": "A Cluster id\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "executionId": {
	//       "description": "An Execution id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "A History id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "A Project id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}/clusters/{clusterId}",
	//   "response": {
	//     "$ref": "ScreenshotCluster"
	//   }
	// }

}

// method id "toolresults.projects.histories.executions.clusters.list":

type ProjectsHistoriesExecutionsClustersListCall struct {
	s            *Service
	projectId    string
	historyId    string
	executionId  string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists Screenshot Clusters
//
// Returns the list of screenshot clusters corresponding to an
// execution. Screenshot clusters are created after the execution is
// finished. Clusters are created from a set of screenshots. Between any
// two screenshots, a matching score is calculated based off their
// metadata that determines how similar they are. Screenshots are placed
// in the cluster that has screens which have the highest matching
// scores.
func (r *ProjectsHistoriesExecutionsClustersService) List(projectId string, historyId string, executionId string) *ProjectsHistoriesExecutionsClustersListCall {
	c := &ProjectsHistoriesExecutionsClustersListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsClustersListCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsClustersListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsHistoriesExecutionsClustersListCall) IfNoneMatch(entityTag string) *ProjectsHistoriesExecutionsClustersListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsClustersListCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsClustersListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsClustersListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsClustersListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}/clusters")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":   c.projectId,
		"historyId":   c.historyId,
		"executionId": c.executionId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.clusters.list" call.
// Exactly one of *ListScreenshotClustersResponse or error will be
// non-nil. Any non-2xx status code is an error. Response headers are in
// either *ListScreenshotClustersResponse.ServerResponse.Header or (if a
// response was returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsHistoriesExecutionsClustersListCall) Do(opts ...googleapi.CallOption) (*ListScreenshotClustersResponse, error) {
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
	ret := &ListScreenshotClustersResponse{
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
	//   "description": "Lists Screenshot Clusters\n\nReturns the list of screenshot clusters corresponding to an execution. Screenshot clusters are created after the execution is finished. Clusters are created from a set of screenshots. Between any two screenshots, a matching score is calculated based off their metadata that determines how similar they are. Screenshots are placed in the cluster that has screens which have the highest matching scores.",
	//   "httpMethod": "GET",
	//   "id": "toolresults.projects.histories.executions.clusters.list",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId"
	//   ],
	//   "parameters": {
	//     "executionId": {
	//       "description": "An Execution id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "A History id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "A Project id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}/clusters",
	//   "response": {
	//     "$ref": "ListScreenshotClustersResponse"
	//   }
	// }

}

// method id "toolresults.projects.histories.executions.steps.create":

type ProjectsHistoriesExecutionsStepsCreateCall struct {
	s           *Service
	projectId   string
	historyId   string
	executionId string
	step        *Step
	urlParams_  gensupport.URLParams
	ctx_        context.Context
	header_     http.Header
}

// Create: Creates a Step.
//
// The returned Step will have the id set.
//
// May return any of the following canonical error codes:
//
// - PERMISSION_DENIED - if the user is not authorized to write to
// project - INVALID_ARGUMENT - if the request is malformed -
// FAILED_PRECONDITION - if the step is too large (more than 10Mib) -
// NOT_FOUND - if the containing Execution does not exist
func (r *ProjectsHistoriesExecutionsStepsService) Create(projectId string, historyId string, executionId string, step *Step) *ProjectsHistoriesExecutionsStepsCreateCall {
	c := &ProjectsHistoriesExecutionsStepsCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	c.step = step
	return c
}

// RequestId sets the optional parameter "requestId": A unique request
// ID for server to detect duplicated requests. For example, a
// UUID.
//
// Optional, but strongly recommended.
func (c *ProjectsHistoriesExecutionsStepsCreateCall) RequestId(requestId string) *ProjectsHistoriesExecutionsStepsCreateCall {
	c.urlParams_.Set("requestId", requestId)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsStepsCreateCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsStepsCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsStepsCreateCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsStepsCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsStepsCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsStepsCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.step)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}/steps")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":   c.projectId,
		"historyId":   c.historyId,
		"executionId": c.executionId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.steps.create" call.
// Exactly one of *Step or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Step.ServerResponse.Header or (if a response was returned at all) in
// error.(*googleapi.Error).Header. Use googleapi.IsNotModified to check
// whether the returned error was because http.StatusNotModified was
// returned.
func (c *ProjectsHistoriesExecutionsStepsCreateCall) Do(opts ...googleapi.CallOption) (*Step, error) {
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
	ret := &Step{
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
	//   "description": "Creates a Step.\n\nThe returned Step will have the id set.\n\nMay return any of the following canonical error codes:\n\n- PERMISSION_DENIED - if the user is not authorized to write to project - INVALID_ARGUMENT - if the request is malformed - FAILED_PRECONDITION - if the step is too large (more than 10Mib) - NOT_FOUND - if the containing Execution does not exist",
	//   "httpMethod": "POST",
	//   "id": "toolresults.projects.histories.executions.steps.create",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId"
	//   ],
	//   "parameters": {
	//     "executionId": {
	//       "description": "A Execution id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "A History id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "A Project id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "requestId": {
	//       "description": "A unique request ID for server to detect duplicated requests. For example, a UUID.\n\nOptional, but strongly recommended.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}/steps",
	//   "request": {
	//     "$ref": "Step"
	//   },
	//   "response": {
	//     "$ref": "Step"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.histories.executions.steps.get":

type ProjectsHistoriesExecutionsStepsGetCall struct {
	s            *Service
	projectId    string
	historyId    string
	executionId  string
	stepId       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Gets a Step.
//
// May return any of the following canonical error codes:
//
// - PERMISSION_DENIED - if the user is not authorized to read project -
// INVALID_ARGUMENT - if the request is malformed - NOT_FOUND - if the
// Step does not exist
func (r *ProjectsHistoriesExecutionsStepsService) Get(projectId string, historyId string, executionId string, stepId string) *ProjectsHistoriesExecutionsStepsGetCall {
	c := &ProjectsHistoriesExecutionsStepsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	c.stepId = stepId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsStepsGetCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsStepsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsHistoriesExecutionsStepsGetCall) IfNoneMatch(entityTag string) *ProjectsHistoriesExecutionsStepsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsStepsGetCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsStepsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsStepsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsStepsGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":   c.projectId,
		"historyId":   c.historyId,
		"executionId": c.executionId,
		"stepId":      c.stepId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.steps.get" call.
// Exactly one of *Step or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Step.ServerResponse.Header or (if a response was returned at all) in
// error.(*googleapi.Error).Header. Use googleapi.IsNotModified to check
// whether the returned error was because http.StatusNotModified was
// returned.
func (c *ProjectsHistoriesExecutionsStepsGetCall) Do(opts ...googleapi.CallOption) (*Step, error) {
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
	ret := &Step{
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
	//   "description": "Gets a Step.\n\nMay return any of the following canonical error codes:\n\n- PERMISSION_DENIED - if the user is not authorized to read project - INVALID_ARGUMENT - if the request is malformed - NOT_FOUND - if the Step does not exist",
	//   "httpMethod": "GET",
	//   "id": "toolresults.projects.histories.executions.steps.get",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId",
	//     "stepId"
	//   ],
	//   "parameters": {
	//     "executionId": {
	//       "description": "A Execution id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "A History id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "A Project id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "stepId": {
	//       "description": "A Step id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}",
	//   "response": {
	//     "$ref": "Step"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.histories.executions.steps.getPerfMetricsSummary":

type ProjectsHistoriesExecutionsStepsGetPerfMetricsSummaryCall struct {
	s            *Service
	projectId    string
	historyId    string
	executionId  string
	stepId       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// GetPerfMetricsSummary: Retrieves a PerfMetricsSummary.
//
// May return any of the following error code(s): - NOT_FOUND - The
// specified PerfMetricsSummary does not exist
func (r *ProjectsHistoriesExecutionsStepsService) GetPerfMetricsSummary(projectId string, historyId string, executionId string, stepId string) *ProjectsHistoriesExecutionsStepsGetPerfMetricsSummaryCall {
	c := &ProjectsHistoriesExecutionsStepsGetPerfMetricsSummaryCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	c.stepId = stepId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsStepsGetPerfMetricsSummaryCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsStepsGetPerfMetricsSummaryCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsHistoriesExecutionsStepsGetPerfMetricsSummaryCall) IfNoneMatch(entityTag string) *ProjectsHistoriesExecutionsStepsGetPerfMetricsSummaryCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsStepsGetPerfMetricsSummaryCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsStepsGetPerfMetricsSummaryCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsStepsGetPerfMetricsSummaryCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsStepsGetPerfMetricsSummaryCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}/perfMetricsSummary")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":   c.projectId,
		"historyId":   c.historyId,
		"executionId": c.executionId,
		"stepId":      c.stepId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.steps.getPerfMetricsSummary" call.
// Exactly one of *PerfMetricsSummary or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *PerfMetricsSummary.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsHistoriesExecutionsStepsGetPerfMetricsSummaryCall) Do(opts ...googleapi.CallOption) (*PerfMetricsSummary, error) {
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
	ret := &PerfMetricsSummary{
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
	//   "description": "Retrieves a PerfMetricsSummary.\n\nMay return any of the following error code(s): - NOT_FOUND - The specified PerfMetricsSummary does not exist",
	//   "httpMethod": "GET",
	//   "id": "toolresults.projects.histories.executions.steps.getPerfMetricsSummary",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId",
	//     "stepId"
	//   ],
	//   "parameters": {
	//     "executionId": {
	//       "description": "A tool results execution ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "A tool results history ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The cloud project",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "stepId": {
	//       "description": "A tool results step ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}/perfMetricsSummary",
	//   "response": {
	//     "$ref": "PerfMetricsSummary"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.histories.executions.steps.list":

type ProjectsHistoriesExecutionsStepsListCall struct {
	s            *Service
	projectId    string
	historyId    string
	executionId  string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists Steps for a given Execution.
//
// The steps are sorted by creation_time in descending order. The
// step_id key will be used to order the steps with the same
// creation_time.
//
// May return any of the following canonical error codes:
//
// - PERMISSION_DENIED - if the user is not authorized to read project -
// INVALID_ARGUMENT - if the request is malformed - FAILED_PRECONDITION
// - if an argument in the request happens to be invalid; e.g. if an
// attempt is made to list the children of a nonexistent Step -
// NOT_FOUND - if the containing Execution does not exist
func (r *ProjectsHistoriesExecutionsStepsService) List(projectId string, historyId string, executionId string) *ProjectsHistoriesExecutionsStepsListCall {
	c := &ProjectsHistoriesExecutionsStepsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	return c
}

// PageSize sets the optional parameter "pageSize": The maximum number
// of Steps to fetch.
//
// Default value: 25. The server will use this default if the field is
// not set or has a value of 0.
func (c *ProjectsHistoriesExecutionsStepsListCall) PageSize(pageSize int64) *ProjectsHistoriesExecutionsStepsListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken": A continuation
// token to resume the query at the next item.
func (c *ProjectsHistoriesExecutionsStepsListCall) PageToken(pageToken string) *ProjectsHistoriesExecutionsStepsListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsStepsListCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsStepsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsHistoriesExecutionsStepsListCall) IfNoneMatch(entityTag string) *ProjectsHistoriesExecutionsStepsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsStepsListCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsStepsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsStepsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsStepsListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}/steps")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":   c.projectId,
		"historyId":   c.historyId,
		"executionId": c.executionId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.steps.list" call.
// Exactly one of *ListStepsResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListStepsResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsHistoriesExecutionsStepsListCall) Do(opts ...googleapi.CallOption) (*ListStepsResponse, error) {
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
	ret := &ListStepsResponse{
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
	//   "description": "Lists Steps for a given Execution.\n\nThe steps are sorted by creation_time in descending order. The step_id key will be used to order the steps with the same creation_time.\n\nMay return any of the following canonical error codes:\n\n- PERMISSION_DENIED - if the user is not authorized to read project - INVALID_ARGUMENT - if the request is malformed - FAILED_PRECONDITION - if an argument in the request happens to be invalid; e.g. if an attempt is made to list the children of a nonexistent Step - NOT_FOUND - if the containing Execution does not exist",
	//   "httpMethod": "GET",
	//   "id": "toolresults.projects.histories.executions.steps.list",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId"
	//   ],
	//   "parameters": {
	//     "executionId": {
	//       "description": "A Execution id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "A History id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "The maximum number of Steps to fetch.\n\nDefault value: 25. The server will use this default if the field is not set or has a value of 0.\n\nOptional.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "A continuation token to resume the query at the next item.\n\nOptional.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "A Project id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}/steps",
	//   "response": {
	//     "$ref": "ListStepsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *ProjectsHistoriesExecutionsStepsListCall) Pages(ctx context.Context, f func(*ListStepsResponse) error) error {
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

// method id "toolresults.projects.histories.executions.steps.patch":

type ProjectsHistoriesExecutionsStepsPatchCall struct {
	s           *Service
	projectId   string
	historyId   string
	executionId string
	stepId      string
	step        *Step
	urlParams_  gensupport.URLParams
	ctx_        context.Context
	header_     http.Header
}

// Patch: Updates an existing Step with the supplied partial
// entity.
//
// May return any of the following canonical error codes:
//
// - PERMISSION_DENIED - if the user is not authorized to write project
// - INVALID_ARGUMENT - if the request is malformed -
// FAILED_PRECONDITION - if the requested state transition is illegal
// (e.g try to upload a duplicate xml file), if the updated step is too
// large (more than 10Mib) - NOT_FOUND - if the containing Execution
// does not exist
func (r *ProjectsHistoriesExecutionsStepsService) Patch(projectId string, historyId string, executionId string, stepId string, step *Step) *ProjectsHistoriesExecutionsStepsPatchCall {
	c := &ProjectsHistoriesExecutionsStepsPatchCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	c.stepId = stepId
	c.step = step
	return c
}

// RequestId sets the optional parameter "requestId": A unique request
// ID for server to detect duplicated requests. For example, a
// UUID.
//
// Optional, but strongly recommended.
func (c *ProjectsHistoriesExecutionsStepsPatchCall) RequestId(requestId string) *ProjectsHistoriesExecutionsStepsPatchCall {
	c.urlParams_.Set("requestId", requestId)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsStepsPatchCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsStepsPatchCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsStepsPatchCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsStepsPatchCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsStepsPatchCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsStepsPatchCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.step)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":   c.projectId,
		"historyId":   c.historyId,
		"executionId": c.executionId,
		"stepId":      c.stepId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.steps.patch" call.
// Exactly one of *Step or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Step.ServerResponse.Header or (if a response was returned at all) in
// error.(*googleapi.Error).Header. Use googleapi.IsNotModified to check
// whether the returned error was because http.StatusNotModified was
// returned.
func (c *ProjectsHistoriesExecutionsStepsPatchCall) Do(opts ...googleapi.CallOption) (*Step, error) {
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
	ret := &Step{
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
	//   "description": "Updates an existing Step with the supplied partial entity.\n\nMay return any of the following canonical error codes:\n\n- PERMISSION_DENIED - if the user is not authorized to write project - INVALID_ARGUMENT - if the request is malformed - FAILED_PRECONDITION - if the requested state transition is illegal (e.g try to upload a duplicate xml file), if the updated step is too large (more than 10Mib) - NOT_FOUND - if the containing Execution does not exist",
	//   "httpMethod": "PATCH",
	//   "id": "toolresults.projects.histories.executions.steps.patch",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId",
	//     "stepId"
	//   ],
	//   "parameters": {
	//     "executionId": {
	//       "description": "A Execution id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "A History id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "A Project id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "requestId": {
	//       "description": "A unique request ID for server to detect duplicated requests. For example, a UUID.\n\nOptional, but strongly recommended.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "stepId": {
	//       "description": "A Step id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}",
	//   "request": {
	//     "$ref": "Step"
	//   },
	//   "response": {
	//     "$ref": "Step"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.histories.executions.steps.publishXunitXmlFiles":

type ProjectsHistoriesExecutionsStepsPublishXunitXmlFilesCall struct {
	s                           *Service
	projectId                   string
	historyId                   string
	executionId                 string
	stepId                      string
	publishxunitxmlfilesrequest *PublishXunitXmlFilesRequest
	urlParams_                  gensupport.URLParams
	ctx_                        context.Context
	header_                     http.Header
}

// PublishXunitXmlFiles: Publish xml files to an existing Step.
//
// May return any of the following canonical error codes:
//
// - PERMISSION_DENIED - if the user is not authorized to write project
// - INVALID_ARGUMENT - if the request is malformed -
// FAILED_PRECONDITION - if the requested state transition is illegal,
// e.g try to upload a duplicate xml file or a file too large. -
// NOT_FOUND - if the containing Execution does not exist
func (r *ProjectsHistoriesExecutionsStepsService) PublishXunitXmlFiles(projectId string, historyId string, executionId string, stepId string, publishxunitxmlfilesrequest *PublishXunitXmlFilesRequest) *ProjectsHistoriesExecutionsStepsPublishXunitXmlFilesCall {
	c := &ProjectsHistoriesExecutionsStepsPublishXunitXmlFilesCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	c.stepId = stepId
	c.publishxunitxmlfilesrequest = publishxunitxmlfilesrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsStepsPublishXunitXmlFilesCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsStepsPublishXunitXmlFilesCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsStepsPublishXunitXmlFilesCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsStepsPublishXunitXmlFilesCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsStepsPublishXunitXmlFilesCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsStepsPublishXunitXmlFilesCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.publishxunitxmlfilesrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}:publishXunitXmlFiles")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":   c.projectId,
		"historyId":   c.historyId,
		"executionId": c.executionId,
		"stepId":      c.stepId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.steps.publishXunitXmlFiles" call.
// Exactly one of *Step or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Step.ServerResponse.Header or (if a response was returned at all) in
// error.(*googleapi.Error).Header. Use googleapi.IsNotModified to check
// whether the returned error was because http.StatusNotModified was
// returned.
func (c *ProjectsHistoriesExecutionsStepsPublishXunitXmlFilesCall) Do(opts ...googleapi.CallOption) (*Step, error) {
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
	ret := &Step{
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
	//   "description": "Publish xml files to an existing Step.\n\nMay return any of the following canonical error codes:\n\n- PERMISSION_DENIED - if the user is not authorized to write project - INVALID_ARGUMENT - if the request is malformed - FAILED_PRECONDITION - if the requested state transition is illegal, e.g try to upload a duplicate xml file or a file too large. - NOT_FOUND - if the containing Execution does not exist",
	//   "httpMethod": "POST",
	//   "id": "toolresults.projects.histories.executions.steps.publishXunitXmlFiles",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId",
	//     "stepId"
	//   ],
	//   "parameters": {
	//     "executionId": {
	//       "description": "A Execution id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "A History id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "A Project id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "stepId": {
	//       "description": "A Step id. Note: This step must include a TestExecutionStep.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}:publishXunitXmlFiles",
	//   "request": {
	//     "$ref": "PublishXunitXmlFilesRequest"
	//   },
	//   "response": {
	//     "$ref": "Step"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.histories.executions.steps.perfMetricsSummary.create":

type ProjectsHistoriesExecutionsStepsPerfMetricsSummaryCreateCall struct {
	s                  *Service
	projectId          string
	historyId          string
	executionId        string
	stepId             string
	perfmetricssummary *PerfMetricsSummary
	urlParams_         gensupport.URLParams
	ctx_               context.Context
	header_            http.Header
}

// Create: Creates a PerfMetricsSummary resource. Returns the existing
// one if it has already been created.
//
// May return any of the following error code(s): - NOT_FOUND - The
// containing Step does not exist
func (r *ProjectsHistoriesExecutionsStepsPerfMetricsSummaryService) Create(projectId string, historyId string, executionId string, stepId string, perfmetricssummary *PerfMetricsSummary) *ProjectsHistoriesExecutionsStepsPerfMetricsSummaryCreateCall {
	c := &ProjectsHistoriesExecutionsStepsPerfMetricsSummaryCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	c.stepId = stepId
	c.perfmetricssummary = perfmetricssummary
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsStepsPerfMetricsSummaryCreateCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsStepsPerfMetricsSummaryCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsStepsPerfMetricsSummaryCreateCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsStepsPerfMetricsSummaryCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsStepsPerfMetricsSummaryCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsStepsPerfMetricsSummaryCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.perfmetricssummary)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}/perfMetricsSummary")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":   c.projectId,
		"historyId":   c.historyId,
		"executionId": c.executionId,
		"stepId":      c.stepId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.steps.perfMetricsSummary.create" call.
// Exactly one of *PerfMetricsSummary or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *PerfMetricsSummary.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsHistoriesExecutionsStepsPerfMetricsSummaryCreateCall) Do(opts ...googleapi.CallOption) (*PerfMetricsSummary, error) {
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
	ret := &PerfMetricsSummary{
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
	//   "description": "Creates a PerfMetricsSummary resource. Returns the existing one if it has already been created.\n\nMay return any of the following error code(s): - NOT_FOUND - The containing Step does not exist",
	//   "httpMethod": "POST",
	//   "id": "toolresults.projects.histories.executions.steps.perfMetricsSummary.create",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId",
	//     "stepId"
	//   ],
	//   "parameters": {
	//     "executionId": {
	//       "description": "A tool results execution ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "A tool results history ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The cloud project",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "stepId": {
	//       "description": "A tool results step ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}/perfMetricsSummary",
	//   "request": {
	//     "$ref": "PerfMetricsSummary"
	//   },
	//   "response": {
	//     "$ref": "PerfMetricsSummary"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.histories.executions.steps.perfSampleSeries.create":

type ProjectsHistoriesExecutionsStepsPerfSampleSeriesCreateCall struct {
	s                *Service
	projectId        string
	historyId        string
	executionId      string
	stepId           string
	perfsampleseries *PerfSampleSeries
	urlParams_       gensupport.URLParams
	ctx_             context.Context
	header_          http.Header
}

// Create: Creates a PerfSampleSeries.
//
// May return any of the following error code(s): - ALREADY_EXISTS -
// PerfMetricSummary already exists for the given Step - NOT_FOUND - The
// containing Step does not exist
func (r *ProjectsHistoriesExecutionsStepsPerfSampleSeriesService) Create(projectId string, historyId string, executionId string, stepId string, perfsampleseries *PerfSampleSeries) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesCreateCall {
	c := &ProjectsHistoriesExecutionsStepsPerfSampleSeriesCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	c.stepId = stepId
	c.perfsampleseries = perfsampleseries
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesCreateCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesCreateCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.perfsampleseries)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}/perfSampleSeries")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":   c.projectId,
		"historyId":   c.historyId,
		"executionId": c.executionId,
		"stepId":      c.stepId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.steps.perfSampleSeries.create" call.
// Exactly one of *PerfSampleSeries or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *PerfSampleSeries.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesCreateCall) Do(opts ...googleapi.CallOption) (*PerfSampleSeries, error) {
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
	ret := &PerfSampleSeries{
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
	//   "description": "Creates a PerfSampleSeries.\n\nMay return any of the following error code(s): - ALREADY_EXISTS - PerfMetricSummary already exists for the given Step - NOT_FOUND - The containing Step does not exist",
	//   "httpMethod": "POST",
	//   "id": "toolresults.projects.histories.executions.steps.perfSampleSeries.create",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId",
	//     "stepId"
	//   ],
	//   "parameters": {
	//     "executionId": {
	//       "description": "A tool results execution ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "A tool results history ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The cloud project",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "stepId": {
	//       "description": "A tool results step ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}/perfSampleSeries",
	//   "request": {
	//     "$ref": "PerfSampleSeries"
	//   },
	//   "response": {
	//     "$ref": "PerfSampleSeries"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.histories.executions.steps.perfSampleSeries.get":

type ProjectsHistoriesExecutionsStepsPerfSampleSeriesGetCall struct {
	s              *Service
	projectId      string
	historyId      string
	executionId    string
	stepId         string
	sampleSeriesId string
	urlParams_     gensupport.URLParams
	ifNoneMatch_   string
	ctx_           context.Context
	header_        http.Header
}

// Get: Gets a PerfSampleSeries.
//
// May return any of the following error code(s): - NOT_FOUND - The
// specified PerfSampleSeries does not exist
func (r *ProjectsHistoriesExecutionsStepsPerfSampleSeriesService) Get(projectId string, historyId string, executionId string, stepId string, sampleSeriesId string) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesGetCall {
	c := &ProjectsHistoriesExecutionsStepsPerfSampleSeriesGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	c.stepId = stepId
	c.sampleSeriesId = sampleSeriesId
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesGetCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesGetCall) IfNoneMatch(entityTag string) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesGetCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesGetCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}/perfSampleSeries/{sampleSeriesId}")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":      c.projectId,
		"historyId":      c.historyId,
		"executionId":    c.executionId,
		"stepId":         c.stepId,
		"sampleSeriesId": c.sampleSeriesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.steps.perfSampleSeries.get" call.
// Exactly one of *PerfSampleSeries or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *PerfSampleSeries.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesGetCall) Do(opts ...googleapi.CallOption) (*PerfSampleSeries, error) {
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
	ret := &PerfSampleSeries{
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
	//   "description": "Gets a PerfSampleSeries.\n\nMay return any of the following error code(s): - NOT_FOUND - The specified PerfSampleSeries does not exist",
	//   "httpMethod": "GET",
	//   "id": "toolresults.projects.histories.executions.steps.perfSampleSeries.get",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId",
	//     "stepId",
	//     "sampleSeriesId"
	//   ],
	//   "parameters": {
	//     "executionId": {
	//       "description": "A tool results execution ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "A tool results history ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The cloud project",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "sampleSeriesId": {
	//       "description": "A sample series id",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "stepId": {
	//       "description": "A tool results step ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}/perfSampleSeries/{sampleSeriesId}",
	//   "response": {
	//     "$ref": "PerfSampleSeries"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.histories.executions.steps.perfSampleSeries.list":

type ProjectsHistoriesExecutionsStepsPerfSampleSeriesListCall struct {
	s            *Service
	projectId    string
	historyId    string
	executionId  string
	stepId       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists PerfSampleSeries for a given Step.
//
// The request provides an optional filter which specifies one or more
// PerfMetricsType to include in the result; if none returns all. The
// resulting PerfSampleSeries are sorted by ids.
//
// May return any of the following canonical error codes: - NOT_FOUND -
// The containing Step does not exist
func (r *ProjectsHistoriesExecutionsStepsPerfSampleSeriesService) List(projectId string, historyId string, executionId string, stepId string) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesListCall {
	c := &ProjectsHistoriesExecutionsStepsPerfSampleSeriesListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	c.stepId = stepId
	return c
}

// Filter sets the optional parameter "filter": Specify one or more
// PerfMetricType values such as CPU to filter the result
//
// Possible values:
//   "cpu"
//   "graphics"
//   "memory"
//   "network"
//   "perfMetricTypeUnspecified"
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesListCall) Filter(filter ...string) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesListCall {
	c.urlParams_.SetMulti("filter", append([]string{}, filter...))
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesListCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesListCall) IfNoneMatch(entityTag string) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesListCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}/perfSampleSeries")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":   c.projectId,
		"historyId":   c.historyId,
		"executionId": c.executionId,
		"stepId":      c.stepId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.steps.perfSampleSeries.list" call.
// Exactly one of *ListPerfSampleSeriesResponse or error will be
// non-nil. Any non-2xx status code is an error. Response headers are in
// either *ListPerfSampleSeriesResponse.ServerResponse.Header or (if a
// response was returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesListCall) Do(opts ...googleapi.CallOption) (*ListPerfSampleSeriesResponse, error) {
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
	ret := &ListPerfSampleSeriesResponse{
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
	//   "description": "Lists PerfSampleSeries for a given Step.\n\nThe request provides an optional filter which specifies one or more PerfMetricsType to include in the result; if none returns all. The resulting PerfSampleSeries are sorted by ids.\n\nMay return any of the following canonical error codes: - NOT_FOUND - The containing Step does not exist",
	//   "httpMethod": "GET",
	//   "id": "toolresults.projects.histories.executions.steps.perfSampleSeries.list",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId",
	//     "stepId"
	//   ],
	//   "parameters": {
	//     "executionId": {
	//       "description": "A tool results execution ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "filter": {
	//       "description": "Specify one or more PerfMetricType values such as CPU to filter the result",
	//       "enum": [
	//         "cpu",
	//         "graphics",
	//         "memory",
	//         "network",
	//         "perfMetricTypeUnspecified"
	//       ],
	//       "enumDescriptions": [
	//         "",
	//         "",
	//         "",
	//         "",
	//         ""
	//       ],
	//       "location": "query",
	//       "repeated": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "A tool results history ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The cloud project",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "stepId": {
	//       "description": "A tool results step ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}/perfSampleSeries",
	//   "response": {
	//     "$ref": "ListPerfSampleSeriesResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.histories.executions.steps.perfSampleSeries.samples.batchCreate":

type ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesBatchCreateCall struct {
	s                             *Service
	projectId                     string
	historyId                     string
	executionId                   string
	stepId                        string
	sampleSeriesId                string
	batchcreateperfsamplesrequest *BatchCreatePerfSamplesRequest
	urlParams_                    gensupport.URLParams
	ctx_                          context.Context
	header_                       http.Header
}

// BatchCreate: Creates a batch of PerfSamples - a client can submit
// multiple batches of Perf Samples through repeated calls to this
// method in order to split up a large request payload - duplicates and
// existing timestamp entries will be ignored. - the batch operation may
// partially succeed - the set of elements successfully inserted is
// returned in the response (omits items which already existed in the
// database).
//
// May return any of the following canonical error codes: - NOT_FOUND -
// The containing PerfSampleSeries does not exist
func (r *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesService) BatchCreate(projectId string, historyId string, executionId string, stepId string, sampleSeriesId string, batchcreateperfsamplesrequest *BatchCreatePerfSamplesRequest) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesBatchCreateCall {
	c := &ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesBatchCreateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	c.stepId = stepId
	c.sampleSeriesId = sampleSeriesId
	c.batchcreateperfsamplesrequest = batchcreateperfsamplesrequest
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesBatchCreateCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesBatchCreateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesBatchCreateCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesBatchCreateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesBatchCreateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesBatchCreateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.batchcreateperfsamplesrequest)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}/perfSampleSeries/{sampleSeriesId}/samples:batchCreate")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":      c.projectId,
		"historyId":      c.historyId,
		"executionId":    c.executionId,
		"stepId":         c.stepId,
		"sampleSeriesId": c.sampleSeriesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.steps.perfSampleSeries.samples.batchCreate" call.
// Exactly one of *BatchCreatePerfSamplesResponse or error will be
// non-nil. Any non-2xx status code is an error. Response headers are in
// either *BatchCreatePerfSamplesResponse.ServerResponse.Header or (if a
// response was returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesBatchCreateCall) Do(opts ...googleapi.CallOption) (*BatchCreatePerfSamplesResponse, error) {
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
	ret := &BatchCreatePerfSamplesResponse{
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
	//   "description": "Creates a batch of PerfSamples - a client can submit multiple batches of Perf Samples through repeated calls to this method in order to split up a large request payload - duplicates and existing timestamp entries will be ignored. - the batch operation may partially succeed - the set of elements successfully inserted is returned in the response (omits items which already existed in the database).\n\nMay return any of the following canonical error codes: - NOT_FOUND - The containing PerfSampleSeries does not exist",
	//   "httpMethod": "POST",
	//   "id": "toolresults.projects.histories.executions.steps.perfSampleSeries.samples.batchCreate",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId",
	//     "stepId",
	//     "sampleSeriesId"
	//   ],
	//   "parameters": {
	//     "executionId": {
	//       "description": "A tool results execution ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "A tool results history ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The cloud project",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "sampleSeriesId": {
	//       "description": "A sample series id",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "stepId": {
	//       "description": "A tool results step ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}/perfSampleSeries/{sampleSeriesId}/samples:batchCreate",
	//   "request": {
	//     "$ref": "BatchCreatePerfSamplesRequest"
	//   },
	//   "response": {
	//     "$ref": "BatchCreatePerfSamplesResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// method id "toolresults.projects.histories.executions.steps.perfSampleSeries.samples.list":

type ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall struct {
	s              *Service
	projectId      string
	historyId      string
	executionId    string
	stepId         string
	sampleSeriesId string
	urlParams_     gensupport.URLParams
	ifNoneMatch_   string
	ctx_           context.Context
	header_        http.Header
}

// List: Lists the Performance Samples of a given Sample Series - The
// list results are sorted by timestamps ascending - The default page
// size is 500 samples; and maximum size allowed 5000 - The response
// token indicates the last returned PerfSample timestamp - When the
// results size exceeds the page size, submit a subsequent request
// including the page token to return the rest of the samples up to the
// page limit
//
// May return any of the following canonical error codes: - OUT_OF_RANGE
// - The specified request page_token is out of valid range - NOT_FOUND
// - The containing PerfSampleSeries does not exist
func (r *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesService) List(projectId string, historyId string, executionId string, stepId string, sampleSeriesId string) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall {
	c := &ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	c.stepId = stepId
	c.sampleSeriesId = sampleSeriesId
	return c
}

// PageSize sets the optional parameter "pageSize": The default page
// size is 500 samples, and the maximum size is 5000. If the page_size
// is greater than 5000, the effective page size will be 5000
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall) PageSize(pageSize int64) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken": Optional, the
// next_page_token returned in the previous response
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall) PageToken(pageToken string) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall) IfNoneMatch(entityTag string) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}/perfSampleSeries/{sampleSeriesId}/samples")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":      c.projectId,
		"historyId":      c.historyId,
		"executionId":    c.executionId,
		"stepId":         c.stepId,
		"sampleSeriesId": c.sampleSeriesId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.steps.perfSampleSeries.samples.list" call.
// Exactly one of *ListPerfSamplesResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListPerfSamplesResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall) Do(opts ...googleapi.CallOption) (*ListPerfSamplesResponse, error) {
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
	ret := &ListPerfSamplesResponse{
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
	//   "description": "Lists the Performance Samples of a given Sample Series - The list results are sorted by timestamps ascending - The default page size is 500 samples; and maximum size allowed 5000 - The response token indicates the last returned PerfSample timestamp - When the results size exceeds the page size, submit a subsequent request including the page token to return the rest of the samples up to the page limit\n\nMay return any of the following canonical error codes: - OUT_OF_RANGE - The specified request page_token is out of valid range - NOT_FOUND - The containing PerfSampleSeries does not exist",
	//   "httpMethod": "GET",
	//   "id": "toolresults.projects.histories.executions.steps.perfSampleSeries.samples.list",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId",
	//     "stepId",
	//     "sampleSeriesId"
	//   ],
	//   "parameters": {
	//     "executionId": {
	//       "description": "A tool results execution ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "A tool results history ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "The default page size is 500 samples, and the maximum size is 5000. If the page_size is greater than 5000, the effective page size will be 5000",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "Optional, the next_page_token returned in the previous response",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "The cloud project",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "sampleSeriesId": {
	//       "description": "A sample series id",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "stepId": {
	//       "description": "A tool results step ID.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}/perfSampleSeries/{sampleSeriesId}/samples",
	//   "response": {
	//     "$ref": "ListPerfSamplesResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *ProjectsHistoriesExecutionsStepsPerfSampleSeriesSamplesListCall) Pages(ctx context.Context, f func(*ListPerfSamplesResponse) error) error {
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

// method id "toolresults.projects.histories.executions.steps.thumbnails.list":

type ProjectsHistoriesExecutionsStepsThumbnailsListCall struct {
	s            *Service
	projectId    string
	historyId    string
	executionId  string
	stepId       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists thumbnails of images attached to a step.
//
// May return any of the following canonical error codes: -
// PERMISSION_DENIED - if the user is not authorized to read from the
// project, or from any of the images - INVALID_ARGUMENT - if the
// request is malformed - NOT_FOUND - if the step does not exist, or if
// any of the images do not exist
func (r *ProjectsHistoriesExecutionsStepsThumbnailsService) List(projectId string, historyId string, executionId string, stepId string) *ProjectsHistoriesExecutionsStepsThumbnailsListCall {
	c := &ProjectsHistoriesExecutionsStepsThumbnailsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.projectId = projectId
	c.historyId = historyId
	c.executionId = executionId
	c.stepId = stepId
	return c
}

// PageSize sets the optional parameter "pageSize": The maximum number
// of thumbnails to fetch.
//
// Default value: 50. The server will use this default if the field is
// not set or has a value of 0.
func (c *ProjectsHistoriesExecutionsStepsThumbnailsListCall) PageSize(pageSize int64) *ProjectsHistoriesExecutionsStepsThumbnailsListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken": A continuation
// token to resume the query at the next item.
func (c *ProjectsHistoriesExecutionsStepsThumbnailsListCall) PageToken(pageToken string) *ProjectsHistoriesExecutionsStepsThumbnailsListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *ProjectsHistoriesExecutionsStepsThumbnailsListCall) Fields(s ...googleapi.Field) *ProjectsHistoriesExecutionsStepsThumbnailsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *ProjectsHistoriesExecutionsStepsThumbnailsListCall) IfNoneMatch(entityTag string) *ProjectsHistoriesExecutionsStepsThumbnailsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *ProjectsHistoriesExecutionsStepsThumbnailsListCall) Context(ctx context.Context) *ProjectsHistoriesExecutionsStepsThumbnailsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *ProjectsHistoriesExecutionsStepsThumbnailsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *ProjectsHistoriesExecutionsStepsThumbnailsListCall) doRequest(alt string) (*http.Response, error) {
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
	urls := googleapi.ResolveRelative(c.s.BasePath, "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}/thumbnails")
	urls += "?" + c.urlParams_.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"projectId":   c.projectId,
		"historyId":   c.historyId,
		"executionId": c.executionId,
		"stepId":      c.stepId,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "toolresults.projects.histories.executions.steps.thumbnails.list" call.
// Exactly one of *ListStepThumbnailsResponse or error will be non-nil.
// Any non-2xx status code is an error. Response headers are in either
// *ListStepThumbnailsResponse.ServerResponse.Header or (if a response
// was returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *ProjectsHistoriesExecutionsStepsThumbnailsListCall) Do(opts ...googleapi.CallOption) (*ListStepThumbnailsResponse, error) {
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
	ret := &ListStepThumbnailsResponse{
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
	//   "description": "Lists thumbnails of images attached to a step.\n\nMay return any of the following canonical error codes: - PERMISSION_DENIED - if the user is not authorized to read from the project, or from any of the images - INVALID_ARGUMENT - if the request is malformed - NOT_FOUND - if the step does not exist, or if any of the images do not exist",
	//   "httpMethod": "GET",
	//   "id": "toolresults.projects.histories.executions.steps.thumbnails.list",
	//   "parameterOrder": [
	//     "projectId",
	//     "historyId",
	//     "executionId",
	//     "stepId"
	//   ],
	//   "parameters": {
	//     "executionId": {
	//       "description": "An Execution id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "historyId": {
	//       "description": "A History id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "The maximum number of thumbnails to fetch.\n\nDefault value: 50. The server will use this default if the field is not set or has a value of 0.\n\nOptional.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "A continuation token to resume the query at the next item.\n\nOptional.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "projectId": {
	//       "description": "A Project id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "stepId": {
	//       "description": "A Step id.\n\nRequired.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "{projectId}/histories/{historyId}/executions/{executionId}/steps/{stepId}/thumbnails",
	//   "response": {
	//     "$ref": "ListStepThumbnailsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/cloud-platform"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *ProjectsHistoriesExecutionsStepsThumbnailsListCall) Pages(ctx context.Context, f func(*ListStepThumbnailsResponse) error) error {
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
