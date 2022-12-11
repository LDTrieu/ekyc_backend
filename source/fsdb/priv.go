package fsdb

import (
	"context"
	"ekyc-app/gcloud"
	"ekyc-app/package/wlog"
	"ekyc-app/source/model"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
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

func getOneEqual(ctx context.Context, obj interface{},
	coll, field string, val interface{}) (
	id string, err error) {
	count := 0
	t0 := time.Now()
	err = gcloud.RunFS(ctx, coll,
		func(collectionRef *firestore.CollectionRef) error {
			iter := collectionRef.Where(field, "==", val).Limit(1).Documents(ctx)
			for {
				doc, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					return err
				}

				id = doc.Ref.ID
				err = doc.DataTo(obj)
				if err != nil {
					return err
				}
				count++
			}
			return nil
		})
	if err != nil {
		wlog.Info(ctx, coll, " getOneEqual ", err)
		return
	}
	//---
	dur := time.Since(t0).Milliseconds()
	wlog.Info(ctx, coll, " getOneEqual ", dur, "| key", field, "==", val, "| result:", count > 0)
	if count == 0 {
		return "", model.ErrDocNotFound
	}
	return id, nil
}

func delete(ctx context.Context, coll, id string) (err error) {
	t0 := time.Now()
	err = gcloud.RunFS(ctx, coll,
		func(collectionRef *firestore.CollectionRef) error {
			_, err := collectionRef.Doc(id).Delete(ctx)
			if err != nil {
				if status.Code(err) == codes.NotFound {
					return model.ErrDocNotFound
				}
				return err
			}
			return nil
		})
	dur := time.Since(t0).Milliseconds()
	wlog.Info(ctx, coll, " delete ", dur)
	return err
}
func updateFields(ctx context.Context, id, coll string, obj map[string]interface{}) error {
	if len(id) == 0 {
		return model.ErrDocIdEmpty
	}
	t0 := time.Now()
	err := gcloud.RunFS(ctx, coll,
		func(collectionRef *firestore.CollectionRef) error {
			_, err := collectionRef.Doc(id).Set(ctx, obj, firestore.MergeAll)
			if err != nil {
				return err
			}
			return nil
		})
	dur := time.Since(t0).Milliseconds()
	wlog.Info(ctx, coll, " set ", dur)
	return err
}
