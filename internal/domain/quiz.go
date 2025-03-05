package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/nrednav/cuid2"
	"gorm.io/gorm"
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

func (b Bitmap) Value() (driver.Value, error) {
	return json.Marshal(b)
}

func (b *Bitmap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to convert database value to []byte")
	}
	return json.Unmarshal(bytes, b)
}

type Quiz struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Title     string         `json:"title"`
	Initial   Bitmap         `json:"initial"`
	Questions []Question     `json:"questions"`
	Expires   time.Time      `json:"expires"`
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
	ID     string `json:"-" gorm:"primaryKey"`
	QuizID string `json:"-"`

	Question string   `json:"question"`
	Hint     string   `json:"hint"`
	Answers  []Answer `json:"answers"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) (err error) {
	q.ID = cuid2.Generate()
	return
}

type Answer struct {
	ID         string `json:"-" gorm:"primaryKey"`
	QuestionID string `json:"-"`
	Text       string `json:"text"`
	Pixel      Pixel  `json:"pixel"`
	Correct    bool   `json:"correct"`
}

func (a *Answer) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = cuid2.Generate()
	return
}

// omit the correct field on marshaling, but not on unmarshaling
func (a Answer) MarshalJSON() ([]byte, error) {
	helper := struct {
		Text  string `json:"text"`
		Pixel Pixel  `json:"pixel"`
	}{
		Text:  a.Text,
		Pixel: a.Pixel,
	}
	return json.Marshal(helper)
}

type Pixel struct {
	ID       string `gorm:"primaryKey"`
	AnswerID string `json:"-"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

func (p *Pixel) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = cuid2.Generate()
	return
}
