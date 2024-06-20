package client

import (
	"context"

	"github.com/pttrulez/go-microservices/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
}
