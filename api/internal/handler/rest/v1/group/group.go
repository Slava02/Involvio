package group

import (
	"context"
	"github.com/Slava02/Involvio/api/internal/entity"
	"github.com/Slava02/Involvio/api/internal/usecase"
	"github.com/Slava02/Involvio/api/internal/usecase/commands"
	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

type IGroupUseCase interface {
	CreateGroup(ctx context.Context, cmd commands.GroupByNameCommand) (*entity.Group, error)
	GetGroup(ctx context.Context, cmd commands.GroupByNameCommand) (*entity.Group, error)
	JoinGroup(ctx context.Context, cmd commands.JoinLeaveGroupCommand) error
	LeaveGroup(ctx context.Context, cmd commands.JoinLeaveGroupCommand) error
	DeleteGroup(ctx context.Context, cmd commands.GroupByNameCommand) error
}

var _ IGroupUseCase = (*usecase.GroupUseCase)(nil)

const tracerName = "group handler"

type GroupHandler struct {
	groupUC IGroupUseCase
}

func NewGroupHandler(uc IGroupUseCase) *GroupHandler {
	return &GroupHandler{groupUC: uc}
}

func (sh *GroupHandler) CreateGroup(ctx context.Context, req *CreateGroupRequest) (*GroupResponse, error) {
	const op = "Handler:CreateGroup"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	cmd := commands.GroupByNameCommand{
		Name: req.Body.Name,
	}

	group, err := sh.groupUC.CreateGroup(ctx, cmd)
	if err != nil {
		switch {
		default:
			log.Error("couldn't create group: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	resp := ToGroupOutputFromEntity(group)

	return resp, nil
}

func (sh *GroupHandler) GetGroup(ctx context.Context, req *GroupByNameRequest) (*GroupResponse, error) {
	const op = "Handler:GetGroup"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
		slog.String("groupName", req.Name),
	)
	log.Debug(op)

	cmd := commands.GroupByNameCommand{
		Name: req.Name,
	}

	group, err := sh.groupUC.GetGroup(ctx, cmd)
	if err != nil {
		switch {
		default:
			log.Error("couldn't get group: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	resp := ToGroupOutputFromEntity(group)

	return resp, nil
}

func (sh *GroupHandler) DeleteGroup(ctx context.Context, req *GroupByNameRequest) (*struct{}, error) {
	const op = "Handler:DeleteGroup"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
		slog.String("groupName", req.Name),
	)
	log.Debug(op)

	cmd := commands.GroupByNameCommand{
		Name: req.Name,
	}

	err := sh.groupUC.DeleteGroup(ctx, cmd)
	if err != nil {
		switch {
		default:
			log.Error("couldn't delete group: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	return &struct{}{}, nil
}

func (sh *GroupHandler) JoinGroup(ctx context.Context, req *JoinLeaveGroupRequest) (*struct{}, error) {
	const op = "Handler:JoinGroup"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
		slog.String("groupName", req.Body.GroupName),
		slog.Int("userID", req.Body.UserId),
	)
	log.Debug(op)

	cmd := commands.JoinLeaveGroupCommand{
		GroupName: req.Body.GroupName,
		UserID:    req.Body.UserId,
	}

	err := sh.groupUC.JoinGroup(ctx, cmd)
	if err != nil {
		switch {
		default:
			log.Error("couldn't join group: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	return &struct{}{}, nil
}

func (sh *GroupHandler) LeaveGroup(ctx context.Context, req *JoinLeaveGroupRequest) (*struct{}, error) {
	const op = "Handler:LeaveGroup"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
		slog.String("groupName", req.Body.GroupName),
		slog.Int("userID", req.Body.UserId),
	)
	log.Debug(op)

	cmd := commands.JoinLeaveGroupCommand{
		GroupName: req.Body.GroupName,
		UserID:    req.Body.UserId,
	}

	err := sh.groupUC.LeaveGroup(ctx, cmd)
	if err != nil {
		switch {
		default:
			log.Error("couldn't leave group: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	return &struct{}{}, nil
}
