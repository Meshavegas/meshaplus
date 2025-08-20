import { colors, spacing, typography } from '@/src/theme'
import React from 'react'
import { View, Text, Dimensions } from 'react-native'
// Temporairement commenté jusqu'à l'installation des dépendances
// import { VictoryPie, VictoryChart, VictoryLine, VictoryArea, VictoryAxis, VictoryTooltip } from 'victory-native'
// import { Svg, Circle, G } from 'react-native-svg'
import Animated, { 
  useSharedValue, 
  useAnimatedStyle, 
  withSpring,
  withTiming,
  withRepeat
} from 'react-native-reanimated'
// import { colors, spacing, typography } from '~/theme'

const { width: screenWidth } = Dimensions.get('window')

interface ExpenseData {
  id: string
  category: string
  amount: number
  color: string
  percentage: number
}

interface ExpenseChartProps {
  data: ExpenseData[]
  type?: 'pie' | 'line' | 'area'
  size?: 'small' | 'medium' | 'large'
  showLabels?: boolean
  animated?: boolean
}

export const ExpenseChart: React.FC<ExpenseChartProps> = ({
  data,
  type = 'pie',
  size = 'medium',
  showLabels = true,
  animated = true,
}) => {
  const animationValue = useSharedValue(0)
  const pulseValue = useSharedValue(1)

  React.useEffect(() => {
    if (animated) {
      animationValue.value = withTiming(1, { duration: 1000 })
      pulseValue.value = withRepeat(
        withTiming(1.05, { duration: 1000 }),
        -1,
        true
      )
    }
  }, [animated])

  const animatedStyle = useAnimatedStyle(() => ({
    transform: [{ scale: pulseValue.value }],
  }))

  const getChartSize = () => {
    switch (size) {
      case 'small':
        return screenWidth * 0.6
      case 'large':
        return screenWidth * 0.9
      default:
        return screenWidth * 0.75
    }
  }

  const chartSize = getChartSize()

  // Version simplifiée sans Victory Native pour l'instant
  const renderSimpleChart = () => (
    <Animated.View style={[animatedStyle]}>
      <View style={{
        width: chartSize,
        height: chartSize,
        backgroundColor: colors.background.secondary,
        borderRadius: 12,
        justifyContent: 'center',
        alignItems: 'center',
      }}>
        <Text style={{
          fontSize: typography.sizes.lg,
          color: colors.text.secondary,
          textAlign: 'center',
        }}>
          Graphique {type}
        </Text>
        <Text style={{
          fontSize: typography.sizes.sm,
          color: colors.text.tertiary,
          marginTop: spacing.sm,
        }}>
          {data.length} données
        </Text>
      </View>
    </Animated.View>
  )

  const renderLegend = () => (
    <View style={{ marginTop: spacing.md }}>
      {data.map((item, index) => (
        <View key={item.id} style={{ 
          flexDirection: 'row', 
          alignItems: 'center', 
          marginBottom: spacing.sm 
        }}>
          <View style={{
            width: 12,
            height: 12,
            borderRadius: 6,
            backgroundColor: item.color,
            marginRight: spacing.sm,
          }} />
          <Text style={{
            fontSize: typography.sizes.sm,
            color: colors.text.secondary,
            flex: 1,
          }}>
            {item.category}
          </Text>
          <Text style={{
            fontSize: typography.sizes.sm,
            fontWeight: typography.weights.semibold,
            color: colors.text.primary,
          }}>
            ${item.amount.toFixed(2)}
          </Text>
        </View>
      ))}
    </View>
  )

  return (
    <View style={{ alignItems: 'center' }}>
      {renderSimpleChart()}
      {showLabels && renderLegend()}
    </View>
  )
}

export default ExpenseChart 