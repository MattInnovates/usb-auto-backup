package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ServiceName defines how the service will appear in Windows
const ServiceName = "USB_Auto_Backup"

// installService installs the app as a Windows service
func installService(configPath string) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exePath, _ = filepath.Abs(exePath)

	cmd := exec.Command("sc", "create", ServiceName,
		"binPath=", fmt.Sprintf("\"%s --config %s\"", exePath, configPath),
		"start=", "auto")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// startService starts the service
func startService() error {
	cmd := exec.Command("sc", "start", ServiceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// stopService stops the service
func stopService() error {
	cmd := exec.Command("sc", "stop", ServiceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// uninstallService removes the service
func uninstallService() error {
	cmd := exec.Command("sc", "delete", ServiceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
