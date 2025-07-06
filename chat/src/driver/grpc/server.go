package grpc

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Server struct {
	UnimplementedChatServiceServer
	mu      *sync.Mutex
	clients map[string]ChatService_ChatServer
	queue   []*Message
}

func New() *Server {
	return &Server{
		mu:      &sync.Mutex{},
		clients: make(map[string]ChatService_ChatServer),
		queue:   []*Message{},
	}
}

func (c *Server) Chat(stream ChatService_ChatServer) error {
	errCh := make(chan error)

	go c.receive(stream, errCh)
	go c.send(errCh)

	return <-errCh
}

func (c *Server) receive(stream ChatService_ChatServer, errCh chan error) {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			continue
		}
		if err != nil {
			errCh <- err
			continue
		}

		if msg.GetMsg() == "LOGGED_IN" || msg.GetMsg() == "LOGGED_OUT" {
			from := msg.GetFrom()
			msg.Id = uuid.New().String()

			var format string
			c.mu.Lock()
			if msg.GetMsg() == "LOGGED_IN" {
				format = "%s has entered."
				c.clients[from] = stream
			} else {
				format = "%s has exited."
				delete(c.clients, from)
			}
			c.mu.Unlock()
			msg.Msg = fmt.Sprintf(format, from)
		} else {
			log.Printf("Received message: {%v}", msg)
		}
		c.mu.Lock()
		c.queue = append(c.queue, msg)
		c.mu.Unlock()
	}
}

func (c *Server) send(errCh chan error) {
	for {
		time.Sleep(100 * time.Millisecond)

		c.mu.Lock()
		msgLen := len(c.queue)
		if msgLen == 0 {
			c.mu.Unlock()
			continue
		}
		c.mu.Unlock()

		for _, msg := range c.queue {
			for _, client := range c.clients {
				if err := client.Send(msg); err != nil {
					log.Printf("Failed to send message: %v", err)
					errCh <- err
					continue
				}
			}
			log.Printf("Sent message: {%v}", msg)
		}

		c.mu.Lock()
		c.queue = []*Message{}
		c.mu.Unlock()
	}
}
