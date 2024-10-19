package templates

import (
	"os"
	"path/filepath"

	"github.com/21TechLabs/factory-be/config"
	"github.com/21TechLabs/factory-be/notifications"
)

func GenericParse(r *notifications.Request, template string, body interface{}) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	return r.ParseTemplate(filepath.Join(wd, "notifications", "templates", template), body)
}

type WelcomeMessage struct {
	Name      string
	BrandName string
}

func (msg *WelcomeMessage) ParseAsHTML(r *notifications.Request) error {
	if len(msg.BrandName) == 0 {
		msg.BrandName = config.Name
	}
	return GenericParse(r, "html/welcome.html", msg)
}

type GoodbyeMessage WelcomeMessage

func (msg *GoodbyeMessage) ParseAsHTML(r *notifications.Request) error {
	if len(msg.BrandName) == 0 {
		msg.BrandName = config.Name
	}
	return GenericParse(r, "html/goodbye.html", msg)
}

type ResetPasswordMessage struct {
	Name      string
	BrandName string
	Link      string
}

func (msg *ResetPasswordMessage) ParseAsHTML(r *notifications.Request) error {
	if len(msg.BrandName) == 0 {
		msg.BrandName = config.Name
	}
	return GenericParse(r, "html/reset-password.html", msg)
}
