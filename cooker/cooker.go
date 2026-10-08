package cooker

import (
	"context"
	"log"
	"sync"

	"main.go/dish"
	dishlist "main.go/dish_list"
	"main.go/station"
)

type Cooker struct {
	s *station.Station
	l *dishlist.List
}

func (c *Cooker) Run(ctx context.Context, wg *sync.WaitGroup, job <-chan int, results chan<- dish.OrderDish) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Println("cooker finished")
			return
		case id, ok := <-job:
			if !ok {
				return
			}
			dis, err := c.l.GetByID(id)
			if err != nil {
				return
			}

			update := c.s.Cook(ctx, dis)

			res, err := c.l.Update(dis, update)
			if err != nil {
				return
			}

			select {
			case <-ctx.Done():
				log.Println("cooker finished")
				return

			case results <- res:
			}

		}
	}
}
