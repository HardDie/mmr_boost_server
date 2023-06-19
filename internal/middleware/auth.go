package middleware

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		// List of public routes
		if r.URL.Path == "/api/v1/auth/login" ||
			r.URL.Path == "/api/v1/auth/register" ||
			r.URL.Path == "/api/v1/auth/validate_email" ||
			r.URL.Path == "/api/v1/auth/send_validation_email" ||
			r.URL.Path == "/api/v1/system/swagger" {
			next.ServeHTTP(w, r)
			return
		}

		// Extract token from cookie
		token := utils.GetCookie(r)

		// If we got no cookie
		if token == "" {
			// Extract token from Authorization
			token = utils.GetBearer(r)
			if token == "" {
				http.Error(w, "Session token is empty", http.StatusUnauthorized)
				return
			}
		}

		// Validate if session is active
		ctx := r.Context()
		user, session, err := m.service.Auth.ValidateCookie(ctx, token)
		if err != nil {
			if status.Convert(err).Code() == codes.PermissionDenied {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else {
				http.Error(w, "Internal error", http.StatusInternalServerError)
			}
			return
		}

		ctx = utils.ContextSetUserID(ctx, session.UserID)
		ctx = utils.ContextSetRoleID(ctx, user.RoleID)
		ctx = utils.ContextSetSession(ctx, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
