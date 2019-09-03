package main

import (
	"context"
	"fmt"

	"go.opencensus.io/trace"
)

func StartSpan(ctx context.Context, name string) (context.Context, *trace.Span) {
	return trace.StartSpan(ctx, fmt.Sprintf("/%s", name))
}
