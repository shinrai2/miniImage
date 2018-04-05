package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/bmp"
)

const (
	hBMP  byte   = 0x42
	hGIF  byte   = 0x47
	hJPG  byte   = 0xff
	hPNG  byte   = 0x89
	dNOT  string = "n"
	dGRAY string = "g"
)

func main() {
	option := flag.String("o", "n", "Wanted operation.")
	srcPath := flag.String("src", "img/t.png", "The source file.")
	tarPath := flag.String("tar", "output/t.png", "The target file.")
	flag.Parse()
	switch *option {
	default:
		check(errors.New("unexpected input"))
	case dNOT:
		fmt.Println("Do nothing. :)")
	case dGRAY:
		ToGray(*srcPath, *tarPath)
	}
}

// ToGray method converts a picture to a grayscale
func ToGray(src string, tar string) {
	var (
		err, err2, err3 error
		fin, fout       *os.File
		im              image.Image
	)
	fin, err = os.Open(src)
	defer fin.Close()
	check(err)
	switch getImgType(fin) {
	case hBMP:
		im, err2 = bmp.Decode(fin)
	case hGIF:
		im, err2 = gif.Decode(fin)
	case hJPG:
		im, err2 = jpeg.Decode(fin)
	case hPNG:
		im, err2 = png.Decode(fin)
	}
	check(err2)
	s2 := getImgTypeSuffix(tar)
	fout, err = os.Create(tar) // Created after judgment to prevent invalid file generation
	defer fin.Close()
	check(err)
	switch s2 {
	case hBMP:
		err3 = bmp.Encode(fout, rgbaToGray(im))
	case hGIF:
		err3 = gif.Encode(fout, rgbaToGray(im), nil)
	case hJPG:
		err3 = jpeg.Encode(fout, rgbaToGray(im), nil)
	case hPNG:
		err3 = png.Encode(fout, rgbaToGray(im))
	}
	check(err3)
}

func getImgType(f *os.File) byte {
	tmp := make([]byte, 1)
	_, err := f.Read(tmp)
	check(err)
	f.Seek(0, os.SEEK_SET) // Reset pointer
	return tmp[0] ^ hBMP ^ hGIF ^ hJPG ^ hPNG ^ hBMP ^ hGIF ^ hJPG ^ hPNG
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
		check(errors.New("unexpected output type"))
		return 0x00
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

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
