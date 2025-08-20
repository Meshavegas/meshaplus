package response

import (
	"encoding/json"
	"net/http"
	"time"
)

// Response représente une réponse API standard
type Response struct {
	Success   bool        `json:"success" example:"true"`
	Message   string      `json:"message" example:"Opération réussie"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp" example:"2023-01-01T12:00:00Z"`
}

// ErrorResponse représente une réponse d'erreur
type ErrorResponse struct {
	Success   bool      `json:"success" example:"false"`
	Message   string    `json:"message" example:"Une erreur est survenue"`
	Error     string    `json:"error,omitempty" example:"Détails de l'erreur"`
	Code      string    `json:"code,omitempty" example:"USER_NOT_FOUND"`
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T12:00:00Z"`
}

// ValidationError représente une erreur de validation
type ValidationError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Email invalide"`
}

// ValidationErrorResponse représente une réponse d'erreur de validation
type ValidationErrorResponse struct {
	Success   bool              `json:"success" example:"false"`
	Message   string            `json:"message" example:"Erreurs de validation"`
	Errors    []ValidationError `json:"errors"`
	Timestamp time.Time         `json:"timestamp" example:"2023-01-01T12:00:00Z"`
}

// Success envoie une réponse de succès
func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}

	json.NewEncoder(w).Encode(response)
}

// Error envoie une réponse d'erreur
func Error(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
	}

	if err != nil {
		response.Error = err.Error()
	}

	json.NewEncoder(w).Encode(response)
}

// ErrorWithCode envoie une réponse d'erreur avec un code
func ErrorWithCode(w http.ResponseWriter, statusCode int, message string, code string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Success:   false,
		Message:   message,
		Code:      code,
		Timestamp: time.Now(),
	}

	if err != nil {
		response.Error = err.Error()
	}

	json.NewEncoder(w).Encode(response)
}

// ValidationError envoie une réponse d'erreur de validation
func ValidationErrors(w http.ResponseWriter, message string, errors []ValidationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := ValidationErrorResponse{
		Success:   false,
		Message:   message,
		Errors:    errors,
		Timestamp: time.Now(),
	}

	json.NewEncoder(w).Encode(response)
}

// JSON envoie une réponse JSON personnalisée
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
