package models

import (
	"racent.com/pkg/database"
)

var ()

type PleskServer struct {
	BaseModel
	Account  string `json:"account"`  // 主机管理员账号
	Password string `json:"password"` // 主机管理员密码
	Status   uint64 `json:"status"`   // 主机状态，默认为0
	Ip       string `json:"ip"`       // 主机ip
	Max      string `json:"max"`      // 主机最大服务
	Service  string `json:"service"`  // 主机服务数量
	Fill     string `json:"fill"`     // 主机负载 = 服务数量/最大服务 *100
	Url      string `json:"url"`      // 主机管理面板地址
	// Port     string `json:"port"`            // 主机管理面板端口
	ServerGroupId string `json:"server_group_id"` // 主机所属区域id
	Type          string `json:"type"`            // 主机所属区域id
}

func (PleskServer *PleskServer) TableName() string {
	return "hostings"
}

func (PleskServer *PleskServer) Create() {
	database.DB.Create(&PleskServer)
}

func (PleskServer *PleskServer) Save() (rowsAffected int64) {
	result := database.DB.Save(&PleskServer)
	return result.RowsAffected
}

func (PleskServer *PleskServer) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&PleskServer)
	return result.RowsAffected
}
