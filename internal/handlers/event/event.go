package event

import (
	"context"
	"errors"

	"github.com/Estriper0/EventService/internal/models"
	"github.com/Estriper0/EventService/internal/service"
	pb "github.com/Estriper0/protobuf/gen/event"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventGRPCService struct {
	pb.UnimplementedEventServer
	eventService service.IEventService
	validate     *validator.Validate
}

func Register(gRPC *grpc.Server, eventService service.IEventService) {
	pb.RegisterEventServer(gRPC, &EventGRPCService{eventService: eventService, validate: validator.New()})
}

func (s *EventGRPCService) GetAll(
	ctx context.Context,
	req *pb.EmptyRequest,
) (*pb.GetAllResponse, error) {
	events, err := s.eventService.GetAll(ctx)
	if errors.Is(err, service.ErrRepositoryError) {
		return nil, status.Error(codes.Internal, "internal error")
	}
	response := &pb.GetAllResponse{
		Events: []*pb.EventElem{},
	}
	for _, event := range events {
		pb_event := &pb.EventElem{
			Id:                int64(event.Id),
			Title:             event.Title,
			About:             event.About,
			StartDate:         timestamppb.New(event.StartDate),
			Location:          event.Location,
			Status:            string(event.Status),
			MaxAttendees:      int32(event.MaxAttendees),
			CurrentAttendance: int32(event.CurrentAttendance),
			Creator:           event.Creator,
		}
		response.Events = append(response.Events, pb_event)
	}
	return response, nil
}

func (s *EventGRPCService) GetById(
	ctx context.Context,
	req *pb.GetByIdRequest,
) (*pb.GetByIdResponse, error) {
	event, err := s.eventService.GetById(ctx, int(req.Id))
	if err != nil {
		if errors.Is(err, service.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.GetByIdResponse{
		Id:                int64(event.Id),
		Title:             event.Title,
		About:             event.About,
		StartDate:         timestamppb.New(event.StartDate),
		Location:          event.Location,
		Status:            string(event.Status),
		MaxAttendees:      int32(event.MaxAttendees),
		CurrentAttendance: int32(event.CurrentAttendance),
		Creator:           event.Creator,
	}, nil
}

func (s *EventGRPCService) Create(
	ctx context.Context,
	req *pb.CreateRequest,
) (*pb.CreateResponse, error) {
	event := &models.EventCreateRequest{
		Title:        req.Title,
		About:        req.About,
		StartDate:    req.StartDate.AsTime(),
		Location:     req.Location,
		Status:       req.Status,
		MaxAttendees: int(req.MaxAttendees),
		Creator:      req.Creator,
	}

	err := s.validate.Struct(event)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := s.eventService.Create(ctx, event)
	if errors.Is(err, service.ErrRepositoryError) {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.CreateResponse{
		Id: int64(id),
	}, nil
}

func (s *EventGRPCService) DeleteById(
	ctx context.Context,
	req *pb.DeleteByIdRequest,
) (*pb.DeleteByIdResponse, error) {
	err := s.eventService.DeleteById(ctx, int(req.Id))
	if err != nil {
		if errors.Is(err, service.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.DeleteByIdResponse{
		Id: req.Id,
	}, nil
}

func (s *EventGRPCService) Update(
	ctx context.Context,
	req *pb.UpdateRequest,
) (*pb.EmptyResponse, error) {
	event_update := &models.EventUpdateRequest{
		Id:           int(req.Id),
		Title:        req.Title,
		About:        req.About,
		StartDate:    req.StartDate.AsTime(),
		Location:     req.Location,
		Status:       req.Status,
		MaxAttendees: int(req.MaxAttendees),
	}
	if err := s.validate.Struct(event_update); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := s.eventService.Update(ctx, event_update)
	if err != nil {
		if errors.Is(err, service.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.EmptyResponse{}, nil
}

func (s *EventGRPCService) GetAllByCreator(
	ctx context.Context,
	req *pb.GetAllByCreatorRequest,
) (*pb.GetAllResponse, error) {
	if err := s.validate.Var(req.Creator, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	events, err := s.eventService.GetAllByCreator(ctx, req.Creator)
	if errors.Is(err, service.ErrRepositoryError) {
		return nil, status.Error(codes.Internal, "internal error")
	}
	response := &pb.GetAllResponse{
		Events: []*pb.EventElem{},
	}
	for _, event := range events {
		pb_event := &pb.EventElem{
			Id:                int64(event.Id),
			Title:             event.Title,
			About:             event.About,
			StartDate:         timestamppb.New(event.StartDate),
			Location:          event.Location,
			Status:            string(event.Status),
			MaxAttendees:      int32(event.MaxAttendees),
			CurrentAttendance: int32(event.CurrentAttendance),
			Creator:           event.Creator,
		}
		response.Events = append(response.Events, pb_event)
	}
	return response, nil
}

func (s *EventGRPCService) GetAllByStatus(
	ctx context.Context,
	req *pb.GetAllByStatusRequest,
) (*pb.GetAllResponse, error) {
	if err := s.validate.Var(req.Status, "required,oneof=draft published ongoing completed cancelled postponed"); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	events, err := s.eventService.GetAllByStatus(ctx, req.Status)
	if errors.Is(err, service.ErrRepositoryError) {
		return nil, status.Error(codes.Internal, "internal error")
	}
	response := &pb.GetAllResponse{
		Events: []*pb.EventElem{},
	}
	for _, event := range events {
		pb_event := &pb.EventElem{
			Id:                int64(event.Id),
			Title:             event.Title,
			About:             event.About,
			StartDate:         timestamppb.New(event.StartDate),
			Location:          event.Location,
			Status:            string(event.Status),
			MaxAttendees:      int32(event.MaxAttendees),
			CurrentAttendance: int32(event.CurrentAttendance),
			Creator:           event.Creator,
		}
		response.Events = append(response.Events, pb_event)
	}
	return response, nil
}

func (s *EventGRPCService) Register(
	ctx context.Context,
	req *pb.RegisterRequest,
) (*pb.EmptyResponse, error) {
	err := s.validate.Var(req.UserId, "uuid,required")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = s.eventService.Register(ctx, req.UserId, int(req.EventId))
	if err != nil {
		if errors.Is(err, service.ErrMaxRegistered) {
			return nil, status.Error(codes.ResourceExhausted, err.Error())
		}
		if errors.Is(err, service.ErrRegistered) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		if errors.Is(err, service.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.EmptyResponse{}, nil
}

func (s *EventGRPCService) CancellRegister(
	ctx context.Context,
	req *pb.CancellRegisterRequest,
) (*pb.EmptyResponse, error) {
	err := s.validate.Var(req.UserId, "uuid,required")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = s.eventService.CancellRegister(ctx, req.UserId, int(req.EventId))
	if err != nil {
		if errors.Is(err, service.ErrNotRegistered) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, service.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.EmptyResponse{}, nil
}

func (s *EventGRPCService) GetAllByUser(
	ctx context.Context,
	req *pb.GetAllByUserRequest,
) (*pb.GetAllByUserResponse, error) {
	err := s.validate.Var(req.UserId, "uuid,required")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	events, err := s.eventService.GetAllByUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	response := &pb.GetAllByUserResponse{
		Events: []*pb.EventElem{},
	}
	for _, event := range events {
		pb_event := &pb.EventElem{
			Id:                int64(event.Id),
			Title:             event.Title,
			About:             event.About,
			StartDate:         timestamppb.New(event.StartDate),
			Location:          event.Location,
			Status:            string(event.Status),
			MaxAttendees:      int32(event.MaxAttendees),
			CurrentAttendance: int32(event.CurrentAttendance),
			Creator:           event.Creator,
		}
		response.Events = append(response.Events, pb_event)
	}
	return response, nil
}

func (s *EventGRPCService) GetAllUsersByEvent(
	ctx context.Context,
	req *pb.GetAllUsersByEventRequest,
) (*pb.GetAllUsersByEventResponse, error) {
	users_id, err := s.eventService.GetAllUsersByEvent(ctx, int(req.EventId))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	response := &pb.GetAllUsersByEventResponse{}
	for _, id := range *users_id {
		response.UsersId = append(response.UsersId, id)
	}
	return response, nil
}
