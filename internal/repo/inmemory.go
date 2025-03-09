package repo

import (
	"context"
	"time"

	"github.com/sekthor/qrquiz/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var _ Repo = &inMemoryRepo{}

type inMemoryRepo struct {
	quizes []domain.Quiz
	tracer trace.Tracer
}

func NewInMemoryRepo() *inMemoryRepo {
	return &inMemoryRepo{
		quizes: []domain.Quiz{},
		tracer: otel.Tracer("repo"),
	}
}

func (i *inMemoryRepo) GetQuiz(ctx context.Context, id string) (domain.Quiz, error) {
	_, span := i.tracer.Start(ctx, "inMemoryRepo.GetQuiz")
	defer span.End()

	for _, quiz := range i.quizes {
		if quiz.ID == id {
			return quiz, nil
		}
	}
	return domain.Quiz{}, ErrQuizNotFound
}

func (i *inMemoryRepo) Save(ctx context.Context, quiz domain.Quiz) error {
	_, span := i.tracer.Start(ctx, "inMemoryRepo.Save")
	defer span.End()

	i.quizes = append(i.quizes, quiz)
	return nil
}

func (i *inMemoryRepo) List(ctx context.Context, page int, size int) ([]domain.Quiz, error) {
	_, span := i.tracer.Start(ctx, "inMemoryRepo.List")
	defer span.End()

	start := (page - 1) * size
	end := start + size

	if len(i.quizes) < start {
		return []domain.Quiz{}, nil
	}

	if len(i.quizes) < end {
		return i.quizes[start:], nil
	}

	return i.quizes[start:end], nil
}

func (i *inMemoryRepo) DeleteExpired(ctx context.Context) error {
	_, span := i.tracer.Start(ctx, "inMemoryRepo.DeleteExpired")
	defer span.End()

	var unexpired []domain.Quiz
	now := time.Now()
	for _, quiz := range i.quizes {
		if quiz.Expires.Before(now) {
			unexpired = append(unexpired, quiz)
		}
	}
	i.quizes = unexpired
	return nil
}
