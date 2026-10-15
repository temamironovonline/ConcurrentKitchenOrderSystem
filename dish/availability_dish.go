package dish

import (
	"time"

	"main.go/consts"
)

// AvailabilityDish предназначен для описания имеющегося блюда в меню
type AvailabilityDish struct {
	name     string
	stat     consts.ClassStation
	timeCook time.Duration
}

// NewAvailabilityDishes создает новое блюдо для меню
func NewAvailabilityDishes(nameDish string, stat consts.ClassStation, t time.Duration) AvailabilityDish {
	return AvailabilityDish{
		name:     nameDish,
		stat:     stat,
		timeCook: t,
	}
}

// GetTimeout предназначен для получения времени за сколько должно приготовиться блюдо
func (ad AvailabilityDish) GetTimeout() time.Duration {
	return ad.timeCook
}
