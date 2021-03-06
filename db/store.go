package db

import (
	"strconv"
	"strings"

	"github.com/startdusk/tiny-redis/aof"
	"github.com/startdusk/tiny-redis/api/resp"
	"github.com/startdusk/tiny-redis/lib/logger"
	"github.com/startdusk/tiny-redis/resp/reply"
	"github.com/startdusk/tiny-redis/structrue/dict"
)

const defaultDBNums = 16

type StandaloneDatabase struct {
	dbSet []*DB

	aof *aof.Handler
}

func NewStandaloneDatabase(dbNum int, aofFilename string, appendOnly bool) *StandaloneDatabase {
	if dbNum <= 0 {
		dbNum = defaultDBNums
	}
	db := StandaloneDatabase{
		dbSet: make([]*DB, dbNum),
	}
	for i := 0; i < dbNum; i++ {
		db.dbSet[i] = &DB{
			index:  i,
			data:   dict.NewSyncDict(),
			addAOF: func(cl CmdLine) {},
		}
	}
	if appendOnly {
		h, err := aof.NewAOFHandler(&db, aofFilename, appendOnly)
		if err != nil {
			panic(err)
		}
		db.aof = h
		for i := 0; i < dbNum; i++ {
			index := i
			db.dbSet[i].addAOF = func(cmdLine CmdLine) {
				db.aof.Add(index, cmdLine)
			}
		}
	}
	return &db
}

func (d *StandaloneDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	cmdName := strings.ToLower(string(args[0]))
	switch cmdName {
	case "select":
		if len(args) != 2 {
			return reply.NewArgNumErrReply(cmdName)
		}
		return execSelect(client, d, args[1:])
	default:
		index := client.GetDBIndex()
		db := d.dbSet[index]
		return db.Exec(client, args)
	}
}

func (d *StandaloneDatabase) Close() error {
	return nil
}

func (d *StandaloneDatabase) AfterClientClose(c resp.Connection) error {
	return nil
}

// select n
func execSelect(c resp.Connection, db *StandaloneDatabase, args [][]byte) resp.Reply {
	index, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.NewStandardErrReply("ERR invalid DB index")
	}
	if index < 0 || index >= len(db.dbSet) {
		return reply.NewStandardErrReply("ERR DB index is out of range")
	}

	c.SelectDB(index)

	return reply.NewOKReply()
}
