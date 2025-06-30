package driverlima

const (
	driverName        = "lima"
	driverDescription = "Kutti driver for Lima"
)

// Driver implements the drivercore.Driver interface for Lima.
type Driver struct {
	limactlpath  string
	validated    bool
	status       string
	errormessage string
}

// Name returns "lima"
func (vd *Driver) Name() string {
	return driverName
}

// Description returns "Kutti driver for Lima"
func (vd *Driver) Description() string {
	return driverDescription
}

// UsesPerClusterNetworking returns false.
// This driver uses the Lima "user-v2" network for all clusters.
func (vd *Driver) UsesPerClusterNetworking() bool {
	return false
}

// UsesNATNetworking returns true.
// This driver uses port forwarding, though not user-defined.
func (vd *Driver) UsesNATNetworking() bool {
	return true
}

// Status returns the current status of the Driver.
func (vd *Driver) Status() string {
	vd.validate()
	return vd.status
}

// Error returns the last error reported by the Driver.
func (vd *Driver) Error() string {
	return vd.errormessage
}
