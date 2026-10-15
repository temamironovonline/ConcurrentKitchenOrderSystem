package handlers

import (
	"time"

	"main.go/dish"
	"main.go/order"
	"main.go/statistics"
)

// ListDishesDTO используется для временного хранения блюд, когда пользователь отправляет список на обработку
type ListDishesDTO struct {
	List []string `json:"dishes"`
}

// ErrorHTTP предназначен для оформления ответа, если получена ошибка во время работы программы
type ErrorHTTP struct {
	Title string    `json:"title"`
	T     time.Time `json:"time"`
}

// NewErrorHTTP используется для создания ErrorHTTP
func NewErrorHTTP(title string, t time.Time) ErrorHTTP {
	return ErrorHTTP{
		Title: title,
		T:     t,
	}
}

// ReturnID используется для возврата ID, когда блюдо было создано
type ReturnID struct {
	ID         int       `json:"id"`
	TimeCreate time.Time `json:"time"`
}

// NewReturnID используется для создания ReturnID
func NewReturnID(id int, t time.Time) ReturnID {
	return ReturnID{
		ID:         id,
		TimeCreate: t,
	}
}

// DishDTO информация о блюде для отправки пользователю
type DishDTO struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	Status       string         `json:"status"`
	TimeCook     *time.Duration `json:"timeCook"`
	StatusFailed string         `json:"statusFailed"`
}

// NewDishDTO создает экземпляр DishDTO
func NewDishDTO(od dish.Snapshot) DishDTO {
	return DishDTO{
		ID:           od.ID,
		Name:         od.Name,
		Status:       string(od.Status),
		TimeCook:     od.CookingTime,
		StatusFailed: od.StatusFailed,
	}
}

// OrderDTO информация о всем блюде для отправки пользователю
type OrderDTO struct {
	ID         int       `json:"id"`
	Status     string    `json:"status"`
	ListDishes []DishDTO `json:"dishes"`
}

// NewOrderDTO создает экземпляр OrderDTO
func NewOrderDTO(ord order.Snapshot) OrderDTO {
	dishes := make([]DishDTO, 0, len(ord.Dishes))
	for _, v := range ord.Dishes {
		dishes = append(dishes, DishDTO{ID: v.ID, Name: v.Name, Status: string(v.Status), TimeCook: v.CookingTime, StatusFailed: v.StatusFailed})
	}

	o := OrderDTO{
		ID:         ord.ID,
		Status:     ord.Status,
		ListDishes: dishes,
	}

	return o
}

// Stats используется для статистики
type Stats struct {
	AllDishCooked   int64 `json:"allDishProccessed"`
	CompletedCooked int64 `json:"completed"`
	FailCooked      int64 `json:"fail"`
	OrderCookingNow int64 `json:"orderCookingNow"`
}

// NewStats для создания экземпляра Stats
func NewStats(s *statistics.Statistics) Stats {
	return Stats{
		AllDishCooked:   s.GetTotalDishesProccesed(),
		CompletedCooked: s.GetCompletedCooked(),
		FailCooked:      s.GetFailedCooked(),
		OrderCookingNow: s.GetOrderCookingNow(),
	}
}
