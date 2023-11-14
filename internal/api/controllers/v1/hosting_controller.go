package v1

import (
	"encoding/json"
	"fmt"
	"plesk/internal/api/requests"
	"plesk/internal/models"
	"plesk/internal/pkg/common"
	"plesk/internal/pkg/response"
	"plesk/internal/services/plesk"
	"strings"

	"github.com/gin-gonic/gin"
	"racent.com/pkg/logger"
)

type HostingController struct {
	BaseAPIController
}

func (ctrl *HostingController) Create(c *gin.Context) {
	//参数校验
	request := &common.CreateHostingRequest{}
	if ok := requests.Validate(c, request, requests.CreateHosting); !ok {
		return
	}
	domain := request.Domain
	plan := request.Plan
	area := request.Area
	module := request.Module

	//根据参入参数选择对应的模块,目前仅支持plesk
	if module != "plesk" {
		response.JsonFail(c, 1, "目前暂不支持 "+module+" 类型的虚拟主机")
		return
	}

	//校验域名 是否已经存在,已经购买过则报错
	exist_domain := plesk.ExistService("domain", domain)
	if exist_domain {
		response.JsonFail(c, 1, "该域名已经购买过服务，域名:"+domain)
		return
	}

	//校验plan
	exist_plan := plesk.ExistPlan("plan", plan)
	if !exist_plan {
		response.JsonFail(c, 1, "plan不存在:"+plan)
		return
	}

	//查询主机组
	server_group := plesk.GetPleskServerGroup("area", area)
	// fmt.Print(server_group)
	//获取主机配置，此处可以进行负载均衡操作
	server := plesk.SelectServer(server_group)
	if server.Account == "" {
		response.JsonFail(c, 1, "未找到符合区域的主机:"+area)
		return
	}

	//准备请求参数
	login := request.Login
	//若未传入login 参数，则默认设置
	if login == "" {
		//去除域名中的 .
		str := strings.ReplaceAll(domain, ".", "")
		// fmt.Print(str)
		// login := str[0:12]
		login = str //取前16位作为账号
	}

	//检查账户是否已经存在
	is_exist_account := plesk.ExistService("account", login)

	//todo账户已经存在，则无法创建
	if is_exist_account {
		response.JsonFail(c, 1, "该登录账户："+login+" 已经存在，请重新配置该参数。")
		return
	}

	//检查email参数，配置默认值
	email := request.Email
	if email == "" {
		email = login + "@email.com"
	}

	//账户默认语言
	locale := "zh-CN"

	//公司
	company := request.Company
	if company == "" {
		company = "plesk"
	}

	//描述
	description := request.Description
	if description == "" {
		description = "description"
	}

	//随机生成密码
	s_pwd := plesk.GetPasswd(8, "advance") //todo 随机生成8位密码

	//唯一id。未提供则自动生成
	ex_id := request.ExternalId
	if ex_id == "" {
		ex_id = plesk.GenUUid()
	}
	// ex_id := plesk.GenUUid() // todo 随机生成的唯一id

	//开始创建账号
	plesk_account := &models.PleskAccount{
		Name:        domain,
		Company:     company,
		Login:       login,
		Status:      "0", //状态 0为启用，16为挂起
		Email:       email,
		Locale:      locale,
		OwnerLogin:  "admin", //由管理员提供服务
		ExternalId:  ex_id,
		Description: description,
		Password:    s_pwd,
		AccountType: "customer", //账号权限为customer
	}
	// fmt.Print(plesk_account)

	account_res := plesk.CreateAccount(server, plesk_account)

	//创建账户成功，开始创建订阅站点
	if account_res.StatusCode == 201 {
		abody := account_res.Body
		data := make(map[string]interface {
		})
		err := json.Unmarshal(abody, &data)

		if err != nil {
			panic(err)
		}

		aid := fmt.Sprintf("%v", data["id"])
		aguid := fmt.Sprintf("%v", data["guid"])

		//获取产品配置
		f_account := "ftp_" + login
		f_pwd := plesk.GetPasswd(8, "advance") //todo 随机生成8位密码

		plsek_service := &models.PleskService{
			ServerId:    fmt.Sprintf("%v", server.ID), //主机id
			Plan:        request.Plan,                 //产品订阅类型
			Domain:      request.Domain,               //站点名称
			Account:     login,                        // 账户名，默认为域名去掉.的前6位，若已经存在则添加数字1，以此类推
			Password:    s_pwd,                        //密码，随机生成的8位密码
			Status:      "0",                          //账户状态，默认为0
			Aid:         aid,                          //account id
			Did:         "0",                          //domain ID ,账户创建成功后，添加的domain ID
			Aguid:       aguid,                        // 账户的 的 guid
			Dguid:       "0",                          // 站点的guid
			FtpAccount:  f_account,                    // ftp账户名，默认为域名去掉.的前6位，若已经存在则添加数字1，以此类推
			FtpPassword: f_pwd,                        //ftp密码，随机生成的8位密码
		}

		plsek_service.Create()

		//开始创建站点订阅

		plesk_domain := &models.PleskDomain{
			Name:        request.Domain,
			Description: description,
			HostingType: "virtual",
			HostingSettings: map[string]string{
				"ftp_login":    f_account,
				"ftp_password": f_pwd,
			},
			OwnerClient: map[string]string{
				"id":    aid,
				"login": login,
				"guid":  aguid,
			},
			IpAddresses: map[string]string{},
			Ipv4:        map[string]string{},
			Ipv6:        map[string]string{},
			Plan: map[string]string{
				"name": plan, //Professional,Unlimited，Business
			},
		}

		// todo 这里修改为异步操作
		service_id := fmt.Sprintf("%v", plsek_service.ID)
		//异步创建domain
		go func() {
			CreateDomain(server, plesk_domain, service_id)
		}()

		response.JsonSuccess(c, gin.H{
			"id": service_id,
		})

		// 以下注释为同步创建订阅，处理速度较慢
		// domain_res := plesk.CreateDomain(server, plesk_domain)
		// if domain_res.StatusCode == 201 {
		// 	//创建成功
		// 	dbody := domain_res.Body
		// 	data1 := make(map[string]interface {
		// 	})
		// 	err := json.Unmarshal(dbody, &data1)

		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	// fmt.Print(data1)

		// 	did := fmt.Sprintf("%v", data1["id"])
		// 	dguid := fmt.Sprintf("%v", data1["guid"])
		// 	//保存信息
		// 	plsek_service := &models.PleskService{
		// 		ServerId:    fmt.Sprintf("%v", server.ID), //主机id
		// 		ProductId:   request.ProductID,            //产品id
		// 		Domain:      request.Domain,               //站点名称
		// 		Account:     login,                        // 账户名，默认为域名去掉.的前6位，若已经存在则添加数字1，以此类推
		// 		Password:    s_pwd,                        //密码，随机生成的8位密码
		// 		Status:      "0",                          //账户状态，默认为0
		// 		Aid:         aid,                          //account id
		// 		Did:         did,                          //domain ID ,账户创建成功后，添加的domain ID
		// 		Aguid:       aguid,                        // 账户的 的 guid
		// 		Dguid:       dguid,                        // 站点的guid
		// 		FtpAccount:  login,                        // ftp账户名，默认为域名去掉.的前6位，若已经存在则添加数字1，以此类推
		// 		FtpPassword: s_pwd,                        //ftp密码，随机生成的8位密码
		// 	}
		// 	// fmt.Print(plsek_service)
		// 	plsek_service.Create()

		// 	logger.InfoJSON("plsek_service", "default", "plsek_service create success")

		// 	// fmt.Print(res.RowsAffected)
		// 	// fmt.Print(res.Error)
		// 	//返回内容
		// 	// data, _ := json.Marshal(plsek_service)
		// 	response.JsonSuccess(c, gin.H{
		// 		"id": plsek_service.ID,
		// 	})
		// } else {
		// 	//todo 是否删除已经创建的账户
		// 	response.JsonFail(c, code.AdminCreateError, "订阅创建失败："+string(domain_res.Body))
		// }
	} else {
		response.JsonFail(c, 1, "账户创建失败："+string(account_res.Body))
	}

}

func CreateDomain(server models.PleskServer, plesk_domain *models.PleskDomain, service_id string) {

	domain_res := plesk.CreateDomain(server, plesk_domain)
	if domain_res.StatusCode == 201 {
		//创建成功
		dbody := domain_res.Body
		data1 := make(map[string]interface {
		})
		err := json.Unmarshal(dbody, &data1)

		if err != nil {
			panic(err)
		}

		did := fmt.Sprintf("%v", data1["id"])
		dguid := fmt.Sprintf("%v", data1["guid"])
		plsek_service := plesk.GetServiceByID(service_id)

		plsek_service.Did = did
		plsek_service.Dguid = dguid

		//更新service信息
		plsek_service.Save()

		logger.InfoJSON("plsek_service", "default", "domain create success service id:"+service_id)

	} else {
		logger.ErrorJSON("plsek_service", "error", "domain create fail service id:"+service_id)
	}

}

func (ctrl *HostingController) Suspend(c *gin.Context) {

	//参数校验
	request := &common.SetPleskRequest{}
	if ok := requests.Validate(c, request, requests.SetPlesk); !ok {
		return
	}
	id := request.ID

	//检查账户是否已经存在
	is_exist_service := plesk.ExistService("id", id)

	//不存在服务
	if !is_exist_service {
		response.JsonFail(c, 1, "该服务不存在，请检查。")
		return
	}

	plesk_service := plesk.GetServiceByID(id)

	if plesk_service.Status == "16" {
		response.JsonFail(c, 1, "服务已经是挂起状态！")
		return
	}

	server := plesk.GetServer(plesk_service.ServerId)

	//请求suspend接口
	suspend_res := plesk.Suspend(server, plesk_service.Aid)

	if suspend_res.StatusCode != 200 {

		response.JsonFail(c, 1, "暂停失败，接口请求错误："+string(suspend_res.Body))
		return
	}

	//修改数据库信息
	plesk_service.Status = "16"
	plesk_service.Save()

	response.JsonSuccess(c, gin.H{
		"id": id,
	})
	return

}

func (ctrl *HostingController) Active(c *gin.Context) {

	//参数校验
	request := &common.SetPleskRequest{}
	if ok := requests.Validate(c, request, requests.SetPlesk); !ok {
		return
	}
	id := request.ID

	//检查账户是否已经存在
	is_exist_service := plesk.ExistService("id", id)

	//不存在服务
	if !is_exist_service {
		response.JsonFail(c, 1, "该服务不存在，请检查。")
		return
	}

	plesk_service := plesk.GetServiceByID(id)

	if plesk_service.Status == "0" {
		response.JsonFail(c, 1, "服务已经是启用状态！")
		return
	}
	server := plesk.GetServer(plesk_service.ServerId)

	//请求suspend接口
	suspend_res := plesk.Active(server, plesk_service.Aid)

	if suspend_res.StatusCode != 200 {

		response.JsonFail(c, 1, "启用失败，接口请求错误："+string(suspend_res.Body))
		return
	}

	//修改数据库信息
	plesk_service.Status = "0"
	plesk_service.Save()

	response.JsonSuccess(c, "启用成功")
	return

}

func (ctrl *HostingController) Upgrade(c *gin.Context) {

	//参数校验
	request := &common.UpgradePleskRequest{}
	if ok := requests.Validate(c, request, requests.UpgradePlesk); !ok {
		return
	}

	service_id := request.ID

	//检查账户是否已经存在
	is_exist_service := plesk.ExistService("id", service_id)

	//todo账户已经存在，则添加一个数字
	if !is_exist_service {
		response.JsonFail(c, 1, "该服务不存在，请检查。")
		return
	}

	plesk_service := plesk.GetServiceByID(service_id)
	if plesk_service.Status != "0" {
		response.JsonFail(c, 1, "服务已挂起！")
		return
	}
	// fmt.Print(plesk_service.ProductId)
	// fmt.Print(request.ProductID)
	if plesk_service.Plan == request.Plan {
		response.JsonFail(c, 1, "服务已是该套餐，请检查。")
		return
	}
	server := plesk.GetServer(plesk_service.ServerId)

	// plan_data := &common.Plan{
	// 	Name: plan,
	// }
	// upgrade_domain := &common.UpgradePlanRequest{
	// 	Name: plesk_service.Domain,
	// 	Plan: *plan_data,
	// }

	//请求接口
	// fmt.Print(product.Config)

	guid := plesk.GetPlanGuid(server, request.Plan)

	// fmt.Print(guid)

	// fmt.Print(guid)
	// upgrade :=
	upgrade_res := plesk.UpgradeService(server, plesk_service.Domain, guid)

	if upgrade_res == "ok" {
		//修改数据库信息
		plesk_service.Plan = request.Plan
		plesk_service.Save()
		// fmt.Print(plk_res)

		response.JsonSuccess(c, gin.H{
			"id": service_id,
		})
		return
	} else {
		response.JsonFail(c, 1, "升级失败，接口请求错误")
		return
	}
}

func (ctrl *HostingController) ChangePassword(c *gin.Context) {

	//参数校验
	request := &common.ChangePasswordRequest{}
	if ok := requests.Validate(c, request, requests.ChangePasswordPlesk); !ok {
		return
	}

	//检查账户是否已经存在
	is_exist_service := plesk.ExistService("id", request.ID)

	//不存在服务
	if !is_exist_service {
		response.JsonFail(c, 1, "该服务不存在，请检查。")
		return
	}

	plesk_service := plesk.GetServiceByID(request.ID)

	if plesk_service.Status != "0" {
		response.JsonFail(c, 1, "服务已挂起！")
		return
	}

	//比对旧密码？

	// if plesk_service.Password != request.Password {
	// 	response.JsonFail(c, 1, "密码错误")
	// 	return
	// }

	server := plesk.GetServer(plesk_service.ServerId)

	change_data := &common.ChangePassword{
		// Login:    request.Account,
		Password: request.Password,
	}

	//请求接口修改密码

	change_res := plesk.ChangePassword(server, change_data, plesk_service.Aid)

	if change_res.StatusCode == 200 {
		//修改数据库信息
		// plesk_service.Account = request.Account
		plesk_service.Password = request.Password
		plesk_service.Save()

		response.JsonSuccess(c, "修改成功")
		return
	} else {
		response.JsonFail(c, 1, "修改失败，接口请求错误")
		return
	}
}

func (ctrl *HostingController) Delete(c *gin.Context) {

	//参数校验
	request := &common.SetPleskRequest{}
	if ok := requests.Validate(c, request, requests.SetPlesk); !ok {
		return
	}
	id := request.ID
	//检查账户是否已经存在
	is_exist_service := plesk.ExistService("id", id)

	//不存在服务
	if !is_exist_service {
		response.JsonFail(c, 1, "该服务不存在，请检查。")
		return
	}
	plesk_service := plesk.GetServiceByID(id)

	server := plesk.GetServer(plesk_service.ServerId)

	// fmt.Print("开始请求接口，aid:" + plesk_service.Aid)
	//请求接口删除plesk
	change_res := plesk.Delete(server, plesk_service.Aid)

	if change_res.StatusCode == 200 {
		//修改数据库信息

		plesk_service.Delete()

		// fmt.Print(plk_res)

		response.JsonSuccess(c, "删除成功")
		return
	} else {
		response.JsonFail(c, 1, "删除失败，接口请求错误："+string(change_res.Body))
		return
	}
}

func (ctrl *HostingController) Info(c *gin.Context) {

	//参数校验
	request := &common.SetPleskRequest{}
	if ok := requests.Validate(c, request, requests.SetPlesk); !ok {
		return
	}
	id := request.ID
	//检查账户是否已经存在
	is_exist_service := plesk.ExistService("id", id)

	//不存在服务
	if !is_exist_service {
		response.JsonFail(c, 1, "该服务不存在，请检查。")
		return
	}
	plesk_service := plesk.GetServiceByID(id)

	server := plesk.GetServer(plesk_service.ServerId)

	// fmt.Print("开始请求接口，aid:" + plesk_service.Aid)
	//需要获取服务用量使用情况
	res := plesk.GetUseage(server, plesk_service)

	if res["traffic"] == "err" {
		response.JsonFail(c, 1, "获取失败")
	}

	res_data := map[string]string{
		"traffic":        res["traffic"],
		"disk_usage":     res["disk_usage"],
		"max_disk_usage": res["max_disk_usage"],
		"max_traffic":    res["max_traffic"],
		"user_name":      plesk_service.Account,
		"server_name":    server.Url,
		"server_ip":      server.Ip,
	}

	response.JsonSuccess(c, res_data)
	return

}

func (ctrl *HostingController) SignLogin(c *gin.Context) {

	//参数校验
	request := &common.SetPleskRequest{}
	if ok := requests.Validate(c, request, requests.SetPlesk); !ok {
		return
	}

	plesk_service := plesk.GetServiceByID(request.ID)

	if plesk_service.Account == "" {
		response.JsonFail(c, 1, "未找到该服务，请检查。")
		return
	}

	if plesk_service.Status != "0" {
		response.JsonFail(c, 1, "服务已挂起！")
		return
	}

	plesk_server := plesk.GetServer(plesk_service.ServerId)

	token := plesk.GetSignToken(plesk_server, plesk_service)

	if token == "1" {
		response.JsonFail(c, 1, "申请token失败")
	}

	url := "https://" + plesk_server.Url + ":8443/enterprise/rsession_init.php?PLESKSESSID="

	response.JsonSuccess(c, gin.H{
		"url": url + token,
	})

	return

}
