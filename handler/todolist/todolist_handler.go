package todolist_handler

import (
	"net/http"
	"todolist/helper"
	"todolist/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var DB *gorm.DB

func RegisterTodoHandlers(db *gorm.DB) {
	DB = db
}

func GetChecklists(c echo.Context) error {
	var checklists []models.Checklist
	if err := DB.Find(&checklists).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not retrieve checklists"})
	}
	return c.JSON(http.StatusOK, checklists)
}

func CreateChecklist(c echo.Context) error {
	var newChecklist models.Checklist
	if err := c.Bind(&newChecklist); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	newChecklist.Username = c.Get("username").(string)
	newChecklist.ID = helper.GenerateRandomString(10)
	if err := DB.Create(&newChecklist).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create checklist"})
	}

	return c.JSON(http.StatusCreated, newChecklist)
}

func DeleteChecklist(c echo.Context) error {
	id := c.Param("id")

	if err := DB.Delete(&models.Checklist{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete checklist"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Checklist deleted successfully"})
}

func CreateChecklistItem(c echo.Context) error {
	checklistID := c.Param("id")

	var newItem models.ChecklistItem
	if err := c.Bind(&newItem); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	newItem.ID = helper.GenerateRandomString(10)
	newItem.IDChecklist = checklistID
	newItem.IsChecked = false
	if err := DB.Create(&newItem).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create checklist item"})
	}

	return c.JSON(http.StatusCreated, newItem)
}

func GetAllChecklistItems(c echo.Context) error {
	checklistID := c.Param("id")

	var items []models.ChecklistItem
	if err := DB.Where("id_checklist = ?", checklistID).Find(&items).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not retrieve checklist items"})
	}

	return c.JSON(http.StatusOK, items)
}

func DeleteChecklistItem(c echo.Context) error {
	itemID := c.Param("idItem")

	if err := DB.Delete(&models.ChecklistItem{}, "id = ?", itemID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete checklist item"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Checklist item deleted successfully"})
}

func GetChecklistItemByID(c echo.Context) error {
	itemID := c.Param("idItem")

	var item models.ChecklistItem
	if err := DB.First(&item, "id = ?", itemID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Checklist item not found"})
	}

	return c.JSON(http.StatusOK, item)
}

func RenameChecklistItem(c echo.Context) error {
	itemID := c.Param("idItem")

	var updatedItem models.ChecklistItem
	if err := c.Bind(&updatedItem); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	if err := DB.Model(&models.ChecklistItem{}).Where("id = ?", itemID).Update("item_name", updatedItem.ItemName).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not rename checklist item"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Checklist item renamed successfully"})
}

func UpdateChecklistItemStatus(c echo.Context) error {
	itemID := c.Param("idItem")

	var item models.ChecklistItem
	if err := DB.First(&item, "id = ?", itemID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Checklist item not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not retrieve checklist item"})
	}

	if err := DB.Model(&item).Update("is_checked", true).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update checklist item status"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Checklist item status updated successfully"})
}
