package firetils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	fauth "firebase.google.com/go/v4/auth"
	"github.com/treeder/gotils/v2"
)

type contextKey string

const (
	TokenContextKey  = contextKey("token")
	UserIDContextKey = contextKey("user_id")
)

var (
	InvalidToken = errors.New("Invalid auth token")
	NoToken      = errors.New("No auth token provided")
	authClient   *fauth.Client
)

// Authenticate checks the Authorization header for a firebase token
func Authenticate(ctx context.Context, firebaseAuth *fauth.Client, w http.ResponseWriter, r *http.Request, hardVerify bool) (*fauth.Token, error) {
	var err error
	idToken := r.Header.Get("Authorization")
	cookie, _ := r.Cookie("session")
	if cookie == nil {
		cookie, _ = r.Cookie("__session") // only cookie allowed with firebase hosting: https://stackoverflow.com/a/44935288/105562
	}
	// gotils.L(ctx).Info().Println("Authorization!!!", idToken, "cookie:", cookie, "--")
	if cookie == nil && idToken == "" {
		return nil, NoToken
	}
	sessionCookie := ""
	var token *fauth.Token

	if idToken != "" {
		splitToken := strings.Split(idToken, " ")
		if len(splitToken) < 2 {
			return nil, InvalidToken
		}
		authType := splitToken[0]
		idToken = splitToken[1]
		if authType == "Bearer" {
			if hardVerify {
				token, err = firebaseAuth.VerifyIDTokenAndCheckRevoked(ctx, idToken)
				if err != nil {
					// gotils.L(ctx).Error().Println(err)
					if err.Error() == "ID token has been revoked" {
						// Token is revoked. Inform the user to reauthenticate or signOut() the user.
						return nil, errors.New("token has been revoked")
					}
					return nil, fmt.Errorf("cannot verify token: %w", err)
				}
			} else {
				token, err = firebaseAuth.VerifyIDToken(ctx, idToken)
				if err != nil {
					return nil, fmt.Errorf("cannot verify token: %w", err)
				}
			}
			return token, nil
		} else if authType == "Cookie" {
			sessionCookie = idToken
		}
	}
	if cookie != nil {
		sessionCookie = cookie.Value
	}
	if sessionCookie != "" {
		if hardVerify {
			token, err = firebaseAuth.VerifySessionCookieAndCheckRevoked(ctx, sessionCookie)
		} else {
			token, err = firebaseAuth.VerifySessionCookie(ctx, sessionCookie)
		}
		if err != nil {
			// gotils.L(ctx).Error().Println(err)
			return nil, fmt.Errorf("cannot verify token: %w", err)
		}
		return token, nil
	}
	return nil, InvalidToken
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
		ctx = context.WithValue(ctx, TokenContextKey, token)
		ctx = context.WithValue(ctx, UserIDContextKey, token.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth this won't guard it, but will set the token in the context if it's there. Will not error out if it's not there.
func OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token, err := Authenticate(ctx, authClient, w, r, false)
		if err != nil {
			if !errors.Is(err, NoToken) {
				gotils.L(ctx).Error().Printf("auth error, but optional so skipping: %v", err)
			}
			next.ServeHTTP(w, r)
			return
		}
		ctx = context.WithValue(ctx, TokenContextKey, token)
		ctx = context.WithValue(ctx, UserIDContextKey, token.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Token(ctx context.Context) *fauth.Token {
	u, _ := ctx.Value(TokenContextKey).(*fauth.Token)
	return u
}
func UserID(ctx context.Context) string {
	u, _ := ctx.Value(UserIDContextKey).(string)
	return u
}

func SetOwned(ctx context.Context, o OwnedI) {
	o.SetUserID(UserID(ctx))
}
