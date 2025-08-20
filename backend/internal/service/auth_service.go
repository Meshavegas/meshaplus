package service

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"backend/pkg/logger"
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService gère l'authentification des utilisateurs
type AuthService struct {
	userRepo   repository.UserRepository
	jwtService *JWTService
	logger     logger.Logger
}

// NewAuthService crée une nouvelle instance de AuthService
func NewAuthService(userRepo repository.UserRepository, jwtService *JWTService, logger logger.Logger) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
		logger:     logger,
	}
}

// LoginRequest représente les données de connexion
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}

// LoginResponse représente la réponse de connexion
type LoginResponse struct {
	User         *entity.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"`
}

// RegisterRequest représente les données d'inscription
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
	Avatar   string `json:"avatar,omitempty" example:"https://example.com/avatar.jpg"`
}

// RefreshTokenRequest représente la demande de rafraîchissement de token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Login authentifie un utilisateur
func (a *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Récupérer l'utilisateur par email
	user, err := a.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		a.logger.Warn("Tentative de connexion avec email inexistant", logger.String("email", req.Email))
		return nil, fmt.Errorf("email ou mot de passe incorrect")
	}

	// Vérifier le mot de passe
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		a.logger.Warn("Tentative de connexion avec mot de passe incorrect",
			logger.String("email", req.Email),
			logger.Error(err),
		)
		return nil, fmt.Errorf("email ou mot de passe incorrect")
	}

	// Générer les tokens
	accessToken, err := a.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		a.logger.Error("Erreur génération token d'accès", logger.Error(err))
		return nil, fmt.Errorf("erreur génération token: %w", err)
	}

	refreshToken, err := a.jwtService.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		a.logger.Error("Erreur génération token de rafraîchissement", logger.Error(err))
		return nil, fmt.Errorf("erreur génération token: %w", err)
	}

	a.logger.Info("Connexion réussie",
		logger.String("user_id", user.ID.String()),
		logger.String("email", user.Email),
	)

	return &LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    a.jwtService.config.ExpirationHours * 3600, // en secondes
	}, nil
}

// Register inscrit un nouvel utilisateur
func (a *AuthService) Register(ctx context.Context, req RegisterRequest) (*LoginResponse, error) {
	// Vérifier si l'email existe déjà
	exists, err := a.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		a.logger.Error("Erreur vérification email", logger.Error(err))
		return nil, fmt.Errorf("erreur vérification email: %w", err)
	}

	if exists {
		a.logger.Warn("Tentative d'inscription avec email existant", logger.String("email", req.Email))
		return nil, fmt.Errorf("email déjà utilisé")
	}

	// Hasher le mot de passe
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		a.logger.Error("Erreur hashage mot de passe", logger.Error(err))
		return nil, fmt.Errorf("erreur hashage mot de passe: %w", err)
	}

	// Créer l'utilisateur
	user := &entity.User{
		Name:         req.Name,
		Email:        req.Email,
		Avatar:       req.Avatar,
		PasswordHash: string(passwordHash),
	}

	// Sauvegarder l'utilisateur
	if err := a.userRepo.Create(ctx, user); err != nil {
		a.logger.Error("Erreur création utilisateur", logger.Error(err))
		return nil, fmt.Errorf("erreur création utilisateur: %w", err)
	}

	// Générer les tokens
	accessToken, err := a.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		a.logger.Error("Erreur génération token d'accès", logger.Error(err))
		return nil, fmt.Errorf("erreur génération token: %w", err)
	}

	refreshToken, err := a.jwtService.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		a.logger.Error("Erreur génération token de rafraîchissement", logger.Error(err))
		return nil, fmt.Errorf("erreur génération token: %w", err)
	}

	a.logger.Info("Inscription réussie",
		logger.String("user_id", user.ID.String()),
		logger.String("email", user.Email),
	)

	return &LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    a.jwtService.config.ExpirationHours * 3600, // en secondes
	}, nil
}

// RefreshToken rafraîchit un token d'accès
func (a *AuthService) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*LoginResponse, error) {
	// Valider le token de rafraîchissement
	claims, err := a.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		a.logger.Warn("Token de rafraîchissement invalide", logger.Error(err))
		return nil, fmt.Errorf("token de rafraîchissement invalide")
	}

	// Vérifier que c'est bien un token de rafraîchissement
	if claims.Type != "refresh" {
		a.logger.Warn("Token n'est pas un token de rafraîchissement", logger.String("type", claims.Type))
		return nil, fmt.Errorf("token n'est pas un token de rafraîchissement")
	}

	// Récupérer l'utilisateur
	user, err := a.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		a.logger.Error("Utilisateur non trouvé lors du rafraîchissement", logger.Error(err))
		return nil, fmt.Errorf("utilisateur non trouvé")
	}

	// Générer un nouveau token d'accès
	accessToken, err := a.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		a.logger.Error("Erreur génération nouveau token d'accès", logger.Error(err))
		return nil, fmt.Errorf("erreur génération token: %w", err)
	}

	// Générer un nouveau token de rafraîchissement
	refreshToken, err := a.jwtService.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		a.logger.Error("Erreur génération nouveau token de rafraîchissement", logger.Error(err))
		return nil, fmt.Errorf("erreur génération token: %w", err)
	}

	a.logger.Info("Token rafraîchi avec succès",
		logger.String("user_id", user.ID.String()),
		logger.String("email", user.Email),
	)

	return &LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    a.jwtService.config.ExpirationHours * 3600, // en secondes
	}, nil
}

// GetCurrentUser récupère l'utilisateur actuel à partir du contexte
func (a *AuthService) GetCurrentUser(ctx context.Context) (*entity.User, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("utilisateur non authentifié")
	}

	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("utilisateur non trouvé: %w", err)
	}

	return user, nil
}
