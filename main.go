package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"
	"os"
	"time"

	"github.com/disintegration/imaging"
)

func loadImage(path string) image.Image {
	reader, err := os.Open(path)
	if err != nil {
		fmt.Print(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		fmt.Print(err)
	}

	return img
}

func randPixel() uint8 {
	return uint8(rand.Intn(255))
}

func generateLine(img *image.Image, newImg *image.RGBA) {
	maxX := newImg.Bounds().Size().X
	maxY := newImg.Bounds().Size().Y

	start := rand.Intn(maxY)
	end := start + rand.Intn(200)

	offset := rand.Intn(100) - 50
	offsetY := rand.Intn(50) - 25
	yintensity := 1.0 + (rand.Float32() * 0.4)
	cbSat := 0.8 + rand.Float32()*0.4
	crSat := 0.8 + rand.Float32()*0.4

	flip := rand.Float32() > 0.5

	for y := start; y < end; y++ {
		for x := 0; x < maxX; x++ {
			col := (*img).At(x+offset, y+offsetY).(color.YCbCr)
			col.Y = uint8(float32(col.Y) * yintensity)

			if flip {
				cr := col.Cr
				col.Cr = col.Cb
				col.Cb = cr
			}

			col.Cb = uint8(float32(col.Cb) * cbSat)
			col.Cr = uint8(float32(col.Cb) * crSat)

			newImg.Set(x, y, col)
		}
	}
}

func main() {
	img := loadImage(os.Args[1])
	newImg := image.NewRGBA(img.Bounds())

	maxX := newImg.Bounds().Size().X
	maxY := newImg.Bounds().Size().Y

	// Seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	// start := rand.Intn(maxY)
	// end := start + rand.Intn(maxY-start)

	cbOffset := rand.Intn(100) - 50
	crOffset := rand.Intn(100) - 50
	breakPos := rand.Intn(maxY)

	for y := 0; y < maxY; y++ {
		messup := rand.Float32() > 0.4

		for x := 0; x < maxX; x++ {
			col := img.At(x, y).(color.YCbCr)

			if messup {
				newCol := img.At(x+cbOffset+rand.Intn(50), y).(color.YCbCr)
				col.Cb = newCol.Cb

				newCol = img.At(x+crOffset+rand.Intn(10), y).(color.YCbCr)
				col.Cr = newCol.Cr

				newCol = img.At(x+rand.Intn(5), y).(color.YCbCr)
				col.Y = newCol.Y
			}

			if y > breakPos {
				newCol := img.At(x+cbOffset-crOffset, y).(color.YCbCr)
				col.Y = newCol.Y
			}

			newImg.Set(x, y, col)
		}
	}

	for i := 0; i < rand.Intn(5)+1; i++ {
		generateLine(&img, newImg)
	}

	processedImg := imaging.AdjustSaturation(newImg, 40)
	processedImg = imaging.Sharpen(processedImg, 5)

	file, err := os.Create("./out.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	jpeg.Encode(file, processedImg, nil)
}
