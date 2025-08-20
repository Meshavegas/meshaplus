# My Go Backend - Clean Architecture

🚀 **Backend Go moderne suivant la Clean Architecture (Hexagonal Architecture)**

Une API REST complète avec gestion d'utilisateurs, upload de fichiers, documentation Swagger automatique, et communication avec des services externes.

## 🏗️ Architecture

```
my-go-backend/
├── cmd/api/                    # Point d'entrée de l'application
├── internal/
│   ├── domain/                 # Couche métier (entities, interfaces)
│   │   ├── entity/             # Entités métier
│   │   └── repository/         # Interfaces des repositories
│   ├── usecase/                # Cas d'usage (logique applicative)
│   ├── handler/                # Contrôleurs HTTP (API REST)
│   ├── repository/             # Implémentations des repositories
│   │   └── postgres/           # Implémentation PostgreSQL
│   ├── service/                # Clients pour services externes
│   └── infra/                  # Infrastructure (DB, stockage)
├── pkg/                        # Packages partagés
│   ├── config/                 # Configuration (Viper)
│   ├── logger/                 # Logger structuré (Zap)
│   ├── middleware/             # Middlewares HTTP
│   └── response/               # Standardisation des réponses
├── test/                       # Tests (unitaires + intégration)
└── docs/                       # Documentation Swagger générée
```

## ✨ Fonctionnalités

- 🏛️ **Clean Architecture** - Séparation claire des responsabilités
- 📚 **Documentation Swagger** - Interface interactive automatique
- 🗄️ **PostgreSQL + GORM** - Base de données relationnelle
- 🚀 **Redis** - Cache et sessions
- 📁 **Upload de fichiers** - Stockage local ou cloud (S3)
- 🌐 **API REST complète** - CRUD utilisateurs + gestion fichiers
- 🔧 **Configuration flexible** - Viper + variables d'environnement
- 📊 **Logging structuré** - Zap avec niveaux configurables
- 🧪 **Tests** - Unitaires et d'intégration
- 🐳 **Docker** - Environnement de développement complet

## 🚀 Démarrage rapide

### Prérequis

- Go 1.21+
- Docker & Docker Compose
- Make (optionnel, mais recommandé)

### Installation

```bash
# Cloner le projet
git clone <votre-repo>
cd my-go-backend

# Configuration initiale complète
make setup

# Ou étape par étape :
# 1. Installer les dépendances
make deps

# 2. Démarrer les services (PostgreSQL, Redis)
make docker-up

# 3. Générer la documentation Swagger
make swagger

# 4. Lancer l'application
make run
```

### Accès aux services

- 🌐 **API** : http://localhost:8080
- 📚 **Swagger UI** : http://localhost:8080/swagger/index.html
- 🗄️ **Adminer** (PostgreSQL) : http://localhost:8081
- 🔧 **PostgreSQL** : localhost:5432
- 📊 **Redis** : localhost:6379

## 📋 API Endpoints

### Utilisateurs

```http
POST   /api/v1/users          # Créer un utilisateur
GET    /api/v1/users          # Lister les utilisateurs (paginé)
GET    /api/v1/users/{id}     # Récupérer un utilisateur
PUT    /api/v1/users/{id}     # Mettre à jour un utilisateur
DELETE /api/v1/users/{id}     # Supprimer un utilisateur
```

### Fichiers

```http
POST   /api/v1/files/upload   # Upload d'un fichier
GET    /api/v1/files/{id}     # Récupérer un fichier
DELETE /api/v1/files/{id}     # Supprimer un fichier
```

### Health Check

```http
GET    /api/v1/health         # Vérification santé de l'API
```

## 🔧 Configuration

Copiez `.env.example` vers `.env` et ajustez selon vos besoins :

```bash
cp .env.example .env
```

### Configurations par environnement

- `configs/config.dev.yaml` - Développement
- `configs/config.prod.yaml` - Production  
- `configs/config.test.yaml` - Tests

## 🧪 Tests

```bash
# Tests unitaires
make test

# Tests avec couverture
make test-coverage

# Tests d'intégration (lance les services Docker)
make test-integration
```

## 📦 Commandes Make

```bash
make help              # Affiche toutes les commandes disponibles
make deps               # Installe les dépendances
make build              # Compile l'application
make run                # Lance l'application
make dev                # Mode développement (hot reload avec Air)
make test               # Lance les tests
make swagger            # Génère la documentation Swagger
make docker-up          # Démarre les services Docker
make docker-down        # Arrête les services Docker
make clean              # Nettoie les fichiers générés
make lint               # Lance le linter
make format             # Formate le code
```

## 🐳 Docker

### Développement

```bash
# Services infrastructure seulement
make docker-up

# Développement de l'app en local
make run
```

###