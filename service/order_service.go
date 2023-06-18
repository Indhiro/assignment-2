package service

import (
	"assignment-2/dto"
	"assignment-2/entity"
	"assignment-2/pkg/errs"
	"assignment-2/repository/order_repository"
	"fmt"
	"net/http"
)

type orderService struct {
	orderRepo   order_repository.OrderRepository
	itemService ItemService
}

type OrderService interface {
	CreateOrder(payload dto.NewOrderRequest) (*dto.NewOrderResponse, errs.MessageErr)
	UpdateOrder(orderId int, payload dto.NewOrderRequest) (*dto.GetOrderResponse, errs.MessageErr)
	GetAllOrders() ([]dto.GetOrderResponse, errs.MessageErr)
	DeleteOrder(orderId int) (*dto.DeleteOrderResponse, errs.MessageErr)
}

func NewOrderService(orderRepo order_repository.OrderRepository, itemService ItemService) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		itemService: itemService,
	}
}

func (o *orderService) CreateOrder(payload dto.NewOrderRequest) (*dto.NewOrderResponse, errs.MessageErr) {
	orderPayload := entity.Order{
		OrderedAt:    payload.OrderedAt,
		CustomerName: payload.CustomerName,
	}

	itemsPayload := []entity.Item{}

	for _, eachItem := range payload.Items {
		item := entity.Item{
			ItemCode:    eachItem.ItemCode,
			Quantity:    eachItem.Quantity,
			Description: eachItem.Description,
		}

		itemsPayload = append(itemsPayload, item)
	}

	newOrder, err := o.orderRepo.CreateOrder(orderPayload, itemsPayload)

	if err != nil {
		return nil, err
	}

	response := &dto.NewOrderResponse{
		StatusCode: http.StatusCreated,
		Message:    "Success",
		Data: dto.NewOrderRequest{
			OrderedAt:    newOrder.OrderedAt,
			CustomerName: newOrder.CustomerName,
			Items:        payload.Items,
		},
	}

	return response, nil
}

func (o *orderService) GetAllOrders() ([]dto.GetOrderResponse, errs.MessageErr) {

	orders, err := o.orderRepo.GetAllOrders()
	if err != nil {
		return nil, err
	}

	ordersResponse := []dto.GetOrderResponse{}

	for _, eachOrder := range orders {
		itemsResponse := []dto.ItemResponse{}

		for _, eachItem := range eachOrder.Items {
			itemResponse := eachItem.ItemToItemResponse()

			itemsResponse = append(itemsResponse, itemResponse)
		}

		orderResponse := dto.GetOrderResponse{
			Code: http.StatusOK,
			Data: dto.OrderResponse{
				Id:           eachOrder.Order.OrderId,
				CreatedAt:    eachOrder.Order.CreatedAt,
				UpdatedAt:    eachOrder.Order.UpdatedAt,
				CustomerName: eachOrder.Order.CustomerName,
				Items:        itemsResponse,
			},
		}

		ordersResponse = append(ordersResponse, orderResponse)
	}

	return ordersResponse, nil
}

func (o *orderService) UpdateOrder(orderId int, payload dto.NewOrderRequest) (*dto.GetOrderResponse, errs.MessageErr) {
	itemCodes := payload.ItemsToItemCode() //[]string{"123", "456"}

	_, err := o.itemService.FindItemsByItemCodes(itemCodes)

	if err != nil {
		return nil, errs.NewInternalServerError("Failed to begin transaction")
	}

	orderPayload := entity.Order{
		OrderId:      orderId,
		OrderedAt:    payload.OrderedAt,
		CustomerName: payload.CustomerName,
	}

	itemsPayload := []entity.Item{}

	for _, eachItem := range payload.Items {
		item := entity.Item{
			ItemCode:    eachItem.ItemCode,
			Quantity:    eachItem.Quantity,
			Description: eachItem.Description,
		}

		itemsPayload = append(itemsPayload, item)
	}

	orderItem, err := o.orderRepo.UpdateOrder(orderPayload, itemsPayload)

	if err != nil {
		return nil, errs.NewNotFound("Order not found")
	}

	itemsResponse := []dto.ItemResponse{}

	for _, eachItem := range orderItem.Items {
		itemResponse := eachItem.ItemToItemResponse()

		itemsResponse = append(itemsResponse, itemResponse)
	}

	result := dto.GetOrderResponse{
		Code: http.StatusOK,
		Data: dto.OrderResponse{
			Id:           orderItem.Order.OrderId,
			CreatedAt:    orderItem.Order.CreatedAt,
			UpdatedAt:    orderItem.Order.UpdatedAt,
			CustomerName: orderItem.Order.CustomerName,
			Items:        itemsResponse,
		},
	}

	return &result, nil
}

func (o *orderService) DeleteOrder(orderId int) (*dto.DeleteOrderResponse, errs.MessageErr) {
	if err := o.orderRepo.DeleteOrder(orderId); err != nil {
		return nil, err
	}

	response := &dto.DeleteOrderResponse{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("Order with id %d has been successfully deleted", orderId),
	}
	return response, nil
}
