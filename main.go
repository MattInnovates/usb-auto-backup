package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Set your backup directory
const backupDir = "C:\\USB_Backups"

func main() {
	fmt.Println("USB Auto Backup Running...")

	// Keep checking for new USB drives
	for {
		drives := getUSBDrives()
		for _, drive := range drives {
			fmt.Println("USB Detected:", drive)
			backupUSB(drive)
		}
		time.Sleep(10 * time.Second) // Check every 10 seconds
	}
}

// Get list of USB drives
func getUSBDrives() []string {
	out, err := exec.Command("wmic", "logicaldisk", "where", "drivetype=2", "get", "deviceid").Output()
	if err != nil {
		fmt.Println("Error detecting USB drives:", err)
		return nil
	}

	var drives []string
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 2 && strings.HasSuffix(line, ":") {
			drives = append(drives, line)
		}
	}
	return drives
}

// Backup all files from USB
func backupUSB(drive string) {
	srcPath := drive + "\\"
	destPath := filepath.Join(backupDir, time.Now().Format("2006-01-02_15-04-05"))

	fmt.Println("Backing up USB from", srcPath, "to", destPath)

	err := copyDir(srcPath, destPath)
	if err != nil {
		fmt.Println("Backup failed:", err)
	} else {
		fmt.Println("Backup completed successfully!")
	}
}

// Copy a directory recursively
func copyDir(src, dest string) error {
	err := os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, destPath)
		} else {
			err = copyFile(srcPath, destPath)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

// Copy a single file
func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}
