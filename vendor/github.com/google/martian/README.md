# Martian Proxy [![Build Status](https://travis-ci.org/google/martian.svg?branch=master)](https://travis-ci.org/google/martian)

Martian Proxy is a programmable HTTP proxy designed to be used for testing.

Martian is a great tool to use if you want to:

* Verify that all (or some subset) of requests are secure
* Mock external services at the network layer
* Inject headers, modify cookies or perform other mutations of HTTP requests
  and responses
* Verify that pingbacks happen when you think they should
* Unwrap encrypted traffic (requires install of CA certificate in browser)

By taking advantage of Go cross-compilation, Martian can be deployed
anywhere that Go can target.

## Latest Version

v3.0.0

## Requirements

Go 1.11

## Go Modules Support

Martian Proxy added support for Go modules since v3.0.0.  If you use a Go version that does not support modules, this will break you.  The latest version without Go modules support was tagged v2.1.0.

## Getting Started

### Installation
Martian Proxy can be installed using `go install`

    go get github.com/google/martian/ && \
    go install github.com/google/martian/cmd/proxy

### Start the Proxy
Assuming you've installed Martian, running the proxy is as simple as

    $GOPATH/bin/proxy

If you want to see system logs as Martian is running, pass in the verbosity
flag:

    $GOPATH/bin/proxy -v=2

By default, Martian will be running on port 8080, and the Martian API will be running on 8181
. The port can be specified via flags:

    $GOPATH/bin/proxy -addr=:9999 -api-addr=:9898

### Logging
For logging of requests and responses a [logging
modifier](https://github.com/google/martian/wiki/Modifier-Reference#logging) is
available or [HAR](http://www.softwareishard.com/blog/har-12-spec/) logs are
available if the `-har` flag is used.

#### HAR Logging
To enable HAR logging in Martian call the binary with the `-har` flag:

    $GOPATH/bin/proxy -har

If the `-har` flag has been enabled two HAR related endpoints will be
available:

    GET http://martian.proxy/logs

Will retrieve the HAR log of all requests and responses seen by the proxy since
the last reset.

    DELETE http://martian.proxy/logs/reset

Will reset the in-memory HAR log. Note that the log will grow unbounded unless
it is periodically reset.

### Configure
Once Martian is running, you need to configure its behavior. Without
configuration, Martian is just proxying without doing anything to the requests
or responses. If enabled, logging will take place without additional
configuration.

Martian is configured by JSON messages sent over HTTP that take the general
form of:

    {
      "header.Modifier": {
        "scope": ["response"],
        "name": "Test-Header",
        "value": "true"
      }
    }

The above configuration tells Martian to inject a header with the name
"Test-Header" and the value "true" on all responses.

Let's break down the parts of this message.

* `[package.Type]`: The package.Type of the modifier that you want to use. In
  this case, it's "header.Modifier", which is the name of the modifier that
  sets headers (to learn more about the `header.Modifier`, please
  refer to the [modifier reference](https://github.com/google/martian/wiki/Modifier-Reference)).

* `[package.Type].scope`: Indicates whether to apply to the modifier to
  requests, responses or both. This can be an array containing "request",
  "response", or both.

* `[package.Type].[key]`: Modifier specific data. In the case of the header
  modifier, we need the `name` and `value` of the header.

This is a simple configuration, for more complex configurations, modifiers are
combined with groups and filters to compose the desired behavior.

To configure Martian, `POST` the JSON to `http://martian.proxy/modifiers`. You'll
want to use whatever mechanism your language of choice provides you to make
HTTP requests, but for demo purposes, curl works (assuming your configuration
is in a file called `modifier.json`).

        curl -x localhost:8080 \
             -X POST \
             -H "Content-Type: application/json" \
             -d @modifier.json \
                "http://martian.proxy/configure"

### Intercepting HTTPS Requests and Responses
Martian supports modifying HTTPS requests and responses if configured to do so.

In order for Martian to intercept HTTPS traffic a custom CA certificate must be
installed in the browser so that connection warnings are not shown.

The easiest way to install the CA certificate is to start the proxy with the
necessary flags to use a custom CA certificate and private key using the `-cert`
and `-key` flags, or to have the proxy generate one using the `-generate-ca-cert`
flag.

After the proxy has started, visit http://martian.proxy/authority.cer in the
browser configured to use the proxy and a prompt will be displayed to install
the certificate.

Several flags are available in `examples/main.go` to help configure MITM
functionality:

    -key=""
      PEM encoded private key file of the CA certificate provided in -cert; used
      to sign certificates that are generated on-the-fly
    -cert=""
      PEM encoded CA certificate file used to generate certificates
    -generate-ca-cert=false
      generates a CA certificate and private key to use for man-in-the-middle;
      most users choosing this option will immediately visit
      http://martian.proxy/authority.cer in the browser whose traffic is to be
      intercepted to install the newly generated CA certificate
    -organization="Martian Proxy"
      organization name set on the dynamically-generated certificates during
      man-in-the-middle
    -validity="1h"
      window of time around the time of request that the dynamically-generated
      certificate is valid for; the duration is set such that the total valid
      timeframe is double the value of validity (1h before & 1h after)

### Check Verifiers
Let's assume that you've configured Martian to verify the presence a specific
header in responses to a specific URL.

Here's a configuration to verify that all requests to `example.com` return
responses with a `200 OK`.

          {
            "url.Filter": {
              "scope": ["request", "response"],
              "host" : "example.com",
              "modifier" : {
                "status.Verifier": {
                  "scope" : ["response"],
                  "statusCode": 200
                }
              }
            }
          }

Once Martian is running, configured and the requests and resultant responses you
wish to verify have taken place, you can verify your expectation that you only
got back `200 OK` responses.

To check verifications, perform

    GET http://martian.proxy/verify

Failed expectations are tracked as errors, and the list of errors are retrieved
by making a `GET` request to `host:port/martian/verify`, which will return
a list of errors:

      {
          "errors" : [
              {
                  "message": "response(http://example.com) status code verify failure: got 500, want 200"
              },
              {
                  "message": "response(http://example.com/foo) status code verify failure: got 500, want 200"
              }
          ]
      }

Verification errors are held in memory until they are explicitly cleared by

    POST http://martian.proxy/verify/reset

## Martian as a Library
Martian can also be included into any Go program and used as a library.

## Modifiers All The Way Down
Martian's request and response modification system is designed to be general
and extensible. The design objective is to provide individual modifier
behaviors that can arranged to build out nearly any desired modification.

When working with Martian to compose behaviors, you'll need to be familiar with
these different types of interactions:

* Modifiers: Changes the state of a request or a response
* Filters: Conditionally allows a contained Modifier to execute
* Groups: Bundles multiple modifiers to be executed in the order specified in
  the group
* Verifiers: Tracks network traffic against expectations

Modifiers, filters and groups all implement `RequestModifier`,
`ResponseModifier` or `RequestResponseModifier` (defined in
[`martian.go`](https://github.com/google/martian/blob/master/martian.go)).

```go
ModifyRequest(req *http.Request) error

ModifyResponse(res *http.Response) error
```

Throughout the code (and this documentation) you'll see the word "modifier"
used as a term that encompasses modifiers, groups and filters. Even though a
group does not modify a request or response, we still refer to it as a
"modifier".

We refer to anything that implements the `modifier` interface as a Modifier.

### Parser Registration
Each modifier must register its own parser with Martian. The parser is
responsible for parsing a JSON message into a Go struct that implements a
modifier interface.

Martian holds modifier parsers as a map of strings to functions that is built
out at run-time. Each modifier is responsible for registering its parser with a
call to `parse.Register` in `init()`.

Signature of parse.Register:

```go
Register(name string, parseFunc func(b []byte) (interface{}, error))
```

Register takes in the key as a string in the form `package.Type`. For
instance, `cookie_modifier` registers itself with the key `cookie.Modifier` and
`query_string_filter` registers itself as `querystring.Filter`. This string is
the same as the value of `name` in the JSON configuration message.

In the following configuration message, `header.Modifier` is how the header
modifier is registered in the `init()` of `header_modifier.go`.

    {
      "header.Modifier": {
        "scope": ["response"],
        "name" : "Test-Header",
        "value" : "true"
      }
    }

Example of parser registration from `header_modifier.go`:

```go
func init() {
  parse.Register("header.Modifier", modifierFromJSON)
}

func modifierFromJSON(b []byte) (interface{}, error) {
  ...
}
```

### Adding Your Own Modifier
If you have a use-case in mind that we have not developed modifiers, filters or
verifiers for, you can easily extend Martian to your very specific needs.

There are 2 mandatory parts of a modifier:

* Implement the modifier interface
* Register the parser

Any Go struct that implements those interfaces can act as a `modifier`.

## Contact
For questions and comments on how to use Martian, features announcements, or
design discussions check out our public Google Group at
https://groups.google.com/forum/#!forum/martianproxy-users.

For security related issues please send a detailed report to our private core
group at martianproxy-core@googlegroups.com.

## Disclaimer
This is not an official Google product (experimental or otherwise), it is just
code that happens to be owned by Google.
