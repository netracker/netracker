package server

import (
	"encoding/json"
	"fmt"
	"github.com/bmizerany/assert"
	"github.com/drewolson/testflight"
	"github.com/drewolson/testflight/ws"
	"netracker/game"
	"netracker/pairing"
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

func TestWebsocketListsCurrentPairingsOnConnect(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		c1 := ws.Connect(r, "/ws")
		defer c1.Close()
		c1.SendMessage("newgame bar")
		c2 := ws.Connect(r, "/ws")
		defer c2.Close()
		c2.SendMessage("newgame baz")
		c3 := ws.Connect(r, "/ws")
		defer c3.Close()
		pairingmessage, _ := c3.ReceiveMessage()
		var pairings []*pairing.Pairing
		json.Unmarshal([]byte(pairingmessage), &pairings)

		assert.Equal(t, 2, len(pairings))
		assert.Equal(t, "bar", pairings[0].Id)
	})
}

func TestWebsocketCanJoinPairingInProgress(t *testing.T) {

	withServer(func(r *testflight.Requester) {
		client1 := ws.Connect(r, "/ws")
		defer client1.Close()

		client2 := ws.Connect(r, "/ws")
		defer client2.Close()

		client1.SendMessage("newgame bar")
		client2.SendMessage("pairings")

		client1.FlushMessages(2)
		client2.FlushMessages(2)

		pairingmessages := client2.ReceivedMessages[1]
		fmt.Println(pairingmessages)

		client2.SendMessage("join bar")

		client1.SendMessage("addcorpcredit")
		client1.FlushMessages(2)
		gamemessage := client1.ReceivedMessages[3]

		var game *game.Game
		json.Unmarshal([]byte(gamemessage), &game)
		assert.Equal(t, game.CorpCredits, 6)

	})
}

func TestWebsocketInitialPairingState(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		connection := ws.Connect(r, "/ws")
		defer connection.Close()

		connection.SendMessage("newgame myid")

		connection.FlushMessages(2)
		message := connection.ReceivedMessages[1]

		game := game.Game{}
		json.Unmarshal([]byte(message), &game)

		assert.Equal(t, 5, game.CorpCredits)
		assert.Equal(t, 5, game.RunnerCredits)
		assert.Equal(t, player.CORP, game.ActivePlayer.Role)
		assert.Equal(t, 0, game.Clicks)

		connection.SendMessage("pairings")
		pairingmessage, _ := connection.ReceiveMessage()

		var pairings []*pairing.Pairing
		json.Unmarshal([]byte(pairingmessage), &pairings)
		fmt.Println(pairings)
		assert.Equal(t, "myid", pairings[0].Id)
	})
}

func TestWebSocketAcceptsMessages(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		connection := ws.Connect(r, "/ws")
		defer connection.Close()
		connection.SendMessage("newgame foo")
		connection.SendMessage("addcorpcredit")

		err := connection.FlushMessages(3)
		if err != nil {
			assert.Equal(t, true, false)
		}
		message := connection.ReceivedMessages[2]

		game := game.Game{}
		json.Unmarshal([]byte(message), &game)

		assert.Equal(t, 6, game.CorpCredits)

	})
}

func TestPairsOfConnectionsAreIsolated(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		client1 := ws.Connect(r, "/ws")
		defer client1.Close()
		client1.SendMessage("newgame bar")

		client2 := ws.Connect(r, "/ws")
		defer client2.Close()
		client2.SendMessage("join bar")

		client3 := ws.Connect(r, "/ws")
		defer client3.Close()
		client3.SendMessage("newgame baz")

		client1.SendMessage("addcorpcredit")
		client2.FlushMessages(3)

		client3.SendMessage("removecorpcredit")
		client3.FlushMessages(3)

		game1 := game.Game{}
		json.Unmarshal([]byte(client2.ReceivedMessages[2]), &game1)

		game2 := game.Game{}
		json.Unmarshal([]byte(client3.ReceivedMessages[2]), &game2)

		assert.Equal(t, 6, game1.CorpCredits)
		assert.Equal(t, 4, game2.CorpCredits)
	})
}
