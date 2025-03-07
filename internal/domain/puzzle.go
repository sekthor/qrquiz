package domain

import (
	"math"
	"math/rand"

	"github.com/skip2/go-qrcode"
)

// A Puzzle corrupts a given QR code.
// All correct answers are needed to recreate it with the missing pixels.
// Fault-Tolerance of QRs is kind of messing with this concept.
// On "Low" RecoveryLevel we can still expect 7% recovery
type Puzzle struct {
	Questions []Question
	Initial   Bitmap
}

type EligiblePixels struct {
	Set   []Pixel
	Unset []Pixel
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

	// get a bitmap without the quiet zone
	qr.DisableBorder = true
	bitmap := qr.Bitmap()

	puzzle, err = assignPixels(questions, bitmap)
	if err != nil {
		return puzzle, err
	}

	return puzzle, nil
}

func getEligiblePixels(bitmap Bitmap) EligiblePixels {
	var pixels EligiblePixels

	cornerSize := PositionSize + ModuleSize

	for y, row := range bitmap {
		for x, isPixelSet := range row {
			// ingore pixels on position corners
			if (y < cornerSize && x < cornerSize) || // top left
				(y < cornerSize && x > len(row)-cornerSize-1) || // top right
				(y > len(bitmap)-cornerSize-1 && x < cornerSize) { // bottom left
				continue
			}

			if isPixelSet {
				pixels.Set = append(pixels.Set, Pixel{X: x, Y: y})
			} else {
				pixels.Unset = append(pixels.Unset, Pixel{X: x, Y: y})
			}
		}
	}
	return pixels
}

// Correct answers are assigned a set pixel from the qr code.
// Wrong answers are assigned an unset pixel.
// From the QR code, we "subtract" all pixels, that were chosen for correct answers.
// This creates our initial Puzzle QR.
// If we add all pixels from all correct answers, we can recreate the original QR.
// Wrong answers will result set pixels, that are supposed to be unset, further
// corrupting the QR.
func assignPixels(questions []Question, bitmap Bitmap) (Puzzle, error) {
	var puzzle Puzzle

	// get set/unset pixels, but ignore the position markers in the corner
	pixels := getEligiblePixels(bitmap)
	pixels.Shuffle()

	// we need to count before assigning to figure out
	// how many pixels to assign per answer
	var correctCount, wrongCount = 0, 0
	for _, question := range questions {
		for _, answer := range question.Answers {
			if answer.Correct {
				correctCount += 1
			} else {
				wrongCount += 1
			}
		}
	}

	// will be filled with all pixels assigned to correct answers.
	// these will be removed from the original QR bitmap
	var subtractMask []Pixel

	// determine the amount of pixels per answer. Since we choose set pixels for
	// correct answers and unset for wrong ones, we distribute the pixels to the
	// respective answer amount. The amount per answer is the smaller of the results
	// as this will work for both.
	pixelsPerAnswer := int(math.Min(
		float64(len(pixels.Set)/correctCount),
		float64(len(pixels.Unset)/wrongCount),
	))

	if pixelsPerAnswer < 1 {
		return puzzle, ErrTooManyAnswers
	}

	var setCursor, unsetCursor = 0, 0
	for i, question := range questions {
		for j, answer := range question.Answers {
			if answer.Correct {
				questions[i].Answers[j].Pixels = pixels.Set[setCursor : setCursor+pixelsPerAnswer]
				subtractMask = append(subtractMask, pixels.Set[setCursor:setCursor+pixelsPerAnswer]...)
				setCursor += pixelsPerAnswer
			} else {
				questions[i].Answers[j].Pixels = pixels.Unset[unsetCursor : unsetCursor+pixelsPerAnswer]
				unsetCursor += pixelsPerAnswer
			}
		}
	}

	puzzle.Questions = questions
	puzzle.Initial = subtract(bitmap, subtractMask)

	return puzzle, nil
}

// deletes set pixels from a bitmap from a given list of pixels
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

// Recreate the corrupted QR code
func (p Puzzle) QR() Bitmap {
	for _, q := range p.Questions {
		for _, a := range q.Answers {
			if a.Correct {
				for _, pixel := range a.Pixels {
					p.Initial[pixel.Y][pixel.X] = true
				}
			}
		}
	}
	return p.Initial
}

func (e *EligiblePixels) Shuffle() {
	e.Set = shuffle(e.Set)
	e.Unset = shuffle(e.Unset)
}
