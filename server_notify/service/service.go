package service

import (
	"errors"
	"log"
	"net"
	"server_notify/libs"

	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"golang.org/x/sys/unix"
)

type service struct {
	epoll *libs.Epoll
	hub   *libs.Hub
}

func New(epoll *libs.Epoll, hub *libs.Hub) *service {
	return &service{epoll, hub}
}

func (s *service) initWebsocketError(connection net.Conn, error_text string) {
	connection.Close()
	log.Println("Failed to " + error_text)
}

func (s *service) WebsocketHandler(c *gin.Context) {
	connection, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
	if err != nil {
		s.initWebsocketError(connection, "upgrade connection "+err.Error())
		return
	}

	if err := s.epoll.Add(connection); err != nil {
		s.initWebsocketError(connection, "add connection "+err.Error())
		return
	}

	client := libs.NewClient(connection)
	s.hub.Register <- client
}

func (s *service) StartReadWebsocket() {
	for {

		connections, err := s.epoll.Wait()
		if err != nil {
			if !errors.Is(err, unix.EINTR) {
				log.Printf("Failed to epoll wait %v", err)
			}
			continue
		}

		for _, connection := range connections {
			if connection == nil {
				break
			}

			if _, op, err := wsutil.ReadClientData(connection); err != nil {
				s.hub.Unregister <- connection
			} else if op == ws.OpClose {
				s.hub.Unregister <- connection
			}
		}
	}
}
