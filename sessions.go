package firetils

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
	"github.com/treeder/gotils/v2"
)

// SessionHandler little helper you can throw in a route to create a session cookie
// This will NOT create a user for you.
// func SessionHandler() http.HandlerFunc {
// 	return gotils.ErrorHandler(CreateSession)
// }

// CreateSession creates a Firebase session cookie
// more: https://firebase.google.com/docs/auth/admin/manage-cookies
// Use in a handler like this:
// func createSession(w http.ResponseWriter, r *http.Request) error {
// 	ctx := r.Context()

// 	idToken, resp, err := CreateSession(globals.App.Firebase.Auth, r)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = EnsureUserDefault(ctx, globals.App.Firebase.Auth, globals.App.Firebase.Firestore, idToken, models.UsersCollection)
// 	if err != nil {
// 		return err
// 	}

//		return gotils.WriteObject(w, http.StatusOK, resp)
//	}
func CreateSession(auth *auth.Client, r *http.Request) (token string, response map[string]any, err error) {
	ctx := r.Context()
	expiresIn := time.Hour * 24 * 14 // max is 2 weeks
	// Create the session cookie. This will also verify the ID token in the process.
	// The session cookie will have the same claims as the ID token.
	// To only allow session cookie setting on recent sign-in, auth_time in ID token
	// can be checked to ensure user was recently signed in before creating a session cookie.
	idTokenIn := r.Header.Get("Authorization")
	splitToken := strings.Split(idTokenIn, " ")
	token = splitToken[1]
	cookie, err := auth.SessionCookie(ctx, token, expiresIn)
	if err != nil {
		return token, nil, gotils.C(ctx).Errorf("Failed to create session cookie: %w", err)
	}
	// return gotils.WriteObject(w, http.StatusOK, map[string]any{"cookie": cookie, "expires": int(expiresIn.Seconds())})
	return token, map[string]any{"cookie": cookie, "expires": int(expiresIn.Seconds())}, nil
}

// This will verify the idToken, fetch the user from firebase auth, then update the 'usersCollection' with, email, name and image IF it doesn't already exist.
// TODO: Will update those fields if they have changed.
func EnsureUserDefault(ctx context.Context, auth *auth.Client, fs *firestore.Client, idToken, usersCollection string) (*auth.UserRecord, error) {
	ctx = gotils.With(ctx, "func", "ensureUser")
	token, err := auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, gotils.C(ctx).Errorf("cannot verify ID token: %w", err)
	}

	userID := token.UID
	firebaseUser, err := auth.GetUser(ctx, userID)
	if err != nil {
		return nil, gotils.C(ctx).Errorf("error getting user from firebase: %w", err)
	}

	user := &FirestoredAndTimeStamped{}
	err = GetByID(ctx, fs, usersCollection, userID, user)
	if err != nil {
		if !errors.Is(err, gotils.ErrNotFound) {
			return nil, gotils.C(ctx).Errorf("Error getting user: %w", err)
		}
	}

	// then we make new user
	_, err = fs.Collection(usersCollection).Doc(userID).Create(ctx, map[string]any{
		"email": firebaseUser.Email,
		"name":  firebaseUser.DisplayName,
		"image": firebaseUser.PhotoURL,
	})
	if err != nil {
		return nil, gotils.C(ctx).Errorf("error updating user: %w", err)
	}
	return firebaseUser, nil
}
