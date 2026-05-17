package server

import (
	"log/slog"
	"net"
	"os"
)

type Config struct {
	Addr   string
	Logger *slog.Logger
}

type TcpServer struct {
	config Config
}

func NewTcpServer(config Config) *TcpServer {
	if config.Addr == "" {
		config.Addr = ":8090"
	}
	if config.Logger == nil {
		config.Logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
	}

	return &TcpServer{
		config: config,
	}
}

func (s *TcpServer) Run() error {
	listener, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		s.config.Logger.Error("failed to listen", "addr", s.config.Addr, "error", err)
		return err
	}

	defer listener.Close()

	s.config.Logger.Info("TCP server running", "address", s.config.Addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go func(connection net.Conn) {
			defer conn.Close()

			data := make([]byte, 1024)

			n, err := connection.Read(data)
			if err != nil {
				s.config.Logger.Error("An error occured while reading the connection data", "error", err)
				return
			}
			clientData := string(data[:n])
			s.config.Logger.Info("client sent", "data", clientData)
			_, err = conn.Write([]byte(clientData))
			if err != nil {
				return
			}
		}(conn)

	}

	return nil
}
