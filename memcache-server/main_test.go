package main

import (
	"net"
	"strings"
	"testing"
	"time"
)

func TestMemcacheServer(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		expected string
	}{
		{
			name:     "SET command",
			commands: []string{"set key 0 60 5\r\nvalue\r\n"},
			expected: "STORED\r\n",
		},
		{
			name:     "GET command",
			commands: []string{"get key\r\n"},
			expected: "VALUE key 0 0\r\nvalue\r\nEND\r\n",
		},
		{
			name:     "ADD command",
			commands: []string{"add key 0 60 5\r\nvalue\r\n"},
			expected: "STORED\r\n",
		},
		{
			name:     "REPLACE command",
			commands: []string{"set key 0 60 5\r\nvalue\r\n", "replace key 0 60 7\r\nnewval\r\n"},
			expected: "STORED\r\n",
		},
		{
			name:     "APPEND command",
			commands: []string{"set key 0 60 5\r\nvalue\r\n", "append key 0 60 3\r\n123\r\n"},
			expected: "STORED\r\n",
		},
		{
			name:     "PREPEND command",
			commands: []string{"set key 0 60 5\r\nvalue\r\n", "prepend key 0 60 3\r\n123\r\n"},
			expected: "STORED\r\n",
		},
		{
			name:     "GET non-existent key",
			commands: []string{"get non_existent_key\r\n"},
			expected: "END\r\n",
		},
		{
			name:     "ADD existing key",
			commands: []string{"set key 0 60 5\r\nvalue\r\n", "add key 0 60 5\r\nvalue\r\n"},
			expected: "NOT_STORED\r\n",
		},
		{
			name:     "REPLACE non-existent key",
			commands: []string{"replace non_existent_key 0 60 5\r\nvalue\r\n"},
			expected: "NOT_STORED\r\n",
		},
		{
			name:     "APPEND non-existent key",
			commands: []string{"append non_existent_key 0 60 5\r\nvalue\r\n"},
			expected: "NOT_STORED\r\n",
		},
		{
			name:     "PREPEND non-existent key",
			commands: []string{"prepend non_existent_key 0 60 5\r\nvalue\r\n"},
			expected: "NOT_STORED\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &MemcacheServer{
				port:    "6666",
				clients: make(map[net.Conn]Command),
				data:    make(map[string]Data),
			}

			go server.server()

			conn, err := net.Dial("tcp", "localhost:6666")
			if err != nil {
				t.Fatalf("failed to connect to server: %v", err)
			}
			defer conn.Close()

			for _, cmd := range tt.commands {
				_, err := conn.Write([]byte(cmd))
				if err != nil {
					t.Fatalf("failed to write command to server: %v", err)
				}
			}

			time.Sleep(100 * time.Millisecond)

			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				t.Fatalf("failed to read response from server: %v", err)
			}

			got := string(buf[:n])
			if !strings.Contains(got, tt.expected) {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
