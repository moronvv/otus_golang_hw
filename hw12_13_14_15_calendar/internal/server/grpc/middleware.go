package internalgrpc

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	internalcontext "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/context"
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

func UnaryAuthInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	uidHeader := md["user-id"]

	if len(uidHeader) != 1 {
		return nil, status.Error(codes.Unauthenticated, "wrong User-ID header")
	}
	userID, err := uuid.Parse(uidHeader[0])
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return handler(internalcontext.SetUserID(ctx, userID), req)
}
