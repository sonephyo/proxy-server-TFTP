package main

import (
	"assignment-2/helper"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func saveImageToFile(imageBytes []byte, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Error creating file: ", err)
	}
	defer file.Close()

	// Write the image bytes to the file
	_, err = file.Write(imageBytes)
	if err != nil {
		log.Fatal("Error writing image to file: ", err)
	}

	fmt.Println("Image saved successfully to", filename)
}

func sendImageURLTOServer(conn net.Conn, imgURL string) error {

	_, err := conn.Write([]byte(imgURL))
	if err != nil {
		log.Fatal("Error: sending image to server: ", err)
		return err
	}
	return nil
}

func flattenByteArray(arr [][]byte) []byte {
	totalLength := 0
	for _, row := range arr {
		totalLength += len(row)
	}

	flattened := make([]byte, 0, totalLength)

	for i, row := range arr {
		fmt.Printf("--> This is row number %v data is %v\n", i, row[:5])
		flattened = append(flattened, row...)
	}
	fmt.Printf("-----------> Data temp to check: %v", flattened[:10])

	return flattened
}

func ReadImagePacket(conn net.Conn, key byte) ([]byte, error) {
	chunk := make([]byte, 516)

	var blocks [][]byte

	for {
		n, err := conn.Read(chunk)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		// Deserialize the received TFTP DATA packet.
		tftpData, err := DeserializeTFTPDATA(chunk[:n])
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		decryptedData := xorEncryptDecrypt(tftpData.Data, key)

		blocks = helper.ReplaceInnerSlice(blocks, int(tftpData.Block), decryptedData)

		// Create and send ACK for the current block.
		tftpAckPacketBytes, err := CreateTFTPACKPacket(tftpData.Block)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		_, err = conn.Write(tftpAckPacketBytes)
		if err != nil {
			log.Fatal("Error sending ACK:", err)
			return nil, err
		}
		fmt.Printf("Sent ACK for block %d\n", tftpData.Block)
		helper.ColorPrintln("cyan", "The value of n is "+fmt.Sprint(n))

		// Per TFTP, the final data packet is less than the full packet size.
		if n != 516 {
			break
		}
	}

	helper.ColorPrintln("green", "Finished receiving all data")

	fullMessage := flattenByteArray(blocks)

	helper.ColorPrintln("green", fmt.Sprintf("Full message length: %v", len(fullMessage)))
	// fmt.Println(fullMessage[:])
	return fullMessage, nil
}

func main() {
	hostAddress := "localhost:3000"

	imgURL := flag.String("link", "https://static.boredpanda.com/blog/wp-content/uploads/2020/07/funny-expressive-dog-corgi-genthecorgi-1-1-5f0ea719ea38a__700.jpg", "a string")
	flag.Parse()

	// Creating connection with the server
	conn, err := net.Dial("tcp", hostAddress)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	clientID := *imgURL
	_, err = conn.Write([]byte(clientID))
	if err != nil {
		fmt.Println("Write error:", err)
		return
	}
	time.Sleep(1 * time.Millisecond)

	sessionNumBytes := make([]byte, 4)
	_, err = io.ReadFull(conn, sessionNumBytes)
	if err != nil {
		if err == io.EOF {
			fmt.Println("Client closed Connection")
			return
		}
		fmt.Println("Read full Error: ", err)
		return
	}
	sessionNum := binary.BigEndian.Uint32(sessionNumBytes)

	if sessionNum == 0 {
		fmt.Println("Invalid chuckSize")
		return
	}

	key := generateKey([]byte(clientID)[0], byte(sessionNum))
	

	sendImageURLTOServer(conn, *imgURL)

	fullMessage, err := ReadImagePacket(conn, key)
	if err != nil {
		return
	}

	fmt.Println(len(fullMessage))
	saveImageToFile(fullMessage, "test"+".jpg")
	fmt.Println("Key is ", key)
	defer conn.Close()
}

// Helper Functions
func generateKey(clientID byte, sessionNum byte) byte{
	return clientID ^ sessionNum
}

func xorEncryptDecrypt(data []byte, key byte) []byte {
	enc := make([]byte, len(data))
	for i := range data {
		enc[i] = data[i] ^ key
	}
	return enc
}
