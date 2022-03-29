# Ports used by CAPA

Name      | Port Number | Description |
---       | ---         | ---
`metrics` |             | Port that exposes the metrics. This can be customized by setting the `--metrics-bind-addr` flag when starting the manager. The default is to only listen on `localhost:8080`
`webhook` | `9443`      | Webhook server port. To disable this set `--webhook-port` flag to `0`.
`health`  | `9440`      | Port that exposes the health endpoint. This can be customized by setting the `--health-addr` flag when starting the manager.
`profiler`|             | Expose the pprof profiler. By default is not configured. Can set the `--profiler-address` flag. e.g. `--profiler-address 6060`
