package main

import (
	"assignment-2/helper"
	"fmt"
)

type tftpRRQPacket struct {
	Opcode   uint16
	Filename string
	Mode     string
}

func CreateRRQPacket() ([]byte, error) {

	request := tftpRRQPacket{
		Opcode:   1,
		Filename: "test.txt",
		Mode:     "octet",
	}

	data, err := request.SerializeTFTPRRQ()
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}

	helper.ColorPrintln("green", "Serialized TFTP RRQ Request: "+string(data))
	fmt.Printf("Hex Dump: % x\n", data)
	return data, nil
}
