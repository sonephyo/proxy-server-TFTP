package main

import (
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

	return data, nil
}

func CreateTFTPDATAPacket(blockNumber uint16, selectedBytes []byte) ([]byte, error) {
	request := tftpDATAPacket{
		Opcode: 1,
		Block:  blockNumber,
		Data:   selectedBytes,
	}

	data, err := request.SerializeTFTPDATA()
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}

	return data, nil
}
