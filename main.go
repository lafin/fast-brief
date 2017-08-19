package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/lafin/fast"
)

func main() {
	img, _ := gg.LoadImage("image_1.png")

	// Create a new grayscale image
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	gray := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			oldColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(oldColor)
			gray.Set(x, y, grayColor)
		}
	}

	pixList := make([]int, w*h)
	for index := 0; index < w*h; index++ {
		pixList[index] = int(gray.Pix[index])
	}

	cornerList := fast.FindCorners(pixList, w, h, 20)
	dc := gg.NewContextForImage(img)
	for i := 0; i < len(cornerList); i += 2 {
		dc.DrawCircle(float64(cornerList[i]), float64(cornerList[i+1]), 2)
	}
	dc.SetHexColor("#0000FF")
	dc.Fill()
	dc.SavePNG("image_2.png")

	fmt.Println("done")
}
