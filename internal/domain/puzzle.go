package domain

import (
	"errors"
	"math/rand"

	"github.com/skip2/go-qrcode"
)

var (
	ErrorNotEnoughPixels = errors.New("not enough pixels for the amount of answers in quiz")
)

// A Puzzle corrupts a given QR code.
// All correct answers are needed to recreate it with the missing pixels.
// Fault-Tolerance of QRs is kind of messing with this concept.
// On "Low" RecoveryLevel we can still expect 7% recovery
type Puzzle struct {
	Questions []Question
	Initial   Bitmap
}

// Recreate the corrupted QR code
func (p Puzzle) QR() Bitmap {
	for _, q := range p.Questions {
		for _, a := range q.Answers {
			if a.Correct {
				p.Initial[a.Pixel.Y][a.Pixel.X] = true
			}
		}
	}
	return p.Initial
}

// Creates a Puzzle by encoding a secret in a QR Code, and then subtracting
// black pixels for every correct answer in a list of questions.
// It assignes these pixels to the corresponding correct answer.
// Wrong answers are assigned a blank pixel from the QR code.
func NewPuzzle(secret string, questions []Question) (Puzzle, error) {
	var puzzle Puzzle
	qr, err := qrcode.New(secret, qrcode.Low)
	if err != nil {
		return puzzle, errEncodeQr
	}

	QR := qr.Bitmap()
	puzzle, err = assignPixels(questions, QR)
	if err != nil {
		return puzzle, err
	}

	return puzzle, nil
}

// Correct answers are assigned a black pixel from the qr code.
// Wrong answers are assigned a white pixel.
// From the QR code, we "subtract" all black pixels, that were chosen for correct answers.
// This creates our initial Puzzle QR, which need the pixels of correct answers to be added
// to recreate the original QR.
// Wrong answers will result add black pixels, where none should be.
// DISCLAIMER: too many answers may result in "running out of pixels"
func assignPixels(questions []Question, QR Bitmap) (Puzzle, error) {
	var puzzle Puzzle

	black, white := groupPixelsByColor(QR, false)

	// make sure the pixels are in random order
	black = shuffle(black)
	white = shuffle(white)

	var correctCount, wrongCount = 0, 0
	var subtractMask []Pixel

	for i := 0; i < len(questions); i++ {
		for j := 0; j < len(questions[i].Answers); j++ {
			if questions[i].Answers[j].Correct {
				pixel := black[correctCount]
				subtractMask = append(subtractMask, pixel)
				questions[i].Answers[j].Pixel = pixel
				correctCount++
			} else {
				questions[i].Answers[j].Pixel = white[wrongCount]
				wrongCount++
			}

			if correctCount > len(black) || wrongCount > len(white) {
				return puzzle, ErrorNotEnoughPixels
			}
		}
	}

	puzzle.Questions = questions
	puzzle.Initial = subtract(QR, subtractMask)

	return puzzle, nil
}

// deletes black pixels from a bitmap from a given list of pixels
func subtract(qr Bitmap, mask []Pixel) Bitmap {
	for _, pixel := range mask {
		qr[pixel.Y][pixel.X] = false
	}
	return qr
}

// Fisher-Yates Shuffle for Pixel slice.
// Shuffles a slice of Pixels into random order.
// Returns a shuffled copy of the slice.
func shuffle(pixels []Pixel) []Pixel {
	shuffled := append([]Pixel(nil), pixels...)
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}

// splits a bitmap into two slices of Pixels. One for black and one for white pixels.
func groupPixelsByColor(bitmap [][]bool, quietzone bool) (black []Pixel, white []Pixel) {
	offset := 0

	if !quietzone {
		bitmap, offset = removeQuietZone(bitmap)
	}

	for y, row := range bitmap {
		for x, isBlack := range row {
			pixel := Pixel{X: x + offset, Y: y + offset}
			if isBlack {
				black = append(black, pixel)
			} else {
				white = append(white, pixel)
			}
		}
	}
	return
}

func detectQuietZoneSize(bitmap Bitmap) int {
	for i, row := range bitmap {
		for _, pixel := range row {
			if pixel {
				return i
			}
		}
	}
	return len(bitmap) / 2
}

func removeQuietZone(bitmap Bitmap) (Bitmap, int) {
	size := detectQuietZoneSize(bitmap)
	var new Bitmap
	for _, row := range bitmap[size : len(bitmap)-size] {
		new = append(new, row[size:len(row)-size])
	}
	return new, size
}
