// Export du composant principal
export { SetupWizard } from './SetupWizard';

// Export des types
export type { WizardData } from './SetupWizard';

// Export de l'exemple
export { default as WizardExample } from './example';

// Export du service API
export { userSetupApi } from '@/src/services/userSetupService/userSetupApi';
export type { UserSetupResponse, UserSetupError } from '@/src/services/userSetupService/userSetupApi';

// Export du hook personnalisé
export { useSetupWizard } from '@/src/hooks/useSetupWizard'; 