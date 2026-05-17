package main

import (
	"log/slog"
	"os"

	"github.com/Berveglieri/tcp-echo/server"
)

func main() {
	s := server.NewTcpServer(server.Config{
		Addr: ":8090",
	})
	err := s.Run()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
