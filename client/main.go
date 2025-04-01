package main

import (
	"encoding/binary"
	"fmt"
	"io"
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

func ReadImagePacket(conn net.Conn) ([]byte, error) {
	// Reading image size
	headerBytes := make([]byte, 4)
	conn.Read(headerBytes)
	fmt.Println("Received Header:", headerBytes)
	imgByteLength := uint32(binary.BigEndian.Uint32(headerBytes))
	fmt.Println(imgByteLength)
	chunk := make([]byte, 1024)

	var fullMessage []byte

	for uint32(len(fullMessage)) < uint32(imgByteLength) {
		remaining := imgByteLength - uint32(len(fullMessage))

		currentChunkSize := uint32(len(chunk))
		if currentChunkSize > remaining {
			currentChunkSize = remaining
		}

		chunk := make([]byte, currentChunkSize)
		n, err := conn.Read(chunk)
		if err != nil {
			if err == io.EOF && uint32(len(fullMessage)) == imgByteLength {
				return nil, err
			}
			fmt.Println("Error reading chuck: ", err)
			return nil, err
		}

		fullMessage = append(fullMessage, chunk[:n]...)
	}

	fmt.Println(len(fullMessage))
	return fullMessage, nil
}

func main() {
	hostAddress := "localhost:3000"
	// imgURL := "https://static.boredpanda.com/blog/wp-content/uploads/2020/07/funny-expressive-dog-corgi-genthecorgi-1-1-5f0ea719ea38a__700.jpg"

	// Creating connection with the server
	conn, err := net.Dial("tcp", hostAddress)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// sendImageURLTOServer(conn, imgURL)

	// fullMessage, err := ReadImagePacket(conn)
	// if err != nil {
	// 	return
	// }

	// saveImageToFile(fullMessage, "test" + ".jpg")

	chunk := make([]byte, 1024)
	n, err := conn.Read(chunk)
	if err != nil {
		log.Fatal(err)
		return
	}

	tftp, err := DeserializeTFTPDATA(chunk[:n])
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println(string(chunk[:n]))
	fmt.Println(tftp.Block, ",", tftp.Opcode)

	tftpAckPacketBytes, err := CreateTFTPACKPacket()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Sending Data: ", tftpAckPacketBytes)
	fmt.Println("Length of the data:", len(tftpAckPacketBytes))

	_, err = conn.Write(tftpAckPacketBytes)
	if err != nil {
		log.Fatal("Error sending ACK:", err)
		return
	}
	// fmt.Printf("Sent ACK for block %d\n", tftp.Block)

	conn.Close()
	// time.Sleep(100 * time.Second)
}
