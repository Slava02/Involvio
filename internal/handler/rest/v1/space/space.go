package space

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

type ISpaceUseCase interface {
	CreateSpace(ctx context.Context, cmd commands.CreateSpaceCommand) (*entity.Space, error)
	GetSpace(ctx context.Context, cmd commands.SpaceByIdCommand) (*entity.Space, error)
	JoinSpace(ctx context.Context, cmd commands.JoinSpaceCommand) error
	UpdateSpace(ctx context.Context, cmd commands.UpdateSpaceCommand) (*entity.Space, error)
	DeleteSpace(ctx context.Context, cmd commands.SpaceByIdCommand) error
}

var _ ISpaceUseCase = (*usecase.SpaceUseCase)(nil)

const tracerName = "space handler"

type SpaceHandler struct {
	spaceUC ISpaceUseCase
}

func NewSpaceHandler(uc ISpaceUseCase) *SpaceHandler {
	return &SpaceHandler{spaceUC: uc}
}

func (sh *SpaceHandler) CreateSpace(ctx context.Context, req *CreateSpaceRequest) (*SpaceResponse, error) {
	const op = "Handler:CreateSpace"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	cmd := commands.CreateSpaceCommand{
		UserID:      req.Body.UserId,
		Name:        req.Body.Name,
		Description: req.Body.Description,
		Tags:        req.Body.Tags,
	}

	space, err := sh.spaceUC.CreateSpace(ctx, cmd)
	if err != nil {
		switch {
		default:
			log.Error("couldn't create space: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	resp := ToSpaceOutputFromEntity(space)

	return resp, nil
}

func (sh *SpaceHandler) GetSpace(ctx context.Context, req *SpaceByIdRequest) (*SpaceResponse, error) {
	const op = "Handler:GetSpace"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
		slog.Int("space id", req.ID),
	)
	log.Debug(op)

	cmd := commands.SpaceByIdCommand{
		ID: req.ID,
	}

	space, err := sh.spaceUC.GetSpace(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEventNotFound):
			log.Info("couldn't get space: ", err.Error())
			return nil, huma.Error404NotFound(err.Error())
		default:
			log.Error("couldn't get space: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	resp := ToSpaceOutputFromEntity(space)

	return resp, nil
}

func (sh *SpaceHandler) JoinSpace(ctx context.Context, req *JoinSpaceRequest) (*struct{}, error) {
	const op = "Handler:JoinSpace"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
		slog.Int("space id", req.Body.SpaceId),
		slog.Int("user id", req.Body.UserId),
	)
	log.Debug(op)

	cmd := commands.JoinSpaceCommand{
		SpaceID: req.Body.SpaceId,
		UserID:  req.Body.UserId,
	}

	err := sh.spaceUC.JoinSpace(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEventNotFound):
			log.Info("couldn't get space: ", err.Error())
			return nil, huma.Error404NotFound(err.Error())
		default:
			log.Error("couldn't get space: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	return &struct{}{}, nil
}

func (sh *SpaceHandler) UpdateSpace(ctx context.Context, req *UpdateSpaceRequest) (*SpaceResponse, error) {
	const op = "Handler:UpdateSpace"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
		slog.Int("space id", req.ID),
	)
	log.Debug(op)

	cmd := commands.UpdateSpaceCommand{
		ID:          req.ID,
		Name:        req.Body.Name,
		Description: req.Body.Description,
	}

	space, err := sh.spaceUC.UpdateSpace(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEventNotFound):
			log.Info("couldn't get space: ", err.Error())
			return nil, huma.Error404NotFound(err.Error())
		default:
			log.Error("couldn't get space: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	resp := ToSpaceOutputFromEntity(space)

	return resp, nil
}

func (sh *SpaceHandler) DeleteSpace(ctx context.Context, req *SpaceByIdRequest) (*struct{}, error) {
	const op = "Handler:DeleteSpace"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
		slog.Int("space id", req.ID),
	)
	log.Debug(op)

	cmd := commands.SpaceByIdCommand{
		ID: req.ID,
	}

	err := sh.spaceUC.DeleteSpace(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEventNotFound):
			log.Info("couldn't get space: ", err.Error())
			return nil, huma.Error404NotFound(err.Error())
		default:
			log.Error("couldn't get space: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	return &struct{}{}, nil
}
