package dish

import (
	"sync"
	"time"

	"main.go/consts"
)

// ListAvailabilityDishes представляет собой список доступных блюд (меню)
type ListAvailabilityDishes struct {
	mtx  sync.Mutex
	list map[string]AvailabilityDish
}

// NewListAvailabilityDishes конструктор для создания, который заполняет блюда
func NewListAvailabilityDishes() *ListAvailabilityDishes {
	newList := &ListAvailabilityDishes{
		list: make(map[string]AvailabilityDish),
	}

	dis := NewAvailabilityDishes("burger", consts.GrillClass, 3*time.Second)
	newList.list[dis.name] = dis

	dis = NewAvailabilityDishes("fried egg", consts.GrillClass, 2*time.Second)
	newList.list[dis.name] = dis

	dis = NewAvailabilityDishes("bacon", consts.GrillClass, 1*time.Second)
	newList.list[dis.name] = dis

	dis = NewAvailabilityDishes("fries", consts.FryerClass, 1*time.Second)
	newList.list[dis.name] = dis

	dis = NewAvailabilityDishes("nuggets", consts.FryerClass, 2*time.Second)
	newList.list[dis.name] = dis

	dis = NewAvailabilityDishes("strips", consts.FryerClass, 3*time.Second)
	newList.list[dis.name] = dis

	dis = NewAvailabilityDishes("chocolate muss", consts.MolecularClass, 5*time.Second)
	newList.list[dis.name] = dis

	dis = NewAvailabilityDishes("strawberry muss", consts.MolecularClass, 3*time.Second)
	newList.list[dis.name] = dis

	dis = NewAvailabilityDishes("blueberry muss", consts.MolecularClass, 3*time.Second)
	newList.list[dis.name] = dis

	dis = NewAvailabilityDishes("whiskey", consts.BarClass, 1*time.Second)
	newList.list[dis.name] = dis

	dis = NewAvailabilityDishes("cognac", consts.BarClass, 1*time.Second)
	newList.list[dis.name] = dis

	dis = NewAvailabilityDishes("milkshake", consts.BarClass, 1*time.Second)
	newList.list[dis.name] = dis

	return newList
}

// Get позволяет получить конкретное блюдо по названию
func (ad *ListAvailabilityDishes) Get(name string) (AvailabilityDish, bool) {
	ad.mtx.Lock()
	defer ad.mtx.Unlock()
	newAD, ok := ad.list[name]
	if !ok {
		return AvailabilityDish{}, false
	}

	return newAD, true
}
