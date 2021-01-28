package firetils

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"github.com/treeder/gotils"
	"go.uber.org/zap"
)

func SaveGeneric(ctx context.Context, client *firestore.Client, collection, id string, ow *Timestamped) (*firestore.DocumentRef, *Timestamped, error) {
	UpdateTimeStamps(ow)
	ref := client.Collection(collection).Doc(id)
	_, err := ref.Set(ctx, ow)
	if err != nil {
		gotils.L(ctx).Error("Failed to save generic object!", zap.Error(err))
		return nil, nil, errors.New("Failed to store object, please try again")
	}
	return ref, ow, nil
}
