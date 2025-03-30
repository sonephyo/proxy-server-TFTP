package main

import (
	"fmt"
	"net"
)

func main() {
	hostAddress := "localhost:3000"

	// Creating connection with the server
	conn, err := net.Dial("tcp", hostAddress)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	chunk := make([]byte, 1024)
	n, err:= conn.Read(chunk)
}