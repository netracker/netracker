package pairing

import (
	"code.google.com/p/go.net/websocket"
	"github.com/bmizerany/assert"
	"testing"
)

func TestAddConn(t *testing.T){
	conn := &websocket.Conn{}
	pairing := New()
	pairing.AddConn(conn)

	assert.Equal(t, conn, pairing.connections[0])
}

func TestAddThirdConnReturnsError(t *testing.T){
	conn1 := &websocket.Conn{}
	conn2 := &websocket.Conn{}
	conn3 := &websocket.Conn{}
	pairing := New()

	pairing.AddConn(conn1)
	pairing.AddConn(conn2)

	err := pairing.AddConn(conn3)
	assert.Equal(t, &FullPairingError{}, err)
}
