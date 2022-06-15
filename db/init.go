package db

func init() {
	// ping
	RegisterCommand("ping", Ping, 1)

	// keys
	{
		// del key1 key2 key3...keyn -> artiy = -n
		// must be 2 args and variable length
		RegisterCommand("del", execDel, -2)
		// exists key1 key2 key3...keyn -> artiy = -n
		// must be 2 args and variable length
		RegisterCommand("exists", execExists, -2)
		// flushdb
		RegisterCommand("flushdb", execFlushDB, 1)
		// type key1
		RegisterCommand("type", execType, 2)
		// rename key1 key2
		RegisterCommand("rename", execRename, 3)
		// renamenx key1 key2
		RegisterCommand("renamenx", execRenameNX, 3)
		// keys *
		RegisterCommand("keys", execKeys, 2)
	}

	// string
	{
		// get key1
		RegisterCommand("get", execGet, 2)
		// set key1 val1
		RegisterCommand("set", execSet, 3)
		// setnx key1 val1
		RegisterCommand("setnx", execSetNX, 3)
		// getset key1 val1
		RegisterCommand("getset", execGetSet, 3)
		// strlen key1
		RegisterCommand("strlen", execStrLen, 2)
	}
}
