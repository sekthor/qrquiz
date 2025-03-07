package repo

import "github.com/sekthor/qrquiz/internal/domain"

var _ Repo = &inMemoryRepo{}

type inMemoryRepo struct {
	quizes []domain.Quiz
}

func NewInMemoryRepo() *inMemoryRepo {
	return &inMemoryRepo{
		quizes: []domain.Quiz{},
	}
}

func (i *inMemoryRepo) GetQuiz(id string) (domain.Quiz, error) {
	for _, quiz := range i.quizes {
		if quiz.ID == id {
			return quiz, nil
		}
	}
	return domain.Quiz{}, ErrQuizNotFound
}

func (i *inMemoryRepo) Save(quiz domain.Quiz) error {
	i.quizes = append(i.quizes, quiz)
	return nil
}

func (i *inMemoryRepo) List() ([]domain.Quiz, error) {
	return i.quizes, nil
}
