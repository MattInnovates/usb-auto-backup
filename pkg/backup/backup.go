package backup

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/MattInnovates/usb-auto-backup/pkg/config"
)

// USBDevice mirrors the struct in main.go
type USBDevice struct {
	Label  string
	Serial string
	Path   string // root path (e.g. "F:\")
}

// BackupDevice copies USB contents into the backup dir, then uploads via SFTP if enabled.
func BackupDevice(dev USBDevice, cfg *config.Config) error {
	timestamp := time.Now().Format("20060102-150405")
	destRoot := filepath.Join(cfg.BackupDir, fmt.Sprintf("%s-%s-%s", dev.Label, dev.Serial, timestamp))

	fmt.Printf("📂 Starting backup for %s (%s) → %s\n", dev.Label, dev.Serial, destRoot)

	sourceRoot := dev.Path

	err := copyDir(sourceRoot, destRoot)
	if err != nil {
		return fmt.Errorf("failed to copy files: %w", err)
	}

	fmt.Printf("✅ Backup complete for %s (%s)\n", dev.Label, dev.Serial)

	// Optional SFTP upload
	if cfg.SFTP.Enabled {
		err = uploadSCP(destRoot, cfg)
		if err != nil {
			return fmt.Errorf("sftp upload failed: %w", err)
		}
		fmt.Printf("🌍 SFTP upload complete for %s\n", dev.Label)
	}

	return nil
}

// uploadSCP uploads a folder to the remote server using scp (requires OpenSSH in PATH).
func uploadSCP(localPath string, cfg *config.Config) error {
	fmt.Printf("🌍 Uploading %s to %s@%s:%s\n",
		localPath, cfg.SFTP.Username, cfg.SFTP.Host, cfg.SFTP.Remote)

	// Build target string like user@host:/remote/path
	target := fmt.Sprintf("%s@%s:%s", cfg.SFTP.Username, cfg.SFTP.Host, cfg.SFTP.Remote)

	// -r = recursive copy, -P = port
	cmd := exec.Command("scp", "-r", "-P", fmt.Sprint(cfg.SFTP.Port), localPath, target)

	// Set password? Normally scp uses SSH keys or will prompt interactively.
	// If password auth is needed, recommend Pageant/ssh-agent or pre-shared key.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// copyDir recursively copies a directory tree, preserving structure.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		return copyFile(path, target, info.Mode())
	})
}

func copyFile(src, dst string, perm os.FileMode) error {
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

	return out.Chmod(perm)
}
