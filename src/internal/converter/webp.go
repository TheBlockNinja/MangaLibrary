package converter

import (
	"bufio"
	"image"
	"image/color"
	"image/png"
	"os"

	"golang.org/x/image/webp"
)

func ConvertNYCbCrA(Color color.NYCbCrA) color.RGBA {
	r, g, b := color.YCbCrToRGB(Color.Y, Color.Cb, Color.Cr)
	return color.RGBA{r, g, b, 255}
}

func ConvertNRGBA(Color color.NRGBA) color.RGBA {
	return color.RGBA{
		Color.R,
		Color.G,
		Color.B,
		255,
	}
}

/************************/
func ConvertWebPToPng(from, to string) error{
	//open the WebP file
	webpFile, _ := os.OpenFile(from, os.O_RDWR, 0777)
	defer webpFile.Close()

	//create a reader
	reader := bufio.NewReader(webpFile)

	//decode the WebP
	imageData, _ := webp.Decode(reader)
	webpFile.Close()

	//create new image
	newImage := image.NewRGBA(imageData.Bounds())

	//fill new image with pixels
	for y := imageData.Bounds().Min.Y; y < imageData.Bounds().Max.Y; y++ {
		for x := imageData.Bounds().Min.X; x < imageData.Bounds().Max.X; x++ {

			//get pixel from imageData
			pixel := imageData.At(x, y)

			//convert pixel to RGBA
			var RGBApixel color.RGBA
			switch imageData.ColorModel() {

			case color.NYCbCrAModel:
				RGBApixel = ConvertNYCbCrA(pixel.(color.NYCbCrA))

			case color.NRGBAModel:
				RGBApixel = ConvertNRGBA(pixel.(color.NRGBA))

			default:
				RGBApixel = color.RGBAModel.Convert(pixel).(color.RGBA)
				RGBApixel.A = 255

			}

			//set new pixel in new image
			newImage.Set(x, y, RGBApixel)

		}
	}
	//create the new PNG file
	pngFile, _ := os.Create(to)
	defer pngFile.Close()

	//write to the PNG file
	return png.Encode(pngFile, newImage)
}
