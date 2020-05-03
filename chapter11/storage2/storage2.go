package storage

import (
	"fmt"
	"log"
	"net/smtp"
)

// Email sender configuration
// NOTE: never put passowrds in source code!
const (
	sender   = "notification@example.com"
	password = "correcthorsebatterystaple"
	hostname = "smtp.example.com"
	template = `Warning you are using %d bytes of storage, %d%% of your quota.`
)

var usage = make(map[string]int64)

func bytesInUse(username string) int64 { return usage[username] }

var notifyUser = func(username, msg string) {
	auth := smtp.PlainAuth("", sender, password, hostname)
	err := smtp.SendMail(hostname+":587", auth, sender, []string{username}, []byte(msg))
	if err != nil {
		log.Printf("smtp.SendMail(%s) failed: %s", username, err)
	}
}

func CheckQuota(username string) {
	used := bytesInUse(username)
	const quota = 1_000_000_000
	percent := 100 * used / quota
	if percent < 90 {
		return
	}
	msg := fmt.Sprintf(template, used, percent)
	notifyUser(username, msg)

}
