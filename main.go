package main

import "C"

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
	dNOT  string = "not"
	dHELP string = "help"
	dGRAY string = "g"
)

func main() {
	option := flag.String("o", dNOT, "Wanted operation.")
	srcPath := flag.String("src", "img/t.png", "The source file.")
	tarPath := flag.String("tar", "output/t.png", "The target file.")
	moreEx := flag.String("more", "", "More expansion operations.")
	flag.Parse()
	fmt.Println("-------------------------------")
	fmt.Println("Option:\t\t", *option)
	fmt.Println("SrcPath:\t", *srcPath)
	fmt.Println("TarPath:\t", *tarPath)
	fmt.Println("moreEx:\t\t", *moreEx)
	fmt.Println("-------------------------------")
	switch *option {
	default:
		// check(errors.New("unexpected input"))
		o, me := parseOaEx(*option, *moreEx)
		image2File(*tarPath, op(file2Image(*srcPath), o, me))
		fmt.Println("All conversions were successful. :)")
	case dNOT:
		fmt.Println("Nothing to do. :)")
	case dHELP:
		oHelp()
	}

}

func op(im image.Image, o []string, mex []string) image.Image {
	if len(o) == 0 {
		return im
	}
	switch o[0] {
	case dGRAY:
		fmt.Println("Successful conversion: rgbaToGray. :)")
		return op(rgbaToGray(im), o[1:], mex[1:])
	default:
		return im
	}
}

func parseOaEx(o string, mex string) ([]string, []string) {
	ro := make([]string, len(o))
	for i := 0; i < len(o); i++ {
		ro[i] = o[i : i+1]
	}
	rmex := strings.Split(mex, "|")
	if len(ro) != len(rmex) {
		panic(errors.New("The input parameters do not match"))
	}
	return ro, rmex
}

func oHelp() {
	fmt.Println("Parameter format:")
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

//export rgb2GrayC
func rgb2GrayC(src *C.char, tar *C.char) {
	image2File(C.GoString(tar), rgbaToGray(file2Image(C.GoString(src))))
}
