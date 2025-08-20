// Exemples d'utilisation des alias d'imports

// ✅ Imports avec alias (recommandé)
import { ScreenContent } from '@/components/ScreenContent';
import '@/global.css';
import { Ionicons } from '@expo/vector-icons';

// ✅ Import depuis le fichier d'index des composants
import { EditScreenInfo } from '@/components';

// ❌ Imports relatifs (à éviter)
// import { ScreenContent } from '../../components/ScreenContent';
// import '../global.css';

export default function ExampleComponent() {
  return (
    <ScreenContent title="Example" path="examples/imports.tsx">
      <p>This component demonstrates the use of import aliases.</p>
    </ScreenContent>
  );
} 