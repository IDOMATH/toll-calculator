package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/idomath/toll-calculator/types"
	"log"
	"net/http"
)

var kafkaTopic = "obudata"

func main() {
	receiver, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", receiver.handleWebSocket)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgch    chan types.ObuData
	conn     *websocket.Conn
	producer DataProducer
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
		msgch:    make(chan types.ObuData, 128),
		producer: p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.ObuData) error {
	return dr.producer.ProduceData(data)
}

func (dr *DataReceiver) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.webSocketReceiveLoop()
}

func (dr *DataReceiver) webSocketReceiveLoop() {
	fmt.Println("New OBU connected to client")
	for {
		var data types.ObuData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue
		}
		fmt.Println("received message", data)
		if err := dr.produceData(data); err != nil {
			fmt.Println("kafka produce error:", err)
		}
	}
}
