package order

import (
	"sync/atomic"

	dishlist "main.go/dish_list"
)

var idOrder atomic.Int64

type Order struct {
	id         atomic.Int64
	listDishes *dishlist.List
}

func NewOrder(list *dishlist.List) *Order {
	order := &Order{
		listDishes: list,
	}
	idOrder.Add(1)
	order.id.Add(idOrder.Load())
	return order
}
