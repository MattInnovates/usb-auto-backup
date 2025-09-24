package notify

import (
	"github.com/go-toast/toast"
)

func Show(title, message string) error {
	notification := toast.Notification{
		AppID:   "USB Auto Backup",
		Title:   title,
		Message: message,
	}
	return notification.Push()
}
