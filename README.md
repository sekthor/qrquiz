# qrquiz

I have a lot of silly ideas.
This is one of them:

> "Solve a quiz, where the answers you give affect the pixels of a QR code. If it scans, you solved it correctly."

This app

- allows you to create a quiz
  - with multiple choice answers
  - and a secret (the quiz's solution)
- encodes the secret in a QR code
- corrupts the QR code by "stealing" some of it's pixels and assigning them to correct answers
- assignes not-needed pixels to wrong answers

this means

- giving right answers, "reassembles" the QR code
- wrong answers further corrupt the QR code

## Credits

- QR Code SVG: [openmoji.org](https://openmoji.org/) ([CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/))
