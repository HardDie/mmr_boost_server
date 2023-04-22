package utils

import (
	"context"

	"github.com/HardDie/mmr_boost_server/internal/entity"
)

type mmrBoostType string

func ContextSetUserID(ctx context.Context, userID int32) context.Context {
	return context.WithValue(ctx, mmrBoostType("userID"), userID)
}
func ContextGetUserID(ctx context.Context) int32 {
	return ctx.Value(mmrBoostType("userID")).(int32)
}

func ContextSetRoleID(ctx context.Context, roleID int32) context.Context {
	return context.WithValue(ctx, mmrBoostType("roleID"), roleID)
}
func ContextGetRoleID(ctx context.Context) int32 {
	return ctx.Value(mmrBoostType("roleID")).(int32)
}

func ContextSetSession(ctx context.Context, session *entity.AccessToken) context.Context {
	return context.WithValue(ctx, mmrBoostType("session"), session)
}
func ContextGetSession(ctx context.Context) *entity.AccessToken {
	return ctx.Value(mmrBoostType("session")).(*entity.AccessToken)
}
