package cluster

import (
	"context"
	"strings"

	pool "github.com/jolestar/go-commons-pool"
	"github.com/startdusk/tiny-redis/api/db"
	"github.com/startdusk/tiny-redis/api/resp"
	"github.com/startdusk/tiny-redis/config"
	database "github.com/startdusk/tiny-redis/db"
	"github.com/startdusk/tiny-redis/lib/consistenthash"
	"github.com/startdusk/tiny-redis/lib/logger"
	"github.com/startdusk/tiny-redis/resp/reply"
)

type Database struct {
	self       string
	nodes      []string
	peerPicker *consistenthash.NodeMap
	peerConns  map[string]*pool.ObjectPool
	db         db.Database
}

func NewDatabase() *Database {
	cluster := Database{
		self:      config.Properties.Self,
		peerConns: make(map[string]*pool.ObjectPool),
		db: database.NewStandaloneDatabase(
			config.Properties.Databases,
			config.Properties.AppendFilename,
			config.Properties.AppendOnly),
		peerPicker: consistenthash.NewNodeMap(nil),
	}
	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	nodes = append(nodes, config.Properties.Peers...)
	nodes = append(nodes, config.Properties.Self)
	cluster.nodes = nodes
	cluster.peerPicker.Add(nodes...)
	ctx := context.Background()
	for _, peer := range config.Properties.Peers {
		cluster.peerConns[peer] = pool.NewObjectPoolWithDefaultConfig(
			ctx, &connection{
				Peer: peer,
			})
	}
	return &cluster
}

type CmdFunc func(cluster *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply

var router = NewRouter()

func (d *Database) Exec(client resp.Connection, args [][]byte) (result resp.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			result = reply.NewUnknowErrReply()
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	cmdFunc, ok := router[cmdName]
	if !ok {
		return reply.NewStandardErrReply("not supported command " + cmdName)
	}

	return cmdFunc(d, client, args)
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) AfterClientClose(c resp.Connection) error {
	return d.db.AfterClientClose(c)
}
