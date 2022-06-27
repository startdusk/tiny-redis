package cluster

import (
	"github.com/startdusk/tiny-redis/api/resp"
	"github.com/startdusk/tiny-redis/resp/reply"
)

// del k1 k2 k3 k4 k5
func del(cluster *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(c, cmdArgs)
	var errReply reply.ErrorReply
	var deleted int64
	for _, r := range replies {
		if reply.IsErrRely(r) {
			errReply = r.(reply.ErrorReply)
			break
		}
		intReply, ok := r.(*reply.NumberReply)
		if !ok {
			errReply = reply.NewStandardErrReply("error")
		}
		deleted += int64(intReply.Code)
	}
	if errReply == nil {
		return reply.NewNumberReply(deleted)
	}
	return reply.NewStandardErrReply("error: " + errReply.Error())
}
