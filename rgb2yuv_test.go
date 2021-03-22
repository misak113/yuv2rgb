package yuv2rgb

import (
	"image"
	"os"
	"testing"

	"github.com/lmittmann/ppm"
	"github.com/stretchr/testify/require"
)

// change this to true if you'd like to generate output snapshots of data/ directory raw data of images
const writeSnapshots = false

// you can setup tolerance on output image colors
// it has very negative performance effect for tests
const tolerance = 0

type testCase struct {
	ppmFilePath string
	yuvFilePath string
}

func TestConvertRGBAToYCbCrImageStandard(t *testing.T) {
	cases := []testCase{
		{"./data/test1.ppm", "./data/test1_standard.yuv"},
		{"./data/test2.ppm", "./data/test2_standard.yuv"},
		{"./data/test3.ppm", "./data/test3_standard.yuv"},
	}
	for _, c := range cases {
		testConvertRGBAToYCbCrImageCase(t, c, ConvertRGBAToYCbCrImageStandard, nil, nil)
	}
}

func TestConvertRGBAToYCbCrImageSSEUnaligned(t *testing.T) {
	cases := []testCase{
		{"./data/test1.ppm", "./data/test1_sseunaligned.yuv"},
		{"./data/test2.ppm", "./data/test2_sseunaligned.yuv"},
		{"./data/test3.ppm", "./data/test3_sseunaligned.yuv"},
	}
	for _, c := range cases {
		testConvertRGBAToYCbCrImageCase(t, c, ConvertRGBAToYCbCrImageSSEUnaligned, nil, nil)
	}
}

func BenchmarkConvertRGBAToYCbCrImageStandardSmall(b *testing.B) {
	b.StopTimer()
	c := testCase{"./data/test1.ppm", "./data/test1_standard.yuv"}
	for i := 0; i < b.N; i++ {
		testConvertRGBAToYCbCrImageCase(b, c, ConvertRGBAToYCbCrImageStandard, b.StartTimer, b.StopTimer)
	}
}

func BenchmarkConvertRGBAToYCbCrImageSSEUnalignedSmall(b *testing.B) {
	c := testCase{"./data/test1.ppm", "./data/test1_sseunaligned.yuv"}
	for i := 0; i < b.N; i++ {
		testConvertRGBAToYCbCrImageCase(b, c, ConvertRGBAToYCbCrImageSSEUnaligned, b.StartTimer, b.StopTimer)
	}
}

func BenchmarkConvertRGBAToYCbCrImageStandardMedium(b *testing.B) {
	b.StopTimer()
	c := testCase{"./data/test2.ppm", "./data/test2_standard.yuv"}
	for i := 0; i < b.N; i++ {
		testConvertRGBAToYCbCrImageCase(b, c, ConvertRGBAToYCbCrImageStandard, b.StartTimer, b.StopTimer)
	}
}

func BenchmarkConvertRGBAToYCbCrImageSSEUnalignedMedium(b *testing.B) {
	c := testCase{"./data/test2.ppm", "./data/test2_sseunaligned.yuv"}
	for i := 0; i < b.N; i++ {
		testConvertRGBAToYCbCrImageCase(b, c, ConvertRGBAToYCbCrImageSSEUnaligned, b.StartTimer, b.StopTimer)
	}
}

func BenchmarkConvertRGBAToYCbCrImageStandardLarge(b *testing.B) {
	b.StopTimer()
	c := testCase{"./data/test3.ppm", "./data/test3_standard.yuv"}
	for i := 0; i < b.N; i++ {
		testConvertRGBAToYCbCrImageCase(b, c, ConvertRGBAToYCbCrImageStandard, b.StartTimer, b.StopTimer)
	}
}

func BenchmarkConvertRGBAToYCbCrImageSSEUnalignedLarge(b *testing.B) {
	c := testCase{"./data/test3.ppm", "./data/test3_sseunaligned.yuv"}
	for i := 0; i < b.N; i++ {
		testConvertRGBAToYCbCrImageCase(b, c, ConvertRGBAToYCbCrImageSSEUnaligned, b.StartTimer, b.StopTimer)
	}
}

func testConvertRGBAToYCbCrImageCase(
	t require.TestingT,
	c testCase,
	convertRGBAToYCbCrImageFn ConvertRGBAToYCbCrImageType,
	beforeFn func(),
	afterFn func(),
) {
	imgReader, err := os.Open(c.ppmFilePath)
	require.NoError(t, err)
	defer imgReader.Close()

	rgbaCfg, err := ppm.DecodeConfig(imgReader)
	require.NoError(t, err)
	imgReader.Seek(0, 0)

	rgbaImg, err := ppm.Decode(imgReader)
	require.NoError(t, err)

	//fmt.Println("RGB image", reflect.TypeOf(rgbaCfg.ColorModel), rgbaCfg.Width, rgbaCfg.Height)
	width := rgbaCfg.Width
	height := rgbaCfg.Height
	yStride := width           // width + (16-width%16)%16
	cStride := (width + 1) / 2 // (width+1)/2 + (16-((width+1)/2)%16)%16
	yLen := yStride * height
	cLen := cStride * ((height + 1) / 2) // ((width + 1) / 2) * ((height + 1) / 2)
	//fmt.Println("YUV image", yStride, cStride, yLen, cLen)

	if beforeFn != nil {
		beforeFn()
	}
	ycbcrImg, err := convertRGBAToYCbCrImageFn(rgbaImg.(*image.RGBA))
	if afterFn != nil {
		afterFn()
	}

	if writeSnapshots {
		ycbcrWriter, err := os.OpenFile(c.yuvFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		require.NoError(t, err)
		defer ycbcrWriter.Close()

		_, err = ycbcrWriter.Write(ycbcrImg.Y)
		require.NoError(t, err)
		_, err = ycbcrWriter.Write(ycbcrImg.Cb)
		require.NoError(t, err)
		_, err = ycbcrWriter.Write(ycbcrImg.Cr)
		require.NoError(t, err)
	}

	ycbcrReader, err := os.Open(c.yuvFilePath)
	require.NoError(t, err)
	defer ycbcrReader.Close()

	expectedY := make([]byte, yLen)
	_, err = ycbcrReader.Read(expectedY)
	require.NoError(t, err)

	expectedCb := make([]byte, cLen)
	_, err = ycbcrReader.Read(expectedCb)
	require.NoError(t, err)

	expectedCr := make([]byte, cLen)
	_, err = ycbcrReader.Read(expectedCr)
	require.NoError(t, err)

	require.Len(t, ycbcrImg.Y, len(expectedY))
	require.Len(t, ycbcrImg.Cb, len(expectedCb))
	require.Len(t, ycbcrImg.Cr, len(expectedCr))

	require.Equal(t, yStride, ycbcrImg.YStride)
	require.Equal(t, cStride, ycbcrImg.CStride)

	if tolerance == 0 {
		require.EqualValues(t, expectedY, ycbcrImg.Y)
		require.EqualValues(t, expectedCb, ycbcrImg.Cb)
		require.EqualValues(t, expectedCr, ycbcrImg.Cr)
	} else {
		for i := 0; i < len(expectedY); i++ {
			require.InDelta(t, expectedY[i], ycbcrImg.Y[i], tolerance, "Index %d of %d", i, len(expectedY))
		}
		for i := 0; i < len(expectedCb); i++ {
			require.InDelta(t, expectedCb[i], ycbcrImg.Cb[i], tolerance, "Index %d of %d", i, len(expectedCb))
		}
		for i := 0; i < len(expectedCr); i++ {
			require.InDelta(t, expectedCr[i], ycbcrImg.Cr[i], tolerance, "Index %d of %d", i, len(expectedCr))
		}
	}
}
