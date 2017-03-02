package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	caFile      = "/etc/federation/ca.crt"
	serverCert  = "/etc/federation/server.crt"
	serverKey   = "/etc/federation/server.key"
	clusterJSON = `{"kind":"Cluster","apiVersion":"federation/v1beta1","metadata":{"name":"%s"},"spec":{"serverAddressByClientCIDRs":[{"clientCIDR":"0.0.0.0/0","serverAddress":"https://%s"}],"secretRef":{"name":"%s"}},"status":{}}`
)

type federationManager struct {
	IP         string
	caFile     string
	serverCert string
	serverKey  string
	client     *http.Client
}

type ClustersResponse struct {
	Entries []ClusterMeta `json:"items"`
}

type ClusterMeta struct {
	Meta ClusterData `json:"metadata"`
}

type ClusterData struct {
	Name string `json:"name"`
}

// NewFederationManager: Creates a new federation manager object which allows to add or remove clusters to/from the federation API
func NewFederationManager(ip string) *federationManager {

	// Load federation client certs
	cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		log.Fatal(err)
	}

	// Load the federation CA cert
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup the HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return &federationManager{
		client: &http.Client{Transport: transport},
		IP:     ip,
	}

}

func (f *federationManager) AllClusters() ClustersResponse {
	var clusterResp ClustersResponse
	allReq, err := http.NewRequest("GET", fmt.Sprintf("https://%s/apis/federation/v1beta1/clusters", f.IP), nil)
	if err != nil {
		log.Println("Error creating GET request for all clusters")
		return clusterResp
	}

	res, err := f.client.Do(allReq)
	if err != nil {
		log.Println("Error executing GET for all cluster", err)
		return clusterResp
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&clusterResp)
	if err != nil {
		log.Println("Error parsing all clusters response:", err)
	}

	log.Println(clusterResp)

	return clusterResp
}

func (f *federationManager) RemoveCluster(cluster string) bool {
	deleteReq, err := http.NewRequest("DELETE", fmt.Sprintf("https://%s/apis/federation/v1beta1/clusters/%s", f.IP, cluster), nil)
	if err != nil {
		log.Println("Error creating DELETE request for cluster", cluster)
		return false
	}

	res, err := f.client.Do(deleteReq)
	if err != nil {
		log.Println("Error executing DELETE for cluster", cluster, err)
		return false
	}
	defer res.Body.Close()

	return res.StatusCode == 200
}

func (f *federationManager) AddCluster(cluster, ip string) bool {
	// format the JSON request body
	jsonBody := fmt.Sprintf(clusterJSON, cluster, ip, cluster)

	// create buffer
	buffer := bytes.NewBuffer([]byte(jsonBody))

	// Post it
	addReq, err := http.NewRequest("POST", fmt.Sprintf("https://%s/apis/federation/v1beta1/clusters", f.IP), buffer)
	addReq.Header.Add("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		log.Println("Error creating POST request for adding cluster", cluster)
		return false
	}

	res, err := f.client.Do(addReq)
	if err != nil {
		log.Println("Error executing POST for adding cluster", cluster, err)
		return false
	}
	defer res.Body.Close()

	return res.StatusCode == 200
}
