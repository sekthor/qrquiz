package domain

import "errors"

var (
	ErrTooManyAnswers = errors.New("there are more answers in the quiz, than available pixels in the QR code")
)
