import React from 'react'
import { TouchableOpacity, Text, StyleSheet, ViewStyle, TextStyle } from 'react-native'
import Animated, { 
  useSharedValue, 
  useAnimatedStyle, 
  withSpring,
  withTiming,
  interpolateColor,
  runOnUI
} from 'react-native-reanimated'
import { colors, spacing, borderRadius, typography, shadows } from '@/src/theme'

export interface ButtonProps {
  title: string
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger'
  size?: 'small' | 'medium' | 'large'
  onPress: () => void
  disabled?: boolean
  loading?: boolean
  icon?: React.ReactNode
  fullWidth?: boolean
  style?: ViewStyle
  textStyle?: TextStyle
}

const AnimatedTouchableOpacity = Animated.createAnimatedComponent(TouchableOpacity)

export const Button: React.FC<ButtonProps> = ({
  title,
  variant = 'primary',
  size = 'medium',
  onPress,
  disabled = false,
  loading = false,
  icon,
  fullWidth = false,
  style,
  textStyle,
}) => {
  const pressed = useSharedValue(0)
  const scale = useSharedValue(1)

  const animatedStyle = useAnimatedStyle(() => {
    'worklet'
    const variantColors = getVariantColor(variant)
    return {
      transform: [
        { scale: scale.value },
      ],
      backgroundColor: interpolateColor(
        pressed.value,
        [0, 1],
        [
          variantColors.background,
          variantColors.pressed,
        ]
      ),
    }
  })

  const handlePressIn = () => {
    pressed.value = withTiming(1, { duration: 100 })
    scale.value = withSpring(0.95)
  }

  const handlePressOut = () => {
    pressed.value = withTiming(0, { duration: 100 })
    scale.value = withSpring(1)
  }

  const buttonStyle = [
    styles.button,
    styles[variant],
    styles[size],
    fullWidth && styles.fullWidth,
    disabled && styles.disabled,
    animatedStyle,
    style,
  ]

  const textStyleComposed = [
    styles.text,
    styles[`${variant}Text`],
    styles[`${size}Text`],
    disabled && styles.disabledText,
    textStyle,
  ]

  return (
    <AnimatedTouchableOpacity
      style={buttonStyle}
      onPress={onPress}
      onPressIn={handlePressIn}
      onPressOut={handlePressOut}
      disabled={disabled || loading}
      activeOpacity={0.8}
    >
      {icon && <Animated.View style={styles.icon}>{icon}</Animated.View>}
      <Text style={textStyleComposed}>
        {loading ? 'Chargement...' : title}
      </Text>
    </AnimatedTouchableOpacity>
  )
}

const getVariantColor = (variant: string) => {
  'worklet'
  switch (variant) {
    case 'primary':
      return {
        background: colors.primary[500],
        pressed: colors.primary[600],
      }
    case 'secondary':
      return {
        background: colors.secondary[500],
        pressed: colors.secondary[600],
      }
    case 'outline':
      return {
        background: 'transparent',
        pressed: colors.neutral[100],
      }
    case 'ghost':
      return {
        background: 'transparent',
        pressed: colors.neutral[100],
      }
    case 'danger':
      return {
        background: colors.error[500],
        pressed: colors.error[600],
      }
    default:
      return {
        background: colors.primary[500],
        pressed: colors.primary[600],
      }
  }
}

const styles = StyleSheet.create({
  button: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: borderRadius.md,
    ...shadows.sm,
  },
  
  // Variants
  primary: {
    backgroundColor: colors.primary[500],
  },
  secondary: {
    backgroundColor: colors.secondary[500],
  },
  outline: {
    backgroundColor: 'transparent',
    borderWidth: 1,
    borderColor: colors.primary[500],
  },
  ghost: {
    backgroundColor: 'transparent',
  },
  danger: {
    backgroundColor: colors.error[500],
  },
  
  // Sizes
  small: {
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    minHeight: 36,
  },
  medium: {
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    minHeight: 44,
  },
  large: {
    paddingHorizontal: spacing.xl,
    paddingVertical: spacing.lg,
    minHeight: 52,
  },
  
  // Text styles
  text: {
    fontWeight: typography.weights.medium,
    textAlign: 'center',
  },
  primaryText: {
    color: colors.text.inverse,
  },
  secondaryText: {
    color: colors.text.inverse,
  },
  outlineText: {
    color: colors.primary[500],
  },
  ghostText: {
    color: colors.primary[500],
  },
  dangerText: {
    color: colors.text.inverse,
  },
  smallText: {
    fontSize: typography.sizes.sm,
  },
  mediumText: {
    fontSize: typography.sizes.md,
  },
  largeText: {
    fontSize: typography.sizes.lg,
  },
  
  // States
  disabled: {
    opacity: 0.5,
  },
  disabledText: {
    opacity: 0.7,
  },
  
  // Layout
  fullWidth: {
    width: '100%',
  },
  
  // Icon
  icon: {
    marginRight: spacing.sm,
  },
})

export default Button 