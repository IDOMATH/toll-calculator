package main

import (
	"github.com/idomath/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

const basePrice = 1.15

type Aggregator interface {
	AggregateDistance(distance types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(distance types.Distance) error
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

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	logrus.WithFields(logrus.Fields{
		"obuId":    distance.ObuId,
		"distance": distance.Value,
		"unix":     distance.Unix,
	}).Info("aggregating distance")
	return i.store.Insert(distance)
}

func (i *InvoiceAggregator) CalculateInvoice(id int) (*types.Invoice, error) {
	dist, err := i.store.Get(id)
	if err != nil {
		return nil, err
	}
	invoice := &types.Invoice{
		ObuId:         id,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}
	return invoice, nil
}
