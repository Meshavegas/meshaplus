package middleware

import (
	"backend/internal/service"
	"backend/pkg/jwt"
	"backend/pkg/logger"
	"context"
	"net/http"
	"strings"
)

// AuthMiddleware gère l'authentification JWT
type AuthMiddleware struct {
	jwtService *service.JWTService
	logger     logger.Logger
}

// NewAuthMiddleware crée une nouvelle instance de AuthMiddleware
func NewAuthMiddleware(jwtService *service.JWTService, logger logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		logger:     logger,
	}
}

// Authenticate middleware pour authentifier les requêtes
func (a *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extraire le token du header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			a.logger.Warn("Header Authorization manquant")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Vérifier le format "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			a.logger.Warn("Format Authorization header invalide")
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]

		// Valider le token
		claims, err := a.jwtService.ValidateToken(tokenString)
		if err != nil {
			a.logger.Warn("Token invalide", logger.Error(err))
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Vérifier que c'est un token d'accès
		if claims.Type != "access" {
			a.logger.Warn("Token n'est pas un token d'accès", logger.String("type", claims.Type))
			http.Error(w, "Invalid token type", http.StatusUnauthorized)
			return
		}

		// Ajouter les informations utilisateur au contexte
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "claims", claims)

		a.logger.Debug("Utilisateur authentifié",
			logger.String("user_id", claims.UserID.String()),
			logger.String("email", claims.Email),
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth middleware pour authentifier optionnellement les requêtes
func (a *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extraire le token du header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// Pas de token, continuer sans authentification
			next.ServeHTTP(w, r)
			return
		}

		// Vérifier le format "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			// Format invalide, continuer sans authentification
			next.ServeHTTP(w, r)
			return
		}

		tokenString := tokenParts[1]

		// Valider le token
		claims, err := a.jwtService.ValidateToken(tokenString)
		if err != nil {
			// Token invalide, continuer sans authentification
			a.logger.Debug("Token invalide dans OptionalAuth", logger.Error(err))
			next.ServeHTTP(w, r)
			return
		}

		// Vérifier que c'est un token d'accès
		if claims.Type != "access" {
			// Type invalide, continuer sans authentification
			next.ServeHTTP(w, r)
			return
		}

		// Ajouter les informations utilisateur au contexte
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "claims", claims)

		a.logger.Debug("Utilisateur authentifié (optionnel)",
			logger.String("user_id", claims.UserID.String()),
			logger.String("email", claims.Email),
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole middleware pour vérifier les rôles (pour une utilisation future)
func (a *AuthMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Récupérer les claims du contexte
			claims, ok := r.Context().Value("claims").(*jwt.Claims)
			if !ok {
				a.logger.Warn("Claims non trouvés dans le contexte")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// TODO: Implémenter la vérification des rôles
			// Pour l'instant, on accepte tous les utilisateurs authentifiés
			a.logger.Debug("Vérification des rôles",
				logger.String("user_id", claims.UserID.String()),
				logger.String("roles", strings.Join(roles, ",")),
			)

			next.ServeHTTP(w, r)
		})
	}
}
