package middleware

import (
	"market_server/common/middleware"
	"net/http"
)

type JwtAuthMiddleware struct {
	Secret string
}

func NewJwtAuthMiddleware(secret string) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{
		Secret: secret,
	}
}

func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.NewCommonJwtAuthMiddleware(m.Secret).Handle(next)(w, r)
	}
}
