package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

func isHostPortOpen2(host, port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 3*time.Second)
	if err == nil {
		conn.Close()
		return true
	}

	return false
}

func main() {
	host := flag.String("host", "", "the host to connect to")
	port := flag.String("port", "", "port")
	flag.Parse()

	if *host == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	hosts := strings.Split(*host, ",")

	if len(hosts) == 1 && strings.Contains(*host, "*") {
		hosts = []string{}
		if strings.Contains(*host, "*") {
			for i := 2; i < 255; i++ {
				newIP := strings.Replace(*host, "*", fmt.Sprint(i), 1)
				hosts = append(hosts, newIP)
			}
		}
	}

	if len(hosts) >= 2 || strings.Contains(*host, "*") {
		hostsAlive := make(map[string]string)
		for _, host := range hosts {
			fmt.Println("Scanning host swap ", host)
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			for i := 0; i < 65535; i++ {
				wg.Add(1)
				go func(host string, i int) {
					defer wg.Done()
					select {
					case <-ctx.Done():
						return
					default:
					}

					valid := isHostPortOpen2(host, fmt.Sprint(i))
					if valid {
						fmt.Fprintf(os.Stdout, "Host %s is alive\n", host)
						hostsAlive[host] = "Host " + host + " is alive"
						cancel()
						return
					}

				}(host, i)
			}
			wg.Wait()
		}

		for _, h := range hostsAlive {
			fmt.Println(h)
		}

		os.Exit(0)
	}

	if *port != "" {

		fmt.Fprintf(os.Stdout, "Connecting to %s:%s\n", *host, *port)

		valid := isHostPortOpen(*host, *port)
		if !valid {
			fmt.Println("Port: " + *port + " is closed")
			os.Exit(1)
		}

		fmt.Println("Port: " + *port + " is open")

	} else {
		var wg sync.WaitGroup

		fmt.Fprintf(os.Stdout, "Scanning host %s\n", *host)
		for i := 1; i <= 65535; i++ {
			wg.Add(1)
			go func(host string, port int) {
				defer wg.Done()
				valid := isHostPortOpen(host, fmt.Sprint(port))
				if valid {
					fmt.Println("Port: " + fmt.Sprintf("%d", port) + " is open")
				}

			}(*host, i)
		}

		wg.Wait()
	}

}

func isHostPortOpen(host string, port string) bool {
	conn, err := net.DialTimeout("tcp", host+":"+port, time.Second*5)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}
