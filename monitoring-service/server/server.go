package main

import (
	"encoding/json"
	"fmt"
	"log"
	"monitoring/db"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	myDB, err := db.NewMysql(host, user, password, port, database)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../public/chart.html")
	})

	http.HandleFunc("/historical-data", func(w http.ResponseWriter, r *http.Request) {
		// get id query parameter
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id parameter is required", http.StatusBadRequest)
			return
		}

		startDate := r.URL.Query().Get("start")
		if startDate == "" {
			http.Error(w, "start parameter is required", http.StatusBadRequest)
			return
		}
		endDate := r.URL.Query().Get("end")
		if endDate == "" {
			http.Error(w, "end parameter is required", http.StatusBadRequest)
			return
		}

		idVal, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "id parameter must be an integer", http.StatusBadRequest)
			return
		}

		data, err := myDB.GetHistoricDataByURLID(idVal, startDate, endDate)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error getting historical data", http.StatusInternalServerError)
			return
		}

		urls, err := myDB.GetURLs()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error getting urls", http.StatusInternalServerError)
			return
		}

		data.URLs = urls

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)

	})

	port = "8080"

	fmt.Printf("Server is running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
