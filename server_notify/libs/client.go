package libs

import (
	"net"
	"sync"

	"github.com/gobwas/ws/wsutil"
)

type Client struct {
	Lock       *sync.Mutex
	Connection net.Conn
}

func NewClient(connection net.Conn) *Client {
	return &Client{Lock: &sync.Mutex{}, Connection: connection}
}

func (c *Client) Write(message []byte) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	wsutil.WriteServerText(c.Connection, message)
}
