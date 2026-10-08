package dish

import (
	"time"

	"main.go/consts"
)

type ListAvailabilityDishes struct {
	list map[string]AvailabilityDish
}

func NewListAvailabilityDishes() *ListAvailabilityDishes {
	return &ListAvailabilityDishes{
		list: make(map[string]AvailabilityDish),
	}
}

func (l *ListAvailabilityDishes) Add() {
	dis := NewAvailabilityDishes("burger", consts.GrillClass, 3*time.Second)
	l.list[dis.name] = dis

	dis = NewAvailabilityDishes("fried egg", consts.GrillClass, 2*time.Second)
	l.list[dis.name] = dis

	dis = NewAvailabilityDishes("bacon", consts.GrillClass, 1*time.Second)
	l.list[dis.name] = dis

	dis = NewAvailabilityDishes("fries", consts.FryerClass, 1*time.Second)
	l.list[dis.name] = dis

	dis = NewAvailabilityDishes("nuggets", consts.FryerClass, 2*time.Second)
	l.list[dis.name] = dis

	dis = NewAvailabilityDishes("strips", consts.FryerClass, 3*time.Second)
	l.list[dis.name] = dis

	dis = NewAvailabilityDishes("chocolate muss", consts.MolecularClass, 5*time.Second)
	l.list[dis.name] = dis

	dis = NewAvailabilityDishes("strawberry muss", consts.MolecularClass, 3*time.Second)
	l.list[dis.name] = dis

	dis = NewAvailabilityDishes("blueberry muss", consts.MolecularClass, 3*time.Second)
	l.list[dis.name] = dis

	dis = NewAvailabilityDishes("whiskey", consts.BarClass, 1*time.Second)
	l.list[dis.name] = dis

	dis = NewAvailabilityDishes("cognac", consts.BarClass, 1*time.Second)
	l.list[dis.name] = dis

	dis = NewAvailabilityDishes("milkshake", consts.BarClass, 1*time.Second)
	l.list[dis.name] = dis
}
