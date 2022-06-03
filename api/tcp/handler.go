package tcp

import (
	"context"
	"io"
)

type Handler interface {
	Handle(ctx context.Context, conn io.ReadWriteCloser)
	Close() error
}
