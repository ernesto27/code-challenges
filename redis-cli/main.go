package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

type RedisClient struct {
	conn      net.Conn
	host      string
	port      string
	commands  map[string]Command
	cmdValues []string
}

type Command struct {
	Arguments []Argument `json:"arguments"`
}

type Argument struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	KeySpecIndex int    `json:"key_spec_index"`
	Optional     bool   `json:"optional"`
	Token        string `json:"token"`
	Arguments    []Argument
}

func NewRedisClient(host string, port string) (*RedisClient, error) {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return nil, err
	}

	cmdValues := []string{"gets", "set"}
	commands := make(map[string]Command)

	for _, cmd := range cmdValues {
		file, err := os.Open("commands/" + cmd + ".json")
		if err != nil {
			fmt.Println("Error opening file:", err)
			continue
		}
		defer file.Close()

		byteValue, err := io.ReadAll(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
			continue
		}
		var data map[string]Command
		err = json.Unmarshal(byteValue, &data)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			continue
		}

		commands[cmd] = data[strings.ToUpper(cmd)]
	}

	return &RedisClient{
		conn:      conn,
		host:      host,
		port:      port,
		commands:  commands,
		cmdValues: cmdValues,
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

func (r *RedisClient) getSuggestion(command string) string {
	getCommand := r.commands[command]

	var suggest string
	for _, arg := range getCommand.Arguments {
		if arg.Type == "oneof" {

			suggest += "["
			for idx, childArg := range arg.Arguments {
				var extra string
				if !strings.EqualFold(childArg.Token, childArg.Name) {
					extra = " " + childArg.Name
				}
				if idx == len(arg.Arguments)-1 {
					suggest += childArg.Token + extra
				} else {
					suggest += fmt.Sprintf("%s%s|", childArg.Token, extra)
				}
			}
			suggest += "] "
		} else {
			if arg.Optional {
				suggest += fmt.Sprintf("[%s] ", arg.Token)
			} else {
				suggest += fmt.Sprintf("%s ", arg.Name)
			}
		}
	}

	if suggest == "" {
		return ""
	}

	return "\033[90m" + suggest + "\033[0m"
}

func main() {

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
	var suggest string
	for {
		fmt.Print("\033[2K\033[G")

		fmt.Print("$localhost:6379> ", value+suggest)
		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("Error getting key:", err)
			break
		}

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

		case keyboard.KeyCtrlC:
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

		lowerValue := strings.ToLower(value)

		if lowerValue == "set " || lowerValue == "get " {
			suggest = redisClient.getSuggestion(lowerValue[:len(lowerValue)-1])
		} else {
			suggest = ""
		}

	}
}
