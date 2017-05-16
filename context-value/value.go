package context_value

import "context"

type traceId struct{}

func ParseTraceId(ctx context.Context) string {
	if value, ok := ctx.Value(traceId{}).(string); ok {
		return value
	}
	return ""
}

func WithTraceId(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, traceId{}, value)
}
