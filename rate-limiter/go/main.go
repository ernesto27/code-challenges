package main

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// TODO ADD MUTEX
type RateLimiter struct {
	capacity int
	bucket   map[string]int
}

func (rateLimite RateLimiter) addToken() {
	for {
		time.Sleep(10 * time.Second)
		for k := range rateLimite.bucket {
			if rateLimite.bucket[k] < rateLimite.capacity {
				rateLimite.bucket[k] += 1
			}
		}
		fmt.Println(rateLimite.bucket)
	}
}

func main() {

	rateLimiter := RateLimiter{
		capacity: 5,
		bucket:   make(map[string]int),
	}

	go rateLimiter.addToken()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// get ip user
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Println(err)
		}

		if _, ok := rateLimiter.bucket[ip]; !ok {
			rateLimiter.bucket[ip] = rateLimiter.capacity
		} else {
			if rateLimiter.bucket[ip] == 0 {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
		}

		rateLimiter.bucket[ip]--
		// check bucket for ip user
		w.Write([]byte("rate-limit"))
	})

	// every minute, reset the bucket

	http.ListenAndServe(":8080", nil)
}
