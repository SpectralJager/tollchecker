package receiver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/SpectralJager/tollchecker/obu"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

const (
	READ_BUFFER_SIZE   = 1024
	WRITE_BUFFER_SIZE  = 1024
	RECEIVER_CHAN_SIZE = 128
	OBU_TOPIC          = "obu_topic"
	PARTITION          = 0
)

type DataReceiver struct {
	msgch chan obu.OBUData
	conn  *websocket.Conn
	prod  *kafka.Conn
}

func NewDataReceiver() DataReceiver {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "broker:9092", OBU_TOPIC, PARTITION)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	return DataReceiver{
		msgch: make(chan obu.OBUData, RECEIVER_CHAN_SIZE),
		prod:  conn,
	}
}

func ProduceData(dr *DataReceiver, data obu.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = dr.prod.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return err
	}
	_, err = dr.prod.WriteMessages(
		kafka.Message{Value: b},
	)
	return err
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
		if err := ProduceData(dr, data); err != nil {
			log.Printf("kafka produce error: %e", err)
		}
	}
}
