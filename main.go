package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/fogleman/gg"
	"github.com/lafin/brief"
	"github.com/lafin/fast"
	"github.com/tajtiattila/blur"
)

func grayImageToPixList(gray *image.Gray, width, height int) []int {
	pixList := make([]int, width*height)
	for index := 0; index < width*height; index++ {
		pixList[index] = int(gray.Pix[index])
	}

	return pixList
}

func toGray(path string) (*image.Gray, int, int) {
	infile, err := os.Open(path)
	if err != nil {
		// replace this with real error handling
		panic(err)
	}
	defer infile.Close()

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	img, _, err := image.Decode(infile)
	if err != nil {
		// replace this with real error handling
		panic(err)
	}

	img = blur.Gaussian(img, 1, blur.ReuseSrc)

	// Create a new grayscale image
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	gray := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			gray.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}

	return gray, width, height
}

func main() {
	img1, _ := gg.LoadImage("fast_1.png")
	gray1, width1, height1 := toGray("fast_1.png")
	pixList1 := grayImageToPixList(gray1, width1, height1)
	corners1 := fast.FindCorners(pixList1, width1, height1, 40)

	dc := gg.NewContextForImage(img1)
	for i := 0; i < len(corners1); i += 2 {
		dc.DrawCircle(float64(corners1[i]), float64(corners1[i+1]), 2)
	}
	dc.SetHexColor("#0000FF")
	dc.Fill()
	dc.SavePNG("fast_2.png")

	gray2, width2, height2 := toGray("brief_1.png")
	gray3, width3, height3 := toGray("brief_2.png")
	pixList2 := grayImageToPixList(gray2, width2, height2)
	pixList3 := grayImageToPixList(gray3, width3, height3)

	randomWindowOffsets := brief.InitOffsets()
	corners2 := fast.FindCorners(pixList2, width2, height2, 40)
	descriptors2 := brief.GetDescriptors(pixList2, width2, corners2, randomWindowOffsets)
	corners3 := fast.FindCorners(pixList3, width3, height3, 40)
	descriptors3 := brief.GetDescriptors(pixList3, width3, corners3, randomWindowOffsets)

	matches := brief.ReciprocalMatch(corners2, descriptors2, corners3, descriptors3)
	for _, match := range matches {
		fmt.Println(match)
	}

	fmt.Println("done")
}
