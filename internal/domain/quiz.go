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
	ModuleSize        = 1
	PositionSize      = 7 * ModuleSize
	DefaultExpiration = time.Hour * 24 * 30 * 3 // three months
)

var (
	errEncodeQr = errors.New("could not encode data to qr code")
)

type Bitmap [][]bool

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

type Answer struct {
	ID         string  `json:"-" gorm:"primaryKey"`
	QuestionID string  `json:"-"`
	Text       string  `json:"text"`
	Pixels     []Pixel `json:"pixels"`
	Correct    bool    `json:"correct"`
}

type Pixel struct {
	ID       string `json:"-" gorm:"primaryKey"`
	AnswerID string `json:"-"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) (err error) {
	q.ID = cuid2.Generate()
	return
}

func (a *Answer) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = cuid2.Generate()
	return
}

func (p *Pixel) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = cuid2.Generate()
	return
}

// omit the correct field on marshaling, but not on unmarshaling
func (a Answer) MarshalJSON() ([]byte, error) {
	helper := struct {
		Text   string  `json:"text"`
		Pixels []Pixel `json:"pixels"`
	}{
		Text:   a.Text,
		Pixels: a.Pixels,
	}
	return json.Marshal(helper)
}

// TODO: implement the valuer & scanner interface in a less
// "verbose" and more performant way than JSON
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
