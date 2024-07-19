package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

type CreateRoomReq struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(c echo.Context) error {
	var req CreateRoomReq
	if err := c.Bind(&req); err != nil {
		fmt.Println("bad request")

		return c.JSON(http.StatusBadRequest, err.Error())
	}

	h.hub.Rooms[req.Id] = &Room{
		Id:      req.Id,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	return c.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// origin := r.Header.Get("Origin")
		return true
	},
}

func (h *Handler) JoinRoom(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	roomId := c.Param("roomId")
	clientId := c.QueryParam("userId")
	user := c.QueryParam("user")

	cl := &Client{
		Conn:    conn,
		Message: make(chan *Message, 10),
		Id:      clientId,
		RoomId:  roomId,
		User:    user,
	}

	m := &Message{
		Content: "A new user join the room",
		RoomId:  roomId,
		User:    user,
	}

	// rigester a new client through the register channel
	h.hub.Register <- cl

	// broadcast that message
	h.hub.Broadcast <- m

	go cl.WriteMessage()
	cl.ReadMessage(h.hub)

	return c.JSON(http.StatusOK, nil)
}

type RoomResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetRooms(c echo.Context) error {
	rooms := make([]RoomResponse, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomResponse{
			Id:   r.Id,
			Name: r.Name,
		})
	}

	return c.JSON(http.StatusOK, rooms)
}

type ClientResponse struct {
	Id   string `json:"id"`
	User string `json:"user"`
}

func (h *Handler) GetClients(c echo.Context) error {
	var Clients []ClientResponse

	roomdId := c.Param("roomId")

	if _, ok := h.hub.Rooms[roomdId]; !ok {
		Clients = make([]ClientResponse, 0)

		return c.JSON(http.StatusOK, Clients)
	}

	for _, cl := range h.hub.Rooms[roomdId].Clients {
		Clients = append(Clients, ClientResponse{
			Id:   cl.Id,
			User: cl.User,
		})
	}

	return c.JSON(http.StatusOK, Clients)
}
