package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Server struct {
	conns     map[*websocket.Conn]User
	db        *DbInstance
	logWriter io.Writer
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewServer() *Server {
	return &Server{
		conns:     make(map[*websocket.Conn]User),
		db:        NewDbInstance(),
		logWriter: os.Stderr,
	}
}

func NewServerWithLog(file *os.File) *Server {
	return &Server{
		conns:     make(map[*websocket.Conn]User),
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
	auth.GET("/chat", s.handleChat)

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

func (s *Server) clientLoop(ws *websocket.Conn) {
	var msg Message
	for {
		err := ws.ReadJSON(&msg)
		if err != nil {
			ws.WriteJSON(gin.H{"Error": err})
			continue
		}
	}
}

func (s *Server) handleChat(c *gin.Context) {
	session, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authenticated Route"})
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.conns[ws] = session.(Session).User

	ws.WriteJSON(map[string]string{"msg": "Connected to the server!"})
	go func(ws *websocket.Conn) {
		s.clientLoop(ws)
	}(ws)
}
