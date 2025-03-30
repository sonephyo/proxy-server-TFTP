package main

import (
	// "assignment-2/helper"
	"assignment-2/helper"
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net"
	"net/http"
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

func getImageFromURL() []byte {
	url := "https://static.boredpanda.com/blog/wp-content/uploads/2020/07/funny-expressive-dog-corgi-genthecorgi-1-1-5f0ea719ea38a__700.jpg"
	response, e := http.Get(url)
	fmt.Println("Response Status code: ", response.StatusCode)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()
	selected_image, _, _ := image.Decode(response.Body)

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, selected_image, nil)
	if err != nil {
		fmt.Println("Something wrong with getting image from URL: ", err.Error())
	}
	return buf.Bytes()
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
		// s.Msgch <- Message{
		// 	From:    conn.RemoteAddr().String(),
		// 	Payload: []byte("This is dummy text"),
		// }

		// tftpRRQPacket, err := CreateRRQPacket()
		// if err != nil {
		// 	helper.ColorPrintln("red", "Error occured: "+err.Error())
		// 	return
		// }
		// conn.Write(tftpRRQPacket)

		// buf := new(bytes.Buffer)
		// err := jpeg.Encode(buf, new_image, nil)
		// send_s3 := buf.Bytes()

		imageBytes := getImageFromURL()

		imageSize := len(imageBytes)
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, int32(imageSize))
		if err != nil {
			helper.ColorPrintln("red", "Something went wrong: " + err.Error())
			return
		}

		fmt.Println(len(imageBytes))
		fmt.Println(buf.Bytes())
		conn.Write(buf.Bytes())
		conn.Write(imageBytes)

		time.Sleep(1 * time.Minute)
	}
}
