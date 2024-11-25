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

			err = myDB.CreateURLHealthCheck(url.ID, statusCode, int(duration.Milliseconds()), isAlive)
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
