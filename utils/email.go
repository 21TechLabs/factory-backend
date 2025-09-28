package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SendEmailRequestType string

const (
	SendEmailRequestTypePasswordReset SendEmailRequestType = "password-reset"
	SendEmailRequestTypeAlarmUpdates  SendEmailRequestType = "alarm updates"
)

type DynamicTemplateData struct {
	Subject string `json:"subject"`
}

type EmailRecepient struct {
	To                  string              `json:"to"`
	TemplateID          string              `json:"templateId,omitempty"`
	DynamicTemplateData DynamicTemplateData `json:"dynamicTemplateData"`
	Body                string              `json:"body"`
}

type SendMailRequest struct {
	Type       SendEmailRequestType `json:"type"`
	Recipients []EmailRecepient     `json:"recipients"`
}

type SendPasswordResetEmailRequest struct {
	To       string `json:"to"`
	ResetKey string `json:"resetKey"`
}

type EmailService struct {
	EmailAPI string
	AuthKey  string
	UIURL    string
	baseURL  string
	source   string
}

type SendMailBody struct {
	AuthKey         string               `json:"authKey"`
	RequestSource   string               `json:"requestSource"`
	RequestType     SendEmailRequestType `json:"requestType"`
	EmailRecepients []EmailRecepient     `json:"emailRecepients"`
}

func NewEmailService(emailAPI, authKey, uiURL string) *EmailService {
	return &EmailService{
		EmailAPI: emailAPI,
		AuthKey:  authKey,
		UIURL:    uiURL,
		baseURL:  fmt.Sprintf("%s/send-emails", emailAPI),
	}
}

func (es *EmailService) HttpClient() *http.Client {
	return &http.Client{Timeout: time.Second * 5} // Set appropriate timeout
}

func (es *EmailService) NewRequest(method string, url string, body interface{}) (*http.Request, error) {
	var reader *bytes.Buffer = bytes.NewBuffer(nil)

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reader = bytes.NewBuffer(jsonBody)
	}

	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create new HTTP request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", es.AuthKey)

	return request, nil
}

func (es *EmailService) SendMail(req SendMailRequest) error {

	body := SendMailBody{
		AuthKey:         es.AuthKey,
		RequestSource:   es.source,
		RequestType:     req.Type,
		EmailRecepients: req.Recipients,
	}

	request, err := es.NewRequest("POST", es.baseURL, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := es.HttpClient()
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send email request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send email, status code: %d", response.StatusCode)
	}
	return nil
}

func (es *EmailService) SendPasswordResetEmail(req SendPasswordResetEmailRequest) error {
	body := SendMailRequest{
		Type: SendEmailRequestTypePasswordReset,
		Recipients: []EmailRecepient{
			{
				To: req.To,
				DynamicTemplateData: DynamicTemplateData{
					Subject: "Password Reset Request",
				},
				Body: fmt.Sprintf("This is your ELIMS password reset email. Please go to %s/sessions/reset/%s to complete the reset.", es.UIURL, req.ResetKey),
			},
		},
	}
	return es.SendMail(body)
}

func (es *EmailService) SendMails(personalizedEmails []EmailRecepient) error {
	if personalizedEmails == nil {
		return fmt.Errorf("personalized emails cannot be nil")
	}
	body := SendMailRequest{
		Type:       SendEmailRequestTypeAlarmUpdates,
		Recipients: personalizedEmails,
	}
	return es.SendMail(body)
}
