package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var hopByHopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"TE",
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func removeHopByHopHeaders(h http.Header) {
	for _, header := range hopByHopHeaders {
		h.Del(header)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	// fmt.Println("Request received from: ", r.RemoteAddr)
	// fmt.Println("Request target: ", r.RequestURI)
	// fmt.Println("Request url: ", r.URL)
	fmt.Printf("Client: %s request URL: %s", r.RemoteAddr, r.RequestURI)
	fmt.Println()

	bytes, err := os.ReadFile("forbidden-hosts.txt")
	if err != nil {
		fmt.Println("Error reading file")
	} else {
		for _, domain := range strings.Split(string(bytes), "\n") {
			if strings.Contains(r.RequestURI, domain) {
				http.Error(w, "Website not allowed: "+domain, http.StatusForbidden)
				return
			}
		}
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", r.RequestURI, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	removeHopByHopHeaders(r.Header)

	req.Header.Add("X-Forwarded-For", r.RemoteAddr)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data []byte
	if err != nil {
		w.Write([]byte("Error"))
		return
	} else {
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			w.Write([]byte("Error"))
			return
		}
	}
	resp.Body.Close()

	forbidden, err := os.ReadFile("banned-words.txt")
	if err != nil {
		fmt.Println("Error reading file")
	} else {
		for _, word := range strings.Split(string(forbidden), "\n") {
			if strings.Contains(strings.ToLower(string(data)), strings.ToLower(word)) {
				http.Error(w, "Website content not allowed.", http.StatusForbidden)
				return
			}
		}

	}

	fmt.Println(r.RemoteAddr, " ", resp.Status)
	fmt.Fprint(w, string(data))

}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8989", nil)

	fmt.Println("Server started on port 8989")
}
