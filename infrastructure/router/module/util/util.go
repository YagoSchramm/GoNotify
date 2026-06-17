package util

import (
	"net/http"

	"github.com/YagoSchramm/GoNotify/domain/entity/derr"
	"github.com/YagoSchramm/GoNotify/infrastructure/foundation/jwt"
	"github.com/google/uuid"
)

func UserIDFromContext(r *http.Request) (uuid.UUID, error) {
	claims, ok := r.Context().Value("user_claims").(*jwt.Claims)
	if !ok {
		return uuid.Nil, derr.UnauthorizedError
	}
	return claims.UserID, nil
}
