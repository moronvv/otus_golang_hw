package internalgrpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	internalerrors "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/errors"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/pb"
)

func getErrorStatus(err error) error {
	switch {
	case errors.Is(err, internalerrors.ErrDocumentOperationForbidden):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, internalerrors.ErrDocumentNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func (s *server) CreateEvent(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {
	eventToCreate := models.Event{}
	if err := eventToCreate.UnmarshalPB(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	createdEvent, err := s.app.CreateEvent(ctx, &eventToCreate)
	if err != nil {
		return nil, getErrorStatus(err)
	}

	return createdEvent.MarshalPB(), nil
}

func (s *server) GetEvents(ctx context.Context, _ *emptypb.Empty) (*pb.EventResponses, error) {
	events, err := s.app.GetEvents(ctx)
	if err != nil {
		return nil, getErrorStatus(err)
	}

	resps := []*pb.EventResponse{}
	for _, event := range events {
		resps = append(resps, event.MarshalPB())
	}

	return &pb.EventResponses{
		Events: resps,
	}, nil
}

func (s *server) GetEvent(ctx context.Context, req *pb.EventId) (*pb.EventResponse, error) {
	event, err := s.app.GetEvent(ctx, req.Id)
	if err != nil {
		return nil, getErrorStatus(err)
	}

	return event.MarshalPB(), nil
}

func (s *server) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.EventResponse, error) {
	eventToUpdate := models.Event{}
	if err := eventToUpdate.UnmarshalPB(req.Event); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	updatedEvent, err := s.app.UpdateEvent(ctx, req.Id, &eventToUpdate)
	if err != nil {
		return nil, getErrorStatus(err)
	}

	return updatedEvent.MarshalPB(), nil
}

func (s *server) DeleteEvent(ctx context.Context, req *pb.EventId) (*emptypb.Empty, error) {
	if err := s.app.DeleteEvent(ctx, req.Id); err != nil {
		return nil, getErrorStatus(err)
	}

	return &emptypb.Empty{}, nil
}