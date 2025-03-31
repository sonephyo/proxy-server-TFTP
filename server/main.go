package main

import (
	"assignment-2/helper"

	"log"
)

func main() {
	server := NewServer(":3000")

	// go func() {
	// 	for msg := range server.Msgch {
	// 		fmt.Printf("recieved message from connection- (%s): %s\n", msg.From, string(msg.Payload))
	// 	}
	// }()
	helper.ColorPrintln("blue", "Starting the server")
	log.Fatal(server.Start())
}
