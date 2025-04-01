package main

type tftpRRQPacket struct {
	Opcode   uint16
	Filename string
	Mode     string
}

type tftpDATAPacket struct {
	Opcode uint16
	Block  uint16
	Data   []byte
}

type tftpACKPacket struct {
	Opcode uint16
	Block  uint16
}

type tftpERRORPacket struct {
	Opcode    uint16
	Errorcode uint16
	ErrMsg    string
}
