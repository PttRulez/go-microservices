package main

import (
	"time"

	"github.com/pttrulez/go-microservices/types"
	"github.com/sirupsen/logrus"
)

// Миддлвары это просто обертки которые эмулируют вызов функции у оборачиваемой 
// сущности. Выполняют свою логику, далее вызывают одноименную функцию у оборачиваемой
// сущности
type LogMiddleware struct {
	next DataProducer
}

// логгирует инфу перед тем как отдать её в дата продюсер. Она просто оборачивает дата
// продюсер, выполняет такую же функцию ProduceData и далее вызывает функицю ProduceData
// у продюсера
func NewLogMiddleware(next DataProducer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) ProduceData(data types.OBUData) error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuID": data.OBUID,
			"lat":   data.Lat,
			"long":  data.Long,
			"took":  time.Since(start),
		}).Info("producing to kafka")
	}(time.Now())

	return l.next.ProduceData(data)
}
