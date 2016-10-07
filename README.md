# Postgres Exporter with own service discovery

Its SD does not depend on Prometheus SD. You only set the exporter address in prometheus config as a target,

then Exporter starts getting metric.

#### SD model

1. Get list of services names from consul

2. get service properties from consul

3. filter with *tag* specified in cli parameters

4. get parameters from the KV. Parameters' path are designed to be /kv/monitoring/*service_name*/*tag*

Advantages:

- Only one service's monitoring agent instead of N agents for N services

- Live service discovery w/o restarting/reloading

- Metrics are exported in background and stored in the cache. So it prevents high resources utilization.

#### Usage

Use aksentyev/postgres_exporter docker image to easy deploy the app.

Add your own metrics to the queries.yaml file. See *exporter/queries.yaml*.
