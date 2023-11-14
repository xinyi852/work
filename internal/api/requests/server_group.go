package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

func DeletePleskServerGroupData(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"group_id": []string{"required"},
	}
	messages := govalidator.MapData{}
	return validate(data, rules, messages)
}

func UpdatePleskServerGroupData(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"group_id":  []string{"required"},
		"name":      []string{"required"},
		"fill_type": []string{"required"},
	}
	messages := govalidator.MapData{}
	return validate(data, rules, messages)
}
