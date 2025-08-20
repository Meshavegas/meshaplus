
import React from "react";
import { View } from "react-native";
import { Svg, Circle, Text as SVGText } from 'react-native-svg';

interface CircularProgressProps{
    size?:number,
    strokeWidth?:number, 
    text?:string, 
    progressPercent?:number, 
    bgColor?:string,
    pgColor?:string,
    textSize?:number,
    textColor?:string
}
const CircularProgress = (props:CircularProgressProps) => {
  const { size=12, strokeWidth=2, text="", progressPercent=50, bgColor, pgColor, textSize, textColor } = props;
  const radius = (size - strokeWidth) / 2;
  const circum = radius * 2 * Math.PI;
  const svgProgress = 100 - progressPercent;

  return (
    <View style={{margin: 10}}>
      <Svg width={size} height={size}>
        {/* Background Circle */}
        <Circle 
          stroke={bgColor ? bgColor : "#f2f2f2"}
          fill="none"
          cx={size / 2}
          cy={size / 2}
          r={radius}
          {...{strokeWidth}}
        />
        
        {/* Progress Circle */}
        <Circle 
          stroke={pgColor ? pgColor : "#3b5998"}
          fill="none"
          cx={size / 2}
          cy={size / 2}
          r={radius}
          strokeDasharray={`${circum} ${circum}`}
          strokeDashoffset={radius * Math.PI * 2 * (svgProgress/100)}
          strokeLinecap="round"
          transform={`rotate(-90, ${size/2}, ${size/2})`}
          {...{strokeWidth}}
        />

        {/* Text */}
        <SVGText
          fontSize={textSize ? textSize : "10"}
          x={size / 2}
          y={size / 2 + (textSize ?  (textSize / 2) - 1 : 5)}
          textAnchor="middle"
          fill={textColor ? textColor : "#333333"}
        >
          {text}
        </SVGText>
      </Svg>
    </View>
  )
}

export default CircularProgress;