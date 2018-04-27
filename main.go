package main

import "C"

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/image/bmp"
)

const (
	hBMP byte = 0x42
	hGIF byte = 0x47
	hJPG byte = 0xff
	hPNG byte = 0x89
)

var once sync.Once
var imageMap = make(map[int]image.Image)

func main() {
	var mor string
	fmt.Println("Only supports calling as dynamic link library.")
	fmt.Print("?> ")
	fmt.Scanf("%s", &mor)
	if mor == "debug" {
		fmt.Println("DEBUG MODE::") // put some codes next here.
	}
}

/**
 * assist function
 */

func getImageObject(key int) image.Image {
	return imageMap[key]
}

func setImageObject(value image.Image) int {
	once.Do(func() { rand.Seed(time.Now().UnixNano()) })
	randi := rand.Int()
	for imageMap[randi] != nil {
		randi = rand.Int()
	}
	imageMap[randi] = value
	return randi
}

func delImageObject(key int) {
	delete(imageMap, key)
}

func getImgType(f *os.File) byte {
	tmp := make([]byte, 1)
	_, err := f.Read(tmp)
	check(err)
	f.Seek(0, os.SEEK_SET) // Reset pointer
	return tmp[0]
}

func getImgTypeSuffix(filename string) byte {
	switch strings.ToLower(filename[len(filename)-4:]) {
	case ".bmp":
		return hBMP
	case ".gif":
		return hGIF
	case ".jpg", "jpeg":
		return hJPG
	case ".png":
		return hPNG
	default:
		return 0x00
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

/**
 * core function
 */

func rgbaToGray(img image.Image) *image.Gray {
	var (
		bounds = img.Bounds()
		gray   = image.NewGray(bounds)
	)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			var rgba = img.At(x, y)
			gray.Set(x, y, rgba)
		}
	}
	return gray
}

func grayToRgba(img image.Image) *image.RGBA {
	var (
		bounds = img.Bounds()
		rgba   = image.NewRGBA(bounds)
	)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			var gray = img.At(x, y)
			rgba.Set(x, y, gray)
		}
	}
	return rgba
}

func moveBounds(img image.Image, left, top, right, bottom int, r, g, b, a uint8) image.Image {
	var (
		bounds = image.Rectangle{
			image.Point{0, 0},
			image.Point{img.Bounds().Dx() + left + right, img.Bounds().Dy() + top + bottom},
		}
		fillColor = color.RGBA{r, g, b, a}
	)
	if verifyGray(img) {
		// grayscale
		img2 := image.NewGray(bounds)
		for x := 0; x < bounds.Max.X; x++ {
			for y := 0; y < bounds.Max.Y; y++ {
				var ori color.Color
				if xi, yi := func() (int, int) { return x - left, y - top }(); xi >= 0 && xi < img.Bounds().Dx() && yi >= 0 && yi < img.Bounds().Dy() {
					ori = img.At(xi, yi)
				} else {
					ori = fillColor
				}
				img2.Set(x, y, ori)
			}
		}
		return img2
	} else {
		// non-grayscale
		img2 := image.NewRGBA(bounds)
		for x := 0; x < bounds.Max.X; x++ {
			for y := 0; y < bounds.Max.Y; y++ {
				var ori color.Color
				if xi, yi := func() (int, int) { return x - left, y - top }(); xi >= 0 && xi < img.Bounds().Dx() && yi >= 0 && yi < img.Bounds().Dy() {
					ori = img.At(xi, yi)
				} else {
					ori = fillColor
				}
				img2.Set(x, y, ori)
			}
		}
		return img2
	}
}

func verifyGray(img image.Image) bool {
	_, ok := img.(*image.Gray)
	_, ok2 := img.(*image.Paletted)
	return ok || ok2
}

/**
 * I/O function
 */

func file2Image(src string) image.Image {
	var (
		err error
		fin *os.File
		im  image.Image
	)
	fin, err = os.Open(src)
	check(err)
	defer fin.Close()
	switch getImgType(fin) {
	case hBMP:
		im, err = bmp.Decode(fin)
	case hGIF:
		im, err = gif.Decode(fin)
	case hJPG:
		im, err = jpeg.Decode(fin)
	case hPNG:
		im, err = png.Decode(fin)
	default:
		panic(errors.New("unexpected input type"))
	}
	check(err)
	return im
}

func image2File(tar string, im image.Image) {
	var (
		err  error
		fout *os.File
	)
	s2 := getImgTypeSuffix(tar)
	if s2 == 0x00 {
		panic(errors.New("unexpected output type"))
	}
	fout, err = os.Create(tar) // Created after judgment to prevent invalid file generation
	check(err)
	defer fout.Close()
	switch s2 {
	case hBMP:
		err = bmp.Encode(fout, im)
	case hGIF:
		err = gif.Encode(fout, im, nil)
	case hJPG:
		err = jpeg.Encode(fout, im, nil)
	case hPNG:
		err = png.Encode(fout, im)
	}
	check(err)
}

/**
 * c-shared function
 */

// //export rgb2GrayC
// func rgb2GrayC(src *C.char, tar *C.char) {
// 	image2File(C.GoString(tar), rgbaToGray(file2Image(C.GoString(src))))
// }

// //export rgb2GrayC2
// func rgb2GrayC2(src *C.char, tar *C.char) {
// 	var wg sync.WaitGroup
// 	srcG := C.GoString(src)
// 	tarG := C.GoString(tar)
// 	finfo, _ := ioutil.ReadDir(srcG)
// 	for _, x := range finfo {
// 		srcPath := srcG + "/" + x.Name()
// 		tarPath := tarG + "/" + x.Name()
// 		if x.IsDir() {
// 			continue
// 		} else {
// 			wg.Add(1)
// 			go func() {
// 				image2File(tarPath, rgbaToGray(file2Image(srcPath)))
// 				wg.Done()
// 			}()
// 		}
// 	}
// 	wg.Wait()
// }

// //export gray2RgbaC
// func gray2RgbaC(src *C.char, tar *C.char) {
// 	image2File(C.GoString(tar), grayToRgba(file2Image(C.GoString(src))))
// }

// //export verifyGrayC
// func verifyGrayC(src *C.char) bool {
// 	return verifyGray(file2Image(C.GoString(src)))
// }

// //export moveBoundsC
// func moveBoundsC(src *C.char, tar *C.char, left, top, right, bottom int, r, g, b, a uint8) {
// 	image2File(C.GoString(tar), moveBounds(file2Image(C.GoString(src)), left, top, right, bottom, r, g, b, a))
// }

//export exportInitialize
func exportInitialize(path *C.char) int {
	return setImageObject(file2Image(C.GoString(path)))
}

//export exportSave
func exportSave(key int, path *C.char) {
	image2File(C.GoString(path), getImageObject(key))
}
