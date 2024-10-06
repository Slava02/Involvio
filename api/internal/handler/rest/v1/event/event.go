package event

import (
	"context"
	"errors"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/repository"
	"github.com/Slava02/Involvio/internal/usecase"
	"github.com/Slava02/Involvio/internal/usecase/commands"
	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

type IEventUseCase interface {
	CreateEvent(ctx context.Context, cmd commands.CreateEventCommand) (*entity.Event, error)
	GetEvent(ctx context.Context, cmd commands.EventByIdCommand) (*entity.Event, error)
	JoinEvent(ctx context.Context, cmd commands.JoinEventCommand) error
	DeleteEvent(ctx context.Context, cmd commands.EventByIdCommand) error
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
	)
	log.Debug(op)

	b := req.Body
	// TODO: добавить максимальное количество участников
	cmd := commands.CreateEventCommand{
		UserId:      b.UserId,
		SpaceId:     b.SpaceId,
		Name:        b.EventInfo.Name,
		Description: b.EventInfo.Description,
		BeginDate:   b.EventInfo.BeginDate,
		EndDate:     b.EventInfo.EndDate,
		Tags:        b.EventInfo.Tags,
	}

	event, err := eh.eventUC.CreateEvent(ctx, cmd)
	if err != nil {
		switch {
		default:
			log.Error("couldn't join event: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	resp := ToEventOutputFromEntity(event)

	return resp, nil
}

func (eh *EventHandler) GetEvent(ctx context.Context, req *EventByIdRequest) (*EventResponse, error) {
	const op = "Handler:JoinEvent"

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

	event, err := eh.eventUC.GetEvent(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEventNotFound):
			log.Info("couldn't get event: ", err.Error())
			return nil, huma.Error404NotFound(err.Error())
		default:
			log.Error("couldn't get event: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	resp := ToEventOutputFromEntity(event)

	return resp, nil
}

func (eh *EventHandler) JoinEvent(ctx context.Context, req *JoinEventRequest) (*struct{}, error) {
	const op = "Handler:JoinEvent"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
		slog.Int("event id", req.EventId),
		slog.Int("user id", req.Body.UserId),
	)
	log.Debug(op)

	cmd := commands.JoinEventCommand{
		EventId: req.EventId,
		UserId:  req.Body.UserId,
	}

	err := eh.eventUC.JoinEvent(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEventAlreadyExists):
			log.Info("couldn't join event: ", err.Error())
			return nil, huma.Error400BadRequest("user in event already exists")
		case errors.Is(err, repository.ErrEventNotFound):
			log.Info("couldn't join event: ", err.Error())
			return nil, huma.Error404NotFound("space not found")
		default:
			log.Error("couldn't join event: ", err.Error())
			return nil, huma.Error500InternalServerError("internal service error")
		}
	}

	return &struct{}{}, nil
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
		case errors.Is(err, repository.ErrEventNotFound):
			log.Info("couldn't delete event: ", err.Error())
			return nil, huma.Error404NotFound(err.Error())
		default:
			log.Error("couldn't delete event: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	return &struct{}{}, nil
}
