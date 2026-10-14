package order

import (
	"context"
	"errors"
	"sync"
	"time"

	"main.go/dish"
	dishlist "main.go/dish_list"
)

var mtxID sync.Mutex
var idCounter int = 1

type OrderStatus string

const (
	StatusPending  OrderStatus = "pending"
	StatusCanceled OrderStatus = "canceled"
	StatusTimeout  OrderStatus = "timeout"
	StatusCooked   OrderStatus = "cooked"
	StatusFailed   OrderStatus = "failed"
)

type Order struct {
	mtx        sync.RWMutex
	id         int
	listDishes *dishlist.List
	cancel     context.CancelFunc
	status     OrderStatus
}

type Snapshot struct {
	ID     int
	Status string
	Dishes []dish.Snapshot
}

func NewOrder(ctx context.Context, timeout time.Duration, menu *dish.ListAvailabilityDishes, dishes []string) (*Order, context.Context, error) {
	specs := make([]dish.AvailabilityDish, 0, len(dishes))
	for _, name := range dishes {
		spec, ok := menu.Get(name)
		if !ok {
			return nil, nil, errors.New("такого блюда не существует в меню")
		}
		specs = append(specs, spec)
	}

	ctxOrder, cancel := context.WithTimeout(ctx, timeout)
	listDishes := dishlist.NewList()

	for _, v := range specs {
		listDishes.Add(dish.NewOrderDish(v))
	}

	mtxID.Lock()
	defer mtxID.Unlock()
	order := &Order{
		id:         idCounter,
		listDishes: listDishes,
		cancel:     cancel,
		status:     StatusPending,
	}
	idCounter++

	return order, ctxOrder, nil
}

func (o *Order) Snapshot() Snapshot {
	o.mtx.RLock()
	defer o.mtx.RUnlock()

	return Snapshot{
		ID:     o.id,
		Status: string(o.status),
		Dishes: o.listDishes.Snapshots(),
	}
}

func (o *Order) GetID() int {
	o.mtx.RLock()
	defer o.mtx.RUnlock()
	return o.id
}

func (o *Order) Dishes() map[int]dish.OrderDish {
	o.mtx.RLock()
	defer o.mtx.RUnlock()
	tmp := o.listDishes.Get()
	return tmp
}

func (o *Order) UpdateDish(d dish.OrderDish, od dish.UpdateOrderDish) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	o.listDishes.Update(d, od)
}

func (o *Order) Finish(err error) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	defer o.cancel()

	if o.status != StatusPending {
		return
	}

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		o.status = StatusTimeout
	case errors.Is(err, context.Canceled):
		o.status = StatusCanceled
	default:
		tmp := o.listDishes.Get()
		for _, v := range tmp {
			if status := v.GetStatus(); status == dish.StatusTimeout {
				o.status = StatusTimeout
				return
			}

			if status := v.GetStatusFailed(); status != "" {
				o.status = StatusFailed
				return
			}
		}
		o.status = StatusCooked
	}
}

func (o *Order) Cancel() {
	o.Finish(context.Canceled)
}
