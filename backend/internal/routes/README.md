# Structure des Routes

Ce dossier contient l'organisation des routes de l'application suivant une approche modulaire par ressource.

## Organisation

### `routes.go`
Fichier principal qui orchestre toutes les routes de l'application. Il configure les routes API v1 et la documentation Swagger.

### `user_routes.go`
Routes spécifiques aux utilisateurs :
- `POST /api/v1/users` - Créer un utilisateur
- `GET /api/v1/users` - Lister tous les utilisateurs
- `GET /api/v1/users/{id}` - Récupérer un utilisateur par ID
- `PUT /api/v1/users/{id}` - Mettre à jour un utilisateur
- `DELETE /api/v1/users/{id}` - Supprimer un utilisateur

### `file_routes.go`
Routes pour la gestion des fichiers (à implémenter) :
- `POST /api/v1/files/upload` - Uploader un fichier
- `GET /api/v1/files/{id}` - Récupérer un fichier
- `DELETE /api/v1/files/{id}` - Supprimer un fichier

### `health_routes.go`
Routes pour la santé de l'application (à implémenter) :
- `GET /api/v1/health` - Vérifier la santé de l'application

### `swagger_routes.go`
Routes pour la documentation Swagger :
- `GET /swagger/*` - Interface Swagger UI

## Utilisation

Dans `main.go`, utilisez simplement :

```go
// Configuration des routes
routes.SetupRoutes(r, userUsecase, loggerInstance)
```

## Ajout de nouvelles routes

Pour ajouter une nouvelle ressource :

1. Créez un nouveau fichier `{resource}_routes.go`
2. Implémentez la fonction `Setup{Resource}Routes`
3. Ajoutez l'appel dans `routes.go`

Exemple :
```go
// Dans routes.go
func SetupRoutes(r chi.Router, userUsecase *usecase.UserUsecase, logger logger.Logger) {
    r.Route("/api/v1", func(r chi.Router) {
        SetupUserRoutes(r, userUsecase, logger)
        SetupProductRoutes(r, productUsecase, logger) // Nouvelle ressource
    })
}
``` 