package main

import (
	"encoding/json"
	"sync"
)

type CounterRegistryService interface {
	Increment(zone string)                 // Increment the counter for a specific zone
	IncrementZones(zones map[string]int64) // Increment the counter for a specific zone
	ToJson() []byte                        // Returns a json representation of the counter
	Get(zone string) int64
	AllZones() map[string]int64 // returns the zones map
	ResetCounter()
}

type counterRegistry struct {
	Zones map[string]int64
	mutex sync.Mutex
}

func NewCounterRegistry() *counterRegistry {
	//	service := new(CounterRegistryService)
	registry := &counterRegistry{
		mutex: sync.Mutex{},
		Zones: make(map[string]int64),
	}
	//service = registry
	return registry
}

func (r *counterRegistry) AllZones() map[string]int64 {

	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.Zones
}

func (r *counterRegistry) ResetCounter() {

	r.mutex.Lock()
	defer r.mutex.Unlock()
	for key, _ := range r.Zones {
		r.Zones[key] = 0
	}
}

func (r *counterRegistry) Increment(zone string) {

	r.mutex.Lock()
	defer r.mutex.Unlock()
	// if the entry exists, then increment its value
	if _, ok := r.Zones[zone]; ok {
		r.Zones[zone]++
	} else {
		// otherwise create a new entry
		r.Zones[zone] = 1
	}
}

// This is used in the admin to increment the counters for each zones,
// arriving from the heartbeats
func (r *counterRegistry) IncrementZones(zones map[string]int64) {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// loop through the data coming from the parameter
	for zone, counter := range zones {

		if _, ok := r.Zones[zone]; ok {
			r.Zones[zone] += counter
		} else {
			// otherwise create a new entry
			r.Zones[zone] = counter
		}

	}
}

func (r *counterRegistry) Get(zone string) int64 {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.Zones[zone]

}

func (r *counterRegistry) ToJson() []byte {
	var data []byte
	r.mutex.Lock()
	defer r.mutex.Unlock()
	data, _ = json.Marshal(r.Zones)
	return data
}
