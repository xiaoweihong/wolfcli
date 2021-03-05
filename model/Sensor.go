package model

type Sensor struct {
	Id        string
	Name      string
	OrdId     string
	OrgName   string
	Url       string
	Longitude float64
	Latitude  float64
	Type      int64
	Status    int64
}

type SensorConfig struct {
	Speed           int
	FramingStrategy int
	VideoStartTime  int64
	VideoEndTime    int64
	BaseTime        int64
}
