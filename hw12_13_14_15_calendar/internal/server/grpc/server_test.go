package internalgrpc_test

import (
	"context"
	"database/sql"
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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
	mockedapp "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app/mocked"
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

type EventsHandlersSuite struct {
	suite.Suite
	mockedApp  *mockedapp.MockApp
	baseSrv    *grpc.Server
	srv        *bufconn.Listener
	clientConn *grpc.ClientConn
	client     pb.EventServiceClient

	eventData *eventTestData
}

func (s *EventsHandlersSuite) SetupSuite() {
	var err error
	t := s.T()

	s.mockedApp = mockedapp.NewMockApp(t)
	s.baseSrv = grpc.NewServer()
	s.srv = newServer(s.mockedApp, s.baseSrv)
	s.clientConn, err = newClient(s.srv)
	require.NoError(t, err)
	s.client = pb.NewEventServiceClient(s.clientConn)

	s.eventData = newEventTestData()
}

func (s *EventsHandlersSuite) TearDownSuite() {
	s.baseSrv.Stop()

	err := s.clientConn.Close()
	require.NoError(s.T(), err)
}

func (s *EventsHandlersSuite) TestCreateEventHandler() {
	t := s.T()
	ctx := context.Background()

	t.Run("0", func(t *testing.T) {
		s.mockedApp.EXPECT().CreateEvent(mock.Anything, s.eventData.before).Return(s.eventData.after, nil).Once()

		resp, err := s.client.CreateEvent(ctx, s.eventData.req)
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(s.eventData.expectedResp, resp, protocmp.Transform()))

		s.mockedApp.AssertExpectations(t)
	})
}

func TestEventsHandlersSuite(t *testing.T) {
	suite.Run(t, new(EventsHandlersSuite))
}
