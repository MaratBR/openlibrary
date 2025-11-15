package email

import "context"

type Service interface {
	Send(
		ctx context.Context,
		recipient, title, content string,
	) error
}
