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

// StatusCook для определения различных статусов готовности блюда
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

// UpdateOrderDish представляет собой тип, который хранит в себе обновленную информацию о блюде, чтобы потом перенести ее в основное блюдо
type UpdateOrderDish struct {
	id           int
	isReady      bool
	cookingTime  *time.Duration
	status       StatusCook
	statusFailed string
}

// NewUpdateOrderDish создает экземпляр типа
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

// Snapshot тип, который используется для передачи информации в handler и ответа пользователю
type Snapshot struct {
	ID           int
	Name         string
	Status       StatusCook
	CookingTime  *time.Duration
	StatusFailed string
}

// Snapshot метод, который создает экземпляр типа Snapshot
func (od OrderDish) Snapshot() Snapshot {
	return Snapshot{ID: od.id, Name: od.GetName(), Status: od.status, CookingTime: od.cookingTime, StatusFailed: od.statusFailed}
}

// GetStatus позволяет получить статус приготовления блюда
func (uod UpdateOrderDish) GetStatus() StatusCook {
	return uod.status
}

// SetID позволяет установить ID для обновляемых данных по блюду
func (uod UpdateOrderDish) SetID(id int) UpdateOrderDish {
	uod.id = id
	return uod
}

// GetStatusFailed предназначен для получения причины провала блюда типа UpdateOrderDish. Если провала не было, то и вернется пустая строка
func (uod UpdateOrderDish) GetStatusFailed() string {
	return uod.statusFailed
}

// GetTimeout предназначен для получения заданного времени приготовления блюда
func (od OrderDish) GetTimeout() time.Duration {
	return od.dish.GetTimeout()
}

// Update предназначен для обновления данных о блюде
func (od OrderDish) Update(newOD UpdateOrderDish) OrderDish {
	od.isReady = newOD.isReady
	od.cookingTime = newOD.cookingTime
	od.status = newOD.status
	od.statusFailed = newOD.statusFailed
	return od
}

// GetID получение ID блюда
func (od OrderDish) GetID() int {
	return od.id
}

// GetStation получение станции, на которой готовится блюдо
func (od OrderDish) GetStation() consts.ClassStation {
	return od.dish.stat
}

// GetName получение названия блюда
func (od OrderDish) GetName() string {
	return od.dish.name
}

// GetStatus получение статуса приготовления блюда
func (od OrderDish) GetStatus() StatusCook {
	return od.status
}

// GetStatusFailed получение причины провала блюда у главного экземпляра блюда
func (od OrderDish) GetStatusFailed() string {
	return od.statusFailed
}

// GetReadyStatus получение информации готово блюдо или нет
func (od OrderDish) GetReadyStatus() bool {
	return od.isReady
}
