backend/
├── cmd/
│   └── api/
│       └── main.go                 # Point d'entrée de l'application
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── user.go            # Entité User
│   │   │   └── file.go            # Entité File
│   │   └── repository/
│   │       ├── user_repository.go  # Interface repository User
│   │       └── file_repository.go  # Interface repository File
│   ├── usecase/
│   │   ├── user_usecase.go        # Logique métier User
│   │   └── file_usecase.go        # Logique métier File
│   ├── handler/
│   │   ├── user_handler.go        # Contrôleur HTTP User
│   │   ├── file_handler.go        # Contrôleur HTTP File
│   │   └── health_handler.go      # Health check
│   ├── repository/
│   │   ├── postgres/
│   │   │   ├── user_repository.go # Implémentation PostgreSQL User
│   │   │   └── file_repository.go # Implémentation PostgreSQL File
│   │   └── redis/
│   │       └── cache_repository.go # Implémentation Redis
│   ├── service/
│   │   └── external_service.go    # Client API externe
│   └── infra/
│       ├── database/
│       │   ├── postgres.go        # Connexion PostgreSQL
│       │   └── redis.go           # Connexion Redis
│       └── storage/
│           ├── local.go           # Stockage local
│           └── s3.go              # Stockage S3
├── pkg/
│   ├── config/
│   │   └── config.go              # Configuration avec Viper
│   ├── logger/
│   │   └── logger.go              # Logger avec Zap
│   ├── middleware/
│   │   ├── cors.go                # Middleware CORS
│   │   ├── logger.go              # Middleware logging
│   │   └── auth.go                # Middleware auth (optionnel)
│   ├── response/
│   │   └── response.go            # Standardisation des réponses
│   └── utils/
│       ├── validator.go           # Validation des données
│       └── helpers.go             # Fonctions utilitaires
├── test/
│   ├── mocks/
│   │   ├── user_repository_mock.go
│   │   └── external_service_mock.go
│   ├── unit/
│   │   ├── usecase/
│   │   │   └── user_usecase_test.go
│   │   └── handler/
│   │       └── user_handler_test.go
│   └── integration/
│       └── api_integration_test.go
├── docs/                          # Documentation Swagger générée
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── configs/
│   ├── config.dev.yaml
│   ├── config.prod.yaml
│   └── config.test.yaml
├── scripts/
│   ├── setup.sh                   # Script d'installation
│   └── migrate.sql                # Migrations SQL
├── .env.example                   # Variables d'environnement exemple
├── .gitignore
├── docker-compose.yml             # Services (PostgreSQL, Redis)
├── Dockerfile
├── Makefile                       # Commandes de développement
├── go.mod
├── go.sum
└── README.md                      # Documentation du projet







