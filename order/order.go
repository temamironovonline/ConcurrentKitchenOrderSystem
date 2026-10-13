package order

import (
	"context"
	"log"
	"sync"
	"time"

	"main.go/dish"
	dishlist "main.go/dish_list"
)

var mtxID sync.Mutex
var idCounter int = 1

type Order struct {
	mtx        sync.RWMutex
	id         int
	listDishes *dishlist.List
	cancel     context.CancelFunc
	status     string
}

type Snapshot struct {
	ID     int
	Status string
	Dishes []dish.Snapshot
}

func NewOrder(ctx context.Context, timeout time.Duration, dishes []string) (*Order, context.Context) {
	mtxID.Lock()
	defer mtxID.Unlock()

	newAvailabilityDishes := dish.NewListAvailabilityDishes()

	ctxOrder, cancel := context.WithTimeout(ctx, timeout)

	listDishes := dishlist.NewList()

	for _, v := range dishes {
		newDish := dish.NewOrderDish(newAvailabilityDishes.Get(v))
		listDishes.Add(newDish)
	}

	order := &Order{
		id:         idCounter,
		listDishes: listDishes,
		cancel:     cancel,
	}
	idCounter++

	return order, ctxOrder
}

func (o *Order) Snapshot() Snapshot {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	return Snapshot{
		ID:     o.id,
		Status: o.status,
		Dishes: o.listDishes.Snapshots(),
	}
}

func (o *Order) GetID() int {
	o.mtx.RLock()
	defer o.mtx.RUnlock()
	return o.id
}

func (o *Order) Dishes() map[int]dish.OrderDish {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	tmp := o.listDishes.Get()
	return tmp
}

func (o *Order) UpdateDish(d dish.OrderDish, od dish.UpdateOrderDish) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	o.listDishes.Update(d, od)
}

func (o *Order) Finish() {
	log.Println("Заглушка. Заказ готов")
}
