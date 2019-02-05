---
title: Customizing your gateway
category: documentation
order: 101
---

# Customizing your gateway

## Message serialization
### Custom serializer

You might want to serialize request/response messages in MessagePack instead of JSON, for example.

1. Write a custom implementation of [`Marshaler`](http://godoc.org/github.com/grpc-ecosystem/grpc-gateway/runtime#Marshaler)
2. Register your marshaler with [`WithMarshalerOption`](http://godoc.org/github.com/grpc-ecosystem/grpc-gateway/runtime#WithMarshalerOption)
   e.g.
   ```go
   var m your.MsgPackMarshaler
   mux := runtime.NewServeMux(runtime.WithMarshalerOption("application/x-msgpack", m))
   ```

You can see [the default implementation for JSON](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/runtime/marshal_jsonpb.go) for reference.

### Using camelCase for JSON

The protocol buffer compiler generates camelCase JSON tags that can be used with jsonpb package. By default jsonpb Marshaller uses `OrigName: true` which uses the exact case used in the proto files. To use camelCase for the JSON representation,
   ```go
   mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName:false}))
   ```

## Mapping from HTTP request headers to gRPC client metadata
You might not like [the default mapping rule](http://godoc.org/github.com/grpc-ecosystem/grpc-gateway/runtime#DefaultHeaderMatcher) and might want to pass through all the HTTP headers, for example.

1. Write a [`HeaderMatcherFunc`](http://godoc.org/github.com/grpc-ecosystem/grpc-gateway/runtime#HeaderMatcherFunc).
2. Register the function with [`WithIncomingHeaderMatcher`](http://godoc.org/github.com/grpc-ecosystem/grpc-gateway/runtime#WithIncomingHeaderMatcher)

   e.g.
   ```go
   func yourMatcher(headerName string) (mdName string, ok bool) {
   	...
   }
   ...
   mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(yourMatcher))

   ```

## Mapping from gRPC server metadata to HTTP response headers
ditto. Use [`WithOutgoingHeaderMatcher`](http://godoc.org/github.com/grpc-ecosystem/grpc-gateway/runtime#WithOutgoingHeaderMatcher)

## Mutate response messages or set response headers
You might want to return a subset of response fields as HTTP response headers; 
You might want to simply set an application-specific token in a header.
Or you might want to mutate the response messages to be returned.

1. Write a filter function.
   ```go
   func myFilter(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
   	w.Header().Set("X-My-Tracking-Token", resp.Token)
   	resp.Token = ""
   	return nil
   }
   ```
2. Register the filter with [`WithForwardResponseOption`](http://godoc.org/github.com/grpc-ecosystem/grpc-gateway/runtime#WithForwardResponseOption)
   
   e.g.
   ```go
   mux := runtime.NewServeMux(runtime.WithForwardResponseOption(myFilter))
   ```

## OpenTracing Support

If your project uses [OpenTracing](https://github.com/opentracing/opentracing-go) and you'd like spans to propagate through the gateway, you can add some middleware which parses the incoming HTTP headers to create a new span correctly.

```go
import (
   ...
   "github.com/opentracing/opentracing-go"
   "github.com/opentracing/opentracing-go/ext"
)

var grpcGatewayTag = opentracing.Tag{Key: string(ext.Component), Value: "grpc-gateway"}

func tracingWrapper(h http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    parentSpanContext, err := opentracing.GlobalTracer().Extract(
      opentracing.HTTPHeaders,
      opentracing.HTTPHeadersCarrier(r.Header))
    if err == nil || err == opentracing.ErrSpanContextNotFound {
      serverSpan := opentracing.GlobalTracer().StartSpan(
        "ServeHTTP",
        // this is magical, it attaches the new span to the parent parentSpanContext, and creates an unparented one if empty.
        ext.RPCServerOption(parentSpanContext),
        grpcGatewayTag,
      )
      r = r.WithContext(opentracing.ContextWithSpan(r.Context(), serverSpan))
      defer serverSpan.Finish()
    }
    h.ServeHTTP(w, r)
  })
}

// Then just wrap the mux returned by runtime.NewServeMux() like this
if err := http.ListenAndServe(":8080", tracingWrapper(mux)); err != nil {
  log.Fatalf("failed to start gateway server on 8080: %v", err)
}
```

## Error handler
http://mycodesmells.com/post/grpc-gateway-error-handler

## Replace a response forwarder per method
You might want to keep the behavior of the current marshaler but change only a message forwarding of a certain API method.

1. write a custom forwarder which is compatible to [`ForwardResponseMessage`](http://godoc.org/github.com/grpc-ecosystem/grpc-gateway/runtime#ForwardResponseMessage) or [`ForwardResponseStream`](http://godoc.org/github.com/grpc-ecosystem/grpc-gateway/runtime#ForwardResponseStream).
2. replace the default forwarder of the method with your one.

   e.g. add `forwarder_overwrite.go` into the go package of the generated code,
   ```go
   package generated
   
   import (
   	"net/http"

   	"github.com/grpc-ecosystem/grpc-gateway/runtime"
   	"github.com/golang/protobuf/proto"
   	"golang.org/x/net/context"
   )

   func forwardCheckoutResp(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, resp proto.Message, opts ...func(context.Context, http.ResponseWriter, proto.Message) error) {
   	if someCondition(resp) {
   		http.Error(w, "not enough credit", http. StatusPaymentRequired)
   		return
   	}
   	runtime.ForwardResponseMessage(ctx, mux, marshaler, w, req, resp, opts...)
   }
   
   func init() {
   	forward_MyService_Checkout_0 = forwardCheckoutResp
   }
   ```
