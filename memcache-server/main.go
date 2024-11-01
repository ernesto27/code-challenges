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
	"time"
)

const GET = "get"
const SET = "set"
const ADD = "add"
const REPLACE = "replace"
const PREPEND = "prepend"
const APPEND = "append"

type MemcacheServer struct {
	port    string
	clients map[net.Conn]Command
	data    map[string]Data
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
	created    time.Time
}

type Data struct {
	value      string
	expiration int
	createAt   time.Time
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

func (m *MemcacheServer) setData(key string, data Data) {
	m.mu.Lock()
	m.data[key] = data
	m.mu.Unlock()
}

func (m *MemcacheServer) getData(key string) (Data, error) {
	data, ok := m.data[key]
	if !ok {
		return Data{}, fmt.Errorf("key not found")
	}
	return data, nil
}

func (m *MemcacheServer) deleteData(key string) {
	m.mu.Lock()
	delete(m.data, key)
	m.mu.Unlock()
}

func (m *MemcacheServer) response(conn net.Conn, response string) {
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Println(err)
	}
}

func (m *MemcacheServer) removeCarriageReturn(s string) string {
	return strings.TrimSuffix(s, "\r")
}

func (m *MemcacheServer) isValidCommand(conn net.Conn) bool {
	resp := m.clients[conn].name != "" && m.clients[conn].value != ""
	return resp
}

func (m *MemcacheServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			continue
		}
		command = strings.TrimSuffix(command, "\n")
		err = m.parseCommand(conn, command)

		if err == nil {
			switch m.clients[conn].name {
			case SET:
				if m.isValidCommand(conn) {
					m.setData(m.clients[conn].key, Data{
						value:      m.clients[conn].value,
						expiration: m.clients[conn].expiration,
						createAt:   time.Now(),
					})

					if !m.clients[conn].noReply {
						response := "STORED\r\n"
						m.response(conn, response)
					}
					m.clients[conn] = Command{}
				}

			case ADD:
				if m.isValidCommand(conn) {
					if _, ok := m.data[m.clients[conn].key]; !ok {
						m.setData(m.clients[conn].key, Data{
							value:      m.clients[conn].value,
							expiration: m.clients[conn].expiration,
							createAt:   time.Now(),
						})
						response := "STORED\r\n"
						m.response(conn, response)
					} else {
						response := "NOT_STORED\r\n"
						m.response(conn, response)
					}
					m.clients[conn] = Command{}
				}

			case GET:
				data, err := m.getData(m.clients[conn].key)
				if err != nil {
					response := "END\r\n"
					m.response(conn, response)
				} else {
					if m.checkExpiration(data) {
						response := "END\r\n"
						m.response(conn, response)
						m.deleteData(m.clients[conn].key)
						m.clients[conn] = Command{}
						continue
					}

					response := fmt.Sprintf("VALUE %s 0 %d\r\n%s\r\nEND\r\n", m.clients[conn].key, m.clients[conn].byteCount, data.value)
					m.response(conn, response)
				}
				m.clients[conn] = Command{}
			case REPLACE:
				if m.isValidCommand(conn) {
					if _, ok := m.data[m.clients[conn].key]; ok {
						m.setData(m.clients[conn].key, Data{
							value:      m.clients[conn].value,
							expiration: m.clients[conn].expiration,
							createAt:   time.Now(),
						})
						response := "STORED\r\n"
						m.response(conn, response)
					} else {
						response := "NOT_STORED\r\n"
						m.response(conn, response)
					}
					m.clients[conn] = Command{}
				}

			case APPEND:
				if m.isValidCommand(conn) {
					data, ok := m.getData(m.clients[conn].key)
					if ok == nil {
						data.value = data.value + m.clients[conn].value
						m.setData(m.removeCarriageReturn(m.clients[conn].key), data)
						response := "STORED\r\n"
						m.response(conn, response)
					} else {
						response := "NOT_STORED\r\n"
						m.response(conn, response)
					}
					m.clients[conn] = Command{}
				}
			case PREPEND:
				if m.isValidCommand(conn) {
					data, ok := m.getData(m.clients[conn].key)
					if ok == nil {
						data.value = m.clients[conn].value + data.value
						m.setData(m.removeCarriageReturn(m.clients[conn].key), data)
						response := "STORED\r\n"
						m.response(conn, response)
					} else {
						response := "NOT_STORED\r\n"
						m.response(conn, response)
					}
					m.clients[conn] = Command{}
				}

			}
		}
	}
}

func (m *MemcacheServer) parseCommand(conn net.Conn, command string) error {
	args := strings.Split(command, " ")

	if len(args) == 2 {
		if args[0] == GET {
			key := args[1]
			key = m.removeCarriageReturn(key)
			cmd := &Command{
				name: GET,
				key:  key,
			}
			m.clients[conn] = *cmd
		} else {
			return fmt.Errorf("invalid command")
		}
		return nil
	}

	if len(args) == 5 || len(args) == 6 {
		if (args[0] == SET || args[0] == ADD ||
			args[0] == REPLACE || args[0] == APPEND || args[0] == PREPEND) &&
			args[1] != "" {

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

			expirationValue, err := strconv.Atoi(args[3])
			if err != nil {
				return err
			}

			cmd := &Command{
				name:       args[0],
				key:        args[1],
				flags:      0,
				expiration: expirationValue,
				byteCount:  byteCountValue,
				noReply:    noReply,
			}
			m.clients[conn] = *cmd
		} else {
			return fmt.Errorf("invalid command")
		}
		return nil
	}

	if m.clients[conn].name != "" {
		cmd := m.clients[conn]
		cmd.value = m.removeCarriageReturn(command)
		cmd.created = time.Now().Add(-3 * time.Hour)
		m.clients[conn] = cmd
	}

	return nil
}

func (m *MemcacheServer) checkExpiration(data Data) bool {
	if data.expiration == 0 {
		return false
	}

	duration := time.Since(data.createAt)
	seconds := int(duration.Seconds())

	return seconds > data.expiration
}

func main() {
	port := flag.String("p", "11211", "Port to listen on")
	flag.Parse()

	memecacheServer := &MemcacheServer{
		port:    *port,
		clients: make(map[net.Conn]Command),
		data:    make(map[string]Data),
	}
	if err := memecacheServer.server(); err != nil {
		log.Fatal(err)
	}
}
