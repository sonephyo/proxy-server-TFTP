package main

import (
	"assignment-2/helper"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
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

func ReadImagePacket(conn net.Conn) ([]byte, error) {
	// 516 bytes per TFTP DATA packet (opcode, block, data)
	chunk := make([]byte, 516)

	// Map to hold the data for each block number.
	// blocks := make(map[uint16][]byte)
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
		// time.Sleep(5 * time.Second)

		blocks = helper.ReplaceInnerSlice(blocks, int(tftpData.Block), tftpData.Data)

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

	sendImageURLTOServer(conn, *imgURL)

	fullMessage, err := ReadImagePacket(conn)
	if err != nil {
		return
	}

	fmt.Println(len(fullMessage))
	saveImageToFile(fullMessage, "test"+".jpg")
	// time.Sleep(10 * time.Second)
	defer conn.Close()
}
