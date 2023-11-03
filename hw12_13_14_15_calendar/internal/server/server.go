package internalserver

import "context"

type Server interface {
	GetType() string
	GetAddress() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
