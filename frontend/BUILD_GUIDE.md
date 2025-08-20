# Guide de Build Android - Mother AI

Ce guide explique comment créer des builds de développement pour Android.

## Prérequis

1. **EAS CLI installé** :
   ```bash
   npm install -g @expo/eas-cli
   ```

2. **Compte Expo configuré** :
   ```bash
   eas login
   ```

3. **Projet configuré** :
   ```bash
   eas build:configure
   ```

## Types de Builds

### 1. Build de Développement (Development)

**Commande :**
```bash
npm run build:dev:android
# ou
eas build --platform android --profile development
```

**Caractéristiques :**
- ✅ Client de développement inclus
- ✅ Debug activé
- ✅ Distribution interne
- ✅ Format APK pour installation directe
- ✅ Variables d'environnement : `EXPO_PUBLIC_ENV=development`

**Utilisation :**
- Tests sur appareil physique
- Développement avec hot reload
- Debugging avancé

### 2. Build de Preview

**Commande :**
```bash
npm run build:preview:android
# ou
eas build --platform android --profile preview
```

**Caractéristiques :**
- ✅ Distribution interne
- ✅ Format APK
- ✅ Variables d'environnement : `EXPO_PUBLIC_ENV=preview`

**Utilisation :**
- Tests de fonctionnalités
- Validation par l'équipe
- Tests de performance

### 3. Build de Production

**Commande :**
```bash
npm run build:prod:android
# ou
eas build --platform android --profile production
```

**Caractéristiques :**
- ✅ Auto-increment des versions
- ✅ Optimisé pour la production
- ✅ Variables d'environnement : `EXPO_PUBLIC_ENV=production`

**Utilisation :**
- Publication sur Google Play Store
- Distribution finale

## Installation du Build de Développement

### Méthode 1 : Téléchargement Direct

1. Lancez le build :
   ```bash
   npm run build:dev:android
   ```

2. Attendez la fin du build (5-10 minutes)

3. Téléchargez l'APK depuis le lien fourni

4. Installez sur votre appareil Android :
   ```bash
   adb install path/to/your-app.apk
   ```

### Méthode 2 : Installation via QR Code

1. Scannez le QR code avec l'app Expo Go
2. L'app se téléchargera automatiquement

## Configuration Avancée

### Variables d'Environnement

Les variables d'environnement sont définies dans `eas.json` :

```json
{
  "build": {
    "development": {
      "env": {
        "EXPO_PUBLIC_ENV": "development"
      }
    }
  }
}
```

### Permissions Android

Les permissions sont configurées dans `app.json` :

```json
{
  "android": {
    "permissions": [
      "INTERNET",
      "ACCESS_NETWORK_STATE"
    ]
  }
}
```

## Dépannage

### Erreur de Build

1. **Vérifiez la configuration** :
   ```bash
   eas build:configure
   ```

2. **Nettoyez le cache** :
   ```bash
   expo r -c
   ```

3. **Vérifiez les logs** :
   ```bash
   eas build:list
   ```

### Problèmes d'Installation

1. **Activez les sources inconnues** sur votre appareil Android
2. **Désinstallez** l'ancienne version si elle existe
3. **Vérifiez l'espace** disponible sur l'appareil

## Scripts Utiles

```bash
# Build de développement
npm run build:dev:android

# Build de preview
npm run build:preview:android

# Build de production
npm run build:prod:android

# Submit vers les stores
npm run submit:android
```

## Monitoring

- **EAS Dashboard** : https://expo.dev/accounts/flowerr/projects/mother-ai
- **Build History** : Voir tous les builds passés
- **Analytics** : Suivi des performances

## Support

Pour toute question ou problème :
1. Consultez la [documentation EAS](https://docs.expo.dev/build/introduction/)
2. Vérifiez les [logs de build](https://docs.expo.dev/build-reference/build-logs/)
3. Contactez l'équipe de développement 