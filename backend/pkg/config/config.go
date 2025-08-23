package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Redis       RedisConfig       `mapstructure:"redis"`
	Logger      LoggerConfig      `mapstructure:"logger"`
	Storage     StorageConfig     `mapstructure:"storage"`
	ExternalAPI ExternalAPIConfig `mapstructure:"external_api"`
	JWT         JWTConfig         `mapstructure:"jwt"`
	AI          AIConfig          `mapstructure:"ai"`
}

type ServerConfig struct {
	Port         int `mapstructure:"port"`
	ReadTimeout  int `mapstructure:"read_timeout"`
	WriteTimeout int `mapstructure:"write_timeout"`
	IdleTimeout  int `mapstructure:"idle_timeout"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type StorageConfig struct {
	Type      string   `mapstructure:"type"` // local, s3
	LocalPath string   `mapstructure:"local_path"`
	S3Config  S3Config `mapstructure:"s3"`
}

type S3Config struct {
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type ExternalAPIConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Timeout int    `mapstructure:"timeout"`
}

type JWTConfig struct {
	SecretKey              string `mapstructure:"secret_key"`
	ExpirationHours        int    `mapstructure:"expiration_hours"`
	RefreshExpirationHours int    `mapstructure:"refresh_expiration_hours"`
}

type AIConfig struct {
	Model string `mapstructure:"model"`
}

func Load() (*Config, error) {
	// Charger le fichier .env si disponible
	if err := godotenv.Load(); err != nil {
		// Le fichier .env n'est pas obligatoire
	}

	// Déterminer l'environnement
	env := getEnv("APP_ENV", "dev")

	// Configuration de Viper
	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../../configs") // Pour les tests

	// Variables d'environnement
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	// Valeurs par défaut
	setDefaults()

	// Lire le fichier de configuration
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("erreur lecture config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("erreur unmarshalling config: %w", err)
	}

	// Override avec les variables d'environnement
	overrideWithEnv(&config)

	return &config, nil
}

func setDefaults() {
	// Server
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 15)
	viper.SetDefault("server.write_timeout", 15)
	viper.SetDefault("server.idle_timeout", 60)

	// Database
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "myapp")
	viper.SetDefault("database.sslmode", "disable")

	// Redis
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// Logger
	viper.SetDefault("logger.level", "info")

	// Storage
	viper.SetDefault("storage.type", "local")
	viper.SetDefault("storage.local_path", "./uploads")

	// External API
	viper.SetDefault("external_api.base_url", "https://jsonplaceholder.typicode.com")
	viper.SetDefault("external_api.timeout", 30)

	// JWT
	viper.SetDefault("jwt.secret_key", "your-secret-key-change-in-production")
	viper.SetDefault("jwt.expiration_hours", 24)
	viper.SetDefault("jwt.refresh_expiration_hours", 168) // 7 jours
}

func overrideWithEnv(config *Config) {
	// Server
	if port := getEnv("SERVER_PORT", ""); port != "" {
		if portInt, err := strconv.Atoi(port); err == nil {
			config.Server.Port = portInt
		}
	}

	// Database
	if host := getEnv("DB_HOST", ""); host != "" {
		config.Database.Host = host
	}
	if user := getEnv("DB_USER", ""); user != "" {
		config.Database.User = user
	}
	if password := getEnv("DB_PASSWORD", ""); password != "" {
		config.Database.Password = password
	}
	if dbname := getEnv("DB_NAME", ""); dbname != "" {
		config.Database.DBName = dbname
	}

	// Redis
	if host := getEnv("REDIS_HOST", ""); host != "" {
		config.Redis.Host = host
	}
	if password := getEnv("REDIS_PASSWORD", ""); password != "" {
		config.Redis.Password = password
	}

	// JWT
	if secretKey := getEnv("JWT_SECRET_KEY", ""); secretKey != "" {
		config.JWT.SecretKey = secretKey
	}
	if expirationHours := getEnv("JWT_EXPIRATION_HOURS", ""); expirationHours != "" {
		if hours, err := strconv.Atoi(expirationHours); err == nil {
			config.JWT.ExpirationHours = hours
		}
	}
	if refreshExpirationHours := getEnv("JWT_REFRESH_EXPIRATION_HOURS", ""); refreshExpirationHours != "" {
		if hours, err := strconv.Atoi(refreshExpirationHours); err == nil {
			config.JWT.RefreshExpirationHours = hours
		}
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
