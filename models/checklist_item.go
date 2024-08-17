package models

type ChecklistItem struct {
	ID          string `gorm:"primaryKey" json:"id"`
	IDChecklist string `json:"id_checklist"`
	ItemName    string `json:"itemName"`
	IsChecked   bool   `json:"is_checked"`
}
