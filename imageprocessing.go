package imageprocessing

import (
	"errors"
	"image"
	"math"
	"net/http"

	"github.com/disintegration/imaging"
)

// ProcessImage processes an image found at provided URL into also a thumbnail
func ProcessImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("Failed to download file")
	}

	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	targetRatio := 1.5
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	sourceRatio := float64(width) / float64(height)
	var cropRectangle image.Rectangle
	if sourceRatio < targetRatio {
		targetHeight := int(float64(width) / targetRatio)
		offset := int(math.Floor(float64(height-targetHeight) / float64(2)))
		cropRectangle = image.Rectangle{
			Min: image.Point{X: 0, Y: offset},
			Max: image.Point{X: width, Y: height - offset},
		}
	} else if sourceRatio > targetRatio {
		targetWidth := int(targetRatio * float64(height))
		offset := int(math.Floor(float64(width-targetWidth) / float64(2)))
		cropRectangle = image.Rectangle{
			Min: image.Point{X: offset, Y: 0},
			Max: image.Point{X: width - offset, Y: height},
		}
	}
	croppedImage := imaging.Crop(img, cropRectangle)

	return imaging.Resize(croppedImage, 800, 534, imaging.Lanczos), nil
}
