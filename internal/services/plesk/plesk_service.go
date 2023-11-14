package plesk

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"plesk/internal/models"
	"plesk/internal/pkg/common"
	"plesk/internal/pkg/httpclient"

	"racent.com/pkg/database"
)

func ExistService(field, value string) bool {
	var count int64
	database.DB.Model(models.PleskService{}).Where(field+" = ?", value).Count(&count)
	return count > 0
}

func CreateAccount(server models.PleskServer, plesk_account *models.PleskAccount) httpclient.ResponseWrapper {

	auth := auth(server)
	//创建账户接口
	api_url := "https://" + server.Url + ":8443/api/v2/clients"

	data, _ := json.Marshal(&plesk_account)

	http := httpclient.RequestData{
		Headers: map[string]string{
			"Authorization": auth,
			"Content-Type":  "application/json",
			"Accept":        "application/json",
		},
		Url:        api_url,
		Params:     data,
		Cert:       "",
		PrivateKey: "",
		Files:      nil,
	}

	res := http.Post()
	return res
}

func CreateDomain(server models.PleskServer, plesk_domain *models.PleskDomain) httpclient.ResponseWrapper {

	auth := auth(server)
	//创建站点订阅接口
	api_url := "https://" + server.Url + ":8443/api/v2/domains"

	data, _ := json.Marshal(plesk_domain)
	// fmt.Print(string(data))

	http := httpclient.RequestData{
		Headers: map[string]string{
			"Authorization": auth,
			"Content-Type":  "application/json",
			"Accept":        "application/json",
		},
		Url:        api_url,
		Params:     data,
		Cert:       "",
		PrivateKey: "",
		Files:      nil,
	}
	// fmt.Print(string(data))
	res := http.Post()

	return res

}

func ChangePassword(server models.PleskServer, change_data *common.ChangePassword, aid string) httpclient.ResponseWrapper {

	auth := auth(server)
	//更新站点订阅接口
	api_url := "https://" + server.Url + ":8443/api/v2/clients/" + aid

	data, _ := json.Marshal(&change_data)
	// fmt.Print(string(data))

	http := httpclient.RequestData{
		Headers: map[string]string{
			"Authorization": auth,
			"Content-Type":  "application/json",
			"Accept":        "application/json",
		},
		Url:        api_url,
		Params:     data,
		Cert:       "",
		PrivateKey: "",
		Files:      nil,
	}
	// fmt.Print(string(data))
	res := http.Put()

	return res

}

func Delete(server models.PleskServer, id string) httpclient.ResponseWrapper {

	auth := auth(server)
	//删除账户接口
	api_url := "https://" + server.Url + ":8443/api/v2/clients/" + id

	http := httpclient.RequestData{
		Headers: map[string]string{
			"Authorization": auth,
			"Content-Type":  "application/json",
			"Accept":        "application/json",
		},
		Url:        api_url,
		Params:     "",
		Cert:       "",
		PrivateKey: "",
		Files:      nil,
	}
	res := http.Delete()

	return res

}

func Suspend(server models.PleskServer, id string) httpclient.ResponseWrapper {

	auth := auth(server)
	//删除账户接口
	api_url := "https://" + server.Url + ":8443/api/v2/clients/" + id + "/suspend"

	http := httpclient.RequestData{
		Headers: map[string]string{
			"Authorization": auth,
			"Content-Type":  "application/json",
			"Accept":        "application/json",
		},
		Url:        api_url,
		Params:     "",
		Cert:       "",
		PrivateKey: "",
		Files:      nil,
	}
	res := http.Put()

	return res

}

func Active(server models.PleskServer, id string) httpclient.ResponseWrapper {

	auth := auth(server)
	//删除账户接口
	api_url := "https://" + server.Url + ":8443/api/v2/clients/" + id + "/activate"

	http := httpclient.RequestData{
		Headers: map[string]string{
			"Authorization": auth,
			"Content-Type":  "application/json",
			"Accept":        "application/json",
		},
		Url:        api_url,
		Params:     "",
		Cert:       "",
		PrivateKey: "",
		Files:      nil,
	}
	res := http.Put()

	return res

}

func UpgradeDomain(server models.PleskServer, plan *common.UpgradePlanRequest, did string) httpclient.ResponseWrapper {

	auth := auth(server)
	//更新站点订阅接口
	api_url := "https://" + server.Url + ":8443/api/v2/domains/" + did

	data, _ := json.Marshal(&plan)
	// fmt.Print(string(data))

	http := httpclient.RequestData{
		Headers: map[string]string{
			"Authorization": auth,
			"Content-Type":  "application/json",
			"Accept":        "application/json",
		},
		Url:        api_url,
		Params:     data,
		Cert:       "",
		PrivateKey: "",
		Files:      nil,
	}
	// fmt.Print(string(data))
	res := http.Put()

	// fmt.Print(res.StatusCode)
	// fmt.Print(string(res.Body))

	return res

}

func auth(server models.PleskServer) string {
	account := server.Account
	pwd := server.Password

	str := account + ":" + pwd

	md5str := base64.StdEncoding.EncodeToString([]byte(str))
	auth := "Basic " + md5str

	return auth
}

func GetPasswd(length int, kind string) string {
	passwd := make([]rune, length)
	var codeModel []rune
	switch kind {
	case "num":
		codeModel = []rune("0123456789")
	case "char":
		codeModel = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	case "mix":
		codeModel = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	case "advance":
		codeModel = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+=-!@#$%*,.[]")
	default:
		codeModel = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	}
	for i := range passwd {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(codeModel))))
		passwd[i] = codeModel[int(index.Int64())]
	}
	return string(passwd)

}

func GenUUid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func SaveData(data *models.PleskService) {
	database.DB.Model(models.PleskService{}).Create(&data)
	return
}

func GetServiceByID(id string) models.PleskService {
	var plesk_serrvice models.PleskService
	database.DB.Model(models.PleskService{}).Where("id = ?", id).First(&plesk_serrvice)
	return plesk_serrvice
}

func GetUseage(plesk_server models.PleskServer, plesk_service models.PleskService) map[string]string {
	url := "https://" + plesk_server.Url + ":8443/enterprise/control/agent.php"

	account := plesk_service.Account
	password := plesk_service.Password

	buf := []byte(`<packet>
	<site>
		<get>
		   <filter>
				<name>` + plesk_service.Domain + `</name>
		   </filter>
		   <dataset>
				<stat/>
				<disk_usage/>
		   </dataset>
		</get>
	</site>
	</packet>`)
	data := bytes.NewBuffer(buf)

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return map[string]string{
			"traffic":    "err",
			"disk_usage": "err",
		}
	}

	// 添加自定义HTTP头
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("HTTP_AUTH_LOGIN", account)
	req.Header.Add("HTTP_AUTH_PASSWD", password)
	req.Header.Add("HTTP_PRETTY_PRINT", "TRUE")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return map[string]string{
			"traffic":    "err",
			"disk_usage": "err",
		}
	}
	defer resp.Body.Close()
	// 读取响应内容
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return map[string]string{
			"traffic":    "err",
			"disk_usage": "err",
		}
	}

	// XML字符串
	xmlStr := string(respData)

	// fmt.Print(xmlStr)
	// id := "1"
	//解析XML字符串
	var res common.UsagePacket
	xml.Unmarshal([]byte(xmlStr), &res)
	if err != nil {
		fmt.Printf("error: %v", err)
		return map[string]string{
			"traffic":    "err",
			"disk_usage": "err",
		}
	}
	traffic := res.UsageSite.UsageGet.Result.Data.UsageStat.Traffic
	disk_usage := res.UsageSite.UsageGet.Result.Data.UsageDisk.Httpdocs

	max_data := GetMaxByProductId(plesk_service.Plan)

	max_traffic := max_data["max_traffic"]
	max_disk_usage := max_data["max_disk_usage"]
	//打印解析结果
	// fmt.Printf("token: %+v", id)

	return map[string]string{
		"traffic":        traffic,
		"max_traffic":    max_traffic,
		"disk_usage":     disk_usage,
		"max_disk_usage": max_disk_usage,
	}

}

// 设置各产品最大硬盘及流量
func GetMaxByProductId(id string) map[string]string {
	switch id {
	case "1":
		return map[string]string{
			"max_traffic":    "31457280",
			"max_disk_usage": "2097152",
		}
	default:
		return map[string]string{
			"max_traffic":    "31457280",
			"max_disk_usage": "2097152",
		}
	}
}

func GetSignToken(plesk_server models.PleskServer, plesk_service models.PleskService) string {

	url := "https://" + plesk_server.Url + ":8443/enterprise/control/agent.php"
	// fmt.Print(url)

	//用户的ip
	ip := "123"
	encoded := base64.StdEncoding.EncodeToString([]byte(ip))
	account := plesk_service.Account
	password := plesk_service.Password

	buf := []byte(`<packet version="1.6.9.1">
	<server>
	  <create_session>
		<login>` + account + `</login>
		<data>
		  <user_ip>` + encoded + `</user_ip>
		  <source_server></source_server>
		</data>
	  </create_session>
	</server>
  </packet>`)
	data := bytes.NewBuffer(buf)

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return "1"
	}

	// 添加自定义HTTP头
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("HTTP_AUTH_LOGIN", account)
	req.Header.Add("HTTP_AUTH_PASSWD", password)
	req.Header.Add("HTTP_PRETTY_PRINT", "TRUE")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return "1"
	}
	defer resp.Body.Close()

	// 输出响应信息
	// fmt.Println("Response Status:", resp.Status)
	// fmt.Println("Response Headers:")
	// for key, value := range resp.Header {
	// 	fmt.Printf("%s: %s\n", key, value)
	// }

	// 读取响应内容
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return "1"
	}

	// 输出响应内容
	// fmt.Println("Response Body:")
	// fmt.Println(string(respData))

	// XML字符串
	xmlStr := string(respData)

	// 解析XML字符串
	var res common.Packet
	xml.Unmarshal([]byte(xmlStr), &res)
	if err != nil {
		fmt.Printf("error: %v", err)
		return "1"
	}
	id := res.Server.CreateSession.Result.ID
	// 打印解析结果
	// fmt.Printf("token: %+v", id)

	return id
}

func GetPlanGuid(plesk_server models.PleskServer, plan string) string {

	url := "https://" + plesk_server.Url + ":8443/enterprise/control/agent.php"

	account := plesk_server.Account
	password := plesk_server.Password

	buf := []byte(`<packet version="1.6.3.0">
	<service-plan>
	<get>
	   <filter>
	   		<name>` + plan + `</name>
	   </filter>
	</get>
	</service-plan>
	</packet>`)
	data := bytes.NewBuffer(buf)

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return "1"
	}

	// 添加自定义HTTP头
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("HTTP_AUTH_LOGIN", account)
	req.Header.Add("HTTP_AUTH_PASSWD", password)
	req.Header.Add("HTTP_PRETTY_PRINT", "TRUE")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return "1"
	}
	defer resp.Body.Close()

	// 输出响应信息
	// fmt.Println("Response Status:", resp.Status)
	// fmt.Println("Response Headers:")
	// for key, value := range resp.Header {
	// 	fmt.Printf("%s: %s\n", key, value)
	// }

	// 读取响应内容
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return "1"
	}

	// 输出响应内容
	// fmt.Println("Response Body:")
	// fmt.Println(string(respData))

	// XML字符串
	xmlStr := string(respData)

	// fmt.Print(xmlStr)

	// logger.InfoJSON("xmlinfo", "default", xmlStr)

	// 解析XML字符串
	var res common.PlanPacket
	xml.Unmarshal([]byte(xmlStr), &res)

	// fmt.Print(res.ServicePlan.Get.Result.GUID)
	// if err != nil {
	// 	fmt.Printf("error: %v", err)
	// 	return "1"
	// }
	// id := res.Server.CreateSession.Result.ID
	// 打印解析结果
	// fmt.Printf("token: %+v", id)

	return res.ServicePlan.Get.Result.GUID
}

func UpgradeService(plesk_server models.PleskServer, name string, guid string) string {

	url := "https://" + plesk_server.Url + ":8443/enterprise/control/agent.php"

	account := plesk_server.Account
	password := plesk_server.Password
	// id :="1"

	// buf := []byte(`<packet version ="1.6.3.0">
	// <webspace>
	//    <switch-subscription>
	// 	  <filter>
	// 		 <name>` + name + `</name>
	// 	  </filter>
	// 	   <plan-external-id>` + guid + `</plan-external-id>
	//    </switch-subscription>
	// </webspace>
	// </packet>`)
	buf := []byte(`<packet version ="1.6.3.0">
	<webspace>
	   <switch-subscription>
		  <filter>
			 <name>` + name + `</name>
		  </filter>
		   <plan-guid>` + guid + `</plan-guid>
	   </switch-subscription>
	</webspace>
	</packet>`)
	data := bytes.NewBuffer(buf)

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return "1"
	}

	// 添加自定义HTTP头
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("HTTP_AUTH_LOGIN", account)
	req.Header.Add("HTTP_AUTH_PASSWD", password)
	req.Header.Add("HTTP_PRETTY_PRINT", "TRUE")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return "1"
	}
	defer resp.Body.Close()

	// 输出响应信息
	// fmt.Println("Response Status:", resp.Status)
	// fmt.Println("Response Headers:")
	// for key, value := range resp.Header {
	// 	fmt.Printf("%s: %s\n", key, value)
	// }

	// 读取响应内容
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return "1"
	}

	// 输出响应内容
	// fmt.Println("Response Body:")
	// fmt.Println(string(respData))

	// XML字符串
	xmlStr := string(respData)

	fmt.Print(xmlStr)

	// logger.InfoJSON("xmlinfo", "default", xmlStr)

	// 解析XML字符串
	var res common.UpgradePacket
	xml.Unmarshal([]byte(xmlStr), &res)

	// fmt.Print(res.WebSpace.SwitchSub.Result.Status)
	// if err != nil {
	// 	fmt.Printf("error: %v", err)
	// 	return "1"
	// }
	// id := res.Server.CreateSession.Result.ID
	// 打印解析结果
	// fmt.Printf("token: %+v", id)

	return res.WebSpace.SwitchSub.Result.Status
}
