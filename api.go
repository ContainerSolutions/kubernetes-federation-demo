package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

// Defines the currently supported cloud providers
// TODO: AWS
const (
	AWS CloudProvider = iota + 1
	GCE
)

var (
	openBrace  = byte('(')
	closeBrace = byte(')')
)

// Cloud Provider ENUM
type CloudProvider uint8

type cluster struct {
	Name   string
	IP     string
	Joined bool
}

type clusters []cluster

// ==============================
type templateData struct {
	RemoteHost string
}

// API Service
type api struct {
	config              *ApiConfig
	serviceConfig       *ServiceConfig
	adminConfig         *AdminConfig
	templateData        templateData
	datacenters         map[string]bool
	traffic             map[string]bool
	mutex               sync.Mutex
	mutexTraffic        sync.Mutex
	mutexTrafficCounter sync.Mutex
	federation          *federationManager
	clusters            clusters
}

// Create a new API struct
func NewAPI(apiConfig *ApiConfig, serviceConfig *ServiceConfig, adminConfig *AdminConfig) *api {

	remoteHost := serviceConfig.adminHost
	if serviceConfig != nil && serviceConfig.adminPort != "" {
		remoteHost = fmt.Sprintf("%s:%s", serviceConfig.adminHost, serviceConfig.adminPort)
	}
	if remoteHost == "" {
		remoteHost = "localhost:8081" // for local development
	}

	api := api{
		config:              apiConfig,
		serviceConfig:       serviceConfig,
		adminConfig:         adminConfig,
		datacenters:         make(map[string]bool),
		traffic:             make(map[string]bool),
		mutex:               sync.Mutex{},
		mutexTraffic:        sync.Mutex{},
		mutexTrafficCounter: sync.Mutex{},
	}

	// set the federation manager if the FEDERATION_IP
	if adminConfig.federationIP != "" {
		log.Println("FEDERATION_IP:", adminConfig.federationIP)
		api.federation = NewFederationManager(adminConfig.federationIP)
	}
	if adminConfig.clusters != "" {
		api.clusters = parseClusters(adminConfig.clusters)
	}

	return &api
}

// Start the API: both the admin and the services
// the difference is set by an environment variable
func (api *api) Start() {

	http.HandleFunc("/favicon.ico", api.faviconHandlerFunc)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// used for signaling that the conatiner is up and running
	http.HandleFunc("/live", api.liveHandlerFunc)

	// used for signaling that the conatiner is ready to receive requests
	http.HandleFunc("/ready", api.readyHandlerFunc)

	if api.config.isAdmin {
		// Add zone handler
		http.HandleFunc("/ping", api.pingHandlerFunc)

		http.HandleFunc("/", api.zoneIndexHandlerFunc)

		http.HandleFunc("/services", api.adminHandlerFunc)

		http.HandleFunc("/disable", api.adminDisableHandlerFunc)

		http.HandleFunc("/enable", api.adminEnableHandlerFunc)

		http.HandleFunc("/startTraffic", api.startTraffic)

		http.HandleFunc("/stopTraffic", api.stopTraffic)

		http.HandleFunc("/trafficSourceActive", api.trafficSourceActive)

		// simple GET entry points
		http.HandleFunc("/federation/clusters", api.clustersHandlerFunc)
		http.HandleFunc("/federation/add", api.federationAddHandlerFunc)
		http.HandleFunc("/federation/remove", api.federationRemoveHandlerFunc)

	} else {

		indexHandler := http.HandlerFunc(api.indexHandlerFunc)

		// Add zone handler
		http.Handle("/", api.requestsMiddleware(indexHandler))

		// Service to get the data
		http.HandleFunc("/location", api.zoneHandlerFunc)

		// disable the readyness
		http.HandleFunc("/disable", api.disableHandlerFunc)

		// enable the readyness
		http.HandleFunc("/enable", api.enableHandlerFunc)

		// kill ther app
		http.HandleFunc("/kill", api.killHandlerFunc)

		// default zone ready to true
		if api.config.zone != nil {
			api.config.zone.Ready = true
		}

	}

	// start the HTTP server
	portHost := fmt.Sprintf("%s:%s", api.config.host, api.config.port)

	// show the vars
	if api.config.isAdmin {
		log.Println("Admin mode - HTTP listening on:", portHost, "for provider", api.config.provider)
	} else {
		log.Println("HTTP listening on:", portHost, "for provider", api.config.provider)
	}

	for _, e := range os.Environ() {
		log.Println(e)
	}

	// socket listening
	log.Fatal(http.ListenAndServe(portHost, nil))
}

// =============================================================
// Middleware for request counter
func (api *api) requestsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		trafficSourceZone := r.Header.Get("X-Origin-Region")
		if trafficSourceZone != "" {
			// if the key already exists, then increment the counter
			api.mutexTrafficCounter.Lock()
			defer api.mutexTrafficCounter.Unlock()
			if currentCounter, ok := api.config.traffic[trafficSourceZone]; ok {
				api.config.traffic[trafficSourceZone] = currentCounter + 1
				log.Printf("Rate counter incremented for zone %s: %s\n", trafficSourceZone, api.config.traffic[trafficSourceZone])
			} else {
				// create the entry
				api.config.traffic[trafficSourceZone] = 1
				log.Printf("Rate counter created for zone %s: %s\n", trafficSourceZone, api.config.traffic[trafficSourceZone])
			}
		}

		next.ServeHTTP(w, r)
	})
}

// ##### ADMIN
// =============================================================

func (api *api) clustersHandlerFunc(w http.ResponseWriter, r *http.Request) {

	all := api.federation.AllClusters()

	for _, federatedCluster := range all.Entries {

		for index, staticCluster := range api.clusters {
			if federatedCluster.Meta.Name == staticCluster.Name {
				log.Println("Found federated cluster:", federatedCluster.Meta.Name)
				api.clusters[index].Joined = true
				break
			} else {
				api.clusters[index].Joined = false
			}
		}
	}

	data := api.clusters.toJson()
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

func (api *api) federationAddHandlerFunc(w http.ResponseWriter, r *http.Request) {
	clusterName := r.URL.Query().Get("name")

	if clusterName != "" {
		// retrieve the IP
		ip := api.clusters.findIP(clusterName)
		if ip != "" && api.federation != nil {
			if ok := api.federation.AddCluster(clusterName, ip); ok {
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}

	w.WriteHeader(http.StatusInternalServerError)

}

func (api *api) federationRemoveHandlerFunc(w http.ResponseWriter, r *http.Request) {
	clusterName := r.URL.Query().Get("name")

	if clusterName != "" {
		// retrieve the IP
		ip := api.clusters.findIP(clusterName)
		if ip != "" && api.federation != nil {
			if ok := api.federation.RemoveCluster(clusterName); ok {
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}

	w.WriteHeader(http.StatusInternalServerError)

}

func (api *api) startTraffic(w http.ResponseWriter, r *http.Request) {
	dc := r.URL.Query().Get("dc")

	if dc != "" {
		api.mutexTraffic.Lock()
		api.traffic[dc] = true
		api.mutexTraffic.Unlock()
		log.Println("Traffic from", dc, "enabled.")
	}
}

func (api *api) stopTraffic(w http.ResponseWriter, r *http.Request) {
	dc := r.URL.Query().Get("dc")

	if dc != "" {
		api.mutexTraffic.Lock()
		api.traffic[dc] = true
		api.mutexTraffic.Unlock()
		log.Println("Traffic from", dc, "disabled.")
	}
}

func (api *api) trafficSourceActive(w http.ResponseWriter, r *http.Request) {
	dc := r.URL.Query().Get("dc")

	if dc != "" {
		api.mutexTraffic.Lock()
		active, ok := api.traffic[dc]
		api.mutexTraffic.Unlock()
		if ok {
			fmt.Fprintf(w, "%t", active)
			return
		}
	}
	fmt.Fprintf(w, "%t", true)
}

func (api *api) adminHandlerFunc(w http.ResponseWriter, r *http.Request) {

	vms := api.adminConfig.adminPanel.getAll()
	data, _ := json.Marshal(&vms)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (api *api) pingHandlerFunc(w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		http.Error(w, "Missing request body", 400)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var z zone
	err := decoder.Decode(&z)

	if err != nil {
		log.Println("Could not deserialize zone", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if z.Name != "" {
		if value, ok := api.datacenters[z.Name]; ok {
			log.Println("Datacenter:", z.Name, value)
			z.Ready = value
		}
	}

	// retrieve the existing traffic for the zone
	existingZone := api.adminConfig.adminPanel.get(z.Name)
	if existingZone != nil {

		// loop through all source of traffic the zone received
		for trafficSourceZoneName, newTrafficCounter := range z.Traffic {
			existingZone.Traffic[trafficSourceZoneName] = newTrafficCounter
		}
	}

	if existingZone == nil {
		existingZone = &z
		log.Printf("NOT Found: %s\n", existingZone)
	}

	//log.Println("Zone", existingZone.Name, "updated traffic:", existingZone.Traffic)

	// update the in memory info
	api.adminConfig.adminPanel.ping(existingZone)

	data := existingZone.toJson()
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

func (api *api) adminDisableHandlerFunc(w http.ResponseWriter, r *http.Request) {
	dc := r.URL.Query().Get("dc")

	if dc != "" {
		api.mutex.Lock()
		api.datacenters[dc] = false
		api.mutex.Unlock()
		log.Println("Datacenter", dc, "disabled.")
	}
}

func (api *api) adminEnableHandlerFunc(w http.ResponseWriter, r *http.Request) {
	dc := r.URL.Query().Get("dc")

	if dc != "" {
		api.mutex.Lock()
		api.datacenters[dc] = true
		api.mutex.Unlock()
		log.Println("Datacenter", dc, "enabled.")
	}
}

// =====================================================================
// ##### END ADMIN

func (api *api) indexHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", api.config.zone.Name)
}

func (api *api) zoneIndexHandlerFunc(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("Error parsing admin template files ", err)
	}
	t.Execute(w, nil)
}

func (api *api) faviconHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func (api *api) killHandlerFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("Used KILL SWITCH")
	os.Exit(1)
}

// checks whether the service is up and runnning
func (api *api) liveHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// check if the service is ready to serve requests
func (api *api) readyHandlerFunc(w http.ResponseWriter, r *http.Request) {
	api.writeReadyness(w)
}

// Enable the service
func (api *api) enableHandlerFunc(w http.ResponseWriter, r *http.Request) {
	api.mutex.Lock()
	api.config.isReady = true
	api.mutex.Unlock()
	log.Println("Service enabled")
	w.WriteHeader(http.StatusOK)
}

// Disable the service
func (api *api) disableHandlerFunc(w http.ResponseWriter, r *http.Request) {
	api.mutex.Lock()
	api.config.isReady = false
	api.mutex.Unlock()
	log.Println("Service disabled")
	w.WriteHeader(http.StatusOK)
}

// Returns the datacenter zone information about the running process
func (api *api) zoneHandlerFunc(w http.ResponseWriter, r *http.Request) {
	api.mutex.Lock()
	defer api.mutex.Unlock()

	if !api.config.isReady {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// serve back the VM info and the client remote IP info
	remoteIp := getIPAdress(r)
	host, _, err := net.SplitHostPort(remoteIp)
	if err != nil {
		log.Println("Error splitting remote IP", err)
		host = remoteIp
	}
	if host != "" {
		coord := getIPCoordinates(host)
		api.config.zone.ClientIpAddress = &coord
		log.Printf("Got client IP coordinates: %s\n", api.config.zone.ClientIpAddress)
	}

	// serve the data back
	if data := api.config.zone.toJson(); len(data) > 0 {
		w.Write(data)
		return
	}
	log.Println("Error serialize VM information to JSON")

	// otherwise return error
	w.WriteHeader(http.StatusServiceUnavailable)

}

func (api *api) writeReadyness(w http.ResponseWriter) {

	api.mutex.Lock()
	if api.config.isReady {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	api.mutex.Unlock()

}

func remove(slice []*zone, s int) []*zone {
	return append(slice[:s], slice[s+1:]...)
}

func parseClusters(clustersEnv string) clusters {
	var clusters []cluster
	if clustersEnv == "" {
		return clusters
	}
	log.Println("Parsing clusters:", clustersEnv)

	r := csv.NewReader(strings.NewReader(clustersEnv))

	for {
		records, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}

		for _, record := range records {
			tokens := strings.Split(record, "=")
			if len(tokens) > 0 {
				clusters = append(clusters, cluster{Name: strings.Replace(tokens[0], "gce", "cluster", 1), IP: tokens[1]})
			}
		}

	}

	log.Println("Parsed clusters", clusters)

	return clusters
}

func (c clusters) toJson() []byte {
	var data []byte

	jsonData, err := json.Marshal(c)
	if err != nil {
		return data
	}

	return jsonData
}

func (c clusters) findIP(name string) string {
	for _, cluster := range c {
		if cluster.Name == name {
			return cluster.IP
		}
	}
	return ""
}
