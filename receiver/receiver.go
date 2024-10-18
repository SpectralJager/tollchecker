package receiver

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/SpectralJager/tollchecker/config"
	"github.com/SpectralJager/tollchecker/obu"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

const (
	READ_BUFFER_SIZE   = 1024
	WRITE_BUFFER_SIZE  = 1024
	RECEIVER_CHAN_SIZE = 128
)

type DataReceiver struct {
	msgch chan obu.OBUData
	conn  *websocket.Conn
	prod  *kafka.Conn
	lg    *slog.Logger
}

func NewDataReceiver() DataReceiver {
	conn, err := kafka.DialLeader(context.Background(), "tcp", config.KAFKA_URL, config.OBU_TOPIC, config.OBU_PARTITION)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	return DataReceiver{
		msgch: make(chan obu.OBUData, RECEIVER_CHAN_SIZE),
		prod:  conn,
		lg:    slog.New(slog.NewTextHandler(os.Stdout, nil)),
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
	dr.lg.Info("obu client connected")
	producer := ProduceDataWithLogging(dr.lg)
	for {
		var data obu.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			dr.lg.Error("read error", "error", err)
			continue
		}
		dr.lg.Info("received obu data", "data", data.String())
		if err := producer(dr.prod, data); err != nil {
			dr.lg.Error("kafka produce error", "error", err)
		}
	}
}
