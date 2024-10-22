package middleware

import (
	"app/internal/infrastructure/responder"
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
)

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			responder.NewResponder().ErrorBadRequest(w, err)
			//http.Error(w, err.Error(), 401)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			//http.Error(w, http.StatusText(401), 401)
			responder.NewResponder().ErrorBadRequest(w, errors.New(http.StatusText(401)))
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}
