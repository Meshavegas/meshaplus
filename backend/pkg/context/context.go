package context

import (
	"backend/pkg/jwt"
	"context"

	"github.com/google/uuid"
)

// GetUserIDFromContext récupère l'ID utilisateur du contexte
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	return userID, ok
}

// GetEmailFromContext récupère l'email utilisateur du contexte
func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value("email").(string)
	return email, ok
}

// GetClaimsFromContext récupère les claims du contexte
func GetClaimsFromContext(ctx context.Context) (*jwt.Claims, bool) {
	claims, ok := ctx.Value("claims").(*jwt.Claims)
	return claims, ok
}
