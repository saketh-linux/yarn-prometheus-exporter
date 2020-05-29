package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
        "os/exec"

	"github.com/prometheus/client_golang/prometheus"
)

type collectorApps struct {
	endpoint          *url.URL
	username          string
	password          string
	queueName         string
	allocatedMB       *prometheus.Desc
	allocatedVCores   *prometheus.Desc
	runningContainers *prometheus.Desc
	elapsedTime       *prometheus.Desc
}

const appMetricsNamespace = "yarn_apps"

func newAppsCollector(endpoint *url.URL, username string, password string, queueName string) *collectorApps {
	return &collectorApps{
		endpoint:          endpoint,
		username:          username,
		password:          password,
		queueName:         queueName,
		allocatedMB:       newFuncMetric(appMetricsNamespace, "allocated_memory", "Allocated memory", []string{"user", "name", "id", "state", "applicationType", "queue"}),
		allocatedVCores:   newFuncMetric(appMetricsNamespace, "allocatedVCores", "Allocated virtual cores", []string{"user", "name", "id", "state", "applicationType", "queue"}),
		runningContainers: newFuncMetric(appMetricsNamespace, "runningContainers", "Running containers", []string{"user", "name", "id", "state", "applicationType", "queue"}),
		elapsedTime:       newFuncMetric(appMetricsNamespace, "elapsedTime", "Elaspsed Time", []string{"user", "name", "id", "state", "applicationType", "queue","mail"}),
	}
}

func (c *collectorApps) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.allocatedMB
	ch <- c.allocatedVCores
	ch <- c.runningContainers
	ch <- c.elapsedTime
}

func (c *collectorApps) Collect(ch chan<- prometheus.Metric) {

	queues := strings.Split(c.queueName, ",")

	for _, yarn_queue := range queues {
		yarn_queue = strings.TrimSpace(yarn_queue)

		fmt.Println(yarn_queue)

		appData, err := fetchAppsMetrics(c.endpoint, c.username, c.password, yarn_queue)
		if err != nil {
			log.Println("Error while collecting apps metrics from YARN: " + err.Error())
		}

		appMetrics := appData["apps"]["app"]

		for _, value := range appMetrics {

			allocatedMB := value["allocatedMB"].(float64)
			allocatedVCores := value["allocatedVCores"].(float64)
			runningContainers := value["runningContainers"].(float64)
			elapsedTime := value["elapsedTime"].(float64)

			user := value["user"].(string)
			name := value["name"].(string)
			state := value["state"].(string)
			applicationType := value["applicationType"].(string)
			id := value["id"].(string)
                        out,_ := exec.Command("sh","-c",fmt.Sprintf("ldapsearch -LLL -h proxyldap -p 55395 -w ldap2015 -D uid=NISMaster,ou=people,dc=uhc,dc=com -b ou=people,dc=uhc,dc=com -s one '(uid=%s)'|grep mail|awk '{print $2}'",user)).Output()
                        mail := strings.TrimSpace(string(out))

			ch <- prometheus.MustNewConstMetric(c.allocatedMB, prometheus.GaugeValue, allocatedMB, user, name, id, state, applicationType, yarn_queue)
			ch <- prometheus.MustNewConstMetric(c.allocatedVCores, prometheus.GaugeValue, allocatedVCores, user, name, id, state, applicationType, yarn_queue)
			ch <- prometheus.MustNewConstMetric(c.runningContainers, prometheus.GaugeValue, runningContainers, user, name, id, state, applicationType, yarn_queue)
			ch <- prometheus.MustNewConstMetric(c.elapsedTime, prometheus.GaugeValue, elapsedTime, user, name, id, state, applicationType, yarn_queue,mail)
		}

	}

	return
}

func fetchAppsMetrics(u *url.URL, username string, password string, queueName string) (map[string]map[string][]map[string]interface{}, error) {

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}

	client := &http.Client{
		Transport: transCfg,
		Timeout:   100 * time.Second}

	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))

	q := req.URL.Query()
	q.Add("queue", queueName)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("unexpected HTTP status: " + string(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)

	// fmt.Println(string(body))

	if err != nil {
		return nil, err
	}

	var data map[string]map[string][]map[string]interface{}

	err = json.Unmarshal(body, &data)

	if err != nil {
		return nil, err
	}

	return data, nil
}
