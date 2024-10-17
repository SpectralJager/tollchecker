package main

import (
	"log"
	"time"

	"github.com/SpectralJager/tollchecker/obu"
	"github.com/gorilla/websocket"
)

func main() {
	obuIDS := obu.GenerateOBUIDS(20)
	conn, _, err := websocket.DefaultDialer.Dial(obu.WSEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for _, id := range obuIDS {
			lat, long := obu.GenLocation()
			data := obu.OBUData{
				OBUID: id,
				Lat:   lat,
				Long:  long,
			}
			// fmt.Printf("%+v\n", data)
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(obu.SendInterval)
	}
}
