package main

import (
	"image"
	"image/color"
	"strings"
)

func CreateImage(time string, size int) *image.RGBA {
	colonWidth := 4
	digitWidth := 8
	digitHeight := 12

	imgWidth := (6*digitWidth + 2*colonWidth) * size
	imgHeight := digitHeight * size
	upLeft := image.Point{X: 0, Y: 0}
	lowRight := image.Point{X: imgWidth, Y: imgHeight}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	offset := 0
	digit := ""

	for _, char := range time {
		switch string(char) {
		case "0":
			digit = Zero
		case "1":
			digit = One
		case "2":
			digit = Two
		case "3":
			digit = Three
		case "4":
			digit = Four
		case "5":
			digit = Five
		case "6":
			digit = Six
		case "7":
			digit = Seven
		case "8":
			digit = Eight
		case "9":
			digit = Nine
		default:
			digit = Colon
		}

		drawDigit(img, digit, size, offset)

		if digit == Colon {
			offset += colonWidth
		} else {
			offset += digitWidth
		}
	}

	return img
}

func drawDigit(img *image.RGBA, digit string, k int, offset int) {
	lines := strings.Split(digit, "\n")

	for h := 0; h < len(lines); h++ {
		for w, char := range lines[h] {
			fillCell(img, w+offset, h, k, string(char))
		}
	}
}

func fillCell(img *image.RGBA, x int, y, size int, char string) {
	for i := x * size; i < (x+1)*size; i++ {
		for j := y * size; j < (y+1)*size; j++ {
			if char != "." {
				img.Set(i, j, Cyan)
			} else {
				img.Set(i, j, color.RGBA{R: 255, G: 255, B: 255, A: 255})
			}
		}
	}
}
