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

// OrderStatus предназначен для обозначения статусов приготовления заказа
type OrderStatus string

const (
	StatusPending  OrderStatus = "pending"
	StatusCanceled OrderStatus = "canceled"
	StatusTimeout  OrderStatus = "timeout"
	StatusCooked   OrderStatus = "cooked"
	StatusFailed   OrderStatus = "failed"
)

// Order структура, хранящая информацию о заказе
type Order struct {
	mtx        sync.RWMutex
	id         int
	listDishes *dishlist.List
	cancel     context.CancelFunc
	status     OrderStatus
}

// Snapshot предназначен для передачи информации в хендлер, с последующей отправкой пользователю
type Snapshot struct {
	ID     int
	Status string
	Dishes []dish.Snapshot
}

// NewOrder создает новый заказ
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

// NewSnapshot функция, которая делает Snapshot
func (o *Order) NewSnapshot() Snapshot {
	o.mtx.RLock()
	defer o.mtx.RUnlock()

	return Snapshot{
		ID:     o.id,
		Status: string(o.status),
		Dishes: o.listDishes.Snapshots(),
	}
}

// GetID позволяет получить ID заказа
func (o *Order) GetID() int {
	o.mtx.RLock()
	defer o.mtx.RUnlock()
	return o.id
}

// Dishes позволяет получить копию блюд в заказе
func (o *Order) Dishes() map[int]dish.OrderDish {
	o.mtx.RLock()
	defer o.mtx.RUnlock()
	tmp := o.listDishes.Get()
	return tmp
}

// UpdateDish предназначен для обновления информации о блюде
func (o *Order) UpdateDish(d dish.OrderDish, od dish.UpdateOrderDish) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	o.listDishes.Update(d, od)
}

// Finish используется, когда по какой-то из причин приготовление блюда было закончено
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

// Cancel предназначен для отмены заказа
func (o *Order) Cancel() {
	o.Finish(context.Canceled)
}
