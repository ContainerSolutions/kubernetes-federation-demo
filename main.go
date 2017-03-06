package main

import (
	"os"
)

var (
	serviceHost = ""
	servicePort = ""
)

func main() {
	var serviceConfig ServiceConfig
	var adminConfig AdminConfig
	provider := GCE

	apiConfig := ApiConfig{
		host:        getEnvOrElse("HOST", "0.0.0.0"),
		port:        getEnvOrElse("PORT", "8080"),
		serviceHost: os.Getenv("GEOSERVER_SERVICE_HOST"),
		servicePort: os.Getenv("GEOSERVER_SERVICE_PORT"),
		provider:    provider,
		zone:        NewZone(provider, ""),
		isAdmin:     getEnvValueBool("ADMIN"),
		isReady:     true,
		traffic:     NewCounterRegistry(),
	}

	if apiConfig.serviceHost != "" {
		apiConfig.zone.IpAddress.Ip = apiConfig.serviceHost
	}

	if apiConfig.servicePort != "" {
		apiConfig.zone.IpAddress.Port = apiConfig.servicePort
	}

	if apiConfig.zone.IpAddress.Port == "" {
		apiConfig.zone.IpAddress.Port = apiConfig.port
	}

	if apiConfig.isAdmin {
		adminConfig.federationIP = os.Getenv("FEDERATION_IP")
		adminConfig.clusters = os.Getenv("CLUSTERS")
		adminConfig.adminPanel = NewAdminPanel()

	} else {

		serviceConfig.adminHost = getEnvOrElse("REMOTE_IP", "")
		serviceConfig.adminPort = getEnvOrElse("REMOTE_PORT", "")
		serviceConfig.interval = getEnvOrElse("INTERVAL", "5")

		heart := NewHeartBeat(&apiConfig, &serviceConfig)
		heart.Start()

	}

	// ==============================================================

	api := NewAPI(&apiConfig, &serviceConfig, &adminConfig)
	api.Start()
}

func getEnvOrElse(envName, defaultValue string) string {
	varValue := os.Getenv(envName)
	if varValue == "" {
		return defaultValue
	}

	return varValue
}

func getEnvValueBool(envName string) bool {
	isAdminConfig := os.Getenv(envName)
	return isAdminConfig == "1"
}
