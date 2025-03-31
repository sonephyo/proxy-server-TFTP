package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func DeserializeTFTPRRQ(data []byte) (*tftpRRQPacket, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("invalid TFTP RRQ packet: too short")
	}

	opCode := binary.BigEndian.Uint16(data[:2])
	data = data[2:]

	filenameEnd := bytes.IndexByte(data, 0)
	filename := string(data[:filenameEnd])
	data = data[filenameEnd+1:]

	modeEnd := bytes.IndexByte(data, 0)
	mode := string(data[:modeEnd])

	return &tftpRRQPacket{
		Opcode:   opCode,
		Filename: filename,
		Mode:     mode,
	}, nil
}

func DeserializeTFTPDATA(data []byte) (*tftpDATAPacket, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("invalid TFTP DATA packet: too short")
	}

	opCode := binary.BigEndian.Uint16(data[:2])
	data = data[2:]

	block := binary.BigEndian.Uint16(data[:2])
	data = data[2:]

	packetData := data[:]

	return &tftpDATAPacket{
		Opcode: opCode,
		Block: block,
		Data: []byte(packetData),
	}, nil
}