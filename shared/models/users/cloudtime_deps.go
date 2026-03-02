package model

type CloudtimeDepartments struct {
	DepSerial int64  `gorm:"column:dep_serial" json:"dep_serial"`
	DepName   string `gorm:"column:dep_name" json:"dep_name"`
	DepNo     string `gorm:"column:dep_no" json:"dep_no"`
}

func (CloudtimeDepartments) TableName() string {
	return "dt_dep"
}
