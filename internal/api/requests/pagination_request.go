package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type PaginationRequest struct {
	Sort    string `valid:"sort" json:"sort" form:"sort"`
	Order   string `valid:"order" json:"order" form:"order"`
	PerPage string `valid:"per_page" json:"per_page" form:"per_page"`
}

func Pagination(data interface{}, c *gin.Context) map[string][]string {

	rules := govalidator.MapData{
		"sort":     []string{"in:id,created_at,updated_at"},
		"order":    []string{"in:asc,desc"},
		"per_page": []string{"numeric_between:0,100"},
	}
	messages := govalidator.MapData{
		"sort": []string{
			"in:排序字段仅支持 id,created_at,updated_at",
		},
		"order": []string{
			"in:排序规则仅支持 asc（正序）,desc（倒序）",
		},
		"per_page": []string{
			"numeric_between:每页条数的值介于 2~100 之间",
		},
	}
	return validate(data, rules, messages)
}

func FreeCertList(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"user_id":  []string{"required"},
		"sort":     []string{"in:id,created_at,updated_at"},
		"order":    []string{"in:asc,desc"},
		"per_page": []string{"numeric_between:0,100"},
	}
	messages := govalidator.MapData{
		"sort": []string{
			"in:排序字段仅支持 id,created_at,updated_at",
		},
		"order": []string{
			"in:排序规则仅支持 asc（正序）,desc（倒序）",
		},
		"per_page": []string{
			"numeric_between:每页条数的值介于 2~100 之间",
		},
	}
	return validate(data, rules, messages)
}

func PleskServerGroupList(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"sort":     []string{"in:id,created_at,updated_at"},
		"order":    []string{"in:asc,desc"},
		"per_page": []string{"numeric_between:0,100"},
	}
	messages := govalidator.MapData{
		"sort": []string{
			"in:排序字段仅支持 id,created_at,updated_at",
		},
		"order": []string{
			"in:排序规则仅支持 asc（正序）,desc（倒序）",
		},
		"per_page": []string{
			"numeric_between:每页条数的值介于 2~100 之间",
		},
	}
	return validate(data, rules, messages)
}

func CreateHosting(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"domain": []string{"required"},
		"plan":   []string{"required"},
		"module": []string{"required"},
		"area":   []string{"required"},
	}
	messages := govalidator.MapData{}
	return validate(data, rules, messages)
}

func UpgradePlesk(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"id":   []string{"required"},
		"plan": []string{"required"},
	}
	messages := govalidator.MapData{}
	return validate(data, rules, messages)
}

func SetPlesk(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"id": []string{"required"},
	}
	messages := govalidator.MapData{}
	return validate(data, rules, messages)
}

func ChangePasswordPlesk(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"id": []string{"required"},
		// "account":  []string{"required"},
		"password": []string{"required"},
		// "new_password": []string{"required"},
	}
	messages := govalidator.MapData{}
	return validate(data, rules, messages)
}

func CreatePleskServer(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"name":         []string{"required"},
		"account":      []string{"required"},
		"password":     []string{"required"},
		"ip":           []string{"required"},
		"max":          []string{"required"},
		"service":      []string{"required"},
		"url":          []string{"required"},
		"server_group": []string{"required"},
		"type":         []string{"required"},
	}
	messages := govalidator.MapData{}
	return validate(data, rules, messages)
}

func UpdatePleskServer(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"id":           []string{"required"},
		"name":         []string{"required"},
		"account":      []string{"required"},
		"password":     []string{"required"},
		"ip":           []string{"required"},
		"max":          []string{"required"},
		"service":      []string{"required"},
		"url":          []string{"required"},
		"server_group": []string{"required"},
		"type":         []string{"required"},
	}
	messages := govalidator.MapData{}
	return validate(data, rules, messages)
}

func CreatePleskServerGroupData(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"name":      []string{"required"},
		"fill_type": []string{"required"},
	}
	messages := govalidator.MapData{}
	return validate(data, rules, messages)
}

func CreatePleskProductRequest(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"name":        []string{"required"},
		"plan":        []string{"required"},
		"max_data":    []string{"required"},
		"max_traffic": []string{"required"},
	}
	messages := govalidator.MapData{}
	return validate(data, rules, messages)
}

func DeletePleskProductRequest(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"id": []string{"required"},
	}
	messages := govalidator.MapData{}
	return validate(data, rules, messages)
}

func DeleteServerRequest(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"id": []string{"required"},
	}
	messages := govalidator.MapData{}
	return validate(data, rules, messages)
}

func UpdatePleskProductRequest(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"id":          []string{"required"},
		"name":        []string{"required"},
		"plan":        []string{"required"},
		"max_data":    []string{"required"},
		"max_traffic": []string{"required"},
	}
	messages := govalidator.MapData{}
	return validate(data, rules, messages)
}

func PleskApiRequest(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"data": []string{"required"},
	}
	messages := govalidator.MapData{}
	return validate(data, rules, messages)
}
