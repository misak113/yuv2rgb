package imgext

import "image"

// ConvertRGBAToBGRA converts *image.RGBA into *imgext.BGRA
func ConvertRGBAToBGRA(rgba *image.RGBA) *BGRA {
	bgra := &BGRA{
		Rect:   rgba.Rect,
		Stride: rgba.Stride,
		Pix:    make([]byte, len(rgba.Pix)),
	}
	for y := 0; y < rgba.Rect.Max.Y; y++ {
		for x := 0; x < rgba.Rect.Max.X; x++ {
			bgra.Set(x, y, rgba.At(x, y))
		}
	}
	return bgra
}
