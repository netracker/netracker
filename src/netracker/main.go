package main

import (
	"netracker/game"
	"netracker/parser"
	"netracker/server"
)

func main() {
	game := game.New()
	messageParser := parser.New(game)
	server := server.New(game, messageParser)

	server.Run()
}
