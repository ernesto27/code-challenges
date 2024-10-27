package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

type MemcacheServer struct {
	port    string
	clients map[net.Conn]Command
}

type Command struct {
	name       string
	key        string
	flags      int
	expiration int
	byteCount  int
	value      string
}

func (m *MemcacheServer) server() error {
	listener, err := net.Listen("tcp", ":"+m.port)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		m.clients[conn] = Command{}
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Println(conn.RemoteAddr())
		go m.handleConnection(conn)

	}
}

func (m *MemcacheServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		command = strings.TrimSuffix(command, "\n")
		m.parseCommand(conn, command)
		fmt.Println(m.clients)

		response := "Hello, client!\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (m *MemcacheServer) parseCommand(conn net.Conn, command string) error {
	// set username 0 0 4

	if m.clients[conn].name != "" {
		cmd := m.clients[conn]
		cmd.value = command
		m.clients[conn] = cmd
	} else {

		args := strings.Split(command, " ")
		fmt.Println(args)

		if len(args) < 4 {
			return fmt.Errorf("invalid command")
		}

		cmd := &Command{
			name:       args[0],
			key:        args[1],
			flags:      0,
			expiration: 0,
			byteCount:  0,
		}

		m.clients[conn] = *cmd
	}

	return nil

}

func main() {
	port := flag.String("p", "11211", "Port to listen on")
	flag.Parse()

	memecacheServer := &MemcacheServer{
		port:    *port,
		clients: make(map[net.Conn]Command),
	}
	if err := memecacheServer.server(); err != nil {
		log.Fatal(err)
	}
}
