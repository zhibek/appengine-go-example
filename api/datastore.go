package main

import (
	"context"
	"errors"
	"log"
	"os"

	"cloud.google.com/go/datastore"

	"github.com/rs/xid"
)

// Demo entity struct
type Demo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// initDatastore Init Datastore client using current PROJECT_ID
func initDatastore(ctx context.Context) (*datastore.Client, error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		// In local environment, any projectID can be used to access local datastore
		projectID = "dev"
	}
	log.Printf("initDatastore() %s", projectID)

	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("initDatastore() Failed to create client - %s", err)
		return nil, errors.New("Failed to create client")
	}

	return client, nil
}

// CreateDemo Create a new "Demo" entity
func CreateDemo(ctx context.Context, demo *Demo) (*Demo, error) {
	log.Printf("createDemo()")
	client, err := initDatastore(ctx)
	if err != nil {
		log.Printf("createDemo() ERR %s", err)
		return nil, err
	}

	demo.ID = xid.New().String()
	demoKey := datastore.NameKey("Demo", demo.ID, nil)

	ret, err := client.Put(ctx, demoKey, demo)
	if err != nil {
		log.Printf("createDemo() ERR %s", err)
		return nil, err
	}
	log.Printf("createDemo() client.Put() %s", ret)

	return demo, nil
}

// FetchDemo Fetch a new "Demo" entity, using ID
func FetchDemo(ctx context.Context, id string) (*Demo, error) {
	log.Printf("fetchDemo()")
	client, err := initDatastore(ctx)
	if err != nil {
		log.Printf("fetchDemo() ERR %s", err)
		return nil, err
	}

	demoKey := datastore.NameKey("Demo", id, nil)
	demo := new(Demo)
	if err := client.Get(ctx, demoKey, demo); err != nil {
		log.Printf("fetchDemo() ERR %s", err)
		return nil, err
	}
	log.Printf("fetchDemo() client.Put() %s", demo)

	return demo, nil
}
