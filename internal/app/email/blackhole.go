package email

import (
	"context"
)

type Blackhole struct{}

func (mg *Blackhole) Send(ctx context.Context, recipient, title, content string) error {
	return nil
}

func NewBlackhole() *Blackhole {
	return &Blackhole{}
}
