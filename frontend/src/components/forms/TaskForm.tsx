import React, { useState, useEffect } from 'react'
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  ScrollView,
  Alert,
  Modal,
  KeyboardAvoidingView,
  Platform,
  Pressable,
} from 'react-native'
import { Button } from '@/src/components/ui/Button'
import { colors, spacing } from '@/src/theme'
import Select from '../Select'
import Icon from '../Icons'
import { isEmpty } from '@/src/utils/stringUtils'

// Types pour les tâches
export interface Task {
  id?: string
  title: string
  description?: string
  priority: 'low' | 'medium' | 'high'
  status: 'todo' | 'in_progress' | 'completed'
  dueDate?: Date
  category?: string
  tags?: string[]
  estimatedTime?: number // en minutes
  actualTime?: number // en minutes
  createdAt?: Date
  updatedAt?: Date
}

export interface CreateTaskRequest {
  title: string
  description?: string
  priority: 'low' | 'medium' | 'high'
  status: 'todo' | 'in_progress' | 'completed'
  dueDate?: Date
  category?: string
  tags?: string[]
  estimatedTime?: number
}

interface TaskFormProps {
  visible: boolean
  onClose: () => void
  onSubmit: (task: CreateTaskRequest) => Promise<void>
  categories?: Array<{ id: string; name: string; color: string; icon: string }>
  defaultStatus?: 'todo' | 'in_progress' | 'completed'
}

export const TaskForm: React.FC<TaskFormProps> = ({
  visible,
  onClose,
  onSubmit,
  categories = [],
  defaultStatus = 'todo',
}) => {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [priority, setPriority] = useState<'low' | 'medium' | 'high'>('medium')
  const [status, setStatus] = useState<'todo' | 'in_progress' | 'completed'>(defaultStatus)
  const [selectedCategory, setSelectedCategory] = useState<{ id: string; name: string; color: string; icon: string } | null>(null)
  const [dueDate, setDueDate] = useState<Date | undefined>(undefined)
  const [estimatedTime, setEstimatedTime] = useState('')
  const [tags, setTags] = useState<string[]>([])
  const [newTag, setNewTag] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    setErrors({})
  }, [status])

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!title.trim()) {
      newErrors.title = 'Le titre est requis'
    }

    if (estimatedTime && parseFloat(estimatedTime) <= 0) {
      newErrors.estimatedTime = 'Le temps estimé doit être supérieur à 0'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async () => {
    if (!validateForm()) {
      return
    }

    setIsLoading(true)
    try {
      const task: CreateTaskRequest = {
        title: title.trim(),
        description: description.trim() || undefined,
        priority,
        status,
        dueDate,
        category: selectedCategory?.id,
        tags: tags.length > 0 ? tags : undefined,
        estimatedTime: estimatedTime ? parseFloat(estimatedTime) : undefined,
      }

      await onSubmit(task)
      
      // Reset form
      setTitle('')
      setDescription('')
      setPriority('medium')
      setStatus('todo')
      setSelectedCategory(null)
      setDueDate(undefined)
      setEstimatedTime('')
      setTags([])
      setNewTag('')
      setErrors({})
      
      onClose()
    } catch (error) {
      Alert.alert(
        'Erreur',
        error instanceof Error ? error.message : 'Erreur lors de la création de la tâche'
      )
    } finally {
      setIsLoading(false)
    }
  }

  const addTag = () => {
    if (newTag.trim() && !tags.includes(newTag.trim())) {
      setTags([...tags, newTag.trim()])
      setNewTag('')
    }
  }

  const removeTag = (tagToRemove: string) => {
    setTags(tags.filter(tag => tag !== tagToRemove))
  }

  const formatDate = (date: Date) => {
    return date.toLocaleDateString('fr-FR')
  }

  const formatTime = (minutes: number) => {
    const hours = Math.floor(minutes / 60)
    const mins = minutes % 60
    return hours > 0 ? `${hours}h${mins > 0 ? ` ${mins}min` : ''}` : `${mins}min`
  }

  // Options pour les priorités
  const priorityOptions = [
    { label: 'Basse', value: 'low', icon: 'fa6:arrow-down', color: '#10b981' },
    { label: 'Moyenne', value: 'medium', icon: 'fa6:minus', color: '#f59e0b' },
    { label: 'Haute', value: 'high', icon: 'fa6:arrow-up', color: '#ef4444' },
  ]

  // Options pour les statuts
  const statusOptions = [
    { label: 'À faire', value: 'todo', icon: 'fa6:circle', color: '#6b7280' },
    { label: 'En cours', value: 'in_progress', icon: 'fa6:clock', color: '#3b82f6' },
    { label: 'Terminé', value: 'completed', icon: 'fa6:check-circle', color: '#10b981' },
  ]

  return (
    <Modal
      visible={visible}
      animationType="slide"
      presentationStyle="pageSheet"
      onRequestClose={onClose}
    >
      <KeyboardAvoidingView 
        style={{ flex: 1 }} 
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      >
        <View className="flex-1 bg-white">
          {/* Header */}
          <View className="flex-row items-center justify-between p-4 border-b border-gray-200">
            <TouchableOpacity onPress={onClose}>
              <Icon name="close" size={24} color={colors.text.primary} />
            </TouchableOpacity>
            <Text className="text-lg font-semibold text-gray-900">
              Nouvelle Tâche
            </Text>
            <View style={{ width: 24 }} />
          </View>

          <ScrollView className="flex-1 px-4 py-6">
            {/* Titre */}
            <View className="mb-6">
              <Text className="text-sm font-medium mb-2 text-gray-900">
                Titre *
              </Text>
              <View className={`border rounded-lg px-4 py-3 ${
                errors.title ? 'border-red-500' : 'border-gray-300'
              } bg-gray-50`}>
                <TextInput
                  className="text-base text-gray-900"
                  placeholder="Titre de la tâche..."
                  placeholderTextColor={colors.text.tertiary}
                  value={title}
                  onChangeText={(text) => {
                    setTitle(text)
                    if (errors.title) {
                      setErrors(prev => ({ ...prev, title: '' }))
                    }
                  }}
                />
              </View>
              {errors.title && (
                <Text className="text-sm mt-1 text-red-500">{errors.title}</Text>
              )}
            </View>

            {/* Statut */}
            <View className="mb-6">
              <Text className="text-sm font-medium mb-3 text-gray-900">
                Statut
              </Text>
              <View className="flex-row space-x-3">
                {statusOptions.map((statusOption) => (
                  <TouchableOpacity
                    key={statusOption.value}
                    className={`flex-1 p-3 rounded-lg border-2 ${
                      status === statusOption.value 
                        ? `border-${statusOption.color} bg-${statusOption.color}10` 
                        : 'border-gray-300'
                    }`}
                    onPress={() => setStatus(statusOption.value as any)}
                  >
                    <View className="items-center">
                      <Icon 
                        name={statusOption.icon} 
                        size={24} 
                        color={status === statusOption.value ? statusOption.color : colors.text.tertiary} 
                      />
                      <Text className={`text-sm font-medium mt-1 ${
                        status === statusOption.value ? `text-${statusOption.color}` : 'text-gray-600'
                      }`}>
                        {statusOption.label}
                      </Text>
                    </View>
                  </TouchableOpacity>
                ))}
              </View>
            </View>

            {/* Priorité */}
            <View className="mb-6">
              <Text className="text-sm font-medium mb-3 text-gray-900">
                Priorité
              </Text>
              <View className="flex-row space-x-3">
                {priorityOptions.map((priorityOption) => (
                  <TouchableOpacity
                    key={priorityOption.value}
                    className={`flex-1 p-3 rounded-lg border-2 ${
                      priority === priorityOption.value 
                        ? `border-${priorityOption.color} bg-${priorityOption.color}10` 
                        : 'border-gray-300'
                    }`}
                    onPress={() => setPriority(priorityOption.value as any)}
                  >
                    <View className="items-center">
                      <Icon 
                        name={priorityOption.icon} 
                        size={24} 
                        color={priority === priorityOption.value ? priorityOption.color : colors.text.tertiary} 
                      />
                      <Text className={`text-sm font-medium mt-1 ${
                        priority === priorityOption.value ? `text-${priorityOption.color}` : 'text-gray-600'
                      }`}>
                        {priorityOption.label}
                      </Text>
                    </View>
                  </TouchableOpacity>
                ))}
              </View>
            </View>

            {/* Catégorie */}
            {categories.length > 0 && (
              <View className="mb-6">
                <Text className="text-sm font-medium mb-2 text-gray-900">
                  Catégorie (optionnel)
                </Text>
                <ScrollView horizontal showsHorizontalScrollIndicator={false}>
                  <View className="flex-row gap-2">
                    {categories.map(category => (
                      <Pressable
                        key={category.id}
                        className={`flex-row items-center gap-2 py-4 px-4 rounded-xl overflow-hidden`}
                        style={{
                          backgroundColor:
                            selectedCategory?.id === category.id
                              ? `#${category.color.replace('#', '')}`
                              : `#${category.color.replace('#', '')}87`,
                          borderRadius: 10,
                          borderWidth: 1,
                          borderColor:
                            selectedCategory?.id === category.id
                              ? `#${category.color.replace('#', '')}`
                              : `#${category.color.replace('#', '')}87`,
                        }}
                        onPress={() => {
                          if (selectedCategory?.id === category.id) {
                            setSelectedCategory(null);
                          } else {
                            setSelectedCategory(category);
                          }
                        }}
                      >
                        <View className="flex-row items-center gap-2">
                          <Icon 
                            name={category.icon as any} 
                            size={20} 
                            color={selectedCategory?.id === category.id ? '#fff' : category.color} 
                          />
                          <Text className="text-white">{category.name}</Text>
                        </View>
                      </Pressable>
                    ))}
                  </View>
                </ScrollView>
              </View>
            )}

            {/* Date d'échéance */}
            <View className="mb-6">
              <Text className="text-sm font-medium mb-2 text-gray-900">
                Date d'échéance (optionnel)
              </Text>
              <View className="flex-row items-center justify-between border border-gray-300 rounded-lg px-4 py-3 bg-gray-50">
                <View className="flex-row items-center">
                  <Icon 
                    name="calendar-outline" 
                    size={20} 
                    color={colors.text.tertiary} 
                    style={{ marginRight: spacing.sm }}
                  />
                  <Text className="text-base text-gray-900">
                    {dueDate ? formatDate(dueDate) : 'Aucune date définie'}
                  </Text>
                </View>
                <TouchableOpacity onPress={() => setDueDate(dueDate ? undefined : new Date())}>
                  <Icon 
                    name={dueDate ? "close-circle" : "add-circle"} 
                    size={24} 
                    color={colors.text.tertiary} 
                  />
                </TouchableOpacity>
              </View>
            </View>

            {/* Temps estimé */}
            <View className="mb-6">
              <Text className="text-sm font-medium mb-2 text-gray-900">
                Temps estimé (optionnel)
              </Text>
              <View className={`flex-row items-center border rounded-lg px-4 py-3 ${
                errors.estimatedTime ? 'border-red-500' : 'border-gray-300'
              } bg-gray-50`}>
                <Icon 
                  name="time-outline" 
                  size={20} 
                  color={colors.text.tertiary} 
                  style={{ marginRight: spacing.sm }}
                />
                <TextInput
                  className="flex-1 text-base text-gray-900"
                  placeholder="0"
                  placeholderTextColor={colors.text.tertiary}
                  value={estimatedTime}
                  onChangeText={(text) => {
                    setEstimatedTime(text.replace(/[^0-9]/g, ''))
                    if (errors.estimatedTime) {
                      setErrors(prev => ({ ...prev, estimatedTime: '' }))
                    }
                  }}
                  keyboardType="numeric"
                />
                <Text className="text-base ml-2 text-gray-900">minutes</Text>
              </View>
              {errors.estimatedTime && (
                <Text className="text-sm mt-1 text-red-500">{errors.estimatedTime}</Text>
              )}
            </View>

            {/* Tags */}
            <View className="mb-6">
              <Text className="text-sm font-medium mb-2 text-gray-900">
                Tags (optionnel)
              </Text>
              <View className="flex-row items-center border border-gray-300 rounded-lg px-4 py-3 bg-gray-50">
                <Icon 
                  name="pricetag-outline" 
                  size={20} 
                  color={colors.text.tertiary} 
                  style={{ marginRight: spacing.sm }}
                />
                <TextInput
                  className="flex-1 text-base text-gray-900"
                  placeholder="Ajouter un tag..."
                  placeholderTextColor={colors.text.tertiary}
                  value={newTag}
                  onChangeText={setNewTag}
                  onSubmitEditing={addTag}
                />
                <TouchableOpacity onPress={addTag} disabled={!newTag.trim()}>
                  <Icon 
                    name="add-circle" 
                    size={24} 
                    color={newTag.trim() ? colors.primary[500] : colors.text.tertiary} 
                  />
                </TouchableOpacity>
              </View>
              {tags.length > 0 && (
                <View className="flex-row flex-wrap gap-2 mt-3">
                  {tags.map((tag, index) => (
                    <View key={index} className="flex-row items-center bg-blue-100 rounded-full px-3 py-1">
                      <Text className="text-sm text-blue-800 mr-2">{tag}</Text>
                      <TouchableOpacity onPress={() => removeTag(tag)}>
                        <Icon name="close-circle" size={16} color="#1e40af" />
                      </TouchableOpacity>
                    </View>
                  ))}
                </View>
              )}
            </View>

            {/* Description */}
            <View className="mb-6">
              <Text className="text-sm font-medium mb-2 text-gray-900">
                Description (optionnel)
              </Text>
              <View className="border border-gray-300 rounded-lg px-4 py-3 bg-gray-50">
                <TextInput
                  className="text-base text-gray-900"
                  placeholder="Description de la tâche..."
                  placeholderTextColor={colors.text.tertiary}
                  value={description}
                  onChangeText={setDescription}
                  multiline
                  numberOfLines={4}
                  textAlignVertical="top"
                />
              </View>
            </View>
          </ScrollView>

          {/* Footer */}
          <View className="p-4 border-t border-gray-200">
            <Button
              title="Créer la tâche"
              onPress={handleSubmit}
              disabled={isLoading}
              loading={isLoading}
              fullWidth
              size="large"
              variant="primary"
            />
          </View>
        </View>
      </KeyboardAvoidingView>
    </Modal>
  )
} 