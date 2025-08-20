import { ExpoConfig, ConfigContext } from 'expo/config';

export default ({ config }: ConfigContext): ExpoConfig => ({
  ...config,
  name: 'Mother AI',
  slug: 'mother-ai',
  version: '1.0.0',
  orientation: 'portrait',
  icon: './assets/icon.png',
  userInterfaceStyle: 'light',
  splash: {
    image: './assets/splash.png',
    resizeMode: 'contain',
    backgroundColor: '#ffffff',
  },
  assetBundlePatterns: ['**/*'],
  ios: {
    supportsTablet: true,
    bundleIdentifier: 'com.meshaplus.app',
  },
  android: {
    package: 'com.meshaplus.app',
    versionCode: 1,
    adaptiveIcon: {
      foregroundImage: './assets/adaptive-icon.png',
      backgroundColor: '#ffffff',
    },
    permissions: ['INTERNET', 'ACCESS_NETWORK_STATE'],
  },
  web: {
    favicon: './assets/favicon.png',
  },
  experiments: {
    tsconfigPaths: true,
  },
  plugins: ['expo-router'],
  extra: {
    router: {},
    eas: {
      projectId: 'a2484208-e073-46b8-a28e-84a8931c4c2a',
    },
  },
  owner: 'flowerr',
}); 