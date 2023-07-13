package utils

import (
	"context"

	"github.com/HardDie/mmr_boost_server/internal/entity"
)

type mmrBoostType string

func ContextSetUserID(ctx context.Context, userID int32) context.Context {
	return context.WithValue(ctx, mmrBoostType("userID"), userID)
}
func ContextGetUserID(ctx context.Context) (int32, bool) {
	userID, ok := ctx.Value(mmrBoostType("userID")).(int32)
	return userID, ok
}

func ContextSetRoleID(ctx context.Context, roleID int32) context.Context {
	return context.WithValue(ctx, mmrBoostType("roleID"), roleID)
}
func ContextGetRoleID(ctx context.Context) (int32, bool) {
	roleID, ok := ctx.Value(mmrBoostType("roleID")).(int32)
	return roleID, ok
}

func ContextSetSession(ctx context.Context, session *entity.AccessToken) context.Context {
	return context.WithValue(ctx, mmrBoostType("session"), session)
}
func ContextGetSession(ctx context.Context) (*entity.AccessToken, bool) {
	session, ok := ctx.Value(mmrBoostType("session")).(*entity.AccessToken)
	return session, ok
}
