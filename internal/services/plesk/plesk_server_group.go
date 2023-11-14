package plesk

import (
	"plesk/internal/models"
	"plesk/internal/pkg/common"
	"plesk/internal/pkg/paginator"

	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"racent.com/pkg/database"
	"racent.com/pkg/helpers"
	"racent.com/pkg/logger"
)

// func Paginate(c *gin.Context, request *common.FreeCertListRequest) (freeCertificates []models.CertBasicInfo, paging paginator.Paging) {
// 	queryDB := database.DB.Model(models.FreeCertificate{})
// 	queryDB.Select(
// 		"id", "user_id", "product_id", "vendor_id", "serial_num", "apply_response",
// 		"subject_common_name", "not_before", "not_after", "status", "created_at",
// 	)
// 	queryDB.Where("user_id = ?", request.UserId)
// 	// keyword := c.DefaultQuery("keyword", "")
// 	if !helpers.Empty(request.Status) {
// 		queryDB.Where("status = ?", request.Status)
// 	}
// 	if !helpers.Empty(request.SerialNum) {
// 		queryDB.Where("serial_num = ?", request.SerialNum)
// 	}
// 	if !helpers.Empty(request.SubjectCommonName) {
// 		queryDB.Where("subject_common_name LIKE ?", "%"+request.SubjectCommonName+"%")
// 	}

// 	if !helpers.Empty(request.BeginTime) && !helpers.Empty(request.EndTime) {
// 		queryDB.Where("not_after BETWEEN ? AND ?", request.BeginTime, request.EndTime)
// 	}

// 	// if keyword != "" {
// 	//     queryDB.Where("name LIKE ?", "%"+keyword+"%")
// 	// }
// 	paging = paginator.Paginate(
// 		c,
// 		queryDB,
// 		&freeCertificates,
// 		25,
// 	)
// 	for i, item := range freeCertificates {
// 		freeCertificates[i].StatusDesc = models.StatusText(item.Status)
// 		freeCertificates[i].CASerialNum = item.GetCaOrderId()
// 	}
// 	return
// }

func CreatePleskServerGroup(request *common.CreatePleskServerGroupData) (*models.PleskServerGroup, error) {
	logger.DebugJSON("CreatePleskServerGroup", "request", request)
	PleskServerGroupModel := &models.PleskServerGroup{
		Name:     request.Name,
		FillType: request.FillType,
	}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		tx.Create(PleskServerGroupModel)
		if PleskServerGroupModel.ID == 0 {
			logger.ErrorJSON("PleskServerGroupModel", "request", request)
			return errors.New("Create plesk server group fail.")
		}
		return nil
	})
	return PleskServerGroupModel, err
}

func ExistPleskServerGroup(field, value string) bool {
	var count int64
	database.DB.Model(models.PleskServerGroup{}).Where(field+" = ?", value).Count(&count)
	return count > 0
}

func GetPleskServerGroup(field string, value string) (plesk_server_group models.PleskServerGroup) {
	database.DB.Model(models.PleskServerGroup{}).Where(field+" = ?", value).First(&plesk_server_group)
	return plesk_server_group
}

func PleskServerGroupPaginate(c *gin.Context, request *common.ServerGropuListRequest) (PleskServerGroups []models.PleskServerGroup, paging paginator.Paging) {
	queryDB := database.DB.Model(models.PleskServerGroup{})
	queryDB.Select(
		"id", "name", "area", "fill_type", "created_at",
	)
	// queryDB.Where("user_id = ?", request.UserId)
	// keyword := c.DefaultQuery("keyword", "")
	// if keyword != "" {
	//     queryDB.Where("name LIKE ?", "%"+keyword+"%")
	// }
	if !helpers.Empty(request.Name) {
		keyword := c.DefaultQuery("keyword", request.Name)
		if keyword != "" {
			queryDB.Where("name LIKE ?", "%"+keyword+"%")
		}
	}

	paging = paginator.Paginate(
		c,
		queryDB,
		&PleskServerGroups,
		25,
	)

	//处理数据
	// for i, item := range freeCertificates {
	// 	freeCertificates[i].StatusDesc = models.StatusText(item.Status)
	// 	freeCertificates[i].CASerialNum = item.GetCaOrderId()
	// }

	return
}

func DeletePleskServerGroup(group_id string) bool {
	var plesk_server_group models.PleskServerGroup
	res := database.DB.Model(models.PleskServerGroup{}).Where("id", group_id).Delete(&plesk_server_group)
	return res.RowsAffected == 1
}

func UpdatePleskServerGroup(request *common.UpdatePleskServerGroupData) bool {
	plesk_server_group := GetPleskServerGroup("id", request.GroupId)
	plesk_server_group.Name = request.Name
	plesk_server_group.FillType = request.FillType
	res := plesk_server_group.Save()
	return res == 1
}
