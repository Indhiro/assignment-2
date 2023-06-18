package handler

import (
	"assignment-2/database"
	"assignment-2/repository/item_repository/item_pg"
	"assignment-2/repository/order_repository/order_pg"
	"assignment-2/service"

	"github.com/gin-gonic/gin"
)

func StartApp() {
	database.InitiliazeDatabase()

	db := database.GetDatabaseInstance()

	itemRepo := item_pg.NewItemPG(db)

	itemService := service.NewItemService(itemRepo)

	orderRepo := order_pg.NewOrderPG(db)

	orderService := service.NewOrderService(orderRepo, itemService)

	orderHandler := NewOrderHandler(orderService)

	r := gin.Default()

	r.POST("/orders", orderHandler.CreateOrder)

	r.GET("/orders", orderHandler.GetAllOrders)

	r.PUT("/orders/:orderId", orderHandler.UpdateOrderById)

	r.DELETE("/orders/:orderId", orderHandler.DeleteOrderById)

	r.Run(":8080")
}
