package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) handleRegister(c *gin.Context) {
	s.log("Starting Register Handle")
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := s.db.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		s.log("Unable to create user err: " + err.Error())
		return
	}

	sessionId, err := s.db.CreateSession(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		s.log("Unable to create session err: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"session_id": sessionId})
}

func (s *Server) handleLogin(c *gin.Context) {
	s.log("Starting Login Handle")
	var req LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, actualPassword, err := s.db.GetUserPassword(req.Username)
	if !MatchPasswords(req.Password, actualPassword) || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		s.log("Passwords don't match")
		return
	}

	sessionId, err := s.db.CreateSession(user.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		s.log("Unable to create session err: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"session_id": sessionId})
}

func (s *Server) authMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId := c.GetHeader("session_id")
		if sessionId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Session Id Needed"})
			return
		}

		session, err := s.db.GetSession(sessionId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Set("session", session)
		c.Next()
	}
}

func (s *Server) handleAuthEcho(c *gin.Context) {
	session, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authenticated Route"})
	}

	c.JSON(http.StatusOK, session)
}
