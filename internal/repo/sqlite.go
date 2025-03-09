package repo

import (
	"log"

	"github.com/sekthor/qrquiz/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ Repo = sqliteRepo{}

type sqliteRepo struct {
	db *gorm.DB
}

func NewSqliteRepo() sqliteRepo {
	db, err := gorm.Open(sqlite.Open("data/qrquiz.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	repo := sqliteRepo{
		db: db,
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

func (s sqliteRepo) GetQuiz(id string) (domain.Quiz, error) {
	var quiz domain.Quiz
	result := s.db.Preload(clause.Associations).
		Preload("Questions.Answers.Pixels").
		First(&quiz, "id = ?", id)
	return quiz, result.Error
}

func (s sqliteRepo) Save(quiz domain.Quiz) error {
	return s.db.Save(&quiz).Error
}

func (s sqliteRepo) List(page int, size int) ([]domain.Quiz, error) {
	var list []domain.Quiz
	result := s.db.Offset((page - 1) * size).Limit(size).Find(&list)
	return list, result.Error
}
