package server

import (
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/sirupsen/logrus"
	"github.com/stenic/ledger/internal/auth"
	"github.com/stenic/ledger/internal/pkg/messagebus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 5 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func wsHandler(authValidator auth.LedgerValidator, msgBuss messagebus.MessageBus) gin.HandlerFunc {
	logger := logrus.WithFields(logrus.Fields{
		"scope": "websockets",
	})

	hub := newHub()
	go hub.run()

	go func() {
		var pubsub <-chan *redis.Message

		for pubsub == nil {
			pubsub = msgBuss.Consume("newVersion")
			time.Sleep(1 * time.Second)
		}

		logger.Debug("Instance of pubsub configured, waiting for messages")
		for msg := range pubsub {
			logger.Trace("Refreshed version count")
			hub.broadcast <- []byte(msg.Payload)
		}
	}()

	return func(c *gin.Context) {
		conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
		if err != nil {
			logger.Warn(err)
		}
		logger.Debug("Client connected")
		client := &Client{
			hub:  hub,
			conn: conn,
			send: make(chan []byte, 256),
		}
		hub.register <- client

		go client.writeHandler()
	}
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	logger := logrus.WithField("scope", "websocket")
	for {
		select {
		case client := <-h.register:
			logger.WithField("id", client.conn.RemoteAddr()).Debug("Registering client")
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				logger.WithField("id", client.conn.RemoteAddr()).Debug("Unremoving client")
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			logger.Debug("Broadcasting message")
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

type Client struct {
	hub  *Hub
	conn net.Conn
	send chan []byte
}

func (c *Client) writeHandler() {
	logger := logrus.WithField("scope", "websocket")

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				return
			}
			err := wsutil.WriteServerMessage(c.conn, ws.OpText, message)
			if err != nil {
				if _, ok := err.(wsutil.ClosedError); ok {
					logger.Debug("Client disconnected")
				} else {
					logger.Error(err)
				}
				return
			}
		case <-ticker.C:
			logger.WithField("id", c.conn.RemoteAddr()).Trace("Sending ping")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := wsutil.WriteServerMessage(c.conn, ws.OpPing, nil); err != nil {
				return
			}
		}
	}
}
