package yuv2rgb

import (
	"image"
	"image/color"

	"github.com/misak113/yuv2rgb/go/imgext"
)

func init() {
	bgra2yuvFns[priorityLow] = append(bgra2yuvFns[priorityLow], ConvertBGRAToYCbCrImageStandard)
	rgba2yuvFns[priorityLow] = append(rgba2yuvFns[priorityLow], ConvertRGBAToYCbCrImageStandard)
}

// ConvertRGBAToYCbCrImageStandard will use Go native color.RGBToYCbCr()
func ConvertRGBAToYCbCrImageStandard(
	rgba *image.RGBA,
) (*image.YCbCr, error) {
	return convertImageToYCbCrImageStandard(rgba)
}

// ConvertBGRAToYCbCrImageStandard will use Go native color.RGBToYCbCr()
func ConvertBGRAToYCbCrImageStandard(
	bgra *imgext.BGRA,
) (*image.YCbCr, error) {
	return convertImageToYCbCrImageStandard(bgra)
}

// convertImageToYCbCrImageStandard will use Go native color.RGBToYCbCr() of image.RGBA or image.BGRA
func convertImageToYCbCrImageStandard(
	img image.Image,
) (*image.YCbCr, error) {
	bounds := img.Bounds()
	ycbcr := image.NewYCbCr(bounds, image.YCbCrSubsampleRatio420)

	for row := 0; row < bounds.Max.Y; row++ {
		for col := 0; col < bounds.Max.X; col++ {
			r, g, b, _ := img.At(col, row).RGBA()
			y, cb, cr := color.RGBToYCbCr(uint8(r), uint8(g), uint8(b))

			ycbcr.Y[ycbcr.YOffset(col, row)] = y
			ycbcr.Cb[ycbcr.COffset(col, row)] = cb
			ycbcr.Cr[ycbcr.COffset(col, row)] = cr
		}
	}
	return ycbcr, nil
}
