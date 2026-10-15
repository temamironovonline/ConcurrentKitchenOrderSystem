package chef

import (
	"context"
	"sync"
	"time"

	"main.go/consts"
	"main.go/dish"
	"main.go/order"
	orderlist "main.go/order_list"
	"main.go/statistics"
)

// Chef является структурой Шефа
type Chef struct {
	listOrder *orderlist.OrderList
	stations  map[consts.ClassStation]StationCook
	menu      *dish.ListAvailabilityDishes
	statistic *statistics.Statistics
}

// StationCook представляет собой интерфейс для станций (в данном случае grill, bar, molecular, fryer)
type StationCook interface {
	Cook(ctx context.Context, d dish.OrderDish) dish.UpdateOrderDish
}

// NewChef создает Шефа, который является связующим звеном между хендлерами и внутренней логикой программы
func NewChef(stations map[consts.ClassStation]StationCook, menu *dish.ListAvailabilityDishes, statistic *statistics.Statistics) *Chef {
	return &Chef{
		listOrder: orderlist.NewOrderList(),
		stations:  stations,
		menu:      menu,
		statistic: statistic,
	}
}

// CreateOrder предназначен для создания заказа с возвращением ID по созданному заказу, не дожидаясь его завершения
func (c *Chef) CreateOrder(ctx context.Context, timeout time.Duration, dishes []string) (int, error) {
	newOrder, ctxOrder, err := order.NewOrder(ctx, timeout, c.menu, dishes)
	if err != nil {
		return 0, err
	}

	c.listOrder.Add(newOrder)
	c.statistic.AddOrderCookingNow()

	wg := &sync.WaitGroup{}
	for _, v := range newOrder.Dishes() {
		wg.Add(1)

		go func(d dish.OrderDish) {
			defer wg.Done()

			st := c.stations[d.GetStation()]
			update := st.Cook(ctxOrder, d)
			newOrder.UpdateDish(d, update)
		}(v)
	}

	go func() {
		defer c.statistic.DecrementOrderCookingNow()
		wg.Wait()
		newOrder.Finish(ctxOrder.Err())

		c.updateStatistics(newOrder)
	}()

	return newOrder.GetID(), nil
}

func (c *Chef) updateStatistics(od *order.Order) {
	tmp := od.Dishes()

	for _, v := range tmp {
		c.statistic.AddTotalDishesProcessed()
		if v.GetReadyStatus() {
			c.statistic.AddCompletedCooked()
		} else {
			c.statistic.AddFailCooked()
		}
	}
}

// GetOrder используется для получения данных по конкретному заказу
func (c *Chef) GetOrder(id int) (order.Snapshot, error) {
	return c.listOrder.GetOrderSnapshot(id)
}

// CancelOrder используется для отмены заказа
func (c *Chef) CancelOrder(id int) (order.Snapshot, error) {
	ord, err := c.listOrder.GetOrder(id)

	if err != nil {
		return order.Snapshot{}, err
	}
	ord.Cancel()
	return ord.NewSnapshot(), nil
}

// GetStatistics используется для получения промежуточной статистики
func (c *Chef) GetStatistics() *statistics.Statistics {
	return c.statistic
}
