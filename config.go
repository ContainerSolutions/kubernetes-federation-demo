package main

type ApiConfig struct {
	// api config
	host string
	port string

	// service envs
	serviceHost string
	servicePort string

	// provider
	provider CloudProvider

	// zone the deployment belongs to
	zone *zone

	isAdmin bool

	isReady bool

	traffic map[string]int64
}

type ServiceConfig struct {

	// the host and port where the service sends the heartbeats to
	adminHost string
	adminPort string
	interval  string
}

type AdminConfig struct {
	adminPanel *AdminPanel

	// registry sevice (used only in case we are in amdin mode)
	registryService RegistryService
}
