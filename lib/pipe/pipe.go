package pipe

import "context"

type PipeConsumer[T any] interface {
	Read(ctx context.Context, dest []T) int
}

type PipeProducer[T any] interface {
	Write(ctx context.Context, src []T) error
}

type Pipe[T any] interface {
	PipeConsumer[T]
	PipeProducer[T]

	Close() error
}
