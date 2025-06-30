package driverlima

import (
	"strings"

	"github.com/kuttiproject/drivercore"
	"github.com/kuttiproject/sshclient"
)

// TODO: Look at parameterizing these
var (
	limaUsername = "kuttiadmin"
	limaPassword = "Pass@word1"
)

// runwithresults allows running commands inside a VM Host.
// It does this by creating an SSH session with the host.
func (vh *Machine) runwithresults(execpath string, paramarray ...string) (string, error) {
	client := sshclient.NewWithPassword(limaUsername, limaPassword)
	params := append([]string{execpath}, paramarray...)
	output, err := client.RunWithResults(vh.SSHAddress(), strings.Join(params, " "))
	if err != nil {
		return "", err
	}
	return output, nil
}

var limaCommands = map[drivercore.PredefinedCommand]func(*Machine, ...string) error{
	drivercore.RenameMachine: renamemachine,
}

func renamemachine(vh *Machine, params ...string) error {
	newname := params[0]
	execname := "set-hostname.sh"

	_, err := vh.runwithresults(
		"/usr/bin/sudo",
		execname,
		newname,
	)

	return err
}
