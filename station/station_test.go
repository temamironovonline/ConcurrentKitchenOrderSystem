package station

import (
	"context"
	"sync"
	"testing"
	"time"

	"main.go/consts"
	"main.go/dish"
)

func newTestDish() dish.OrderDish {
	return dish.NewOrderDish(dish.NewAvailabilityDishes("test", consts.GrillClass, 5*time.Second))
}

func newTestDishWithTimeout(timeout time.Duration) dish.OrderDish {
	return dish.NewOrderDish(dish.NewAvailabilityDishes("test", consts.GrillClass, timeout))
}

func TestStation_RespectsLimit(t *testing.T) {
	s := newStation(consts.GrillClass, 2)
	s.cookTime = func() time.Duration { return 100 * time.Millisecond }
	d := newTestDish()

	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Cook(context.Background(), d)
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	// 4 блюда по 100 мс при лимите 2 дают две «волны», то есть около 200 мс.
	// Если бы лимит не работал, всё закончилось бы за ~100 мс.
	if elapsed < 190*time.Millisecond {
		t.Errorf("лимит не работает: 4 блюда заняли %v, ожидали не меньше 200ms", elapsed)
	}
	if elapsed > 400*time.Millisecond {
		t.Errorf("слишком долго: %v", elapsed)
	}
}

func TestStation_CancelWhileWaitingForCook(t *testing.T) {
	s := newStation(consts.GrillClass, 1)
	d := newTestDish()
	s.cookTime = func() time.Duration { return 300 * time.Millisecond }

	go s.Cook(context.Background(), d)
	time.Sleep(50 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	s.Cook(ctx, d)
	elapsed := time.Since(start)

	if elapsed > 200*time.Millisecond {
		t.Errorf("блюдо ждало повара %v вместо отмены по контексту", elapsed)
	}
}

func TestStation_DishTimeout(t *testing.T) {
	s := newStation(consts.GrillClass, 1)
	d := newTestDishWithTimeout(50 * time.Millisecond)
	s.cookTime = func() time.Duration { return 500 * time.Millisecond }

	start := time.Now()
	s.Cook(context.Background(), d)
	elapsed := time.Since(start)

	if elapsed > 200*time.Millisecond {
		t.Errorf("cook вернулся через %v, ожидали около 50 мс", elapsed)
	}
}

func TestStation_CookInTime(t *testing.T) {
	s := newStation(consts.GrillClass, 1)
	s.cookTime = func() time.Duration { return 50 * time.Millisecond }
	d := newTestDishWithTimeout(500 * time.Millisecond)

	start := time.Now()
	s.Cook(context.Background(), d)
	elapsed := time.Since(start)

	if elapsed < 40*time.Millisecond {
		t.Errorf("блюдо приготовилось за %v, а ожидалось около 50ms", elapsed)
	}

	if elapsed > 200*time.Millisecond {
		t.Errorf("блюдо готовится больше %v, а ожидалось около 50ms", elapsed)
	}

}

func TestStation_CookResult(t *testing.T) {
	tests := []struct {
		name        string
		cookTime    time.Duration
		dishTimeout time.Duration
		want        dish.StatusCook
	}{
		{"успевает в срок", 50 * time.Millisecond, 500 * time.Millisecond, dish.StatusCooked},
		{"не успевает в срок", 500 * time.Millisecond, 50 * time.Millisecond, dish.StatusTimeout},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newStation(consts.GrillClass, 1)
			s.cookTime = func() time.Duration { return tc.cookTime }
			d := newTestDishWithTimeout(tc.dishTimeout)

			update := s.Cook(context.Background(), d)

			if update.GetStatus() != tc.want {
				t.Errorf("статус: %v, хотели: %v", update.GetStatus(), tc.want)
			}
		})
	}
}
