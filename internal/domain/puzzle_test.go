package domain

import (
	"testing"

	"github.com/skip2/go-qrcode"
)

func Test_getEligiblePixels(t *testing.T) {
	helloQR, _ := qrcode.New("hello", qrcode.Low)
	helloQR.DisableBorder = true
	bitmap := helloQR.Bitmap()

	var ineligible = []Pixel{}

	// compute pixels of position quares -> ineligible
	positionSize := PositionSize + ModuleSize
	qrCodeSize := len(bitmap)
	for y := 0; y < positionSize; y++ {
		for x := 0; x < positionSize; x++ {
			ineligible = append(ineligible, Pixel{X: x, Y: y})                  // top left
			ineligible = append(ineligible, Pixel{X: qrCodeSize - x - 1, Y: y}) // top right
			ineligible = append(ineligible, Pixel{X: x, Y: qrCodeSize - y - 1}) // bottom left
		}
	}

	type args struct {
		bitmap Bitmap
	}
	tests := []struct {
		name       string
		args       args
		ineligible []Pixel
	}{
		{
			name:       "ensure no ineligible pixels in EligiblePixels",
			ineligible: ineligible,
			args:       args{bitmap: bitmap},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getEligiblePixels(tt.args.bitmap)
			for _, pixel := range tt.ineligible {
				for _, set := range got.Set {
					if (set.X == pixel.X) && (set.Y == pixel.Y) {
						t.Errorf("set contains ineligible Pixel at {x=%d,y=%d}", set.X, set.Y)
					}
				}
				for _, unset := range got.Unset {
					if (unset.X == pixel.X) && (unset.Y == pixel.Y) {
						t.Errorf("unset contains ineligible Pixel at {x=%d,y=%d}", unset.X, unset.Y)
					}
				}
			}
		})
	}
}
