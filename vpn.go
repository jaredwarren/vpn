package vpn

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// VPN info about vpn connection
type VPN struct {
	Name    string
	Timeout time.Duration
	Status  string
}

// Dial connects to VPN
func Dial(name string) (v *VPN, err error) {
	v = &VPN{
		Name:    name,
		Timeout: 60 * time.Second,
		Status:  "Unknown",
	}

	// Connect
	err = v.Connect()
	return
}

// Connect to VPN, wait for status to read "Connected", or timeout
func (v *VPN) Connect() error {
	v.Status = "Unknown"
	_, err := runBash(fmt.Sprintf(`scutil --nc start "%s"`, v.Name), nil)
	if err != "" {
		return errors.New(err)
	}
	// Wait for VPN Connection.
	c1 := make(chan bool, 1)
	go func() {
		for {
			out, _ := runBash(fmt.Sprintf(`scutil --nc status "%s" | sed -n 1p `, v.Name), nil)
			if out == "Connected" {
				v.Status = "Connected"
				c1 <- true
				return
			}
		}
	}()
	select {
	case <-c1:
	case <-time.After(v.Timeout):
		return errors.New("timeout waiting for VPN connection")
	}
	return nil
}

// Close from VPN, wait for status to read "Disconnected", or timeout
func (v *VPN) Close() error {
	v.Status = "Unknown"
	_, err := runBash(fmt.Sprintf(`scutil --nc stop "%s"`, v.Name), nil)
	if err != "" {
		return errors.New(err)
	}
	// Wait for VPN Connection.
	c1 := make(chan bool, 1)
	go func() {
		for {
			out, _ := runBash(fmt.Sprintf(`scutil --nc status "%s" | sed -n 1p `, v.Name), nil)
			if out == "Disconnected" {
				v.Status = "Disconnected"
				c1 <- true
				return
			}
		}
	}()
	select {
	case <-c1:
	case <-time.After(v.Timeout):
		return errors.New("timeout waiting for VPN to disconnect")
	}
	return nil
}

// GetVPN return VPN status by name
func GetVPN(name string) string {
	out, _ := runBash(fmt.Sprintf(`scutil --nc status "%s" | sed -n 1p `, name), nil)
	return out
}

// runBash run a bash command an returns stdout and stderr
func runBash(commandString string, env []string) (stdOut, stdErr string) {
	cmd := exec.Command("bash", "-c", commandString)
	if len(env) > 0 {
		cmd.Env = env
	}
	var stdOutBuf bytes.Buffer
	var stdErrBuf bytes.Buffer
	cmd.Stdout = &stdOutBuf
	cmd.Stderr = &stdErrBuf
	cmd.Run()
	stdOut = strings.TrimSuffix(stdOutBuf.String(), "\n")
	stdErr = strings.TrimSuffix(stdErrBuf.String(), "\n")
	return
}
