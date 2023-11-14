// Package validators 存放自定义规则和验证器
package validators

import (
	"racent.com/pkg/database"
	"racent.com/pkg/helpers"
	"strconv"
)

// ExistsIds 验证 ID数组 是否存在
func ExistsIds(tableName, filed string, ids []uint64, errs map[string][]string) map[string][]string {
	if helpers.Empty(ids) {
		return errs
	}
	query := database.DB.Table(tableName).Select("id").Where(filed+" IN (?)", ids)
	var count int64
	query.Count(&count)

	if strconv.Itoa(len(ids)) != strconv.FormatInt(count, 10) {
		errs[filed] = append(errs[filed], "数据不合法！")
	}
	return errs
}

func CheckQuota(userId, productId uint64, errs map[string][]string) map[string][]string {
	var count int64
	database.DB.Table("user_certificate_quota").
		Where("`user_id` = ? AND `product_id` = ? AND `remain_quota` > 0", userId, productId).Count(&count)
	if count == 0 {
		errs["user_id"] = append(errs["user_id"], "免费证书配额不足")
	}
	return errs
}
