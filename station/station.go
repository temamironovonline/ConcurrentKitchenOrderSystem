package station

import (
	"context"
	"math/rand"
	"sync/atomic"
	"time"

	"main.go/consts"
	"main.go/dish"
)

// Station представляет собой структуру, которая содержит информация о станции на кухне
type Station struct {
	class             consts.ClassStation
	limitCountWorkers atomic.Int64
}

func newStation(class consts.ClassStation, limitWorkers int) *Station {
	station := &Station{
		class: class,
	}
	station.limitCountWorkers.Add(int64(limitWorkers))

	return station
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

// SetLimitCountWorkers позволяет изменять лимит количества работников
func (s *Station) SetLimitCountWorkers(count int) {
	s.limitCountWorkers.Add(int64(count))
}

// Cook представляет собой сам процесс приготовления блюда
func (s *Station) Cook(ctx context.Context, d dish.OrderDish) dish.UpdateOrderDish {
	ctxCook, cancel := context.WithTimeout(ctx, d.GetTimeout())
	defer cancel()

	now := time.Now()

	select {
	case <-ctxCook.Done():
		return dish.NewUpdateOrderDish(false, nil, ctxCook.Err())
	case <-time.After(generateTime()):
		elapsed := time.Since(now)
		return dish.NewUpdateOrderDish(true, &elapsed, nil)
	}

}

func generateTime() time.Duration {
	return time.Duration(rand.Intn(10)) * time.Second
}
