package models

import (
	"gorm.io/gorm"
	"racent.com/pkg/database"
)

var ()

type PleskService struct {
	BaseModel
	ServerId    string `json:"server_id"`    //主机id
	Plan        string `json:"plan"`         //plan
	Domain      string `json:"domain"`       //站点名称
	Account     string `json:"account"`      // 账户名，默认为域名去掉.的前6位，若已经存在则添加数字1，以此类推
	Password    string `json:"password"`     //密码，随机生成的8位密码
	Status      string `json:"status"`       //账户状态，默认为0
	Aid         string `json:"aid"`          //账户id
	Aguid       string `json:"aguid"`        // 账户的 的 guid
	Did         string `json:"did"`          //domain ID ,账户创建成功后，添加的domain ID
	Dguid       string `json:"dguid"`        // 站点的guid
	FtpAccount  string `json:"ftp_account"`  // ftp账户名，默认为域名去掉.的前6位，若已经存在则添加数字1，以此类推
	FtpPassword string `json:"ftp_password"` //ftp密码，随机生成的8位密码
}

type PleskAccount struct {
	BaseModel
	Name        string `json:"name"`
	Company     string `json:"company"`
	Login       string `json:"login"`
	Status      string `json:"status"`
	Email       string `json:"email"`
	Locale      string `json:"locale"`
	OwnerLogin  string `json:"owner_login"`
	ExternalId  string `json:"external_id"`
	Description string `json:"description"`
	Password    string `json:"password"`
	AccountType string `json:"type"`
}

type PleskDomain struct {
	BaseModel
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	HostingType     string            `json:"hosting_type"`
	HostingSettings map[string]string `json:"hosting_settings"`
	// BaseDomain      string `json:"base_domain"`
	// ParentDomain    string `json:"parent_domain"`
	OwnerClient map[string]string `json:"owner_client"`
	IpAddresses map[string]string `json:"ip_addresses"`
	Ipv4        map[string]string `json:"ipv4"`
	Ipv6        map[string]string `json:"ipv6"`
	Plan        map[string]string `json:"plan"`
}

func (PleskService *PleskService) TableName() string {
	return "hosting_services"
}

func (PleskService *PleskService) Create() (rowsAffected *gorm.DB) {
	res := database.DB.Create(&PleskService)
	return res
}

func (PleskService *PleskService) Exist() {
	database.DB.Where("domain = ?", PleskService.Domain).First(&PleskService)
}

func (PleskService *PleskService) Save() (rowsAffected int64) {
	result := database.DB.Save(&PleskService)
	return result.RowsAffected
}

func (PleskService *PleskService) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&PleskService)
	return result.RowsAffected
}
