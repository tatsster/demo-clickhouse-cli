package repo

import (
	"context"

	"github.com/tikivn/clickhousectl/internal/pkg/entity"
)

type UserRepo interface {
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, username string) error
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
}
