package cfg

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

type Info struct {
	PrePrivKey string `firestore:"pre_priv_key"`
	//FacepayPrePublic string `firestore:"facepay_account_service_public_key"`
	LogBucket          string `firestore:"log_bucket"`
	TerminalFileBucket string `firestore:"terminal_file_bucket"`
}

var (
	fsClient *firestore.Client
)

const (
	devProjectId = "ekyc-services"
)

func getDbClient(ctx context.Context) (client *firestore.Client, err error) {
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

func Get(ctx context.Context) (info *Info, err error) {
	err = RunFirestore(ctx, func(client *firestore.Client) error {
		docRef := client.Collection("cfg").Doc("server")
		snap, err := docRef.Get(ctx)
		if err != nil {
			return err
		}
		info = &Info{}
		err = snap.DataTo(info)
		return err
	})
	return
}
