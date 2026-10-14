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

func (s *Statistics) AddTotalDishesProcessed() {
	s.allDishCooked.Add(1)
}

func (s *Statistics) GetTotalDishesProccesed() int64 {
	return s.allDishCooked.Load()
}

func (s *Statistics) AddCompletedCooked() {
	s.completedCooked.Add(1)
}

func (s *Statistics) GetCompletedCooked() int64 {
	return s.completedCooked.Load()
}

func (s *Statistics) AddFailCooked() {
	s.failCooked.Add(1)
}

func (s *Statistics) GetFailedCooked() int64 {
	return s.failCooked.Load()
}

func (s *Statistics) AddOrderCookingNow() {
	s.orderCookingNow.Add(1)
}

func (s *Statistics) GetOrderCookingNow() int64 {
	return s.orderCookingNow.Load()
}

func (s *Statistics) DecrementOrderCookingNow() {
	s.orderCookingNow.Add(-1)
}
