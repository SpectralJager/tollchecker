package calculator

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"time"

	"github.com/SpectralJager/tollchecker/config"
	"github.com/SpectralJager/tollchecker/obu"
	"github.com/segmentio/kafka-go"
)

type Distance struct {
	Value float64 `json:"value"`
	OBUID int     `json:"obuid"`
	Unix  int64   `json:"unix"`
}

func OBUConsumer() {
	TopicConsumer(config.OBU_TOPIC, config.OBU_PARTITION)
}

func TopicConsumer(topic string, partition int) {
	conn, err := kafka.DialLeader(context.Background(), config.KAFKA_NETWORK, config.KAFKA_URL, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	ConsumeLoop(conn)
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}

func ConsumeLoop(conn *kafka.Conn) {
	var prev obu.OBUData
	for {
		msg, err := conn.ReadMessage(1024)
		if err != nil {
			log.Println("failed to read messages:", err)
			continue
		}
		var data obu.OBUData
		err = json.Unmarshal(msg.Value, &data)
		if err != nil {
			log.Println("failed to unmarshal messages:", err)
			continue
		}
		log.Printf("received from kafka: %s", data.String())
		if prev.OBUID == 0 {
			prev = data
		}
		distance, err := CalculateDistance(prev, data)
		if err != nil {
			log.Println("failed to calculate distance:", err)
			continue
		}
		log.Printf("calculated distance: %.2f", distance.Value)

	}
}

func CalculateDistance(d1, d2 obu.OBUData) (Distance, error) {
	res := math.Sqrt(math.Pow(d2.Lat-d1.Lat, 2) + math.Pow(d2.Long-d1.Long, 2))
	return Distance{
		Value: res,
		OBUID: d2.OBUID,
		Unix:  time.Now().Unix(),
	}, nil
}
