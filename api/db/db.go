package db

import "github.com/startdusk/tiny-redis/api/resp"

type CmdLine = [][]byte

type Database interface {
	Exec(client resp.Connection, args [][]byte) resp.Reply
	Close() error
	AfterClientClose(c resp.Connection) error
}

type DataEntity struct {
	Data any
}
