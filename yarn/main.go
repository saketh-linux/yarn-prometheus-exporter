package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr         string
	endpoint     *url.URL
	appsEndPoint *url.URL
	username     string
	password     string
	queue        string
)

func main() {
	loadEnv()

	c := newCollector(endpoint, username, password)
	prometheus.Register(c)

	appsCollector := newAppsCollector(appsEndPoint, username, password, queue)
	prometheus.Register(appsCollector)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}

func loadEnv() {
	addr = getEnvOr("YARN_PROMETHEUS_LISTEN_ADDR", ":9113")
	scheme := getEnvOr("YARN_RM_ENDPOINT_SCHEME", "http")
	host := getEnvOr("YARN_RM_ENDPOINT_HOST", "localhost")
	port := getEnvOr("YARN_RM_ENDPOINT_PORT", "8088")
	path := getEnvOr("YARN_RM_CLUSTER_PATH", "ws/v1/cluster/metrics")
	appsPath := getEnvOr("YARN_RM_APPS_PATH", "ws/v1/cluster/apps")
	user := getEnvOr("YARN_RM_USER_NAME", " ")
	pwd := getEnvOr("YARN_RM_PASSWORD", " ")
	rmQueue := getEnvOr("YARN_RM_QUEUE", " ")

	e, err := url.Parse(scheme + "://" + host + ":" + port + "/" + path)
	if err != nil {
		log.Fatal()
	}
	endpoint = e

	appsMetricsEndpoint, err := url.Parse(scheme + "://" + host + ":" + port + "/" + appsPath)
	if err != nil {
		log.Fatal()
	}
	appsEndPoint = appsMetricsEndpoint

	username = user
	password = pwd
	queue = rmQueue
}

func getEnvOr(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}
