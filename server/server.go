package main

import (
	// "assignment-2/helper"
	"assignment-2/helper"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net"
	"net/http"
	"slices"
	"strings"
	"sync"
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

	response, _ := http.Get(url)
	fmt.Println("Response Status code: ", response.StatusCode)
	// if e != nil {
	// 	log.Fatal(e)
	// }
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
	helper.ColorPrintln("magenta", "Sending: "+fmt.Sprint(blockNumber))
	_, err = conn.Write(dataPacket)
	if err != nil {
		fmt.Println(err.Error())
		close(s.quitch)
		return err
	}

	return nil
}

func recieveTFTPACKPacket(conn net.Conn) (int, error) {
	buf := make([]byte, 4)
	n, err := conn.Read(buf)
	if err != nil {
		// log.Fatal(err.Error())
		return -1, err
	}

	tftpACKPacket, err := DeserializeTFTPACK(buf[:n])
	if err != nil {
		return -1, err
	}
	fmt.Println("Recieved Ack block number: ", tftpACKPacket.Block)
	return int(tftpACKPacket.Block), nil
}

func operateServerSideImage(conn net.Conn, imgURL string, s *Server) error {

	imageBytes := getImageFromURL(imgURL)

	imageLen := len(imageBytes)

	chunkSize := 512

	imageBytesBlocks := make(map[uint16][]byte)
	currentIdx := 0
	var i uint16 = 0
	for {
		if currentIdx+chunkSize >= imageLen {
			imageBytesBlocks[i] = imageBytes[currentIdx:]
			break
		}
		imageBytesBlocks[i] = imageBytes[currentIdx : currentIdx+chunkSize]
		currentIdx += chunkSize
		i += 1
	}

	for item, value := range imageBytesBlocks {
		fmt.Printf("%v Key: %v\n", item, len(value))
	}

	var blockNumber uint16 = 0

	helper.ColorPrintln("yellow", fmt.Sprintf("The image Length: %v", imageLen))

	// Preping for mutux
	var mu sync.Mutex
	var wg sync.WaitGroup
	maxInFlight := 10
	var inFlight int
	windowNotFull := sync.NewCond(&mu)

	var acknowledgedBlocks []int

	i = 0
	for key := range len(imageBytesBlocks) - 1 {
		dataToSend := imageBytesBlocks[uint16(key)]
		helper.ColorPrintln("red", "Key: "+fmt.Sprint(key))

		mu.Lock()
		for inFlight >= maxInFlight {
			windowNotFull.Wait()
		}
		inFlight++
		currentBlock := blockNumber
		blockNumber++
		mu.Unlock()

		wg.Add(1)
		go func(block uint16, data []byte) {
			defer wg.Done()
			sendTFTPDATAPacket(conn, s, blockNumber, dataToSend)

			recievedBlock, _ := recieveTFTPACKPacket(conn)

			mu.Lock()
			inFlight--
			windowNotFull.Signal()
			acknowledgedBlocks = append(acknowledgedBlocks, recievedBlock)
			mu.Unlock()
		}(currentBlock, dataToSend)

		time.Sleep(1 * time.Millisecond)

	}
	keys := make([]uint16, 0, len(imageBytesBlocks))
	var tempL int
	for k, v := range imageBytesBlocks {
		keys = append(keys, k)
		tempL += len(v)
	}
	helper.ColorPrintln("blue", "Here:"+fmt.Sprint(tempL))
	slices.Sort(keys)
	fmt.Println(keys)

	operateUnAckBlocks(conn, s, imageBytesBlocks, acknowledgedBlocks)

	fmt.Println(imageBytes[:200])
	wg.Wait()
	return nil
}

func operateUnAckBlocks(conn net.Conn, s *Server, imageBytesBlocks map[uint16][]byte, acknowledgedBlocks []int) {
	var unAckBlocks []int
	for key := range imageBytesBlocks {
		found := false
		for _, ackBlock := range acknowledgedBlocks {
			if int(key) == ackBlock {
				found = true
				break
			}
		}
		if !found {
			unAckBlocks = append(unAckBlocks, int(key))
		}
	}

	helper.ColorPrintln("red", "Unack blocks: "+strings.Trim(strings.Join(strings.Fields(fmt.Sprint(unAckBlocks)), ","), "[]"))

	sendUnAcknowledgedBlocks(conn, s, imageBytesBlocks, unAckBlocks)
}

func sendUnAcknowledgedBlocks(conn net.Conn, s *Server, imageBytesBlocks map[uint16][]byte, blocks []int) error {

	for _, blockNumber := range blocks {
		helper.ColorPrintln("cyan", fmt.Sprint(blockNumber)+", "+fmt.Sprint(len(imageBytesBlocks[uint16(blockNumber)])))

		dataToSend := imageBytesBlocks[uint16(blockNumber)]
		sendTFTPDATAPacket(conn, s, uint16(blockNumber), dataToSend)

		recieveTFTPACKPacket(conn)
		time.Sleep(1 * time.Millisecond)
	}
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
