package email

import (
	"context"
	"fmt"
)

type Console struct{}

func (mg *Console) Send(ctx context.Context, recipient, title, content string) error {
	fmt.Printf("EMAIL SENT!\nTo: %s\nTitle: %s\nContent:\n%s\n", recipient, title, content)
	return nil
}

func NewConsole() *Console {
	return &Console{}
}
