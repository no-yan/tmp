package main

import (
	"fmt"
	"net/http"
	"time"
)

func slowServer(w http.ResponseWriter, r *http.Request) {
	data := []byte("This is a test of slow data transmission.\n")
	for i := 0; i < 10; i++ {
		w.Write(data)
		w.(http.Flusher).Flush()
		time.Sleep(1 * time.Second)
	}
	fmt.Fprintln(w, "End of data.")
}

func main() {
	http.HandleFunc("/", slowServer)
	fmt.Println("Starting slow server on :8080")
	http.ListenAndServe(":8080", nil)
}
