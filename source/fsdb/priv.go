package fsdb

import (
	"context"
	"ekyc-app/gcloud"
	"ekyc-app/package/wlog"
	"ekyc-app/source/model"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Add generates a DocumentRef with a unique ID.
func add(ctx context.Context, coll string, obj interface{}) (id string, err error) {
	t0 := time.Now()
	err = gcloud.RunFS(ctx, coll,
		func(collectionRef *firestore.CollectionRef) error {
			ref, _, err := collectionRef.Add(ctx, obj)
			if err != nil {
				return err
			}
			id = ref.ID
			return nil
		})
	dur := time.Since(t0).Milliseconds()
	wlog.Info(ctx, coll, " add ", dur)
	return
}

func getById(ctx context.Context, obj interface{}, coll, id string) (
	err error) {
	t0 := time.Now()
	err = gcloud.RunFS(ctx, coll, func(collectionRef *firestore.CollectionRef) error {
		docRef, err := collectionRef.Doc(id).Get(ctx)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return model.ErrDocNotFound
			}
			return err
		}
		if err = docRef.DataTo(obj); err != nil {
			return err
		}
		return nil
	})

	dur := time.Since(t0).Milliseconds()
	wlog.Info(ctx, coll, " getByID ", dur, " milliseconds")
	return
}
