package notifications

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/21TechLabs/factory-backend/utils"
)

var emailFrom string
var emailHOST string
var emailPORT string
var emailUSER string
var emailPASSWORD string
var emailAuth smtp.Auth

func init() {
	utils.LoadEnv()

	emailFrom = utils.GetEnv("EMAIL_FROM", false)
	emailHOST = utils.GetEnv("EMAIL_HOST", false)
	emailPORT = utils.GetEnv("EMAIL_PORT", false)
	emailUSER = utils.GetEnv("EMAIL_USER", false)
	emailPASSWORD = utils.GetEnv("EMAIL_PASSWORD", false)
	emailAuth = smtp.PlainAuth("", emailUSER, emailPASSWORD, emailHOST)
}

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

// Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)

	if err := smtp.SendMail(fmt.Sprintf("%s:%s", emailHOST, emailPORT), emailAuth, emailFrom, r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
