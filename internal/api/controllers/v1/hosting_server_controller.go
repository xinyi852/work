package v1

import (
	"plesk/internal/api/requests"
	"plesk/internal/pkg/common"
	"plesk/internal/pkg/response"
	"plesk/internal/services/plesk"

	"github.com/gin-gonic/gin"
)

type PleskServerController struct {
	BaseAPIController
}

func (ctrl *PleskServerController) Create(c *gin.Context) {
	//参数校验
	request := &common.CreatePleskServerRequest{}
	if ok := requests.Validate(c, request, requests.CreatePleskServer); !ok {
		return
	}

	//todo::检查该plesk主机是否已经存在，通过ip判断
	is_exist := plesk.ExistServer("ip", request.Ip)
	if is_exist {
		response.JsonFail(c, 1, "已经添加过该主机，ip:"+request.Ip)
		return
	}
	//测试连接
	// res := plesk.TestServer(request)

	res := 200
	if res == 200 {
		//保存数据到数据库

		// 修改数据库
		plesk_server, err := plesk.CreateServer(request)
		if err != nil {
			response.JsonFail(c, response.ErrDefault, err.Error())
			return
		}

		response.JsonSuccess(c, plesk_server)
		//返回成功

	} else {
		response.JsonFail(c, 1, "连接测试失败，添加主机未成功，请检查")
	}

}

func (ctrl *PleskServerController) Update(c *gin.Context) {
	//参数校验
	request := &common.UpdatePleskServerRequest{}
	if ok := requests.Validate(c, request, requests.UpdatePleskServer); !ok {
		return
	}

	//todo::检查该plesk主机是否已经存在，通过ip判断
	is_exist := plesk.ExistServer("id", request.Id)
	if !is_exist {
		response.JsonFail(c, 1, "该主机不存在，ip:"+request.Ip)
		return
	}
	//测试连接
	// res := plesk.TestServer(request)

	res := 200
	if res == 200 {
		//保存数据到数据库

		// 修改数据库
		res := plesk.UpdateServer(request)
		if !res {
			response.JsonFail(c, 1, "更新主机组失败")
			return
		}

		response.JsonSuccess(c, gin.H{
			"server_id": request.Id,
		})
		return
		//返回成功

	} else {
		response.JsonFail(c, 1, "连接测试失败，添加主机未成功，请检查")
	}

}

// 删除plesk主机
func (ctrl *PleskServerController) Delete(c *gin.Context) {

	//参数校验
	request := &common.DeleteServerRequest{}
	if ok := requests.Validate(c, request, requests.DeleteServerRequest); !ok {
		return
	}
	server_id := request.Id

	// fmt.Print(server_id)

	pleskServer := plesk.GetServer(server_id)
	if pleskServer.ID == 0 {
		response.Abort404(c)
		return
	}

	rowsAffected := pleskServer.Delete()
	if rowsAffected > 0 {
		response.JsonSuccess(c, gin.H{
			"server_id": server_id,
		})
		return
	}
	response.JsonFail(c, response.ErrDefault)
}

// 主机列表
func (ctrl *PleskServerController) ServerList(c *gin.Context) {
	pleskServers, pager := plesk.ServerPaginate(c, 10)

	response.JsonSuccess(c, gin.H{
		"list":  pleskServers,
		"pager": pager,
	})
}
