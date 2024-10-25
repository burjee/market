package libs

import (
	"errors"
	"log"
	"net"
	"time"

	"github.com/panjf2000/ants/v2"
	"golang.org/x/sys/unix"
)

type BroadcastPack struct {
	HandleFunc func(*Client) func()
}

type Hub struct {
	epoll   *Epoll
	go_pool *ants.Pool

	connections map[net.Conn]*Client

	Register   chan *Client
	Unregister chan net.Conn
	Broadcast  chan *BroadcastPack
}

func NewHub(epoll *Epoll, go_pool *ants.Pool) *Hub {
	return &Hub{
		epoll:   epoll,
		go_pool: go_pool,

		connections: make(map[net.Conn]*Client),

		Register:   make(chan *Client, 1000),
		Unregister: make(chan net.Conn, 1000),
		Broadcast:  make(chan *BroadcastPack, 1000),
	}
}

func (h *Hub) Run() {
	ticker := time.NewTicker(time.Minute)

	for {
		select {

		case <-ticker.C:
			for _, client := range h.connections {
				h.go_pool.Submit(func() { client.Write([]byte("ping")) })
			}

		case client := <-h.Register:
			h.connections[client.Connection] = client

		case connection := <-h.Unregister:
			if client, ok := h.connections[connection]; ok {
				h.removeFromEpoll(client.Connection)
				h.removeConnection(client)
			}

		case broadcast_pack := <-h.Broadcast:
			for _, client := range h.connections {
				h.go_pool.Submit(broadcast_pack.HandleFunc(client))
			}
		}
	}
}

func (h *Hub) removeFromEpoll(connection net.Conn) {
	max_retries := 10
	for retries := 0; retries < max_retries; retries += 1 {
		if err := h.epoll.Remove(connection); err == nil || errors.Is(err, unix.ENOENT) || errors.Is(err, unix.EBADF) {
			break
		} else {
			log.Printf("Failed to remove from epoll %v (attempt %d)", err, retries+1)
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func (h *Hub) removeConnection(client *Client) {
	delete(h.connections, client.Connection)
	client.Connection.Close()
}

func (h *Hub) Close() {
	for connection := range h.connections {
		connection.Close()
	}
}
