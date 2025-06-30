package driverlima

import (
	"github.com/kuttiproject/drivercore"
)

func init() {
	// The lima packages require the following
	//limaVersion.Version = "v1.0.7"

	driver := &Driver{}

	drivercore.RegisterDriver(driverName, driver)
}
