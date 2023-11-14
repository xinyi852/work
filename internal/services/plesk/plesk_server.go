package plesk

import (
	"encoding/base64"
	"plesk/internal/models"
	"plesk/internal/pkg/common"

	"errors"
	"plesk/internal/pkg/httpclient"

	"plesk/internal/pkg/paginator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"racent.com/pkg/database"
	"racent.com/pkg/logger"
)

func GetServer(idStr string) (PleskServer models.PleskServer) {
	database.DB.Where("id", idStr).First(&PleskServer)
	return PleskServer
}

func CreateServer(request *common.CreatePleskServerRequest) (*models.PleskServer, error) {
	logger.DebugJSON("CreateServer", "request", request)
	pleskServerModel := &models.PleskServer{
		Account:       request.Account,
		Password:      request.Password,
		Status:        0,
		Ip:            request.Ip,
		Max:           request.Max,
		Service:       request.Service,
		Url:           request.Url,
		Fill:          "0", //todo
		ServerGroupId: request.AreaGroup,
		Type:          request.Type,
	}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		tx.Create(pleskServerModel)
		if pleskServerModel.ID == 0 {
			logger.ErrorJSON("CreateServer", "request", request)
			return errors.New("保存主机失败")
		}
		return nil
	})
	return pleskServerModel, err
}

func GetServerByID(id string) models.PleskServer {
	var plesk_serrver models.PleskServer
	database.DB.Model(models.PleskServer{}).Where("id = ?", id).First(&plesk_serrver)
	return plesk_serrver
}

func UpdateServer(request *common.UpdatePleskServerRequest) bool {
	logger.DebugJSON("update server", "request", request)

	server := GetServerByID(request.Id)
	server.Account = request.Account
	server.Password = request.Password

	server.Ip = request.Ip
	server.Max = request.Max
	server.Service = request.Service
	server.Url = request.Url
	server.ServerGroupId = request.AreaGroup
	server.Type = request.Type

	res := server.Save()
	return res == 1

}

func TestServer(request *common.CreatePleskServerRequest) int {
	//获取ip，测试接口是否能正常访问
	//绑定json和结构体
	//获取json中的key,注意使用 . 访问
	account := request.Account
	pwd := request.Password
	url := request.Url

	str := account + ":" + pwd

	md5str := base64.StdEncoding.EncodeToString([]byte(str))
	auth := "Basic " + md5str

	//测试接口
	api_url := "https://" + url + ":8443/api/v2/cli/commands"

	http := httpclient.RequestData{
		Headers: map[string]string{
			"Authorization": auth,
			"Content-Type":  "application/json",
			"Accept":        "application/json",
		},
		Url:    api_url,
		Params: map[string]string{},
		//Params     map[string]interface{}
		Cert:       "",
		PrivateKey: "",
		Files:      nil,
	}

	res := http.Get().StatusCode

	return res

}

func DelServer(serverID string) (PleskServer models.PleskServer) {
	database.DB.Model(models.PleskServer{}).Where("id", serverID).Delete(&PleskServer)
	return
}

func ExistServer(field, value string) bool {
	var count int64
	database.DB.Model(models.PleskServer{}).Where(field+" = ?", value).Count(&count)
	return count > 0
}

func SelectServer(server_group models.PleskServerGroup) (PleskServer models.PleskServer) {

	//负载均衡
	if server_group.FillType == "bl" {

		database.DB.Model(models.PleskServer{}).Where("server_group_id=?", server_group.ID).Where("fill < 100").Order("fill DESC").First(&PleskServer)

	} else {
		//单个满载
		database.DB.Model(models.PleskServer{}).Where("server_group_id=?", server_group.ID).Where("fill < 100").Order("order ASC").First(&PleskServer)

	}

	return PleskServer
}

func ServerPaginate(c *gin.Context, perPage int) (PleskServers []models.PleskServer, paging paginator.Paging) {
	queryDB := database.DB.Model(models.PleskServer{})
	// keyword := c.DefaultQuery("keyword", "")
	// if keyword != "" {
	//     queryDB.Where("name LIKE ?", "%"+keyword+"%")
	// }
	paging = paginator.Paginate(
		c,
		queryDB,
		&PleskServers,
		perPage,
	)

	//隐去敏感信息
	for i, _ := range PleskServers {
		PleskServers[i].Password = "******"
	}

	return
}
