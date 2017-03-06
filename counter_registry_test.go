package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCounterIncrement(t *testing.T) {
	zone := "us-east-1"
	cReg := NewCounterRegistry()
	cReg.Increment(zone)

	counter := cReg.Get(zone)

	assert.Equal(t, 1, counter)
}

func TestCounterIncrementMany(t *testing.T) {
	zone := "us-east-1"
	cReg := NewCounterRegistry()
	cReg.Increment(zone)

	// Add two more
	cReg.Increment(zone)
	cReg.Increment(zone)

	counter := cReg.Get(zone)

	assert.Equal(t, 3, counter)
}

func TestCounterIncrementJson(t *testing.T) {
	zone := "us-east-1"
	cReg := NewCounterRegistry()
	cReg.Increment(zone)

	// Add two more
	cReg.Increment(zone)
	cReg.Increment(zone)

	counter := cReg.Get(zone)

	assert.Equal(t, 3, counter)
	assert.Equal(t, "{\"us-east-1\":3}", string(cReg.ToJson()))

}

func TestCounterIncrementZonesInitial(t *testing.T) {
	// Data preparation
	zones := make(map[string]int64)
	zone := "us-east-1"
	zones[zone] = 12345

	// Create the registry
	cReg := NewCounterRegistry()

	// increment of many the zone
	cReg.IncrementZones(zones)

	// verify the counter
	counter := cReg.Get(zone)

	assert.Equal(t, 12345, counter)
}

func TestCounterIncrementZones(t *testing.T) {
	// Data preparation
	zones := make(map[string]int64)
	zone := "us-east-1"
	zones[zone] = 12345

	// Create the registry
	cReg := NewCounterRegistry()
	// increment of once the zone
	cReg.Increment(zone)

	// increment of many the zone
	cReg.IncrementZones(zones)

	// verify the counter
	counter := cReg.Get(zone)

	assert.Equal(t, 12345+1, counter)
}

func TestresetCounter(t *testing.T) {
	// Data preparation
	zones := make(map[string]int64)
	zone := "us-east-1"
	zones[zone] = 12345

	// Create the registry
	cReg := NewCounterRegistry()
	// increment of once the zone
	cReg.Increment(zone)

	// increment of many the zone
	cReg.IncrementZones(zones)

	// verify the counter
	counter := cReg.Get(zone)

	assert.Equal(t, 12345+1, counter)

	cReg.ResetCounter()

	counter = cReg.Get(zone)

	assert.Equal(t, 0, counter)

}
