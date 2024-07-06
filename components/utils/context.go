package utils

import "context"

type ContextWrapper struct {
	Context context.Context
	Cancel  context.CancelFunc
}

func (c *ContextWrapper) ResetContext() (context.Context, context.CancelFunc) {
	ctx, cancel := c.Context, c.Cancel
	c.Context, c.Cancel = context.WithCancel(context.Background())
	return ctx, cancel
}

func (c *ContextWrapper) ToCancel(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
