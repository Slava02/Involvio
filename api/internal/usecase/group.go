package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/Slava02/Involvio/api/internal/entity"
	"github.com/Slava02/Involvio/api/internal/repository"
	"github.com/Slava02/Involvio/api/internal/usecase/commands"
	"github.com/Slava02/Involvio/api/pkg/hexid"
	"log/slog"
)

type IGroupRepository interface {
	GetGroup(ctx context.Context, name string) (*entity.Group, error)
	DeleteGroup(ctx context.Context, name string) error
	InsertGroup(ctx context.Context, group *entity.Group) error
	AddUser(ctx context.Context, userId int, groupName string) error
	RemoveUser(ctx context.Context, userId int, groupName string) error
}

func NewGroupUseCase(ur IGroupRepository) *GroupUseCase {
	return &GroupUseCase{groupRepo: ur}
}

type GroupUseCase struct {
	groupRepo IGroupRepository
}

func (sc *GroupUseCase) DeleteGroup(ctx context.Context, cmd commands.GroupByNameCommand) error {
	const op = "UseCase:DeleteGroup"

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	_, err := sc.GetGroup(ctx, commands.GroupByNameCommand{Name: cmd.Name})
	if err != nil {
		log.Debug("couldn't get group: ", err.Error())
		return fail(err)
	}

	err = sc.groupRepo.DeleteGroup(ctx, cmd.Name)
	if err != nil {
		log.Debug("couldn't delete group: ", err.Error())
		return fail(err)
	}

	return nil
}

func (sc *GroupUseCase) CreateGroup(ctx context.Context, cmd commands.GroupByNameCommand) (*entity.Group, error) {
	const op = "UseCase:CreateGroup"

	fail := func(err error) (*entity.Group, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	// TODO: вынести в зависимости
	groupId, err := hexid.Generate()
	if err != nil {
		log.Error("couldn't generate id: ", err.Error())
		return fail(err)
	}

	group := &entity.Group{
		ID:   groupId,
		Name: cmd.Name,
	}

	err = sc.groupRepo.InsertGroup(ctx, group)
	if err != nil {
		log.Debug("couldn't insert group: ", err.Error())
		return fail(err)
	}

	return group, nil
}

func (sc *GroupUseCase) GetGroup(ctx context.Context, cmd commands.GroupByNameCommand) (*entity.Group, error) {
	const op = "UseCase:GetGroup"

	fail := func(err error) (*entity.Group, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
		slog.String("groupName", cmd.Name),
	)
	log.Debug(op)

	group, err := sc.groupRepo.GetGroup(ctx, cmd.Name)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Info("couldn't get event: ", err.Error())
		} else {
			log.Debug("couldn't get event: ", err.Error())
		}
		return fail(err)
	}

	return group, nil
}

func (sc *GroupUseCase) JoinGroup(ctx context.Context, cmd commands.JoinLeaveGroupCommand) error {
	const op = "UseCase:JoinGroup"

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	_, err := sc.GetGroup(ctx, commands.GroupByNameCommand{Name: cmd.GroupName})
	if err != nil {
		log.Debug("couldn't get group: ", err.Error())
		return fail(err)
	}

	err = sc.groupRepo.AddUser(ctx, cmd.UserID, cmd.GroupName)
	if err != nil {
		return fail(err)
	}

	return nil
}

func (sc *GroupUseCase) LeaveGroup(ctx context.Context, cmd commands.JoinLeaveGroupCommand) error {
	const op = "UseCase:LeaveGroup"

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
		slog.String("groupName", cmd.GroupName),
	)
	log.Debug(op)

	_, err := sc.GetGroup(ctx, commands.GroupByNameCommand{Name: cmd.GroupName})
	if err != nil {
		log.Debug("couldn't get group: ", err.Error())
		return fail(err)
	}

	err = sc.groupRepo.RemoveUser(ctx, cmd.UserID, cmd.GroupName)
	if err != nil {
		return fail(err)
	}

	return nil
}
