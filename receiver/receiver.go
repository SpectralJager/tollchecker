package receiver

import (
	"log"
	"net/http"

	"github.com/SpectralJager/tollchecker/obu"
	"github.com/gorilla/websocket"
)

const (
	READ_BUFFER_SIZE   = 1024
	WRITE_BUFFER_SIZE  = 1024
	RECEIVER_CHAN_SIZE = 128
)

type DataReceiver struct {
	msgch chan obu.OBUData
	conn  *websocket.Conn
}

func NewDataReceiver() DataReceiver {
	return DataReceiver{
		msgch: make(chan obu.OBUData, RECEIVER_CHAN_SIZE),
	}
}

func WSHandler(dr *DataReceiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  READ_BUFFER_SIZE,
			WriteBufferSize: WRITE_BUFFER_SIZE,
		}
		conn, err := upgrader.Upgrade(w, r, http.Header{})
		if err != nil {
			log.Fatal(err)
		}
		dr.conn = conn
		go WSReceiveLoop(dr)
	}
}

func WSReceiveLoop(dr *DataReceiver) {
	log.Println("obu client connected")
	for {
		var data obu.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Printf("read error: %e", err)
			continue
		}
		log.Printf("received obu data [%d]<lat %.2f :: long %.2f>", data.OBUID, data.Lat, data.Long)
		// dr.msgch <- data
	}
}
