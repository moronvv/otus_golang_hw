package internalgrpc_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
	mockedapp "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app/mocked"
	internalerrors "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/errors"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/pb"
	internalgrpc "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/server/grpc"
)

const buffSize = 1024 * 1024

func newServer(app app.App, baseSrv *grpc.Server) *bufconn.Listener {
	internalgrpc.NewServer(nil, app, nil, baseSrv)

	lis := bufconn.Listen(buffSize)
	go func() {
		if err := baseSrv.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	return lis
}

func newClient(lis *bufconn.Listener) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(context.Background(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

type eventTestData struct {
	req          *pb.EventRequest
	incorrectReq *pb.EventRequest
	expectedResp *pb.EventResponse
	before       *models.Event
	after        *models.Event
}

func newEventTestData() *eventTestData {
	dt := time.Now().UTC()
	duration, _ := time.ParseDuration("1m")
	emptyDuration, _ := time.ParseDuration("0s")
	userID := uuid.New()

	return &eventTestData{
		req: &pb.EventRequest{
			Title:       "test",
			Description: "description",
			Datetime:    timestamppb.New(dt),
			Duration:    durationpb.New(duration),
			UserId:      userID.String(),
		},
		incorrectReq: &pb.EventRequest{
			Title:       "t",
			Description: "desc",
			UserId:      "incorrect uuid",
		},
		expectedResp: &pb.EventResponse{
			Id:           1,
			Title:        "test",
			Description:  "description",
			Datetime:     timestamppb.New(dt),
			Duration:     durationpb.New(duration),
			UserId:       userID.String(),
			NotifyBefore: durationpb.New(emptyDuration),
		},
		before: &models.Event{
			Title: "test",
			Description: sql.NullString{
				String: "description",
				Valid:  true,
			},
			Datetime: dt,
			Duration: duration,
			UserID:   userID,
		},
		after: &models.Event{
			ID:    1,
			Title: "test",
			Description: sql.NullString{
				String: "description",
				Valid:  true,
			},
			Datetime: dt,
			Duration: duration,
			UserID:   userID,
		},
	}
}

type EventHandlersSuite struct {
	suite.Suite
	mockedApp  *mockedapp.MockApp
	baseSrv    *grpc.Server
	srv        *bufconn.Listener
	clientConn *grpc.ClientConn
	client     pb.EventServiceClient

	eventData *eventTestData
	mdCtx     context.Context
}

func (s *EventHandlersSuite) SetupSuite() {
	var err error
	t := s.T()

	s.mockedApp = mockedapp.NewMockApp(t)
	s.baseSrv = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			internalgrpc.UnaryAuthInterceptor,
		),
	)
	s.srv = newServer(s.mockedApp, s.baseSrv)
	s.clientConn, err = newClient(s.srv)
	require.NoError(t, err)
	s.client = pb.NewEventServiceClient(s.clientConn)

	s.eventData = newEventTestData()
	s.mdCtx = metadata.NewOutgoingContext(
		context.Background(),
		metadata.New(map[string]string{"user-id": s.eventData.req.UserId}),
	)
}

func (s *EventHandlersSuite) TearDownSuite() {
	s.baseSrv.Stop()

	err := s.clientConn.Close()
	require.NoError(s.T(), err)
}

func (s *EventHandlersSuite) TestCreateEventHandler() {
	t := s.T()

	t.Run("OK", func(t *testing.T) {
		s.mockedApp.EXPECT().CreateEvent(mock.Anything, s.eventData.before).Return(s.eventData.after, nil).Once()

		resp, err := s.client.CreateEvent(s.mdCtx, s.eventData.req)
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(s.eventData.expectedResp, resp, protocmp.Transform()))

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("UNAUTHENTICATED", func(t *testing.T) {
		resp, err := s.client.CreateEvent(context.Background(), s.eventData.req)
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.Unauthenticated, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("INVALID_ARGUMENT", func(t *testing.T) {
		resp, err := s.client.CreateEvent(s.mdCtx, s.eventData.incorrectReq)
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.InvalidArgument, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("INTERNAL", func(t *testing.T) {
		s.mockedApp.EXPECT().CreateEvent(mock.Anything, s.eventData.before).Return(nil, fmt.Errorf("test")).Once()

		resp, err := s.client.CreateEvent(s.mdCtx, s.eventData.req)
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.Internal, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})
}

func (s *EventHandlersSuite) TestGetEventsHandler() {
	t := s.T()

	t.Run("OK", func(t *testing.T) {
		s.mockedApp.EXPECT().GetEvents(mock.Anything).Return([]models.Event{*s.eventData.after}, nil).Once()

		resp, err := s.client.GetEvents(s.mdCtx, &emptypb.Empty{})
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(&pb.EventResponses{
			Events: []*pb.EventResponse{s.eventData.expectedResp},
		}, resp, protocmp.Transform()))

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("UNAUTHENTICATED", func(t *testing.T) {
		resp, err := s.client.GetEvents(context.Background(), &emptypb.Empty{})
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.Unauthenticated, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("INTERNAL", func(t *testing.T) {
		s.mockedApp.EXPECT().GetEvents(mock.Anything).Return(nil, fmt.Errorf("test")).Once()

		resp, err := s.client.GetEvents(s.mdCtx, &emptypb.Empty{})
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.Internal, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})
}

func (s *EventHandlersSuite) TestGetEventHandler() {
	t := s.T()

	t.Run("OK", func(t *testing.T) {
		s.mockedApp.EXPECT().GetEvent(mock.Anything, int64(1)).Return(s.eventData.after, nil).Once()

		resp, err := s.client.GetEvent(s.mdCtx, &pb.EventId{Id: 1})
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(s.eventData.expectedResp, resp, protocmp.Transform()))

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("UNAUTHENTICATED", func(t *testing.T) {
		resp, err := s.client.GetEvent(context.Background(), &pb.EventId{Id: 1})
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.Unauthenticated, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("PERMISSION_DENIED", func(t *testing.T) {
		s.mockedApp.EXPECT().
			GetEvent(mock.Anything, int64(1)).
			Return(nil, internalerrors.ErrDocumentOperationForbidden).
			Once()

		resp, err := s.client.GetEvent(s.mdCtx, &pb.EventId{Id: 1})
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.PermissionDenied, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("NOT_FOUND", func(t *testing.T) {
		s.mockedApp.EXPECT().GetEvent(mock.Anything, int64(1)).Return(nil, internalerrors.ErrDocumentNotFound).Once()

		resp, err := s.client.GetEvent(s.mdCtx, &pb.EventId{Id: 1})
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.NotFound, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("INTERNAL", func(t *testing.T) {
		s.mockedApp.EXPECT().GetEvent(mock.Anything, int64(1)).Return(nil, fmt.Errorf("test")).Once()

		resp, err := s.client.GetEvent(s.mdCtx, &pb.EventId{Id: 1})
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.Internal, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})
}

func (s *EventHandlersSuite) TestUpdateEventHandler() {
	t := s.T()

	t.Run("OK", func(t *testing.T) {
		s.mockedApp.EXPECT().UpdateEvent(mock.Anything, int64(1), s.eventData.before).Return(s.eventData.after, nil).Once()

		resp, err := s.client.UpdateEvent(s.mdCtx, &pb.UpdateEventRequest{
			Id:    1,
			Event: s.eventData.req,
		})
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(s.eventData.expectedResp, resp, protocmp.Transform()))

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("UNAUTHENTICATED", func(t *testing.T) {
		resp, err := s.client.UpdateEvent(context.Background(), &pb.UpdateEventRequest{
			Id:    1,
			Event: s.eventData.req,
		})
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.Unauthenticated, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("INVALID_ARGUMENT", func(t *testing.T) {
		resp, err := s.client.UpdateEvent(s.mdCtx, &pb.UpdateEventRequest{
			Id:    1,
			Event: s.eventData.incorrectReq,
		})
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.InvalidArgument, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("PERMISSION_DENIED", func(t *testing.T) {
		s.mockedApp.EXPECT().
			UpdateEvent(mock.Anything, int64(1), s.eventData.before).
			Return(nil, internalerrors.ErrDocumentOperationForbidden).
			Once()

		resp, err := s.client.UpdateEvent(s.mdCtx, &pb.UpdateEventRequest{
			Id:    1,
			Event: s.eventData.req,
		})
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.PermissionDenied, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("NOT_FOUND", func(t *testing.T) {
		s.mockedApp.EXPECT().
			UpdateEvent(mock.Anything, int64(1), s.eventData.before).
			Return(nil, internalerrors.ErrDocumentNotFound).
			Once()

		resp, err := s.client.UpdateEvent(s.mdCtx, &pb.UpdateEventRequest{
			Id:    1,
			Event: s.eventData.req,
		})
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.NotFound, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("INTERNAL", func(t *testing.T) {
		s.mockedApp.EXPECT().UpdateEvent(mock.Anything, int64(1), s.eventData.before).Return(nil, fmt.Errorf("test")).Once()

		resp, err := s.client.UpdateEvent(s.mdCtx, &pb.UpdateEventRequest{
			Id:    1,
			Event: s.eventData.req,
		})
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.Internal, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})
}

func (s *EventHandlersSuite) TestDeleteEventHandler() {
	t := s.T()

	t.Run("OK", func(t *testing.T) {
		s.mockedApp.EXPECT().DeleteEvent(mock.Anything, int64(1)).Return(nil).Once()

		resp, err := s.client.DeleteEvent(s.mdCtx, &pb.EventId{Id: 1})
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(&emptypb.Empty{}, resp, protocmp.Transform()))

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("UNAUTHENTICATED", func(t *testing.T) {
		resp, err := s.client.DeleteEvent(context.Background(), &pb.EventId{Id: 1})
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.Unauthenticated, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("PERMISSION_DENIED", func(t *testing.T) {
		s.mockedApp.EXPECT().DeleteEvent(mock.Anything, int64(1)).Return(internalerrors.ErrDocumentOperationForbidden).Once()

		resp, err := s.client.DeleteEvent(s.mdCtx, &pb.EventId{Id: 1})
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.PermissionDenied, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("NOT_FOUND", func(t *testing.T) {
		s.mockedApp.EXPECT().DeleteEvent(mock.Anything, int64(1)).Return(internalerrors.ErrDocumentNotFound).Once()

		resp, err := s.client.DeleteEvent(s.mdCtx, &pb.EventId{Id: 1})
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.NotFound, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("INTERNAL", func(t *testing.T) {
		s.mockedApp.EXPECT().DeleteEvent(mock.Anything, int64(1)).Return(fmt.Errorf("test")).Once()

		resp, err := s.client.DeleteEvent(s.mdCtx, &pb.EventId{Id: 1})
		grpcErr, ok := status.FromError(err)
		require.True(t, ok)
		require.Equalf(t, codes.Internal, grpcErr.Code(), grpcErr.Message())
		require.Nil(t, resp)

		s.mockedApp.AssertExpectations(t)
	})
}

func TestEventHandlersSuite(t *testing.T) {
	suite.Run(t, new(EventHandlersSuite))
}
