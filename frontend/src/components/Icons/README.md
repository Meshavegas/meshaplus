# Composant Icon Unifié

Ce composant vous permet d'utiliser toutes les bibliothèques d'icônes d'Expo Vector Icons avec un seul import et un système de préfixes simple.

## Installation

```bash
# Assurez-vous que @expo/vector-icons est installé
npm install @expo/vector-icons
```

## Utilisation

### Import principal

```tsx
import { Icon } from '@/src/components/Icons'
```

### Utilisation avec préfixes

```tsx
// Ionicons (par défaut)
<Icon name="home" size={24} color="#000" />

// Avec préfixe explicite
<Icon name="io:home" size={24} color="#000" />
<Icon name="md:add" size={24} color="#000" />
<Icon name="fa:user" size={24} color="#000" />
<Icon name="ant:heart" size={24} color="#000" />
```

### Composants spécialisés

```tsx
import { 
  IoIcon,    // Ionicons
  MdIcon,    // Material Icons
  McIcon,    // Material Community Icons
  FaIcon,    // FontAwesome
  Fa5Icon,   // FontAwesome5
  AntIcon,   // AntDesign
  EntIcon,   // Entypo
  EvilIcon,  // EvilIcons
  FeatherIcon, // Feather
  FontistoIcon, // Fontisto
  FoundationIcon, // Foundation
  OctIcon,   // Octicons
  SliIcon,   // SimpleLineIcons
  ZocialIcon // Zocial
} from '@/src/components/Icons'

// Utilisation
<IoIcon name="home" size={24} color="#000" />
<MdIcon name="add" size={24} color="#000" />
<FaIcon name="user" size={24} color="#000" />
<AntIcon name="heart" size={24} color="#000" />
```

## Préfixes disponibles

| Préfixe | Bibliothèque | Exemples |
|---------|-------------|----------|
| `io`, `ion`, `ionic` | Ionicons | `home`, `add`, `person` |
| `md`, `mat`, `material` | Material Icons | `add`, `delete`, `edit` |
| `mc`, `mci`, `mdc` | Material Community Icons | `account`, `home`, `settings` |
| `fa`, `fontawesome` | FontAwesome | `user`, `heart`, `star` |
| `fa5`, `fas`, `far`, `fab` | FontAwesome5 | `user`, `heart`, `star` |
| `ant`, `antd` | AntDesign | `heart`, `star`, `like` |
| `ent`, `entypo` | Entypo | `home`, `user`, `mail` |
| `evil`, `evi` | EvilIcons | `heart`, `star`, `like` |
| `feather`, `fea` | Feather | `home`, `user`, `mail` |
| `fontisto`, `fst` | Fontisto | `home`, `user`, `mail` |
| `foundation`, `fou` | Foundation | `home`, `user`, `mail` |
| `oct`, `octicons` | Octicons | `home`, `user`, `mail` |
| `sli`, `simpleline` | SimpleLineIcons | `home`, `user`, `mail` |
| `zocial`, `zoc` | Zocial | `home`, `user`, `mail` |

## Exemples d'utilisation

```tsx
import React from 'react'
import { View, Text } from 'react-native'
import { Icon, IoIcon, MdIcon, FaIcon } from '@/src/components/Icons'

export default function IconExample() {
  return (
    <View style={{ padding: 20 }}>
      {/* Utilisation avec préfixe */}
      <Icon name="io:home" size={24} color="#007AFF" />
      <Icon name="md:add" size={24} color="#FF3B30" />
      <Icon name="fa:user" size={24} color="#34C759" />
      
      {/* Utilisation sans préfixe (Ionicons par défaut) */}
      <Icon name="home" size={24} color="#007AFF" />
      
      {/* Composants spécialisés */}
      <IoIcon name="home" size={24} color="#007AFF" />
      <MdIcon name="add" size={24} color="#FF3B30" />
      <FaIcon name="user" size={24} color="#34C759" />
    </View>
  )
}
```

## Props

| Prop | Type | Défaut | Description |
|------|------|--------|-------------|
| `name` | `string` | - | Nom de l'icône (avec ou sans préfixe) |
| `size` | `number` | `24` | Taille de l'icône |
| `color` | `string` | `#000` | Couleur de l'icône |
| `style` | `any` | - | Styles supplémentaires |

## Avantages

- ✅ **Un seul import** pour toutes les bibliothèques d'icônes
- ✅ **Système de préfixes simple** et intuitif
- ✅ **Ionicons par défaut** (pas besoin de préfixe)
- ✅ **Composants spécialisés** pour chaque bibliothèque
- ✅ **TypeScript support** complet
- ✅ **Performance optimisée** avec React.createElement

## Migration depuis les imports directs

**Avant :**
```tsx
import { Ionicons } from '@expo/vector-icons'
import { MaterialIcons } from '@expo/vector-icons'
import { FontAwesome } from '@expo/vector-icons'

<Ionicons name="home" size={24} color="#000" />
<MaterialIcons name="add" size={24} color="#000" />
<FontAwesome name="user" size={24} color="#000" />
```

**Après :**
```tsx
import { Icon } from '@/src/components/Icons'

<Icon name="io:home" size={24} color="#000" />
<Icon name="md:add" size={24} color="#000" />
<Icon name="fa:user" size={24} color="#000" />
``` 