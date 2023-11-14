package v1

import (
	"plesk/internal/api/requests"
	"plesk/internal/pkg/common"
	"plesk/internal/pkg/response"
	"plesk/internal/services/plesk"

	"github.com/gin-gonic/gin"
)

type PleskServerGroupController struct {
	BaseAPIController
}

// 创建新的server group
func (ctrl *PleskServerGroupController) Create(c *gin.Context) {
	//参数校验 data
	request := &common.CreatePleskServerGroupData{}

	if ok := requests.Validate(c, request, requests.CreatePleskServerGroupData); !ok {
		return
	}

	plesk_server, err := plesk.CreatePleskServerGroup(request)

	if err != nil {
		response.JsonFail(c, response.ErrDefault, err.Error())
		return
	}

	response.JsonSuccess(c, gin.H{
		"server_id": plesk_server.ID,
	})
	return

}

func (ctrl *PleskServerGroupController) Delete(c *gin.Context) {

	//参数校验
	request := &common.DeletePleskServerGroupData{}

	if ok := requests.Validate(c, request, requests.DeletePleskServerGroupData); !ok {
		return
	}

	group_id := request.GroupId
	//检查组是否已经存在
	exist_group := plesk.ExistPleskServerGroup("id", group_id)

	//不存在
	if !exist_group {
		response.JsonFail(c, 1, "该主机组不存在，请检查。")
		return
	}

	//判断主机组下是否存在主机
	exist_server := plesk.ExistServer("group_id", group_id)
	//存在主机
	if exist_server {
		response.JsonFail(c, 1, "主机组存在未删除的主机，请先删除主机。")
		return
	}

	res := plesk.DeletePleskServerGroup(group_id)

	if !res {
		response.JsonFail(c, 1, "删除主机组失败！")
		return
	}

	response.JsonSuccess(c, gin.H{
		"server_id": group_id,
	})
	return

}

func (ctrl *PleskServerGroupController) Update(c *gin.Context) {
	//参数校验
	request := &common.UpdatePleskServerGroupData{}
	if ok := requests.Validate(c, request, requests.UpdatePleskServerGroupData); !ok {
		return
	}
	group_id := request.GroupId
	//检查服务器区组是否已经存在
	exist_group := plesk.ExistPleskServerGroup("id", group_id)

	if !exist_group {
		response.JsonFail(c, 1, "该主机组不存在，请检查。")
		return
	}

	res := plesk.UpdatePleskServerGroup(request)

	if !res {
		response.JsonFail(c, 1, "更新主机组失败")
		return
	}

	response.JsonSuccess(c, gin.H{
		"server_id": group_id,
	})
	return
}

func (ctrl *PleskServerGroupController) List(c *gin.Context) {

	request := &common.ServerGropuListRequest{}
	if ok := requests.Validate(c, request, requests.PleskServerGroupList); !ok {
		return
	}

	pleskServerGroups, pager := plesk.PleskServerGroupPaginate(c, request)

	response.JsonSuccess(c, gin.H{
		"list":  pleskServerGroups,
		"pager": pager,
	})
}
