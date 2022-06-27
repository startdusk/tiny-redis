package cluster

import (
	"github.com/startdusk/tiny-redis/api/resp"
	"github.com/startdusk/tiny-redis/resp/reply"
)

func flushdb(cluster *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(c, cmdArgs)
	var errReply reply.ErrorReply
	for _, r := range replies {
		if reply.IsErrRely(r) {
			errReply = r.(reply.ErrorReply)
			break
		}
	}
	if errReply == nil {
		return reply.NewOKReply()
	}
	return reply.NewStandardErrReply("error: " + errReply.Error())
}
