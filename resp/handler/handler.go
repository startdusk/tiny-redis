package handler

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/startdusk/tiny-redis/api/db"
	"github.com/startdusk/tiny-redis/lib/logger"
	"github.com/startdusk/tiny-redis/lib/sync/atomic"
	"github.com/startdusk/tiny-redis/resp/conn"
	"github.com/startdusk/tiny-redis/resp/parser"
	"github.com/startdusk/tiny-redis/resp/reply"
)

var (
	ErrUnknownReply = []byte("-ERR unknown\r\n")
)

type emptyKey struct{}

func NewHandler(db db.Database) *Handler {
	return &Handler{
		db: db,
	}
}

type Handler struct {
	activeConn sync.Map
	db         db.Database
	closing    atomic.Boolean
}

// closes a single client
func (r *Handler) closeClient(client *conn.Conn) {
	client.Close()
	r.db.AfterClientClose(client)
	r.activeConn.Delete(client)
}

func (r *Handler) Handle(ctx context.Context, c net.Conn) {
	if r.closing.Get() {
		c.Close()
	}

	client := conn.NewConn(c)
	r.activeConn.Store(client, emptyKey{})
	stream := parser.ParseStream(c)
	for payload := range stream {
		// error
		if payload.Err != nil {
			// TODO: use of closed network connection
			if errors.Is(payload.Err, io.EOF) ||
				errors.Is(payload.Err, io.ErrUnexpectedEOF) ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				r.closeClient(client)
				logger.Info("conneciton closed: " + client.RemoteAddr().String())
				return
			}

			// protocol error
			errReply := reply.NewStandardErrReply(payload.Err.Error())
			err := client.Write(errReply.Bytes())
			if err != nil {
				r.closeClient(client)
				logger.Info("conneciton closed: " + client.RemoteAddr().String())
				return
			}
			continue
		}

		// exec
		if payload.Data == nil {
			continue
		}
		bulkReply, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk reply")
			continue
		}
		res := r.db.Exec(client, bulkReply.Args())
		if res == nil {
			client.Write(ErrUnknownReply)
		} else {
			client.Write(res.Bytes())
		}
	}
}

func (r *Handler) Close() error {
	logger.Info("handler shutting down")
	r.closing.Set(true)
	r.activeConn.Range(func(key, value any) bool {
		if client, ok := key.(*conn.Conn); ok {
			client.Close()
		}
		return true
	})
	r.db.Close()
	return nil
}
