import React, { useState, useEffect } from 'react'
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  Alert,
  ActivityIndicator,
} from 'react-native'
import { router } from 'expo-router'
import { Ionicons } from '@expo/vector-icons'
import { useAuthStore } from '@/src/stores/authStore'
import { Button } from '@/src/components/ui/Button'
import { colors, spacing, typography, shadows } from '@/src/theme'

export default function Login() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [emailError, setEmailError] = useState('')
  const [passwordError, setPasswordError] = useState('')

  const { login, isLoading, error, clearError, user, require_preferences } = useAuthStore()

  // Clear error when component mounts
  useEffect(() => {
    clearError()
  }, [])

  const validateEmail = (email: string) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!email) {
      return 'L\'email est requis'
    }
    if (!emailRegex.test(email)) {
      return 'Format d\'email invalide'
    }
    return ''
  }

  const validatePassword = (password: string) => {
    if (!password) {
      return 'Le mot de passe est requis'
    }
    if (password.length < 6) {
      return 'Le mot de passe doit contenir au moins 6 caractères'
    }
    return ''
  }

  const handleEmailChange = (text: string) => {
    setEmail(text)
    if (emailError) {
      setEmailError(validateEmail(text))
    }
  }

  const handlePasswordChange = (text: string) => {
    setPassword(text)
    if (passwordError) {
      setPasswordError(validatePassword(text))
    }
  }

  const validateForm = () => {
    const emailValidation = validateEmail(email)
    const passwordValidation = validatePassword(password)

    setEmailError(emailValidation)
    setPasswordError(passwordValidation)

    return !emailValidation && !passwordValidation
  }

  const handleLogin = async () => {
    if (!validateForm()) {
      return
    }

    try {
      const { isAuthenticated, require_preferences } = await login(email, password)
      console.log(isAuthenticated, "isAuthenticated");
      console.log(require_preferences, "require_preferences");
      
      if (isAuthenticated) {
       
        if (require_preferences) {
          router.replace('/dashboard/wizard')
        } else {
          router.replace('/dashboard/tabs/overview')
        }
      } else {
        Alert.alert('Erreur', 'Email ou mot de passe incorrect')
      }
      // Navigation will be handled by the auth state change
    } catch (error) {
      // Error is handled by the store
    }
  }

  const handleRegister = () => {
    router.push('/auth/register')
  }

  const handleForgotPassword = () => {
    Alert.alert(
      'Mot de passe oublié',
      'Cette fonctionnalité sera bientôt disponible.',
      [{ text: 'OK' }]
    )
  }

  return (
    <KeyboardAvoidingView
      style={{ flex: 1 }}
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
    >
      <ScrollView
        contentContainerStyle={{ flexGrow: 1 }}
        keyboardShouldPersistTaps="handled"
      >
        <View className="flex-1 bg-white px-6">
          {/* Header */}
          <View className="flex-1 justify-center">
            <View className="items-center mb-8">
              <View
                className="w-20 h-20 rounded-full items-center justify-center mb-4"
                style={{
                  backgroundColor: colors.primary[100],
                  ...shadows.md
                }}
              >
                <Ionicons
                  name="wallet-outline"
                  size={40}
                  color={colors.primary[500]}
                />
              </View>
              <Text
                className="text-3xl font-bold mb-2"
                style={{ color: colors.text.primary }}
              >
                MeshaPlus
              </Text>
              <Text
                className="text-base text-center"
                style={{ color: colors.text.secondary }}
              >
                Connectez-vous à votre compte
              </Text>
            </View>

            {/* Form */}
            <View className="space-y-4">
              {/* Email Input */}
              <View>
                <Text
                  className="text-sm font-medium mb-2"
                  style={{ color: colors.text.primary }}
                >
                  Email
                </Text>
                <View
                  className={`flex-row items-center border rounded-lg px-4 py-3 ${emailError ? 'border-red-500' : 'border-gray-300'
                    }`}
                  style={{ backgroundColor: colors.background.secondary }}
                >
                  <Ionicons
                    name="mail-outline"
                    size={20}
                    color={colors.text.tertiary}
                    style={{ marginRight: spacing.sm }}
                  />
                  <TextInput
                    className="flex-1 text-base"
                    placeholder="votre@email.com"
                    placeholderTextColor={colors.text.tertiary}
                    value={email}
                    onChangeText={handleEmailChange}
                    keyboardType="email-address"
                    autoCapitalize="none"
                    autoCorrect={false}
                    style={{ color: colors.text.primary }}
                  />
                </View>
                {emailError ? (
                  <Text
                    className="text-sm mt-1"
                    style={{ color: colors.error[500] }}
                  >
                    {emailError}
                  </Text>
                ) : null}
              </View>

              {/* Password Input */}
              <View>
                <Text
                  className="text-sm font-medium mb-2"
                  style={{ color: colors.text.primary }}
                >
                  Mot de passe
                </Text>
                <View
                  className={`flex-row items-center border rounded-lg px-4 py-3 ${passwordError ? 'border-red-500' : 'border-gray-300'
                    }`}
                  style={{ backgroundColor: colors.background.secondary }}
                >
                  <Ionicons
                    name="lock-closed-outline"
                    size={20}
                    color={colors.text.tertiary}
                    style={{ marginRight: spacing.sm }}
                  />
                  <TextInput
                    className="flex-1 text-base"
                    placeholder="Votre mot de passe"
                    placeholderTextColor={colors.text.tertiary}
                    value={password}
                    onChangeText={handlePasswordChange}
                    secureTextEntry={!showPassword}
                    autoCapitalize="none"
                    autoCorrect={false}
                    style={{ color: colors.text.primary }}
                  />
                  <TouchableOpacity
                    onPress={() => setShowPassword(!showPassword)}
                    hitSlop={{ top: 10, bottom: 10, left: 10, right: 10 }}
                  >
                    <Ionicons
                      name={showPassword ? "eye-off-outline" : "eye-outline"}
                      size={20}
                      color={colors.text.tertiary}
                    />
                  </TouchableOpacity>
                </View>
                {passwordError ? (
                  <Text
                    className="text-sm mt-1"
                    style={{ color: colors.error[500] }}
                  >
                    {passwordError}
                  </Text>
                ) : null}
              </View>

              {/* Forgot Password */}
              <TouchableOpacity
                onPress={handleForgotPassword}
                className="self-end"
              >
                <Text
                  className="text-sm"
                  style={{ color: colors.primary[500] }}
                >
                  Mot de passe oublié ?
                </Text>
              </TouchableOpacity>

              {/* Error Message */}
              {error ? (
                <View
                  className="p-3 rounded-lg"
                  style={{ backgroundColor: colors.error[50] }}
                >
                  <Text
                    className="text-sm text-center"
                    style={{ color: colors.error[600] }}
                  >
                    {error}
                  </Text>
                </View>
              ) : null}

              {/* Login Button */}
              <Button
                title="Se connecter"
                onPress={handleLogin}
                disabled={isLoading}
                loading={isLoading}
                fullWidth
                size="large"
                style={{ marginTop: spacing.md }}
              />

              {/* Divider */}
              <View className="flex-row items-center my-6">
                <View
                  className="flex-1 h-px"
                  style={{ backgroundColor: colors.neutral[200] }}
                />
                <Text
                  className="mx-4 text-sm"
                  style={{ color: colors.text.tertiary }}
                >
                  ou
                </Text>
                <View
                  className="flex-1 h-px"
                  style={{ backgroundColor: colors.neutral[200] }}
                />
              </View>

              {/* Register Link */}
              <View className="items-center">
                <Text
                  className="text-sm"
                  style={{ color: colors.text.secondary }}
                >
                  Pas encore de compte ?{' '}
                </Text>
                <TouchableOpacity onPress={handleRegister}>
                  <Text
                    className="text-sm font-medium"
                    style={{ color: colors.primary[500] }}
                  >
                    Créer un compte
                  </Text>
                </TouchableOpacity>
              </View>
            </View>
          </View>

          {/* Footer */}
          <View className="py-6">
            <Text
              className="text-xs text-center"
              style={{ color: colors.text.tertiary }}
            >
              En vous connectant, vous acceptez nos conditions d'utilisation
            </Text>
          </View>
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  )
}