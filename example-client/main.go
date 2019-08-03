package main

import (
  "fmt"
  "cloud.google.com/go/pubsub"
  "google.golang.org/api/option"
  "context"
  "os"
)

func main() {
  _, err := pubsub.NewClient(context.Background(),
    "experimental-berlin", option.WithCredentialsFile("keyfile.json"),
  )
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to create pubsub client: %s\n", err)
    os.Exit(1)
  }
}
