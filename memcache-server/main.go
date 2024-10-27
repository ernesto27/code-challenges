package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

type MemcacheServer struct {
	port    string
	clients map[net.Conn]Command
	data    map[string]string
	mu      sync.Mutex
}

type Command struct {
	name       string
	key        string
	flags      int
	expiration int
	byteCount  int
	value      string
	noReply    bool
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

func (m *MemcacheServer) setData(key, value string) {
	m.mu.Lock()
	m.data[key] = value
	m.mu.Unlock()
}

func (m *MemcacheServer) getData(key string) (string, error) {
	value, ok := m.data[key]
	if !ok {
		return "", fmt.Errorf("key not found")
	}
	return value, nil
}

func (m *MemcacheServer) removeCarriageReturn(s string) string {
	return strings.TrimSuffix(s, "\r")
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
		err = m.parseCommand(conn, command)

		if err == nil {
			switch m.clients[conn].name {
			case "set":
				if m.clients[conn].name != "" && m.clients[conn].value != "" {
					m.setData(m.clients[conn].key, m.clients[conn].value)

					if !m.clients[conn].noReply {
						response := "STORED\r\n"
						_, err = conn.Write([]byte(response))
						if err != nil {
							log.Println(err)
							return
						}
					}
					m.clients[conn] = Command{}
				}

			case "get":
				value, err := m.getData(m.clients[conn].key)
				if err != nil {
					response := "END\r\n"
					_, err = conn.Write([]byte(response))
					if err != nil {
						log.Println(err)
						return
					}
				} else {
					response := fmt.Sprintf("VALUE %s 0 %d\r\n%s\r\nEND\r\n", m.clients[conn].key, m.clients[conn].byteCount, value)
					_, err = conn.Write([]byte(response))
					if err != nil {
						log.Println(err)
						return
					}
				}
				m.clients[conn] = Command{}
			}
		}
	}
}

func (m *MemcacheServer) parseCommand(conn net.Conn, command string) error {
	// set username 0 0 4
	args := strings.Split(command, " ")

	if len(args) == 2 {
		if args[0] == "get" {
			key := args[1]
			key = m.removeCarriageReturn(key)
			cmd := &Command{
				name: "get",
				key:  key,
			}
			m.clients[conn] = *cmd
		} else {
			return fmt.Errorf("invalid command")
		}

	} else if len(args) == 5 || len(args) == 6 {
		if args[0] == "set" && args[1] != "" {
			byteCount := args[4]
			byteCount = m.removeCarriageReturn(byteCount)
			byteCountValue, err := strconv.Atoi(byteCount)
			if err != nil {
				return err
			}

			noReply := false
			if len(args) == 6 {
				noReplyValue := m.removeCarriageReturn(args[5])
				if noReplyValue == "noreply" {
					noReply = true
				}
			}
			cmd := &Command{
				name:       args[0],
				key:        args[1],
				flags:      0,
				expiration: 0,
				byteCount:  byteCountValue,
				noReply:    noReply,
			}
			m.clients[conn] = *cmd
		} else {
			return fmt.Errorf("invalid command")
		}
	} else {
		if m.clients[conn].name != "" {
			cmd := m.clients[conn]
			cmd.value = command
			m.clients[conn] = cmd
		}
	}

	return nil

}

func main() {
	port := flag.String("p", "11211", "Port to listen on")
	flag.Parse()

	memecacheServer := &MemcacheServer{
		port:    *port,
		clients: make(map[net.Conn]Command),
		data:    make(map[string]string),
	}
	if err := memecacheServer.server(); err != nil {
		log.Fatal(err)
	}
}
