package connection_manager

import (
	"code.google.com/p/go.net/websocket"
)

type ConnectionManager struct {
	connections []*websocket.Conn
	register    chan *websocket.Conn
}

func New() *ConnectionManager {
	return &ConnectionManager{
		connections: []*websocket.Conn{},
		register:    make(chan *websocket.Conn),
	}
}

func (manager *ConnectionManager) AddConn(conn *websocket.Conn) {
	manager.register <- conn
}

func (manager *ConnectionManager) Broadcast(message string) {
	for _, conn := range manager.connections {
		websocket.Message.Send(conn, message)
	}
}

func (manager *ConnectionManager) Run() {
	for {
		conn := <-manager.register
		manager.connections = append(manager.connections, conn)
	}
}
