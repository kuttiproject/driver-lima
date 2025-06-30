package driverlima

import "github.com/kuttiproject/drivercore"

// QualifiedNetworkName returns a unique Network name for a cluster.
func (vd *Driver) QualifiedNetworkName(clustername string) string {
	panic("QualifiedNetworkName not implemented for lima driver")
}

// DeleteNetwork deletes the Network for a cluster.
func (vd *Driver) DeleteNetwork(clustername string) error {
	panic("DeleteNetwork not implemented for lima driver")
}

// NewNetwork creates a new Network for a cluster.
func (vd *Driver) NewNetwork(clustername string) (drivercore.Network, error) {
	panic("NewNetwork not implemented for lima driver")
}
