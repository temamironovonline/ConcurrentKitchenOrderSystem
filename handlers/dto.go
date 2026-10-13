package handlers

import (
	"time"

	"main.go/dish"
	"main.go/order"
)

type ListDishesDTO struct {
	List []string `json:"dishes"`
}

type ErrorHTTP struct {
	title string
	t     time.Time
}

func NewErrorHTTP(title string, t time.Time) ErrorHTTP {
	return ErrorHTTP{
		title: title,
		t:     t,
	}
}

type ReturnID struct {
	ID         int       `json:"id"`
	TimeCreate time.Time `json:"time"`
}

func NewCreateOrderAnswerDTO(id int, t time.Time) ReturnID {
	return ReturnID{
		ID:         id,
		TimeCreate: t,
	}
}

type IntermediateDishDTO struct {
	id     int
	name   string
	status string
}

func NewIntermediateDishDTO(od dish.OrderDish) IntermediateDishDTO {
	return IntermediateDishDTO{
		id:     od.GetID(),
		name:   od.GetName(),
		status: string(od.GetStatus()),
	}
}

type OrderDTO struct {
	ID         int             `json:"id"`
	Status     string          `json:"status"`
	ListDishes []dish.Snapshot `json:"dishes"`
}

func newOrderDTO(ord order.Snapshot) OrderDTO {
	return OrderDTO{
		ID:         ord.ID,
		Status:     ord.Status,
		ListDishes: ord.Dishes,
	}
}
