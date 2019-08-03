package imageprocessing

import (
	"fmt"
	"image"
	"math"
	"net/http"

	"github.com/disintegration/imaging"
)

func validateResponse(resp *http.Response, url string) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Failed to download %v", url)
	}

	if resp.ContentLength < 0 {
		return fmt.Errorf("Couldn't deduce content length from response")
	}

	sizeMb := float64(resp.ContentLength) / 1024 / 1024
	maxSizeMb := 3.0
	if sizeMb > maxSizeMb {
		return fmt.Errorf("Image is greater than %v MB in size: %v", maxSizeMb, url)
	}

	return nil
}

func downloadImage(url string) (image.Image, error) {
	resp, err := http.Head(url)
	if err != nil {
		return nil, err
	}
	if err = validateResponse(resp, url); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resp, err = http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err = validateResponse(resp, url); err != nil {
		return nil, err
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// ProcessImage processes an image found at provided URL.
//
// A thumbnail of the image gets produced.
func ProcessImage(url string) (image.Image, error) {
	img, err := downloadImage(url)
	if err != nil {
		return nil, err
	}

	targetRatio := 1.5
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	sourceRatio := float64(width) / float64(height)
	var cropRectangle image.Rectangle
	if sourceRatio < targetRatio {
		// Crop the image height wise
		targetHeight := int(float64(width) / targetRatio)
		offset := int(math.Floor(float64(height-targetHeight) / float64(2)))
		cropRectangle = image.Rectangle{
			Min: image.Point{X: 0, Y: offset},
			Max: image.Point{X: width, Y: height - offset},
		}
	} else if sourceRatio > targetRatio {
		// Crop the image width wise
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
