package domain

import (
	"reflect"
	"testing"

	"github.com/skip2/go-qrcode"
)

func Test_assignPixels(t *testing.T) {

	qr, _ := qrcode.New("hello", qrcode.Low)
	helloBitmap := qr.Bitmap()

	type args struct {
		questions []Question
		QR        Bitmap
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				QR: helloBitmap,
				questions: []Question{
					{
						Question: "What is your name",
						Answers: []Answer{
							{
								Text:    "Sir Lancelot",
								Correct: true,
							},
							{
								Text:    "Sir Robin the Brave",
								Correct: false,
							},
						},
					},
					{
						Question: "What is your Quest",
						Answers: []Answer{
							{
								Text:    "To seek the holy grail",
								Correct: true,
							},
							{
								Text:    "To find a shrubbery",
								Correct: false,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := assignPixels(tt.args.questions, tt.args.QR)
			if (err != nil) != tt.wantErr {
				t.Errorf("assignPixels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.QR(), tt.args.QR) {
				t.Errorf("assignPixels() = %v, want %v", got, tt.args.QR)
			}
		})
	}
}
