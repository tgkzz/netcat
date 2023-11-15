package internal

import (
	"net"
	"sync"
)

var (
	MaxConnections = 10
	DefaultPort    = "8989"
)

type Message struct {
	from net.Conn
	body string
	info string
}

var (
	MessageHistory []string
	Clients        = make(map[string]net.Conn)
	Connections    = 0
	Mutex          sync.Mutex
)
