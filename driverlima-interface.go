package driverlima

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kuttiproject/kuttilog"
	"github.com/kuttiproject/workspace"
)

type limaInfo struct {
	Name          string `json:"name"`
	Hostname      string `json:"hostname"`
	Status        string `json:"status"`
	Dir           string `json:"dir"`
	SSHLocalPort  int    `json:"sshLocalPort"`
	SSHConfigFile string `json:"sshConfigFile"`
}

type logEntry struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
	Time  string `json:"time"`
}

type limaResult struct {
	logEntries    []logEntry
	machineInfos  []limaInfo
	rawResult     string
	isRaw         bool
	isLogEntry    bool
	isMachineInfo bool
}

func (lr *limaResult) LastLogErrorMessage() string {
	if lr == nil {
		return ""
	}

	lastIndex := len(lr.logEntries) - 1

	if lastIndex >= 0 {
		lastEntry := lr.logEntries[lastIndex]
		if lastEntry.Level == "fatal" || lastEntry.Level == "error" {
			return lastEntry.Msg
		}
	}

	return ""
}

// ParseLogString reads a multi-line string, parsing each line as either
// a logEntry or a limaInfo JSON object. It returns two slices
// containing the parsed objects of each kind, and an error if any
// line fails to parse or doesn't strictly match one of the two expected shapes.
func newLimaResult(input string) (*limaResult, error) {
	// Initialize result to store the parsed objects
	result := &limaResult{
		logEntries:   []logEntry{},
		machineInfos: []limaInfo{},
	}
	// Initialize slices to store the parsed objects
	// var logEntries []logEntry
	// var machineInfos []limaInfo

	// Create a new scanner to read the input string line by line
	scanner := bufio.NewScanner(strings.NewReader(input))
	lineNumber := 0 // Keep track of the current line number for error reporting

	// Iterate over each line in the input string
	for scanner.Scan() {
		lineNumber++           // Increment line number for each new line
		line := scanner.Text() // Get the current line content as a string

		// Skip empty lines to avoid parsing errors
		if strings.TrimSpace(line) == "" {
			continue
		}

		// First, try to unmarshal the line into a generic map to inspect its structure (keys)
		var rawMap map[string]any
		if err := json.Unmarshal([]byte(line), &rawMap); err != nil {
			// If unmarshaling to a map fails, it's not a valid JSON.
			// If no JSON lines have been read so far, return the
			// full input string as a raw result.
			if !result.isLogEntry && !result.isMachineInfo {
				result.isRaw = true
				result.rawResult = input
				return result, nil
			}

			//return nil, fmt.Errorf("line %d: invalid JSON format: %w", lineNumber, err)
			// Otherwise, add the line to result.rawresult. No error.
			result.rawResult += line + "\n"
		}

		// Define the expected keys for each shape for strict matching
		expectedLogEntryKeys := map[string]bool{
			"level": true,
			"msg":   true,
			"time":  true,
		}
		expectedLimaInfoKeys := map[string]bool{
			"name":          true,
			"hostname":      true,
			"status":        true,
			"dir":           true,
			"sshConfigFile": true,
		}

		// Check if the current line's JSON map matches the LogEntry shape
		isLogEntryShape := true
		for key := range expectedLogEntryKeys { // Check if all keys match
			if _, ok := rawMap[key]; !ok {
				isLogEntryShape = false
				break
			}
		}

		// Check if the current line's JSON map matches the MachineInfo shape
		isMachineInfoShape := true
		for key := range expectedLimaInfoKeys { // Check if all keys match
			if _, ok := rawMap[key]; !ok {
				isMachineInfoShape = false
				break
			}
		}

		// Determine which type the line matches and unmarshal accordingly
		if isLogEntryShape && !isMachineInfoShape {
			// It matches LogEntry shape exclusively
			var le logEntry
			if err := json.Unmarshal([]byte(line), &le); err != nil {
				// This error should ideally not happen if rawMap unmarshal was successful,
				// but it's a good safeguard.
				return nil, fmt.Errorf("line %d: failed to unmarshal into LogEntry: %w", lineNumber, err)
			}
			result.logEntries = append(result.logEntries, le)
			result.isLogEntry = true
		} else if isMachineInfoShape && !isLogEntryShape {
			// It matches MachineInfo shape exclusively
			var mi limaInfo
			if err := json.Unmarshal([]byte(line), &mi); err != nil {
				// Similar safeguard as above for MachineInfo
				return nil, fmt.Errorf("line %d: failed to unmarshal into MachineInfo: %w", lineNumber, err)
			}
			result.machineInfos = append(result.machineInfos, mi)
			result.isMachineInfo = true
		} else {
			// The line does not strictly match either shape, or it somehow matches both
			// (which is unlikely given distinct key sets, but handled for robustness).
			return nil, fmt.Errorf("line %d: does not strictly match either expected JSON shape. Content: %s", lineNumber, line)
		}
	}

	// Check for any errors that occurred during scanning (e.g., I/O errors)
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input string: %w", err)
	}

	// Return the result and no error
	return result, nil
}

func machineDir() (string, error) {
	return workspace.CacheSubDir("driver-lima-machines")
}

func machineFilePath(machineqname string) (string, error) {
	machinedir, err := machineDir()
	if err != nil {
		return "", err
	}

	machinefile := filepath.Join(machinedir, machineqname+".yaml")
	return machinefile, nil
}

func findLimaCtl() (string, error) {
	// First, try looking up limactl on the path
	toolpath, err := exec.LookPath("limactl")
	if err == nil {
		return toolpath, nil
	}

	return "", errors.New("limactl not found")
}

func (d *Driver) validate() error {
	if d.validated {
		return nil
	}

	limactlpath, err := findLimaCtl()
	if err != nil {
		d.status = "Error"
		d.errormessage = err.Error()
		return err
	}

	d.limactlpath = limactlpath
	d.status = "Ready"
	d.errormessage = ""
	d.validated = true

	return nil
}

func (d *Driver) runwithresults(args ...string) (*limaResult, error) {
	limactlargs := []string{
		"--tty=false",
		"--log-format=json",
	}

	switch kuttilog.LogLevel() {
	case kuttilog.Error:
		limactlargs = append(limactlargs, "--log-level", "error")
	case kuttilog.Minimal:
		limactlargs = append(limactlargs, "--log-level", "warn")
	case kuttilog.Info:
		limactlargs = append(limactlargs, "--log-level", "info")
	case kuttilog.Verbose:
		limactlargs = append(limactlargs, "--log-level", "trace")
	case kuttilog.Debug:
		limactlargs = append(limactlargs, "--log-level", "debug")
	}
	limactlargs = append(limactlargs, args...)
	resultstring, err := workspace.RunWithResults(d.limactlpath, limactlargs...)
	result, err2 := newLimaResult(resultstring)
	if err2 != nil {
		err = errors.Join(err, err2)
	}
	return result, err
}

//go:embed assets/knode.yaml
var manifest string

func writemanifest(manifestpath string, imagesourceurl string) error {
	manifestFile, err := os.Create(manifestpath)
	if err != nil {
		return err
	}

	defer manifestFile.Close()

	newmanifest := strings.Replace(manifest, "{{ .ImageSourceUrl }}", imagesourceurl, 1)
	_, err = manifestFile.WriteString(newmanifest)
	if err != nil {
		return err
	}

	return nil
}
