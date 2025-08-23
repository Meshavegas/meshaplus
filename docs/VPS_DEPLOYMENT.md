# 🚀 Déploiement VPS - MeshaPlus API

Ce guide explique comment déployer l'API MeshaPlus sur un VPS avec un workflow GitHub Actions simple et robuste, inspiré de l'article de [Lanre](https://lanre.wtf/blog/2024/05/11/go-vps).

## 📋 Prérequis

- Un VPS avec accès SSH
- Une clé SSH configurée
- Go 1.23+ installé sur le VPS
- PostgreSQL installé et configuré
- Accès administrateur au repository GitHub

## 🔧 Configuration initiale

### 1. Configuration des secrets GitHub

Exécutez le script de configuration des secrets :

```bash
./scripts/setup-github-secrets.sh
```

Ou configurez manuellement ces secrets dans GitHub (Settings > Secrets and variables > Actions) :

- `VPS_HOST` : Adresse IP ou nom de domaine de votre VPS
- `VPS_USER` : Nom d'utilisateur SSH (généralement `root`)
- `VPS_SSH_KEY` : Contenu de votre clé SSH privée
- `VPS_PORT` : Port SSH (optionnel, défaut: 22)

### 2. Configuration du service systemd sur le VPS

Connectez-vous à votre VPS et exécutez :

```bash
# Copier le fichier de service
sudo cp scripts/meshaplus-api.service /etc/systemd/system/

# Créer les répertoires nécessaires
sudo mkdir -p /opt/meshaplus-api/bin
sudo mkdir -p /opt/meshaplus-api/backups
sudo mkdir -p /opt/meshaplus-api/logs

# Recharger systemd
sudo systemctl daemon-reload

# Activer le service
sudo systemctl enable meshaplus-api
```

## 🚀 Déploiement automatique

### Déclenchement

Le déploiement se déclenche automatiquement quand :
- Un push est effectué sur la branche `main`
- Les fichiers modifiés sont dans `backend/` ou le workflow lui-même
- Un déploiement manuel est déclenché via GitHub Actions

### Processus de déploiement

1. **Build** : Compilation de l'application Go
2. **Upload** : Transfert du binaire vers le VPS
3. **Backup** : Sauvegarde de l'ancienne version
4. **Déploiement** : Remplacement du binaire et redémarrage du service
5. **Test** : Vérification que le service fonctionne
6. **Rollback** : Retour à l'ancienne version si échec

## 🧪 Test du déploiement

Testez votre déploiement avec :

```bash
./scripts/test-deployment.sh [VPS_HOST] [VPS_USER]
```

Exemple :
```bash
./scripts/test-deployment.sh 192.168.1.100 root
```

## 📁 Structure des fichiers

```
/opt/meshaplus-api/
├── bin/
│   └── api                    # Binaire de l'application
├── backups/
│   ├── 20241201_143022_api   # Sauvegardes automatiques
│   └── ...
└── logs/                      # Logs de l'application
```

## 🔧 Commandes utiles

### Service systemd

```bash
# Démarrer le service
sudo systemctl start meshaplus-api

# Arrêter le service
sudo systemctl stop meshaplus-api

# Redémarrer le service
sudo systemctl restart meshaplus-api

# Vérifier le statut
sudo systemctl status meshaplus-api

# Voir les logs en temps réel
sudo journalctl -u meshaplus-api -f

# Voir les logs récents
sudo journalctl -u meshaplus-api --no-pager -n 50
```

### Déploiement manuel

```bash
# Déclencher un déploiement manuel
gh workflow run vps-deploy.yml

# Voir l'historique des déploiements
gh run list --workflow=vps-deploy.yml
```

## 🔍 Monitoring et logs

### Logs du service

```bash
# Logs en temps réel
sudo journalctl -u meshaplus-api -f

# Logs des dernières 24h
sudo journalctl -u meshaplus-api --since "24 hours ago"

# Logs d'erreur uniquement
sudo journalctl -u meshaplus-api -p err
```

### Vérification de l'API

```bash
# Health check
curl http://localhost:8080/health

# Documentation Swagger
curl http://localhost:8080/swagger/index.html

# Test de l'API
curl http://localhost:8080/api/v1
```

## 🛠️ Dépannage

### Service ne démarre pas

1. Vérifiez les logs :
   ```bash
   sudo journalctl -u meshaplus-api --no-pager -n 50
   ```

2. Vérifiez la configuration :
   ```bash
   sudo systemctl status meshaplus-api
   ```

3. Testez le binaire manuellement :
   ```bash
   /opt/meshaplus-api/bin/api
   ```

### Rollback manuel

Si le déploiement échoue et le rollback automatique ne fonctionne pas :

```bash
# Lister les sauvegardes
ls -la /opt/meshaplus-api/backups/

# Restaurer une version précédente
sudo systemctl stop meshaplus-api
sudo cp /opt/meshaplus-api/backups/YYYYMMDD_HHMMSS_api /opt/meshaplus-api/bin/api
sudo systemctl start meshaplus-api
```

### Problèmes de permissions

```bash
# Vérifier les permissions
ls -la /opt/meshaplus-api/bin/api

# Corriger les permissions si nécessaire
sudo chmod +x /opt/meshaplus-api/bin/api
sudo chown root:root /opt/meshaplus-api/bin/api
```

## 🔒 Sécurité

- Le service s'exécute avec des privilèges limités
- Les sauvegardes sont automatiquement nettoyées (gardent les 5 dernières)
- Le service redémarre automatiquement en cas de crash
- Les logs sont centralisés via systemd

## 📚 Ressources

- [Article de référence - Lanre](https://lanre.wtf/blog/2024/05/11/go-vps)
- [Documentation systemd](https://systemd.io/)
- [GitHub Actions SSH](https://github.com/appleboy/ssh-action)
- [GitHub Actions SCP](https://github.com/appleboy/scp-action) 