package main

import (
	"assignment-2/helper"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"time"
)

func saveImageToFile(imageBytes []byte, filename string) {
	file, err := os.Create("/tmp/" + filename)
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

	for _, row := range arr {
		flattened = append(flattened, row...)
	}

	return flattened
}

func ReadImagePacket(conn net.Conn, key byte, dropPercentage *float64) ([]byte, error) {
	chunk := make([]byte, 516)

	var blocks [][]byte

	for {
		n, err := conn.Read(chunk)
		if err != nil {
			return nil, err
		}

		tftpData, err := DeserializeTFTPDATA(chunk[:n])
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		if rand.Float64() < *dropPercentage {
			helper.ColorPrintln("red", fmt.Sprintf("Simulating drop: NOT sending ACK for block %d", tftpData.Block))
			decryptedData := xorEncryptDecrypt(tftpData.Data, key)
			blocks = helper.ReplaceInnerSlice(blocks, int(tftpData.Block), decryptedData)
			continue
		}
		decryptedData := xorEncryptDecrypt(tftpData.Data, key)

		blocks = helper.ReplaceInnerSlice(blocks, int(tftpData.Block), decryptedData)

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

		if len(tftpData.Data) < 512 {
			break
		}
	}

	helper.ColorPrintln("green", "Finished receiving all data")

	fullMessage := flattenByteArray(blocks)
	return fullMessage, nil
}

func main() {
	hostAddress := "localhost:3000"

	imgURL := flag.String("link", "https://static.boredpanda.com/blog/wp-content/uploads/2020/07/funny-expressive-dog-corgi-genthecorgi-1-1-5f0ea719ea38a__700.jpg", "a string")
	dropPercentage := flag.Float64("drop", 0.0, "droprate (for ignoring some packets)")
	flag.Parse()

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

	fullMessage, err := ReadImagePacket(conn, key, dropPercentage)
	if err != nil {
		return
	}

	saveImageToFile(fullMessage, "soney-"+fmt.Sprint(rand.IntN(1000))+".jpg")
	helper.ColorPrintln("cyan", "Encryption Key: "+fmt.Sprint(key))
	helper.ColorPrintln("yellow", "Bytes Recieved: "+fmt.Sprint(len(fullMessage)))
	defer conn.Close()
}

// Helper Functions
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
