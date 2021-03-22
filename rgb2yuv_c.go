package yuv2rgb

import (
	"C"
)

type ycbcrType int

// YCbCrType enum of conversion type
const (
	YCbCrTypeJpeg ycbcrType = 0
	YCbCrType601  ycbcrType = 1
	YCbCrType709  ycbcrType = 2
)
