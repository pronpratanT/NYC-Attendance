package model

type CloudtimeDepartments struct {
	DepSerial int64  `json:"dep_serial"`
	DepName   string `json:"dep_name"`
	DepNo     string `json:"dep_no"`
}

func (CloudtimeDepartments) TableName() string {
	return "dt_dep"
}
