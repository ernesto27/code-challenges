package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
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

type WindowSize struct {
	seconds   int
	counter   int
	threshold int
	mu        sync.Mutex
}

func (window *WindowSize) addCounter() bool {
	window.mu.Lock()
	defer window.mu.Unlock()

	window.counter++
	if window.counter > window.threshold {
		fmt.Println("Too Many Requests")
		return false
	}
	return true
}

func (window *WindowSize) resetCounter() {
	for {
		time.Sleep(time.Duration(window.seconds) * time.Second)
		window.mu.Lock()
		window.counter = 0
		window.mu.Unlock()
	}
}

func main() {
	// window size
	window := WindowSize{
		seconds:   20,
		counter:   0,
		threshold: 5,
	}

	go window.resetCounter()

	http.HandleFunc("/window", func(w http.ResponseWriter, r *http.Request) {
		if window.addCounter() {
			w.Write([]byte("request success"))
		} else {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}
	})

	// token bucket
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
