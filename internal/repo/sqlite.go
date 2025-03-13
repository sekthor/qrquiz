package repo

import (
	"context"
	"log"
	"time"

	"github.com/sekthor/qrquiz/internal/domain"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ Repo = sqliteRepo{}

type sqliteRepo struct {
	db     *gorm.DB
	tracer trace.Tracer
}

func NewSqliteRepo() sqliteRepo {
	db, err := gorm.Open(sqlite.Open("data/qrquiz.db"), &gorm.Config{Logger: NewLogger()})
	if err != nil {
		log.Fatal(err)
	}

	repo := sqliteRepo{
		db:     db,
		tracer: otel.Tracer("repo"),
	}

	err = repo.db.AutoMigrate(
		&domain.Quiz{},
		&domain.Question{},
		&domain.Answer{},
		&domain.Pixel{},
	)

	if err != nil {
		log.Fatal(err)
	}

	return repo
}

func (s sqliteRepo) GetQuiz(ctx context.Context, id string) (domain.Quiz, error) {
	_, span := s.tracer.Start(ctx, "sqliteRepo.GetQuiz")
	defer span.End()

	var quiz domain.Quiz
	result := s.db.WithContext(ctx).
		Preload(clause.Associations).
		Preload("Questions.Answers.Pixels").
		First(&quiz, "id = ?", id)
	return quiz, result.Error
}

func (s sqliteRepo) Save(ctx context.Context, quiz domain.Quiz) error {
	ctx, span := s.tracer.Start(ctx, "sqliteRepo.Save")
	defer span.End()

	return s.db.WithContext(ctx).Save(&quiz).Error
}

func (s sqliteRepo) List(ctx context.Context, page int, size int) ([]domain.Quiz, error) {
	ctx, span := s.tracer.Start(ctx, "sqliteRepo.List")
	defer span.End()

	var list []domain.Quiz
	result := s.db.WithContext(ctx).
		Offset((page - 1) * size).
		Limit(size).Find(&list)
	return list, result.Error
}

func (s sqliteRepo) DeleteExpired(ctx context.Context) error {
	ctx, span := s.tracer.Start(ctx, "sqliteRepo.DeleteExpired")
	defer span.End()

	result := s.db.
		WithContext(ctx).
		Where("expires < ?", time.Now()).
		Delete(&domain.Quiz{})

	if result.RowsAffected != 0 {
		logrus.WithContext(ctx).Debugf("deleted %d expired quizzes", result.RowsAffected)
	} else {
		logrus.WithContext(ctx).Debug("no expired quizzes to delete")
	}

	return result.Error
}
