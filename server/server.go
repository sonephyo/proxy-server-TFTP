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

func getImageFromURL(url string) []byte {
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

func operateServerSideImage(conn net.Conn, imgURL string) error {
	imageBytes := getImageFromURL(imgURL)

	imageSize := len(imageBytes)
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(imageSize))
	if err != nil {
		return err
	}

	fmt.Println(len(imageBytes))
	fmt.Println(buf.Bytes())
	conn.Write(buf.Bytes())
	conn.Write(imageBytes)

	time.Sleep(1 * time.Minute)
	return nil
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	helper.ColorPrintln("green", "New client connected: "+conn.RemoteAddr().String())

	// for {

	// send operation for DATA packet
	dataPacket, err := CreateTFTPDATAPacket()
	if err != nil {
		helper.ColorPrintln("red", "Error occured: "+err.Error())
		return
	}
	_, err = conn.Write(dataPacket)
	if err != nil {
		fmt.Println(err.Error())
		close(s.quitch)
		return
	}
	time.Sleep(3 * time.Second)
	fmt.Println("3 seconds passes by ...")


	buf := make([]byte, 4)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	fmt.Println(buf[:n])

	// ack := make([]byte, 4) // TFTP ACK packet size
	// n, err := conn.Read(ack)
	// if err != nil {
	// 	helper.ColorPrintln("red", "Read Error (ACK expected): "+err.Error())
	// 	close(s.quitch)
	// 	return
	// }
	// fmt.Println(ack[:n])

	// buf := make([]byte, 2048)
	// n, err := conn.Read(buf)
	// if err != nil {
	// 	fmt.Println("read error: ", err)
	// 	continue
	// }

	// imgURL := string(buf[:n])
	// err = operateServerSideImage(conn, imgURL)
	// if err != nil {
	// 	helper.ColorPrintln("red", "Something went wrong: "+err.Error())
	// 	return
	// }

	// tftpRRQPacket, err := CreateTFTPRRQPacket()
	// if err != nil {
	// 	helper.ColorPrintln("red", "Error occured: "+err.Error())
	// 	return
	// }
	// conn.Write(tftpRRQPacket)

	// }
}
