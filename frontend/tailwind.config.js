/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./app/**/*.{js,ts,tsx}', './components/**/*.{js,ts,tsx}'],

  presets: [require('nativewind/preset')],
  theme: {
    extend: {
      colors: {
        primary: '#319795', // Teal-500 as primary
        secondary: '#805AD5', // Purple as secondary
        accent: '#F6AD55', // Orange as accent
        success: '#48BB78', // Green
        warning: '#ECC94B', // Yellow
        error: '#FF4D4F', // Red
        info: '#4299E1', // Blue
        
        // Neutral colors
        neutral: {
          50: '#F7FAFC',
          100: '#EDF2F7',
          200: '#E2E8F0',
          300: '#CBD5E0',
          400: '#A0AEC0',
          500: '#718096',
          600: '#4A5568',
          700: '#2D3748',
          800: '#1A202C',
          900: '#171923',
          950: '#0D1117'
        },
        
        // Original red color
        
        // Teal color palette
        teal: {
          50: '#E6FFFA',
          100: '#B2F5EA',
          200: '#81E6D9',
          300: '#4FD1C5',
          400: '#38B2AC',
          500: '#319795',
          DEFAULT: '#2C7A7B',
          700: '#285E61',
          800: '#234E52',
          900: '#1D4044',
          950: '#0F2C2D'
        },
        
        // Purple color palette
        purple: {
          50: '#FAF5FF',
          100: '#E9D8FD',
          200: '#D6BCFA',
          300: '#B794F4',
          DEFAULT: '#9F7AEA',
          500: '#805AD5',
          600: '#6B46C1',
          700: '#553C9A',
          800: '#44337A',
          900: '#322659',
          950: '#1F1637'
        },
        
        // Orange color palette
        orange: {
          50: '#FFFAF0',
          100: '#FEEBC8',
          200: '#FBD38D',
          300: '#F6AD55',
          DEFAULT: '#ED8936',
          500: '#DD6B20',
          600: '#C05621',
          700: '#9C4221',
          800: '#7B341E',
          900: '#652B19',
          950: '#451C11'
        },
      },
    },
  },
  plugins: [],
};
