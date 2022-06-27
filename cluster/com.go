package cluster

import (
	"context"
	"errors"
	"strconv"

	"github.com/startdusk/tiny-redis/api/resp"
	"github.com/startdusk/tiny-redis/lib/conv"
	"github.com/startdusk/tiny-redis/resp/client"
	"github.com/startdusk/tiny-redis/resp/reply"
)

func (d *Database) getPeerClient(targetPeer string) (*client.Client, error) {
	pool, ok := d.peerConns[targetPeer]
	if !ok {
		return nil, errors.New("connection not found")
	}

	obj, err := pool.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}
	c, ok := obj.(*client.Client)
	if !ok {
		return nil, errors.New("wrong type")
	}
	return c, nil
}

func (d *Database) pushbackPeerClient(peer string, c *client.Client) error {
	pool, ok := d.peerConns[peer]
	if !ok {
		return errors.New("connection not found")
	}
	return pool.ReturnObject(context.Background(), c)
}

func (d *Database) relay(peer string, c resp.Connection, args [][]byte) resp.Reply {
	if peer == d.self {
		return d.db.Exec(c, args)
	}

	peerClient, err := d.getPeerClient(peer)
	if err != nil {
		return reply.NewStandardErrReply(err.Error())
	}
	defer d.pushbackPeerClient(peer, peerClient)
	cmd := conv.ToCmdLine("SELECT", strconv.Itoa(c.GetDBIndex()))
	peerClient.Send(cmd)
	return peerClient.Send(args)
}

func (d *Database) broadcast(c resp.Connection, args [][]byte) map[string]resp.Reply {
	replys := make(map[string]resp.Reply)
	for _, node := range d.nodes {
		replys[node] = d.relay(node, c, args)
	}
	return replys
}
