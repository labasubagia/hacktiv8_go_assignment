package main

import (
	"assignment2/controllers"
	"assignment2/database"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

const Port = 8000

func main() {
	db, err := database.NewPostgres()
	if err != nil {
		log.Fatalf("cannot connect to database: %s", err)
		return
	}

	router := gin.Default()

	order := router.Group("/orders")
	orderCtrl := controllers.NewOrderController(db)
	order.POST("", orderCtrl.CreateOrder)
	order.GET("", orderCtrl.GetOrders)
	order.PUT("/:orderId", orderCtrl.UpdateOrder)
	order.DELETE("/:orderId", orderCtrl.DeleteOrder)

	if err := router.Run(fmt.Sprintf(":%d", Port)); err != nil {
		log.Fatalf("Server failed to start: %s", err)
		return
	}
}
