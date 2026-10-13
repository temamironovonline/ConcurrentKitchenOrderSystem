package chef

import (
	"context"
	"sync"
	"time"

	"main.go/consts"
	"main.go/dish"
	"main.go/order"
	orderlist "main.go/order_list"
	"main.go/station"
)

type Chef struct {
	mtx       sync.Mutex
	listOrder *orderlist.OrderList
	stations  map[consts.ClassStation]StationCook
}

type StationCook interface {
	Cook(ctx context.Context, d dish.OrderDish) dish.UpdateOrderDish
}

func NewChef() *Chef {
	c := &Chef{
		listOrder: orderlist.NewOrderList(),
		stations:  make(map[consts.ClassStation]StationCook),
	}

	c.stations[consts.GrillClass] = station.NewGrillStation(3)
	c.stations[consts.BarClass] = station.NewBarStation(3)
	c.stations[consts.FryerClass] = station.NewFryerStation(3)
	c.stations[consts.MolecularClass] = station.NewMolecularStation(3)
	return c
}

func (c *Chef) CreateOrder(ctx context.Context, timeout time.Duration, dishes []string) int {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	newOrder, ctxOrder := order.NewOrder(ctx, timeout, dishes)

	c.listOrder.Add(newOrder)

	wg := &sync.WaitGroup{}
	for _, v := range newOrder.Dishes() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			st := c.stations[v.GetStation()]
			new := st.Cook(ctxOrder, v)
			newOrder.UpdateDish(v, new)
		}()
	}

	go func() {
		wg.Wait()
		newOrder.Finish()
	}()

	return newOrder.GetID()
}

func (c *Chef) GetOrder(id int) (order.Snapshot, error) {
	c.listOrder.GetCertainOrder(id)
}
