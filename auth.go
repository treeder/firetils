package firetils

import (
	"context"
	"errors"
	"net/http"
	"strings"

	fauth "firebase.google.com/go/v4/auth"
	"github.com/treeder/gotils"
	"go.uber.org/zap"
)

// Authenticate checks the Authorization header for a firebase token
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
