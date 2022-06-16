package db

import (
	"strings"

	"github.com/startdusk/tiny-redis/api/db"
	"github.com/startdusk/tiny-redis/api/resp"
	"github.com/startdusk/tiny-redis/model/dict"
	"github.com/startdusk/tiny-redis/resp/reply"
)

type DB struct {
	index  int
	data   dict.Dict
	addAOF func(CmdLine)
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply

type CmdLine = [][]byte

func (d *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply {
	if len(cmdLine) == 0 {
		return reply.NewStandardErrReply("ERR none command")
	}
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.NewStandardErrReply("ERR unknown command " + cmdName)
	}

	// eg. expect send: [SET KEY VALUE] but got [SET KEY]
	if !validateArity(cmd.arity, cmdLine) {
		return reply.NewArgNumErrReply(cmdName)
	}

	exec := cmd.exector
	// SET KEY VALUE -> KEY VALUE
	return exec(d, cmdLine[1:])
}

func (d *DB) GetEntity(key string) (*db.DataEntity, bool) {
	raw, ok := d.data.Get(key)
	if !ok {
		return nil, false
	}
	entity, ok := raw.(*db.DataEntity)
	if !ok {
		return nil, false
	}
	return entity, true
}

func (d *DB) PutEntity(key string, entity *db.DataEntity) int {
	return d.data.Put(key, entity)
}

func (d *DB) PutIfExists(key string, entity *db.DataEntity) int {
	return d.data.PutIfExists(key, entity)
}

func (d *DB) PutIfAbsent(key string, entity *db.DataEntity) int {
	return d.data.PutIfAbsent(key, entity)
}

func (d *DB) Remove(key string) {
	d.data.Remove(key)
}

func (d *DB) Removes(keys ...string) (deleted int) {
	for _, key := range keys {
		_, exists := d.data.Get(key)
		if exists {
			d.Remove(key)
			deleted++
		}
	}
	return
}

func (d *DB) Flush() {
	d.data.Clear()
}

func (d *DB) Index() int {
	return d.index
}

// SET KEY VALUE -> arity = 3
// EXISTS KEY1 KEY2 KEY3...KEYn -> artiy = -n
func validateArity(arity int, cmdArgs [][]byte) bool {
	argNum := len(cmdArgs)
	if arity >= 0 {
		return argNum == arity
	}
	return argNum >= -arity
}
