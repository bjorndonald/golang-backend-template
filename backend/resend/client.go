package resend

import (
	"github.com/resend/resend-go/v2"
)

type Client struct {
	resend *resend.Client
}

type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

func NewClient(apiKey string) *Client {
	client := resend.NewClient(apiKey)
	return &Client{
		resend: client,
	}
}

func (c *Client) Send(emails []string, from, fromName, subject, content string) (string, error) {

	params := &resend.SendEmailRequest{
		From:    from,
		To:      emails,
		Subject: subject,
		Html:    content,
	}
	res, err := c.resend.Emails.Send(params)
	if err != nil {
		return "", err
	}

	return res.Id, nil
}
