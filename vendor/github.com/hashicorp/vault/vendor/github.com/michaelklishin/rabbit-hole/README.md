# Rabbit Hole, a RabbitMQ HTTP API Client for Go

This library is a [RabbitMQ HTTP API](https://raw.githack.com/rabbitmq/rabbitmq-management/rabbitmq_v3_6_0/priv/www/api/index.html) client for the Go language.

## Supported Go Versions

Rabbit Hole supports 3 most recent Go releases.


## Supported RabbitMQ Versions

 * RabbitMQ 3.x

All versions require [RabbitMQ Management UI plugin](https://www.rabbitmq.com/management.html) to be installed and enabled.


## Project Maturity

Rabbit Hole is a mature library (first released in late 2013)
designed after a couple of other RabbitMQ HTTP API clients with stable
APIs. Breaking API changes are not out of the question but not without
a reasonable version bump.

It is largely (80-90%) feature complete and decently documented.


## Installation

```
go get github.com/michaelklishin/rabbit-hole
```


## Documentation

### API Reference

[API reference](http://godoc.org/github.com/michaelklishin/rabbit-hole) is available on [godoc.org](http://godoc.org).

Continue reading for a list of example snippets.

### Overview

To import the package:

``` go
import (
       "github.com/michaelklishin/rabbit-hole"
)
```

All HTTP API operations are accessible via `rabbithole.Client`, which
should be instantiated with `rabbithole.NewClient`:

``` go
// URI, username, password
rmqc, _ = NewClient("http://127.0.0.1:15672", "guest", "guest")
```

TLS (HTTPS) can be enabled by adding an HTTP transport to the parameters
of `rabbithole.NewTLSClient`:

``` go
transport := &http.Transport{TLSClientConfig: tlsConfig}
rmqc, _ := NewTLSClient("https://127.0.0.1:15672", "guest", "guest", transport)
```

RabbitMQ HTTP API has to be [configured to use TLS](http://www.rabbitmq.com/management.html#single-listener-https).


### Getting Overview

``` go
res, err := rmqc.Overview()
```

### Node and Cluster Status

``` go
xs, err := rmqc.ListNodes()
// => []NodeInfo, err

node, err := rmqc.GetNode("rabbit@mercurio")
// => NodeInfo, err
```


### Operations on Connections

``` go
xs, err := rmqc.ListConnections()
// => []ConnectionInfo, err

conn, err := rmqc.GetConnection("127.0.0.1:50545 -> 127.0.0.1:5672")
// => ConnectionInfo, err

// Forcefully close connection
_, err := rmqc.CloseConnection("127.0.0.1:50545 -> 127.0.0.1:5672")
// => *http.Response, err
```


### Operations on Channels

``` go
xs, err := rmqc.ListChannels()
// => []ChannelInfo, err

ch, err := rmqc.GetChannel("127.0.0.1:50545 -> 127.0.0.1:5672 (1)")
// => ChannelInfo, err
```


### Operations on Vhosts

``` go
xs, err := rmqc.ListVhosts()
// => []VhostInfo, err

// information about individual vhost
x, err := rmqc.GetVhost("/")
// => VhostInfo, err

// creates or updates individual vhost
resp, err := rmqc.PutVhost("/", VhostSettings{Tracing: false})
// => *http.Response, err

// deletes individual vhost
resp, err := rmqc.DeleteVhost("/")
// => *http.Response, err
```


### Managing Users

``` go
xs, err := rmqc.ListUsers()
// => []UserInfo, err

// information about individual user
x, err := rmqc.GetUser("my.user")
// => UserInfo, err

// creates or updates individual user
resp, err := rmqc.PutUser("my.user", UserSettings{Password: "s3krE7", Tags: "management,policymaker"})
// => *http.Response, err

// creates or updates individual user with no password
resp, err := rmqc.PutUserWithoutPassword("my.user", UserSettings{Tags: "management,policymaker"})
// => *http.Response, err

// deletes individual user
resp, err := rmqc.DeleteUser("my.user")
// => *http.Response, err
```

``` go
// creates or updates individual user with a SHA256 password hash
hash := SaltedPasswordHashSHA256("password-s3krE7")
resp, err := rmqc.PutUser("my.user", UserSettings{
  PasswordHash: hash,
  HashingAlgorithm: HashingAlgorithmSHA256,
  Tags: "management,policymaker"})
// => *http.Response, err
```


### Managing Permissions

``` go
xs, err := rmqc.ListPermissions()
// => []PermissionInfo, err

// permissions of individual user
x, err := rmqc.ListPermissionsOf("my.user")
// => []PermissionInfo, err

// permissions of individual user in vhost
x, err := rmqc.GetPermissionsIn("/", "my.user")
// => PermissionInfo, err

// updates permissions of user in vhost
resp, err := rmqc.UpdatePermissionsIn("/", "my.user", Permissions{Configure: ".*", Write: ".*", Read: ".*"})
// => *http.Response, err

// revokes permissions in vhost
resp, err := rmqc.ClearPermissionsIn("/", "my.user")
// => *http.Response, err
```


### Operations on Exchanges

``` go
xs, err := rmqc.ListExchanges()
// => []ExchangeInfo, err

// list exchanges in a vhost
xs, err := rmqc.ListExchangesIn("/")
// => []ExchangeInfo, err

// information about individual exchange
x, err := rmqc.GetExchange("/", "amq.fanout")
// => ExchangeInfo, err

// declares an exchange
resp, err := rmqc.DeclareExchange("/", "an.exchange", ExchangeSettings{Type: "fanout", Durable: false})
// => *http.Response, err

// deletes individual exchange
resp, err := rmqc.DeleteExchange("/", "an.exchange")
// => *http.Response, err
```


### Operations on Queues

``` go
qs, err := rmqc.ListQueues()
// => []QueueInfo, err

// list queues in a vhost
qs, err := rmqc.ListQueuesIn("/")
// => []QueueInfo, err

// information about individual queue
q, err := rmqc.GetQueue("/", "a.queue")
// => QueueInfo, err

// declares a queue
resp, err := rmqc.DeclareQueue("/", "a.queue", QueueSettings{Durable: false})
// => *http.Response, err

// deletes individual queue
resp, err := rmqc.DeleteQueue("/", "a.queue")
// => *http.Response, err

// purges all messages in queue
resp, err := rmqc.PurgeQueue("/", "a.queue")
// => *http.Response, err
```


### Operations on Bindings

``` go
bs, err := rmqc.ListBindings()
// => []BindingInfo, err

// list bindings in a vhost
bs, err := rmqc.ListBindingsIn("/")
// => []BindingInfo, err

// list bindings of a queue
bs, err := rmqc.ListQueueBindings("/", "a.queue")
// => []BindingInfo, err

// list all bindings having the exchange as source
bs1, err := rmqc.ListExchangeBindingsWithSource("/", "an.exchange")
// => []BindingInfo, err

// list all bindings having the exchange as destinattion
bs2, err := rmqc.ListExchangeBindingsWithDestination("/", "an.exchange")
// => []BindingInfo, err

// declare a binding
resp, err := rmqc.DeclareBinding("/", BindingInfo{
	Source: "an.exchange",
	Destination: "a.queue",
	DestinationType: "queue",
	RoutingKey: "#",
})
// => *http.Response, err

// deletes individual binding
resp, err := rmqc.DeleteBinding("/", BindingInfo{
	Source: "an.exchange",
	Destination: "a.queue",
	DestinationType: "queue",
	RoutingKey: "#",
	PropertiesKey: "%23",
})
// => *http.Response, err
```

### Operations on Shovels

``` go
qs, err := rmqc.ListShovels()
// => []ShovelInfo, err

// list shovels in a vhost
qs, err := rmqc.ListShovelsIn("/")
// => []ShovelInfo, err

// information about an individual shovel
q, err := rmqc.GetShovel("/", "a.shovel")
// => ShovelInfo, err

// declares a shovel
shovelDetails := rabbithole.ShovelDefinition{SourceURI: "amqp://sourceURI", SourceQueue: "mySourceQueue", DestinationURI: "amqp://destinationURI", DestinationQueue: "myDestQueue", AddForwardHeaders: true, AckMode: "on-confirm", DeleteAfter: "never"}
resp, err := rmqc.DeclareShovel("/", "a.shovel", shovelDetails)
// => *http.Response, err

// deletes an individual shovel
resp, err := rmqc.DeleteShovel("/", "a.shovel")
// => *http.Response, err

```

### Operations on cluster name
``` go
// Get cluster name
cn, err := rmqc.GetClusterName()
// => ClusterName, err

// Rename cluster
resp, err := rmqc.SetClusterName(ClusterName{Name: "rabbitmq@rabbit-hole"})
// => *http.Response, err

```

### HTTPS Connections

``` go
var tlsConfig *tls.Config

...

transport := &http.Transport{TLSClientConfig: tlsConfig}

rmqc, err := NewTLSClient("https://127.0.0.1:15672", "guest", "guest", transport)
```

### Changing Transport Layer

``` go
var transport *http.Transport

...

rmqc.SetTransport(transport)
```


## CI Status

[![Build Status](https://travis-ci.org/michaelklishin/rabbit-hole.svg?branch=master)](https://travis-ci.org/michaelklishin/rabbit-hole)


## Contributing

See [CONTRIBUTING.md](https://github.com/michaelklishin/rabbit-hole/blob/master/CONTRIBUTING.md)


## License & Copyright

2-clause BSD license.

(c) Michael S. Klishin, 2013-2019.
