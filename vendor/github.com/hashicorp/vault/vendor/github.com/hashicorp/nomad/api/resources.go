package api

import (
	"strconv"
)

// Resources encapsulates the required resources of
// a given task or task group.
type Resources struct {
	CPU      *int
	MemoryMB *int `mapstructure:"memory"`
	DiskMB   *int `mapstructure:"disk"`
	Networks []*NetworkResource
	Devices  []*RequestedDevice

	// COMPAT(0.10)
	// XXX Deprecated. Please do not use. The field will be removed in Nomad
	// 0.10 and is only being kept to allow any references to be removed before
	// then.
	IOPS *int
}

// Canonicalize will supply missing values in the cases
// where they are not provided.
func (r *Resources) Canonicalize() {
	defaultResources := DefaultResources()
	if r.CPU == nil {
		r.CPU = defaultResources.CPU
	}
	if r.MemoryMB == nil {
		r.MemoryMB = defaultResources.MemoryMB
	}
	for _, n := range r.Networks {
		n.Canonicalize()
	}
	for _, d := range r.Devices {
		d.Canonicalize()
	}
}

// DefaultResources is a small resources object that contains the
// default resources requests that we will provide to an object.
// ---  THIS FUNCTION IS REPLICATED IN nomad/structs/structs.go
// and should be kept in sync.
func DefaultResources() *Resources {
	return &Resources{
		CPU:      intToPtr(100),
		MemoryMB: intToPtr(300),
	}
}

// MinResources is a small resources object that contains the
// absolute minimum resources that we will provide to an object.
// This should not be confused with the defaults which are
// provided in DefaultResources() ---  THIS LOGIC IS REPLICATED
// IN nomad/structs/structs.go and should be kept in sync.
func MinResources() *Resources {
	return &Resources{
		CPU:      intToPtr(20),
		MemoryMB: intToPtr(10),
	}
}

// Merge merges this resource with another resource.
func (r *Resources) Merge(other *Resources) {
	if other == nil {
		return
	}
	if other.CPU != nil {
		r.CPU = other.CPU
	}
	if other.MemoryMB != nil {
		r.MemoryMB = other.MemoryMB
	}
	if other.DiskMB != nil {
		r.DiskMB = other.DiskMB
	}
	if len(other.Networks) != 0 {
		r.Networks = other.Networks
	}
	if len(other.Devices) != 0 {
		r.Devices = other.Devices
	}
}

type Port struct {
	Label string
	Value int `mapstructure:"static"`
}

// NetworkResource is used to describe required network
// resources of a given task.
type NetworkResource struct {
	Device        string
	CIDR          string
	IP            string
	MBits         *int
	ReservedPorts []Port
	DynamicPorts  []Port
}

func (n *NetworkResource) Canonicalize() {
	if n.MBits == nil {
		n.MBits = intToPtr(10)
	}
}

// NodeDeviceResource captures a set of devices sharing a common
// vendor/type/device_name tuple.
type NodeDeviceResource struct {

	// Vendor specifies the vendor of device
	Vendor string

	// Type specifies the type of the device
	Type string

	// Name specifies the specific model of the device
	Name string

	// Instances are list of the devices matching the vendor/type/name
	Instances []*NodeDevice

	Attributes map[string]*Attribute
}

func (r NodeDeviceResource) ID() string {
	return r.Vendor + "/" + r.Type + "/" + r.Name
}

// NodeDevice is an instance of a particular device.
type NodeDevice struct {
	// ID is the ID of the device.
	ID string

	// Healthy captures whether the device is healthy.
	Healthy bool

	// HealthDescription is used to provide a human readable description of why
	// the device may be unhealthy.
	HealthDescription string

	// Locality stores HW locality information for the node to optionally be
	// used when making placement decisions.
	Locality *NodeDeviceLocality
}

// Attribute is used to describe the value of an attribute, optionally
// specifying units
type Attribute struct {
	// Float is the float value for the attribute
	FloatVal *float64 `json:"Float,omitempty"`

	// Int is the int value for the attribute
	IntVal *int64 `json:"Int,omitempty"`

	// String is the string value for the attribute
	StringVal *string `json:"String,omitempty"`

	// Bool is the bool value for the attribute
	BoolVal *bool `json:"Bool,omitempty"`

	// Unit is the optional unit for the set int or float value
	Unit string
}

func (a Attribute) String() string {
	switch {
	case a.FloatVal != nil:
		str := formatFloat(*a.FloatVal, 3)
		if a.Unit != "" {
			str += " " + a.Unit
		}
		return str
	case a.IntVal != nil:
		str := strconv.FormatInt(*a.IntVal, 10)
		if a.Unit != "" {
			str += " " + a.Unit
		}
		return str
	case a.StringVal != nil:
		return *a.StringVal
	case a.BoolVal != nil:
		return strconv.FormatBool(*a.BoolVal)
	default:
		return "<unknown>"
	}
}

// NodeDeviceLocality stores information about the devices hardware locality on
// the node.
type NodeDeviceLocality struct {
	// PciBusID is the PCI Bus ID for the device.
	PciBusID string
}

// RequestedDevice is used to request a device for a task.
type RequestedDevice struct {
	// Name is the request name. The possible values are as follows:
	// * <type>: A single value only specifies the type of request.
	// * <vendor>/<type>: A single slash delimiter assumes the vendor and type of device is specified.
	// * <vendor>/<type>/<name>: Two slash delimiters assume vendor, type and specific model are specified.
	//
	// Examples are as follows:
	// * "gpu"
	// * "nvidia/gpu"
	// * "nvidia/gpu/GTX2080Ti"
	Name string

	// Count is the number of requested devices
	Count *uint64

	// Constraints are a set of constraints to apply when selecting the device
	// to use.
	Constraints []*Constraint

	// Affinities are a set of affinites to apply when selecting the device
	// to use.
	Affinities []*Affinity
}

func (d *RequestedDevice) Canonicalize() {
	if d.Count == nil {
		d.Count = uint64ToPtr(1)
	}

	for _, a := range d.Affinities {
		a.Canonicalize()
	}
}
