package receiver

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/SpectralJager/tollchecker/obu"
	"github.com/segmentio/kafka-go"
)

type ProducerDataFunc func(*kafka.Conn, obu.OBUData) error

func ProduceData(conn *kafka.Conn, data obu.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return err
	}
	_, err = conn.WriteMessages(
		kafka.Message{Value: b},
	)
	return err

}

func ProduceDataWithLogging(l *slog.Logger) ProducerDataFunc {
	return func(conn *kafka.Conn, data obu.OBUData) error {
		l.Info("produced to kafka", "data", data.String())
		return ProduceData(conn, data)
	}
}
