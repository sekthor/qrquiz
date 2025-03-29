# QR Codes for self-assembly

> A silly app, that corrupts a QR code containing a secret.
> To reveal it, you must reassemble the QR by correctly solving a quiz.

## How it works

This app...

- allows you to create a quiz
  - with multiple choice answers
  - and a secret (the quiz's solution)
- encodes the secret in a QR code
- corrupts the QR code by "stealing" some of it's pixels and assigning them to correct answers
- assignes not-needed pixels to wrong answers

This means:

- giving right answers, "reassembles" the QR code
- wrong answers further corrupt the QR code

## Credits

- QR Code SVG: [openmoji.org](https://openmoji.org/) ([CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/))
