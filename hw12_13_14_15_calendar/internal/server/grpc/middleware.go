package internalgrpc

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type loggerInterceptor struct {
	logger *slog.Logger
}

func (i *loggerInterceptor) UnaryLoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	p, _ := peer.FromContext(ctx)
	md, _ := metadata.FromIncomingContext(ctx)

	start := time.Now()
	resp, err := handler(ctx, req)
	st, _ := status.FromError(err)

	i.logger.Info("",
		slog.Group("grpc",
			"client_ip", p.Addr.String(),
			"user_agent", md.Get("User-Agent"),
			"protocol", "HTTP/2",
			"status", st.Code(),
			"method", "POST",
			"path", info.FullMethod,
		),
		"latency", time.Since(start),
	)

	return resp, err
}

func newLoggerInterceptor(logger *slog.Logger) *loggerInterceptor {
	return &loggerInterceptor{
		logger: logger,
	}
}
