package yuv2rgb

import (
	"fmt"
	"image"
	"reflect"

	"github.com/misak113/yuv2rgb/go/imgext"
)

type priority int

const (
	priorityLow    priority = 0
	priorityMedium priority = 1
	priorityHigh   priority = 2
)

type bgraToYCbCrFnType = func(bgra *imgext.BGRA) (*image.YCbCr, error)
type rgbaToYCbCrFnType = func(rgba *image.RGBA) (*image.YCbCr, error)

var bgra2yuvFns = map[priority][]bgraToYCbCrFnType{}
var rgba2yuvFns = map[priority][]rgbaToYCbCrFnType{}

// ConvertImageToYCbCrImage will try use the most fast algorithm for input image.RGBA or imgext.BGRA
// Use IPP intel API unaligned to convert image into YCbCr on amd64
// Use SSE intel API unaligned to convert image into YCbCr on amd64 if IPP not available
// if above algorithm not available, Go native standard color.RGBToYCbCr() is used
func ConvertImageToYCbCrImage(
	img image.Image,
) (*image.YCbCr, error) {
	switch typedImg := img.(type) {
	case *image.RGBA:
		return ConvertRGBAToYCbCrImage(typedImg)
	case *imgext.BGRA:
		return ConvertBGRAToYCbCrImage(typedImg)
	}
	return nil, fmt.Errorf("Not supported image.Image struct %v", reflect.TypeOf(img))
}

// ConvertRGBAToYCbCrImage see ConvertImageToYCbCrImage
func ConvertRGBAToYCbCrImage(
	rgba *image.RGBA,
) (*image.YCbCr, error) {
	for pri := priorityHigh; pri >= priorityLow; pri-- {
		if fns, exists := rgba2yuvFns[pri]; exists && len(fns) > 0 {
			return fns[0](rgba)
		}
	}

	return nil, fmt.Errorf("Not available any conversion function")
}

// ConvertBGRAToYCbCrImage see ConvertImageToYCbCrImage
func ConvertBGRAToYCbCrImage(
	bgra *imgext.BGRA,
) (*image.YCbCr, error) {
	for pri := priorityHigh; pri >= priorityLow; pri-- {
		if fns, exists := bgra2yuvFns[pri]; exists && len(fns) > 0 {
			return fns[0](bgra)
		}
	}

	return nil, fmt.Errorf("Not available any conversion function")
}
