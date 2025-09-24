package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// DeviceRule defines a USB device we can back up
type DeviceRule struct {
	Label  string `json:"label"`
	Serial string `json:"serial"`
	Upload string `json:"upload"` // e.g. "sftp" or "local"
}

// SFTPConfig holds remote upload settings
type SFTPConfig struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Remote   string `json:"remotePath"`
}

// Config holds global settings
type Config struct {
	BackupDir string       `json:"backupDir"`
	Notify    bool         `json:"notify"`
	Devices   []DeviceRule `json:"devices"`
	SFTP      SFTPConfig   `json:"sftp"`
}

var configPath = filepath.Join("configs", "config.json")

// Load reads the config.json file
func Load() (*Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save writes the current config back to config.json
func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// EnrolDevice adds a new USB device if not already present
func (c *Config) EnrolDevice(label, serial string) (bool, error) {
	// check if already exists
	for _, dev := range c.Devices {
		if dev.Serial == serial {
			return false, nil // already enrolled
		}
	}

	// add new device with default settings
	newDev := DeviceRule{
		Label:  label,
		Serial: serial,
		Upload: "local", // default to local backup
	}
	c.Devices = append(c.Devices, newDev)

	// save updated config
	if err := c.Save(); err != nil {
		return false, err
	}

	return true, nil
}
