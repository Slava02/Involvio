package event

import (
	"context"
	"errors"
	"github.com/Slava02/Involvio/api/internal/entity"
	"github.com/Slava02/Involvio/api/internal/repository"
	"github.com/Slava02/Involvio/api/internal/usecase"
	"github.com/Slava02/Involvio/api/internal/usecase/commands"
	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

type IEventUseCase interface {
	CreateEvent(ctx context.Context, cmd commands.CreateEventCommand) (*entity.Event, error)
	GetUserEvents(ctx context.Context, cmd commands.EventByUserIdCommand) ([]entity.Event, error)
	JoinEvent(ctx context.Context, cmd commands.JoinEventCommand) error
	DeleteEvent(ctx context.Context, cmd commands.EventByIdCommand) error
	ReviewEvent(ctx context.Context, cmd commands.ReviewEventCommand) error
}

var _ IEventUseCase = (*usecase.EventUseCase)(nil)

const tracerName = "event handler"

type EventHandler struct {
	eventUC IEventUseCase
}

func NewEventHandler(uc IEventUseCase) *EventHandler {
	return &EventHandler{eventUC: uc}
}

func (eh *EventHandler) CreateEvent(ctx context.Context, req *CreateEventRequest) (*EventResponse, error) {
	const op = "Handler:CreateEvent"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
		slog.String("name:", req.Body.Name),
	)
	log.Debug(op)

	cmd := commands.CreateEventCommand{
		Name:        req.Body.Name,
		Description: req.Body.Description,
		Users:       req.Body.Users,
		Date:        req.Body.Date,
	}

	event, err := eh.eventUC.CreateEvent(ctx, cmd)
	if err != nil {
		switch {
		default:
			log.Error("couldn't create event: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	resp := ToEventOutputFromEntity(*event)

	return resp, nil
}

func (eh *EventHandler) GetUserEvents(ctx context.Context, req *EventByUserIdRequest) (*UserEventsResponse, error) {
	const op = "Handler:GetUserEvents"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
		slog.Int("event id", req.ID),
	)
	log.Debug(op)

	cmd := commands.EventByUserIdCommand{
		ID: req.ID,
	}

	events, err := eh.eventUC.GetUserEvents(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			log.Info("couldn't get event: ", err.Error())
			return nil, huma.Error404NotFound(err.Error())
		default:
			log.Error("couldn't get event: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	resp := ToEventsOutputFromEntity(events)

	return resp, nil
}

func (eh *EventHandler) DeleteEvent(ctx context.Context, req *EventByIdRequest) (*struct{}, error) {
	const op = "Handler:DeleteEvent"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
		slog.Int("event id", req.ID),
	)
	log.Debug(op)

	cmd := commands.EventByIdCommand{
		ID: req.ID,
	}

	err := eh.eventUC.DeleteEvent(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			log.Info("couldn't delete event: ", err.Error())
			return nil, huma.Error404NotFound(err.Error())
		default:
			log.Error("couldn't delete event: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	return &struct{}{}, nil
}

func (eh *EventHandler) JoinEvent(ctx context.Context, req *JoinEventRequest) (*struct{}, error) {
	const op = "Handler:JoinEvent"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
		slog.Int("eventID", req.EventId),
		slog.Int("userID", req.Body.UserId),
	)
	log.Debug(op)

	cmd := commands.JoinEventCommand{
		EventId: req.EventId,
		UserId:  req.Body.UserId,
	}

	err := eh.eventUC.JoinEvent(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrAlreadyExists):
			log.Info("couldn't join event: ", err.Error())
			return nil, huma.Error400BadRequest("user in event already exists")
		case errors.Is(err, repository.ErrNotFound):
			log.Info("couldn't join event: ", err.Error())
			return nil, huma.Error404NotFound("group not found")
		default:
			log.Error("couldn't join event: ", err.Error())
			return nil, huma.Error500InternalServerError("internal service error")
		}
	}

	return &struct{}{}, nil
}

func (eh *EventHandler) ReviewEvent(ctx context.Context, req *ReviewEventRequest) (*struct{}, error) {
	const op = "Handler:ReviewEvent"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	cmd := commands.ReviewEventCommand{
		EventID: req.Body.EventID,
		WhoID:   req.Body.WhoID,
		WhomID:  req.Body.WhomID,
		Grade:   req.Body.Grade,
	}

	err := eh.eventUC.ReviewEvent(ctx, cmd)
	if err != nil {
		switch {
		default:
			log.Error("couldn't review event: ", err.Error())
			return nil, huma.Error500InternalServerError("internal service error")
		}
	}

	return &struct{}{}, nil
}
