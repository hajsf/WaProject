package main

import (
	"fmt"
	"net/http"
)

type sseData struct {
	event, message string
}
type DataPasser struct {
	data       chan sseData
	logs       chan string
	connection chan struct{} // To control maximum allowed clients connections
}

func (p *DataPasser) HandleSignal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	setupCORS(&w, r)

	fmt.Println("Client connected from IP:", r.RemoteAddr)

	p.connection <- struct{}{}
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Internal error", 500)
		return
	}

	fmt.Fprint(w, "event: notification\ndata: Connection to WhatsApp server ...\n\n")
	flusher.Flush()

	// Connect to the WhatsApp client
	go Connect()

	for {
		select {
		case data := <-p.data:
			fmt.Println("recieved")

			switch {
			case len(data.event) > 0:
				fmt.Fprintf(w, "event: %v\ndata: %v\n\n", data.event, data.message)
			case len(data.event) == 0:
				fmt.Fprintf(w, "data: %v\n\n", data.message)
			}
			flusher.Flush()
		case <-r.Context().Done():
			<-p.connection
			fmt.Println("Connection closed")
			return
		}
	}
}

func setupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Cache-Control", "no-cache")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
