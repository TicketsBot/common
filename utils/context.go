package utils

import (
	"context"
	"time"
)

func DefaultContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	return ctx
}
