package tcp

import (
	"bufio"
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/startdusk/tiny-redis/lib/logger"
	"github.com/startdusk/tiny-redis/lib/sync/atomic"
	"github.com/startdusk/tiny-redis/lib/sync/wait"
)

type EchoClient struct {
	Conn    io.ReadWriteCloser
	Waiting wait.Wait
}

func (c *EchoClient) Close() error {
	c.Waiting.WaitWithTimeout(10 * time.Second)
	return c.Conn.Close()
}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

func (h *EchoHandler) Handle(ctx context.Context, conn io.ReadWriteCloser) {
	if h.closing.Get() {
		conn.Close()
		return
	}

	client := &EchoClient{
		Conn: conn,
	}
	h.activeConn.Store(client, struct{}{})
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.Info("connectoin close")
				h.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}

		client.Waiting.Add(1)
		b := []byte(msg)
		conn.Write(b)
		client.Waiting.Done()
	}
}

func (h *EchoHandler) Close() error {
	logger.Info("handler shutting down")
	h.closing.Set(true)
	h.activeConn.Range(func(key, _ any) bool {
		if conn, ok := key.(*EchoClient); ok {
			conn.Close()
		}
		// return true for range any
		return true
	})
	return nil
}
