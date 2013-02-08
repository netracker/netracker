package server

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"log"
	"net/http"
	"netracker/pairing"
	"netracker/util"
	"strings"
)

type Server struct {
	pairings []*pairing.Pairing
}

func New() *Server {

	server := &Server{}

	return server
}

type method func(string, *Server, *websocket.Conn) *pairing.Pairing

func pairings(message string, server *Server, conn *websocket.Conn) *pairing.Pairing {
	pairingjson, _ := json.Marshal(server.pairings)
	websocket.Message.Send(conn, string(pairingjson))
	return nil
}

func newGame(message string, server *Server, conn *websocket.Conn) *pairing.Pairing {
	pairId := strings.SplitN(message, " ", 2)[1]
	paired := server.addConnToPairings(pairId, conn)
	paired.Broadcast(paired.Game.ToJson())
	return paired
}

func join(message string, server *Server, conn *websocket.Conn) *pairing.Pairing {
	pairId := strings.SplitN(message, " ", 2)[1]
	paired := server.addConnToPairings(pairId, conn)
	paired.Broadcast(paired.Game.ToJson())
	return paired
}

func parseMessage(message string) method {
	if message == "pairings" {
		return pairings
	} else if strings.Contains(message, "newgame") {
		return newGame
	} else if strings.Contains(message, "join") {
		return join
	}

	return nil
}

func (server *Server) reader(conn *websocket.Conn) {

	var paired *pairing.Pairing = nil
	for {
		var message string
		err := websocket.Message.Receive(conn, &message)

		if err != nil {
			return
		}

		log.Printf("Got message: %v", message)
		action := parseMessage(message)
		if action != nil {
			paired = action(message, server, conn)
		} else {
			paired.Parser.Parse(message)
			paired.Broadcast(paired.Game.ToJson())
		}
	}

	conn.Close()
}

func (server *Server) handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(util.RelativePath("/../../../public"))))
	mux.Handle("/ws", websocket.Handler(func(conn *websocket.Conn) {
		pairingsjson, _ := json.Marshal(server.pairings)
		websocket.Message.Send(conn, string(pairingsjson))
		server.reader(conn)
	}))

	return mux
}

func (server *Server) findOrCreatePairing(pairId string) (newpairing *pairing.Pairing) {
	for _, pairing := range server.pairings {
		if pairing.Id == pairId {
			newpairing = pairing
			return
		}
	}
	newpairing = pairing.New(pairId)
	server.pairings = append(server.pairings, newpairing)
	return
}

func (server *Server) addConnToPairings(pairId string, conn *websocket.Conn) *pairing.Pairing {
	pairing := server.findOrCreatePairing(pairId)
	err := pairing.AddConn(conn)
	if err != nil {
		return nil
	}
	return pairing
}

func (server *Server) Run() {
	log.Print("Starting Netracker Server: http://localhost:3000")

	err := http.ListenAndServe(":3000", server.handler())

	if err != nil {
		log.Printf("error %v", err)
	}
}
