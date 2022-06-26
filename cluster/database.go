package cluster

import (
	"context"

	pool "github.com/jolestar/go-commons-pool"
	"github.com/startdusk/tiny-redis/api/db"
	"github.com/startdusk/tiny-redis/api/resp"
	"github.com/startdusk/tiny-redis/config"
	database "github.com/startdusk/tiny-redis/db"
	"github.com/startdusk/tiny-redis/lib/consistenthash"
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
		pool.NewObjectPoolWithDefaultConfig(
			ctx, &connection{
				Peer: peer,
			})
	}
	return &cluster
}

func (d *Database) Exec(client resp.Connection, args [][]byte) resp.Reply {
	panic("not implemented") // TODO: Implement
}

func (d *Database) Close() error {
	panic("not implemented") // TODO: Implement
}

func (d *Database) AfterClientClose(c resp.Connection) error {
	panic("not implemented") // TODO: Implement
}
