# Experimental Berlin Image Processing Service
Experimental Berlin's service for processing images.

This is a serverless function that receives an image URL from the client, fetches the image, generates a thumbnail and uploads it to Google Cloud Storage.

## Testing
```$ go test```

## Deployment

```
$ gcloud functions deploy imageProcessing --region europe-west1 --entry-point ProcessImage --memory 1024 --runtime go111 --trigger-topic imageProcessingRequest
```

## Reading Logs
```
$ gcloud functions logs read imageProcessing --region europe-west1
```

## Publishing an Image Processing Request
```
$ gcloud pubsub topics publish imageProcessingRequest --project experimental-berlin --message "http://example.com"
```
