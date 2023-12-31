package types

type ObuData struct {
	ObuId int     `json:"obuId"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

type Distance struct {
	Value float64 `json:"value"`
	ObuId int     `json:"obuId"`
	Unix  int64   `json:"unix"`
}

type Invoice struct {
	ObuId         int     `json:"obuId"`
	TotalDistance float64 `json:"totalDistance"`
	TotalAmount   float64 `json:"totalAmount"`
}
