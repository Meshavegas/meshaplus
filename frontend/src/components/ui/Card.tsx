import React from 'react'
import { View, StyleSheet, ViewStyle, Pressable } from 'react-native'
import Animated, { 
  useSharedValue, 
  useAnimatedStyle, 
  withSpring,
  withTiming
} from 'react-native-reanimated'
import { colors, spacing, borderRadius, shadows } from '@/src/theme'

export interface CardProps {
  children: React.ReactNode
  variant?: 'default' | 'elevated' | 'outlined' | 'flat'
  padding?: 'none' | 'small' | 'medium' | 'large'
  margin?: 'none' | 'small' | 'medium' | 'large'
  onPress?: () => void
  style?: ViewStyle
  animated?: boolean
}

const AnimatedView = Animated.createAnimatedComponent(View)

export const Card: React.FC<CardProps> = ({
  children,
  variant = 'default',
  padding = 'medium',
  margin = 'none',
  onPress,
  style,
  animated = true,
}) => {
  const scale = useSharedValue(1)
  const pressed = useSharedValue(0)

  const animatedStyle = useAnimatedStyle(() => {
    if (!animated) return {}
    
    return {
      transform: [
        { scale: scale.value },
      ],
      opacity: withTiming(pressed.value === 1 ? 0.9 : 1, { duration: 100 }),
    }
  })

  const handlePressIn = () => {
    if (animated) {
      scale.value = withSpring(0.98)
      pressed.value = 1
    }
  }

  const handlePressOut = () => {
    if (animated) {
      scale.value = withSpring(1)
      pressed.value = 0
    }
  }

  const getPaddingStyle = () => {
    switch (padding) {
      case 'none': return {}
      case 'small': return { padding: spacing.sm }
      case 'medium': return { padding: spacing.md }
      case 'large': return { padding: spacing.lg }
      default: return { padding: spacing.md }
    }
  }

  const getMarginStyle = () => {
    switch (margin) {
      case 'none': return {}
      case 'small': return { margin: spacing.sm }
      case 'medium': return { margin: spacing.md }
      case 'large': return { margin: spacing.lg }
      default: return {}
    }
  }

  const cardStyle = [
    styles.card,
    styles[variant],
    getPaddingStyle(),
    getMarginStyle(),
    animatedStyle,
    style,
  ]

  if (onPress) {
    return (
      <AnimatedView
        style={cardStyle}
        onTouchStart={handlePressIn}
        onTouchEnd={handlePressOut}
      >
        <Pressable onPress={onPress}>   
          {children}
        </Pressable>
      </AnimatedView>
    )
  }

  return (
    <AnimatedView style={cardStyle}>
      {children}
    </AnimatedView>
  )
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.background.primary,
    borderRadius: borderRadius.lg,
  },
  
  // Variants
  default: {
    ...shadows.sm,
  },
  elevated: {
    ...shadows.lg,
  },
  outlined: {
    borderWidth: 1,
    borderColor: colors.neutral[200],
  },
  flat: {
    backgroundColor: colors.background.secondary,
  },
})

export default Card 