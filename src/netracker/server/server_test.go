package server

import (
	"encoding/json"
	"github.com/bmizerany/assert"
	"github.com/drewolson/testflight"
	"github.com/drewolson/testflight/ws"
	"netracker/game"
	"netracker/player"
	"testing"
)

func withServer(f func(*testflight.Requester)) {
	netrackerServer := New()
	testflight.WithServer(netrackerServer.handler(), f)
}

func TestRoot(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		response := r.Get("/")

		assert.Equal(t, 200, response.StatusCode)
	})
}

func TestWebsocketInitialGameState(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		connection := ws.Connect(r, "/ws")
		defer connection.Close()

		message, _ := connection.ReceiveMessage()

		game := game.Game{}
		json.Unmarshal([]byte(message), &game)

		assert.Equal(t, 5, game.CorpCredits)
		assert.Equal(t, 5, game.RunnerCredits)
		assert.Equal(t, player.CORP, game.ActivePlayer.Role)
		assert.Equal(t, 0, game.Clicks)
	})
}

func TestWebSocketAcceptsMessages(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		connection := ws.Connect(r, "/ws")
		defer connection.Close()
		connection.ReceiveMessage()
		connection.SendMessage("addcorpcredit")
		message, _ := connection.ReceiveMessage()
		connection.Close()

		game := game.Game{}
		json.Unmarshal([]byte(message), &game)

		assert.Equal(t, 6, game.CorpCredits)

	})
}

func TestPairsOfConnectionsAreIsolated(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		client1 := ws.Connect(r, "/ws")
		defer client1.Close()

		client2 := ws.Connect(r, "/ws")
		defer client2.Close()

		client3 := ws.Connect(r, "/ws")
		defer client3.Close()

		client1.SendMessage("addcorpcredit")
		client2.FlushMessages(2)

		client3.SendMessage("removecorpcredit")
		client3.FlushMessages(2)

		game1 := game.Game{}
		json.Unmarshal([]byte(client2.ReceivedMessages[1]), &game1)

		game2 := game.Game{}
		json.Unmarshal([]byte(client3.ReceivedMessages[1]), &game2)

		assert.Equal(t, 6, game1.CorpCredits)
		assert.Equal(t, 4, game2.CorpCredits)
	})
}
