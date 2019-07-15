---
layout: "docs"
page_title: "Performance Standby Nodes - Vault Enterprise"
sidebar_title: "Performance Standbys"
sidebar_current: "docs-vault-enterprise-perf-standbys"
description: |-
  Performance Standby Nodes - Vault Enterprise
---

# Performance Standby Nodes

Vault supports a multi-server mode for high availability. This mode protects
against outages by running multiple Vault servers. High availability mode
is automatically enabled when using a data store that supports it. You can
learn more about HA mode on the [Concepts](/docs/concepts/ha.html) page.

Vault Enterprise offers additional features that allow HA nodes to service
read-only requests on the local standby node. Read-only requests are requests
that do not modify Vault's storage.

## Server-to-Server Communication

Performance Standbys require the request forwarding method described in the [HA
Server-to-Server](/docs/concepts/ha.html#server-to-server-communication) docs.

A performance standby will connect to the active node over the existing request
forwarding connection. If selected by the active node to be promoted to a
performance standby it will be handed a newly-generated private key and certificate
for use in creating a new mutually-authenticated TLS connection to the cluster
port. This connection will be used to send updates from the active node to the
standby.

## Request Forwarding

A Performance Standby will attempt to process requests that come in. If a
storage write is detected the standby will forward the request over the cluster
port connection to the active node. If the request is read-only the Performance
Standby will handle the requests locally.

Sending requests to Performance Standbys that result in forwarded writes will be
slightly slower than going directly to the active node. A client that has
advanced knowledge of the behavior of the call can choose to point the request
to the appropriate node.

### Direct Access

A Performance Standby will tag itself as such in consul if service registration
is enabled. To access the set of Performance Standbys the `performance-standby`
tag can be used. For example to send requests to only the performance standbys
`https://performance-standby.vault.dc1.consul` could be used (host name may vary
based on consul configuration).

### Behind Load Balancers

Additionally, if you wish to point your load balancers at performance standby
nodes, the `sys/health` endpoint can be used to determine if a node is a
performance standby. See the [sys/health API](/api/system/health.html) docs for
more info.

## Disabling Performance Standbys

To disable performance standbys the `disable_performance_standby` flag should be
set to true in the Vault config file. This will both tell a standby not to
attempt to enable performance mode and an active node to not allow any
performance standby connections.

This setting should be synced across all nodes in the cluster.

## Monitoring Performance Standbys

To verify your node is a performance standby the `vault status` command can be
used:

```
$ vault status
Key                                    Value
---                                    -----
Seal Type                              shamir
Sealed                                 false
Total Shares                           1
Threshold                              1
Version                                0.11.0+prem
Cluster Name                           vault-cluster-d040e74c
Cluster ID                             9f82e03b-71fb-97a6-9c5a-46fa6715d6e4
HA Enabled                             true
HA Cluster                             https://127.0.0.1:8201
HA Mode                                standby
Active Node Address                    http://127.0.0.1:8200
Performance Standby Node               true
Performance Standby Last Remote WAL    380329
```
