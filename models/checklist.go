package models

type Checklist struct {
	ID            string `gorm:"primaryKey" json:"id"`
	ChecklistName string `json:"name"`
	Username      string `json:"username"`
}
