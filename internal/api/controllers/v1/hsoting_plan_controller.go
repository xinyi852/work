package v1

import (
	"plesk/internal/api/requests"
	"plesk/internal/pkg/common"
	"plesk/internal/pkg/response"
	"plesk/internal/services/plesk"

	"github.com/gin-gonic/gin"
)

type PleskProductController struct {
	BaseAPIController
}

func (ctrl *PleskProductController) Create(c *gin.Context) {
	//参数校验
	request := &common.CreatePlanRequest{}
	if ok := requests.Validate(c, request, requests.CreatePleskProductRequest); !ok {
		return
	}

	plesk_server, err := plesk.CreatePlan(request)

	if err != nil {
		response.JsonFail(c, response.ErrDefault, err.Error())
		return
	}

	response.JsonSuccess(c, gin.H{
		"plan_id": plesk_server.ID,
	})
	return

}

func (ctrl *PleskProductController) Delete(c *gin.Context) {

	//参数校验
	//参数校验
	request := &common.DeletePlanRequest{}
	if ok := requests.Validate(c, request, requests.DeletePleskProductRequest); !ok {
		return
	}
	id := request.Id
	//检查账户是否已经存在
	is_exist_product := plesk.ExistPlan("id", id)

	//不存在服务
	if !is_exist_product {
		response.JsonFail(c, 1, "该订阅类型不存在，请检查。")
		return
	}

	plesk.DelPlan(id)

	response.JsonSuccess(c, gin.H{
		"plan_id": id,
	})

	return
}

func (ctrl *PleskProductController) Update(c *gin.Context) {

	//参数校验
	request := &common.UpdatePlanRequest{}
	if ok := requests.Validate(c, request, requests.UpdatePleskProductRequest); !ok {
		return
	}
	id := request.Id
	//检查账户是否已经存在
	is_exist_product := plesk.ExistPlan("id", id)

	if !is_exist_product {
		response.JsonFail(c, 1, "该订阅类型不存在，请检查。")
		return
	}

	plesk_product := plesk.GetPlan(id)
	plesk_product.Name = request.Name
	plesk_product.Plan = request.Plan
	plesk_product.MaxData = request.MaxData
	plesk_product.MaxTraffic = request.MaxTraffic
	plesk_product.Save()

	response.JsonSuccess(c, gin.H{
		"plan_id": id,
	})

	return
}

func (ctrl *PleskProductController) List(c *gin.Context) {
	pleskProducts, pager := plesk.ProductPaginate(c, 10)

	response.JsonSuccess(c, gin.H{
		"list":  pleskProducts,
		"pager": pager,
	})
}
