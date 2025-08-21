// @title My Go Backend API
// @version 1.0
// @description API REST suivant la Clean Architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"backend/internal/handler"
	"backend/internal/infra/database"
	"backend/internal/repository/postgres"
	"backend/internal/routes"
	"backend/internal/service"
	"backend/internal/usecase"
	"backend/pkg/config"
	"backend/pkg/logger"

	"backend/pkg/middleware"
)

func main() {
	// Charger la configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la configuration: %v", err)
	}

	// Initialiser le logger
	loggerInstance := logger.New(cfg.Logger.Level)
	defer func() {
		if err := loggerInstance.Sync(); err != nil {
			log.Printf("Error syncing logger: %v", err)
		}
	}()

	// Connexion à la base de données
	db, err := database.NewPostgresConnection(cfg.Database, loggerInstance)
	if err != nil {
		loggerInstance.Fatal("Erreur connexion PostgreSQL", logger.Error(err))
	}
	defer database.ClosePostgresConnection(db, loggerInstance)

	// Exécuter les migrations
	if err := database.RunMigrations(db, loggerInstance); err != nil {
		loggerInstance.Fatal("Erreur lors des migrations", logger.Error(err))
	}

	// Connexion à Redis (optionnel) - TODO: implement Redis connection
	// redisClient, err := database.NewRedisConnection(cfg.Redis)
	// if err != nil {
	// 	loggerInstance.Warn("Redis non disponible, cache désactivé", logger.Error(err))
	// }

	// Initialisation du stockage - TODO: implement storage
	// fileStorage := storage.NewLocalStorage(cfg.Storage.LocalPath)

	// Injection des dépendances - Repositories
	userRepo := postgres.NewUserRepository(db)
	taskRepo := postgres.NewTaskRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)
	accountRepo := postgres.NewAccountRepository(db)
	budgetRepo := postgres.NewBudgetRepository(db)
	savingGoalRepo := postgres.NewSavingGoalRepository(db)
	// fileRepo := postgres.NewFileRepository(db) // TODO: implement file repository
	// externalService := service.NewExternalService(cfg.ExternalAPI.BaseURL) // TODO: implement external service

	// Services
	jwtService := service.NewJWTService(cfg.JWT, loggerInstance)
	authService := service.NewAuthService(userRepo, jwtService, loggerInstance)
	taskService := service.NewTaskService(taskRepo, loggerInstance)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo, loggerInstance)
	accountService := service.NewAccountService(accountRepo, loggerInstance)
	budgetService := service.NewBudgetService(budgetRepo, loggerInstance)
	savingGoalService := service.NewSavingGoalService(savingGoalRepo, loggerInstance)

	// Usecases
	userUsecase := usecase.NewUserUsecase(userRepo, loggerInstance)
	// fileUsecase := usecase.NewFileUsecase(fileRepo, fileStorage, loggerInstance) // TODO: implement file usecase

	// Middlewares
	authMiddleware := middleware.NewAuthMiddleware(jwtService, loggerInstance)

	// Handlers
	// authHandler := handler.NewAuthHandler(authService, loggerInstance) // TODO: implement auth routes
	taskHandler := handler.NewTaskHandler(taskService, loggerInstance)
	transactionHandler := handler.NewTransactionHandler(transactionService, loggerInstance)
	accountHandler := handler.NewAccountHandler(accountService, loggerInstance)
	budgetHandler := handler.NewBudgetHandler(budgetService, loggerInstance)
	savingGoalHandler := handler.NewSavingGoalHandler(savingGoalService, loggerInstance)
	// userHandler := handler.NewUserHandler(userUsecase, loggerInstance)
	// fileHandler := handler.NewFileHandler(fileUsecase, loggerInstance) // TODO: implement file handler
	// healthHandler := handler.NewHealthHandler(db, redisClient) // TODO: implement health handler

	// Configuration du routeur
	r := chi.NewRouter()

	// Middlewares globaux
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	// r.Use(customMiddleware.Logger(loggerInstance)) // TODO: implement logger middleware
	r.Use(chimiddleware.Recoverer)
	// r.Use(customMiddleware.CORS()) // TODO: implement CORS middleware

	// Configuration des routes
	// TODO: Implement routes setup
	routes.SetupRoutes(r, userUsecase, authService, authMiddleware, taskHandler, transactionHandler, accountHandler, budgetHandler, savingGoalHandler, loggerInstance)

	// Configuration du serveur
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Démarrage du serveur en arrière-plan
	go func() {
		loggerInstance.Info("Serveur démarré",
			logger.String("port", fmt.Sprintf(":%d", cfg.Server.Port)),
			logger.String("swagger", "http://localhost:8080/swagger/index.html"),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			loggerInstance.Fatal("Erreur démarrage serveur", logger.Error(err))
		}
	}()

	// Attendre le signal d'arrêt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	loggerInstance.Info("Arrêt du serveur en cours...")

	// Arrêt propre avec timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		loggerInstance.Fatal("Erreur lors de l'arrêt", logger.Error(err))
	}

	loggerInstance.Info("Serveur arrêté proprement")
}
