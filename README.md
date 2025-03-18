# USB Auto Backup

## 📌 What is USB Auto Backup?
USB Auto Backup is a lightweight, command-line tool written in **Go** that automatically detects USB drives when they are inserted and backs up their contents to a designated folder on your system.

## 🔥 Why Does This Exist?
Ever plugged in a USB drive, intending to back up files, but forgot? This tool ensures that never happens again. 

With USB Auto Backup, you:
- **Never lose important files** due to forgetfulness.
- **Automate USB backups** without manual intervention.
- **Save time** with an efficient and hands-free process.
- **Avoid complex backup software** and bloated UI applications.

This tool is designed to run quietly in the background, checking for new USB drives and backing up data without needing user interaction. It’s perfect for professionals, students, and anyone who regularly uses USB storage.

## 🚀 Features
✅ Automatically detects USB drives when plugged in.
✅ Copies all files from USB to a designated backup folder.
✅ Creates timestamped backups to prevent overwrites.
✅ Runs efficiently in the background without user input.
✅ Simple and lightweight—just run and forget!

## ⚙️ How It Works
1. The tool continuously monitors your system for new USB drives.
2. When a USB drive is detected, it identifies the drive letter.
3. It then copies all files from the USB to a backup folder (`C:\USB_Backups\[timestamp]`).
4. The tool loops and waits for the next USB insertion.

## 🎯 Roadmap
- 🔹 Incremental backups (only copy new/changed files).
- 🔹 Configurable settings (custom backup location, file filters, etc.).
- 🔹 Logging system for backup tracking.
- 🔹 Windows system tray notifications.
- 🔹 Potential cross-platform support.

## 💡 Contributions & Feedback
This is an open-source project, and contributions are welcome! If you have feature requests, bug reports, or ideas, feel free to open an issue or submit a pull request.

Let’s make USB backups **effortless and foolproof**! 🚀
