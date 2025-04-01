package main

import (
	"assignment-2/helper"
	"fmt"
)

func CreateTFTPRRQPacket() ([]byte, error) {

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

func CreateTFTPDATAPacket() ([]byte, error) {
	request := tftpDATAPacket{
		Opcode: 1,
		Block:  1,
		Data:   []byte("Hello"),
	}

	data, err := request.SerializeTFTPDATA()
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}

	helper.ColorPrintln("green", "Serialized TFTP DATA Request: "+string(data))
	fmt.Printf("Hex Dump: % x\n", data)
	return data, nil
}
