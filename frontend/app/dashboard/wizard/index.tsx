import { useAuthStore } from "@/src";
import { SetupWizard, userSetupApi, WizardData } from "@/src/components/wizard"; 
import { router } from "expo-router";
import { useState } from "react";
import { Text, View } from "react-native";

export default function Wizard() {

    const { setRequirePreferences } = useAuthStore()
    const [userData, setUserData] = useState<WizardData | null>(null);

  const handleWizardComplete = async (data: WizardData) => {
    console.log('Données du wizard:', JSON.stringify(data, null, 2));
    
    // Ici vous enverriez les données à votre backend
    await userSetupApi.saveUserSetup(data);
    router.push("/dashboard/tabs/overview");
    setUserData(data);
    setRequirePreferences(false);   
  };

  const closeWizard = () => {
    setRequirePreferences(false);
    router.push("/dashboard/tabs/overview");
  };

  return (
    <View className="flex-1 bg-gray-50 mt-safe justify-center items-center">
            <SetupWizard
        visible={true}
        onClose={closeWizard}
        onComplete={handleWizardComplete}
      />
      <Text>
        {JSON.stringify(userData, null, 2)}
      </Text>
    </View>
  )
}