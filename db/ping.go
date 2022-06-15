package db

import (
	"github.com/startdusk/tiny-redis/api/resp"
	"github.com/startdusk/tiny-redis/resp/reply"
)

func Ping(db *DB, args [][]byte) resp.Reply {
	return reply.NewPongReply()
}
