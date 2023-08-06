package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/idomath/toll-calculator/types"
	"log"
	"net/http"
)

func main() {
	receiver := NewDataReceiver()
	http.HandleFunc("/ws", receiver.handleWebSocket)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgch chan types.ObuData
	conn  *websocket.Conn
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgch: make(chan types.ObuData, 128),
	}
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
		fmt.Printf("received OBU data from [%d] :: <lat %.2f, long %.2f> \n", data.ObuId, data.Lat, data.Long)
		dr.msgch <- data
	}
}
