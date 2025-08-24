import React from 'react'
import { Ionicons } from '@expo/vector-icons'
import { MaterialIcons } from '@expo/vector-icons'
import { MaterialCommunityIcons } from '@expo/vector-icons'
import { FontAwesome } from '@expo/vector-icons'
import { FontAwesome5 } from '@expo/vector-icons'
import { FontAwesome6 } from '@expo/vector-icons'
import { AntDesign } from '@expo/vector-icons'
import { Entypo } from '@expo/vector-icons'
import { EvilIcons } from '@expo/vector-icons'
import { Feather } from '@expo/vector-icons'
import { Fontisto } from '@expo/vector-icons'
import { Foundation } from '@expo/vector-icons'
import { Octicons } from '@expo/vector-icons'
import { SimpleLineIcons } from '@expo/vector-icons'
import { Zocial } from '@expo/vector-icons'
import { IONICONS, FONTAWESOME_6, type IconName } from './iconNames'

// Types pour les noms d'icônes de chaque bibliothèque
export type IconLibraryType = 
  | typeof Ionicons
  | typeof MaterialIcons
  | typeof MaterialCommunityIcons
  | typeof FontAwesome
  | typeof FontAwesome5
  | typeof FontAwesome6
  | typeof AntDesign
  | typeof Entypo
  | typeof EvilIcons
  | typeof Feather
  | typeof Fontisto
  | typeof Foundation
  | typeof Octicons
  | typeof SimpleLineIcons
  | typeof Zocial

// Types pour les préfixes d'icônes
export type IconPrefix = 
  | 'io' | 'ion' | 'ionic'
  | 'md' | 'mat' | 'material'
  | 'mc' | 'mci' | 'mdc'
  | 'fa' | 'fontawesome'
  | 'fa5' | 'fas' | 'far' | 'fab' | 'fal' | 'fad'
  | 'fa6' | 'fa-solid' | 'fa-regular' | 'fa-brands' | 'fa-light' | 'fa-duotone'
  | 'ant' | 'antd'
  | 'ent' | 'entypo'
  | 'evil' | 'evi'
  | 'feather' | 'fea'
  | 'fontisto' | 'fst'
  | 'foundation' | 'fou'
  | 'oct' | 'octicons'
  | 'sli' | 'simpleline'
  | 'zocial' | 'zoc'

// Types pour les props du composant Icon
export interface IconProps {
  name: string | IconName
  size?: number
  color?: string
  style?: React.ComponentProps<any>['style']
}

// Interface pour le résultat du parsing
interface ParsedIcon {
  prefix: IconPrefix
  name: string
  library: IconLibraryType
}

// Mapping des préfixes vers les bibliothèques d'icônes avec typage strict
const iconLibraries: Record<IconPrefix, IconLibraryType> = {
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
  
  // FontAwesome6
  'fa6': FontAwesome6,
  'fa-solid': FontAwesome6,
  'fa-regular': FontAwesome6,
  'fa-brands': FontAwesome6,
  'fa-light': FontAwesome6,
  'fa-duotone': FontAwesome6,
  
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

// Fonction pour extraire le préfixe et le nom de l'icône avec typage strict
const parseIconName = (fullName: string): ParsedIcon => {  
  // Si le nom contient un préfixe (ex: "io:home", "md:add", "fa:user")
  if (fullName.includes(':')) {
    const [prefix, name] = fullName.split(':')
    const normalizedPrefix = prefix.toLowerCase()
    
    // Vérifier si le préfixe existe dans notre mapping
    if (normalizedPrefix in iconLibraries) {
      const library = iconLibraries[normalizedPrefix as IconPrefix]
      return { 
        prefix: normalizedPrefix as IconPrefix, 
        name, 
        library 
      }
    }
  }
  
  // Par défaut, utiliser Ionicons
  return { 
    prefix: 'io', 
    name: fullName, 
    library: Ionicons 
  }
}

// Composant Icon principal avec typage strict
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

// Types pour les props des composants spécialisés
export interface SpecializedIconProps extends Omit<IconProps, 'name'> {
  name: string
}

// Composants spécialisés pour chaque bibliothèque avec typage strict
export const IoIcon: React.FC<SpecializedIconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `io:${props.name}` })
}

export const MdIcon: React.FC<SpecializedIconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `md:${props.name}` })
}

export const McIcon: React.FC<SpecializedIconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `mc:${props.name}` })
}

export const FaIcon: React.FC<SpecializedIconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `fa:${props.name}` })
}

export const Fa5Icon: React.FC<SpecializedIconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `fa5:${props.name}` })
}

export const Fa6Icon: React.FC<SpecializedIconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `fa6:${props.name}` })
}

export const AntIcon: React.FC<SpecializedIconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `ant:${props.name}` })
}

export const EntIcon: React.FC<SpecializedIconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `ent:${props.name}` })
}

export const EvilIcon: React.FC<SpecializedIconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `evil:${props.name}` })
}

export const FeatherIcon: React.FC<SpecializedIconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `feather:${props.name}` })
}

export const FontistoIcon: React.FC<SpecializedIconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `fontisto:${props.name}` })
}

export const FoundationIcon: React.FC<SpecializedIconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `foundation:${props.name}` })
}

export const OctIcon: React.FC<SpecializedIconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `oct:${props.name}` })
}

export const SliIcon: React.FC<SpecializedIconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `sli:${props.name}` })
}

export const ZocialIcon: React.FC<SpecializedIconProps> = (props) => {
  return React.createElement(Icon, { ...props, name: `zocial:${props.name}` })
}

// Export des constantes d'icônes
export { IONICONS, FONTAWESOME_6, FONTAWESOME_5, MATERIAL_COMMUNITY_ICONS } from './iconNames'

// Export par défaut
export default Icon 