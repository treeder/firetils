package firetils

import (
	"context"

	firebase "firebase.google.com/go/v4"
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
		return nil, err
	}
	AuthClient, err = app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return app, nil
}
