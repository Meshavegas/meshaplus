# 🚀 Déploiement Binaire MeshaPlus Backend

Ce guide explique comment déployer MeshaPlus Backend directement sur un VPS sans Docker, en utilisant les services locaux (PostgreSQL, Redis).

## 📋 Prérequis

### Sur votre machine locale
- Go 1.21+
- SSH/SCP
- Accès SSH au VPS

### Sur le VPS
- Ubuntu/Debian
- Accès root ou sudo
- Connexion Internet

## 🔧 Installation des Services Locaux

### 1. Installation automatique (recommandé)

```bash
# Depuis votre machine locale
make setup-vps-services VPS_HOST=your-vps.com VPS_USER=root
```

### 2. Installation manuelle

Si vous préférez installer manuellement :

```bash
# Se connecter au VPS
ssh root@your-vps.com

# Exécuter le script d'installation
chmod +x /tmp/setup-local-services.sh
./setup-local-services.sh
```

### Services installés

- **PostgreSQL 15** : Base de données principale
  - Port : 5432
  - Base : `meshaplus`
  - Utilisateur : `meshaplus`
  - Mot de passe : `meshaplus_password`

- **Redis 7** : Cache et sessions
  - Port : 6379
  - Pas de mot de passe (à configurer en production)

- **Nginx** : Reverse proxy (optionnel)
  - Port : 80/443
  - Configuration automatique

## 🚀 Déploiement de l'Application

### Déploiement automatique

```bash
# Depuis votre machine locale
make deploy-binary VPS_HOST=your-vps.com VPS_USER=root
```

### Déploiement manuel

```bash
# 1. Compiler le binaire
make build-prod

# 2. Transférer et déployer
./scripts/build-and-deploy.sh -h your-vps.com -u root
```

## 📁 Structure de Déploiement

```
/opt/meshaplus/
├── bin/
│   └── api                    # Binaire de l'application
├── configs/
│   └── config.prod.yaml      # Configuration de production
├── uploads/                   # Fichiers uploadés
└── backups/                   # Sauvegardes des versions

/var/log/meshaplus/
└── app.log                    # Logs de l'application
```

## 🔧 Configuration

### Fichier de configuration

Le fichier `configs/config.prod.yaml` contient la configuration de production :

```yaml
server:
  port: 8080
  host: "0.0.0.0"

database:
  host: "localhost"
  port: 5432
  name: "meshaplus"
  user: "meshaplus"
  password: "meshaplus_password"

redis:
  host: "localhost"
  port: 6379
```

### Service systemd

L'application est gérée par systemd :

```bash
# Statut du service
sudo systemctl status meshaplus-backend

# Démarrer le service
sudo systemctl start meshaplus-backend

# Arrêter le service
sudo systemctl stop meshaplus-backend

# Redémarrer le service
sudo systemctl restart meshaplus-backend

# Voir les logs
sudo journalctl -u meshaplus-backend -f
```

## 🌐 Accès à l'Application

### URLs d'accès

- **API** : `http://your-vps.com:8080`
- **Health Check** : `http://your-vps.com:8080/api/v1/health`
- **Swagger** : `http://your-vps.com:8080/swagger/index.html`

### Avec Nginx (recommandé)

Si Nginx est installé, l'application sera accessible sur le port 80 :

- **API** : `http://your-vps.com`
- **Health Check** : `http://your-vps.com/api/v1/health`
- **Swagger** : `http://your-vps.com/swagger/index.html`

## 🔄 Mise à Jour

### Mise à jour automatique

```bash
# Le script de déploiement gère automatiquement :
# - Sauvegarde de l'ancienne version
# - Arrêt du service
# - Déploiement de la nouvelle version
# - Redémarrage du service
# - Vérification de santé

make deploy-binary VPS_HOST=your-vps.com VPS_USER=root
```

### Rollback manuel

```bash
# Se connecter au VPS
ssh root@your-vps.com

# Arrêter le service
sudo systemctl stop meshaplus-backend

# Restaurer la sauvegarde
sudo cp /opt/meshaplus/bin/api.backup.YYYYMMDD-HHMMSS /opt/meshaplus/bin/api

# Redémarrer le service
sudo systemctl start meshaplus-backend
```

## 🔍 Monitoring et Logs

### Logs de l'application

```bash
# Logs systemd
sudo journalctl -u meshaplus-backend -f

# Logs fichier
sudo tail -f /var/log/meshaplus/app.log
```

### Monitoring des services

```bash
# Statut PostgreSQL
sudo systemctl status postgresql

# Statut Redis
sudo systemctl status redis-server

# Statut Nginx
sudo systemctl status nginx
```

### Logs des services

```bash
# Logs PostgreSQL
sudo tail -f /var/log/postgresql/postgresql-*.log

# Logs Redis
sudo tail -f /var/log/redis/redis-server.log

# Logs Nginx
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

## 🔒 Sécurité

### Recommandations de production

1. **Changer les mots de passe par défaut**
   ```bash
   # PostgreSQL
   sudo -u postgres psql
   ALTER USER meshaplus WITH PASSWORD 'nouveau_mot_de_passe_fort';
   ```

2. **Configurer Redis avec un mot de passe**
   ```bash
   # Éditer /etc/redis/redis.conf
   sudo nano /etc/redis/redis.conf
   # Ajouter : requirepass votre_mot_de_passe_redis
   sudo systemctl restart redis-server
   ```

3. **Configurer le firewall**
   ```bash
   # Le script configure automatiquement UFW
   # Vérifier : sudo ufw status
   ```

4. **SSL/TLS avec Let's Encrypt**
   ```bash
   # Installer Certbot
   sudo apt install certbot python3-certbot-nginx

   # Obtenir un certificat
   sudo certbot --nginx -d your-domain.com
   ```

## 🛠️ Dépannage

### Problèmes courants

1. **Service ne démarre pas**
   ```bash
   # Vérifier les logs
   sudo journalctl -u meshaplus-backend -n 50

   # Vérifier les permissions
   sudo chown -R meshaplus:meshaplus /opt/meshaplus
   ```

2. **Base de données inaccessible**
   ```bash
   # Vérifier PostgreSQL
   sudo systemctl status postgresql
   pg_isready -h localhost -p 5432
   ```

3. **Redis inaccessible**
   ```bash
   # Vérifier Redis
   sudo systemctl status redis-server
   redis-cli ping
   ```

4. **Port déjà utilisé**
   ```bash
   # Vérifier les ports utilisés
   sudo netstat -tlnp | grep :8080
   ```

## 📚 Commandes Utiles

### Makefile

```bash
# Build pour production
make build-prod

# Déploiement complet
make deploy-binary VPS_HOST=your-vps.com VPS_USER=root

# Installation des services
make setup-vps-services VPS_HOST=your-vps.com VPS_USER=root
```

### Scripts

```bash
# Build et déploiement
./scripts/build-and-deploy.sh -h your-vps.com -u root

# Installation des services
./scripts/setup-local-services.sh

# Déploiement du binaire
./scripts/deploy-binary-vps.sh
```

## 🔄 Migration depuis Docker

Si vous migrez depuis un déploiement Docker :

1. **Arrêter les conteneurs Docker**
   ```bash
   docker-compose down
   ```

2. **Sauvegarder les données**
   ```bash
   # Sauvegarder PostgreSQL
   docker exec meshaplus-postgres pg_dump -U postgres meshaplus > backup.sql
   ```

3. **Installer les services locaux**
   ```bash
   make setup-vps-services VPS_HOST=your-vps.com VPS_USER=root
   ```

4. **Restaurer les données**
   ```bash
   # Restaurer PostgreSQL
   psql -h localhost -U meshaplus -d meshaplus < backup.sql
   ```

5. **Déployer l'application**
   ```bash
   make deploy-binary VPS_HOST=your-vps.com VPS_USER=root
   ``` 