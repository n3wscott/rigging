package testbed

import (
	"context"
	"k8s.io/client-go/dynamic"
)

type testBed struct {
	client dynamic.Interface
}

func (b *testBed) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}
