package space

import (
	"context"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/usecase"
	"github.com/Slava02/Involvio/internal/usecase/commands"
	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type ISpaceUseCase interface {
	CreateSpace(ctx context.Context, cmd commands.SpaceCommand) (*entity.Space, error)
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
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "CreateSpace", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := commands.SpaceCommand{
		Name:        req.Body.Name,
		Description: req.Body.Description,
		Tags:        req.Body.Tags,
	}

	space, err := sh.spaceUC.CreateSpace(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resp := ToSpaceOutputFromEntity(space)

	return resp, nil
}

func (sh *SpaceHandler) GetSpace(ctx context.Context, req *SpaceByIdRequest) (*SpaceResponse, error) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "GetSpace", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := commands.SpaceByIdCommand{
		ID: req.ID,
	}

	space, err := sh.spaceUC.GetSpace(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resp := ToSpaceOutputFromEntity(space)

	return resp, nil
}

func (sh *SpaceHandler) JoinSpace(ctx context.Context, req *JoinSpaceRequest) (*struct{}, error) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "JoinSpace", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := commands.JoinSpaceCommand{
		SpaceID: req.Body.SpaceId,
		UserID:  req.Body.UserId,
	}

	err := sh.spaceUC.JoinSpace(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &struct{}{}, nil
}

func (sh *SpaceHandler) UpdateSpace(ctx context.Context, req *UpdateSpaceRequest) (*SpaceResponse, error) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "UpdateSpace", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := commands.UpdateSpaceCommand{
		ID:          req.ID,
		Name:        req.Body.Name,
		Description: req.Body.Description,
	}

	space, err := sh.spaceUC.UpdateSpace(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resp := ToSpaceOutputFromEntity(space)

	return resp, nil
}

func (sh *SpaceHandler) DeleteSpace(ctx context.Context, req *SpaceByIdRequest) (*struct{}, error) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "DeleteSpace", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := commands.SpaceByIdCommand{
		ID: req.ID,
	}

	err := sh.spaceUC.DeleteSpace(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &struct{}{}, nil
}
