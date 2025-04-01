package main

import (
	"assignment-2/helper"
)

func CreateTFTPACKPacket(blockNumber uint16) ([]byte, error) {
	tftpAckPacket := tftpACKPacket{
		Opcode: 4,
		Block:  blockNumber,
	}

	tftpAckBytes, err := tftpAckPacket.SerializeTFTPACK()
	if err != nil {
		helper.ColorPrintln("red", "Error serializing ACK packet: "+err.Error())
		return nil, err
	}

	return tftpAckBytes, nil
}
