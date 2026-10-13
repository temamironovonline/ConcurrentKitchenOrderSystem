package orderlist

import (
	"errors"
	"sync"

	"main.go/order"
)

type OrderList struct {
	mtx  sync.Mutex
	list map[int]*order.Order
}

func NewOrderList() *OrderList {
	return &OrderList{
		list: make(map[int]*order.Order),
	}
}

func (od *OrderList) Add(order *order.Order) {
	od.mtx.Lock()
	defer od.mtx.Unlock()
	od.list[order.GetID()] = order
}

func (od *OrderList) GetCertainOrder(id int) (order.Snapshot, error) {
	od.mtx.Lock()
	defer od.mtx.Unlock()
	o, ok := od.list[id]
	if !ok {
		return order.Snapshot{}, errors.New("такого заказа не существует")
	}

	return o.Snapshot(), nil

}
