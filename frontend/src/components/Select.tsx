import React, { useState } from 'react';
import { View, Text, TouchableOpacity, ScrollView, Modal, StyleSheet, TextInput } from 'react-native';
import { Icon, FONTAWESOME_6 } from '@/src/components/Icons';

// Types pour les options du Select
export interface SelectOption<T = string> {
  label: string;
  value: T;
  disabled?: boolean;
  icon?: string;
  description?: string;
}

// Props du composant Select
export interface SelectProps<T = string> {
  // Options et valeur
  options: SelectOption<T>[];
  value?: T;
  onValueChange?: (value: T) => void;
  
  // Configuration
  placeholder?: string;
  disabled?: boolean;
  searchable?: boolean;
  multiple?: boolean;
  
  // Apparence
  size?: 'sm' | 'md' | 'lg';
  variant?: 'default' | 'outline' | 'filled';
  width?: number | string;
  
  // Labels et descriptions
  label?: string;
  error?: string;
  helperText?: string;
  
  // Callbacks
  onOpen?: () => void;
  onClose?: () => void;
  
  // Styles personnalisés (pour compatibilité future)
  // className?: string;
  // triggerClassName?: string;
  // contentClassName?: string;
  // itemClassName?: string;
}

// Composant Select principal
export function Select<T = string>({
  options,
  value,
  onValueChange,
  placeholder = "Sélectionner une option",
  disabled = false,
  searchable = false,
  multiple = false,
  size = 'md',
  variant = 'default',
  width = '100%',
  label,
  error,
  helperText,
  onOpen,
  onClose,
}: SelectProps<T>) {
  const [isOpen, setIsOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');

  // Trouver l'option sélectionnée
  const selectedOption = options.find(option => option.value === value);
  
  // Filtrer les options si searchable est activé
  const filteredOptions = searchable 
    ? options.filter(option => 
        option.label.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : options;

  // Gérer l'ouverture/fermeture
  const handleOpen = () => {
    if (!disabled) {
      setIsOpen(true);
      onOpen?.();
    }
  };

  const handleClose = () => {
    setIsOpen(false);
    setSearchQuery('');
    onClose?.();
  };

  // Gérer la sélection
  const handleSelect = (option: SelectOption<T>) => {
    if (!option.disabled) {
      onValueChange?.(option.value);
      if (!multiple) {
        handleClose();
      }
    }
  };

  // Styles dynamiques
  const getTriggerStyles = () => {
    const baseStyles = {
      width: width as any,
      borderWidth: 1,
      borderRadius: 8,
      paddingHorizontal: 16,
      flexDirection: 'row' as const,
      alignItems: 'center' as const,
      justifyContent: 'space-between' as const,
    };

    // Taille
    const sizeStyles = {
      sm: { height: 32, paddingVertical: 6 },
      md: { height: 40, paddingVertical: 8 },
      lg: { height: 48, paddingVertical: 12 },
    };

    // Variant
    const variantStyles = {
      default: {
        backgroundColor: '#ffffff',
        borderColor: '#d1d5db',
      },
      outline: {
        backgroundColor: 'transparent',
        borderColor: '#3b82f6',
      },
      filled: {
        backgroundColor: '#f3f4f6',
        borderColor: 'transparent',
      },
    };

    return {
      ...baseStyles,
      ...sizeStyles[size],
      ...variantStyles[variant],
      ...(disabled && { opacity: 0.5 }),
      ...(error && { borderColor: '#ef4444' }),
    };
  };

  const getContentStyles = () => ({
    backgroundColor: '#ffffff',
    borderRadius: 8,
    borderWidth: 1,
    borderColor: '#d1d5db',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 5,
    maxHeight: 300,
  });

  const getItemStyles = (isSelected: boolean, isDisabled: boolean) => ({
    paddingVertical: 12,
    paddingHorizontal: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#f3f4f6',
    flexDirection: 'row' as const,
    alignItems: 'center' as const,
    backgroundColor: isSelected ? '#eff6ff' : 'transparent',
    opacity: isDisabled ? 0.5 : 1,
  });

  return (
    <View style={styles.container}>
      {/* Label */}
      {label && (
        <Text style={styles.label}>{label}</Text>
      )}

      {/* Trigger */}
      <TouchableOpacity
        style={getTriggerStyles()}
        onPress={handleOpen}
        disabled={disabled}
        activeOpacity={0.7}
      >
        <View style={styles.triggerContent}>
          {selectedOption?.icon && (
            <Icon 
              name={selectedOption.icon} 
              size={16} 
              color="#6b7280" 
              style={styles.triggerIcon}
            />
          )}
          <Text style={[
            styles.triggerText,
            !selectedOption && styles.placeholderText
          ]}>
            {selectedOption?.label || placeholder}
          </Text>
        </View>
        <Icon 
          name={`fa6:${FONTAWESOME_6.CHEVRON_DOWN}`} 
          size={16} 
          color="#6b7280" 
        />
      </TouchableOpacity>

      {/* Error */}
      {error && (
        <Text style={styles.errorText}>{error}</Text>
      )}

      {/* Helper text */}
      {helperText && !error && (
        <Text style={styles.helperText}>{helperText}</Text>
      )}

      {/* Modal */}
      <Modal
        visible={isOpen}
        transparent
        animationType="fade"
        onRequestClose={handleClose}
      >
        <TouchableOpacity
          style={styles.overlay}
          activeOpacity={1}
          onPress={handleClose}
        >
          <View style={getContentStyles()}>
            {/* Header */}
            <View style={styles.header}>
              <Text style={styles.headerTitle}>
                {label || 'Sélectionner'}
              </Text>
              <TouchableOpacity onPress={handleClose} style={styles.closeButton}>
                <Icon 
                  name={`fa6:${FONTAWESOME_6.XMARK}`} 
                  size={20} 
                  color="#6b7280" 
                />
              </TouchableOpacity>
            </View>

            {/* Search (si activé) */}
            {searchable && (
              <View style={styles.searchContainer}>
                <Icon 
                  name={`fa6:${FONTAWESOME_6.MAGNIFYING_GLASS}`} 
                  size={16} 
                  color="#6b7280" 
                  style={styles.searchIcon}
                />
                <TextInput
                  style={styles.searchInput}
                  placeholder="Rechercher..."
                  value={searchQuery}
                  onChangeText={setSearchQuery}
                  autoFocus
                />
              </View>
            )}

            {/* Options */}
            <ScrollView style={styles.optionsContainer}>
              {filteredOptions.map((option, index) => {
                const isSelected = option.value === value;
                const isLast = index === filteredOptions.length - 1;

                return (
                  <TouchableOpacity
                    key={`${option.value}-${index}`}
                    style={[
                      getItemStyles(isSelected, !!option.disabled),
                      isLast && styles.lastItem,
                    ]}
                    onPress={() => handleSelect(option)}
                    disabled={option.disabled}
                    activeOpacity={0.7}
                  >
                    {/* Option icon */}
                    {option.icon && (
                      <Icon 
                        name={option.icon} 
                        size={16} 
                        color="#6b7280" 
                        style={styles.optionIcon}
                      />
                    )}

                    {/* Option content */}
                    <View style={styles.optionContent}>
                      <Text style={[
                        styles.optionLabel,
                        isSelected && styles.selectedOptionLabel
                      ]}>
                        {option.label}
                      </Text>
                      {option.description && (
                        <Text style={styles.optionDescription}>
                          {option.description}
                        </Text>
                      )}
                    </View>

                    {/* Selection indicator */}
                    {isSelected && (
                      <Icon 
                        name={`fa6:${FONTAWESOME_6.CHECK}`} 
                        size={16} 
                        color="#3b82f6" 
                        style={styles.checkIcon}
                      />
                    )}
                  </TouchableOpacity>
                );
              })}
            </ScrollView>

            {/* Empty state */}
            {filteredOptions.length === 0 && (
              <View style={styles.emptyState}>
                <Text style={styles.emptyStateText}>
                  {searchable ? 'Aucun résultat trouvé' : 'Aucune option disponible'}
                </Text>
              </View>
            )}
          </View>
        </TouchableOpacity>
      </Modal>
    </View>
  );
}

// Styles
const styles = StyleSheet.create({
  container: {
    width: '100%',
  },
  label: {
    fontSize: 14,
    fontWeight: '500',
    color: '#374151',
    marginBottom: 4,
  },
  triggerContent: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
  },
  triggerIcon: {
    marginRight: 8,
  },
  triggerText: {
    fontSize: 14,
    color: '#111827',
    flex: 1,
  },
  placeholderText: {
    color: '#9ca3af',
  },
  errorText: {
    fontSize: 12,
    color: '#ef4444',
    marginTop: 4,
  },
  helperText: {
    fontSize: 12,
    color: '#6b7280',
    marginTop: 4,
  },
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#f3f4f6',
  },
  headerTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#111827',
  },
  closeButton: {
    padding: 4,
  },
  searchContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#f3f4f6',
  },
  searchIcon: {
    marginRight: 8,
  },
  searchInput: {
    flex: 1,
    fontSize: 14,
    color: '#111827',
    padding: 0,
  },
  optionsContainer: {
    maxHeight: 200,
  },
  optionIcon: {
    marginRight: 12,
  },
  optionContent: {
    flex: 1,
  },
  optionLabel: {
    fontSize: 14,
    color: '#111827',
  },
  selectedOptionLabel: {
    fontWeight: '500',
    color: '#3b82f6',
  },
  optionDescription: {
    fontSize: 12,
    color: '#6b7280',
    marginTop: 2,
  },
  checkIcon: {
    marginLeft: 8,
  },
  lastItem: {
    borderBottomWidth: 0,
  },
  emptyState: {
    padding: 20,
    alignItems: 'center',
  },
  emptyStateText: {
    fontSize: 14,
    color: '#6b7280',
    textAlign: 'center',
  },
});

// Export par défaut
export default Select;
