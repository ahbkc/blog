package utils

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"regexp"
	"testing"
)

type PathRegex struct {
	Regex string
	Path  []string
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

var paths = []PathRegex{
	{`^Admin*`, []string{"AdminIndexGet", "adminIndexGet"}},
	{`^/admin/[^login]`, []string{"/admin/login", "/admin/index", "/admin/category/index"}},
}

//admin path regexp
func TestPathRegexp(t *testing.T) {
	r := &regexp.Regexp{}
	for _, v := range paths {
		r = regexp.MustCompile(v.Regex)
		for _, k := range v.Path {
			l := r.MatchString(k)
			t.Logf("regex: %s, data: %s, result: %v\n", v.Regex, k, l)
		}
	}
}
