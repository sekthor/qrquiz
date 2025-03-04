package domain

import (
	"errors"
	"time"

	"github.com/nrednav/cuid2"
)

const (
	CorrectAnswer = true
	WrongAnswer   = false

	// three months
	DefaultExpiration = time.Hour * 24 * 30 * 3
)

var (
	errEncodeQr = errors.New("could not encode data to qr code")
)

type Bitmap [][]bool

type Quiz struct {
	ID        string
	Title     string
	Initial   Bitmap
	Questions []Question
	Expires   time.Time
}

func NewQuiz(title string, secret string, questions []Question) (Quiz, error) {
	var quiz Quiz

	quiz.ID = cuid2.Generate()
	quiz.Title = title
	quiz.Expires = time.Now().Add(DefaultExpiration)

	puzzle, err := NewPuzzle(secret, questions)
	if err != nil {
		return quiz, err
	}

	quiz.Initial = puzzle.Initial
	quiz.Questions = puzzle.Questions

	return quiz, nil
}

type Question struct {
	Question string
	Hint     string
	Answers  []Answer
}

type Answer struct {
	Text    string
	Pixel   Pixel
	Correct bool
}

type Pixel struct {
	X int
	Y int
}
