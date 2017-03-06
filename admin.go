package main

import (
	"fmt"
	"sync"
)

type AdminPanel struct {
	zones map[string]*zone
	mutex sync.RWMutex
}

func NewAdminPanel() *AdminPanel {

	return &AdminPanel{
		zones: make(map[string]*zone),
		mutex: sync.RWMutex{},
	}
}

func (adm *AdminPanel) Enable(zoneName string) {

	if !adm.exists(zoneName) {
		return
	}

	adm.mutex.Lock()
	defer adm.mutex.Unlock()
	adm.zones[zoneName].Ready = true

}

func (adm *AdminPanel) Disable(zoneName string) {

	if !adm.exists(zoneName) {
		return
	}

	adm.mutex.Lock()
	defer adm.mutex.Unlock()
	adm.zones[zoneName].Ready = false

}

func (adm *AdminPanel) add(z *zone) *zone {

	adm.mutex.Lock()
	defer adm.mutex.Unlock()
	adm.zones[z.Name] = z

	return z
}

func (adm *AdminPanel) get(zoneName string) *zone {

	z := new(zone)
	adm.mutex.RLock()
	defer adm.mutex.RUnlock()
	z = adm.zones[zoneName]
	return z
}

func (adm *AdminPanel) exists(zoneName string) bool {

	adm.mutex.RLock()
	defer adm.mutex.RUnlock()
	_, ok := adm.zones[zoneName]

	return ok
}

func (adm *AdminPanel) set(z *zone) {

	adm.mutex.Lock()
	defer adm.mutex.Unlock()
	adm.zones[z.Name] = z
}

func (adm *AdminPanel) Ping(z zone) {

	// check if the zone exists and if it does then increment the zone traffic
	if adm.exists(z.Name) {
		//fmt.Printf("zone exists: %q\n", *z)
		existing := adm.get(z.Name)
		fmt.Println("EXISTING:", existing)
		existing.IncrementZones(z.Traffic)
	} else {
		// if not add a new zone and set the current traffic counter coming from the parameter
		adm.add(&z)
	}
}

func (adm *AdminPanel) AllZones() []*zone {

	var zones []*zone
	adm.mutex.RLock()
	defer adm.mutex.RUnlock()
	for _, z := range adm.zones {
		zones = append(zones, z)
	}

	return zones

}
