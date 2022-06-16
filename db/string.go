package db

import (
	dbent "github.com/startdusk/tiny-redis/api/db"
	"github.com/startdusk/tiny-redis/api/resp"
	"github.com/startdusk/tiny-redis/lib/utils"
	"github.com/startdusk/tiny-redis/resp/reply"
)

// GET
func execGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.NewNullBulkReply()
	}

	data, ok := entity.Data.([]byte)
	if !ok {
		return reply.NewStandardErrReply("key must be string type")
	}
	return reply.NewBulkReply(data)
}

// SET key1 val1
func execSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity := &dbent.DataEntity{
		Data: args[1],
	}
	db.PutEntity(key, entity)
	db.addAOF(utils.ToCmdLineWithCmdName("set", args...))
	return reply.NewOKReply()
}

// SETNX key1 val1
func execSetNX(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity := &dbent.DataEntity{
		Data: args[1],
	}
	res := db.PutIfAbsent(key, entity)
	db.addAOF(utils.ToCmdLineWithCmdName("setnx", args...))
	return reply.NewNumberReply(int64(res))
}

// GETSET key1 val1
func execGetSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	old, exists := db.GetEntity(key)

	entity := &dbent.DataEntity{
		Data: args[1],
	}
	db.PutEntity(key, entity)

	db.addAOF(utils.ToCmdLineWithCmdName("getset", args...))
	if exists {
		return reply.NewBulkReply(old.Data.([]byte))
	}
	return reply.NewNullBulkReply()
}

// SETLEN key1
func execStrLen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.NewNullBulkReply()
	}
	bytes := entity.Data.([]byte)
	return reply.NewNumberReply(int64(len(bytes)))
}
