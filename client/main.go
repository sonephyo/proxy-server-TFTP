package main

import (
	"assignment-2/helper"
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
	for {
		n, err:= conn.Read(chunk)
		if err != nil {
			helper.ColorPrintln("red", "Error occured while sending from server to client")
			return
		}
		tftp, err := DeserializeTFTPRRQ(chunk[:n])
		if err != nil {
			fmt.Println("Error: ", err)
			return	
		}
		fmt.Println(string(chunk[:n]))
		fmt.Println(tftp.Filename, ", ", tftp.Mode, ", ", tftp.Opcode)
	}
}