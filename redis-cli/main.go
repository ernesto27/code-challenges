package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

type RedisClient struct {
	conn net.Conn
	host string
	port string
}

func NewRedisClient(host string, port string) (*RedisClient, error) {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return nil, err
	}
	return &RedisClient{
		conn: conn,
		host: host,
		port: port,
	}, nil
}

func (r *RedisClient) Close() {
	r.conn.Close()
}

func (r *RedisClient) buildCommand(args ...string) string {
	command := "*" + fmt.Sprint(len(args)) + "\r\n"
	for _, arg := range args {
		command += "$" + fmt.Sprint(len(arg)) + "\r\n" + arg + "\r\n"
	}

	return command
}

func (r *RedisClient) execute(args ...string) (string, error) {
	command := r.buildCommand(args...)

	_, err := r.conn.Write([]byte(command))
	if err != nil {
		return "", err
	}

	r.conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	reader := bufio.NewReader(r.conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(response, "$") {
		data, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		response = data
	}

	return response, nil
}

func main() {
	// Connect to Redis server
	redisClient, err := NewRedisClient("localhost", "6379")
	if err != nil {
		panic(err)
	}
	defer redisClient.Close()

	if err := keyboard.Open(); err != nil {
		fmt.Println("Error opening keyboard:", err)
		os.Exit(1)
	}
	defer keyboard.Close()

	var value string
	for {
		fmt.Print("\033[2K\033[G")

		suggest := "\033[90msuggest\033[0m"
		fmt.Print("$localhost:6379> ", value+" "+suggest)
		var err error
		var key keyboard.Key

		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("Error getting key:", err)
			break
		}
		//fmt.Printf("You pressed: rune %q, key %X\r\n", char, key)

		switch key {
		case keyboard.KeyEnter:
			parts := strings.Fields(value)

			if len(parts) == 0 {
				continue
			}

			resp, err := redisClient.execute(parts...)
			if err != nil {
				fmt.Println("Error executing command:", err)
				continue
			}
			fmt.Println()
			fmt.Println(resp)
			value = ""

		case keyboard.KeyEsc:
			fmt.Println("ESC pressed. Exiting...")
			return

		case keyboard.KeyBackspace2:
			if len(value) != 0 {
				value = value[:len(value)-1]
			}
		case keyboard.KeySpace:
			value += " "

		default:
			value += string(char)
		}

	}
}
