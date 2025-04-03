package endpoint

import (
	"errors"
	"strconv"

	"github.com/ariebrainware/ltt-inventory-service/config"
	"github.com/ariebrainware/ltt-inventory-service/model"
	"github.com/ariebrainware/ltt-inventory-service/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getDBConnectionOrHandleError(c *gin.Context) *gorm.DB {
	db, err := config.ConnectMySQL()
	if err != nil {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "Failed to connect to MySQL",
			Err: err,
		})
		c.Abort()
		return nil
	}
	return db
}

func parseQueryParams(c *gin.Context) (int, int, string) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	keyword := c.Query("keyword")
	return limit, offset, keyword
}

func fetchInventories(c *gin.Context, limit, offset int, keyword string) ([]model.InventoryMaster, error) {
	var inventories []model.InventoryMaster

	db := getDBConnectionOrHandleError(c)
	if db == nil {
		return nil, errors.New("failed to connect to database")
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

	inventories, err := fetchInventories(c, limit, offset, keyword)
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

func CreateInventory(c *gin.Context) {
	var inventory model.InventoryMaster

	if err := c.ShouldBindJSON(&inventory); err != nil {
		util.CallUserError(c, util.APIErrorParams{
			Msg: "Invalid request payload",
			Err: err,
		})
		return
	}

	db := getDBConnectionOrHandleError(c)
	if db == nil {
		return
	}

	// Check if the inventory already exists
	var existingInventory model.InventoryMaster
	if err := db.Where("product_name = ?", inventory.ProductName).First(&existingInventory).Error; err == nil {
		util.CallUserError(c, util.APIErrorParams{
			Msg: "Inventory already exists",
			Err: errors.New("inventory already exists"),
		})
		return
	}

	// Create the new inventory
	if err := db.Create(&inventory).Error; err != nil {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "Failed to create inventory",
			Err: err,
		})
		return
	}

	util.CallSuccessOK(c, util.APISuccessParams{
		Msg:  "Inventory created successfully",
		Data: inventory,
	})
}

func GetInventory(c *gin.Context) {
	id := c.Param("id")
	var inventory model.InventoryMaster

	db := getDBConnectionOrHandleError(c)
	if db == nil {
		return
	}

	if err := db.First(&inventory, id).Error; err != nil {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "Failed to retrieve inventory",
			Err: err,
		})
		return
	}
	if inventory.ID == 0 {
		util.CallUserError(c, util.APIErrorParams{
			Msg: "Inventory not found",
			Err: errors.New("inventory not found"),
		})
		return
	}
	util.CallSuccessOK(c, util.APISuccessParams{
		Msg:  "Inventory retrieved successfully",
		Data: inventory,
	})
	return
}

func UpdateInventory(c *gin.Context) {
	var inventory model.InventoryMaster
	id := c.Param("id")

	if err := c.ShouldBindJSON(&inventory); err != nil {
		util.CallUserError(c, util.APIErrorParams{
			Msg: "Invalid request payload",
			Err: err,
		})
		return
	}

	db := getDBConnectionOrHandleError(c)
	if db == nil {
		return
	}

	if err := db.Model(&inventory).Where("id = ?", id).Updates(inventory).Error; err != nil {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "Failed to update inventory",
			Err: err,
		})
		return
	}
	util.CallSuccessOK(c, util.APISuccessParams{
		Msg:  "Inventory updated successfully",
		Data: inventory,
	})
	return
}

func DeleteInventory(c *gin.Context) {
	id := c.Param("id")

	db := getDBConnectionOrHandleError(c)
	if db == nil {
		return
	}

	if err := db.Delete(&model.InventoryMaster{}, id).Error; err != nil {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "Failed to delete inventory",
			Err: err,
		})
		return
	}

	util.CallSuccessOK(c, util.APISuccessParams{
		Msg: "Inventory deleted successfully",
	})
}
