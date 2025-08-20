package entity

import "errors"

// Erreurs du domaine User
var (
	ErrUserNotFound      = errors.New("utilisateur non trouvé")
	ErrUserAlreadyExists = errors.New("utilisateur déjà existant")
	ErrInvalidUserData   = errors.New("données utilisateur invalides")
	ErrUserInactive      = errors.New("utilisateur inactif")
)

// Erreurs du domaine File
var (
	ErrFileNotFound     = errors.New("fichier non trouvé")
	ErrFileUploadFailed = errors.New("échec de l'upload du fichier")
	ErrInvalidFileType  = errors.New("type de fichier non autorisé")
	ErrFileTooLarge     = errors.New("fichier trop volumineux")
)

// Erreurs génériques
var (
	ErrInternalServer = errors.New("erreur interne du serveur")
	ErrValidation     = errors.New("erreur de validation")
	ErrUnauthorized   = errors.New("non autorisé")
	ErrForbidden      = errors.New("accès interdit")
	ErrBadRequest     = errors.New("requête invalide")
)
