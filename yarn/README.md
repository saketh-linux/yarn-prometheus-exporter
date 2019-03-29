#YARN prometheus exporter
Export YARN metrics in Prometheus format.

##Build
Requires Go. Tested with Go 1.9+.

    go get
    go build -o yarn-prometheus-exporter .
##Run

The exporter can be configured using environment variables. These are the defaults:

    YARN_PROMETHEUS_LISTEN_ADDR=:9113
    YARN_PROMETHEUS_ENDPOINT_SCHEME=http
    YARN_PROMETHEUS_ENDPOINT_HOST=localhost
    YARN_PROMETHEUS_ENDPOINT_PORT=8088
    YARN_PROMETHEUS_ENDPOINT_PATH=ws/v1/cluster/metrics
    
##Run the exporter:

    ./yarn-prometheus-exporter

##The metrics can be scraped from:

    http://localhost:9113/metrics

##Run using docker:


    docker run -p 9113:9113 pbweb/yarn-prometheus-exporter