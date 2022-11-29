package gcloud

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

var (
	fsClient *firestore.Client
)

func getDbClient(ctx context.Context) (
	client *firestore.Client, err error) {
	if fsClient == nil {
		conf := &firebase.Config{}
		app, err := firebase.NewApp(ctx, conf)
		if err != nil {
			return nil, err
		}
		fsClient, err = app.Firestore(ctx)
		if err != nil {
			conf.ProjectID = devProjectId
			app, _ = firebase.NewApp(ctx, conf)
			fsClient, err = app.Firestore(ctx)
			if err != nil {
				return nil, err
			}
		}
	}
	return fsClient, nil
}
func RunFirestore(ctx context.Context,
	f func(client *firestore.Client) error) (err error) {
	client, err := getDbClient(ctx)
	if err != nil {
		return
	}
	err = f(client)
	return
}
func RunFS(ctx context.Context, collName string,
	f func(coll *firestore.CollectionRef) error) error {
	client, err := getDbClient(ctx)
	if err != nil {
		return err
	}
	ref := client.Collection(collName)
	return f(ref)
}
