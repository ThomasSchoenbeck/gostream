package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// type sse struct {
// 	Id     int    `json:"id"`
// 	Method int    `json:"method"`
// 	Body   string `json:"body"`
// 	Data   string `json:"data"`
// 	Time   string `json:"time"`
// 	Type   string `json:"type"`
// }

func randomNumber() int {
	min := 10
	max := 250

	return rand.Intn(max-min) + min
}

func StreamHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	for i := 0; i < 20; i++ {
		randNum := randomNumber()
		fmt.Fprintf(w, "Index: %d  ->  waited for %dms\n", i, randNum)
		// msg := sse{Id: i, Method: http.StatusOK, Body: "body data", Data: fmt.Sprintf("Index: %d  ->  waited for %dms\n", i, randNum), Time: time.Now().String(), Type: "message"}
		// json.NewEncoder(w).Encode(msg)
		flusher.Flush()
		time.Sleep(time.Duration(randNum) * time.Millisecond)
	}
	// w.WriteHeader(http.StatusInternalServerError)

	fmt.Fprintln(w, "done")
}

func init() {
}

func main() {
	r := mux.NewRouter()
	// r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/stream", StreamHandler).Methods("GET")
	// r.HandleFunc("/articles", ArticlesHandler)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server Running at Port 8000")
	log.Fatal(srv.ListenAndServe())
}
