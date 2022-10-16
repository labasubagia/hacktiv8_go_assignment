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

	photoCtrl := controllers.NewPhotoController(db)
	photoRouter := router.Group("/photos", middlewares.Authentication())
	{
		photoRouter.POST("", photoCtrl.Create)
		photoRouter.GET("", photoCtrl.List)
		photoRouter.PUT("/:photoId", photoCtrl.Update)
		photoRouter.DELETE("/:photoId", photoCtrl.Delete)
	}

	commentCtrl := controllers.NewCommentController(db)
	commentRouter := router.Group("/comments", middlewares.Authentication())
	{
		commentRouter.POST("", commentCtrl.Create)
		commentRouter.GET("", commentCtrl.List)
		commentRouter.PUT("/:commentId", commentCtrl.Update)
		commentRouter.DELETE("/:commentId", commentCtrl.Delete)
	}

	socialMediaCtrl := controllers.NewSocialMediaController(db)
	socialMediaRouter := router.Group("/socialmedias", middlewares.Authentication())
	{
		socialMediaRouter.POST("", socialMediaCtrl.Create)
		socialMediaRouter.GET("", socialMediaCtrl.List)
		socialMediaRouter.PUT("/:socialMediaId", socialMediaCtrl.Update)
		socialMediaRouter.DELETE("/:socialMediaId", socialMediaCtrl.Delete)
	}

	if err := router.Run(fmt.Sprintf(":%v", os.Getenv("PORT"))); err != nil {
		panic(err)
	}
}
