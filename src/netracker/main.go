package main

import (
	"netracker/server"
)

func main() {
	server := server.New()

	server.Run()
}
