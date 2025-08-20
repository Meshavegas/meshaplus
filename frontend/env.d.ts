declare module '@env' {
  export const EXPO_PUBLIC_ENV: 'development' | 'preview' | 'production';
  export const EXPO_PUBLIC_API_URL: string;
  export const EXPO_PUBLIC_APP_NAME: string;
}

declare global {
  namespace NodeJS {
    interface ProcessEnv {
      EXPO_PUBLIC_ENV: 'development' | 'preview' | 'production';
      EXPO_PUBLIC_API_URL: string;
      EXPO_PUBLIC_APP_NAME: string;
    }
  }
} 