package entity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User représente un utilisateur du système
type User struct {
	tableName    struct{}   `pg:"users"` // Utilisé par go-pg pour spécifier le nom de la table
	ID           uuid.UUID  `pg:"id,pk,type:uuid,default:gen_random_uuid()" json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name         string     `pg:"name,notnull" json:"name" validate:"required,min=2,max=100" example:"John Doe"`
	Email        string     `pg:"email,unique,notnull" json:"email" validate:"required,email" example:"john.doe@example.com"`
	Avatar       string     `pg:"avatar,default:''" json:"avatar,omitempty" example:"https://example.com/avatar.jpg"`
	PasswordHash string     `pg:"password_hash,notnull" json:"-" validate:"required"` // Le hash du mot de passe n'est jamais exposé en JSON
	CreatedAt    time.Time  `pg:"created_at,default:now()" json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt    time.Time  `pg:"updated_at,default:now()" json:"updated_at" example:"2023-01-01T12:00:00Z"`
	DeletedAt    *time.Time `pg:"deleted_at,soft_delete" json:"-"`
}

// CreateUserRequest représente les données pour créer un utilisateur
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Avatar   string `json:"avatar,omitempty" example:"https://example.com/avatar.jpg"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
}

// UpdateUserRequest représente les données pour mettre à jour un utilisateur
type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=2,max=100" example:"John Doe"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email" example:"john.doe@example.com"`
	Avatar   *string `json:"avatar,omitempty" example:"https://example.com/avatar.jpg"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6" example:"newpassword123"`
}

// UserListResponse représente une liste paginée d'utilisateurs
type UserListResponse struct {
	Users      []User `json:"users"`
	Total      int64  `json:"total" example:"100"`
	Page       int    `json:"page" example:"1"`
	PageSize   int    `json:"page_size" example:"10"`
	TotalPages int    `json:"total_pages" example:"10"`
}

// UserQuery représente les paramètres de recherche d'utilisateurs
type UserQuery struct {
	Page     int    `json:"page" query:"page" validate:"min=1" example:"1"`
	PageSize int    `json:"page_size" query:"page_size" validate:"min=1,max=100" example:"10"`
	Search   string `json:"search,omitempty" query:"search" example:"john"`
}

// IsValidForCreation vérifie si l'utilisateur peut être créé
func (u *User) IsValidForCreation() error {
	if u.Name == "" {
		return ErrInvalidUserData
	}
	if u.Email == "" {
		return ErrInvalidUserData
	}
	if u.PasswordHash == "" {
		return ErrInvalidUserData
	}
	return nil
}

// BeforeInsert hook go-pg pour générer l'UUID et mettre à jour les timestamps
func (u *User) BeforeInsert(ctx context.Context) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	return nil
}

// BeforeUpdate hook go-pg pour mettre à jour le timestamp
func (u *User) BeforeUpdate(ctx context.Context) error {
	u.UpdatedAt = time.Now()
	return nil
}
