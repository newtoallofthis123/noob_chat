package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/newtoallofthis123/ranhash"
)

type Server struct {
	rooms     Rooms
	db        *DbInstance
	logWriter io.Writer
}

type Room struct {
	id    string
	name  string
	conns map[*websocket.Conn]User
}

func NewRoom(id string) Room {
	return Room{
		id:    id,
		conns: make(map[*websocket.Conn]User),
	}
}

type Rooms struct {
	rooms map[string]Room
}

func NewRooms() Rooms {
	return Rooms{
		rooms: map[string]Room{},
	}
}

func (r *Rooms) GetRoom(id string) Room {
	_, exists := r.rooms[id]
	if !exists {
		r.rooms[id] = NewRoom(id)
	}
	return r.rooms[id]
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewServer() *Server {
	return &Server{
		rooms:     NewRooms(),
		db:        NewDbInstance(),
		logWriter: os.Stderr,
	}
}

func NewServerWithLog(file *os.File) *Server {
	return &Server{
		rooms:     NewRooms(),
		db:        NewDbInstance(),
		logWriter: file,
	}
}

func (s *Server) StartServer(addr string) {
	r := gin.Default()

	r.GET("/", s.handleVersion)
	r.POST("/register", s.handleRegister)
	r.POST("/login", s.handleLogin)

	auth := r.Group("/api")
	auth.Use(s.authMiddleWare())
	auth.GET("/echo", s.handleAuthEcho)

	room := auth.Group("/chat")
	room.Use(s.roomMiddleWare())
	room.GET("/:roomId", s.handleChat)

	s.log("Starting server with addr: " + addr)
	r.Run(addr)
}

func (s *Server) handleVersion(c *gin.Context) {
	c.String(http.StatusOK, "NoobChat v0.0.1")
}

func (s *Server) log(msg string) {
	fmt.Fprintf(s.logWriter, "%s - %s", time.Now().String(), msg)
}

type Message struct {
	Content   string `json:"content,omitempty"`
	UserId    string `json:"user_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

func (s *Server) clientLoop(room Room, ws *websocket.Conn) {
	var msg Message
	for {
		err := ws.ReadJSON(&msg)
		if err != nil {
			ws.WriteJSON(gin.H{"Error": err})
			continue
		}

		//TODO: Add in json validation for the userid and sanitize the content
		var req = CreateChatRequest{
			Content: msg.Content,
			UserId:  msg.UserId,
			RoomId:  room.id,
		}

		id, err := s.db.CreateChat(req)
		if err != nil {
			ws.WriteJSON(gin.H{"Error": err})
			continue
		}
		s.log("Msg with id: " + id + " created")
		s.broadcast(msg, room)
	}
}

func (s *Server) broadcast(msg Message, room Room) {
	for conn := range room.conns {
		conn.WriteJSON(msg)
	}
}

func (s *Server) roomMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomId := c.Param("roomId")
		if roomId == "" {
			roomId = ranhash.RanHash(8)
		}

		room := s.rooms.GetRoom(roomId)
		c.Set("room", room)
	}
}

func (s *Server) handleChat(c *gin.Context) {
	session, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authenticated Route"})
		return
	}

	roomObj, exists := c.Get("room")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authenticated Route"})
		return
	}

	room := roomObj.(Room)

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room.conns[ws] = session.(Session).User

	ws.WriteJSON(map[string]string{"msg": "Connected to room: " + room.id})
	go func(ws *websocket.Conn) {
		s.clientLoop(room, ws)
	}(ws)
}
