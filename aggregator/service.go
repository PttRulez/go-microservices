package main

import (
	"github.com/pttrulez/go-microservices/types"
)

const basePrice = 3.15

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalcualateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

// Сохраняет дистанцию в хранилище
func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	return i.store.Insert(distance)
}

// Рассчитывает инвойс для оплаты дистанции, к-ю проехал автомобиль с OBU на борту
func (i *InvoiceAggregator) CalcualateInvoice(obuid int) (*types.Invoice, error) {
	dist, err := i.store.Get(obuid)
	if err != nil {
		return nil, err
	}
	inv := &types.Invoice{
		OBUID:         obuid,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}
	return inv, nil
}
