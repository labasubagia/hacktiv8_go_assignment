package main

import (
	"final_project/controllers"
	"final_project/database"
	"final_project/middlewares"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	userCtrl := controllers.NewUserController(db)
	userRouter := router.Group("/users")
	{
		userRouter.POST("/register", userCtrl.Register)
		userRouter.POST("/login", userCtrl.Login)
		userRouter.PUT("", middlewares.Authentication(), userCtrl.Update)
		userRouter.DELETE("", middlewares.Authentication(), userCtrl.Delete)
	}

	if err := router.Run(fmt.Sprintf(":%v", os.Getenv("PORT"))); err != nil {
		panic(err)
	}
}
