package main

import (
	"log"
	"net/http"

	"github.com/SpectralJager/tollchecker/receiver"
)

func main() {
	recv := receiver.NewDataReceiver()
	http.HandleFunc("/ws", receiver.WSHandler(&recv))
	log.Println("receiver started")
	if err := http.ListenAndServe(":30000", nil); err != nil {
		log.Fatal(err)
	}
}
