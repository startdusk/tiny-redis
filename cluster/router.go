package cluster

import (
	"github.com/startdusk/tiny-redis/api/resp"
)

func NewRouter() map[string]CmdFunc {
	routerMap := make(map[string]CmdFunc)
	routerMap["exists"] = defaultCmdFunc // exists k1
	routerMap["type"] = defaultCmdFunc   // type k1
	routerMap["set"] = defaultCmdFunc
	routerMap["setnx"] = defaultCmdFunc
	routerMap["get"] = defaultCmdFunc
	routerMap["getset"] = defaultCmdFunc
	routerMap["ping"] = ping
	routerMap["rename"] = rename
	routerMap["renamenx"] = rename
	routerMap["flushdb"] = flushdb
	routerMap["del"] = del
	routerMap["select"] = sel
	return routerMap
}

func defaultCmdFunc(cluster *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	key := string(cmdArgs[1])
	peer := cluster.peerPicker.Pick(key)
	return cluster.relay(peer, c, cmdArgs)
}
