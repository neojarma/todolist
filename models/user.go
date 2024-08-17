package models

type User struct {
	ID       string `gorm:"primaryKey" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
