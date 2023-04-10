package airflow

import (
	"context"
	"fmt"

	"github.com/tikivn/clickhousectl/internal/pkg/entity"
	"github.com/tikivn/clickhousectl/internal/pkg/entity/operator"
	"github.com/tikivn/clickhousectl/internal/pkg/service"
	"github.com/tikivn/clickhousectl/internal/utils/readfile"
)

func NewApp(
	airflowService service.AirFlowService,
) App {
	return &app{airflowService: airflowService}
}

type App interface {
	InsertData(ctx context.Context, rawParams ...string) error
}

type app struct {
	airflowService service.AirFlowService
}

func (a *app) prepareParams(requires int, params ...string) ([]string, error) {
	var results = make([]string, 0, len(params))
	for idx, param := range params {
		if idx < requires && param == "" {
			return nil, fmt.Errorf("param%d is missing", idx+1)
		}
		if param != "" {
			results = append(results, param)
		}
	}
	return results, nil
}

func (a *app) InsertData(ctx context.Context, rawParams ...string) error {
	params, err := a.prepareParams(1, rawParams...)
	if err != nil {
		return fmt.Errorf("fail to parse params: %v", err)
	}

	// Read table DDL here
	var tableConfig entity.TableDDL
	readfile.ReadJSON(params[0], &tableConfig)

	// Read query info here
	var req operator.OperatorInsertData
	readfile.ReadJSON(params[1], &req)

	if err != nil {
		return fmt.Errorf("fail to query from json file: %v", err)
	}

	if err := a.airflowService.Execute(ctx, tableConfig, &req); err != nil {
		return fmt.Errorf("fail to execute insert query: %v", err)
	}

	return nil
}

// ! Cannot use now bc Clickhouse doesnt support Delete from distributed table but only replicate
func (a *app) DeleteData(ctx context.Context, rawParams ...string) error {
	params, err := a.prepareParams(1, rawParams...)
	if err != nil {
		return fmt.Errorf("fail to parse params: %v", err)
	}

	// Read table DDL here
	var tableConfig entity.TableDDL
	readfile.ReadJSON(params[0], &tableConfig)

	// Read query info here
	var req operator.OperatorDeleteData
	readfile.ReadJSON(params[1], &req)

	return nil
}
