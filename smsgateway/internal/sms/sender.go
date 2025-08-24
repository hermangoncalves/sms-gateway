package sms

import (
	"log"
	"os/exec"
)

func SendSMS(number, message string) error {
	cmd := exec.Command("termux-sms-send", "-n", number, message)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to send SMS to %s: %v", number, err)
		return err
	}
	return nil
}
