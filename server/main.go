package main

import (
	"assignment-2/helper"

	"log"
)

func main() {
	server := NewServer(":3000")

	helper.ColorPrintln("blue", "Starting the server")
	log.Fatal(server.Start())
}
