package firetils

import (
	"context"
	"net/http"

	firebase "firebase.google.com/go/v4"
	fauth "firebase.google.com/go/v4/auth"
)

type fireHolder struct {
	App  *firebase.App
	auth *fauth.Client
}

// NewFireHolder don't use this yet
// This is intended to be a way to have multiple firetils instance with different configs, rather than globals
func NewFireHolder(ctx context.Context, app *firebase.App) *fireHolder {
	return &fireHolder{App: app}
}

func (fh *fireHolder) Auth(ctx context.Context) *fauth.Client {
	if fh.auth == nil {
		var err error
		fh.auth, err = fh.App.Auth(ctx)
		if err != nil {
			panic(err)
		}
	}
	return fh.auth
}

// Middleware that will check the Authorization header and return a forbidden error response if not valid
func (fh *fireHolder) FireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token, err := Authenticate(ctx, fh.Auth(ctx), w, r, false)
		if err != nil {
			return
		}
		// fmt.Printf("authed %v\n", token.UID)
		// ctx = lWith(ctx, "user_id", token.UID)
		ctx = context.WithValue(ctx, TokenContextKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
