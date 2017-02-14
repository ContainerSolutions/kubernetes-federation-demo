// This file allows the geo server to report its status to a central admin service
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var urlFormatWithPort = "http://%s:%s/ping"
var urlFormat = "http://%s/ping"

type heartbeat struct {
	apiConfig     *ApiConfig
	serviceConfig *ServiceConfig
	interval      int
	client        http.Client
}

func NewHeartBeat(apiConfig *ApiConfig, serviceConfig *ServiceConfig) *heartbeat {

	intervalNum := 10

	if intNum, err := strconv.Atoi(serviceConfig.interval); err == nil {
		intervalNum = intNum
	}

	return &heartbeat{
		apiConfig:     apiConfig,
		serviceConfig: serviceConfig,
		interval:      intervalNum,
		client:        http.Client{},
	}

}

func (h *heartbeat) Start() {

	if h.serviceConfig.adminHost == "" {
		log.Println("REMOTE_IP env variable is not set. Not sending heartbeats.")
		return
	}

	// build the url
	var url string
	if h.serviceConfig.adminPort == "" {
		url = fmt.Sprintf(urlFormat, h.serviceConfig.adminHost)
	} else {
		url = fmt.Sprintf(urlFormatWithPort, h.serviceConfig.adminHost, h.serviceConfig.adminPort)
	}

	// start the ticker
	ticker := time.NewTicker(time.Duration(h.interval) * time.Second)

	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				h.ping(url)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (h *heartbeat) ping(url string) {
	// add the traffic information
	// traffic := make(map[string]int64)
	// for zoneId, counter := range h.apiConfig.traffic {
	// 	traffic[zoneId] = counter
	// 	log.Printf("Sending traffic for zone: %s => %d\n", zoneId, counter)
	// }

	h.apiConfig.zone.Traffic = h.apiConfig.traffic

	if data := h.apiConfig.zone.toJson(); len(data) > 0 {

		// set the timestamp
		h.apiConfig.zone.Timestamp = time.Now().UTC()

		req, err := http.NewRequest("POST", url, bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")

		resp, err := h.client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			//parse JSON and set the ready flag
			log.Println("Response Status:", url, resp.Status)

			var z zone
			if err := json.NewDecoder(resp.Body).Decode(&z); err != nil {
				log.Println(err)
				return
			}

			// set the readyness based on the response coming from the admin
			h.apiConfig.isReady = z.Ready
			log.Printf("Ping response: %s\n", z)

		} else {
			log.Println("Could not send heartbeat", err)
		}
		return
	}
	log.Println("Error serialize VM information to JSON")

}
