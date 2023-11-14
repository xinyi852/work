package models

import (
	"racent.com/pkg/database"
)

var ()

type PleskServerGroup struct {
	BaseModel
	Name     string `json:"name"`      // 区域名称
	FillType string `json:"fill_type"` // 策略
	Area     string `json:"area"`      //区域
	// Server   string `json:"server"`    // 包含的主机
}

func (PleskServerGroup *PleskServerGroup) TableName() string {
	return "hosting_groups"
}

func (PleskServerGroup *PleskServerGroup) Create() {
	database.DB.Create(&PleskServerGroup)
}

func (PleskServerGroup *PleskServerGroup) Save() (rowsAffected int64) {
	result := database.DB.Save(&PleskServerGroup)
	return result.RowsAffected
}

func (PleskServerGroup *PleskServerGroup) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&PleskServerGroup)
	return result.RowsAffected
}
