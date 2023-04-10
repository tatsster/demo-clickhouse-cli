package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tikivn/clickhousectl/internal/pkg/entity"
	"github.com/tikivn/clickhousectl/internal/pkg/entity/operator"
	"github.com/tikivn/clickhousectl/internal/pkg/infra"
	"github.com/tikivn/clickhousectl/internal/pkg/repo"
	"golang.org/x/sync/errgroup"
)

func NewTableOperatorService(
	clickHouseServerRepo repo.ClickHouseServerRepo,
) TableOperatorService {
	return &tableOperatorServiceImpl{
		clickHouseServerRepo: clickHouseServerRepo,
	}
}

type TableOperatorService interface {
	Execute(ctx context.Context, op operator.OperatorTable) error
}

type tableOperatorServiceImpl struct {
	// Repo get clickhouse server
	clickHouseServerRepo repo.ClickHouseServerRepo
}

func (svc *tableOperatorServiceImpl) Execute(ctx context.Context, opEntity operator.OperatorTable) error {
	clickHouseServers, err := svc.clickHouseServerRepo.FindAll(ctx, opEntity.ServerId())

	if err != nil {
		return err
	}

	// Create operator with table
	op, err := svc.operator(clickHouseServers, opEntity)
	if err != nil {
		return err
	}

	// execute
	var eg errgroup.Group
	for idx := range clickHouseServers {
		// This to bypass goroutine in loop only last item
		id := idx
		eg.Go(func() error {
			conn := infra.NewClickHouseConnection(clickHouseServers[id].Host, clickHouseServers[id].Port).
				WithAuthentication(clickHouseServers[id].Username, clickHouseServers[id].Password)

			clickHouseSession, cleanup, err := infra.NewClickhouseSession(conn)
			if err != nil {
				return err
			}

			defer cleanup()

			stmts := op.ExecuteStmts[clickHouseServers[id].Id]
			for _, stmt := range stmts {
				// err := conn.ExecuteDryRun(ctx, stmt)
				if err := clickHouseSession.Exec(stmt).Error; err != nil {
					return err
				}
			}
			return nil
		})
	}

	execErr := eg.Wait()
	if execErr != nil {
		return execErr
	}

	return nil
}

func (svc *tableOperatorServiceImpl) operator(servers []*entity.ClickHouseServer, opEntity operator.OperatorTable) (*operator.Operator, error) {
	executes, err := opEntity.Stmts(servers)
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
