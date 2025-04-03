package endpoint

import (
	"strconv"

	"github.com/ariebrainware/ltt-inventory-service/config"
	"github.com/ariebrainware/ltt-inventory-service/model"
	"github.com/ariebrainware/ltt-inventory-service/util"
	"github.com/gin-gonic/gin"
)

func parseQueryParams(c *gin.Context) (int, int, string) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	keyword := c.Query("keyword")
	return limit, offset, keyword
}

func fetchInventories(limit, offset int, keyword string) ([]model.InventoryMaster, error) {
	var inventories []model.InventoryMaster

	db, err := config.ConnectMySQL()
	if err != nil {
		return nil, err
	}

	query := db.Offset(offset).Order("created_at ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if keyword != "" {
		query = query.Where("product_name LIKE ? OR category LIKE ? OR brand LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Find(&inventories).Error; err != nil {
		return nil, err
	}

	return inventories, nil
}

func ListInventory(c *gin.Context) {
	limit, offset, keyword := parseQueryParams(c)

	inventories, err := fetchInventories(limit, offset, keyword)
	if err != nil {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "Failed to retrieve inventories",
			Err: err,
		})
		return
	}

	util.CallSuccessOK(c, util.APISuccessParams{
		Msg:  "Inventories retrieved",
		Data: map[string]interface{}{"inventories": inventories},
	})
}
