package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const sendInterval = time.Second

const wsEndpoint = "wd://127.0.0.1:3000/ws"

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

func genObuIds(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}

func main() {
	for {
		fmt.Println(genLocation())
		time.Sleep(sendInterval)
	}

}

func init() {
	rand.Seed(time.Now().UnixNano())
}
