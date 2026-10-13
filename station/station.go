package station

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"main.go/consts"
	"main.go/dish"
)

var defectCook = []string{"закончился ингридиент", "подгорело"}

// Station представляет собой структуру, которая содержит информация о станции на кухне
type Station struct {
	class    consts.ClassStation
	cookers  chan struct{}
	cookTime func() time.Duration
	defect   func() error
}

func newStation(class consts.ClassStation, limitWorkers int) *Station {
	return &Station{
		class:    class,
		cookers:  make(chan struct{}, limitWorkers),
		cookTime: generateTime,
		defect:   randomDefect,
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

	ctxCook, cancel := context.WithTimeout(ctx, d.GetTimeout())
	defer cancel()

	// Ждем повара, но не дольше, чем ждем заказ
	select {
	case <-ctx.Done():
		return dish.NewUpdateOrderDish(false, nil, ctx.Err())
	case <-ctxCook.Done():
		return dish.NewUpdateOrderDish(false, nil, ctxCook.Err())
	case s.cookers <- struct{}{}:
		defer func() { <-s.cookers }() // освобождаем место после окончания работы повара

		// Если ctx.Done уже сработал, но select выбрал второй кейс
		if err := ctx.Err(); err != nil {
			return dish.NewUpdateOrderDish(false, nil, err)
		}

		if err := ctxCook.Err(); err != nil {
			return dish.NewUpdateOrderDish(false, nil, err)
		}
	}

	now := time.Now()

	select {
	case <-ctxCook.Done():
		return dish.NewUpdateOrderDish(false, nil, ctxCook.Err())
	case <-time.After(s.cookTime()):
		err := s.defect()
		isReady := err == nil
		elapsed := time.Since(now)
		if !isReady {
			elapsed = 0
		}

		return dish.NewUpdateOrderDish(isReady, &elapsed, err)
	}
}

func randomDefect() error {
	if rand.Intn(100) < 20 {
		return errors.New(defectCook[rand.Intn(len(defectCook))])
	}
	return nil
}

func generateTime() time.Duration {
	return time.Duration(rand.Intn(10)+1) * time.Second
}
