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
	addr     string
	endpoint *url.URL
	username string
	password string
)

func main() {
	loadEnv()

	c := newCollector(endpoint, username, password)
	prometheus.Register(c)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}

func loadEnv() {
	addr = getEnvOr("YARN_PROMETHEUS_LISTEN_ADDR", ":9113")
	scheme := getEnvOr("YARN_RM_ENDPOINT_SCHEME", "http")
	host := getEnvOr("YARN_RM_ENDPOINT_HOST", "localhost")
	port := getEnvOr("YARN_RM_ENDPOINT_PORT", "8088")
	path := getEnvOr("YARN_RM_PATH", "ws/v1/cluster/metrics")
	user := getEnvOr("YARN_RM_USER_NAME", " ")
	pwd := getEnvOr("YARN_RM_PASSWORD", " ")

	e, err := url.Parse(scheme + "://" + host + ":" + port + "/" + path)
	if err != nil {
		log.Fatal()
	}

	endpoint = e
	username = user
	password = pwd
}

func getEnvOr(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}
