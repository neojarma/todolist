package main

import (
	"log"
	"os"
	"todolist/auth"
	todolist_handler "todolist/handler/todolist"
	user_handler "todolist/handler/user"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	user_handler.RegisterUserHandlers(db)
	todolist_handler.RegisterTodoHandlers(db)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/register", user_handler.Register)
	e.POST("/login", user_handler.Login)

	e.GET("/checklist", todolist_handler.GetChecklists, auth.ExtractUsername)
	e.POST("/checklist", todolist_handler.CreateChecklist, auth.ExtractUsername)
	e.DELETE("/checklist/:id", todolist_handler.DeleteChecklist, auth.ExtractUsername)

	e.GET("/checklist/:id/item", todolist_handler.GetAllChecklistItems, auth.ExtractUsername)
	e.POST("/checklist/:id/item", todolist_handler.CreateChecklistItem, auth.ExtractUsername)
	e.DELETE("/checklist/:id/item/:idItem", todolist_handler.DeleteChecklistItem, auth.ExtractUsername)
	e.GET("/checklist/:id/item/:idItem", todolist_handler.GetChecklistItemByID, auth.ExtractUsername)
	e.PUT("/checklist/:id/item/rename/:idItem", todolist_handler.RenameChecklistItem, auth.ExtractUsername)
	e.PUT("/checklist/:id/item/:idItem", todolist_handler.UpdateChecklistItemStatus, auth.ExtractUsername)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
