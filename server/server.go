package server

import (
	"errors"
	"io"
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

		go func() {
			err := s.handleConnection(conn)
			if err != nil {
				s.config.Logger.Error("connection failed", "error", err)
			}
		}()

	}
}

func (s *TcpServer) handleConnection(connection net.Conn) error {
	defer connection.Close()

	buffer := make([]byte, 1024)
	var received []byte

	for {
		n, err := connection.Read(buffer)
		if n > 0 {
			received = append(received, buffer[:n]...)
		}
		if err == io.EOF {
			sent := 0

			for sent < len(received) {
				n, err := connection.Write(received[sent:])
				if err != nil {
					s.config.Logger.Error("failed to write connection data", "error", err)
					return err
				}
				if n == 0 {
					s.config.Logger.Error("nothing to write", "n", n)
					return errors.New("write have nothing to write n == 0")
				}
				sent += n
			}
			break

		}
		if err != nil {
			s.config.Logger.Error("failed to read connection data", "error", err)
			return err
		}
	}

	return nil
}

func (s *TcpServer) Close() {
	// TODO

}
