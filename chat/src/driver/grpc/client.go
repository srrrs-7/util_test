package grpc

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	ChatServiceClient
	From string
}

func (c *Client) NewClient() {
	stream, err := c.Chat(context.Background())
	if err != nil {
		log.Printf("Failed to chat: %v", err)
		return
	}
	defer stream.CloseSend()

	waitCh := make(chan struct{})

	go c.receive(stream, waitCh)
	go c.send(stream, waitCh)
	c.sendMsg(stream, "LOGGED_IN")

	<-waitCh
}

func (c *Client) receive(stream ChatService_ChatClient, waitCh chan struct{}) {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Failed to receive message: %v", err)
			continue
		}
		log.Printf("[%s] : %s (%s)\n", msg.GetFrom(), msg.GetMsg(), msg.GetCreatedAt())
	}
	waitCh <- struct{}{}
}

func (c *Client) send(stream ChatService_ChatClient, waitCh chan struct{}) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Failed to read: %v", err)
			continue
		}
		trimmed := strings.TrimSpace(text)
		if len(trimmed) == 0 {
			continue
		}
		if trimmed == "exit" {
			c.sendMsg(stream, "LOGGED_OUT")
			break
		}
		c.sendMsg(stream, trimmed)
	}
	waitCh <- struct{}{}
}

func (c *Client) sendMsg(stream ChatService_ChatClient, body string) {
	msg := &Message{
		From:      c.From,
		Id:        uuid.New().String(),
		Msg:       body,
		CreatedAt: timestamppb.Now(),
	}

	if err := stream.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
