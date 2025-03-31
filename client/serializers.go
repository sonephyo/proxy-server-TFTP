package main

import (
	"bytes"
	"encoding/binary"
)

// func (req *tftpDATAPacket) SerializeTFTPDATA() ([]byte, error) {
// 	buf := new(bytes.Buffer)

// 	if err := binary.Write(buf, binary.BigEndian, req.Opcode); err != nil {
// 		return nil, err
// 	}
// 	if err := binary.Write(buf, binary.BigEndian, req.Block); err != nil {
// 		return nil, err
// 	}

// 	if _, err := buf.Write(req.Data); err != nil {
// 		return nil, err
// 	}

// 	return buf.Bytes(), nil
// }

func (req *tftpACKPacket) SerializeTFTPACK() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, req.Opcode); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, req.Block); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (req *tftpERRORPacket) SerializeTFTPERROR() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, req.Opcode); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, req.Errorcode); err != nil {
		return nil, err
	}

	buf.WriteString(req.ErrMsg)
	buf.WriteByte(0)

	return buf.Bytes(), nil
}


