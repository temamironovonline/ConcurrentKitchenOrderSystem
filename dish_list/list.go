package dishlist

import (
	"errors"
	"log"
	"sync"

	"main.go/dish"
)

type List struct {
	mtx  sync.Mutex
	list map[int]dish.OrderDish
}

func NewList() *List {
	return &List{
		list: make(map[int]dish.OrderDish),
	}
}

func (l *List) GetByID(id int) (dish.OrderDish, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	dis, ok := l.list[id]
	if !ok {
		log.Println("блюда в заказе с таким ID не существует")
		return dish.OrderDish{}, errors.New("блюда в заказе с таким ID не существует")
	}

	return dis, nil
}

func (l *List) Add(d dish.OrderDish) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if _, ok := l.list[d.GetID()]; ok {
		log.Println("блюдо с таким ID уже есть в списке")
		return
	}
	l.list[d.GetID()] = d
}

func (l *List) Update(d dish.OrderDish, od dish.UpdateOrderDish) (dish.OrderDish, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if _, ok := l.list[d.GetID()]; !ok {
		log.Println("такого блюда в списке не существует")
		return dish.OrderDish{}, errors.New("такого блюда в списке не существует")
	}

	new := d.Update(od)
	l.list[d.GetID()] = new
	return new, nil
}
