package database

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"go.uber.org/fx"
)

// ProvideDB provides a firestore client
func ProvideDB() *firestore.Client {
	projectID := "caffy-beans-example"

	client, err := firestore.NewClient(context.TODO(), projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

var Module = fx.Provide(ProvideDB)
