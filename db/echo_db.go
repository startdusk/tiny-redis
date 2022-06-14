package db

import (
	"github.com/startdusk/tiny-redis/api/resp"
	"github.com/startdusk/tiny-redis/resp/reply"
)

type EchoDB struct{}

func (e EchoDB) Exec(client resp.Connection, args [][]byte) resp.Reply {
	return reply.NewMultiBulkReply(args)
}

func (e EchoDB) Close() error {
	return nil
}

func (e EchoDB) AfterClientClose(c resp.Connection) error {
	return nil
}
