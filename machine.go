package driverlima

import (
	"fmt"
	"strings"

	"github.com/kuttiproject/drivercore"
	"github.com/kuttiproject/kuttilog"
	"github.com/pkg/errors"
)

type Machine struct {
	driver *Driver

	name        string
	clustername string
	sshhostport int
	// savedipaddress string
	limainfo     *limaInfo
	status       drivercore.MachineStatus
	errormessage string
}

// Name is the name of the machine.
// The operating system hostname should match this.
func (m *Machine) Name() string {
	return m.name
}

// Status can be drivercore.MachineStatusRunning, drivercore.MachineStatusStopped
// drivercore.MachineStatusUnknown or drivercore.MachineStatusError.
func (m *Machine) Status() drivercore.MachineStatus {
	m.get()
	return m.status
}

// Error returns the last error caused when manipulating this machine.
// A valid value can be expected only when Status() returns
// drivercore.MachineStatusError.
func (m *Machine) Error() string {
	if m.limainfo == nil {
		m.get()
	}
	return m.errormessage
}

// IPAddress returns the current IP Address of this Machine.
// A valid value can be expected only when Status() returns
// drivercore.MachineStatusRunning.
func (m *Machine) IPAddress() string {
	kuttilog.Println(kuttilog.MaxLevel(), "In ipaddress 1")

	status := m.Status()
	if status != drivercore.MachineStatusRunning {
		return ""
	}

	kuttilog.Println(kuttilog.MaxLevel(), "In ipaddress 2")

	//return "0.0.1.0"
	limactlargs := []string{
		"shell",
		m.qName(),
		"get-primary-ip.sh",
	}

	result, err := m.driver.runwithresults(limactlargs...)
	if err != nil {
		kuttilog.Printf(kuttilog.Error, "Error fetching ipaddess: %v", err)
		m.status = drivercore.MachineStatusError
		m.errormessage = err.Error()
		return err.Error()
	}

	if !result.isRaw {
		return "unexpected format error fetching ipaddress"
	}

	return strings.TrimRight(result.rawResult, "\n")
}

// SSHAddress returns the host address and port number to SSH into this Machine.
// For drivers that use NAT netwoking, the host address will be 'localhost'.
func (m *Machine) SSHAddress() string {
	return fmt.Sprintf("localhost:%v", m.sshhostport)
}

// Start starts a Machine.
// Note that a Machine may not be ready for further operations at the end of this,
// and therefore its status may not change immediately.
// See WaitForStateChange().
func (m *Machine) Start() error {
	limctlparams := []string{
		"start",
		m.qName(),
	}

	_, err := m.driver.runwithresults(limctlparams...)
	if err != nil {
		return err
	}

	m.limainfo = nil
	return nil
}

// Stop stops a Machine.
// Note that a Machine may not be ready for further operations at the end of this,
// and therefore its status will not change immediately.
// See WaitForStateChange().
func (m *Machine) Stop() error {
	limctlparams := []string{
		"stop",
		m.qName(),
	}

	_, err := m.driver.runwithresults(limctlparams...)
	if err != nil {
		return err
	}

	m.limainfo = nil
	return nil
}

// ForceStop stops a Machine forcibly.
// This operation should set the status to drivercore.MachineStatusStopped.
func (m *Machine) ForceStop() error {
	limctlparams := []string{
		"stop",
		"-f",
		m.qName(),
	}

	_, err := m.driver.runwithresults(limctlparams...)
	if err != nil {
		return err
	}

	m.limainfo = nil
	return nil
}

// WaitForStateChange waits the specified number of seconds, or until the Machine
// status changes.
// WaitForStateChange should be called after calls to Start() or Stop(), before
// any other operation. It should not be called _before_ Stop().
// The lima driver silently ignores WaitForStateChanged.
func (m *Machine) WaitForStateChange(timeoutinseconds int) {

}

// ForwardPort creates a rule to forward the specified Machine port to the
// specified physical host port.
func (m *Machine) ForwardPort(hostport int, machineport int) error {
	if machineport == 22 {
		return m.ForwardSSHPort(hostport)
	}

	kuttilog.Println(kuttilog.Verbose, "ForwardPort not explicitly required by lima driver")
	return nil
}

// UnforwardPort removes the rule which forwarded the specified Machine port.
func (m *Machine) UnforwardPort(machineport int) error {
	kuttilog.Println(kuttilog.Verbose, "UnforwardPort not explicitly required by lima driver")
	return nil
}

// ForwardSSHPort forwards the SSH port of this Machine to the specified
// physical host port.
func (m *Machine) ForwardSSHPort(hostport int) error {
	status := m.Status()
	if status != drivercore.MachineStatusStopped {
		return fmt.Errorf("can only forward ports when machine is stopped")
	}

	machinefile, err := machineFilePath(m.qName())
	if err != nil {
		return err
	}

	vd := m.driver

	// Set hostPort in manifest file
	limactlparams := []string{
		"edit",
		machinefile,
		"--set",
		fmt.Sprintf(".ssh.localPort = %v", hostport),
	}
	_, err = vd.runwithresults(limactlparams...)
	if err != nil {
		return errors.Wrap(err, "could not update port forwarding in machine file")
	}

	// Set hostPort in created VM
	limactlparams[1] = m.qName()
	_, err = vd.runwithresults(limactlparams...)
	if err != nil {
		return errors.Wrap(err, "could not update port forwarding in lima vm")
	}

	return nil
}

// ImplementsCommand returns true if the driver implements the specified
// predefined operation.
func (m *Machine) ImplementsCommand(command drivercore.PredefinedCommand) bool {
	_, ok := limaCommands[command]
	return ok
}

// ExecuteCommand executes the specified predefined operation.
func (m *Machine) ExecuteCommand(command drivercore.PredefinedCommand, params ...string) error {
	commandfunc, ok := limaCommands[command]
	if !ok {
		return fmt.Errorf(
			"command '%v' not implemented",
			command,
		)
	}

	return commandfunc(m, params...)
}

func (m *Machine) qName() string {
	return m.driver.QualifiedMachineName(m.name, m.clustername)
}

func (m *Machine) get() {
	limactlargs := []string{
		"list",
		m.qName(),
		"--format",
		"json",
	}

	resultobj, err := m.driver.runwithresults(limactlargs...)
	if err != nil {
		m.status = drivercore.MachineStatusError
		m.errormessage = err.Error()
		return
	}

	result := resultobj.machineInfos[0]

	m.status = drivercore.MachineStatus(result.Status)
	m.errormessage = ""
}
