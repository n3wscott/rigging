package main

import (
	"fmt"

	"github.com/n3wscott/rigging/pkg/testbed"
	"knative.dev/pkg/injection/sharedmain"
)

func main() {
	ctx := sharedmain.EnableInjectionOrDie(nil, nil)

	// Start the testbed.
	tb := testbed.New(ctx)
	if err := tb.Start(ctx); err != nil {
		panic(fmt.Sprintf("Failed running testbed - %s", err))
	}
}
