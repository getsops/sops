---
layout: "api"
page_title: "/sys/metrics - HTTP API"
sidebar_title: "<code>/sys/metrics</code>"
sidebar_current: "api-http-system-metrics"
description: |-
  The `/sys/metrics` endpoint is used to get telemetry metrics for Vault.
---

# `/sys/metrics`

The `/sys/metrics` endpoint is used to get telemetry metrics for Vault.

## Read Telemetry Metrics

This endpoint returns the telemetry metrics for Vault. It can be used by metrics
collections systems like [Prometheus](https://prometheus.io) that use a pull
model for metrics collection.

| Method   | Path             |
| :------- | :--------------- |
| `GET`    | `/sys/metrics`   |

### Parameters

- `format` `(string: "")` – Specifies the format used for the returned metrics. The
  default metrics format is JSON. Setting `format` to `prometheus` will return the
  metrics in [Prometheus format](https://prometheus.io/docs/instrumenting/exposition_formats/#text-based-format).
  
### Sample Request

```
$ curl -H "X-Vault-Token: f3b09679-3001-009d-2b80-9c306ab81aa6" \
    http://127.0.0.1:8200/v1/sys/metrics?format=prometheus
```

### Sample Response

This response is only returned for a `GET` request.

```
# HELP vault_audit_log_request vault_audit_log_request
# TYPE vault_audit_log_request summary
vault_audit_log_request{quantile="0.5"} 0.005927000194787979
vault_audit_log_request{quantile="0.9"} 0.005927000194787979
vault_audit_log_request{quantile="0.99"} 0.005927000194787979
vault_audit_log_request_sum 0.014550999738276005
vault_audit_log_request_count 2
# HELP vault_audit_log_request_failure vault_audit_log_request_failure
# TYPE vault_audit_log_request_failure counter
vault_audit_log_request_failure 0
# HELP vault_audit_log_response vault_audit_log_response
# TYPE vault_audit_log_response summary
vault_audit_log_response{quantile="0.5"} NaN
vault_audit_log_response{quantile="0.9"} NaN
vault_audit_log_response{quantile="0.99"} NaN
vault_audit_log_response_sum 0.0057669999077916145
vault_audit_log_response_count 1
# HELP vault_audit_log_response_failure vault_audit_log_response_failure
# TYPE vault_audit_log_response_failure counter
vault_audit_log_response_failure 0
# HELP vault_barrier_get vault_barrier_get
# TYPE vault_barrier_get summary
vault_barrier_get{quantile="0.5"} 0.011938000097870827
vault_barrier_get{quantile="0.9"} 0.011938000097870827
vault_barrier_get{quantile="0.99"} 0.011938000097870827
vault_barrier_get_sum 0.1814980012131855
vault_barrier_get_count 36
```
