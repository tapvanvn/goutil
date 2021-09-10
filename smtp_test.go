package goutil_test

import (
	"fmt"
	"testing"

	"github.com/tapvanvn/goutil"
)

func TestSendEmail(t *testing.T) {
	var fromEmail = ""
	var password = ""
	smtpServer := goutil.NewSmtpServer("smtp.gmail.com", "465", fromEmail, password)
	if err := smtpServer.SendEmail(fromEmail, "tapvanvn@yahoo.com", "hello", " this is a test email"); err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
