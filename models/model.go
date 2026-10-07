package models

type Human struct {
	Sex              string  `json:"sex"`
	Age              int     `json:"age"`
	IPK              float32 `json:"ipk"`
	SimpleEmployment int     `json:"simpleemployment"`
	IsInvalid        bool    `json:"isinvalid"`
	InvalidGroup     int     `json:"invalidgroup"`
	OperatorID       int     `json:"operatorid"`
}

// type SimpleEmployment struct {
// 	Years     int `json:"years"`
// 	Months    int `json:"months"`
// 	Days      int `json:"days"`
// 	TotalDays int `json:"totaldays"`
// }
