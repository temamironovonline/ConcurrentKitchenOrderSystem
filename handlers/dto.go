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
	Title string    `json:"title"`
	T     time.Time `json:"time"`
}

func NewErrorHTTP(title string, t time.Time) ErrorHTTP {
	return ErrorHTTP{
		Title: title,
		T:     t,
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

type DishDTO struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	Status       string         `json:"status"`
	TimeCook     *time.Duration `json:"timeCook"`
	StatusFailed string         `json:"statusFailed"`
}

func NewIntermediateDishDTO(od dish.Snapshot) DishDTO {
	return DishDTO{
		ID:           od.ID,
		Name:         od.Name,
		Status:       string(od.Status),
		TimeCook:     od.CookingTime,
		StatusFailed: od.StatusFailed,
	}
}

type OrderDTO struct {
	ID         int       `json:"id"`
	Status     string    `json:"status"`
	ListDishes []DishDTO `json:"dishes"`
}

func newOrderDTO(ord order.Snapshot) OrderDTO {
	dishes := make([]DishDTO, 0, len(ord.Dishes))
	for _, v := range ord.Dishes {
		dishes = append(dishes, DishDTO{ID: v.ID, Name: v.Name, Status: string(v.Status)})
	}

	o := OrderDTO{
		ID:         ord.ID,
		Status:     ord.Status,
		ListDishes: dishes,
	}

	return o
}
