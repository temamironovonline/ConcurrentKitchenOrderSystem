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

// NewAvailabilityDishes создает новое блюдо
func NewAvailabilityDishes(nameDish string, stat consts.ClassStation, t time.Duration) AvailabilityDish {
	return AvailabilityDish{
		name:     nameDish,
		stat:     stat,
		timeCook: t * time.Second,
	}
}

func (ad AvailabilityDish) GetTimeout() time.Duration {
	return ad.timeCook
}
