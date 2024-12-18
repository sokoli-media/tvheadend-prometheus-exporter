# tvheadend-prometheus-exporter

Simple project to export tvheadend metrics to Prometheus.

Docker image: `sokolimedia/tvheadend-prometheus-exporter:latest`

Project exports http api on `:9000` with metrics at `/metrics` url.

Required environmental variables:
* `TVHEADEND_HOST`: host of your tvheadend instance
* `TVHEADEND_USERNAME`: username for your tvheadend account (must have access to the api)
* `TVHEADEND_PASSWORD`: password for the same account
