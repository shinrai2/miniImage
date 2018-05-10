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
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/bmp"
)

const (
	hBMP byte    = 0x42
	hGIF byte    = 0x47
	hJPG byte    = 0xff
	hPNG byte    = 0x89
	iMAG byte    = 0x22
	fONT byte    = 0x33
	sDPI float64 = 72
)

var once sync.Once
var imageMap = make(map[int]image.Image)
var fontMap = make(map[int]*truetype.Font)

func main() {
	var mor string
	fmt.Println("Only supports calling as dynamic link library.")
	fmt.Print("?> ")
	fmt.Scanf("%s", &mor)
	if mor == "debug" {
		fmt.Println("DEBUG MODE::") // put some codes next here.
		img := paintNew(200, 200, 128, 128, 128, 255)
		font := loadFont("fonts/times.ttf")
		drawFont(img, font, 20, 255, 255, 255, 255, 100, 100, "123")
		image2File("output/debug.bmp", img)
		fmt.Println("Finished. :)")
	} else {
		fmt.Println("Unknown command. :(")
	}
}

/**
 * assist function
 */

func getImageObject(key int) image.Image {
	return imageMap[key]
}

func newImageObject(value image.Image) int {
	once.Do(func() { rand.Seed(time.Now().UnixNano()) })
	randi := rand.Int()
	for imageMap[randi] != nil {
		randi = rand.Int()
	}
	imageMap[randi] = value
	return randi
}

func setImageObject(key int, value image.Image) {
	imageMap[key] = value
}

func delImageObject(key int) {
	delete(imageMap, key)
}

func getFontObject(key int) *truetype.Font {
	return fontMap[key]
}

func newFontObject(value *truetype.Font) int {
	once.Do(func() { rand.Seed(time.Now().UnixNano()) })
	randi := rand.Int()
	for fontMap[randi] != nil {
		randi = rand.Int()
	}
	fontMap[randi] = value
	return randi
}

func setFontObject(key int, value *truetype.Font) {
	fontMap[key] = value
}

func delFontObject(key int) {
	delete(fontMap, key)
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

func paintNew(width, height int, r, g, b, a uint8) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	backColor := color.RGBA{r, g, b, a}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, backColor)
		}
	}
	return img
}

func loadFont(path string) *truetype.Font {
	fontBytes, err := ioutil.ReadFile(path)
	check(err)
	font, err := freetype.ParseFont(fontBytes)
	check(err)
	return font
}

func drawFont(img image.Image, font *truetype.Font,
	fontSize float64, r, g, b, a uint8, x, y int, content string) {
	c := freetype.NewContext()
	c.SetDPI(sDPI)
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	switch value := img.(type) {
	case *image.RGBA:
		c.SetDst(value)
	case *image.Gray:
		c.SetDst(value)
	case *image.Paletted:
		c.SetDst(value)
	default:
		fmt.Errorf("Unsupport format")
	}
	c.SetSrc(image.NewUniform(color.RGBA{r, g, b, a}))
	pt := freetype.Pt(x, y+int(c.PointToFixed(fontSize*1.5)>>8)) // move the y dim to line center

	_, err := c.DrawString(content, pt)
	check(err)
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

//export exportFromFile
func exportFromFile(symbol byte, path *C.char) int {
	switch symbol {
	case iMAG:
		return newImageObject(file2Image(C.GoString(path)))
	case fONT:
		return newFontObject(loadFont(C.GoString(path)))
	default:
		fmt.Errorf("Unknown symbol: %x", symbol)
	}
	return -1
}

//export exportSave
func exportSave(key int, path *C.char) {
	image2File(C.GoString(path), getImageObject(key))
}

//export exportRelease
func exportRelease(symbol byte, key int) {
	switch symbol {
	case iMAG:
		delImageObject(key)
	case fONT:
		delFontObject(key)
	default:
		fmt.Errorf("Unknown symbol: %x", symbol)
	}
}

//export exportIsGray
func exportIsGray(key int) bool {
	return verifyGray(getImageObject(key))
}

//export exportToGray
func exportToGray(key int, self bool) int {
	tmp := rgbaToGray(getImageObject(key))
	if self == true {
		setImageObject(key, tmp)
		return key
	} else {
		return newImageObject(tmp)
	}
}

//export exportToRgba
func exportToRgba(key int, self bool) int {
	tmp := grayToRgba(getImageObject(key))
	if self == true {
		setImageObject(key, tmp)
		return key
	} else {
		return newImageObject(tmp)
	}
}

//export exportMoveBounds
func exportMoveBounds(key, left, top, right, bottom int, r, g, b, a uint8, self bool) int {
	tmp := moveBounds(getImageObject(key), left, top, right, bottom, r, g, b, a)
	if self == true {
		setImageObject(key, tmp)
		return key
	} else {
		return newImageObject(tmp)
	}
}

//export exportNewBlank
func exportNewBlank(width, height int, r, g, b, a uint8) int {
	return newImageObject(paintNew(width, height, r, g, b, a))
}

//export exportDrawString
func exportDrawString(keyOfFont, keyOfImg int, fontSize float64,
	x, y int, content *C.char, r, g, b, a uint8) {
	drawFont(getImageObject(keyOfImg), getFontObject(keyOfFont), fontSize, r, g, b, a, x, y, C.GoString(content))
}

// //export exportSplice
// func exportSplice(keyOfSplice, keyOfImg int) {

// }

// //export exportSpliced
// func exportSpliced(keyOfSplice int) {

// }
