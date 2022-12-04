package chat

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

const (
	CreateTextMessage = "CreateText"
	DeleteMessage     = "Delete"
)

type Messenger struct {
	Mx    sync.Mutex
	Chats map[string]*Chat
}

func (m *Messenger) HasChat(chatID string) bool {
	m.Mx.Lock()
	defer func() {
		m.Mx.Unlock()
	}()

	if _, ok := m.Chats[chatID]; ok {
		return true
	}

	return false
}

func (m *Messenger) AppendChat(chatID string) {
	m.Mx.Lock()
	defer func() {
		m.Mx.Unlock()
	}()

	m.Chats[chatID] = &Chat{
		ChatID:  chatID,
		Mx:      sync.Mutex{},
		Clients: make(map[int64]*Client, 0),
	}
}

func (m *Messenger) AppendClientToChat(chatID string, ws *websocket.Conn, userID int64) {
	m.Mx.Lock()
	defer func() {
		m.Mx.Unlock()
	}()

	m.Chats[chatID].AppendClient(ws, userID)
}

func (m *Messenger) RemoveClientFromChat(chatID string, userID int64) {
	m.Mx.Lock()
	defer func() {
		m.Mx.Unlock()
	}()

	m.Chats[chatID].RemoveClient(userID)
}

func (c *Chat) RemoveClient(userID int64) {
	c.Mx.Lock()
	defer func() {
		c.Mx.Unlock()
	}()

	delete(c.Clients, userID)
}

type Chat struct {
	ChatID  string
	Mx      sync.Mutex
	Clients map[int64]*Client
}

func (c *Chat) AppendClient(ws *websocket.Conn, userID int64) {
	c.Mx.Lock()
	defer func() {
		c.Mx.Unlock()
	}()

	c.Clients[userID] = &Client{
		Socket: ws,
		UserID: userID,
	}
	fmt.Printf("client %d connected\n", userID)
}

type Client struct {
	Socket *websocket.Conn
	UserID int64
}

func NewMessenger() *Messenger {
	return &Messenger{
		Chats: make(map[string]*Chat, 0),
	}
}
