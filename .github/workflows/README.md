# 🚀 GitHub Actions Workflows - MeshaPlus Backend

Ce dossier contient les workflows GitHub Actions pour automatiser le développement, les tests, la sécurité et le déploiement du backend MeshaPlus.

## 📋 Workflows Disponibles

### 1. 🚀 Deploy Backend (`deploy-backend.yml`)

**Déclencheurs :**
- Push sur `main` ou `develop`
- Pull Request vers `main` ou `develop`
- Modifications dans `backend/` ou le workflow lui-même

**Fonctionnalités :**
- ✅ Tests et validation complète
- 🐳 Construction d'image Docker multi-architecture
- 🔒 Scan de sécurité avec Trivy
- 📈 Tests de performance
- 🚀 Déploiement automatique (staging/production)
- 📢 Notifications

**Jobs :**
1. **🧪 Tests et Validation** - Tests unitaires, linting, couverture
2. **🐳 Build Docker Image** - Construction et push vers GHCR
3. **🚀 Deploy to Production** - Déploiement automatique (main)
4. **🚀 Deploy to Staging** - Déploiement automatique (develop)
5. **🔒 Security Scan** - Scan de vulnérabilités
6. **📈 Performance Test** - Tests de performance
7. **📢 Notifications** - Notifications de succès/échec

### 2. 🚀 Manual Deploy (`manual-deploy.yml`)

**Déclencheurs :**
- Déclenchement manuel (workflow_dispatch)

**Paramètres :**
- `environment` : staging ou production
- `version` : version spécifique (optionnel)
- `force` : forcer le déploiement même si les tests échouent

**Fonctionnalités :**
- 🧪 Tests rapides (sauf si forcé)
- 🐳 Construction et push d'image Docker
- 🚀 Déploiement manuel
- 📢 Notifications

### 3. 🔒 Security & Quality Scan (`security-scan.yml`)

**Déclencheurs :**
- Programmé (tous les lundis à 9h UTC)
- Push sur `main` ou `develop`
- Pull Request vers `main` ou `develop`
- Déclenchement manuel

**Fonctionnalités :**
- 🔍 Analyse de qualité du code
- 🔒 Vérification des dépendances
- 🐳 Sécurité des conteneurs
- 🔐 Détection de secrets
- 📊 Couverture de code
- 📢 Rapport de sécurité

**Jobs :**
1. **🔍 Code Quality** - golangci-lint, govulncheck
2. **🔒 Dependency Security** - Nancy, govulncheck
3. **🐳 Container Security** - Trivy
4. **🔐 Secrets Detection** - TruffleHog
5. **📊 Code Coverage** - Tests avec couverture
6. **📢 Security Summary** - Rapport et notifications

## 🛠️ Configuration

### Variables d'Environnement

```yaml
env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}/backend
  GO_VERSION: '1.21'
  DOCKER_BUILDKIT: 1
```

### Secrets Requis

- `GITHUB_TOKEN` : Automatiquement fourni par GitHub
- `DATABASE_URL` : URL de la base de données (si nécessaire)
- `REDIS_URL` : URL de Redis (si nécessaire)

### Services Utilisés

- **PostgreSQL 15** : Base de données de test
- **Redis 7** : Cache de test
- **GitHub Container Registry** : Stockage des images Docker

## 🚀 Utilisation

### Déploiement Automatique

1. **Staging** : Push sur `develop`
2. **Production** : Push sur `main`

### Déploiement Manuel

1. Aller dans l'onglet "Actions" de GitHub
2. Sélectionner "Manual Deploy Backend"
3. Cliquer sur "Run workflow"
4. Configurer les paramètres :
   - Environment : staging ou production
   - Version : version spécifique (optionnel)
   - Force : forcer le déploiement (optionnel)

### Scan de Sécurité

1. **Automatique** : Tous les lundis à 9h UTC
2. **Manuel** : Via l'interface GitHub Actions
3. **Sur Push/PR** : Automatique sur les branches main/develop

## 📊 Monitoring

### Artifacts Générés

- `coverage-report` : Rapport de couverture HTML
- `security-report` : Rapport de sécurité Markdown
- `trivy-results.sarif` : Résultats de scan Trivy

### Intégrations

- **Codecov** : Couverture de code
- **GitHub Security** : Vulnérabilités
- **GitHub Issues** : Création automatique d'issues pour les problèmes de sécurité

## 🔧 Personnalisation

### Ajout de Notifications

Modifiez les sections `📢 Notifications` dans les workflows :

```yaml
- name: 📢 Send Slack notification
  uses: 8398a7/action-slack@v3
  with:
    status: ${{ job.status }}
    webhook_url: ${{ secrets.SLACK_WEBHOOK }}
```

### Ajout de Déploiement

Modifiez les sections `🚀 Deploy` dans les workflows :

```yaml
- name: 🚀 Deploy to Kubernetes
  run: |
    kubectl set image deployment/meshaplus-backend \
      backend=${{ needs.build.outputs.image-tag }}
```

### Ajout de Tests de Performance

Modifiez les sections `📈 Performance Test` :

```yaml
- name: 🧪 Run k6 tests
  run: |
    k6 run scripts/performance-test.js
```

## 🚨 Dépannage

### Problèmes Courants

1. **Échec de build Docker** : Vérifier le Dockerfile
2. **Échec de tests** : Vérifier les dépendances et la base de données
3. **Échec de déploiement** : Vérifier les permissions et secrets
4. **Échec de scan de sécurité** : Vérifier les vulnérabilités dans les dépendances

### Logs et Debugging

- Consulter les logs dans l'onglet "Actions"
- Télécharger les artifacts pour plus de détails
- Vérifier les issues créées automatiquement

## 📚 Ressources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Buildx](https://docs.docker.com/buildx/)
- [Trivy Security Scanner](https://aquasecurity.github.io/trivy/)
- [golangci-lint](https://golangci-lint.run/)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry) 