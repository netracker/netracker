package server

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
	"netracker/connection_manager"
	"netracker/game"
	"netracker/parser"
)

type Server struct {
	connectionManager *connection_manager.ConnectionManager
	messageParser     *parser.MessageParser
	websocketHandler  func(conn *websocket.Conn)
	game              *game.Game
}

func New(game *game.Game, messageParser *parser.MessageParser) *Server {
	server := &Server{
		connectionManager: connection_manager.New(),
		messageParser:     messageParser,
		game:              game,
	}

	server.websocketHandler = func(conn *websocket.Conn) {
		server.connectionManager.AddConn(conn)
		websocket.Message.Send(conn, server.game.ToJson())
		server.reader(conn)
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
		log.Printf("Got message: %v", message)

		server.messageParser.Parse(message)
		server.connectionManager.Broadcast(server.game.ToJson())
	}

	conn.Close()
}

func (server *Server) Run() {
	log.Print("Starting Netracker Server: http://localhost:3000")
	go server.connectionManager.Run()

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/ws", websocket.Handler(server.websocketHandler))

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Printf("error %v", err)
	}
}
