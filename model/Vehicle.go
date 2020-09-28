package model

const VehicleTableName = "vehicle_capture"

type VehicleTable struct {
	Ts               int64
	ImageUri         string
	CutboardImageUri string
}

func (this *VehicleTable) TableName() string {
	return VehicleTableName
}
