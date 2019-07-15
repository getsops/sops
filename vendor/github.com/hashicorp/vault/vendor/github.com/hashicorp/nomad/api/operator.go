package api

import "strconv"

// Operator can be used to perform low-level operator tasks for Nomad.
type Operator struct {
	c *Client
}

// Operator returns a handle to the operator endpoints.
func (c *Client) Operator() *Operator {
	return &Operator{c}
}

// RaftServer has information about a server in the Raft configuration.
type RaftServer struct {
	// ID is the unique ID for the server. These are currently the same
	// as the address, but they will be changed to a real GUID in a future
	// release of Nomad.
	ID string

	// Node is the node name of the server, as known by Nomad, or this
	// will be set to "(unknown)" otherwise.
	Node string

	// Address is the IP:port of the server, used for Raft communications.
	Address string

	// Leader is true if this server is the current cluster leader.
	Leader bool

	// Voter is true if this server has a vote in the cluster. This might
	// be false if the server is staging and still coming online, or if
	// it's a non-voting server, which will be added in a future release of
	// Nomad.
	Voter bool

	// RaftProtocol is the version of the Raft protocol spoken by this server.
	RaftProtocol string
}

// RaftConfiguration is returned when querying for the current Raft configuration.
type RaftConfiguration struct {
	// Servers has the list of servers in the Raft configuration.
	Servers []*RaftServer

	// Index has the Raft index of this configuration.
	Index uint64
}

// RaftGetConfiguration is used to query the current Raft peer set.
func (op *Operator) RaftGetConfiguration(q *QueryOptions) (*RaftConfiguration, error) {
	r, err := op.c.newRequest("GET", "/v1/operator/raft/configuration")
	if err != nil {
		return nil, err
	}
	r.setQueryOptions(q)
	_, resp, err := requireOK(op.c.doRequest(r))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out RaftConfiguration
	if err := decodeBody(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RaftRemovePeerByAddress is used to kick a stale peer (one that it in the Raft
// quorum but no longer known to Serf or the catalog) by address in the form of
// "IP:port".
func (op *Operator) RaftRemovePeerByAddress(address string, q *WriteOptions) error {
	r, err := op.c.newRequest("DELETE", "/v1/operator/raft/peer")
	if err != nil {
		return err
	}
	r.setWriteOptions(q)

	r.params.Set("address", address)

	_, resp, err := requireOK(op.c.doRequest(r))
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}

// RaftRemovePeerByID is used to kick a stale peer (one that is in the Raft
// quorum but no longer known to Serf or the catalog) by ID.
func (op *Operator) RaftRemovePeerByID(id string, q *WriteOptions) error {
	r, err := op.c.newRequest("DELETE", "/v1/operator/raft/peer")
	if err != nil {
		return err
	}
	r.setWriteOptions(q)

	r.params.Set("id", id)

	_, resp, err := requireOK(op.c.doRequest(r))
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}

type SchedulerConfiguration struct {
	// PreemptionConfig specifies whether to enable eviction of lower
	// priority jobs to place higher priority jobs.
	PreemptionConfig PreemptionConfig

	// CreateIndex/ModifyIndex store the create/modify indexes of this configuration.
	CreateIndex uint64
	ModifyIndex uint64
}

// SchedulerConfigurationResponse is the response object that wraps SchedulerConfiguration
type SchedulerConfigurationResponse struct {
	// SchedulerConfig contains scheduler config options
	SchedulerConfig *SchedulerConfiguration

	QueryMeta
}

// SchedulerSetConfigurationResponse is the response object used
// when updating scheduler configuration
type SchedulerSetConfigurationResponse struct {
	// Updated returns whether the config was actually updated
	// Only set when the request uses CAS
	Updated bool

	WriteMeta
}

// PreemptionConfig specifies whether preemption is enabled based on scheduler type
type PreemptionConfig struct {
	SystemSchedulerEnabled bool
}

// SchedulerGetConfiguration is used to query the current Scheduler configuration.
func (op *Operator) SchedulerGetConfiguration(q *QueryOptions) (*SchedulerConfigurationResponse, *QueryMeta, error) {
	var resp SchedulerConfigurationResponse
	qm, err := op.c.query("/v1/operator/scheduler/configuration", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// SchedulerSetConfiguration is used to set the current Scheduler configuration.
func (op *Operator) SchedulerSetConfiguration(conf *SchedulerConfiguration, q *WriteOptions) (*SchedulerSetConfigurationResponse, *WriteMeta, error) {
	var out SchedulerSetConfigurationResponse
	wm, err := op.c.write("/v1/operator/scheduler/configuration", conf, &out, q)
	if err != nil {
		return nil, nil, err
	}
	return &out, wm, nil
}

// SchedulerCASConfiguration is used to perform a Check-And-Set update on the
// Scheduler configuration. The ModifyIndex value will be respected. Returns
// true on success or false on failures.
func (op *Operator) SchedulerCASConfiguration(conf *SchedulerConfiguration, q *WriteOptions) (*SchedulerSetConfigurationResponse, *WriteMeta, error) {
	var out SchedulerSetConfigurationResponse
	wm, err := op.c.write("/v1/operator/scheduler/configuration?cas="+strconv.FormatUint(conf.ModifyIndex, 10), conf, &out, q)
	if err != nil {
		return nil, nil, err
	}

	return &out, wm, nil
}
