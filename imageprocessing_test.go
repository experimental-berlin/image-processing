package imageprocessing

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"path"
	"testing"

	"github.com/disintegration/imaging"
)

func absint(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func compareImages(img1, img2 image.Image) bool {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()
	if !bounds1.Eq(bounds2) {
		return false
	}

	for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
		for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
			r1, g1, b1, a1 := img1.At(x, y).RGBA()
			r2, g2, b2, a2 := img2.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				fmt.Printf("r: %v != %v, g: %v != %v, b: %v != %v", r1, r2, g1, g2, b1, b2)
				return false
			}
		}
	}

	return true
}

func openImage(fpath string) image.Image {
	img, err := imaging.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func TestProcess(t *testing.T) {
	testCases := []struct {
		name      string
		srcFname  string
		wantFname string
	}{
		{
			"Image",
			"tall-image.png",
			"tall-image-thumbnail.png",
		},
		{
			"Image",
			"wide-image.png",
			"wide-image-thumbnail.png",
		},
	}
	for _, tc := range testCases {
		src := openImage(path.Join("testdata", tc.srcFname))
		want := openImage(path.Join("testdata", tc.wantFname))
		t.Run(tc.name, func(t *testing.T) {
			got := ProcessImage(src)
			if !compareImages(got, want) {
				t.Error("Didn't get expected thumbnail")
			}
		})
	}
}
