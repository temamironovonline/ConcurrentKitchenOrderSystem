package cooker

/*
import (
	"context"
	"log"
	"sync"

	"main.go/dish"
	dishlist "main.go/dish_list"
)

type Cooker struct {
	s StationCook
	l *dishlist.List
}



func NewCooker(s StationCook, l *dishlist.List) *Cooker {
	return &Cooker{
		s: s,
		l: l,
	}
}

func (c *Cooker) Run(ctx context.Context, wg *sync.WaitGroup, job <-chan int, results chan<- dish.UpdateOrderDish) {
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
			update = update.SetID(id)

			select {
			case <-ctx.Done():
				log.Println("cooker finished")
				return

			case results <- update:
			}

		}
	}
}
*/
