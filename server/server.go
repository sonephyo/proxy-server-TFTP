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
	"math/rand/v2"
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

func recieveTFTPACKPacket(conn net.Conn) (uint16, error) {
	buf := make([]byte, 4)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err.Error())
		return 0, err
	}

	tftpACKPacket, err := DeserializeTFTPACK(buf[:n])
	if err != nil {
		return 0, err
	}
	fmt.Println("Recieved Ack block number: ", tftpACKPacket.Block)
	return tftpACKPacket.Block, nil
}

func getImageBytesBlocks(imageBytes []byte, blockSize int) [][]byte {
	if blockSize <= 0 {
		panic("blockSize must be greater than 0")
	}

	var blocks [][]byte
	for i := 0; i < len(imageBytes); i += blockSize {
		end := i + blockSize
		if end > len(imageBytes) {
			end = len(imageBytes)
		}
		blocks = append(blocks, imageBytes[i:end])
	}

	return blocks
}

func operateServerSideImage(conn net.Conn, imgURL string, s *Server, key byte) error {

	imageBytes := getImageFromURL(imgURL)

	imageBytesBlocks := getImageBytesBlocks(imageBytes, 512)

	windowSize := 4
	timeOut := time.Second * 2
	acks := make(map[int]bool)

	for i := 0; i < len(imageBytesBlocks); i += windowSize {
		end := min(i+windowSize, len(imageBytesBlocks))

		for j := i; j < end; j++ {
			encryptedBytes := xorEncryptDecrypt(imageBytesBlocks[j], key)
			sendTFTPDATAPacket(conn, s, uint16(j), encryptedBytes)
		}
		ackCount := 0
		deadline := time.Now().Add(timeOut)

		for ackCount < (end-i) && time.Now().Before(deadline) {
			conn.SetReadDeadline(deadline)
			ack, err := recieveTFTPACKPacket(conn)
			if err == nil {
				if !acks[int(ack)] {
					acks[int(ack)] = true
					ackCount++
				}
			}
		}

		if ackCount < (end - i) {
			i -= windowSize
		}
	}

	return nil
}

func (s *Server) readLoop(conn net.Conn) {

	defer conn.Close()
	helper.ColorPrintln("green", "New client connected: "+conn.RemoteAddr().String())

	clientID := make([]byte, 2048)
	_, err := conn.Read(clientID)
	if err != nil {
		log.Fatal("Error: reading clientID", err)
		return
	}
	clientIDByte := clientID[0]

	sessionNum := rand.IntN(100)
	lengthBuffer := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuffer, uint32(sessionNum))
	conn.Write(lengthBuffer)

	key := generateKey(clientIDByte, byte(sessionNum))

	time.Sleep(1 * time.Millisecond)

	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("read error: ", err)
		return
	}

	imgURL := string(buf[:n])
	err = operateServerSideImage(conn, imgURL, s, key)
	if err != nil {
		helper.ColorPrintln("red", "Something went wrong: "+err.Error())
		return
	}
	fmt.Println("Key is ", key)
}

// Helper functions
func generateKey(clientID byte, sessionNum byte) byte {
	return clientID ^ sessionNum
}

func xorEncryptDecrypt(data []byte, key byte) []byte {
	enc := make([]byte, len(data))
	for i := range data {
		enc[i] = data[i] ^ key
	}
	return enc
}
