package dish

import (
	"context"
	"errors"
	"sync"
	"time"

	"main.go/consts"
)

var idCounter int = 1
var mtxID sync.Mutex

type StatusCook string

const (
	StatusPending  StatusCook = "pending"
	StatusCooked   StatusCook = "cooked"
	StatusCanceled StatusCook = "cancelled"
	StatusTimeout  StatusCook = "timeout"
	StatusFailed   StatusCook = "failed"
)

// OrderDish представляет собой структуру, которая хранит информацию о конкретном готовящемся блюде
type OrderDish struct {
	id           int
	dish         AvailabilityDish // Ссылка на имеющиеся блюда
	isReady      bool
	cookingTime  *time.Duration
	status       StatusCook
	statusFailed string
}

// NewOrderDish представляет собой конструктор для создания экземпляра структуры
func NewOrderDish(dish AvailabilityDish) OrderDish {
	mtxID.Lock()
	defer mtxID.Unlock()
	newDish := OrderDish{
		id:          idCounter,
		dish:        dish,
		cookingTime: nil,
		status:      StatusPending,
	}
	idCounter++

	return newDish
}

type UpdateOrderDish struct {
	id           int
	isReady      bool
	cookingTime  *time.Duration
	status       StatusCook
	statusFailed string
}

func NewUpdateOrderDish(isReady bool, cookingTime *time.Duration, errText error) UpdateOrderDish {
	var status StatusCook
	var statusFailed string

	switch {
	case errors.Is(errText, context.DeadlineExceeded):
		status = StatusTimeout
	case errors.Is(errText, context.Canceled):
		status = StatusCanceled
	default:
		if errText != nil {
			status = StatusFailed
			statusFailed = errText.Error()
			cookingTime = nil
			break
		}
		status = StatusCooked
	}

	return UpdateOrderDish{
		isReady:      isReady,
		cookingTime:  cookingTime,
		status:       status,
		statusFailed: statusFailed,
	}
}

type Snapshot struct {
	ID           int
	Name         string
	Status       StatusCook
	CookingTime  *time.Duration
	StatusFailed string
}

func (od OrderDish) Snapshot() Snapshot {
	return Snapshot{ID: od.id, Name: od.GetName(), Status: od.status, CookingTime: od.cookingTime, StatusFailed: od.statusFailed}
}

func (uod UpdateOrderDish) GetStatus() StatusCook {
	return uod.status
}

func (uod UpdateOrderDish) SetID(id int) UpdateOrderDish {
	uod.id = id
	return uod
}

func (uod UpdateOrderDish) GetStatusFailed() string {
	return uod.statusFailed
}

func (od OrderDish) GetTimeout() time.Duration {
	return od.dish.GetTimeout()
}

func (od OrderDish) Update(newOD UpdateOrderDish) OrderDish {
	od.isReady = newOD.isReady
	od.cookingTime = newOD.cookingTime
	od.status = newOD.status
	od.statusFailed = newOD.statusFailed
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
