package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello World")

	file, err := os.OpenFile("log_ignore_.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	s := NewServerWithLog(file)
	s.StartServer("localhost:1234")
}
