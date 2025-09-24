package main

import (
	"fmt"
	"os"

	"github.com/MattInnovates/usb-auto-backup/pkg/config"
	"github.com/MattInnovates/usb-auto-backup/pkg/notify"
)

type USBDevice struct {
	Label  string
	Serial string
}

func detectUSBs() []USBDevice {
	// Fake detection for now
	return []USBDevice{
		{Label: "WORK_USB", Serial: "1234-ABCD"},
		{Label: "KINGSTON", Serial: "9XYZ-9876"},
	}
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("❌ Failed to load config:", err)
		os.Exit(1)
	}

	// Show startup toast
	notify.Show("USB Auto Backup", "Service started and monitoring USB devices...")

	fmt.Println("✅ USB Auto Backup starting...")
	fmt.Println("Backup directory:", cfg.BackupDir)

	usbList := detectUSBs()
	if len(usbList) == 0 {
		fmt.Println("No USBs detected.")
		return
	}

	for _, dev := range usbList {
		fmt.Printf("🔍 Detected USB: %s (Serial %s)\n", dev.Label, dev.Serial)

		added, err := cfg.EnrolDevice(dev.Label, dev.Serial)
		if err != nil {
			fmt.Println("❌ Failed to enrol device:", err)
			continue
		}

		if added {
			fmt.Printf("✅ Enrolled new device: %s (%s)\n", dev.Label, dev.Serial)
			notify.Show("USB Auto Backup", "New device enrolled: "+dev.Label)
		} else {
			fmt.Printf("ℹ️ Device already enrolled: %s (%s)\n", dev.Label, dev.Serial)
			// 🚫 No toast here to avoid spam — only notify during backup
		}
	}
	fmt.Println("✅ Config updated successfully.")
}
