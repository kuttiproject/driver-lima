package driverlima

import (
	"fmt"
	"os"

	"github.com/kuttiproject/drivercore"
	"github.com/pkg/errors"
)

// QualifiedMachineName returns a unique name for a Machine in a cluster.
// This name is usually used internally by a Driver.
func (vd *Driver) QualifiedMachineName(machinename string, clustername string) string {
	return clustername + "-" + machinename
}

// GetMachine returns a Machine in a cluster.
func (vd *Driver) GetMachine(machinename string, clustername string) (drivercore.Machine, error) {
	err := vd.validate()
	if err != nil {
		return nil, err
	}

	return &Machine{
		driver:      vd,
		name:        machinename,
		clustername: clustername,
	}, nil
}

// DeleteMachine deletes a Machine in a cluster.
func (vd *Driver) DeleteMachine(machinename string, clustername string) error {
	err := vd.validate()
	if err != nil {
		return err
	}

	machinefile, err := machineFilePath(vd.QualifiedMachineName(machinename, clustername))
	if err != nil {
		return errors.Wrap(err, "machine file not accessible")
	}

	limactlparams := []string{
		"rm",
		vd.QualifiedMachineName(machinename, clustername),
	}

	_, err = vd.runwithresults(limactlparams...)
	if err != nil {
		return errors.Wrap(err, "could not delete lima vm")
	}

	err = os.Remove(machinefile)
	if err != nil {
		return errors.Wrap(err, "machine file not deleted")
	}

	return nil
}

// NewMachine creates a new Machine in a cluster, usually using an Image
// for the supplied Kubernetes version.
func (vd *Driver) NewMachine(machinename string, clustername string, k8sversion string) (drivercore.Machine, error) {
	err := vd.validate()
	if err != nil {
		return nil, err
	}

	image, err := vd.GetImage(k8sversion)
	if err != nil {
		return nil, err
	}

	localimage, ok := image.(*Image)
	if !ok {
		return nil, fmt.Errorf("unknown error verifying image for Kubernetes version %v", k8sversion)
	}

	machinefile, err := machineFilePath(vd.QualifiedMachineName(machinename, clustername))
	if err != nil {
		return nil, errors.Wrap(err, "machine file not accessible")
	}

	err = writemanifest(machinefile, localimage.ImageSourceURL)
	if err != nil {
		return nil, errors.Wrap(err, "machine file not written")
	}

	limactlparams := []string{
		"create",
		"--name=" + vd.QualifiedMachineName(machinename, clustername),
		machinefile,
	}

	result, err := vd.runwithresults(limactlparams...)
	if err != nil {
		// TODO: Consider doing a compensatory `limactl rm`.
		// Not risking it in the current version.
		errMsg := result.LastLogErrorMessage()
		if errMsg == "" {
			errMsg = "error during limactl create"
		}
		return nil, errors.Wrap(err, errMsg)
	}

	return &Machine{
		driver:      vd,
		name:        machinename,
		clustername: clustername,
		status:      drivercore.MachineStatusStopped,
	}, nil
}
