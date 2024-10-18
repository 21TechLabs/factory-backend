package notifications

import (
	"fmt"
	"net/smtp"

	"github.com/21TechLabs/factory-be/utils"
)

var emailFrom = utils.GetEnv("EMAIL_FROM", false)
var emailHOST = utils.GetEnv("EMAIL_HOST", false)
var emailPORT = utils.GetEnv("EMAIL_PORT", false)
var emailUSER = utils.GetEnv("EMAIL_USER", false)
var emailPASSWORD = utils.GetEnv("EMAIL_PASSWORD", false)
var emailAuth = smtp.PlainAuth("", emailUSER, emailPASSWORD, emailHOST)

func SendEmail(_to string, _subject string, _body string) error {
	var to = []string{_to}

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, _subject, _body))

	err := smtp.SendMail(fmt.Sprintf("%s:%s", emailHOST, emailPORT), emailAuth, emailFrom, to, msg)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
