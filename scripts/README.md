# 🔧 Scripts SSH pour MeshaPlus

Ce dossier contient des scripts pour configurer et diagnostiquer les connexions SSH nécessaires au déploiement de MeshaPlus.

## 📁 Scripts Disponibles

### 1. `setup-ssh.sh` - Configuration SSH Automatique

**Usage:**
```bash
./setup-ssh.sh <host> <username> [key_name]
```

**Exemple:**
```bash
./setup-ssh.sh 192.168.1.100 ubuntu meshaplus_key
```

**Fonctionnalités:**
- Génère une nouvelle paire de clés SSH
- Configure les permissions appropriées
- Affiche les informations pour GitHub Secrets
- Crée une configuration SSH locale
- Teste la connexion

### 2. `ssh-diagnostic.sh` - Diagnostic SSH

**Usage:**
```bash
./ssh-diagnostic.sh [host] [username] [ssh_key_path]
```

**Exemple:**
```bash
./ssh-diagnostic.sh 192.168.1.100 ubuntu ~/.ssh/meshaplus_key
```

**Fonctionnalités:**
- Vérifie l'existence et la validité de la clé SSH
- Teste les permissions
- Vérifie la connectivité réseau
- Teste la connexion SSH
- Vérifie l'environnement sur le VPS

## 🚀 Démarrage Rapide

### Étape 1: Configuration SSH

```bash
# Rendre les scripts exécutables
chmod +x scripts/*.sh

# Configurer SSH pour votre VPS
./scripts/setup-ssh.sh votre-vps-ip votre-utilisateur
```

### Étape 2: Configuration GitHub Secrets

Le script `setup-ssh.sh` affichera les informations nécessaires. Ajoutez ces secrets dans GitHub :

1. Allez dans votre repository GitHub
2. Settings > Secrets and variables > Actions
3. Ajoutez les secrets suivants :

```yaml
VPS_HOST: "votre-vps-ip-ou-domaine"
VPS_USERNAME: "votre-utilisateur"
VPS_SSH_KEY: "-----BEGIN OPENSSH PRIVATE KEY-----\n...\n-----END OPENSSH PRIVATE KEY-----"
VPS_DEPLOY_PATH: "/opt/meshaplus"  # Optionnel
```

### Étape 3: Test de Diagnostic

```bash
# Tester la configuration
./scripts/ssh-diagnostic.sh votre-vps-ip votre-utilisateur ~/.ssh/meshaplus_key
```

## 🔍 Résolution de Problèmes

### Erreur: "Permission denied (publickey)"

1. **Vérifiez la clé SSH:**
   ```bash
   ./scripts/ssh-diagnostic.sh votre-vps-ip votre-utilisateur ~/.ssh/votre_cle
   ```

2. **Ajoutez la clé publique au VPS:**
   ```bash
   # Copiez la clé publique
   cat ~/.ssh/votre_cle.pub
   
   # Sur le VPS, ajoutez-la à authorized_keys
   echo "votre_clé_publique" >> ~/.ssh/authorized_keys
   chmod 600 ~/.ssh/authorized_keys
   ```

3. **Vérifiez les permissions:**
   ```bash
   chmod 600 ~/.ssh/votre_cle
   chmod 700 ~/.ssh
   ```

### Erreur: "Host key verification failed"

```bash
# Ajoutez l'option StrictHostKeyChecking=no
ssh -o StrictHostKeyChecking=no -i ~/.ssh/votre_cle utilisateur@vps-ip
```

### Erreur: "Connection timeout"

1. Vérifiez que le VPS est accessible
2. Vérifiez que le port SSH (22) est ouvert
3. Vérifiez votre pare-feu

## 📋 Checklist de Vérification

- [ ] Clé SSH privée générée
- [ ] Clé publique ajoutée au VPS
- [ ] Permissions correctes (600 pour clé privée, 644 pour clé publique)
- [ ] Utilisateur SSH existe sur le VPS
- [ ] Service SSH actif sur le VPS
- [ ] Secrets GitHub configurés
- [ ] Test de connexion réussi

## 🛠️ Commandes Utiles

### Génération manuelle de clés SSH

```bash
# Générer une nouvelle clé
ssh-keygen -t ed25519 -f ~/.ssh/meshaplus_key -N ""

# Afficher la clé publique
cat ~/.ssh/meshaplus_key.pub

# Afficher la clé privée (pour GitHub Secrets)
cat ~/.ssh/meshaplus_key
```

### Test de connexion

```bash
# Test basique
ssh -i ~/.ssh/meshaplus_key utilisateur@vps-ip

# Test avec debug
ssh -v -i ~/.ssh/meshaplus_key utilisateur@vps-ip

# Test avec timeout
ssh -o ConnectTimeout=30 -i ~/.ssh/meshaplus_key utilisateur@vps-ip
```

### Vérification sur le VPS

```bash
# Vérifier les logs SSH
sudo journalctl -u ssh -f

# Vérifier la configuration SSH
sudo sshd -T | grep -E "(pubkey|authorized)"

# Vérifier les permissions
ls -la ~/.ssh/
```

## 📚 Documentation Complète

Pour plus de détails sur la résolution de problèmes SSH, consultez :
- [Guide de Dépannage SSH](../docs/SSH_TROUBLESHOOTING.md)
- [Documentation SSH officielle](https://www.openssh.com/manual.html)

## 🆘 Support

Si vous rencontrez des problèmes :

1. Utilisez le script de diagnostic
2. Consultez le guide de dépannage
3. Vérifiez les logs GitHub Actions
4. Contactez l'équipe de développement

---

**Note:** Ces scripts sont conçus pour faciliter la configuration SSH pour le déploiement de MeshaPlus. Assurez-vous de suivre les bonnes pratiques de sécurité. 