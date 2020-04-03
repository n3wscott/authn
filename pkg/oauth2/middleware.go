package oauth2

import (
	"context"
	"github.com/coreos/go-oidc"
	"net/http"
	"strings"
)

func Middleware(ctx context.Context, wrapped http.Handler, verifier *oidc.IDTokenVerifier) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authZHeader := r.Header.Get("authorization")
		if !strings.HasPrefix(authZHeader, "Bearer ") {
			http.Error(w, "expected HTTP bearer authorization header", http.StatusUnauthorized)
			return
		}

		token, err := verifier.Verify(r.Context(), strings.TrimPrefix(authZHeader, "Bearer "))
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		var user User
		if err := token.Claims(&user); err != nil {
			http.Error(w, "error extracting ID token claims", http.StatusInternalServerError)
			return
		}
		wrapped.ServeHTTP(w, r.WithContext(WithAuthenticatedUser(r.Context(), &user)))
	})
}
