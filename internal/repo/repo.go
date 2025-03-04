package repo

import "github.com/sekthor/puzzleinvite/internal/domain"

type Repo interface {
	GetQuiz(id string) (domain.Quiz, error)
	Save(quiz domain.Quiz) error
}
