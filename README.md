# yuv2rgb
Go and C library for fast image conversion between yuv420p and rgb24.

This is a simple library for optimized image conversion between YUV420p and rgb24.
It was done mainly as an exercise to learn to use sse instrinsics, so there may still be room for optimization.

For each conversion, a standard c optimized function and two sse function (with aligned and unaligned memory) are implemented.
The sse version requires only SSE2, which is available on any reasonnably recent CPU.
The library also supports the three different YUV (YCrCb to be correct) color spaces that exist (see comments in code), and others can be added simply.

There is a simple test program, that convert a raw YUV file to rgb ppm format, and measure computation time.
Optionnaly, it also compares the result and computation time with the ffmpeg implementation (that uses MMX), and with the IPP functions.

## GoLang target
### Usage
```sh
go get github.com/misak113/yuv2rgb
```

```go
import (
    yuv2rgb "github.com/misak113/yuv2rgb/go"
)

func main() {
    // ...
    ycbcrImg, err := yuv2rgb.ConvertImageToYCbCrImage(rgbaImg)
}
```

> You have to add following environment variables configuration to your `go build` compilation.

> IPP (of oneapi) is usually installed in different than standard location `/opt/intel/oneapi/ipp/latest/`. So GoLang compiler cannot find it by default.

> It should be compiled with static linking to prevent user's of application to have installed oneapi which is kind a large.

> Also, see [IPP (Intel Integrated Performance Primitives)](#IPP (Intel Integrated Performance Primitives)) notes below.

```sh
export CGO_CFLAGS="-g -O2 -I/opt/intel/oneapi/ipp/latest/include"
export CGO_LDFLAGS="-Wl,-Bstatic -lippcc -lippcore -g -O2 -L/opt/intel/oneapi/ipp/latest/lib/intel64 -Wl,-Bdynamic"
```

### Test
To run tests, simply do:
```sh
cd go
make test
```

To run benchmarks, simply do:
```sh
cd go
make bench
```
#### benchmarks output
```
goos: linux
goarch: amd64
BenchmarkConvertRGBAToYCbCrImageStandardSmall-4             4202            285576 ns/op
BenchmarkConvertRGBAToYCbCrImageSSEUnalignedSmall-4        16311             70406 ns/op
BenchmarkConvertRGBAToYCbCrImageIPPSmall-4                165366              6579 ns/op
BenchmarkConvertRGBAToYCbCrImageStandardMedium-4               3         405089247 ns/op
BenchmarkConvertRGBAToYCbCrImageSSEUnalignedMedium-4           6         189779372 ns/op
BenchmarkConvertRGBAToYCbCrImageIPPMedium-4             1000000000               1.03 ns/op
BenchmarkConvertRGBAToYCbCrImageStandardLarge-4                3         444012528 ns/op
BenchmarkConvertRGBAToYCbCrImageSSEUnalignedLarge-4            7         166926969 ns/op
BenchmarkConvertBGRAToYCbCrImageIPPLarge-4                    69          15423310 ns/op
```

## C
To compile, simply do :

    mkdir build
    cd build
    cmake -DCMAKE_BUILD_TYPE=Release ..
    make

The test program only support raw YUV files for the YUV420 format, and ppm for the RGB24 format.
To generate a raw yuv file, you can use ffmpeg:

    ffmpeg -i example.jpg -c:v rawvideo -pix_fmt yuv420p example.yuv

To generate the rgb file, you can use the ImageMagick convert program:

    convert example.jpg example.ppm

Then, for YUV420 to RGB24 conversion, use the test program like that:

    ./test_yuv_rgb yuv2rgb image.yuv 4096 2160 image
  
The second and third parameters are image width and height (that are needed because not available in the raw YUV file), and fourth parameter is the output filename template (several output files will be generated, named for example output_sse.ppm, output_av.ppm, etc.)

Similarly, for RGB24 to YUV420 conversion:

    ./test_yuv_rgb rgb2yuv image.ppm image

On my computer, the test program on a 4K image give the following for yuv2rgb:

    Time will be measured in each configuration for 100 iterations...
    Processing time (std) : 2.630193 sec
    Processing time (sse2_unaligned) : 0.704394 sec
    Processing time (ffmpeg_unaligned) : 1.221432 sec
    Processing time (ipp_unaligned) : 0.636274 sec
    Processing time (sse2_aligned) : 0.606648 sec
    Processing time (ffmpeg_aligned) : 1.227100 sec
    Processing time (ipp_aligned) : 0.636951 sec

And for rgb2yuv:

    Time will be measured in each configuration for 100 iterations...
    Processing time (std) : 2.588675 sec
    Processing time (sse2_unaligned) : 0.676625 sec
    Processing time (ffmpeg_unaligned) : 3.385816 sec
    Processing time (ipp_unaligned) : 0.593890 sec
    Processing time (sse2_aligned) : 0.640630 sec
    Processing time (ffmpeg_aligned) : 3.397952 sec
    Processing time (ipp_aligned) : 0.579043 sec

configuration : gcc 4.9.2, swscale 3.0.0, IPP 9.0.1, intel i7-5500U

## FFMPEG
Using ffmpeg implementation, compile with `USE_FFMPEG=1` cflag:

```sh
cmake -DUSE_FFMPEG=1 -DCMAKE_BUILD_TYPE=Release ..
```

The requirements is to have installed `ffmpeg`. Debian installation:
```sh
sudo apt install ffmpeg
```

## IPP (Intel Integrated Performance Primitives)
Using IPP implementation, compile with `USE_IPP=1` and optionally `IPP_ROOT=/opt/intel/oneapi/ipp/latest` cflags:

```sh
cmake -DUSE_IPP=1 -DIPP_ROOT=/opt/intel/oneapi/ipp/latest -DCMAKE_BUILD_TYPE=Release ..
```

And then, don't forget to add path of libraries to `LD_LIBRARY_PATH`.

```sh
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/opt/intel/oneapi/ipp/latest/lib/intel64
```

The requirements is to have installed `intel-basekit`. See debian manual here: https://software.intel.com/content/www/us/en/develop/articles/installing-intel-oneapi-toolkits-via-apt.html

## Combined
You can combine more `USE_` cflags together.
