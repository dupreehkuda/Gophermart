package main

import (
	"github.com/dupreehkuda/Gophermart/internal/server"
)

func main() {
	srv := server.NewByConfig()
	srv.Run()
}
