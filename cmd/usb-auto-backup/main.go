package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/MattInnovates/usb-auto-backup/pkg/backup"
	"github.com/MattInnovates/usb-auto-backup/pkg/config"
	"github.com/MattInnovates/usb-auto-backup/pkg/notify"
)

// USBDevice now includes Path (e.g. "F:\\")
type USBDevice struct {
	Label  string
	Serial string
	Path   string
}

// Windows drive type constants
const (
	DRIVE_UNKNOWN     = 0
	DRIVE_NO_ROOT_DIR = 1
	DRIVE_REMOVABLE   = 2
	DRIVE_FIXED       = 3
	DRIVE_REMOTE      = 4
	DRIVE_CDROM       = 5
	DRIVE_RAMDISK     = 6
)

const ServiceName = "USB_Auto_Backup"
const InstallDir = "C:\\usb-auto-backup"

// -----------------------------
// Helpers
// -----------------------------

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func serviceExists() bool {
	cmd := exec.Command("sc", "query", ServiceName)
	err := cmd.Run()
	return err == nil
}

// -----------------------------
// USB Detection
// -----------------------------

func detectUSBs() []USBDevice {
	var devices []USBDevice
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	getDriveType := kernel32.MustFindProc("GetDriveTypeW")
	getVolumeInfo := kernel32.MustFindProc("GetVolumeInformationW")

	for letter := 'D'; letter <= 'Z'; letter++ {
		rootPath := fmt.Sprintf("%c:\\", letter)

		if _, err := os.Stat(rootPath); os.IsNotExist(err) {
			continue
		}

		lpRoot, _ := syscall.UTF16PtrFromString(rootPath)
		ret, _, _ := getDriveType.Call(uintptr(unsafe.Pointer(lpRoot)))
		if ret != DRIVE_REMOVABLE {
			continue
		}

		volName := make([]uint16, syscall.MAX_PATH+1)
		var serial, maxCompLen, fsFlags uint32
		fsName := make([]uint16, syscall.MAX_PATH+1)

		_, _, _ = getVolumeInfo.Call(
			uintptr(unsafe.Pointer(lpRoot)),
			uintptr(unsafe.Pointer(&volName[0])),
			uintptr(len(volName)),
			uintptr(unsafe.Pointer(&serial)),
			uintptr(unsafe.Pointer(&maxCompLen)),
			uintptr(unsafe.Pointer(&fsFlags)),
			uintptr(unsafe.Pointer(&fsName[0])),
			uintptr(len(fsName)),
		)

		label := syscall.UTF16ToString(volName)
		if label == "" {
			label = fmt.Sprintf("USB_%c", letter)
		}

		devices = append(devices, USBDevice{
			Label:  label,
			Serial: fmt.Sprintf("%X", serial),
			Path:   rootPath,
		})
	}

	return devices
}

// -----------------------------
// Agent Logic (continuous loop)
// -----------------------------

func runAgent() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("❌ Failed to load config:", err)
		os.Exit(1)
	}

	notify.Show("USB Auto Backup", "Service running and monitoring USB devices...")
	fmt.Println("✅ USB Auto Backup started")
	fmt.Println("Backup directory:", cfg.BackupDir)

	seen := make(map[string]USBDevice)

	for {
		usbList := detectUSBs()
		current := make(map[string]USBDevice)

		// Detect new devices
		for _, dev := range usbList {
			current[dev.Serial] = dev
			if _, ok := seen[dev.Serial]; !ok {
				fmt.Printf("🔌 New USB detected: %s (%s)\n", dev.Label, dev.Serial)
				notify.Show("USB Auto Backup", "New device detected: "+dev.Label)

				added, err := cfg.EnrolDevice(dev.Label, dev.Serial)
				if err == nil && added {
					fmt.Printf("✅ Enrolled device: %s\n", dev.Label)
				}

				// Run backup
				if err := backup.BackupDevice(backup.USBDevice(dev), cfg); err != nil {
					fmt.Printf("❌ Backup failed for %s: %v\n", dev.Label, err)
					notify.Show("USB Auto Backup", "Backup failed for "+dev.Label)
				} else {
					notify.Show("USB Auto Backup", "Backup complete for "+dev.Label)
				}
			}
		}

		// Detect removed devices
		for serial, dev := range seen {
			if _, ok := current[serial]; !ok {
				fmt.Printf("❌ Device removed: %s (%s)\n", dev.Label, serial)
				notify.Show("USB Auto Backup", "Device removed: "+dev.Label)
			}
		}

		seen = current
		time.Sleep(5 * time.Second)
	}
}

// -----------------------------
// Main Entry
// -----------------------------

func main() {
	// Handle service commands explicitly to avoid prompt recursion
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "service":
			if len(os.Args) < 3 {
				fmt.Println("Usage: usb-auto-backup service [install|start|stop|uninstall] [--config path]")
				return
			}
			action := os.Args[2]
			switch action {
			case "install":
				// Ensure install dir exists
				os.MkdirAll(InstallDir, 0755)

				exePath, _ := os.Executable()
				dstExe := filepath.Join(InstallDir, "usb-auto-backup.exe")
				copyFile(exePath, dstExe)

				// Copy config if available
				cfgPath := "config.json"
				if len(os.Args) >= 5 && os.Args[3] == "--config" {
					cfgPath = os.Args[4]
				}
				if _, err := os.Stat(cfgPath); err == nil {
					copyFile(cfgPath, filepath.Join(InstallDir, "config.json"))
				}

				cmd := exec.Command("sc", "create", ServiceName,
					"binPath=", fmt.Sprintf("\"%s --config %s\"", dstExe, filepath.Join(InstallDir, "config.json")),
					"start=", "auto")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					fmt.Println("❌ Failed to install service:", err)
					return
				}
				fmt.Println("✅ Service installed.")
			case "start":
				cmd := exec.Command("sc", "start", ServiceName)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				_ = cmd.Run()
			case "stop":
				cmd := exec.Command("sc", "stop", ServiceName)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				_ = cmd.Run()
			case "uninstall":
				cmd := exec.Command("sc", "delete", ServiceName)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				_ = cmd.Run()
			default:
				fmt.Println("Unknown service command:", action)
			}
			return
		}
	}

	// Normal run → if service not installed, prompt user
	if !serviceExists() {
		fmt.Println("USB Auto Backup is not installed as a Windows service.")
		fmt.Print("Would you like to install it now? (y/n): ")

		var resp string
		fmt.Scanln(&resp)
		if strings.ToLower(resp) == "y" {
			// Relaunch self with service install
			exePath, _ := os.Executable()
			cmd := exec.Command(exePath, "service", "install", "--config", "config.json")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			_ = cmd.Run()

			cmd = exec.Command(exePath, "service", "start")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			_ = cmd.Run()

			fmt.Println("✅ Service installed and started.")
			return
		}
	}

	// If service already exists or user said no → run normally
	runAgent()
}
