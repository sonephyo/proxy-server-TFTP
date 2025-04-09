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

func getImageBytesBlocks(imageBytes []byte, blockSize int) map[int][]byte {
	if blockSize <= 0 {
		panic("blockSize must be greater than 0")
	}

	blocks := make(map[int][]byte)
	for i := 0; i < len(imageBytes); i += blockSize {
		end := i + blockSize
		if end > len(imageBytes) {
			end = len(imageBytes)
		}
		blocks[i/blockSize] = imageBytes[i:end]
	}
	return blocks
}

func operateServerSideImage(conn net.Conn, imgURL string, s *Server) error {

	imageBytes := getImageFromURL(imgURL)

	imageBytesBlocks := getImageBytesBlocks(imageBytes, 512)

	imageLen := len(imageBytes)

	// chunkSize := 512
	// var blockNumber uint16 = 1

	helper.ColorPrintln("yellow", fmt.Sprintf("The image Length: %v", imageLen))

	for key := range imageBytesBlocks {
		if key == 0 || key == 1 || key == 2 {
			fmt.Println(key, imageBytesBlocks[key])
		} else if key == 99 {
			fmt.Println(key, imageBytesBlocks[key])
		} else{
			fmt.Println(key)
		}
	}

	// Preping for mutux
	// var mu sync.Mutex
	// var wg sync.WaitGroup
	// maxInFlight := 10
	// var inFlight int
	// windowNotFull := sync.NewCond(&mu)

	// for i := 0; i < imageLen; i += chunkSize {

	// 	helper.ColorPrintln("white", fmt.Sprintf("Remaining Length: %v", imageLen-i))
	// 	isLastChunk := false
	// 	var dataToSend []byte

	// 	if imageLen-i < 512 {
	// 		// Do something with the remaining which will close the connection
	// 		remainingBytes := imageLen - i
	// 		dataToSend = imageBytes[i : i+remainingBytes]
	// 		isLastChunk = true
	// 	} else {
	// 		dataToSend = imageBytes[i : i+512]
	// 	}

	// 	//Checking and modifying shared mutux if window available to send
	// 	mu.Lock()
	// 	for inFlight >= maxInFlight {
	// 		windowNotFull.Wait()
	// 	}
	// 	inFlight++
	// 	currentBlock := blockNumber
	// 	blockNumber++
	// 	mu.Unlock()

	// 	wg.Add(1)
	// 	go func(block uint16, data []byte, isLast bool) {
	// 		defer wg.Done()
	// 		sendTFTPDATAPacket(conn, s, blockNumber, dataToSend)

	// 		recieveTFTPACKPacket(conn)

	// 		mu.Lock()
	// 		inFlight--
	// 		windowNotFull.Signal()
	// 		mu.Unlock()
	// 	}(currentBlock, dataToSend, isLastChunk)
	// 	time.Sleep(1 * time.Millisecond)
	// 	if isLastChunk {
	// 		break
	// 	}
	// }

	// fmt.Println(imageBytes[:100])
	// wg.Wait()
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
