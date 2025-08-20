package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger interface pour l'abstraction du logging
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Sync() error
}

// Field représente un champ de log
type Field = zapcore.Field

// logger implémente l'interface Logger avec Zap
type logger struct {
	zap *zap.Logger
}

// New crée une nouvelle instance de logger
func New(level string) Logger {
	config := zap.NewProductionConfig()

	// Configuration du niveau de log
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Configuration de l'encodage (plus lisible en développement)
	config.Encoding = "console"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	zapLogger, err := config.Build()
	if err != nil {
		panic("Impossible d'initialiser le logger: " + err.Error())
	}

	return &logger{
		zap: zapLogger,
	}
}

// Debug log un message de debug
func (l *logger) Debug(msg string, fields ...Field) {
	l.zap.Debug(msg, fields...)
}

// Info log un message d'information
func (l *logger) Info(msg string, fields ...Field) {
	l.zap.Info(msg, fields...)
}

// Warn log un message d'avertissement
func (l *logger) Warn(msg string, fields ...Field) {
	l.zap.Warn(msg, fields...)
}

// Error log un message d'erreur
func (l *logger) Error(msg string, fields ...Field) {
	l.zap.Error(msg, fields...)
}

// Fatal log un message fatal et arrête l'application
func (l *logger) Fatal(msg string, fields ...Field) {
	l.zap.Fatal(msg, fields...)
}

// Sync synchronise les logs
func (l *logger) Sync() error {
	return l.zap.Sync()
}

// Fonctions helper pour créer des champs de log

// String crée un champ string
func String(key, value string) Field {
	return zap.String(key, value)
}

// Int crée un champ int
func Int(key string, value int) Field {
	return zap.Int(key, value)
}

// Int64 crée un champ int64
func Int64(key string, value int64) Field {
	return zap.Int64(key, value)
}

// Float64 crée un champ float64
func Float64(key string, value float64) Field {
	return zap.Float64(key, value)
}

// Bool crée un champ bool
func Bool(key string, value bool) Field {
	return zap.Bool(key, value)
}

// Error crée un champ error
func Error(err error) Field {
	return zap.Error(err)
}

// Duration crée un champ duration
func Duration(key string, value interface{}) Field {
	return zap.Any(key, value)
}

// Any crée un champ de type any
func Any(key string, value interface{}) Field {
	return zap.Any(key, value)
}
