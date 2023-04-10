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

func NewAirFlowService(
	clickHouseServerRepo repo.ClickHouseServerRepo,
) AirFlowService {
	return &airflowService{
		clickHouseServerRepo: clickHouseServerRepo,
	}
}

type AirFlowService interface {
	Execute(ctx context.Context, tableDDL entity.TableDDL, opEntity operator.OperatorArgs) error
}

type airflowService struct {
	clickHouseServerRepo repo.ClickHouseServerRepo
}

func (svc *airflowService) Execute(ctx context.Context, tableDDL entity.TableDDL, opEntity operator.OperatorArgs) error {
	clickHouseServers, err := svc.clickHouseServerRepo.FindAll(ctx, opEntity.ServerId())

	if err != nil {
		return err
	}

	rand.Seed(time.Now().Unix())
	server := clickHouseServers[rand.Intn(len(clickHouseServers))]

	// execute
	var eg errgroup.Group
	conn := infra.NewClickHouseConnection(server.Host, server.Port).
		WithAuthentication(server.Username, server.Password)

	clickHouseSession, cleanup, err := infra.NewClickhouseSession(conn)
	if err != nil {
		return err
	}

	defer cleanup()

	opInsert := opEntity.(*operator.OperatorInsertData)
	for i := range opInsert.Data {
		op, err := svc.operator(server, tableDDL, opEntity, i)
		if err != nil {
			return err
		}

		tx := clickHouseSession.Begin()
		err = tx.Exec(op.ExecuteStmt, op.Argument...).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Commit().Error
		if err != nil {
			return err
		}

		execErr := eg.Wait()
		if execErr != nil {
			return execErr
		}
	}

	return nil
}

func (svc *airflowService) operator(
	server *entity.ClickHouseServer,
	tableDDL entity.TableDDL,
	opEntity operator.OperatorArgs,
	recordIdx int,
) (*operator.ArgumentOperator, error) {
	execute, err := opEntity.Stmt(server, tableDDL, recordIdx)
	if err != nil {
		return nil, err
	}

	return &operator.ArgumentOperator{
		Id:          uuid.New().String(),
		Owner:       "",
		Name:        opEntity.OperatorName(),
		ExecuteStmt: execute,
		Argument:    opEntity.ValueField(tableDDL, recordIdx),
		Status:      operator.OperatorStatus_EXECUTING,
		ExecutedAt:  time.Now().Unix(),
	}, nil
}
