package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	// Connect to Redis server
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	}

	// Send a PING command to Redis using RESP protocol
	// command := "*1\r\n$4\r\nping\r\n"
	// command := "*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n"
	command := "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n"
	_, err = conn.Write([]byte(command))
	if err != nil {
		fmt.Println("Error sending command to Redis:", err)
		return
	}

	// Read the response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	fmt.Println(response)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// if strings.HasPrefix(response, "$") {
	// 	// Read the actual bulk string data
	// 	data, err := reader.ReadString('\n')
	// 	if err != nil {
	// 		fmt.Println("Error reading bulk string data:", err)
	// 		return
	// 	}
	// 	response = data
	// }

	// Print the response
	fmt.Println("Response from Redis:", strings.TrimSpace(response))
}
