package chat

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Messenger struct {
	Chat Chat
}

type Chat struct {
	// ChatID  int64
	Mx      sync.Mutex
	Clients []*Client
}

type Client struct {
	Socket *websocket.Conn
	UserID int64
}

func NewMessenger() *Messenger {
	return &Messenger{
		Chat: Chat{
			Mx: sync.Mutex{},
		},
	}
}
