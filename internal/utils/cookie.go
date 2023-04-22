package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/HardDie/mmr_boost_server/internal/logger"
)

func SetGRPCSessionCookie(ctx context.Context, session string) {
	cookie := http.Cookie{
		Name:     "session",
		Path:     "/",
		Value:    session,
		HttpOnly: true,
	}
	err := grpc.SetHeader(ctx, metadata.Pairs("Set-Cookie", cookie.String()))
	if err != nil {
		logger.Error.Println("error set cookie:", err.Error())
	}
}
func DeleteGRPCSessionCookie(ctx context.Context) {
	cookie := http.Cookie{
		Name:     "session",
		Path:     "/",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	err := grpc.SetHeader(ctx, metadata.Pairs("Set-Cookie", cookie.String()))
	if err != nil {
		logger.Error.Println("error delete cookie:", err.Error())
	}
}

func GetCookie(r *http.Request) string {
	cookie, err := r.Cookie("session")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func GetBearer(r *http.Request) string {
	header := r.Header.Get("Authorization")
	return strings.ReplaceAll(header, "Bearer ", "")
}

func GenerateSessionKey() (string, error) {
	sessionLen := 32
	b := make([]byte, sessionLen)
	nRead, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("read random: %w", err)
	}
	if nRead != sessionLen {
		return "", fmt.Errorf("bad length")
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
