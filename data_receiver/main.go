package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pttrulez/go-microservices/types"
)

func main() {
	// Produce messages to topic (asynchronously)
	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	recv.produceData(types.OBUData{})
	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":8080", nil)
}

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn

	// Продьюсер данных отправляет данные дальше. В нашем случае это кафка
	// Мы ожидаем, что продьюсер реализует интерфейс DataProducer:
	prod DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p          DataProducer
		err        error
		kafkaTopic = "obudata"
	)
	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}
	p = NewLogMiddleware(p)
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
		prod:  p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.prod.ProduceData(data)
}

// Хэндлер для получения данных от OBU
func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	// создаем websocket соединение
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	dr.conn = conn

	// и запускаем бесконечный цикл для получения данных из сокета
	go dr.wsReceiveLoop()
}

// Цикл для получения данных от OBU
func (dr *DataReceiver) wsReceiveLoop() {
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue
		}

		// отправляем данные далее в кафку
		if err := dr.produceData(data); err != nil {
			fmt.Println("kafka produce error:", err)
		}
	}
}
