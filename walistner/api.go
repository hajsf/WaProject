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

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Internal error", 500)
		return
	}

	fmt.Println("Client connected from IP:", r.RemoteAddr)
	// fmt.Println(len(p.connection), "new connection recieved")
	if len(p.connection) > 0 {
		fmt.Fprint(w, "event: notification\ndata: Connection is opened in another browser/tap ...\n\n")
		flusher.Flush()
		/*	w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			resp := make(map[string]string)
			resp["message"] = "Resource Not Found"
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResp)
			return
		*/
	}
	p.connection <- struct{}{}

	fmt.Fprint(w, "event: notification\ndata: Connecting to WhatsApp server ...\n\n")
	flusher.Flush()

	// Connect to the WhatsApp client
	go Connect()

	for {
		select {
		case data := <-p.data:
			// fmt.Println("SSE data recieved")

			switch {
			case len(data.event) > 0:
				fmt.Fprintf(w, "event: %v\ndata: %v\n\n", data.event, data.message)
			case len(data.event) == 0:
				fmt.Fprintf(w, "data: %v\n\n", data.message)
			}
			flusher.Flush()
		case <-r.Context().Done():
			<-p.connection
			fmt.Println("Connection closed from IP:", r.RemoteAddr)
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
