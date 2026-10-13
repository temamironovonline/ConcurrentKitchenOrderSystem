package station

import (
	"context"
	"math/rand"
	"time"

	"main.go/consts"
	"main.go/dish"
)

// Station представляет собой структуру, которая содержит информация о станции на кухне
type Station struct {
	class    consts.ClassStation
	cookers  chan struct{}
	cookTime func() time.Duration
}

func newStation(class consts.ClassStation, limitWorkers int) *Station {
	return &Station{
		class:    class,
		cookers:  make(chan struct{}, limitWorkers),
		cookTime: generateTime,
	}
}

// NewGrillStation создает экземпляр станции для гриля
func NewGrillStation(countWorkers int) *Station {
	return newStation(consts.GrillClass, countWorkers)
}

// NewFryerStation создает экземпляр станции для фритюра
func NewFryerStation(countWorkers int) *Station {
	return newStation(consts.FryerClass, countWorkers)
}

// NewMolecularStation создает экземпляр станции для молекулярной кухни
func NewMolecularStation(countWorkers int) *Station {
	return newStation(consts.MolecularClass, countWorkers)
}

// NewBarStation создает экзмепляр станции для бара
func NewBarStation(countWorkers int) *Station {
	return newStation(consts.BarClass, countWorkers)
}

// Cook представляет собой сам процесс приготовления блюда
func (s *Station) Cook(ctx context.Context, d dish.OrderDish) dish.UpdateOrderDish {
	// Ждем повара, но не дольше, чем ждем заказ
	select {
	case <-ctx.Done():
		return dish.NewUpdateOrderDish(false, nil, ctx.Err())
	case s.cookers <- struct{}{}:
		defer func() { <-s.cookers }() // освобождаем место после окончания работы повара
	}

	ctxCook, cancel := context.WithTimeout(ctx, d.GetTimeout())
	defer cancel()

	now := time.Now()

	select {
	case <-ctxCook.Done():
		return dish.NewUpdateOrderDish(false, nil, ctxCook.Err())
	case <-time.After(s.cookTime()):
		elapsed := time.Since(now)
		return dish.NewUpdateOrderDish(true, &elapsed, nil)
	}
}

func generateTime() time.Duration {
	return time.Duration(rand.Intn(10)) * time.Second
}
