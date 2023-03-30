package main

import (
	"fmt"

	"github.com/sendgrid/rest"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var (
	SG_HOST = "https://api.sendgrid.com"
)

type SendGridClient struct {
}

func NewSendGridClient() *SendGridClient {
	return &SendGridClient{}
}

func (*SendGridClient) TriggerWebhookTest() (*rest.Response, error) {
	url := "/v3/user/webhooks/event/test"
	request := sendgrid.GetRequest(SG_API_KEY, url, SG_HOST)
	request.Method = "POST"
	request.Body = []byte(`{
  "url": "mollit non ipsum magna",
}`)
	response, err := sendgrid.API(request)
	if err != nil {
		return nil, fmt.Errorf("failed TriggerWebhookTest: %v", err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	return response, nil
}

func (*SendGridClient) SendEmail() (*rest.Response, error) {
	from := mail.NewEmail("Example User", "truongan.phan@manabie.com")
	subject := "Sending with Twilio SendGrid is Fun"
	// to := mail.NewEmail("Example User", "panhhuu@gmail.com")
	to := mail.NewEmail("Example User", "vinhnguyenhoang96@@manabie.com")
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go. <a href=\"https://manabie.com\">click here</a></strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(SG_API_KEY)
	response, err := client.Send(message)
	if err != nil {
		return nil, fmt.Errorf("failed SendEmail: %v", err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	return response, nil
}
