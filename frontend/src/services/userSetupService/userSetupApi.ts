import { apiClient, apiHelpers } from '../api/client';
import { WizardData } from '@/src/components/wizard/SetupWizard';

export interface UserSetupResponse {
  success: boolean;
  message: string;
  userId?: string;
  setupId?: string;
}

export interface UserSetupError {
  error: string;
  details?: string;
  field?: string;
}

class UserSetupApi {
  private baseUrl = '/user/setup';

  /**
   * Sauvegarde la configuration utilisateur du wizard
   */
  async saveUserSetup(data: WizardData): Promise<UserSetupResponse> {
    try {
      const response = await apiHelpers.post("/preferences", data);
      console.log(response,"response of saveUserSetup");
      return response.data;
    } catch (error: any) {
      throw this.handleError(error);
    }
  }

  /**
   * Récupère la configuration existante de l'utilisateur
   */
  async getUserSetup(): Promise<WizardData> {
    try {
      const response = await apiClient.get(this.baseUrl);
      return response.data;
    } catch (error: any) {
      throw this.handleError(error);
    }
  }

  /**
   * Met à jour la configuration utilisateur
   */
  async updateUserSetup(data: Partial<WizardData>): Promise<UserSetupResponse> {
    try {
      const response = await apiClient.put(this.baseUrl, data);
      return response.data;
    } catch (error: any) {
      throw this.handleError(error);
    }
  }

  /**
   * Supprime la configuration utilisateur
   */
  async deleteUserSetup(): Promise<UserSetupResponse> {
    try {
      const response = await apiClient.delete(this.baseUrl);
      return response.data;
    } catch (error: any) {
      throw this.handleError(error);
    }
  }

  /**
   * Valide les données du wizard côté serveur
   */
  async validateSetup(data: WizardData): Promise<{ valid: boolean; errors?: UserSetupError[] }> {
    try {
      const response = await apiClient.post(`${this.baseUrl}/validate`, data);
      return response.data;
    } catch (error: any) {
      throw this.handleError(error);
    }
  }

  /**
   * Gestion centralisée des erreurs
   */
  private handleError(error: any): UserSetupError {
    if (error.response) {
      // Erreur de réponse du serveur
      const { status, data } = error.response;
      
      switch (status) {
        case 400:
          return {
            error: 'Données invalides',
            details: data.message || 'Vérifiez vos informations',
            field: data.field
          };
        case 401:
          return {
            error: 'Non autorisé',
            details: 'Vous devez être connecté pour effectuer cette action'
          };
        case 403:
          return {
            error: 'Accès refusé',
            details: 'Vous n\'avez pas les permissions nécessaires'
          };
        case 404:
          return {
            error: 'Configuration introuvable',
            details: 'Aucune configuration trouvée pour cet utilisateur'
          };
        case 409:
          return {
            error: 'Configuration existante',
            details: 'Une configuration existe déjà pour cet utilisateur'
          };
        case 422:
          return {
            error: 'Validation échouée',
            details: data.message || 'Certaines données sont invalides',
            field: data.field
          };
        case 500:
          return {
            error: 'Erreur serveur',
            details: 'Une erreur interne s\'est produite'
          };
        default:
          return {
            error: 'Erreur inconnue',
            details: `Erreur ${status}: ${data.message || 'Erreur inconnue'}`
          };
      }
    } else if (error.request) {
      // Erreur de réseau
      return {
        error: 'Erreur de connexion',
        details: 'Impossible de se connecter au serveur'
      };
    } else {
      // Erreur locale
      return {
        error: 'Erreur locale',
        details: error.message || 'Une erreur s\'est produite'
      };
    }
  }
}

export const userSetupApi = new UserSetupApi(); 