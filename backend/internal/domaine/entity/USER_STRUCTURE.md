# Structure de l'Entité User

## Propriétés

L'entité User a été mise à jour avec les propriétés suivantes :

### Champs Principaux

| Champ | Type | Description | Validation |
|-------|------|-------------|------------|
| `id` | `uuid.UUID` | Identifiant unique de l'utilisateur | Généré automatiquement |
| `name` | `string` | Nom complet de l'utilisateur | Requis, 2-100 caractères |
| `email` | `string` | Adresse email unique | Requis, format email valide |
| `avatar` | `string` | URL de l'avatar de l'utilisateur | Optionnel |
| `passwordHash` | `string` | Hash du mot de passe (bcrypt) | Requis, jamais exposé en JSON |
| `createdAt` | `time.Time` | Date de création | Généré automatiquement |
| `updatedAt` | `time.Time` | Date de dernière modification | Mis à jour automatiquement |
| `deletedAt` | `gorm.DeletedAt` | Date de suppression (soft delete) | Géré par GORM |

### Structures de Requête

#### CreateUserRequest
```go
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Avatar   string `json:"avatar,omitempty"`
    Password string `json:"password" validate:"required,min=6"`
}
```

#### UpdateUserRequest
```go
type UpdateUserRequest struct {
    Name     *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
    Email    *string `json:"email,omitempty" validate:"omitempty,email"`
    Avatar   *string `json:"avatar,omitempty"`
    Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
}
```

## Sécurité

- **Mot de passe** : Les mots de passe sont hashés avec bcrypt avant d'être stockés
- **Exposition** : Le hash du mot de passe n'est jamais exposé dans les réponses JSON
- **Validation** : Toutes les données sont validées avant traitement

## Migration

Pour mettre à jour votre base de données avec la nouvelle structure :

```bash
make migrate
```

Ou manuellement :

```bash
./scripts/run_migration.sh
```

## Exemples d'Utilisation

### Créer un utilisateur
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "avatar": "https://example.com/avatar.jpg"
  }'
```

### Mettre à jour un utilisateur
```bash
curl -X PUT http://localhost:8080/api/v1/users/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "avatar": "https://example.com/new-avatar.jpg"
  }'
```

### Récupérer un utilisateur
```bash
curl http://localhost:8080/api/v1/users/{id}
```

### Lister les utilisateurs
```bash
curl "http://localhost:8080/api/v1/users?page=1&page_size=10&search=john"
``` 