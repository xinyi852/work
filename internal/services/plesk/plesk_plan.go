package plesk

import (
	"plesk/internal/models"
	"plesk/internal/pkg/common"
	"plesk/internal/pkg/paginator"

	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"racent.com/pkg/database"
	"racent.com/pkg/logger"
)

func CreatePlan(request *common.CreatePlanRequest) (*models.PleskPlan, error) {
	logger.DebugJSON("CreatePleskPlan", "request", request)
	pleskProductModel := &models.PleskPlan{
		Name:       request.Name,
		Plan:       request.Plan,
		MaxData:    request.MaxData,
		MaxTraffic: request.MaxTraffic,
	}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		tx.Create(pleskProductModel)
		if pleskProductModel.ID == 0 {
			logger.ErrorJSON("pleskProductModel", "request", request)
			return errors.New("Create plesk plan fail.")
		}
		return nil
	})
	return pleskProductModel, err
}

func ExistPlan(field, value string) bool {
	var count int64
	database.DB.Model(models.PleskPlan{}).Where(field+" = ?", value).Count(&count)
	return count > 0
}

func DelPlan(planId string) (PleskPlan models.PleskPlan) {
	database.DB.Model(models.PleskPlan{}).Where("id", planId).Delete(&PleskPlan)
	return
}

func GetPlan(idStr string) (PleskPlan models.PleskPlan) {
	database.DB.Where("id", idStr).First(&PleskPlan)
	return PleskPlan
}

func ProductPaginate(c *gin.Context, perPage int) (PleskPlans []models.PleskPlan, paging paginator.Paging) {
	queryDB := database.DB.Model(models.PleskPlan{})
	// keyword := c.DefaultQuery("keyword", "")
	// if keyword != "" {
	//     queryDB.Where("name LIKE ?", "%"+keyword+"%")
	// }
	paging = paginator.Paginate(
		c,
		queryDB,
		&PleskPlans,
		perPage,
	)

	return
}
