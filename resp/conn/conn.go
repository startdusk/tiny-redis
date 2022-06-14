package conn

import (
	"net"
	"sync"
	"time"

	"github.com/startdusk/tiny-redis/lib/sync/wait"
)

type Conn struct {
	conn         net.Conn
	waitingReply wait.Wait
	mu           sync.Mutex
	selectedDB   int
}

func NewConn(conn net.Conn) *Conn {
	return &Conn{
		conn: conn,
	}
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Conn) Write(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	c.mu.Lock()
	c.waitingReply.Add(1)
	defer func() {
		c.waitingReply.Done()
		c.mu.Unlock()
	}()

	_, err := c.conn.Write(bytes)
	return err
}

func (c *Conn) GetDBIndex() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.selectedDB
}

func (c *Conn) SelectDB(dbNum int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.selectedDB = dbNum
}

func (c *Conn) Close() error {
	c.waitingReply.WaitWithTimeout(10 * time.Second)
	return c.conn.Close()
}
