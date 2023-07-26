package counter

import (
	"context"

	"cosmossdk.io/collections"
	v1 "cosmossdk.io/x/accounts/examples/counter/v1"
	"cosmossdk.io/x/accounts/sdk"
)

func NewCounter(sb *collections.SchemaBuilder) Counter {
	return Counter{
		Counter: collections.NewSequence(sb, collections.NewPrefix(0), "counter"),
	}
}

type Counter struct {
	Counter collections.Sequence
}

func (a Counter) Init(ctx context.Context, msg v1.MsgInit) (v1.MsgInitResponse, error) {
	err := a.Counter.Set(ctx, msg.CounterValue)
	if err != nil {
		return v1.MsgInitResponse{}, err
	}
	return v1.MsgInitResponse{}, nil
}

func (a Counter) GetCounterValue(ctx context.Context) (uint64, error) {
	return a.Counter.Peek(ctx)
}

func (a Counter) IncreaseCounterValue(ctx context.Context) (uint64, error) {
	return a.Counter.Next(ctx)
}

func (a Counter) RegisterQueryHandlers(router *sdk.QueryRouter) error {
	err := sdk.RegisterQueryHandler(router, func(ctx context.Context, msg v1.QueryCounterRequest) (v1.QueryCounterResponse, error) {
		value, err := a.GetCounterValue(ctx)
		return v1.QueryCounterResponse{CounterValue: value}, err
	})
	if err != nil {
		return err
	}

	return nil
}

func (a Counter) RegisterExecuteHandlers(router *sdk.ExecuteRouter) error {
	err := sdk.RegisterExecuteHandler(router,
		func(ctx context.Context, msg v1.MsgIncreaseCounter) (v1.MsgIncreaseCounterResponse, error) {
			newValue, err := a.IncreaseCounterValue(ctx)
			return v1.MsgIncreaseCounterResponse{CounterValue: newValue}, err
		})
	return err
}