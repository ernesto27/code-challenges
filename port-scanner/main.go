package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	host := flag.String("host", "", "the host to connect to")
	port := flag.String("port", "", "port")
	flag.Parse()

	if *host == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	hosts := strings.Split(*host, ",")
	if *port != "" {
		for _, h := range hosts {
			fmt.Fprintf(os.Stdout, "Connecting to %s:%s\n", h, *port)

			conn, err := net.DialTimeout("tcp", h+":"+*port, time.Second*5)
			if err != nil {
				continue
			}
			defer conn.Close()

			fmt.Println("Port: " + *port + " is open")
		}
	} else {
		var wg sync.WaitGroup
		for _, h := range hosts {
			fmt.Fprintf(os.Stdout, "Scanning host %s\n", h)
			for i := 1; i <= 65535; i++ {
				wg.Add(1)
				go func(host string, port int) {
					defer wg.Done()
					conn, err := net.DialTimeout("tcp", host+":"+fmt.Sprintf("%d", port), time.Second*2)
					if err != nil {
						return
					}
					defer conn.Close()

					fmt.Println("Port: " + fmt.Sprintf("%d", port) + " is open")
				}(h, i)
			}
		}

		wg.Wait()
	}

}
