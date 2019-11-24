package imageprocessing

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"math"
	"mime"
	"net/http"

	"cloud.google.com/go/storage"
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

func downloadImage(imgUrl string) (image.Image, error) {
	resp, err := http.Head(imgUrl)
	if err != nil {
		return nil, err
	}
	if err = validateResponse(resp, imgUrl); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resp, err = http.Get(imgUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err = validateResponse(resp, imgUrl); err != nil {
		return nil, err
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data string `json:"data"`
}

type processingMessage struct {
	Url     string
	EventID string
}

// ProcessImage processes an image found at provided URL.
//
// A thumbnail of the image gets produced.
func ProcessImage(ctx context.Context, m PubSubMessage) error {
	payload, err := b64.StdEncoding.DecodeString(m.Data)
	if err != nil {
		return err
	}
	var message processingMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		return err
	}
	fmt.Printf("received request to download %s\n", message.Url)
	img, err := downloadImage(message.Url)
	if err != nil {
		return err
	}
	thumb, err := resizeImage(img)
	if err != nil {
		return err
	}

	if err := uploadImages(message.EventID, img, thumb); err != nil {
		return err
	}

	return nil
}

func resizeImage(img image.Image) (*image.NRGBA, error) {
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

func uploadImage(ctx context.Context, name string, eventID string, bucket *storage.BucketHandle, img image.Image) error {
	objPath := fmt.Sprintf("images/events/%s/%s.jpg", eventID, name)
	fmt.Printf("Uploading %q\n", objPath)
	obj := bucket.Object(objPath)
	w := obj.NewWriter(ctx)
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.CacheControl = "public, max-age=86400"
	w.ContentType = mime.TypeByExtension(".jpg")
	if err := jpeg.Encode(w, img, &jpeg.Options{Quality: 100}); err != nil {
		w.Close()
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

func uploadImages(eventID string, img, thumb image.Image) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	bucket := client.Bucket("arve.experimental.berlin")
	if err := uploadImage(ctx, "main", eventID, bucket, img); err != nil {
		return err
	}
	if err := uploadImage(ctx, "thumb", eventID, bucket, thumb); err != nil {
		return err
	}
	return nil
}
