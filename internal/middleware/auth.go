package middleware

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/service"
	"github.com/HardDie/mmr_boost_server/internal/utils"
	pb "github.com/HardDie/mmr_boost_server/pkg/proto/server"
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

// AdminMiddleware Validate that current user is admin.
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleID := utils.ContextGetRoleID(r.Context())
		if roleID != int32(pb.UserRoleID_admin) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ManagementMiddleware Validate that current user is admin or manager.
func ManagementMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleID := utils.ContextGetRoleID(r.Context())
		if roleID != int32(pb.UserRoleID_admin) &&
			roleID != int32(pb.UserRoleID_manager) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
