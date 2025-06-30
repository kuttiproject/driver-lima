package driverlima_test

import (
	"testing"

	"github.com/kuttiproject/drivercore/drivercoretest"
	"github.com/kuttiproject/kuttilog"
)

const (
	TESTK8SVERSION = "1.33"
)

func TestDriverLima(t *testing.T) {
	kuttilog.SetLogLevel(kuttilog.MaxLevel())
	drivercoretest.TestDriver(t, "lima", TESTK8SVERSION)
}
