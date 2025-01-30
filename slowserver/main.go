package main

import (
	"fmt"
	"net/http"
	"time"
)

func slowServer(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s: %s\n", time.Now(), r.URL)
	data := []byte("This is a test of slow data transmission.\n")
	for i := 0; i < 10; i++ {
		w.Write(data)
		w.(http.Flusher).Flush()
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Fprintln(w, "End of data.")
}

func errorServer(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s: %s\n", time.Now(), r.URL)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func main() {
	http.HandleFunc("/error", errorServer)
	http.HandleFunc("/", slowServer)
	fmt.Println("Starting slow server on :8080")
	http.ListenAndServe(":8080", nil)
}
