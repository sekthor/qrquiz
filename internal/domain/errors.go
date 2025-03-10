package domain

import "errors"

var (
	ErrTooManyAnswers = errors.New("there are more answers in the quiz, than available pixels in the QR code")
	ErrNoAnswers      = errors.New("quiz has no answers")
	ErrEncodeQr       = errors.New("could not encode data to qr code")
)
