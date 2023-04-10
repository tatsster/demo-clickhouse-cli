package mysql

import (
	"context"
	"encoding/json"

	"github.com/tikivn/clickhousectl/internal/pkg/entity"
	"github.com/tikivn/clickhousectl/internal/pkg/repo"
	"gorm.io/gorm"
)

func NewUserRepo(conn *gorm.DB) repo.UserRepo {
	return &userRepo{conn: conn}
}

type userRepo struct {
	conn *gorm.DB
}

func (repo *userRepo) Create(ctx context.Context, user *entity.User) error {
	dao, err := repo.serialize(user)
	if err != nil {
		return err
	}

	res := repo.conn.Create(dao)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (repo *userRepo) Update(ctx context.Context, user *entity.User) error {
	panic("implement me")
}

func (repo *userRepo) Delete(ctx context.Context, username string) error {
	panic("implement me")
}

func (repo *userRepo) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	params := map[string]interface{}{
		"username": username,
	}

	var dao User
	if find := repo.conn.Where(params).Find(&dao); find.Error != nil {
		return nil, find.Error
	}

	user, err := repo.deserialize(&dao)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *userRepo) serialize(user *entity.User) (*User, error) {
	databases, err := json.Marshal(user.AllowDatabases)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:  user.Username,
		Password:  user.Password,
		Databases: string(databases),
		Status:    string(user.Status),
	}, nil
}

func (repo *userRepo) deserialize(user *User) (*entity.User, error) {
	var allowDatabases []string
	err := json.Unmarshal([]byte(user.Databases), &allowDatabases)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		Username:       user.Username,
		Password:       user.Password,
		AllowDatabases: allowDatabases,
		Status:         entity.UserStatus(user.Status),
	}, nil
}

// ! From MySQL
type User struct {
	Username  string `gorm:"column:username;type:varchar(255);primary_key"`
	Password  string `gorm:"column:password;type:text"`
	Databases string `gorm:"column:allow_databases;type:text"`
	Status    string `gorm:"column:status;type:varchar(64)"`
}

func (u *User) TableName() string {
	return "clickhouse_users"
}
