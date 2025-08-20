package entity

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

// CustomValidator étend le validateur standard avec des règles personnalisées
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator crée une nouvelle instance de validateur personnalisé
func NewCustomValidator() *CustomValidator {
	v := validator.New()

	// Enregistrer les validations personnalisées
	registerCustomValidations(v)

	return &CustomValidator{
		validator: v,
	}
}

// Validate valide une structure
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// registerCustomValidations enregistre les validations personnalisées
func registerCustomValidations(v *validator.Validate) {
	// Validation pour les montants positifs
	v.RegisterValidation("positive_amount", validatePositiveAmount)

	// Validation pour les dates futures
	v.RegisterValidation("future_date", validateFutureDate)

	// Validation pour les UUIDs valides
	v.RegisterValidation("valid_uuid", validateUUID)

	// Validation pour les priorités de tâches
	v.RegisterValidation("task_priority", validateTaskPriority)

	// Validation pour les statuts de tâches
	v.RegisterValidation("task_status", validateTaskStatus)

	// Validation pour les types de transactions
	v.RegisterValidation("transaction_type", validateTransactionType)

	// Validation pour les types de comptes
	v.RegisterValidation("account_type", validateAccountType)

	// Validation pour les types de catégories
	v.RegisterValidation("category_type", validateCategoryType)

	// Validation pour les fréquences
	v.RegisterValidation("frequency", validateFrequency)
}

// validatePositiveAmount valide qu'un montant est positif
func validatePositiveAmount(fl validator.FieldLevel) bool {
	if amount, ok := fl.Field().Interface().(float64); ok {
		return amount > 0
	}
	return false
}

// validateFutureDate valide qu'une date est dans le futur
func validateFutureDate(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(time.Time); ok {
		return date.After(time.Now())
	}
	return false
}

// validateUUID valide qu'un UUID est valide
func validateUUID(fl validator.FieldLevel) bool {
	if uuidStr, ok := fl.Field().Interface().(string); ok {
		_, err := time.Parse("2006-01-02T15:04:05Z07:00", uuidStr)
		return err == nil
	}
	return false
}

// validateTaskPriority valide les priorités de tâches
func validateTaskPriority(fl validator.FieldLevel) bool {
	if priority, ok := fl.Field().Interface().(string); ok {
		validPriorities := []string{"low", "medium", "high"}
		for _, valid := range validPriorities {
			if priority == valid {
				return true
			}
		}
	}
	return false
}

// validateTaskStatus valide les statuts de tâches
func validateTaskStatus(fl validator.FieldLevel) bool {
	if status, ok := fl.Field().Interface().(string); ok {
		validStatuses := []string{"expired", "done", "incoming", "running"}
		for _, valid := range validStatuses {
			if status == valid {
				return true
			}
		}
	}
	return false
}

// validateTransactionType valide les types de transactions
func validateTransactionType(fl validator.FieldLevel) bool {
	if txType, ok := fl.Field().Interface().(string); ok {
		validTypes := []string{"income", "expense"}
		for _, valid := range validTypes {
			if txType == valid {
				return true
			}
		}
	}
	return false
}

// validateAccountType valide les types de comptes
func validateAccountType(fl validator.FieldLevel) bool {
	if accountType, ok := fl.Field().Interface().(string); ok {
		validTypes := []string{"checking", "savings", "mobile_money"}
		for _, valid := range validTypes {
			if accountType == valid {
				return true
			}
		}
	}
	return false
}

// validateCategoryType valide les types de catégories
func validateCategoryType(fl validator.FieldLevel) bool {
	if categoryType, ok := fl.Field().Interface().(string); ok {
		validTypes := []string{"income", "expense", "task"}
		for _, valid := range validTypes {
			if categoryType == valid {
				return true
			}
		}
	}
	return false
}

// validateFrequency valide les fréquences
func validateFrequency(fl validator.FieldLevel) bool {
	if frequency, ok := fl.Field().Interface().(string); ok {
		validFrequencies := []string{"weekly", "monthly", "yearly"}
		for _, valid := range validFrequencies {
			if frequency == valid {
				return true
			}
		}
	}
	return false
}

// ValidationError représente une erreur de validation
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// ValidationErrors représente une liste d'erreurs de validation
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// NewValidationErrors crée une nouvelle instance d'erreurs de validation
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]ValidationError, 0),
	}
}

// AddError ajoute une erreur de validation
func (ve *ValidationErrors) AddError(field, tag, value, message string) {
	ve.Errors = append(ve.Errors, ValidationError{
		Field:   field,
		Tag:     tag,
		Value:   value,
		Message: message,
	})
}

// HasErrors vérifie s'il y a des erreurs
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

// ErrorMessages retourne les messages d'erreur
func (ve *ValidationErrors) ErrorMessages() []string {
	messages := make([]string, len(ve.Errors))
	for i, err := range ve.Errors {
		messages[i] = err.Message
	}
	return messages
}

// GetFieldError retourne l'erreur pour un champ spécifique
func (ve *ValidationErrors) GetFieldError(field string) *ValidationError {
	for _, err := range ve.Errors {
		if err.Field == field {
			return &err
		}
	}
	return nil
}

// String implémente l'interface Stringer
func (ve *ValidationErrors) String() string {
	if len(ve.Errors) == 0 {
		return "No validation errors"
	}

	result := "Validation errors:\n"
	for _, err := range ve.Errors {
		result += fmt.Sprintf("  %s: %s (value: %s)\n", err.Field, err.Message, err.Value)
	}
	return result
}

// Error implémente l'interface error
func (ve *ValidationErrors) Error() string {
	return ve.String()
}

// Messages d'erreur personnalisés
var errorMessages = map[string]string{
	"required":        "Ce champ est obligatoire",
	"min":             "La valeur minimale est %s",
	"max":             "La valeur maximale est %s",
	"email":           "Format d'email invalide",
	"uuid":            "Format UUID invalide",
	"oneof":           "Valeur non autorisée. Valeurs acceptées: %s",
	"gt":              "La valeur doit être supérieure à %s",
	"gte":             "La valeur doit être supérieure ou égale à %s",
	"lt":              "La valeur doit être inférieure à %s",
	"lte":             "La valeur doit être inférieure ou égale à %s",
	"len":             "La longueur doit être exactement %s",
	"datetime":        "Format de date invalide",
	"positive_amount": "Le montant doit être positif",
	"future_date":     "La date doit être dans le futur",
}

// GetErrorMessage retourne le message d'erreur pour un tag donné
func GetErrorMessage(tag string, param string) string {
	if message, exists := errorMessages[tag]; exists {
		if param != "" {
			return fmt.Sprintf(message, param)
		}
		return message
	}
	return fmt.Sprintf("Erreur de validation: %s", tag)
}
