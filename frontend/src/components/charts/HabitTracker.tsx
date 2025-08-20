import React from 'react'
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native'
import Animated, { 
  useSharedValue, 
  useAnimatedStyle, 
  withSpring,
  withTiming,
  withSequence,
  withDelay
} from 'react-native-reanimated'
import { colors, spacing, borderRadius, typography, shadows } from '@/src/theme'

interface HabitDay {
  date: Date
  completed: boolean
  streak: number
}

interface HabitTrackerProps {
  habit: {
    id: string
    title: string
    color: string
    icon: string
    currentStreak: number
    longestStreak: number
  }
  days: HabitDay[]
  onDayPress: (date: Date) => void
  size?: 'small' | 'medium' | 'large'
  animated?: boolean
}

export const HabitTracker: React.FC<HabitTrackerProps> = ({
  habit,
  days,
  onDayPress,
  size = 'medium',
  animated = true,
}) => {
  const scaleValues = days.map(() => useSharedValue(1))
  const opacityValues = days.map(() => useSharedValue(0))
  const streakScale = useSharedValue(1)

  React.useEffect(() => {
    if (animated) {
      // Animation d'entrée en cascade
      days.forEach((_, index) => {
        opacityValues[index].value = withDelay(
          index * 50,
          withTiming(1, { duration: 300 })
        )
        scaleValues[index].value = withDelay(
          index * 50,
          withSpring(1, { damping: 15, stiffness: 150 })
        )
      })

      // Animation du streak
      streakScale.value = withSequence(
        withTiming(1.1, { duration: 200 }),
        withTiming(1, { duration: 200 })
      )
    }
  }, [animated])

  const getDaySize = () => {
    switch (size) {
      case 'small':
        return 24
      case 'large':
        return 40
      default:
        return 32
    }
  }

  const daySize = getDaySize()

  const renderDay = (day: HabitDay, index: number) => {
    const animatedStyle = useAnimatedStyle(() => ({
      transform: [{ scale: scaleValues[index].value }],
      opacity: opacityValues[index].value,
    }))

    const handlePress = () => {
      if (animated) {
        scaleValues[index].value = withSequence(
          withTiming(0.8, { duration: 100 }),
          withTiming(1.1, { duration: 100 }),
          withSpring(1, { damping: 15, stiffness: 150 })
        )
      }
      onDayPress(day.date)
    }

    return (
      <Animated.View key={day.date.toISOString()} style={[animatedStyle]}>
        <TouchableOpacity
          style={[
            styles.day,
            {
              width: daySize,
              height: daySize,
              borderRadius: daySize / 2,
              backgroundColor: day.completed ? habit.color : colors.neutral[100],
              borderColor: day.completed ? habit.color : colors.neutral[200],
            },
          ]}
          onPress={handlePress}
          activeOpacity={0.7}
        >
          {day.completed && (
            <Animated.View style={styles.checkmark}>
              <Text style={[styles.checkmarkText, { color: colors.text.inverse }]}>
                ✓
              </Text>
            </Animated.View>
          )}
        </TouchableOpacity>
      </Animated.View>
    )
  }

  const streakAnimatedStyle = useAnimatedStyle(() => ({
    transform: [{ scale: streakScale.value }],
  }))

  return (
    <View style={styles.container}>
      {/* Header avec titre et streak */}
      <View style={styles.header}>
        <View style={styles.titleContainer}>
          <Text style={styles.title}>{habit.title}</Text>
          <Text style={styles.subtitle}>
            Streak actuel: {habit.currentStreak} jours
          </Text>
        </View>
        
        <Animated.View style={[styles.streakContainer, streakAnimatedStyle]}>
          <Text style={styles.streakNumber}>{habit.currentStreak}</Text>
          <Text style={styles.streakLabel}>jours</Text>
        </Animated.View>
      </View>

      {/* Grille des jours */}
      <View style={styles.daysGrid}>
        {days.map((day, index) => renderDay(day, index))}
      </View>

      {/* Statistiques */}
      <View style={styles.stats}>
        <View style={styles.stat}>
          <Text style={styles.statLabel}>Plus long streak</Text>
          <Text style={styles.statValue}>{habit.longestStreak} jours</Text>
        </View>
        <View style={styles.stat}>
          <Text style={styles.statLabel}>Taux de réussite</Text>
          <Text style={styles.statValue}>
            {Math.round((days.filter(d => d.completed).length / days.length) * 100)}%
          </Text>
        </View>
      </View>
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.background.primary,
    borderRadius: borderRadius.lg,
    padding: spacing.lg,
    ...shadows.md,
  },
  
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: spacing.lg,
  },
  
  titleContainer: {
    flex: 1,
  },
  
  title: {
    fontSize: typography.sizes.lg,
    fontWeight: typography.weights.bold,
    color: colors.text.primary,
    marginBottom: spacing.xs,
  },
  
  subtitle: {
    fontSize: typography.sizes.sm,
    color: colors.text.secondary,
  },
  
  streakContainer: {
    alignItems: 'center',
    backgroundColor: colors.primary[50],
    padding: spacing.sm,
    borderRadius: borderRadius.md,
    minWidth: 60,
  },
  
  streakNumber: {
    fontSize: typography.sizes.xl,
    fontWeight: typography.weights.bold,
    color: colors.primary[600],
  },
  
  streakLabel: {
    fontSize: typography.sizes.xs,
    color: colors.primary[500],
    fontWeight: typography.weights.medium,
  },
  
  daysGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'space-between',
    marginBottom: spacing.lg,
  },
  
  day: {
    justifyContent: 'center',
    alignItems: 'center',
    borderWidth: 2,
    margin: 2,
  },
  
  checkmark: {
    justifyContent: 'center',
    alignItems: 'center',
  },
  
  checkmarkText: {
    fontSize: typography.sizes.sm,
    fontWeight: typography.weights.bold,
  },
  
  stats: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    borderTopWidth: 1,
    borderTopColor: colors.neutral[200],
    paddingTop: spacing.md,
  },
  
  stat: {
    alignItems: 'center',
  },
  
  statLabel: {
    fontSize: typography.sizes.xs,
    color: colors.text.secondary,
    marginBottom: spacing.xs,
  },
  
  statValue: {
    fontSize: typography.sizes.md,
    fontWeight: typography.weights.semibold,
    color: colors.text.primary,
  },
})

export default HabitTracker 