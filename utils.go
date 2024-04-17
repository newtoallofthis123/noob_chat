package main

import (
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func GetDbUrl() string {
	godotenv.Load()

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		panic("Env variable DB_URL not set")
	}

	return dbUrl
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func MatchPasswords(toCheck string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(toCheck))
	return err == nil
}
