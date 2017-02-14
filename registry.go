package main

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"sync"
)

type RegistryService interface {
	Add(host string)                        // Add an endpoint to our registry
	Delete(host string)                     // Remove an endpoint to our registry
	Get(host string) *httputil.ReverseProxy // Return the endpoint list for the given service name/version
	Exists(host string) bool
}

type registry struct {
	entries map[string]*httputil.ReverseProxy
	mutex   sync.Mutex
}

func NewRegistryService() RegistryService {
	r := registry{
		mutex:   sync.Mutex{},
		entries: make(map[string]*httputil.ReverseProxy),
	}
	return r
}

func (r registry) Add(host string) {

	r.mutex.Lock()
	defer r.mutex.Unlock()
	urlString := fmt.Sprintf("http://%s", host)
	log.Printf("Parsing url: %s\n", urlString)
	remote, err := url.Parse(urlString)
	if err != nil {
		log.Println(err)
		return
	}
	r.entries[host] = httputil.NewSingleHostReverseProxy(remote)
	log.Println("Added reverse proxy", host)

}

func (r registry) Exists(host string) bool {

	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.entries[host] != nil

}

func (r registry) Delete(host string) {

	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.entries, "host")
}

func (r registry) Get(host string) *httputil.ReverseProxy {

	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.entries[host]
}
