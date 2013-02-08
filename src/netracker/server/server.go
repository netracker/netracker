package server

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"log"
	"net/http"
	"netracker/pairing"
	"netracker/util"
)

type Server struct {
	pairings []*pairing.Pairing
}

func New() *Server {

	server := &Server{}

	return server
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

		if message == "pairings" {
			pairingjson, _ := json.Marshal(server.pairings)
			websocket.Message.Send(conn, string(pairingjson))
		} else if message == "newgame" && paired == nil {
			paired = server.addConnToPairings(conn)
			paired.Broadcast(paired.Game.ToJson())
		} else if message == "join id" && paired == nil {
			paired = server.addConnToPairings(conn)
			paired.Broadcast(paired.Game.ToJson())
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

func (server *Server) addConnToPairings(conn *websocket.Conn) *pairing.Pairing {
	if len(server.pairings) < 1 {
		server.pairings = append(server.pairings, pairing.New("id"))
	}
	lastpairing := server.pairings[len(server.pairings)-1]
	err := lastpairing.AddConn(conn)
	if err != nil {
		newpairing := pairing.New("id")
		newpairing.AddConn(conn)
		server.pairings = append(server.pairings, newpairing)
		return newpairing
	}
	return lastpairing
}

func (server *Server) Run() {
	log.Print("Starting Netracker Server: http://localhost:3000")

	err := http.ListenAndServe(":3000", server.handler())

	if err != nil {
		log.Printf("error %v", err)
	}
}
