// +build amd64

package yuv2rgb

import (
	/*
		#cgo CFLAGS: -DUSE_IPP=1
		#cgo LDFLAGS: -lippcc -lippcore
		#include "../yuv_rgb.h"
	*/
	"C"

	"image"

	"github.com/misak113/yuv2rgb/go/imgext"
)
import (
	"reflect"
	"unsafe"
)

func init() {
	bgra2yuvFns[priorityHigh] = append(bgra2yuvFns[priorityHigh], ConvertBGRAToYCbCrImageIPP)
	// TODO Implement accelerated RGBA IPP
	// rgba2yuvFns[priorityHigh] = append(rgba2yuvFns[priorityHigh], ConvertRGBAToYCbCrImageIPP)
}

// ConvertBGRAToYCbCrImageIPP will use IPP intel API to convert imgext.BGRA image into image.YCbCr
func ConvertBGRAToYCbCrImageIPP(
	bgra *imgext.BGRA,
) (*image.YCbCr, error) {
	bounds := bgra.Bounds()
	ycbcr := image.NewYCbCr(bounds, image.YCbCrSubsampleRatio420)
	convertBGR32ToYUV420IPP(
		uint(bounds.Max.X),
		uint(bounds.Max.Y),
		bgra.Pix,
		uint(bgra.Stride),
		ycbcr.Y,
		ycbcr.Cb,
		ycbcr.Cr,
		uint(ycbcr.YStride),
		uint(ycbcr.CStride),
		YCbCrType601,
	)
	return ycbcr, nil
}

// convertRGB32ToYUV420IPP will use IPP intel API to convert color space of RGB 32 bit data into YUV420
func convertBGR32ToYUV420IPP(
	width uint,
	height uint,
	bgraBuf []byte,
	bgraStride uint,
	yBuf []byte,
	uBuf []byte,
	vBuf []byte,
	yStride uint,
	cStride uint,
	ycbcrType ycbcrType,
) error {
	// TODO check lengths of buffers
	bgraSliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bgraBuf))
	ySliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&yBuf))
	uSliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&uBuf))
	vSliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&vBuf))
	C.bgr32_yuv420_ipp(
		C.uint(width),
		C.uint(height),
		(*C.uchar)(unsafe.Pointer(bgraSliceHeader.Data)),
		C.uint(bgraStride),
		(*C.uchar)(unsafe.Pointer(ySliceHeader.Data)),
		(*C.uchar)(unsafe.Pointer(uSliceHeader.Data)),
		(*C.uchar)(unsafe.Pointer(vSliceHeader.Data)),
		C.uint(yStride),
		C.uint(cStride),
		C.YCbCrType(ycbcrType),
	)
	return nil
}
