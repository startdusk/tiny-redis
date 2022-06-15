package db

import (
	"github.com/startdusk/tiny-redis/api/resp"
	"github.com/startdusk/tiny-redis/lib/wildcard"
	"github.com/startdusk/tiny-redis/resp/reply"
)

// DEL key1 key2 key3...
func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	deleted := db.Removes(keys...)
	return reply.NewNumberReply(int64(deleted))
}

// EXISTS key1 key2 key3...
func execExists(db *DB, args [][]byte) resp.Reply {
	var exists int64
	for _, arg := range args {
		key := string(arg)
		_, have := db.GetEntity(key)
		if have {
			exists++
		}
	}
	return reply.NewNumberReply(exists)
}

// FLUSHDB
func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	return reply.NewOKReply()
}

// TYPE key1
func execType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.NewStatusReply("none") // TCP :none\r\n
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.NewStatusReply("string")
		// TODO: another type support...
	}
	return reply.NewUnknowErrReply()
}

// RENAME key1 key2 -> change key1 name to key2 name
func execRename(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dst := string(args[1])

	entity, exists := db.GetEntity(src)
	if !exists {
		return reply.NewStandardErrReply("no such key " + src)
	}
	db.PutEntity(dst, entity)
	db.Remove(src)
	return reply.NewOKReply()
}

// RENAMENX key1 key2 -> change key1 name to key2 name if key2 not exists
func execRenameNX(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dst := string(args[1])

	_, ok := db.GetEntity(dst)
	if ok {
		return reply.NewNumberReply(0)
	}
	entity, exists := db.GetEntity(src)
	if !exists {
		return reply.NewStandardErrReply("no such key " + src)
	}
	db.PutEntity(dst, entity)
	db.Remove(src)
	return reply.NewNumberReply(1)
}

// KEYS *
func execKeys(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.data.Range(func(key string, _ any) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})

	return reply.NewMultiBulkReply(result)
}
