package imageprocessing

//
// import (
// 	"fmt"
// 	"image"
// 	_ "image/png"
// 	"log"
// 	"net/http"
// 	"net/http/httptest"
// 	"path"
// 	"testing"
//
// 	"github.com/disintegration/imaging"
// )
//
// func absint(i int) int {
// 	if i < 0 {
// 		return -i
// 	}
// 	return i
// }
//
// func compareImages(img1, img2 image.Image) bool {
// 	bounds1 := img1.Bounds()
// 	bounds2 := img2.Bounds()
// 	if !bounds1.Eq(bounds2) {
// 		return false
// 	}
//
// 	for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
// 		for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
// 			r1, g1, b1, a1 := img1.At(x, y).RGBA()
// 			r2, g2, b2, a2 := img2.At(x, y).RGBA()
// 			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
// 				return false
// 			}
// 		}
// 	}
//
// 	return true
// }
//
// func openImage(fpath string) image.Image {
// 	img, err := imaging.Open(fpath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	return img
// }
//
// func TestProcess(t *testing.T) {
// 	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fpath := path.Join("testdata", r.URL.Path)
// 		http.ServeFile(w, r, fpath)
// 	}))
// 	defer fakeServer.Close()
//
// 	testCases := []struct {
// 		name      string
// 		srcFname  string
// 		wantFname string
// 		wantErr   string
// 	}{
// 		{
// 			"Tall image",
// 			"tall-image.png",
// 			"tall-image-thumbnail.png",
// 			"",
// 		},
// 		{
// 			"Wide image",
// 			"wide-image.png",
// 			"wide-image-thumbnail.png",
// 			"",
// 		},
// 		{
// 			"Too large image",
// 			"too-large-image.png",
// 			"",
// 			"Image is greater than 3 MB in size",
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			url := fmt.Sprintf("%v/%v", fakeServer.URL, tc.srcFname)
// 			got, err := ProcessImage(url)
// 			if err != nil {
// 				if tc.wantErr == "" {
// 					t.Fatalf("Unexpected error: %s", err)
// 					return
// 				}
//
// 				expectedErr := fmt.Sprintf("%v: %v", tc.wantErr, url)
// 				if err.Error() != expectedErr {
// 					t.Fatalf("Expected error '%v', got '%v'", expectedErr, err)
// 				}
//
// 				return
// 			}
//
// 			want := openImage(path.Join("testdata", tc.wantFname))
// 			// Get thumbnail image from upload mock
// 			if !compareImages(got, want) {
// 				t.Error("Didn't get expected thumbnail")
// 			}
// 		})
// 	}
// }
