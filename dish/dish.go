package dish

import (
	"context"
	"sync"
	"time"

	"main.go/consts"
)

var idCounter int = 1
var mtxID sync.Mutex

type StatusCook string

const (
	StatusCooking  StatusCook = "cooking"
	StatusCooked   StatusCook = "cooked"
	StatusCanceled StatusCook = "cancelled"
	StatusTimeout  StatusCook = "timeout"
)

// OrderDish представляет собой структуру, которая хранит информацию о конкретном готовящемся блюде
type OrderDish struct {
	id          int
	dish        AvailabilityDish // Ссылка на имеющиеся блюда
	isReady     bool
	cookingTime *time.Duration
	status      StatusCook
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
	id          int
	isReady     bool
	cookingTime *time.Duration
	status      StatusCook
}

func NewUpdateOrderDish(isReady bool, cookingTime *time.Duration, errText error) UpdateOrderDish {
	var status StatusCook

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

type Snapshot struct {
	ID     int
	Name   string
	Status StatusCook
}

func (od OrderDish) Snapshot() Snapshot {
	return Snapshot{ID: od.id, Name: od.GetName(), Status: od.status}
}

func (uod UpdateOrderDish) GetStatus() StatusCook {
	return uod.status
}

func (uod UpdateOrderDish) SetID(id int) UpdateOrderDish {
	uod.id = id
	return uod
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

func (od OrderDish) GetStation() consts.ClassStation {
	return od.dish.stat
}

func (od OrderDish) GetName() string {
	return od.dish.name
}

func (od OrderDish) GetStatus() StatusCook {
	return od.status
}
