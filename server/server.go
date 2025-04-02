package main

import (
	// "assignment-2/helper"
	"assignment-2/helper"
	"bytes"
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

var imageCache = make(map[string][]byte)

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

	if imageData, exists := imageCache[url]; exists {
		helper.ColorPrintln("green", "URL found in cache. Returning cached data ...")
		return imageData
	}

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

	imageCache[url] = buf.Bytes()
	return buf.Bytes()
}

func sendTFTPDATAPacket(conn net.Conn, s *Server, blockNumber uint16, selectedBytes []byte) error {
	dataPacket, err := CreateTFTPDATAPacket(blockNumber, selectedBytes)
	if err != nil {
		helper.ColorPrintln("red", "Error occured: "+err.Error())
		return err
	}
	_, err = conn.Write(dataPacket)
	if err != nil {
		fmt.Println(err.Error())
		close(s.quitch)
		return err
	}

	return nil
}

func recieveTFTPACKPacket(conn net.Conn) error {
	buf := make([]byte, 4)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	tftpACKPacket, err := DeserializeTFTPACK(buf[:n])
	if err != nil {
		return err
	}
	fmt.Println("Recieved Ack block number: ", tftpACKPacket.Block)
	return nil
}

func operateServerSideImage(conn net.Conn, imgURL string, s *Server) error {



	imageBytes := getImageFromURL(imgURL)

	imageLen := len(imageBytes)

	chunkSize := 512
	var blockNumber uint16 = 1

	helper.ColorPrintln("yellow", fmt.Sprintf("The image Length: %v", imageLen))

	for i := 0; i < imageLen; i += chunkSize {

		helper.ColorPrintln("white", fmt.Sprintf("Remaining Length: %v", imageLen - i))

		if imageLen-i < 512 {
			// Do something with the remaining which will close the connection
			remainingBytes := imageLen - i
			sendTFTPDATAPacket(conn, s, blockNumber, imageBytes[i:i+remainingBytes])

			recieveTFTPACKPacket(conn)

			time.Sleep(3 * time.Second)
			fmt.Println("3 seconds passes by ...")

			break
		}

		sendTFTPDATAPacket(conn, s, blockNumber, imageBytes[i:i+512]) // Writing

		blockNumber++

		recieveTFTPACKPacket(conn) // Reading

		time.Sleep(1 * time.Millisecond)
		fmt.Println("1 miliseconds passes by ...")
	}

	// imageSize := len(imageBytes)
	// buf := new(bytes.Buffer)
	// err := binary.Write(buf, binary.BigEndian, int32(imageSize))
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(len(imageBytes))
	// fmt.Println(buf.Bytes())
	// conn.Write(buf.Bytes())
	// conn.Write(imageBytes)

	// time.Sleep(1 * time.Minute)
	return nil
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	helper.ColorPrintln("green", "New client connected: "+conn.RemoteAddr().String())

	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("read error: ", err)
		return
	}

	imgURL := string(buf[:n])
	err = operateServerSideImage(conn, imgURL, s)
	if err != nil {
		helper.ColorPrintln("red", "Something went wrong: "+err.Error())
		return
	}

}
