package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/tikivn/clickhousectl/internal/pkg/entity"
	"github.com/tikivn/clickhousectl/internal/pkg/entity/operator"
	"github.com/tikivn/clickhousectl/internal/pkg/infra"
	"github.com/tikivn/clickhousectl/internal/pkg/repo"
	"golang.org/x/sync/errgroup"
)

func NewColumnOperatorService(
	clickHouseServerRepo repo.ClickHouseServerRepo,
) ColumnOperatorService {
	return &columnOperatorServiceImpl{
		clickHouseServerRepo: clickHouseServerRepo,
	}
}

type ColumnOperatorService interface {
	Execute(ctx context.Context, opEntity operator.OperatorColumn) error
}

type columnOperatorServiceImpl struct {
	clickHouseServerRepo repo.ClickHouseServerRepo
}

func (svc *columnOperatorServiceImpl) Execute(ctx context.Context, opEntity operator.OperatorColumn) error {
	clickHouseServers, err := svc.clickHouseServerRepo.FindAll(ctx, opEntity.ServerId())

	if err != nil {
		return err
	}

	// Random choose 1 server to get replica table
	rand.Seed(time.Now().Unix())
	server := clickHouseServers[rand.Intn(len(clickHouseServers))]

	// Create operator with table
	op, err := svc.operator(clickHouseServers, server, opEntity)
	if err != nil {
		return err
	}

	// execute
	var eg errgroup.Group
	for id := range clickHouseServers {
		conn := infra.NewClickHouseConnection(clickHouseServers[id].Host, clickHouseServers[id].Port).
			WithAuthentication(clickHouseServers[id].Username, clickHouseServers[id].Password)

		clickHouseSession, cleanup, err := infra.NewClickhouseSession(conn)
		if err != nil {
			return err
		}

		defer cleanup()

		stmts := op.ExecuteStmts[clickHouseServers[id].Id]
		for _, stmt := range stmts {
			if err := clickHouseSession.Exec(stmt).Error; err != nil {
				return err
			}
		}
	}

	execErr := eg.Wait()
	if execErr != nil {
		return execErr
	}

	return nil
}

func (svc *columnOperatorServiceImpl) operator(
	servers []*entity.ClickHouseServer,
	server *entity.ClickHouseServer,
	opEntity operator.OperatorColumn,
) (*operator.Operator, error) {
	executes, err := opEntity.Stmts(servers, server.Id)
	if err != nil {
		return nil, err
	}

	return &operator.Operator{
		Id:           uuid.New().String(),
		Owner:        "",
		Name:         opEntity.OperatorName(),
		ExecuteStmts: executes,
		Status:       operator.OperatorStatus_EXECUTING,
		ExecutedAt:   time.Now().Unix(),
	}, nil
}
