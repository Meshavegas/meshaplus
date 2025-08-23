#!/bin/bash

# Script pour résoudre les problèmes de rate limit Docker
# Usage: ./fix-docker-rate-limit.sh [host] [username] [key_path]

set -e

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Fonction pour afficher les messages
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Paramètres
HOST=${1:-""}
USERNAME=${2:-""}
KEY_PATH=${3:-"~/.ssh/meshaplus_deploy_key"}

# Vérification des paramètres
if [ -z "$HOST" ] || [ -z "$USERNAME" ]; then
    log_error "Usage: $0 <host> <username> [key_path]"
    log_error "Example: $0 192.168.1.100 ubuntu ~/.ssh/meshaplus_deploy_key"
    exit 1
fi

KEY_PATH=$(eval echo "$KEY_PATH")

log_info "🔧 Résolution des problèmes de rate limit Docker"
log_info "Host: $HOST"
log_info "Username: $USERNAME"
log_info "SSH Key: $KEY_PATH"
echo

# 1. Vérifier que la clé SSH existe
if [ ! -f "$KEY_PATH" ]; then
    log_error "Clé SSH introuvable: $KEY_PATH"
    exit 1
fi

log_success "Clé SSH trouvée: $KEY_PATH"
echo

# 2. Tester la connexion SSH
log_info "1. Test de connexion SSH..."
if ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" "echo 'Connexion SSH réussie'"; then
    log_success "Connexion SSH réussie"
else
    log_error "Échec de la connexion SSH"
    exit 1
fi
echo

# 3. Configurer Docker pour éviter les rate limits
log_info "2. Configuration de Docker pour éviter les rate limits..."
ssh -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" << 'EOF'
    echo "🔧 Configuration de Docker..."
    
    # Créer le répertoire de configuration Docker
    sudo mkdir -p /etc/docker
    
    # Créer ou modifier le fichier daemon.json avec des mirrors alternatifs
    sudo tee /etc/docker/daemon.json > /dev/null << 'DOCKER_CONFIG'
{
  "registry-mirrors": [
    "https://mirror.gcr.io",
    "https://registry-1.docker.io",
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com",
    "https://mirror.baidubce.com"
  ],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  },
  "max-concurrent-downloads": 3,
  "max-concurrent-uploads": 3
}
DOCKER_CONFIG
    
    # Redémarrer Docker
    sudo systemctl restart docker
    
    echo "✅ Docker configuré avec des mirrors alternatifs"
EOF
echo

# 4. Pré-télécharger les images nécessaires
log_info "3. Pré-téléchargement des images Docker..."
ssh -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" << 'EOF'
    echo "📦 Pré-téléchargement des images de base..."
    
    # Attendre que Docker soit prêt
    sleep 5
    
    # Télécharger les images avec retry
    for image in "golang:1.23-alpine" "alpine:latest" "postgres:15" "redis:7-alpine"; do
        echo "📥 Téléchargement de $image..."
        for attempt in {1..5}; do
            if docker pull "$image"; then
                echo "✅ $image téléchargé avec succès"
                break
            else
                echo "⚠️ Tentative $attempt/5 échouée pour $image"
                if [ $attempt -lt 5 ]; then
                    echo "⏳ Attente avant nouvelle tentative..."
                    sleep 30
                fi
            fi
        done
    done
    
    echo "✅ Pré-téléchargement terminé"
EOF
echo

# 5. Créer un Dockerfile optimisé
log_info "4. Création d'un Dockerfile optimisé..."
ssh -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" << 'EOF'
    echo "📝 Création d'un Dockerfile optimisé..."
    
    # Créer un Dockerfile alternatif qui utilise des images déjà téléchargées
    cat > /tmp/Dockerfile.optimized << 'DOCKERFILE'
# Build stage
FROM golang:1.23-alpine AS builder

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o bin/api \
    ./cmd/api

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/bin/api .

# Copy configs and scripts
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/scripts ./scripts

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Run the application
CMD ["./api"]
DOCKERFILE
    
    echo "✅ Dockerfile optimisé créé"
EOF
echo

# 6. Tester la construction d'image
log_info "5. Test de construction d'image..."
ssh -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" << 'EOF'
    echo "🧪 Test de construction d'image..."
    
    # Aller dans le répertoire de déploiement
    cd /home/ubuntu/meshaplus/backend
    
    # Copier le Dockerfile optimisé
    cp /tmp/Dockerfile.optimized Dockerfile
    
    # Tester la construction
    if docker buildx build --platform linux/amd64 -t meshaplus-backend:test . --load; then
        echo "✅ Construction d'image réussie"
        
        # Nettoyer l'image de test
        docker rmi meshaplus-backend:test
        
        echo "✅ Test de construction terminé avec succès"
    else
        echo "❌ Échec de la construction d'image"
        exit 1
    fi
EOF
echo

# 7. Vérification finale
log_info "6. Vérification finale..."
ssh -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" << 'EOF'
    echo "🔍 Vérification finale..."
    
    echo "🐳 Images Docker disponibles:"
    docker images
    
    echo "🐳 Configuration Docker:"
    docker info | grep -E "(Registry Mirrors|Docker Root Dir)"
    
    echo "📊 Espace disque:"
    df -h /
    
    echo "✅ Configuration terminée avec succès"
EOF
echo

log_success "🎉 Problèmes de rate limit Docker résolus!"
log_info "Le VPS est maintenant prêt pour le déploiement MeshaPlus"
log_info "Vous pouvez relancer le workflow GitHub Actions" 