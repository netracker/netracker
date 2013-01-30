package pairing

import (
	"code.google.com/p/go.net/websocket"
	"netracker/game"
	"netracker/parser"
)

type Pairing struct {
	connections []*websocket.Conn
	Game        *game.Game
	Parser      *parser.MessageParser
}

type FullPairingError struct {
}

func (error *FullPairingError) Error() string {
	return "Pairing is full"
}

func New() *Pairing {
	game := game.New()
	return &Pairing{
		Game:   game,
		Parser: parser.New(game),
	}
}

func (pairing *Pairing) AddConn(conn *websocket.Conn) *FullPairingError {
	if len(pairing.connections) >= 2 {
		return &FullPairingError{}
	}
	pairing.connections = append(pairing.connections, conn)
	return nil
}

func (pairing *Pairing) Broadcast(message string) {
	for _, conn := range pairing.connections {
		websocket.Message.Send(conn, message)
	}
}
