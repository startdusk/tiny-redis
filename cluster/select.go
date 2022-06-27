package cluster

import "github.com/startdusk/tiny-redis/api/resp"

func sel(cluster *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	return cluster.db.Exec(c, cmdArgs)
}
