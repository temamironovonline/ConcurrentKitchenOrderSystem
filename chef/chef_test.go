package chef

import (
	"context"
	"sync"
	"testing"
	"time"

	"main.go/consts"
	"main.go/dish"
	"main.go/order"
	"main.go/station"
	"main.go/statistics"
)

func init() {}

// TestChef_CorrectingIDConcurrentCreate проверяет ID конкурентно созданных заказов
func TestChef_CorrectingIDConcurrentCreate(t *testing.T) {
	stations := make(map[consts.ClassStation]StationCook)

	stations[consts.BarClass] = station.NewBarStation(3)
	stations[consts.FryerClass] = station.NewFryerStation(3)
	stations[consts.GrillClass] = station.NewGrillStation(3)
	stations[consts.MolecularClass] = station.NewMolecularStation(3)

	menu := dish.NewListAvailabilityDishes()
	stat := statistics.NewStatistics()
	c := NewChef(stations, menu, stat)

	dishes := []string{"burger", "burger", "burger", "burger", "burger"}

	type result struct {
		id  int
		err error
	}

	countOrders := 100
	results := make(chan result, countOrders)

	wg := sync.WaitGroup{}
	for i := 0; i < countOrders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := c.CreateOrder(context.Background(), 10*time.Second, dishes)
			t.Logf("Создан заказ %v", id)
			results <- result{id, err}
		}()
	}

	wg.Wait()
	close(results)

	ids := make(map[int]struct{})

	for res := range results {
		if res.err != nil {
			t.Errorf("CreateOrder() вернул ошибку = %v", res.err)
			continue
		}

		if _, exist := ids[res.id]; exist {
			t.Errorf("id: %v уже существует", res.id)
		}

		ids[res.id] = struct{}{}
	}

	if len(ids) != countOrders {
		t.Errorf("количество записей %v, хотя должно быть 100", len(ids))
	}

	for i := 1; i <= countOrders; i++ {
		if _, exist := ids[i]; !exist {
			t.Errorf("id %v отсутствует, хотя должен быть", i)
		}
	}
}

// TestChef_ConcurrentCancelOrder проверяет конкурентную отмену заказов
func TestChef_ConcurrentCancelOrder(t *testing.T) {
	stations := make(map[consts.ClassStation]StationCook)

	stations[consts.BarClass] = station.NewBarStation(3)
	stations[consts.FryerClass] = station.NewFryerStation(3)
	stations[consts.GrillClass] = station.NewGrillStation(3)
	stations[consts.MolecularClass] = station.NewMolecularStation(3)

	menu := dish.NewListAvailabilityDishes()
	stat := statistics.NewStatistics()
	c := NewChef(stations, menu, stat)

	dishes := []string{"strips", "strips", "strips", "strips", "strips"}

	type result struct {
		id  int
		err error
	}

	var wg sync.WaitGroup
	countOrders := 10
	results := make(chan result, countOrders)

	for i := 0; i < countOrders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := c.CreateOrder(context.Background(), 10*time.Second, dishes)
			t.Logf("создан заказ %v", id)
			results <- result{id, err}
		}()
	}

	wg.Wait()
	close(results)
	ids := make(map[int]struct{})

	for res := range results {
		if res.err != nil {
			t.Fatalf("CreateOrder() выдал ошибку %v", res.err)
		}
		ids[res.id] = struct{}{}
	}

	if len(ids) != countOrders {
		t.Fatalf("получено %d уникальных id, а ожидалось %d", len(ids), countOrders)
	}

	type ResultSnap struct {
		snap order.Snapshot
		err  error
	}

	resultSnap := make(chan ResultSnap, countOrders)

	for i := range ids {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			snap, err := c.CancelOrder(i)
			resultSnap <- ResultSnap{snap: snap, err: err}
		}(i)
	}

	wg.Wait()
	close(resultSnap)

	for v := range resultSnap {
		if v.err != nil {
			t.Errorf("CancelOrder выдал ошибку %v", v.err)
		}

		if v.snap.Status != string(order.StatusCanceled) {
			t.Errorf("Заказ %v имеет статус %v, хотя должен %v", v.snap.ID, v.snap.Status, order.StatusCanceled)
		}
	}
}
