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
	netrackerServer.connectionManager.Run()
	defer netrackerServer.connectionManager.Stop()
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

		message := connection.ReceiveMessage()
		connection.WriteMessage("quit")

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
		_ = connection.ReceiveMessage()
		connection.WriteMessage("addcorpcredit")
		message := connection.ReceiveMessage()
		connection.WriteMessage("quit")

		game := game.Game{}
		json.Unmarshal([]byte(message), &game)

		assert.Equal(t, 6, game.CorpCredits)

	})
}
