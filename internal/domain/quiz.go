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
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Initial   Bitmap     `json:"initial"`
	Questions []Question `json:"questions"`
	Expires   time.Time  `json:"expires"`
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
	Question string   `json:"question"`
	Hint     string   `json:"hint"`
	Answers  []Answer `json:"answers"`
}

type Answer struct {
	Text    string `json:"text"`
	Pixel   Pixel  `json:"pixel"`
	Correct bool   `json:"-"`
}

type Pixel struct {
	X int `json:"x"`
	Y int `json:"y"`
}
