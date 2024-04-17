package main

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/newtoallofthis123/ranhash"
)

type DbInstance struct {
	db *sql.DB
}

func NewDbInstance() *DbInstance {
	db, err := sql.Open("postgres", GetDbUrl())
	if err != nil {
		panic("Unable to establish conn with db")
	}

	query := `
	CREATE TABLE IF NOT EXISTS users(
		id TEXT PRIMARY KEY,
		name TEXT,
		username TEXT,
		password TEXT,
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS sessions(
		id TEXT PRIMARY KEY,
		user_id TEXT REFERENCES users("id"),
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS chats(
		id TEXT PRIMARY KEY,
		content TEXT,
		user_id TEXT REFERENCES users("id"),
		room_id TEXT,
		created_at TIMESTAMP DEFAULT NOW()
	);
	`

	_, err = db.Exec(query)
	if err != nil {
		panic("Unable to Init the DB")
	}

	return &DbInstance{
		db,
	}
}

type User struct {
	Id        string `json:"user"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type CreateUserRequest struct {
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginUserRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (pq *DbInstance) CreateUser(req CreateUserRequest) (string, error) {
	query := `
	INSERT INTO users(id, name, username, password)
	VALUES ($1, $2, $3, $4);
	`

	id := ranhash.RanHash(16)
	password, err := HashPassword(req.Password)
	if err != nil {
		return "", err
	}
	_, err = pq.db.Exec(query, id, req.Name, req.Username, password)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (pq *DbInstance) GetUser(id string) (User, error) {
	query := `
	SELECT * from users where id=$1;
	`

	var user User
	rows := pq.db.QueryRow(query, id)
	password := ""
	err := rows.Scan(&user.Id, &user.Name, &user.Username, &password, &user.CreatedAt)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (pq *DbInstance) GetUserPassword(username string) (User, string, error) {
	query := `
	SELECT * from users where username=$1;
	`

	var user User
	rows := pq.db.QueryRow(query, username)
	password := ""
	err := rows.Scan(&user.Id, &user.Name, &user.Username, &password, &user.CreatedAt)
	if err != nil {
		return User{}, "", err
	}

	return user, password, nil
}

type Session struct {
	Id        string `json:"id,omitempty"`
	User      User   `json:"user,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

func (pq *DbInstance) CreateSession(userId string) (string, error) {
	query := `
	INSERT INTO sessions(id, user_id)
	VALUES ($1, $2);
	`

	id := ranhash.RanHash(16)
	_, err := pq.db.Exec(query, id, userId)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (pq *DbInstance) GetSession(id string) (Session, error) {
	query := `
	SELECT * from sessions where id=$1;
	`

	var session Session
	rows := pq.db.QueryRow(query, id)
	userId := ""
	err := rows.Scan(&session.Id, &userId, &session.CreatedAt)
	if err != nil {
		return Session{}, err
	}

	user, err := pq.GetUser(userId)
	if err != nil {
		return Session{}, err
	}

	session.User = user

	return session, nil
}

func (pq *DbInstance) IsSession(id string) (bool, error) {
	query := `
	SELECT * from sessions where id=$1;
	`

	var session Session
	rows := pq.db.QueryRow(query, id)
	userId := ""
	err := rows.Scan(&session.Id, &userId)
	if err != nil {
		return false, err
	}
	//TODO: Implement Session Expiry Logic
	return true, err
}

type Chat struct {
	Id        string `json:"id,omitempty"`
	Content   string `json:"content,omitempty"`
	UserId    string `json:"user_id,omitempty"`
	RoomId    string `json:"room_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

type CreateChatRequest struct {
	Content string `json:"content,omitempty"`
	UserId  string `json:"user_id,omitempty"`
	RoomId  string `json:"room_id,omitempty"`
}

func (pq *DbInstance) CreateChat(req CreateChatRequest) (string, error) {
	query := `
	INSERT INTO chats(id, content, user_id, room_id)
	VALUES($1, $2, $3, $4);
	`

	id := ranhash.RanHash(8)
	_, err := pq.db.Exec(query, id, req.Content, req.UserId, req.RoomId)
	if err != nil {
		return "", err
	}

	return id, nil
}
