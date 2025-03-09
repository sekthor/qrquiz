package repo

import (
	"context"

	"github.com/sekthor/qrquiz/internal/domain"
)

type Repo interface {
	GetQuiz(ctx context.Context, id string) (domain.Quiz, error)
	Save(ctx context.Context, quiz domain.Quiz) error
	List(ctx context.Context, page int, size int) ([]domain.Quiz, error)
	DeleteExpired(ctx context.Context) error
}
