package main

import (
	"bytes"
	"encoding/binary"
)

func (req *tftpRRQPacket) SerializeTFTPRRQ() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, req.Opcode); err != nil {
		return nil, err
	}

	buf.WriteString(req.Filename)
	buf.WriteByte(0)
	buf.WriteString(req.Mode)
	buf.WriteByte(0)
	return buf.Bytes(), nil
}

func (req *tftpDATAPacket) SerializeTFTPDATA() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, req.Opcode); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, req.Block); err != nil {
		return nil, err
	}

	if _, err := buf.Write(req.Data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
