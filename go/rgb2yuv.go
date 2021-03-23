package yuv2rgb

import (
	"fmt"
	"image"
	"image/color"
)

// ConvertRGBAToYCbCrImageType is func type of ConvertRGBAToYCbCrImage
type ConvertRGBAToYCbCrImageType func(rgba *image.RGBA) (*image.YCbCr, error)

// ConvertRGBAToYCbCrImage will try use the most fast algorithm
// first try IPP intel API unaligned to convert RGBA image into YCbCr
// second try SSE intel API unaligned to convert RGBA image into YCbCr
// if above algorithm not available, Go native color.RGBToYCbCr() is used
func ConvertRGBAToYCbCrImage(
	rgba *image.RGBA,
) (*image.YCbCr, error) {
	// TODO implement compile time detection of the best algorithm
	return nil, fmt.Errorf("Not supported")
}

// ConvertRGBAToYCbCrImageStandard will use Go native color.RGBToYCbCr()
func ConvertRGBAToYCbCrImageStandard(
	rgba *image.RGBA,
) (*image.YCbCr, error) {
	bounds := rgba.Bounds()
	ycbcr := image.NewYCbCr(bounds, image.YCbCrSubsampleRatio420)

	for row := 0; row < bounds.Max.Y; row++ {
		for col := 0; col < bounds.Max.X; col++ {
			r, g, b, _ := rgba.At(col, row).RGBA()
			y, cb, cr := color.RGBToYCbCr(uint8(r), uint8(g), uint8(b))

			ycbcr.Y[ycbcr.YOffset(col, row)] = y
			ycbcr.Cb[ycbcr.COffset(col, row)] = cb
			ycbcr.Cr[ycbcr.COffset(col, row)] = cr
		}
	}
	return ycbcr, nil
}
