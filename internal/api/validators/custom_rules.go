package validators

import (
	"errors"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"racent.com/pkg/database"
	"strconv"
	"strings"
	"unicode/utf8"
)

// 此方法会在初始化时执行，注册自定义表单验证规则
func Init() {
	// 自定义规则 exists，验证请求数据必须存在于数据库中。
	// 常用于验证数据库某个字段的值存在，如用户名、邮箱、手机号、或者分类的名称。
	// exist:users,email 检查数据库表里是否存在同一条信息
	govalidator.AddCustomRule("exist", func(field string, rule string, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "exist:"), ",")
		// 第一个参数，表名，如 users
		tableName := rng[0]
		// 第二个参数，字段名称，如 email 或者 phone
		dbFiled := rng[1]

		// 用户请求过来的数据
		requestValue := value.(uint64)

		// 拼接 SQL
		var result struct{ Found bool }
		database.DB.Raw(fmt.Sprintf("SELECT EXISTS (SELECT `id` FROM %s WHERE %s = ? LIMIT 1) AS `found`", tableName, dbFiled), requestValue).Find(&result)

		// 查询数据
		// 验证不同过，数据库能找到对应的数据
		if result.Found == false {
			// 如果有自定义错误消息的话
			if message != "" {
				return errors.New(message)
			}
			// 默认的错误消息
			return fmt.Errorf("%v 不存在", requestValue)
		}
		// 验证通过
		return nil
	})

	// 自定义规则 not_exists，验证请求数据必须不存在于数据库中。
	// 常用于保证数据库某个字段的值唯一，如用户名、邮箱、手机号、或者分类的名称。
	// not_exists 参数可以有两种，一种是 2 个参数，一种是 3 个参数：
	// not_exists:users,email 检查数据库表里是否存在同一条信息
	// not_exists:users,email,32 排除用户掉 id 为 32 的用户
	govalidator.AddCustomRule("not_exists", func(field string, rule string, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "not_exists:"), ",")

		// 第一个参数，表名称，如 users
		tableName := rng[0]
		// 第二个参数，字段名称，如 email 或者 phone
		dbFiled := rng[1]

		// 第三个参数，排除 ID
		var exceptID string
		if len(rng) > 2 {
			exceptID = rng[2]
		}

		// 用户请求过来的数据
		requestValue := value.(string)

		// 拼接 SQL
		//var result struct{ Found bool }
		//database.DB.Raw(fmt.Sprintf("SELECT EXISTS (SELECT `id` FROM %s WHERE %s = ? LIMIT 1) AS `found`", tableName, dbFiled), requestValue).Find(&result)
		query := database.DB.Table(tableName).Select("id").Where(dbFiled+" = ?", requestValue)

		// 如果传参第三个参数，加上 SQL Where 过滤
		if len(exceptID) > 0 {
			query.Where("id NOT IN (?)", exceptID)
		}

		// 查询数据库
		var count int64
		query.Count(&count)

		// 验证不通过，数据库能找到对应的数据
		if count != 0 {
			// 如果有自定义错误消息的话
			if message != "" {
				return errors.New(message)
			}
			// 默认的错误消息
			return fmt.Errorf("%v 已被占用", requestValue)
		}
		// 验证通过
		return nil
	})

	// 自定义规则 max_cn，验证请求中文长度最长范围
	// max_cn:8 中文长度设定不超过 8
	govalidator.AddCustomRule("max_cn", func(field string, rule string, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "max_cn:"))
		if valLength > l {
			// 如果有自定义错误消息的话，使用自定义消息
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("The %s field must be maximum %d char", field, l)
		}
		return nil
	})

	// 自定义规则 min_cn，验证请求中文长度最短范围
	// min_cn:2 中文长度设定不小于 2
	govalidator.AddCustomRule("min_cn", func(field string, rule string, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "min_cn:"))
		if valLength < l {
			// 如果有自定义错误消息的话，使用自定义消息
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长度需大于 %d 个字", l)
		}
		return nil
	})

}
