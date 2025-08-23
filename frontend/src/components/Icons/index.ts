import React from 'react'
import { Ionicons } from '@expo/vector-icons'
import { MaterialIcons } from '@expo/vector-icons'
import { MaterialCommunityIcons } from '@expo/vector-icons'
import { FontAwesome } from '@expo/vector-icons'
import { FontAwesome5 } from '@expo/vector-icons'
import { AntDesign } from '@expo/vector-icons'
import { Entypo } from '@expo/vector-icons'
import { EvilIcons } from '@expo/vector-icons'
import { Feather } from '@expo/vector-icons'
import { Fontisto } from '@expo/vector-icons'
import { Foundation } from '@expo/vector-icons'
import { Octicons } from '@expo/vector-icons'
import { SimpleLineIcons } from '@expo/vector-icons'
import { Zocial } from '@expo/vector-icons'

// Types pour les props du composant Icon
export interface IconProps {
  name: string
  size?: number
  color?: string
  style?: any
}

// Mapping des préfixes vers les bibliothèques d'icônes
const iconLibraries: Record<string, any> = {
  // Ionicons (par défaut)
  'io': Ionicons,
  'ion': Ionicons,
  'ionic': Ionicons,
  
  // Material Icons
  'md': MaterialIcons,
  'mat': MaterialIcons,
  'material': MaterialIcons,
  
  // Material Community Icons
  'mc': MaterialCommunityIcons,
  'mci': MaterialCommunityIcons,
  'mdc': MaterialCommunityIcons,
  
  // FontAwesome
  'fa': FontAwesome,
  'fontawesome': FontAwesome,
  
  // FontAwesome5
  'fa5': FontAwesome5,
  'fas': FontAwesome5,
  'far': FontAwesome5,
  'fab': FontAwesome5,
  'fal': FontAwesome5,
  'fad': FontAwesome5,
  
  // AntDesign
  'ant': AntDesign,
  'antd': AntDesign,
  
  // Entypo
  'ent': Entypo,
  'entypo': Entypo,
  
  // EvilIcons
  'evil': EvilIcons,
  'evi': EvilIcons,
  
  // Feather
  'feather': Feather,
  'fea': Feather,
  
  // Fontisto
  'fontisto': Fontisto,
  'fst': Fontisto,
  
  // Foundation
  'foundation': Foundation,
  'fou': Foundation,
  
  // Octicons
  'oct': Octicons,
  'octicons': Octicons,
  
  // SimpleLineIcons
  'sli': SimpleLineIcons,
  'simpleline': SimpleLineIcons,
  
  // Zocial
  'zocial': Zocial,
  'zoc': Zocial,
}

// Fonction pour extraire le préfixe et le nom de l'icône
const parseIconName = (fullName: string): { prefix: string; name: string; library: any } => {
  // Si le nom contient un préfixe (ex: "io:home", "md:add", "fa:user")
  if (fullName.includes(':')) {
    const [prefix, name] = fullName.split(':')
    const library = iconLibraries[prefix.toLowerCase()]
    
    if (library) {
      return { prefix, name, library }
    }
  }
  
  // Par défaut, utiliser Ionicons
  return { 
    prefix: 'io', 
    name: fullName, 
    library: Ionicons 
  }
}

// Composant Icon principal
export const Icon: React.FC<IconProps> = ({ 
  name, 
  size = 24, 
  color = '#000', 
  style 
}) => {
  const { library: IconLibrary, name: iconName } = parseIconName(name)
  
  return React.createElement(IconLibrary, {
    name: iconName,
    size,
    color,
    style
  })
}

// Composants spécialisés pour chaque bibliothèque
export const IoIcon: React.FC<IconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `io:${props.name}` })
}

export const MdIcon: React.FC<IconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `md:${props.name}` })
}

export const McIcon: React.FC<IconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `mc:${props.name}` })
}

export const FaIcon: React.FC<IconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `fa:${props.name}` })
}

export const Fa5Icon: React.FC<IconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `fa5:${props.name}` })
}

export const AntIcon: React.FC<IconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `ant:${props.name}` })
}

export const EntIcon: React.FC<IconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `ent:${props.name}` })
}

export const EvilIcon: React.FC<IconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `evil:${props.name}` })
}

export const FeatherIcon: React.FC<IconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `feather:${props.name}` })
}

export const FontistoIcon: React.FC<IconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `fontisto:${props.name}` })
}

export const FoundationIcon: React.FC<IconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `foundation:${props.name}` })
}

export const OctIcon: React.FC<IconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `oct:${props.name}` })
}

export const SliIcon: React.FC<IconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `sli:${props.name}` })
}

export const ZocialIcon: React.FC<IconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `zocial:${props.name}` })
}

// Export par défaut
export default Icon 