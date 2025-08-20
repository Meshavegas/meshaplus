package database

import (
	"backend/pkg/config"
	"fmt"

	"backend/pkg/logger"

	"github.com/go-pg/pg/v10"
)

// NewPostgresConnection crée une nouvelle connexion PostgreSQL avec go-pg
func NewPostgresConnection(cfg config.DatabaseConfig, loggerInstance logger.Logger) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.DBName,
		PoolSize: 20, // Taille du pool de connexions
	})

	// Tester la connexion
	if err := db.Ping(db.Context()); err != nil {
		loggerInstance.Error("Erreur connexion PostgreSQL", logger.Error(err))
		return nil, fmt.Errorf("erreur connexion PostgreSQL: %w", err)
	}

	loggerInstance.Info("Connexion PostgreSQL établie avec succès",
		logger.String("host", cfg.Host),
		logger.Int("port", cfg.Port),
		logger.String("database", cfg.DBName),
	)

	return db, nil
}

// ClosePostgresConnection ferme la connexion PostgreSQL
func ClosePostgresConnection(db *pg.DB, loggerInstance logger.Logger) {
	if db != nil {
		if err := db.Close(); err != nil {
			loggerInstance.Error("Erreur fermeture connexion PostgreSQL", logger.Error(err))
		} else {
			loggerInstance.Info("Connexion PostgreSQL fermée avec succès")
		}
	}
}
