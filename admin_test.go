package main

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	provider        = GCE
	zoneName        = "us-east-1"
	trafficZoneName = "us-west-1"
)

func TestNewAdminPanel(t *testing.T) {
	p := NewAdminPanel()
	assert.NotNil(t, p)
}

func TestAdminPanelAddZone(t *testing.T) {
	p := NewAdminPanel()
	assert.NotNil(t, p)

	z := NewZone(provider, "")
	z.Name = zoneName

	p.add(z)

	assert.True(t, len(p.AllZones()) > 0)
	assert.True(t, p.exists(zoneName))
}

func TestAdminPanelPingZone(t *testing.T) {
	p := NewAdminPanel()
	assert.NotNil(t, p)

	z := NewZone(provider, "")
	z.Name = zoneName
	z.Traffic[trafficZoneName] = 123456

	p.Ping(*z)

	assert.True(t, len(p.AllZones()) > 0)
	assert.True(t, p.exists(zoneName))
	assert.Equal(t, 123456, p.get(zoneName).Traffic[trafficZoneName])
}

func TestAdminPanelPingExistingZone(t *testing.T) {
	p := NewAdminPanel()
	assert.NotNil(t, p)

	z := NewZone(provider, "")
	z.Name = zoneName
	z.Traffic[trafficZoneName] = 123456

	p.Ping(*z)

	assert.True(t, len(p.AllZones()) > 0)
	assert.True(t, p.exists(zoneName))
	assert.Equal(t, 123456, p.get(zoneName).Traffic[trafficZoneName])

	z2 := NewZone(provider, "")
	z2.Name = zoneName
	z2.Traffic[trafficZoneName] = 4

	p.Ping(*z2)
	assert.True(t, len(p.AllZones()) > 0)
	assert.True(t, p.exists(zoneName))
	assert.Equal(t, 123460, p.get(zoneName).Traffic[trafficZoneName])
}

func TestAdminPanelEnableZone(t *testing.T) {
	p := NewAdminPanel()
	assert.NotNil(t, p)

	z := NewZone(provider, "")
	z.Name = zoneName
	z.Traffic[trafficZoneName] = 123456
	p.add(z)

	p.Enable(z.Name)

	assert.True(t, p.get(zoneName).Ready)

}

func TestAdminPanelDisableZone(t *testing.T) {
	p := NewAdminPanel()
	assert.NotNil(t, p)

	z := NewZone(provider, "")
	z.Name = zoneName
	z.Traffic[trafficZoneName] = 123456
	p.add(z)

	p.Disable(z.Name)

	assert.False(t, p.get(zoneName).Ready)

}

func TestAdminPanelEnablePingZone(t *testing.T) {
	p := NewAdminPanel()
	assert.NotNil(t, p)

	z := NewZone(provider, "")
	z.Name = zoneName
	z.Traffic[trafficZoneName] = 123456

	p.Ping(*z)

	p.Enable(z.Name)

	assert.True(t, p.get(zoneName).Ready)

	assert.True(t, len(p.AllZones()) > 0)
	assert.True(t, p.exists(zoneName))
	assert.Equal(t, 123456, p.get(zoneName).Traffic[trafficZoneName])

}
