package main

import (
	"time"

	"github.com/pttrulez/go-microservices/types"
	"github.com/sirupsen/logrus"
)

// логгирует инфу перед тем как там пойдет в аггрегатор
// миддлвары это просто обертки которые эмулируют вызов функции у оборачиваемой сущности
// Выполняют свою логику, далее вызывают одноименную функцию у оборачиваемой сущности
type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{next: next}
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"err":      err,
			"distance": distance,
		}).Info()
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return err
}

func (m *LogMiddleware) CalcualateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
		)
		if inv != nil {
			distance = inv.TotalDistance
			amount = inv.TotalAmount
		}
		logrus.WithFields(logrus.Fields{
			"took":      time.Since(start),
			"err":       err,
			"obuID":     obuID,
			"totalDist": distance,
			"totalAmt":  amount,
		}).Info()
	}(time.Now())
	inv, err = m.next.CalcualateInvoice(obuID)
	if err != nil {
		return nil, err
	}
	return inv, nil
}
