package handler

import (
	"assignment-2/dto"
	"assignment-2/pkg/errs"
	"assignment-2/pkg/helpers"
	"assignment-2/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *orderHandler {
	return &orderHandler{orderService: orderService}
}

func (o *orderHandler) CreateOrder(ctx *gin.Context) {
	var requestBody dto.NewOrderRequest

	if ErrMessage := ctx.ShouldBindJSON(&requestBody); ErrMessage != nil {
		newError := errs.NewUnprocessableEntity(ErrMessage.Error())
		ctx.JSON(newError.StatusCode(), newError)
		return
	}

	newOrder, ErrMessage := o.orderService.CreateOrder(requestBody)
	if ErrMessage != nil {
		ctx.JSON(ErrMessage.StatusCode(), ErrMessage)
		return
	}

	ctx.JSON(newOrder.StatusCode, newOrder)
}

func (o *orderHandler) GetAllOrders(ctx *gin.Context) {
	orders, ErrMessage := o.orderService.GetAllOrders()
	if ErrMessage != nil {
		ctx.JSON(ErrMessage.StatusCode(), ErrMessage)
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (o *orderHandler) GetOrderById(ctx *gin.Context) {
	orders, ErrMessage := o.orderService.GetAllOrders()
	if ErrMessage != nil {
		ctx.JSON(ErrMessage.StatusCode(), ErrMessage)
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

func (o *orderHandler) UpdateOrderById(ctx *gin.Context) {
	orderId, ErrMessage := helpers.GetParamId(ctx, "orderId")
	if ErrMessage != nil {
		newError := errs.NewBadRequest("orderId should be an unsigned integer")
		ctx.JSON(newError.StatusCode(), newError)
		return
	}

	var requestBody dto.NewOrderRequest

	if ErrMessage := ctx.ShouldBindJSON(&requestBody); ErrMessage != nil {
		newError := errs.NewUnprocessableEntity(ErrMessage.Error())
		ctx.JSON(newError.StatusCode(), newError)
		return
	}

	updatedOrder, errOrder := o.orderService.UpdateOrder(orderId, requestBody)
	if errOrder != nil {
		ctx.JSON(errOrder.StatusCode(), errOrder)
		return
	}

	ctx.JSON(updatedOrder.Code, updatedOrder)
}

func (o *orderHandler) DeleteOrderById(ctx *gin.Context) {
	orderId, err := helpers.GetParamId(ctx, "orderId")

	if err != nil {
		newError := errs.NewBadRequest("orderId should be an unsigned integer")
		ctx.JSON(newError.StatusCode(), newError)
		return
	}

	response, err := o.orderService.DeleteOrder(orderId)

	if err != nil {
		ctx.JSON(400, err)
		return
	}

	ctx.JSON(response.StatusCode, response)
}
