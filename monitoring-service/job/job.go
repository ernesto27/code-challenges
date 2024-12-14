package main

import (
	"fmt"
	"monitoring/db"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	myDB, dbErr := db.NewMysql(host, user, password, port, database)
	if dbErr != nil {
		panic(dbErr)
	}

	checkURLs(myDB)

	os.Exit(1)
	// ticker := time.NewTicker(1 * time.Minute)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			checkURLs(myDB)
		case <-quit:
			fmt.Println("Shutting down")
			return
		}
	}

}

func checkURLs(myDB *db.Mysql) {
	urls, err := myDB.GetURLs()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)

		go func(url db.URL) {
			fmt.Println("Checking URL: ", url.URL)
			defer wg.Done()
			resp, duration, err := sendHEADRequest(url.URL)
			isAlive := 1

			var statusCode int
			if err != nil {
				fmt.Println(err)
				statusCode = http.StatusBadRequest
				isAlive = 0
			} else {
				statusCode = resp.StatusCode
			}

			ttbf, whole := getTTFBWholeResponse(url.URL)

			err = myDB.CreateURLHealthCheck(url.ID, statusCode, int(duration.Milliseconds()), ttbf, whole, isAlive)
			if err != nil {
				fmt.Println(err)
			}

		}(url)
	}

	wg.Wait()
	fmt.Println("All URLs checked")
}

func sendHEADRequest(url string) (*http.Response, time.Duration, error) {
	start := time.Now()
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	duration := time.Since(start)

	return resp, duration, nil
}

func getTTFBWholeResponse(url string) (int, int) {
	start := time.Now()

	transport := &customTransport{Transport: http.DefaultTransport}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return 0, 0
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return 0, 0
	}
	defer resp.Body.Close()

	end := time.Now()

	totalResponseTime := end.Sub(start)
	fmt.Printf("TTFB: %v\n", transport.TTFB)
	fmt.Printf("TTFB (milliseconds): %d\n", int(transport.TTFB.Milliseconds()))
	fmt.Printf("Total Response Time: %v\n", totalResponseTime)
	fmt.Printf("Total Response Time (milliseconds): %d\n", int(totalResponseTime.Milliseconds()))

	return int(transport.TTFB.Milliseconds()), int(totalResponseTime.Milliseconds())
}

type customTransport struct {
	Transport http.RoundTripper
	TTFB      time.Duration
}

func (c *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	startTTFB := time.Now()
	resp, err := c.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	c.TTFB = time.Since(startTTFB)
	return resp, nil
}
