package cluster

import (
	"github.com/startdusk/tiny-redis/api/resp"
	"github.com/startdusk/tiny-redis/resp/reply"
)

// rename k1 k2
func rename(cluster *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.NewStandardErrReply("Err wrong number args")
	}

	src := string(cmdArgs[1])
	dst := string(cmdArgs[2])

	srcPeer := cluster.peerPicker.Pick(src)
	dstPeer := cluster.peerPicker.Pick(dst)

	if srcPeer != dstPeer {
		return reply.NewStandardErrReply("Err rename must within on peer")
	}
	return cluster.relay(srcPeer, c, cmdArgs)
}
