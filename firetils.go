package firetils

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	fauth "firebase.google.com/go/v4/auth"
	"github.com/treeder/gotils"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

var (
// App      *firebase.App
// fireauth *fauth.Client
)

type contextKey string

const (
	TokenContextKey = contextKey("token")
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

// Authenticate can be used to check the Authorization header
func Authenticate(ctx context.Context, firebaseAuth *fauth.Client, w http.ResponseWriter, r *http.Request, hardVerify bool) (*fauth.Token, error) {
	idToken := r.Header.Get("Authorization")
	if idToken == "" {
		gotils.WriteError(w, http.StatusForbidden, errors.New("Invalid token"))
		return nil, errors.New("invalid Authorization token")
	}
	splitToken := strings.Split(idToken, " ")
	if len(splitToken) < 2 {
		gotils.WriteError(w, http.StatusForbidden, errors.New("Invalid token"))
		return nil, errors.New("invalid Authorization token")
	}
	idToken = splitToken[1]

	var err error
	var token *fauth.Token
	if hardVerify {
		token, err = firebaseAuth.VerifyIDTokenAndCheckRevoked(ctx, idToken)
		if err != nil {
			if err.Error() == "ID token has been revoked" {
				// Token is revoked. Inform the user to reauthenticate or signOut() the user.
				gotils.L(ctx).Warn("ID token was revoked", zap.Error(err))
				gotils.WriteError(w, http.StatusForbidden, errors.New("token has been revoked"))
				return nil, errors.New("token has been revoked")
			}
			gotils.L(ctx).Warn("error verifying ID token with firebase", zap.Error(err))
			gotils.WriteError(w, http.StatusForbidden, errors.New("cannot verify token"))
			return nil, errors.New("cannot verify token")
		}
	} else {
		token, err = firebaseAuth.VerifyIDToken(ctx, idToken)
		if err != nil {
			gotils.L(ctx).Warn("error verifying ID token", zap.Error(err))
			gotils.WriteError(w, http.StatusForbidden, errors.New("cannot verify token"))
			return nil, errors.New("cannot verify token")
		}
	}
	return token, nil
}
