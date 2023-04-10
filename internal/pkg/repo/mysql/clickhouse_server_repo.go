package mysql

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/tikivn/clickhousectl/internal/pkg/entity"
	"github.com/tikivn/clickhousectl/internal/pkg/repo"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func NewClickHouseServerRepo(conn *gorm.DB, userRepo repo.UserRepo) repo.ClickHouseServerRepo {
	return &clickHouseServerRepoImpl{conn: conn, userRepo: userRepo}
}

type clickHouseServerRepoImpl struct {
	conn     *gorm.DB
	userRepo repo.UserRepo
}

// Get Clickhouse servers host port
func (repo *clickHouseServerRepoImpl) FindAll(ctx context.Context, orgId string) ([]*entity.ClickHouseServer, error) {
	// Find all clickhouse servers are belong to orgId 
	params := map[string]interface{}{
		"org_id": orgId,
	}

	var daos []*ClickHouseServer
	if find := repo.conn.Where(params).Find(&daos); find.Error != nil {
		return nil, find.Error
	}

	if len(daos) == 0 {
		return nil, errors.New("empty Server")
	}

	username := daos[0].Username
	user, err := repo.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	var servers = make([]*entity.ClickHouseServer, 0, len(daos))
	for _, dao := range daos {
		server, err := repo.deserialize(dao)
		if err != nil {
			return nil, err
		}
		server.Password = user.Password
		servers = append(servers, server)
	}

	return servers, nil
}

func (repo *clickHouseServerRepoImpl) deserialize(dao *ClickHouseServer) (*entity.ClickHouseServer, error) {
	var shards []string
	err := json.Unmarshal(dao.Shards, &shards)
	if err != nil {
		return nil, err
	}

	return &entity.ClickHouseServer{
		Id:       dao.Id,
		OrgId:    dao.OrgId,
		Host:     dao.Host,
		Port:     dao.Port,
		Cluster:  dao.Cluster,
		Username: dao.Username,
		Shards:   shards,
	}, nil
}

// ! From MySQL
type ClickHouseServer struct {
	Id       string         `gorm:"column:id;type:varchar(36);primary_key"`
	OrgId    string         `gorm:"column:org_id;type:varchar(36)"`
	Host     string         `gorm:"column:host;type:varchar(255)"`
	Port     string         `gorm:"column:port;type:varchar(16)"`
	Cluster  string         `gorm:"column:cluster;type:varchar(255)"`
	Username string         `gorm:"column:username;type:varchar(255)"`
	Shards   datatypes.JSON `gorm:"column:shards;type:json"`
}

func (chs *ClickHouseServer) TableName() string {
	return "clickhouse_servers"
}
