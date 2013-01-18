package server

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
	"netracker/connection_manager"
	"netracker/game"
	"netracker/parser"
	"netracker/util"
)

type Server struct {
	connectionManager *connection_manager.ConnectionManager
	messageParser     *parser.MessageParser
	game              *game.Game
}

func New() *Server {
	game := game.New()
	messageParser := parser.New(game)
	connectionManager := connection_manager.New()

	server := &Server{
		connectionManager: connectionManager,
		messageParser:     messageParser,
		game:              game,
	}

	return server
}

func (server *Server) receiveWebsocketMessage(conn *websocket.Conn) (message string) {
	err := websocket.Message.Receive(conn, &message)
	if err != nil {
		log.Printf("Got error sending WS message: %v", err)
		return
	}

	return
}

func (server *Server) reader(conn *websocket.Conn) {
	for {
		message := server.receiveWebsocketMessage(conn)

		if message == "quit" {
			return
		}

		log.Printf("Got message: %v", message)

		server.messageParser.Parse(message)
		server.connectionManager.Broadcast(server.game.ToJson())
	}

	conn.Close()
}

func (server *Server) handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(util.RelativePath("/../../../public"))))
	mux.Handle("/ws", websocket.Handler(func(conn *websocket.Conn) {
		server.connectionManager.AddConn(conn)
		websocket.Message.Send(conn, server.game.ToJson())
		server.reader(conn)
	}))

	return mux
}

func (server *Server) Run() {
	log.Print("Starting Netracker Server: http://localhost:3000")

	server.connectionManager.Run()
	err := http.ListenAndServe(":3000", server.handler())

	if err != nil {
		log.Printf("error %v", err)
	}
}
