package main

import (
	"encoding/base64"

	"github.com/prometheus/client_golang/prometheus"
)

func newFuncMetric(metricNameSpace string, metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(metricNameSpace, "", metricName), docString, labels, nil)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
