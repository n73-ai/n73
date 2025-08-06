package handlers

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	isClosing bool
	mu        sync.Mutex
	UserID    string
}

type RegisterPayload struct {
	Conn   *websocket.Conn
	Client *Client
}

var clients = make(map[*websocket.Conn]*Client)
var register = make(chan *RegisterPayload)
var broadcast = make(chan string)
var unregister = make(chan *websocket.Conn)

func SendToUser(userID string, msg string) {
	for conn, client := range clients {
		if client.UserID == userID {
			client.mu.Lock()
			defer client.mu.Unlock()
			if client.isClosing {
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				client.isClosing = true
				conn.Close()
				unregister <- conn
			}
		}
	}
}

func RunHub() {
	for {
		select {
		case payload := <-register:
			clients[payload.Conn] = payload.Client

		case message := <-broadcast:
			for conn, c := range clients {
				go func(conn *websocket.Conn, c *Client) {
					c.mu.Lock()
					defer c.mu.Unlock()
					if c.isClosing {
						return
					}
					if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
						c.isClosing = true
						conn.WriteMessage(websocket.CloseMessage, []byte{})
						conn.Close()
						unregister <- conn
					}
				}(conn, c)
			}

		case conn := <-unregister:
			delete(clients, conn)
		}
	}
}

func FeedChat(c *websocket.Conn) {
	userID := c.Query("user_id")

	client := &Client{UserID: userID}
	register <- &RegisterPayload{Conn: c, Client: client}

	defer func() {
		unregister <- c
		c.Close()
	}()

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}
