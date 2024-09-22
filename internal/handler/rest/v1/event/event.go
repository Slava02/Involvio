package event

import (
	"context"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/usecase"
	"github.com/Slava02/Involvio/internal/usecase/commands"
	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
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
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "CreateEvent", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

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
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resp := ToEventOutputFromEntity(event)

	return resp, nil
}

func (eh *EventHandler) GetEvent(ctx context.Context, req *EventByIdRequest) (*EventResponse, error) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "GetEvent", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := commands.EventByIdCommand{
		ID: req.ID,
	}

	event, err := eh.eventUC.GetEvent(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resp := ToEventOutputFromEntity(event)

	return resp, nil
}

func (eh *EventHandler) JoinEvent(ctx context.Context, req *JoinEventRequest) (*struct{}, error) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "JoinEvent", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := commands.JoinEventCommand{
		EventId: req.EventId,
		UserId:  req.Body.UserId,
	}

	err := eh.eventUC.JoinEvent(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &struct{}{}, nil
}

func (eh *EventHandler) DeleteEvent(ctx context.Context, req *EventByIdRequest) (*struct{}, error) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "DeleteEvent", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := commands.EventByIdCommand{
		ID: req.ID,
	}

	err := eh.eventUC.DeleteEvent(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &struct{}{}, nil
}
