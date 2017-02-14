package main

import (
	"sync"
)

type AdminPanel struct {
	zones map[string]*zone
	mutex sync.Mutex
}

func NewAdminPanel() *AdminPanel {

	return &AdminPanel{
		zones: make(map[string]*zone),
		mutex: sync.Mutex{},
	}
}

func (s *AdminPanel) get(zoneName string) *zone {

	z := new(zone)
	s.mutex.Lock()
	z = s.zones[zoneName]
	s.mutex.Unlock()
	return z
}

func (s *AdminPanel) ping(z *zone) {

	s.mutex.Lock()
	s.zones[z.Name] = z
	s.mutex.Unlock()

}

func (s *AdminPanel) getAll() []*zone {

	var zones []*zone
	s.mutex.Lock()
	for _, z := range s.zones {
		zones = append(zones, z)
	}
	s.mutex.Unlock()
	return zones

}
