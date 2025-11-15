package email

import (
	"context"
	"errors"
	"fmt"

	"github.com/mailgun/mailgun-go/v4"
)

type MailgunEmailService struct {
	impl       *mailgun.MailgunImpl
	senderName string
}

func NewMailgunEmailService(domain, apiKey, senderName string, isEU bool) (*MailgunEmailService, error) {
	if apiKey == "" {
		return nil, errors.New("apiKey is empty")
	}

	if domain == "" {
		return nil, errors.New("domain is empty")
	}

	impl := mailgun.NewMailgun(domain, apiKey)
	if isEU {
		impl.SetAPIBase(mailgun.APIBaseEU)
	}

	return &MailgunEmailService{
		impl:       impl,
		senderName: senderName,
	}, nil
}

func (mg *MailgunEmailService) Send(ctx context.Context, recipient, title, content string) error {
	m := mailgun.NewMessage(
		fmt.Sprintf("%s <%s>", mg.senderName, "no-reply@"+mg.impl.Domain()),
		title,
		content,
		fmt.Sprintf("<%s>", recipient),
	)

	_, _, err := mg.impl.Send(ctx, m)
	return err
}
