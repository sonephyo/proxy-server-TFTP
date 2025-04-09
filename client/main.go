package main

import (
	"assignment-2/helper"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
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
func ReadImagePacket(conn net.Conn) ([]byte, error) {
	// 516 bytes per TFTP DATA packet (opcode, block, data)
	chunk := make([]byte, 1024)

	// Map to hold the data for each block number.
	blocks := make(map[uint16][]byte)

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

		// Store the block data by its block number.
		blocks[tftpData.Block] = tftpData.Data

		helper.ColorPrintln("blue", fmt.Sprintf("Opcode: %d, Block: %d, Data Length: %v", tftpData.Opcode, tftpData.Block, len(tftpData.Data)))
		helper.ColorPrintln("red", fmt.Sprintf("Stored block %d; current total blocks: %v", tftpData.Block, len(blocks)))

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
	// Reorder blocks by their block number.
	var keys []int
	for key := range blocks {
		keys = append(keys, int(key))
	}
	sort.Ints(keys)

	// Reassemble full message in proper order.
	var fullMessage []byte
	for _, k := range keys {
		fullMessage = append(fullMessage, blocks[uint16(k)]...)
	}

	helper.ColorPrintln("green", fmt.Sprintf("Full message length: %v", len(fullMessage)))
	fmt.Println(fullMessage[:100])
	return fullMessage, nil
}

func main() {
	hostAddress := "localhost:3000"
	imgURL := "https://static.boredpanda.com/blog/wp-content/uploads/2020/07/funny-expressive-dog-corgi-genthecorgi-1-1-5f0ea719ea38a__700.jpg"

	// Creating connection with the server
	conn, err := net.Dial("tcp", hostAddress)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	sendImageURLTOServer(conn, imgURL)

	fullMessage, err := ReadImagePacket(conn)
	if err != nil {
		return
	}

	fmt.Println(len(fullMessage))
	saveImageToFile(fullMessage, "test"+".jpg")

	conn.Close()
	// time.Sleep(100 * time.Second)
}
