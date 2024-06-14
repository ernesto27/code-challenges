package main

import (
	"fmt"
	"io"
	"net/http"
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

	fmt.Println("Request received from: ", r.RemoteAddr)
	fmt.Println("Request target: ", r.RequestURI)

	client := &http.Client{}

	// Create a new request using http
	req, err := http.NewRequest("GET", r.RequestURI, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	removeHopByHopHeaders(r.Header)

	// add header to the forwarded request
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
		fmt.Println(string(data))
		if err != nil {
			w.Write([]byte("Error"))
			return
		}
	}
	resp.Body.Close()

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(string(data))

	fmt.Fprint(w, string(data))

}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8989", nil)
}
