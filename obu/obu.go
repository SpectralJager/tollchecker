package obu

import (
	"math"
	"math/rand/v2"
	"time"
)

const (
	SendInterval = time.Second * 60
	WSEndpoint   = "ws://receiver:30000/ws"
)

type OBUData struct {
	OBUID int     `json:"obuid"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

func GenLocation() (float64, float64) {
	return GenCoord(), GenCoord()
}

func GenCoord() float64 {
	return rand.Float64()*180 + 1
}

func GenerateOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.IntN(math.MaxInt)
	}
	return ids
}
