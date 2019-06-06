# Azure SDK for Go

[![godoc](https://godoc.org/github.com/Azure/azure-sdk-for-go?status.svg)](https://godoc.org/github.com/Azure/azure-sdk-for-go)

[![Build Status](https://travis-ci.org/Azure/azure-sdk-for-go.svg?branch=master)](https://travis-ci.org/Azure/azure-sdk-for-go)

[![Go Report Card](https://goreportcard.com/badge/github.com/Azure/azure-sdk-for-go)](https://goreportcard.com/report/github.com/Azure/azure-sdk-for-go)

azure-sdk-for-go provides Go packages for managing and using Azure services.
It officially supports the last two major releases of Go.  Older versions of
Go will be kept running in CI until they no longer work due to changes in any
of the SDK's external dependencies.  The CHANGELOG will be updated when a
version of Go is removed from CI.

To be notified about updates and changes, subscribe to the [Azure update
feed](https://azure.microsoft.com/updates/).

Users may prefer to jump right in to our samples repo at
[github.com/Azure-Samples/azure-sdk-for-go-samples][samples_repo].

Questions and feedback? Chat with us in the **[#Azure SDK
channel](https://gophers.slack.com/messages/CA7HK8EEP)** on the [Gophers
Slack](https://gophers.slack.com/). Sign up
[here](https://invite.slack.golangbridge.org) first if necessary.


## Package Updates

Most packages in the SDK are generated from [Azure API specs][azure_rest_specs]
using [Azure/autorest.go][] and [Azure/autorest][]. These generated packages
depend on the HTTP client implemented at [Azure/go-autorest][].

[azure_rest_specs]: https://github.com/Azure/azure-rest-api-specs
[azure/autorest]: https://github.com/Azure/autorest
[azure/autorest.go]: https://github.com/Azure/autorest.go
[azure/go-autorest]: https://github.com/Azure/go-autorest

The SDK codebase adheres to [semantic versioning](https://semver.org) and thus
avoids breaking changes other than at major (x.0.0) releases. Because Azure's
APIs are updated frequently, we release a **new major version at the end of
each month** with a full changelog. For more details and background see [SDK Update
Practices](https://github.com/Azure/azure-sdk-for-go/wiki/SDK-Update-Practices).

To more reliably manage dependencies like the Azure SDK in your applications we
recommend [golang/dep](https://github.com/golang/dep).

Packages that are still in public preview can be found under the ./services/preview
directory. Please be aware that since these packages are in preview they are subject
to change, including breaking changes outside of a major semver bump.

## Other Azure Go Packages

Azure provides several other packages for using services from Go, listed below.
If a package you need isn't available please open an issue and let us know.

| Service              | Import Path/Repo                                                                                   |
| -------------------- | -------------------------------------------------------------------------------------------------- |
| Storage - Blobs      | [github.com/Azure/azure-storage-blob-go](https://github.com/Azure/azure-storage-blob-go)           |
| Storage - Files      | [github.com/Azure/azure-storage-file-go](https://github.com/Azure/azure-storage-file-go)           |
| Storage - Queues     | [github.com/Azure/azure-storage-queue-go](https://github.com/Azure/azure-storage-queue-go)         |
| Service Bus          | [github.com/Azure/azure-service-bus-go](https://github.com/Azure/azure-service-bus-go)             |
| Event Hubs           | [github.com/Azure/azure-event-hubs-go](https://github.com/Azure/azure-event-hubs-go)               |
| Application Insights | [github.com/Microsoft/ApplicationInsights-go](https://github.com/Microsoft/ApplicationInsights-go) |

# Install and Use:

## Install

```sh
$ go get -u github.com/Azure/azure-sdk-for-go/...
```

or if you use dep, within your repo run:

```sh
$ dep ensure -add github.com/Azure/azure-sdk-for-go
```

If you need to install Go, follow [the official instructions](https://golang.org/dl/).

## Use

For many more scenarios and examples see
[Azure-Samples/azure-sdk-for-go-samples][samples_repo].

Apply the following general steps to use packages in this repo. For more on
authentication and the `Authorizer` interface see [the next
section](#authentication).

1.  Import a package from the [services][services_dir] directory.
2.  Create and authenticate a client with a `New*Client` func, e.g.
    `c := compute.NewVirtualMachinesClient(...)`.
3.  Invoke API methods using the client, e.g.
    `res, err := c.CreateOrUpdate(...)`.
4.  Handle responses and errors.

[services_dir]: https://github.com/Azure/azure-sdk-for-go/tree/master/services

For example, to create a new virtual network (substitute your own values for
strings in angle brackets):

```go
package main

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"

	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

func main() {
	// create a VirtualNetworks client
	vnetClient := network.NewVirtualNetworksClient("<subscriptionID>")

	// create an authorizer from env vars or Azure Managed Service Idenity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		vnetClient.Authorizer = authorizer
	}

	// call the VirtualNetworks CreateOrUpdate API
	vnetClient.CreateOrUpdate(context.Background(),
		"<resourceGroupName>",
		"<vnetName>",
		network.VirtualNetwork{
			Location: to.StringPtr("<azureRegion>"),
			VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
				AddressSpace: &network.AddressSpace{
					AddressPrefixes: &[]string{"10.0.0.0/8"},
				},
				Subnets: &[]network.Subnet{
					{
						Name: to.StringPtr("<subnet1Name>"),
						SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
							AddressPrefix: to.StringPtr("10.0.0.0/16"),
						},
					},
					{
						Name: to.StringPtr("<subnet2Name>"),
						SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
							AddressPrefix: to.StringPtr("10.1.0.0/16"),
						},
					},
				},
			},
		})
}
```

## Authentication

Typical SDK operations must be authenticated and authorized. The _Authorizer_
interface allows use of any auth style in requests, such as inserting an OAuth2
Authorization header and bearer token received from Azure AD.

The SDK itself provides a simple way to get an authorizer which first checks
for OAuth client credentials in environment variables and then falls back to
Azure's [Managed Service Identity]() when available, e.g. when on an Azure
VM. The following snippet from [the previous section](#use) demonstrates
this helper.

```go
import github.com/Azure/go-autorest/autorest/azure/auth

// create a VirtualNetworks client
vnetClient := network.NewVirtualNetworksClient("<subscriptionID>")

// create an authorizer from env vars or Azure Managed Service Idenity
authorizer, err := auth.NewAuthorizerFromEnvironment()
if err == nil {
    vnetClient.Authorizer = authorizer
}

// call the VirtualNetworks CreateOrUpdate API
vnetClient.CreateOrUpdate(context.Background(),
// ...
```

The following environment variables help determine authentication configuration:

- `AZURE_ENVIRONMENT`: Specifies the Azure Environment to use. If not set, it
  defaults to `AzurePublicCloud`. Not applicable to authentication with Managed
  Service Identity (MSI).
- `AZURE_AD_RESOURCE`: Specifies the AAD resource ID to use. If not set, it
  defaults to `ResourceManagerEndpoint` for operations with Azure Resource
  Manager. You can also choose an alternate resource programatically with
  `auth.NewAuthorizerFromEnvironmentWithResource(resource string)`.

### More Authentication Details

The previous is the first and most recommended of several authentication
options offered by the SDK because it allows seamless use of both service
principals and [Azure Managed Service Identity][]. Other options are listed
below.

> Note: If you need to create a new service principal, run `az ad sp create-for-rbac -n "<app_name>"` in the
> [azure-cli](https://github.com/Azure/azure-cli). See [these
> docs](https://docs.microsoft.com/cli/azure/create-an-azure-service-principal-azure-cli?view=azure-cli-latest)
> for more info. Copy the new principal's ID, secret, and tenant ID for use in
> your app, or consider the `--sdk-auth` parameter for serialized output.

[azure managed service identity]: https://docs.microsoft.com/en-us/azure/active-directory/msi-overview

- The `auth.NewAuthorizerFromEnvironment()` described above creates an authorizer
  from the first available of the following configuration:

      1. **Client Credentials**: Azure AD Application ID and Secret.

          - `AZURE_TENANT_ID`: Specifies the Tenant to which to authenticate.
          - `AZURE_CLIENT_ID`: Specifies the app client ID to use.
          - `AZURE_CLIENT_SECRET`: Specifies the app secret to use.

      2. **Client Certificate**: Azure AD Application ID and X.509 Certificate.

          - `AZURE_TENANT_ID`: Specifies the Tenant to which to authenticate.
          - `AZURE_CLIENT_ID`: Specifies the app client ID to use.
          - `AZURE_CERTIFICATE_PATH`: Specifies the certificate Path to use.
          - `AZURE_CERTIFICATE_PASSWORD`: Specifies the certificate password to use.

      3. **Resource Owner Password**: Azure AD User and Password. This grant type is *not
         recommended*, use device login instead if you need interactive login.

          - `AZURE_TENANT_ID`: Specifies the Tenant to which to authenticate.
          - `AZURE_CLIENT_ID`: Specifies the app client ID to use.
          - `AZURE_USERNAME`: Specifies the username to use.
          - `AZURE_PASSWORD`: Specifies the password to use.

      4. **Azure Managed Service Identity**: Delegate credential management to the
         platform. Requires that code is running in Azure, e.g. on a VM. All
         configuration is handled by Azure. See [Azure Managed Service
         Identity](https://docs.microsoft.com/en-us/azure/active-directory/msi-overview)
         for more details.

- The `auth.NewAuthorizerFromFile()` method creates an authorizer using
  credentials from an auth file created by the [Azure CLI][]. Follow these
  steps to utilize:

  1. Create a service principal and output an auth file using `az ad sp create-for-rbac --sdk-auth > client_credentials.json`.
  2. Set environment variable `AZURE_AUTH_LOCATION` to the path of the saved
     output file.
  3. Use the authorizer returned by `auth.NewAuthorizerFromFile()` in your
     client as described above.

- The `auth.NewAuthorizerFromCLI()` method creates an authorizer which
  uses [Azure CLI][] to obtain its credentials. To use this method follow
  these steps:

  1. Install [Azure CLI v2.0.12](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) or later. Upgrade earlier versions.
  2. Use `az login` to sign in to Azure.

  If you receive an error, use `az account get-access-token` to verify access.

  If Azure CLI is not installed to the default directory, you may receive an error
  reporting that `az` cannot be found.  
  Use the `AzureCLIPath` environment variable to define the Azure CLI installation folder.

  If you are signed in to Azure CLI using multiple accounts or your account has
  access to multiple subscriptions, you need to specify the specific subscription
  to be used. To do so, use:

  ```
  az account set --subscription <subscription-id>
  ```

  To verify the current account settings, use:

  ```
  az account list
  ```

[azure cli]: https://github.com/Azure/azure-cli

- Finally, you can use OAuth's [Device Flow][] by calling
  `auth.NewDeviceFlowConfig()` and extracting the Authorizer as follows:

  ```go
  config := auth.NewDeviceFlowConfig(clientID, tenantID)
  a, err = config.Authorizer()
  ```

[device flow]: https://oauth.net/2/device-flow/

# Versioning

azure-sdk-for-go provides at least a basic Go binding for every Azure API. To
provide maximum flexibility to users, the SDK even includes previous versions of
Azure APIs which are still in use. This enables us to support users of the
most updated Azure datacenters, regional datacenters with earlier APIs, and
even on-premises installations of Azure Stack.

**SDK versions** apply globally and are tracked by git
[tags](https://github.com/Azure/azure-sdk-for-go/tags). These are in x.y.z form
and generally adhere to [semantic versioning](https://semver.org) specifications.

**Service API versions** are generally represented by a date string and are
tracked by offering separate packages for each version. For example, to choose the
latest API versions for Compute and Network, use the following imports:

```go
import (
    "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
    "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
)
```

Occasionally service-side changes require major changes to existing versions.
These cases are noted in the changelog, and for this reason `Service API versions`
cannot be used alone to ensure backwards compatibility.

All available services and versions are listed under the `services/` path in
this repo and in [GoDoc][services_godoc]. Run `find ./services -type d -mindepth 3` to list all available service packages.

[services_godoc]: https://godoc.org/github.com/Azure/azure-sdk-for-go/services

### Profiles

Azure **API profiles** specify subsets of Azure APIs and versions. Profiles can provide:

- **stability** for your application by locking to specific API versions; and/or
- **compatibility** for your application with Azure Stack and regional Azure datacenters.

In the Go SDK, profiles are available under the `profiles/` path and their
component API versions are aliases to the true service package under
`services/`. You can use them as follows:

```go
import "github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/compute/mgmt/compute"
import "github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/network/mgmt/network"
import "github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/storage/mgmt/storage"
```

The following profiles are available for hybrid Azure and Azure Stack environments.
- 2017-03-09
- 2018-03-01

In addition to versioned profiles, we also provide two special profiles
`latest` and `preview`. The `latest` profile contains the latest API version
of each service, excluding any preview versions and/or content.  The `preview`
profile is similar to the `latest` profile but includes preview API versions.

The `latest` and `preview` profiles can help you stay up to date with API
updates as you build applications. Since they are by definition not stable,
however, they **should not** be used in production apps. Instead, choose the
latest specific API version (or an older one if necessary) from the `services/`
path.

As an example, to automatically use the most recent Compute APIs, use one of
the following imports:

```go
import "github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
import "github.com/Azure/azure-sdk-for-go/profiles/preview/compute/mgmt/compute"
```

### Avoiding Breaking Changes

To avoid breaking changes, when specifying imports you should specify a `Service API Version` or `Profile`, as well as lock (using [dep](https://github.com/golang/dep) and soon with [Go Modules](https://github.com/golang/go/wiki/Modules)) to a specific SDK version.

For example, in your source code imports, use a `Service API Version` (`2017-12-01`):

```go
import "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
```

or `Profile` version (`2017-03-09`):

```go
import "github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/compute/mgmt/compute"
```

As well as, for dep, a `Gopkg.toml` file with:

```toml
[[constraint]]
  name = "github.com/Azure/azure-sdk-for-go"
  version = "21.0.0"
```

Combined, these techniques will ensure that breaking changes should not occur. If you are extra sensitive to changes, adding an additional [version pin](https://golang.github.io/dep/docs/Gopkg.toml.html#version-rules) in your SDK Version should satisfy your needs:

```toml
[[constraint]]
  name = "github.com/Azure/azure-sdk-for-go"
  version = "=21.3.0"
```

## Inspecting and Debugging

### Built-in Basic Request/Response Logging

Starting with `go-autorest v10.15.0` you can enable basic logging of requests and responses through setting environment variables.
Setting `AZURE_GO_SDK_LOG_LEVEL` to `INFO` will log request/response without their bodies. To include the bodies set the log level to `DEBUG`.

By default the logger writes to strerr, however it can also write to stdout or a file
if specified in `AZURE_GO_SDK_LOG_FILE`. Note that if the specified file already exists it will be truncated.

**IMPORTANT:** by default the logger will redact the Authorization and Ocp-Apim-Subscription-Key
headers. Any other secrets will _not_ be redacted.

### Writing Custom Request/Response Inspectors

All clients implement some handy hooks to help inspect the underlying requests being made to Azure.

- `RequestInspector`: View and manipulate the go `http.Request` before it's sent
- `ResponseInspector`: View the `http.Response` received

Here is an example of how these can be used with `net/http/httputil` to see requests and responses.

```go
vnetClient := network.NewVirtualNetworksClient("<subscriptionID>")
vnetClient.RequestInspector = LogRequest()
vnetClient.ResponseInspector = LogResponse()

...

func LogRequest() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err != nil {
				log.Println(err)
			}
			dump, _ := httputil.DumpRequestOut(r, true)
			log.Println(string(dump))
			return r, err
		})
	}
}

func LogResponse() autorest.RespondDecorator {
	return func(p autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(r *http.Response) error {
			err := p.Respond(r)
			if err != nil {
				log.Println(err)
			}
			dump, _ := httputil.DumpResponse(r, true)
			log.Println(string(dump))
			return err
		})
	}
}
```

## Tracing and Metrics

All packages and the runtime are instrumented using [OpenCensus](https://opencensus.io/).

### Enable

As of now, tracing is disabled by default. There are 2 ways to enable tracing:

- set the environment variable `AZURE_SDK_TRACING_ENABLED` (_Recommended_)
- alternatively, import the `github.com/Azure/go-autorest/tracing` package and call the `tracing.Enable()` function or `tracing.EnableWithAIForwarding()` if using the [App Insights Forwarder](https://docs.microsoft.com/en-us/azure/application-insights/opencensus-local-forwarder).

**Note**: In future major releases of the SDK, tracing may become enabled by default.

### Usage

Once enabled, all SDK calls will emit traces and metrics and the traces will correlate the SDK calls with the raw http calls made to Azure API's. To consume those traces, if are not doing it yet, you need to register an exporter of your choice such as [Azure App Insights](https://docs.microsoft.com/en-us/azure/application-insights/opencensus-local-forwarder) or [Zipkin](https://opencensus.io/quickstart/go/tracing/#exporting-traces).

To correlate the SDK calls between them and with the rest of your code, pass in a context that has a span initiated using the [opencensus-go library](https://github.com/census-instrumentation/opencensus-go) using the `trace.Startspan(ctx context.Context, name string, o ...StartOption)` function. Here is an example:

```go
func doAzureCalls() {
    // The resulting context will be initialized with a root span as the context passed to
    // trace.StartSpan() has no existing span.
    ctx, span := trace.StartSpan(context.Background(),"doAzureCalls", trace.WithSampler(trace.AlwaysSample()))
    defer span.End()

    // The traces from the SDK calls will be correlated under the span inside the context that is passed in.
    zone, _ := zonesClient.CreateOrUpdate(ctx, rg, zoneName, dns.Zone{Location: to.StringPtr("global")}, "", "")
    zone, _ = zonesClient.Get(ctx, rg, *zone.Name)
    for i := 0; i < rrCount; i++ {
        rr, _ := recordsClient.CreateOrUpdate(ctx, rg, zoneName, fmt.Sprintf("rr%d", i), dns.CNAME, rdSet{
            RecordSetProperties: &dns.RecordSetProperties{
                TTL: to.Int64Ptr(3600),
                CnameRecord: &dns.CnameRecord{
                    Cname: to.StringPtr("vladdbCname"),
                },
            },
        },
            "",
            "",
        )
    }
}
```

# Resources

- SDK docs are at [godoc.org](https://godoc.org/github.com/Azure/azure-sdk-for-go/).
- SDK samples are at [Azure-Samples/azure-sdk-for-go-samples](https://github.com/Azure-Samples/azure-sdk-for-go-samples).
- SDK notifications are published via the [Azure update feed](https://azure.microsoft.com/updates/).
- Azure API docs are at [docs.microsoft.com/rest/api](https://docs.microsoft.com/rest/api/).
- General Azure docs are at [docs.microsoft.com/azure](https://docs.microsoft.com/azure).

## License

Apache 2.0, see [LICENSE](./LICENSE).

## Contribute

See [CONTRIBUTING.md](./CONTRIBUTING.md).

[samples_repo]: https://github.com/Azure-Samples/azure-sdk-for-go-samples
