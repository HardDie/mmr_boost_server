package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/HardDie/mmr_boost_server/internal/errs"
	"github.com/HardDie/mmr_boost_server/internal/service"
	"github.com/HardDie/mmr_boost_server/internal/utils"
)

type AuthMiddleware struct {
	service *service.Service
}

func NewAuthMiddleware(srvc *service.Service) *AuthMiddleware {
	return &AuthMiddleware{
		service: srvc,
	}
}
func (m *AuthMiddleware) RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := utils.GetCookie(r)

		// If we got no cookie
		if cookie == nil || cookie.Value == "" {
			http.Error(w, "Session token is empty", http.StatusBadRequest)
			return
		}

		// Validate if session is active
		ctx := r.Context()
		session, err := m.service.AuthValidateCookie(ctx, cookie.Value)
		if err != nil {
			if errors.Is(err, errs.SessionInvalid) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else {
				http.Error(w, "Internal error", http.StatusInternalServerError)
			}
			return
		}

		ctx = context.WithValue(ctx, "userID", session.UserID)
		ctx = context.WithValue(ctx, "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
