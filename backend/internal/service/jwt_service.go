package service

import (
	"backend/pkg/config"
	jwtpkg "backend/pkg/jwt"
	"backend/pkg/logger"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTService gère l'authentification JWT
type JWTService struct {
	config config.JWTConfig
	logger logger.Logger
}

// Claims représente les claims JWT
type Claims = jwtpkg.Claims

// NewJWTService crée une nouvelle instance de JWTService
func NewJWTService(config config.JWTConfig, logger logger.Logger) *JWTService {
	return &JWTService{
		config: config,
		logger: logger,
	}
}

// GenerateAccessToken génère un token d'accès
func (j *JWTService) GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(j.config.ExpirationHours) * time.Hour)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "meshaplus-api",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		j.logger.Error("Erreur génération token d'accès", logger.Error(err))
		return "", fmt.Errorf("erreur génération token: %w", err)
	}

	j.logger.Info("Token d'accès généré avec succès",
		logger.String("user_id", userID.String()),
		logger.String("email", email),
	)

	return tokenString, nil
}

// GenerateRefreshToken génère un token de rafraîchissement
func (j *JWTService) GenerateRefreshToken(userID uuid.UUID, email string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(j.config.RefreshExpirationHours) * time.Hour)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "meshaplus-api",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		j.logger.Error("Erreur génération token de rafraîchissement", logger.Error(err))
		return "", fmt.Errorf("erreur génération token de rafraîchissement: %w", err)
	}

	j.logger.Info("Token de rafraîchissement généré avec succès",
		logger.String("user_id", userID.String()),
		logger.String("email", email),
	)

	return tokenString, nil
}

// ValidateToken valide un token JWT
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Vérifier la méthode de signature
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode de signature inattendue: %v", token.Header["alg"])
		}
		return []byte(j.config.SecretKey), nil
	})

	if err != nil {
		j.logger.Error("Erreur validation token", logger.Error(err))
		return nil, fmt.Errorf("token invalide: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token invalide")
	}

	return claims, nil
}

// RefreshToken rafraîchit un token d'accès avec un token de rafraîchissement
func (j *JWTService) RefreshToken(refreshToken string) (string, error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("token de rafraîchissement invalide: %w", err)
	}

	// Vérifier que c'est bien un token de rafraîchissement
	if claims.Type != "refresh" {
		return "", errors.New("token n'est pas un token de rafraîchissement")
	}

	// Générer un nouveau token d'accès
	accessToken, err := j.GenerateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		return "", fmt.Errorf("erreur génération nouveau token d'accès: %w", err)
	}

	return accessToken, nil
}
