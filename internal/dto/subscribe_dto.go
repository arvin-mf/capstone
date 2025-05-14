package dto

type SubscribePeriodicData struct {
	Bpm                float32 `json:"bpm"`
	BodyTemperature    float32 `json:"suhu_objek"`
	AmbientTemperature float32 `json:"suhu_lingkungan"`
	Status             int     `json:"status"`
}

type SubscribePerpetualData struct {
	RawEcg float32 `json:"raw_ecg"`
}
