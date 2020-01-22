package firetils

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/treeder/gotils"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

// New creates a new firestore client
// Call defer client.Close() after this if you can
func New(ctx context.Context, projectID string, opts []option.ClientOption) (*firebase.App, error) {
	// Use the application default credentials
	var err error
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf, opts...)
	if err != nil {
		// gotils.L(ctx).Sugar().Fatalf("couldn't init firebase newapp: %v\n", err)
		return nil, err
	}
	return app, nil
}

// func FirebaseApp() *firebase.App {
// 	return app
// }

// func FirebaseAuth() *auth.Client {
// 	return fireauth
// }

// func DefaultClient() *firestore.Client {
// 	return client
// }

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
