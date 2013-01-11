package server

import (
        "netracker/connection_manager"
        "code.google.com/p/go.net/websocket"
        "net/http"
        "log"
)

var globalManager = connection_manager.New()

func receiveWebsocketMessage(conn *websocket.Conn) (message string) {
        err := websocket.Message.Receive(conn, &message)
        if err != nil {
                log.Printf("Got error sending WS message: %v", err)
                return
        }

        return
}

func reader(conn *websocket.Conn) {
        for {
                message := receiveWebsocketMessage(conn)
								log.Printf("Got message: %v", message)

        }

        conn.Close()
}

func wsHandler(conn *websocket.Conn) {
        globalManager.AddConn(conn)
        reader(conn)
}


func Run() {
        log.Print("Starting Netracker Server: http://localhost:3000")
        go globalManager.Run()

        http.Handle("/", http.FileServer(http.Dir("public")))
        http.Handle("/ws", websocket.Handler(wsHandler))

        err := http.ListenAndServe(":3000", nil)
        if err != nil {
                log.Printf("error %v", err)
        }
}
