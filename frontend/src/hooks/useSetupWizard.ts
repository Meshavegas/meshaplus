import { useState, useCallback } from 'react';
import { Alert } from 'react-native';
import { WizardData } from '@/src/components/wizard/SetupWizard';
import { userSetupApi, UserSetupError } from '@/src/services/userSetupService/userSetupApi';

interface UseSetupWizardReturn {
  showWizard: boolean;
  isLoading: boolean;
  userData: WizardData | null;
  openWizard: () => void;
  closeWizard: () => void;
  handleWizardComplete: (data: WizardData) => Promise<void>;
  loadUserSetup: () => Promise<void>;
  updateUserSetup: (data: Partial<WizardData>) => Promise<void>;
  resetUserSetup: () => Promise<void>;
}

export const useSetupWizard = (): UseSetupWizardReturn => {
  const [showWizard, setShowWizard] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [userData, setUserData] = useState<WizardData | null>(null);

  const openWizard = useCallback(() => {
    setShowWizard(true);
  }, []);

  const closeWizard = useCallback(() => {
    setShowWizard(false);
  }, []);

  const handleWizardComplete = useCallback(async (data: WizardData) => {
    setIsLoading(true);
    try {
      console.log('Configuration terminée:', JSON.stringify(data, null, 2));
      
      // Envoyer les données au backend
      const response = await userSetupApi.saveUserSetup(data);
      
      // Sauvegarder localement
      setUserData(data);
      
      Alert.alert(
        'Succès ! 🎉',
        'Votre configuration a été sauvegardée avec succès. Votre expérience est maintenant personnalisée !',
        [
          {
            text: 'Continuer',
            onPress: () => setShowWizard(false)
          }
        ]
      );
      
    } catch (error: any) {
      const setupError = error as UserSetupError;
      Alert.alert(
        'Erreur',
        setupError.details || 'Impossible de sauvegarder la configuration. Veuillez réessayer.',
        [
          {
            text: 'OK',
            onPress: () => setShowWizard(false)
          }
        ]
      );
    } finally {
      setIsLoading(false);
    }
  }, []);

  const loadUserSetup = useCallback(async () => {
    setIsLoading(true);
    try {
      const data = await userSetupApi.getUserSetup();
      setUserData(data);
    } catch (error: any) {
      const setupError = error as UserSetupError;
      console.warn('Impossible de charger la configuration:', setupError.details);
      // Ne pas afficher d'alerte pour cette erreur car c'est normal si l'utilisateur n'a pas encore configuré
    } finally {
      setIsLoading(false);
    }
  }, []);

  const updateUserSetup = useCallback(async (data: Partial<WizardData>) => {
    setIsLoading(true);
    try {
      const response = await userSetupApi.updateUserSetup(data);
      
      // Mettre à jour les données locales
      if (userData) {
        setUserData({ ...userData, ...data });
      }
      
      Alert.alert('Succès', 'Configuration mise à jour avec succès !');
      
    } catch (error: any) {
      const setupError = error as UserSetupError;
      Alert.alert(
        'Erreur',
        setupError.details || 'Impossible de mettre à jour la configuration.'
      );
    } finally {
      setIsLoading(false);
    }
  }, [userData]);

  const resetUserSetup = useCallback(async () => {
    try {
      await userSetupApi.deleteUserSetup();
      setUserData(null);
      Alert.alert('Succès', 'Configuration réinitialisée avec succès !');
    } catch (error: any) {
      const setupError = error as UserSetupError;
      Alert.alert(
        'Erreur',
        setupError.details || 'Impossible de réinitialiser la configuration.'
      );
    }
  }, []);

  return {
    showWizard,
    isLoading,
    userData,
    openWizard,
    closeWizard,
    handleWizardComplete,
    loadUserSetup,
    updateUserSetup,
    resetUserSetup,
  };
}; 