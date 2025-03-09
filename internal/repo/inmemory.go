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

func (i *inMemoryRepo) List(page int, size int) ([]domain.Quiz, error) {

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
