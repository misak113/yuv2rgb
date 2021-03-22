package yuv2rgb

import (
	/*
		#include <yuv_rgb.c>
	*/
	"C"

	"image"
	"reflect"
	"unsafe"
)

// ConvertRGBAToYCbCrImageSSEUnaligned will use SSE intel API unaligned to convert RGBA image into YCbCr
func ConvertRGBAToYCbCrImageSSEUnaligned(
	rgba *image.RGBA,
) (*image.YCbCr, error) {
	bounds := rgba.Bounds()
	ycbcr := image.NewYCbCr(bounds, image.YCbCrSubsampleRatio420)
	convertRGB32ToYUV420SSEUnaligned(
		uint(bounds.Max.X),
		uint(bounds.Max.Y),
		rgba.Pix,
		uint(rgba.Stride),
		ycbcr.Y,
		ycbcr.Cb,
		ycbcr.Cr,
		uint(ycbcr.YStride),
		uint(ycbcr.CStride),
		YCbCrType601,
	)
	return ycbcr, nil
}

type cRgbToYuvType func(C.uint, C.uint, *C.uchar, C.uint, *C.uchar, *C.uchar, *C.uchar, C.uint, C.uint, C.YCbCrType)

// convertRGB32ToYUV420SSEUnaligned will use SSE intel API unaligned to convert color space of RGB 32 bit data into YUV420
func convertRGB32ToYUV420SSEUnaligned(
	width uint,
	height uint,
	rgbBuf []byte,
	rgbStride uint,
	yBuf []byte,
	uBuf []byte,
	vBuf []byte,
	yStride uint,
	cStride uint,
	ycbcrType ycbcrType,
) error {
	// TODO check lengths of buffers
	rgbSliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&rgbBuf))
	ySliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&yBuf))
	uSliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&uBuf))
	vSliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&vBuf))
	C.rgb32_yuv420_sseu(
		C.uint(width),
		C.uint(height),
		(*C.uchar)(unsafe.Pointer(rgbSliceHeader.Data)),
		C.uint(rgbStride),
		(*C.uchar)(unsafe.Pointer(ySliceHeader.Data)),
		(*C.uchar)(unsafe.Pointer(uSliceHeader.Data)),
		(*C.uchar)(unsafe.Pointer(vSliceHeader.Data)),
		C.uint(yStride),
		C.uint(cStride),
		C.YCbCrType(ycbcrType),
	)
	return nil
}
