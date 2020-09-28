package model

const FaceTableName = "faces"

type FaceTable struct {
	Ts               int64
	ImageUri         string
	CutboardImageUri string
}

func (this *FaceTable) TableName() string {
	return FaceTableName
}
