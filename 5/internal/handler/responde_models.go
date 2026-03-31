package handler

import (
	"order-service/internal/database"
	"time"
)

type OrderResponse struct {
	OrderUID    string           `json:"order_uid"`
	TrackNumber string           `json:"track_number"`
	Status      string           `json:"status"`
	DateCreated time.Time        `json:"date_created"`
	Delivery    DeliveryResponse `json:"delivery"`
	Payment     PaymentResponse  `json:"payment"`
	Items       []ItemResponse   `json:"items"`
}

type DeliveryResponse struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	City    string `json:"city"`
	Address string `json:"address"`
	Zip     string `json:"zip"`
	Email   string `json:"email"`
}

type PaymentResponse struct {
	Transaction string `json:"transaction"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	Provider    string `json:"provider"`
	Bank        string `json:"bank"`
}

type ItemResponse struct {
	Name    string `json:"name"`
	Brand   string `json:"brand"`
	Price   int    `json:"price"`
	Size    string `json:"size"`
	Article int    `json:"article"`
}

func toOrderResponse(order database.Order) OrderResponse {
	items := make([]ItemResponse, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, ItemResponse{
			Name:    item.Name,
			Brand:   item.Brand,
			Price:   item.Price,
			Size:    item.Size,
			Article: item.ChrtID,
		})
	}

	return OrderResponse{
		OrderUID:    order.OrderUID,
		TrackNumber: order.TrackNumber,
		Status:      "Обработан",
		DateCreated: order.DateCreated,
		Delivery: DeliveryResponse{
			Name:    order.Delivery.Name,
			Phone:   order.Delivery.Phone,
			City:    order.Delivery.City,
			Address: order.Delivery.Address,
			Zip:     order.Delivery.Zip,
			Email:   order.Delivery.Email,
		},
		Payment: PaymentResponse{
			Transaction: order.Payment.Transaction,
			Amount:      order.Payment.Amount,
			Currency:    order.Payment.Currency,
			Provider:    order.Payment.Provider,
			Bank:        order.Payment.Bank,
		},
		Items: items,
	}
}
