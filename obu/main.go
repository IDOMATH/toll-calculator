package main

import (
	"github.com/gorilla/websocket"
	"github.com/idomath/toll-calculator/types"
	"log"
	"math"
	"math/rand"
	"time"
)

const sendInterval = time.Second

const wsEndpoint = "ws://127.0.0.1:30000/ws"

type ObuData struct {
	ObuId int     `json:"obuId"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

func genCoordinate() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func genLocation() (float64, float64) {
	return genCoordinate(), genCoordinate()
}

func generateObuIds(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}

func main() {
	obuIds := generateObuIds(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i < len(obuIds); i++ {
			lat, long := genLocation()
			data := types.ObuData{
				ObuId: obuIds[i],
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

func init() {
	rand.Seed(time.Now().UnixNano())
}
