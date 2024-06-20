package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pttrulez/go-microservices/types"
)

const wsEndpoint = "ws://localhost:8080/ws"
const sendInterval = time.Second * 5

var r *rand.Rand

// Генерирует мок данные от OBU (от машин)
func main() {
	obuids := generateOBUIDS(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		fmt.Println("portion of data")
		for i := 0; i < len(obuids); i++ {
			lat, long := genLatLong()
			data := types.OBUData{
				OBUID: obuids[i],
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(sendInterval)
	}
}

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func genCoord() float64 {
	n := float64(r.Intn(100) + 1)
	f := r.Float64()
	return n + f
}

func generateOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = r.Intn(math.MaxInt)
	}
	return ids
}

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}
