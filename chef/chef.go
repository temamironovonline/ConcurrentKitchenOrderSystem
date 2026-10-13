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

type Chef struct {
	listOrder *orderlist.OrderList
	stations  map[consts.ClassStation]StationCook
	menu      *dish.ListAvailabilityDishes
	statistic *statistics.Statistics
}

type StationCook interface {
	Cook(ctx context.Context, d dish.OrderDish) dish.UpdateOrderDish
}

func NewChef(stations map[consts.ClassStation]StationCook, menu *dish.ListAvailabilityDishes, statistic *statistics.Statistics) *Chef {
	return &Chef{
		listOrder: orderlist.NewOrderList(),
		stations:  stations,
		menu:      menu,
		statistic: statistic,
	}
}

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
		defer c.statistic.SubtractionOrderCookingNow()
		wg.Wait()
		newOrder.Finish(ctxOrder.Err())
	}()

	return newOrder.GetID(), nil
}

func (c *Chef) GetOrder(id int) (order.Snapshot, error) {
	return c.listOrder.GetOrderSnapshot(id)
}

func (c *Chef) CancelOrder(id int) (order.Snapshot, error) {
	ord, err := c.listOrder.GetOrder(id)

	if err != nil {
		return order.Snapshot{}, err
	}
	ord.Cancel()
	return ord.Snapshot(), nil
}
