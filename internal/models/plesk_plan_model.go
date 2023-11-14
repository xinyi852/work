package models

import (
	"racent.com/pkg/database"
)

var ()

type PleskPlan struct {
	BaseModel
	Name       string `json:"name"` // 产品名称
	Plan       string `json:"plan"` //配置参数
	MaxData    string `json:"max_data"`
	MaxTraffic string `json:"max_traffic"`
}

func (PleskPlan *PleskPlan) TableName() string {
	return "hosting_plans"
}

func (PleskPlan *PleskPlan) Create() {
	database.DB.Create(&PleskPlan)
}

func (PleskPlan *PleskPlan) Save() (rowsAffected int64) {
	result := database.DB.Save(&PleskPlan)
	return result.RowsAffected
}

func (PleskPlan *PleskPlan) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&PleskPlan)
	return result.RowsAffected
}
