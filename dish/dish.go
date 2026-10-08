package dish

import (
	"context"
	"sync"
	"time"
)

var idCounter int = 1
var mtxID sync.Mutex

const (
	StatusCooking  = "cooking"
	StatusCooked   = "cooked"
	StatusCanceled = "cancelled"
	StatusTimeout  = "timeout"
)

// OrderDish представляет собой структуру, которая хранит информацию о конкретном готовящемся блюде
type OrderDish struct {
	id          int
	dish        AvailabilityDish // Ссылка на имеющиеся блюда
	isReady     bool
	cookingTime *time.Duration
	status      string
}

// NewOrderDish представляет собой конструктор для создания экземпляра структуры
func NewOrderDish(dish AvailabilityDish) OrderDish {
	mtxID.Lock()
	defer mtxID.Unlock()
	newDish := OrderDish{
		id:          idCounter,
		dish:        dish,
		cookingTime: nil,
		status:      StatusCooking,
	}
	idCounter++

	return newDish
}

type UpdateOrderDish struct {
	isReady     bool
	cookingTime *time.Duration
	status      string
}

func NewUpdateOrderDish(isReady bool, cookingTime *time.Duration, errText error) UpdateOrderDish {
	var status string

	switch errText {
	case context.DeadlineExceeded:
		status = StatusTimeout
	case context.Canceled:
		status = StatusCanceled
	default:
		status = StatusCooked
	}

	return UpdateOrderDish{
		isReady:     isReady,
		cookingTime: cookingTime,
		status:      status,
	}
}

func (od OrderDish) GetTimeout() time.Duration {
	return od.dish.GetTimeout()
}

func (od OrderDish) Update(newOD UpdateOrderDish) OrderDish {
	od.isReady = newOD.isReady
	od.cookingTime = newOD.cookingTime
	od.status = newOD.status
	return od
}

func (od OrderDish) GetID() int {
	return od.id
}
