package pipe

import "context"

type channelPipe[T any] struct {
	ch chan T
}

func (p *channelPipe[T]) Read(ctx context.Context, dest []T) int {

}

func (p *channelPipe[T]) Write(ctx context.Context, src []T) error {
	for _, item := range src {
		p.ch <- item
	}
}

func NewChannelPipe[T any]() Pipe[T] {
	return &channelPipe[T]{
		ch: make(chan T),
	}
}
