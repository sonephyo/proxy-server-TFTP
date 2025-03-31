package main

import (
	"assignment-2/helper"
	"fmt"
)

func CreateTFTPACKPacket() ([]byte, error) {
	tftpAckPacket := tftpACKPacket{
		Opcode: 1,
		Block: 1,
	}

	tftpAckBytes, err := tftpAckPacket.SerializeTFTPACK()
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}

	helper.ColorPrintln("green", "Serialized TFTP DATA Request: "+string(tftpAckBytes))
	return tftpAckBytes, nil	
}