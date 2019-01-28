package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type collector struct {
	endpoint              *url.URL
	username              string
	password              string
	up                    *prometheus.Desc
	applicationsSubmitted *prometheus.Desc
	applicationsCompleted *prometheus.Desc
	applicationsPending   *prometheus.Desc
	applicationsRunning   *prometheus.Desc
	applicationsFailed    *prometheus.Desc
	applicationsKilled    *prometheus.Desc
	memoryReserved        *prometheus.Desc
	memoryAvailable       *prometheus.Desc
	memoryAllocated       *prometheus.Desc
	memoryTotal           *prometheus.Desc
	virtualCoresReserved  *prometheus.Desc
	virtualCoresAvailable *prometheus.Desc
	virtualCoresAllocated *prometheus.Desc
	virtualCoresTotal     *prometheus.Desc
	containersAllocated   *prometheus.Desc
	containersReserved    *prometheus.Desc
	containersPending     *prometheus.Desc
	nodesTotal            *prometheus.Desc
	nodesLost             *prometheus.Desc
	nodesUnhealthy        *prometheus.Desc
	nodesDecommissioned   *prometheus.Desc
	nodesDecommissioning  *prometheus.Desc
	nodesRebooted         *prometheus.Desc
	nodesActive           *prometheus.Desc
	scrapeFailures        *prometheus.Desc
	failureCount          int
}

const metricsNamespace = "yarn"

func newCollector(endpoint *url.URL, username string, password string) *collector {
	return &collector{
		endpoint:              endpoint,
		username:              username,
		password:              password,
		up:                    newFuncMetric(metricsNamespace, "up", "Able to contact YARN", nil),
		applicationsSubmitted: newFuncMetric(metricsNamespace, "applications_submitted", "Total applications submitted", nil),
		applicationsCompleted: newFuncMetric(metricsNamespace, "applications_completed", "Total applications completed", nil),
		applicationsPending:   newFuncMetric(metricsNamespace, "applications_pending", "Applications pending", nil),
		applicationsRunning:   newFuncMetric(metricsNamespace, "applications_running", "Applications running", nil),
		applicationsFailed:    newFuncMetric(metricsNamespace, "applications_failed", "Total application failed", nil),
		applicationsKilled:    newFuncMetric(metricsNamespace, "applications_killed", "Total application killed", nil),
		memoryReserved:        newFuncMetric(metricsNamespace, "memory_reserved", "Memory reserved", nil),
		memoryAvailable:       newFuncMetric(metricsNamespace, "memory_available", "Memory available", nil),
		memoryAllocated:       newFuncMetric(metricsNamespace, "memory_allocated", "Memory allocated", nil),
		memoryTotal:           newFuncMetric(metricsNamespace, "memory_total", "Total memory", nil),
		virtualCoresReserved:  newFuncMetric(metricsNamespace, "virtual_cores_reserved", "Virtual cores reserved", nil),
		virtualCoresAvailable: newFuncMetric(metricsNamespace, "virtual_cores_available", "Virtual cores available", nil),
		virtualCoresAllocated: newFuncMetric(metricsNamespace, "virtual_cores_allocated", "Virtual cores allocated", nil),
		virtualCoresTotal:     newFuncMetric(metricsNamespace, "virtual_cores_total", "Total virtual cores", nil),
		containersAllocated:   newFuncMetric(metricsNamespace, "containers_allocated", "Containers allocated", nil),
		containersReserved:    newFuncMetric(metricsNamespace, "containers_reserved", "Containers reserved", nil),
		containersPending:     newFuncMetric(metricsNamespace, "containers_pending", "Containers pending", nil),
		nodesTotal:            newFuncMetric(metricsNamespace, "nodes_total", "Nodes total", nil),
		nodesLost:             newFuncMetric(metricsNamespace, "nodes_lost", "Nodes lost", nil),
		nodesUnhealthy:        newFuncMetric(metricsNamespace, "nodes_unhealthy", "Nodes unhealthy", nil),
		nodesDecommissioned:   newFuncMetric(metricsNamespace, "nodes_decommissioned", "Nodes decommissioned", nil),
		nodesDecommissioning:  newFuncMetric(metricsNamespace, "nodes_decommissioning", "Nodes decommissioning", nil),
		nodesRebooted:         newFuncMetric(metricsNamespace, "nodes_rebooted", "Nodes rebooted", nil),
		nodesActive:           newFuncMetric(metricsNamespace, "nodes_active", "Nodes active", nil),
		scrapeFailures:        newFuncMetric(metricsNamespace, "scrape_failures_total", "Number of errors while scraping YARN metrics", nil),
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.applicationsSubmitted
	ch <- c.applicationsCompleted
	ch <- c.applicationsPending
	ch <- c.applicationsRunning
	ch <- c.applicationsFailed
	ch <- c.applicationsKilled
	ch <- c.memoryReserved
	ch <- c.memoryAvailable
	ch <- c.memoryAllocated
	ch <- c.memoryTotal
	ch <- c.virtualCoresReserved
	ch <- c.virtualCoresAvailable
	ch <- c.virtualCoresAllocated
	ch <- c.virtualCoresTotal
	ch <- c.containersAllocated
	ch <- c.containersReserved
	ch <- c.containersPending
	ch <- c.nodesTotal
	ch <- c.nodesLost
	ch <- c.nodesUnhealthy
	ch <- c.nodesDecommissioned
	ch <- c.nodesDecommissioning
	ch <- c.nodesRebooted
	ch <- c.nodesActive
	ch <- c.scrapeFailures
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	up := 1.0

	data, err := fetch(c.endpoint, c.username, c.password)
	if err != nil {
		up = 0.0
		c.failureCount++

		log.Println("Error while collecting data from YARN: " + err.Error())
	}

	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, up)
	ch <- prometheus.MustNewConstMetric(c.scrapeFailures, prometheus.CounterValue, float64(c.failureCount))

	if up == 0.0 {
		return
	}

	metrics := data["clusterMetrics"]

	ch <- prometheus.MustNewConstMetric(c.applicationsSubmitted, prometheus.CounterValue, metrics["appsSubmitted"])
	ch <- prometheus.MustNewConstMetric(c.applicationsCompleted, prometheus.CounterValue, metrics["appsCompleted"])
	ch <- prometheus.MustNewConstMetric(c.applicationsPending, prometheus.GaugeValue, metrics["appsPending"])
	ch <- prometheus.MustNewConstMetric(c.applicationsRunning, prometheus.GaugeValue, metrics["appsRunning"])
	ch <- prometheus.MustNewConstMetric(c.applicationsFailed, prometheus.CounterValue, metrics["appsFailed"])
	ch <- prometheus.MustNewConstMetric(c.applicationsKilled, prometheus.CounterValue, metrics["appsKilled"])
	ch <- prometheus.MustNewConstMetric(c.memoryReserved, prometheus.GaugeValue, metrics["reservedMB"])
	ch <- prometheus.MustNewConstMetric(c.memoryAvailable, prometheus.GaugeValue, metrics["availableMB"])
	ch <- prometheus.MustNewConstMetric(c.memoryAllocated, prometheus.GaugeValue, metrics["allocatedMB"])
	ch <- prometheus.MustNewConstMetric(c.memoryTotal, prometheus.GaugeValue, metrics["totalMB"])
	ch <- prometheus.MustNewConstMetric(c.virtualCoresReserved, prometheus.GaugeValue, metrics["reservedVirtualCores"])
	ch <- prometheus.MustNewConstMetric(c.virtualCoresAvailable, prometheus.GaugeValue, metrics["availableVirtualCores"])
	ch <- prometheus.MustNewConstMetric(c.virtualCoresAllocated, prometheus.GaugeValue, metrics["allocatedVirtualCores"])
	ch <- prometheus.MustNewConstMetric(c.virtualCoresTotal, prometheus.GaugeValue, metrics["totalVirtualCores"])
	ch <- prometheus.MustNewConstMetric(c.containersAllocated, prometheus.GaugeValue, metrics["containersAllocated"])
	ch <- prometheus.MustNewConstMetric(c.containersReserved, prometheus.GaugeValue, metrics["containersReserved"])
	ch <- prometheus.MustNewConstMetric(c.containersPending, prometheus.GaugeValue, metrics["containersPending"])
	ch <- prometheus.MustNewConstMetric(c.nodesTotal, prometheus.GaugeValue, metrics["totalNodes"])
	ch <- prometheus.MustNewConstMetric(c.nodesLost, prometheus.GaugeValue, metrics["lostNodes"])
	ch <- prometheus.MustNewConstMetric(c.nodesUnhealthy, prometheus.GaugeValue, metrics["unhealthyNodes"])
	ch <- prometheus.MustNewConstMetric(c.nodesDecommissioned, prometheus.GaugeValue, metrics["decommissionedNodes"])
	ch <- prometheus.MustNewConstMetric(c.nodesDecommissioning, prometheus.GaugeValue, metrics["decommissioningNodes"])
	ch <- prometheus.MustNewConstMetric(c.nodesRebooted, prometheus.GaugeValue, metrics["rebootedNodes"])
	ch <- prometheus.MustNewConstMetric(c.nodesActive, prometheus.GaugeValue, metrics["activeNodes"])

	return
}

func fetch(u *url.URL, username string, password string) (map[string]map[string]float64, error) {

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}

	client := &http.Client{
		Transport: transCfg,
		Timeout:   100 * time.Second}

	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("unexpected HTTP status: " + string(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var data map[string]map[string]float64

	err = json.Unmarshal(body, &data)

	if err != nil {
		return nil, err
	}

	return data, nil
}
