# Variables
APP_NAME=backend
BINARY_NAME=api
MAIN_PATH=./cmd/api
GO_FILES=$(shell find . -name "*.go" -type f -not -path "./vendor/*")

# Couleurs pour les messages
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help build run dev test test-coverage clean deps swagger docker docker-up docker-down migrate lint format test-auth check-swagger

# Aide
help: ## Affiche cette aide
	@echo "$(BLUE)Commandes disponibles:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'

# Installation des dépendances
deps: ## Installe les dépendances Go
	@echo "$(YELLOW)Installation des dépendances...$(NC)"
	@go mod download
	@go mod tidy
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)Dépendances installées avec succès$(NC)"

# Build
build: ## Compile l'application
	@echo "$(YELLOW)Compilation de l'application...$(NC)"
	@go build -o bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Application compilée: bin/$(BINARY_NAME)$(NC)"

# Run
run: swagger ## Lance l'application
	@echo "$(YELLOW)Démarrage de l'application...$(NC)"
	@go run $(MAIN_PATH)

# Mode développement avec rechargement automatique
dev: swagger ## Lance l'application en mode développement
	@echo "$(YELLOW)Mode développement (Air required: go install github.com/cosmtrek/air@latest)$(NC)"
	@air -c .air.toml

# Tests
test: ## Lance tous les tests
	@echo "$(YELLOW)Exécution des tests...$(NC)"
	@go test -v ./...

# Tests avec couverture
test-coverage: ## Lance les tests avec couverture
	@echo "$(YELLOW)Tests avec couverture...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Rapport de couverture: coverage.html$(NC)"

# Tests d'intégration
test-integration: docker-up ## Lance les tests d'intégration
	@echo "$(YELLOW)Tests d'intégration...$(NC)"
	@sleep 5 # Attendre que les services soient prêts
	@go test -v ./test/integration/...
	@$(MAKE) docker-down

# Documentation Swagger
swagger: ## Génère la documentation Swagger
	@echo "$(YELLOW)Génération de la documentation Swagger...$(NC)"
	@~/go/bin/swag init -g $(MAIN_PATH)/main.go -o ./docs
	@echo "$(GREEN)Documentation générée: http://localhost:8080/swagger/index.html$(NC)"

# Linting
lint: ## Lance le linter
	@echo "$(YELLOW)Linting du code...$(NC)"
	@golangci-lint run ./...

# Formatage
format: ## Formate le code Go
	@echo "$(YELLOW)Formatage du code...$(NC)"
	@gofmt -s -w $(GO_FILES)
	@goimports -w $(GO_FILES)
	@echo "$(GREEN)Code formaté$(NC)"

# Migration de base de données
migrate: ## Lance les migrations
	@echo "$(YELLOW)Migration de la base de données...$(NC)"
	@./scripts/run_migration.sh
	@echo "$(GREEN)Migration terminée$(NC)"

# Tests d'authentification
test-auth: ## Lance les tests d'authentification JWT
	@echo "$(YELLOW)Tests d'authentification JWT...$(NC)"
	@./scripts/test_auth.sh

# Tests des routes Expense
test-expense: ## Lance les tests des routes Expense
	@echo "$(YELLOW)Tests des routes Expense...$(NC)"
	@./scripts/test_expense.sh

# Tests des routes Revenue
test-revenue: ## Lance les tests des routes Revenue
	@echo "$(YELLOW)Tests des routes Revenue...$(NC)"
	@./scripts/test_revenue.sh

# Tests des routes SavingStrategy
test-saving-strategy: ## Lance les tests des routes SavingStrategy
	@echo "$(YELLOW)Tests des routes SavingStrategy...$(NC)"
	@./scripts/test_saving_strategy.sh

# Vérification Swagger
check-swagger: ## Vérifie que la documentation Swagger fonctionne
	@echo "$(YELLOW)Vérification de la documentation Swagger...$(NC)"
	@./scripts/check_swagger.sh

# Docker
docker: ## Build l'image Docker
	@echo "$(YELLOW)Construction de l'image Docker...$(NC)"
	@docker build -t $(APP_NAME):latest .
	@echo "$(GREEN)Image Docker construite: $(APP_NAME):latest$(NC)"

# Docker Compose - Up
docker-up: ## Lance les services avec Docker Compose
	@echo "$(YELLOW)Démarrage des services Docker...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)Services démarrés$(NC)"
	@echo "$(BLUE)PostgreSQL: localhost:5432$(NC)"
	@echo "$(BLUE)Redis: localhost:6379$(NC)"
	@echo "$(BLUE)Adminer: http://localhost:8081$(NC)"

# Docker Compose - Down
docker-down: ## Arrête les services Docker Compose
	@echo "$(YELLOW)Arrêt des services Docker...$(NC)"
	@docker-compose down
	@echo "$(GREEN)Services arrêtés$(NC)"

# Docker Compose - Logs
docker-logs: ## Affiche les logs des services
	@docker-compose logs -f

# Nettoyage
clean: ## Nettoie les fichiers générés
	@echo "$(YELLOW)Nettoyage...$(NC)"
	@rm -rf bin/
	@rm -rf coverage.out coverage.html
	@rm -rf docs/docs.go docs/swagger.json docs/swagger.yaml
	@go clean -testcache
	@echo "$(GREEN)Nettoyage terminé$(NC)"

# Setup du projet
setup: deps docker-up migrate swagger ## Configuration initiale du projet
	@echo "$(GREEN)Projet configuré avec succès!$(NC)"
	@echo "$(BLUE)Pour démarrer: make run$(NC)"
	@echo "$(BLUE)Swagger: http://localhost:8080/swagger/index.html$(NC)"

# Commande par défaut
all: clean deps swagger build test ## Clean, deps, swagger, build et test