package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn    *websocket.Conn
	Message chan *Message
	Id      string `json:"id"`
	RoomId  string `json:"roomId"`
	User    string `json:"user"`
}

type Message struct {
	Content    string `json:"content"`
	RoomId     string `json:"roomId"`
	User       string `json:"user"`
	Playground string `json:"playground"`
}

func (c *Client) WriteMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}

			break
		}

		UpdatePlayground(m, c.RoomId, hub)

		msg := &Message{
			Content:    string(m),
			RoomId:     c.RoomId,
			User:       c.User,
			Playground: hub.Rooms[c.RoomId].Playground,
		}

		hub.Broadcast <- msg
	}
}
