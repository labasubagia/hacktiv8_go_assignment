package controllers

import (
	"assignment2/models"
	"assignment2/types"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderController struct {
	db *gorm.DB
}

func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{
		db: db,
	}
}

func (ctrl *OrderController) CreateOrder(c *gin.Context) {
	body := types.ReqWriteOrder{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid payload: %s", err))
		return
	}
	order := reqWriteOrderToOrderModel(body)
	if err := ctrl.db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Cannot create order: %s", err))
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (ctrl *OrderController) UpdateOrder(c *gin.Context) {
	IDParam := c.Param("orderId")
	ID, err := strconv.Atoi(IDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid param: %s", err))
		return
	}
	body := types.ReqWriteOrder{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid payload: %s", err))
		return
	}
	order := reqWriteOrderToOrderModel(body)
	order.OrderID = uint(ID)
	if err := ctrl.db.Select("Items").Updates(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Cannot update order: %s", err))
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (ctrl *OrderController) GetOrders(c *gin.Context) {
	orders := []models.Order{}
	if err := ctrl.db.Preload("Items").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Cannot get orders: %s", err))
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (ctrl *OrderController) DeleteOrder(c *gin.Context) {
	IDParam := c.Param("orderId")
	ID, err := strconv.Atoi(IDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid param: %s", err))
		return
	}
	order := models.Order{OrderID: uint(ID)}
	if err := ctrl.db.Select("Items").Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Cannot delete order: %s", err))
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("Success delete order id=%d", ID))
}

func reqWriteOrderToOrderModel(req types.ReqWriteOrder) models.Order {
	order := models.Order{
		OrderedAt:    req.OrderedAt,
		CustomerName: req.CustomerName,
	}
	for _, v := range req.Items {
		item := models.Item{
			ItemID:      v.ItemID,
			ItemCode:    v.ItemCode,
			Description: v.Description,
			Quantity:    v.Quantity,
			OrderID:     order.OrderID,
		}
		order.Items = append(order.Items, item)
	}
	return order
}
