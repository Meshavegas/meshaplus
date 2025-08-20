const { getDefaultConfig } = require('expo/metro-config');
const { withNativeWind } = require('nativewind/metro');

const config = getDefaultConfig(__dirname);

// Add support for Expo Router
config.resolver.assetExts.push('cjs');

module.exports = withNativeWind(config, { input: './global.css' });
