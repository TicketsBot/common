package rpc

import "context"

type Listener interface {
	BuildContext() (context.Context, context.CancelFunc)
	HandleMessage(ctx context.Context, message []byte)
}
