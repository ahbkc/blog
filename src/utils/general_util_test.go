package utils

import (
	"testing"
	"fmt"
	"image"
	"image/color"
	"os"
	"image/png"
	"regexp"
)

func TestCutString(t *testing.T) {
	param := "/qwe/afasf"
	sep := "/"
	result, err := CutString(param, sep, 2, 0)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result)
}

func TestRandNum(t *testing.T) {
	for i := 0; i < 20; i ++ {
		t.Log(GetRandCode())
	}
}

func TestCreatePng(t *testing.T) {
	const width, height = 256, 256

	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8((x + y) & 255),
				G: uint8((x + y) << 1 & 255),
				B: uint8((x + y) << 2 & 255),
				A: 255,
			})
		}
	}

	f, err := os.Create("image.png")
	if err != nil {
		t.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		t.Fatal(err)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestUint8(t *testing.T) {
	t.Log(uint8((1 + 2) << 2 & 255))
}

//admin path regexp
func TestPathRegexp(t *testing.T) {
	r, err := regexp.Compile(`^Admin*`)
	if err != nil {
		t.Fatal(err)
	}
	var array = []string{"AdminIndexGet", "adminIndexGet"}
	for _, v := range array {
		str := r.FindString(v)
		fmt.Println(str)
	}
}