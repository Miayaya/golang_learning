package main

import (
	"blog-api/config"
	"blog-api/internal/router"
	"blog-api/model"
)

func main() {
	db := config.InitDB()
	db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})

	r := router.SetupRouter(db)
	r.Run(":8080")
}
