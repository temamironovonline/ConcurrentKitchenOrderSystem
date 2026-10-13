package orderlist

import (
	"errors"
	"sync"

	"main.go/order"
)

var ErrOrderNotFound error = errors.New("заказ не найден")

type OrderList struct {
	mtx  sync.RWMutex
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

func (od *OrderList) GetOrderSnapshot(id int) (order.Snapshot, error) {
	od.mtx.RLock()
	defer od.mtx.RUnlock()
	o, ok := od.list[id]
	if !ok {
		return order.Snapshot{}, ErrOrderNotFound
	}

	return o.Snapshot(), nil
}

func (od *OrderList) GetOrder(id int) (*order.Order, error) {
	od.mtx.RLock()
	defer od.mtx.RUnlock()
	o, ok := od.list[id]
	if !ok {
		return nil, ErrOrderNotFound
	}
	return o, nil
}
