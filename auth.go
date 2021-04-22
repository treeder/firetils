package firetils

import (
	"context"
	"errors"
	"net/http"
	"strings"

	fauth "firebase.google.com/go/v4/auth"
	"github.com/treeder/gotils"
)

type contextKey string

const (
	tokenContextKey  = contextKey("token")
	userIDContextKey = contextKey("user_id")
)

var (
	authClient *fauth.Client
)

// Authenticate checks the Authorization header for a firebase token
func Authenticate(ctx context.Context, firebaseAuth *fauth.Client, w http.ResponseWriter, r *http.Request, hardVerify bool) (*fauth.Token, error) {
	idToken := r.Header.Get("Authorization")
	if idToken == "" {
		return nil, errors.New("invalid Authorization token")
	}
	splitToken := strings.Split(idToken, " ")
	if len(splitToken) < 2 {
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
				return nil, errors.New("token has been revoked")
			}
			return nil, errors.New("cannot verify token")
		}
	} else {
		token, err = firebaseAuth.VerifyIDToken(ctx, idToken)
		if err != nil {
			return nil, errors.New("cannot verify token")
		}
	}
	return token, nil
}

// FireAuth middleware to guard endpoints
func FireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := Authenticate(r.Context(), authClient, w, r, false)
		if err != nil {
			gotils.WriteError(w, http.StatusForbidden, err)
			return
		}
		// fmt.Printf("authed %v\n", token.UID)
		ctx := r.Context()
		ctx = context.WithValue(ctx, tokenContextKey, token)
		ctx = context.WithValue(ctx, userIDContextKey, token.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth this won't guard it, but will set the token in the context if it's there. Will not error out if it's not there.
func OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := Authenticate(r.Context(), authClient, w, r, false)
		if err != nil {
			// just ignore it
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, tokenContextKey, token)
		ctx = context.WithValue(ctx, userIDContextKey, token.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Token(ctx context.Context) *fauth.Token {
	u, _ := ctx.Value(tokenContextKey).(*fauth.Token)
	return u
}
func UserID(ctx context.Context) string {
	u, _ := ctx.Value(userIDContextKey).(string)
	return u
}

func SetOwned(ctx context.Context, o OwnedI) {
	o.SetUserID(UserID(ctx))
}
