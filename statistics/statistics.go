package statistics

import "sync/atomic"

type Statistics struct {
	allDishCooked   atomic.Int64
	completedCooked atomic.Int64
	failCooked      atomic.Int64
	orderCookingNow atomic.Int64
}

func NewStatistics() *Statistics {
	return &Statistics{}
}

func (s *Statistics) AddAllDishCooked() {
	s.allDishCooked.Add(1)
}

func (s *Statistics) AddCompletedCooked() {
	s.completedCooked.Add(1)
}

func (s *Statistics) AddFailCooked() {
	s.failCooked.Add(1)
}

func (s *Statistics) AddOrderCookingNow() {
	s.orderCookingNow.Add(1)
}

func (s *Statistics) SubtractionOrderCookingNow() {
	s.orderCookingNow.Add(-1)
}
