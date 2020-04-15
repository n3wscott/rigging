package main

import (
	"context"
	"fmt"

	"github.com/n3wscott/rigging/pkg/testbed"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/signals"
)

func main() {
	ctx := signals.NewContext()
	cfg := sharedmain.ParseAndGetConfigOrDie()
	ctx, informers := injection.Default.SetupInformers(ctx, cfg)

	// Start the injection clients and informers.
	go func(ctx context.Context) {
		if err := controller.StartInformers(ctx.Done(), informers...); err != nil {
			panic(fmt.Sprintf("Failed to start informers - %s", err))
		}
		<-ctx.Done()
	}(ctx)

	// Start the testbed.
	tb := testbed.New(ctx)
	if err := tb.Start(ctx); err != nil {
		panic(fmt.Sprintf("Failed running testbed - %s", err))
	}
}
