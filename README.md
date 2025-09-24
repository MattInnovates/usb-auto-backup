# USB Auto Backup

A Windows background service + CLI tool that automatically backs up USB drives 
and optionally uploads to SFTP.

## Features (planned)
- 🔍 Detect USB drive insertions
- 📂 Auto-backup to local folder
- 🌍 Optional upload to SFTP
- 🔔 Windows notifications
- 🛠 Runs as service or CLI

## Development
Current development has fake drive's to test Toasts and auto add to `config.json`

Clone the repo:

```bash
git clone https://github.com/MattInnovates/usb-auto-backup.git
cd usb-auto-backup
go mod tidy

go run ./cmd/usb-auto-backup