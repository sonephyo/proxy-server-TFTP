package main

import (
	"assignment-2/helper"
	"fmt"
	"net"
	"time"
)

type Message struct {
	From    string
	Payload []byte
}

// The server is for creating a new server that include the port that is going to listen ln I believe is the listener and quitch which is
type Server struct {
	listenAddress string
	ln            net.Listener
	quitch        chan struct{}
	Msgch         chan Message
}

// The NetServer represent the method in which creating a server and giving a default values
func NewServer(listenAddress string) *Server {
	return &Server{
		listenAddress: listenAddress,
		quitch:        make(chan struct{}),
		Msgch:         make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitch
	close(s.Msgch)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}

		fmt.Println("new connection to the server: ", conn.RemoteAddr())

		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	// buf := make([]byte, 2048)
	for {
		// n, err := conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("read error: ", err)
		// 	continue
		// }
		s.Msgch <- Message{
			From:    conn.RemoteAddr().String(),
			Payload: []byte("This is dummy text"),
		}

		tftpRRQPacket, err := CreateRRQPacket()
		if err != nil {
			helper.ColorPrintln("red", "Error occured: "+err.Error())
			return
		}
		conn.Write(tftpRRQPacket)

		time.Sleep(2 * time.Second)
	}
}
