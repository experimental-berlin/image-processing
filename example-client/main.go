package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"os"
)

func main() {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx,
		"experimental-berlin", option.WithCredentialsFile("keyfile.json"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create pubsub client: %s\n", err)
		os.Exit(1)
	}

	defer client.Close()

	topic, err := client.CreateTopic(ctx, "imageProcessingResponses")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create topic: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created topic: %s\n", topic)
}
