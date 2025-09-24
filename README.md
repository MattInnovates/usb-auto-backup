# USB Auto Backup  

A Windows background service + CLI tool that automatically backs up USB drives  
and optionally uploads them to SFTP.  

![Go Report Card](https://goreportcard.com/badge/github.com/MattInnovates/usb-auto-backup)  
![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)  

---

## ✨ Features  

- 🔍 Detects USB drive insertions in real-time  
- 📂 Auto-backup to a local folder  
- 🌍 Optional upload to SFTP  
- 🔔 Windows notifications (toasts)  
- 🛠 Runs as a **Windows service** or as a **CLI** tool  

---

## 🚀 Quick Start  

Clone the repo:  

```bash
git clone https://github.com/MattInnovates/usb-auto-backup.git
cd usb-auto-backup
go mod tidy
```

Run directly for testing:  

```bash
go run ./cmd/usb-auto-backup
```

---

## ⚙️ Configuration  

The tool uses a `config.json` file. Here’s an example:  

```json
{
  "backupPath": "D:/USBBackups",
  "sftp": {
    "enabled": true,
    "host": "sftp.example.com",
    "port": 22,
    "username": "myuser",
    "password": "mypassword",
    "remotePath": "/backups"
  }
}
```

- `backupPath`: Where local backups are stored.  
- `sftp`: Optional SFTP upload settings.  

---

## 📖 Usage  

### Run as CLI  

```bash
usb-auto-backup --config config.json
```

### Install as a Windows Service  

```bash
usb-auto-backup service install --config config.json
usb-auto-backup service start
```

### Uninstall Service  

```bash
usb-auto-backup service stop
usb-auto-backup service uninstall
```

---

## 🧭 Roadmap  

- [x] Detect USB drive insertions  
- [x] Auto-backup to local folder  
- [ ] Optional SFTP upload  
- [x] Windows notifications  
- [ ] Configurable include/exclude rules  
- [ ] Hash verification for backup integrity  
- [ ] Dry-run mode for testing  

---

## 🤝 Contributing  

Contributions, issues, and feature requests are welcome!  
Feel free to open an issue or submit a PR.  

---

## 📜 License  

This project is licensed under the MIT License — see [LICENSE](LICENSE) for details.  
