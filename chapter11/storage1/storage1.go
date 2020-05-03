package storage

import (
	"fmt"
	"log"
	"net/smtp"
)

var usage = make(map[string]int64)

func bytesInUse(username string) int64 { return usage[username] }

// Email sender configuration
// NOTE: never put passowrds in source code!
const (
	sender   = "notification@example.com"
	password = "correcthorsebatterystaple"
	hostname = "smtp.example.com"
	template = `Warning you are using %d bytes of storage, %d%% of your quota.`
)

func CheckQuota(username string) {
	used := bytesInUse(username)
	const quota = 1_000_000_000 // 1 GB
	percent := 100 * used / quota
	if percent < 90 {
		return // OK
	}
	msg := fmt.Sprintf(template, used, percent)
	auth := smtp.PlainAuth("", sender, password, hostname)
	err := smtp.SendMail(hostname+":587", auth, sender, []string{username}, []byte(msg))
	if err != nil {
		log.Printf("smtp.SendMail(%s) failed; %s", username, err)
	}
}
