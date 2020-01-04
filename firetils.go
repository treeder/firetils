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

// var (
// 	app      *firebase.App // todo: move this somewhere more general
// 	fireauth *auth.Client
// 	client   *firestore.Client
// )

// New creates a new firestore client
// Call defer client.Close() after this if you can
func New(ctx context.Context, projectID string, opts []option.ClientOption) (*firestore.Client, error) {
	// Use the application default credentials
	var err error
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf, opts...)
	if err != nil {
		gotils.L(ctx).Sugar().Fatalf("couldn't init firebase newapp: %v\n", err)
		return nil, err
	}
	// fireauth, err := app.Auth(ctx)
	// if err != nil {
	// 	gotils.L(ctx).Sugar().Fatalf("error getting Auth client: %v\n", err)
	// 	return nil, err
	// }

	client, err := app.Firestore(ctx)
	if err != nil {
		gotils.L(ctx).Sugar().Fatalf("couldn't init firestore: %v\n", err)
		return nil, err
	}
	return client, nil
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

func SaveGeneric(ctx context.Context, collection, id string, ow *Timestamped) (*firestore.DocumentRef, *Timestamped, error) {
	UpdateTimeStamps(ow)
	ref := client.Collection(collection).Doc(id)
	_, err := ref.Set(ctx, ow)
	if err != nil {
		gotils.L(ctx).Error("Failed to save generic object!", zap.Error(err))
		return nil, nil, errors.New("Failed to store object, please try again")
	}
	return ref, ow, nil
}
