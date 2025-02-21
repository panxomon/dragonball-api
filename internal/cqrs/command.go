package cqrs

import "context"

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

type CommandWithResultHandler[C any, R any] interface {
	Handle(ctx context.Context, c C) (R, error)
}
