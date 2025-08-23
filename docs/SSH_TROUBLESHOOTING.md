# 🔧 Guide de Dépannage SSH - MeshaPlus

Ce guide vous aide à résoudre les problèmes de connexion SSH lors du déploiement de MeshaPlus.

## 🚨 Erreur: "Permission denied (publickey)"

### Causes possibles et solutions

#### 1. **Clé SSH manquante ou invalide**

**Symptômes:**
- `Permission denied (publickey)`
- `SSH key is empty or invalid`

**Solutions:**
```bash
# Vérifier que la clé SSH existe dans GitHub Secrets
# Aller dans Settings > Secrets and variables > Actions
# Vérifier que VPS_SSH_KEY est défini et contient la clé privée complète

# Format attendu de la clé SSH:
-----BEGIN OPENSSH PRIVATE KEY-----
[contenu de la clé]
-----END OPENSSH PRIVATE KEY-----
```

#### 2. **Clé publique non ajoutée au VPS**

**Symptômes:**
- Connexion SSH échoue même avec une clé valide

**Solutions:**
```bash
# Sur votre machine locale, générer la clé publique
ssh-keygen -y -f ~/.ssh/vps_key > ~/.ssh/vps_key.pub

# Copier la clé publique sur le VPS
cat ~/.ssh/vps_key.pub | ssh user@vps-host "mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys"

# Ou manuellement:
# 1. Copier le contenu de ~/.ssh/vps_key.pub
# 2. Se connecter au VPS
# 3. Ajouter à ~/.ssh/authorized_keys
echo "votre_clé_publique" >> ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys
```

#### 3. **Permissions incorrectes**

**Symptômes:**
- SSH refuse la clé même si elle est valide

**Solutions:**
```bash
# Sur le VPS, vérifier les permissions
chmod 700 ~/.ssh
chmod 600 ~/.ssh/authorized_keys
chown -R $USER:$USER ~/.ssh

# Sur votre machine locale
chmod 600 ~/.ssh/vps_key
```

#### 4. **Configuration SSH incorrecte**

**Symptômes:**
- Connexion SSH échoue avec des erreurs de configuration

**Solutions:**
```bash
# Vérifier la configuration SSH sur le VPS
sudo nano /etc/ssh/sshd_config

# Assurez-vous que ces lignes sont présentes et non commentées:
PubkeyAuthentication yes
AuthorizedKeysFile .ssh/authorized_keys
PasswordAuthentication no  # Optionnel, pour plus de sécurité

# Redémarrer le service SSH
sudo systemctl restart sshd
```

#### 5. **Utilisateur inexistant ou incorrect**

**Symptômes:**
- `Permission denied` même avec une clé valide

**Solutions:**
```bash
# Vérifier que l'utilisateur existe sur le VPS
sudo useradd -m -s /bin/bash username
sudo usermod -aG sudo username

# Ou vérifier un utilisateur existant
id username
```

## 🔍 Diagnostic Automatique

Utilisez notre script de diagnostic pour identifier automatiquement les problèmes :

```bash
# Rendre le script exécutable
chmod +x scripts/ssh-diagnostic.sh

# Exécuter le diagnostic
./scripts/ssh-diagnostic.sh your-vps-host your-username ~/.ssh/vps_key
```

## 📋 Checklist de Vérification

### Avant le déploiement

- [ ] Clé SSH privée générée et sauvegardée
- [ ] Clé publique ajoutée à `~/.ssh/authorized_keys` sur le VPS
- [ ] Permissions correctes sur les fichiers SSH
- [ ] Utilisateur SSH existe sur le VPS
- [ ] Service SSH actif sur le VPS
- [ ] Secrets GitHub configurés correctement

### Secrets GitHub Requis

```yaml
VPS_HOST: "votre-vps-ip-ou-domaine"
VPS_USERNAME: "votre-utilisateur"
VPS_SSH_KEY: "-----BEGIN OPENSSH PRIVATE KEY-----\n...\n-----END OPENSSH PRIVATE KEY-----"
VPS_DEPLOY_PATH: "/opt/meshaplus"  # Optionnel, défaut: /opt/meshaplus
```

## 🛠️ Commandes Utiles

### Génération d'une nouvelle clé SSH

```bash
# Générer une nouvelle paire de clés
ssh-keygen -t ed25519 -f ~/.ssh/vps_key -C "meshaplus-deployment"

# Afficher la clé publique
cat ~/.ssh/vps_key.pub

# Afficher la clé privée (pour GitHub Secrets)
cat ~/.ssh/vps_key
```

### Test de connexion manuel

```bash
# Test basique
ssh -i ~/.ssh/vps_key username@vps-host

# Test avec debug
ssh -v -i ~/.ssh/vps_key username@vps-host

# Test avec timeout
ssh -o ConnectTimeout=30 -i ~/.ssh/vps_key username@vps-host
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

## 🚀 Résolution Rapide

Si vous avez besoin d'une solution rapide :

1. **Générer une nouvelle clé SSH**
2. **Ajouter la clé publique au VPS**
3. **Mettre à jour le secret GitHub**
4. **Relancer le workflow**

```bash
# 1. Générer une nouvelle clé
ssh-keygen -t ed25519 -f ~/.ssh/meshaplus_key -N ""

# 2. Ajouter au VPS
ssh-copy-id -i ~/.ssh/meshaplus_key.pub username@vps-host

# 3. Copier la clé privée dans GitHub Secrets
cat ~/.ssh/meshaplus_key
```

## 📞 Support

Si les problèmes persistent :

1. Vérifiez les logs du workflow GitHub Actions
2. Utilisez le script de diagnostic
3. Consultez les logs SSH sur le VPS
4. Contactez l'équipe de développement

---

**Note:** Ce guide couvre les problèmes les plus courants. Pour des problèmes spécifiques, consultez la documentation SSH officielle. 