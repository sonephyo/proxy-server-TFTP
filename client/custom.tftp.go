package main

import (
	"assignment-2/helper"
	"fmt"
)

func CreateTFTPACKPacket() ([]byte, error) {
	tftpAckPacket := tftpACKPacket{
		Opcode: 4,
		Block:  1,
	}

	tftpAckBytes, err := tftpAckPacket.SerializeTFTPACK()
	if err != nil {
		helper.ColorPrintln("red", "Error serializing ACK packet: "+err.Error())
		return nil, err
	}

	helper.ColorPrintln("green", fmt.Sprintf("Serialized TFTP ACK Packet: % X", tftpAckBytes))
	return tftpAckBytes, nil
}
