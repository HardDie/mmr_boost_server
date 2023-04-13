package utils

import (
	"context"

	"github.com/HardDie/mmr_boost_server/internal/entity"
)

func GetUserIDFromContext(ctx context.Context) int32 {
	return ctx.Value("userID").(int32)
}
func GetAccessTokenFromContext(ctx context.Context) *entity.AccessToken {
	return ctx.Value("session").(*entity.AccessToken)
}
