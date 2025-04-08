package main

import (
	"assignment-2/helper"
	"fmt"
	"log"
	"net"
	"os"
	"slices"
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

	chunk := make([]byte, 516)

	var fullMessage []byte
	imageBytesBlocks := make(map[uint16][]byte)

	for {

		n, err := conn.Read(chunk)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		tftpData, err := DeserializeTFTPDATA(chunk[:n])
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		helper.ColorPrintln("blue", fmt.Sprintf("Opcode: %d, Block: %d, Data Length: %v", tftpData.Opcode, tftpData.Block, len(tftpData.Data)))

		imageBytesBlocks[tftpData.Block] = tftpData.Data

		tftpAckPacketBytes, err := CreateTFTPACKPacket(tftpData.Block)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		fmt.Println("Sending Data: ", tftpAckPacketBytes)

		_, err = conn.Write(tftpAckPacketBytes)
		if err != nil {
			log.Fatal("Error sending ACK:", err)
			return nil, err
		}
		fmt.Printf("Sent ACK for block %d\n", tftpData.Block)
		helper.ColorPrintln("cyan", "The value of n is "+fmt.Sprint(n))
		if n != 516 {
			break
		}
	}

	// unAckLength := make([]byte, 4)
	// _, err := io.ReadFull(conn, unAckLength)
	// if err != nil {
	// 	if err == io.EOF {
	// 		fmt.Println("Client closed Connection")
	// 	}
	// 	fmt.Println("Read full Error: ", err)
	// }

	// chuckLength := binary.BigEndian.Uint32(unAckLength)
	// if chuckLength == 0 {
	// 	fmt.Println("Invalid chuckSize")
	// }

	// fmt.Println("current len: ", len(imageBytesBlocks))

	for {
		helper.ColorPrintln("red", "checkpoint1")
		n, err := conn.Read(chunk)
		if err != nil {
			break
		}

		helper.ColorPrintln("red", "checkpoint2")
		tftpData, err := DeserializeTFTPDATA(chunk[:n])
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		imageBytesBlocks[tftpData.Block] = tftpData.Data

		helper.ColorPrintln("red", "checkpoint3")
		tftpAckPacketBytes, err := CreateTFTPACKPacket(tftpData.Block)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		fmt.Println("Sending Data: ", tftpAckPacketBytes)

		helper.ColorPrintln("red", "checkpoint4")
		_, err = conn.Write(tftpAckPacketBytes)
		if err != nil {
			log.Fatal("Error sending ACK:", err)
			return nil, err
		}
		fmt.Printf("Sent ACK for block %d\n", tftpData.Block)
		helper.ColorPrintln("cyan", "The value of n is "+fmt.Sprint(n))
		fmt.Println("current len: ", len(imageBytesBlocks))
	}

	// To do next :
	// 1. What if the block number are duplicated
	// 2. What if an error packet is sent
	helper.ColorPrintln("green", "End of recieving all data")

	for key, value := range imageBytesBlocks {
		fmt.Printf("%v Key: %v\n", key, len(value))
	}

	fullMessage = reconstructImage(imageBytesBlocks)

	return fullMessage, nil
}

func reconstructImage(imageBytesBlock map[uint16][]byte) []byte {
	keys := make([]uint16, 0, len(imageBytesBlock))
	var tempL int
	for k, v := range imageBytesBlock {
		keys = append(keys, k)
		tempL += len(v)
	}
	helper.ColorPrintln("blue", "Here:"+fmt.Sprint(tempL))
	slices.Sort(keys)
	fmt.Println(keys)

	var result []byte
	for _, key := range keys {
		result = append(result, imageBytesBlock[key]...)

	}

	return result
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
	fmt.Println(fullMessage[:100])
	conn.Close()
	// time.Sleep(100 * time.Second)
}
