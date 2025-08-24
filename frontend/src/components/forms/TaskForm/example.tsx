import React, { useState } from 'react';
import { View, Text, TouchableOpacity, ScrollView } from 'react-native';
import { TaskForm, CreateTaskRequest } from '../TaskForm';

// Exemple d'utilisation du composant TaskForm
export default function TaskFormExample() {
  const [showTaskForm, setShowTaskForm] = useState(false);
  const [tasks, setTasks] = useState<CreateTaskRequest[]>([]);

  // Catégories de tâches
  const taskCategories = [
    { id: 'work', name: 'Travail', color: '#3b82f6', icon: 'fa6:briefcase' },
    { id: 'personal', name: 'Personnel', color: '#10b981', icon: 'fa6:user' },
    { id: 'health', name: 'Santé', color: '#ef4444', icon: 'fa6:heart-pulse' },
    { id: 'finance', name: 'Finance', color: '#f59e0b', icon: 'fa6:wallet' },
    { id: 'education', name: 'Éducation', color: '#8b5cf6', icon: 'fa6:graduation-cap' },
    { id: 'home', name: 'Maison', color: '#06b6d4', icon: 'fa6:house' },
  ];

  const handleCreateTask = async (task: CreateTaskRequest) => {
    // Simuler une création de tâche
    console.log('Nouvelle tâche créée:', task);
    
    // Ajouter la tâche à la liste
    setTasks(prev => [...prev, task]);
    
    // Ici vous appelleriez votre API pour sauvegarder la tâche
    // await taskService.createTask(task);
  };

  const formatDate = (date?: Date) => {
    if (!date) return 'Aucune date';
    return date.toLocaleDateString('fr-FR');
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'high': return '#ef4444';
      case 'medium': return '#f59e0b';
      case 'low': return '#10b981';
      default: return '#6b7280';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed': return '#10b981';
      case 'in_progress': return '#3b82f6';
      case 'todo': return '#6b7280';
      default: return '#6b7280';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'completed': return 'Terminé';
      case 'in_progress': return 'En cours';
      case 'todo': return 'À faire';
      default: return 'À faire';
    }
  };

  return (
    <ScrollView style={{ flex: 1, padding: 20, backgroundColor: '#f5f5f5' }}>
      <Text style={{ fontSize: 24, fontWeight: 'bold', marginBottom: 20 }}>
        Gestionnaire de Tâches
      </Text>

      {/* Bouton pour ouvrir le formulaire */}
      <TouchableOpacity
        style={{
          backgroundColor: '#3b82f6',
          padding: 16,
          borderRadius: 12,
          alignItems: 'center',
          marginBottom: 20,
        }}
        onPress={() => setShowTaskForm(true)}
      >
        <Text style={{ color: 'white', fontSize: 16, fontWeight: '600' }}>
          + Nouvelle Tâche
        </Text>
      </TouchableOpacity>

      {/* Liste des tâches */}
      <View style={{ marginBottom: 20 }}>
        <Text style={{ fontSize: 18, fontWeight: '600', marginBottom: 12 }}>
          Mes Tâches ({tasks.length})
        </Text>
        
        {tasks.length === 0 ? (
          <View style={{
            backgroundColor: 'white',
            padding: 20,
            borderRadius: 12,
            alignItems: 'center',
          }}>
            <Text style={{ color: '#6b7280', textAlign: 'center' }}>
              Aucune tâche créée pour le moment.{'\n'}
              Créez votre première tâche !
            </Text>
          </View>
        ) : (
          <View style={{ gap: 12 }}>
            {tasks.map((task, index) => (
              <View
                key={index}
                style={{
                  backgroundColor: 'white',
                  padding: 16,
                  borderRadius: 12,
                  borderLeftWidth: 4,
                  borderLeftColor: getPriorityColor(task.priority),
                }}
              >
                {/* En-tête de la tâche */}
                <View style={{ flexDirection: 'row', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: 8 }}>
                  <Text style={{ fontSize: 16, fontWeight: '600', flex: 1, marginRight: 8 }}>
                    {task.title}
                  </Text>
                  <View style={{
                    backgroundColor: getStatusColor(task.status),
                    paddingHorizontal: 8,
                    paddingVertical: 4,
                    borderRadius: 12,
                  }}>
                    <Text style={{ color: 'white', fontSize: 12, fontWeight: '500' }}>
                      {getStatusText(task.status)}
                    </Text>
                  </View>
                </View>

                {/* Description */}
                {task.description && (
                  <Text style={{ color: '#6b7280', marginBottom: 8 }}>
                    {task.description}
                  </Text>
                )}

                {/* Métadonnées */}
                <View style={{ flexDirection: 'row', flexWrap: 'wrap', gap: 8 }}>
                  {/* Priorité */}
                  <View style={{
                    backgroundColor: `${getPriorityColor(task.priority)}20`,
                    paddingHorizontal: 8,
                    paddingVertical: 4,
                    borderRadius: 8,
                  }}>
                    <Text style={{ color: getPriorityColor(task.priority), fontSize: 12 }}>
                      {task.priority === 'high' ? 'Haute' : 
                       task.priority === 'medium' ? 'Moyenne' : 'Basse'} priorité
                    </Text>
                  </View>

                  {/* Date d'échéance */}
                  {task.dueDate && (
                    <View style={{
                      backgroundColor: '#f3f4f6',
                      paddingHorizontal: 8,
                      paddingVertical: 4,
                      borderRadius: 8,
                    }}>
                      <Text style={{ color: '#374151', fontSize: 12 }}>
                        Échéance: {formatDate(task.dueDate)}
                      </Text>
                    </View>
                  )}

                  {/* Temps estimé */}
                  {task.estimatedTime && (
                    <View style={{
                      backgroundColor: '#f3f4f6',
                      paddingHorizontal: 8,
                      paddingVertical: 4,
                      borderRadius: 8,
                    }}>
                      <Text style={{ color: '#374151', fontSize: 12 }}>
                        {task.estimatedTime} min
                      </Text>
                    </View>
                  )}
                </View>

                {/* Tags */}
                {task.tags && task.tags.length > 0 && (
                  <View style={{ flexDirection: 'row', flexWrap: 'wrap', gap: 4, marginTop: 8 }}>
                    {task.tags.map((tag, tagIndex) => (
                      <View
                        key={tagIndex}
                        style={{
                          backgroundColor: '#dbeafe',
                          paddingHorizontal: 6,
                          paddingVertical: 2,
                          borderRadius: 6,
                        }}
                      >
                        <Text style={{ color: '#1e40af', fontSize: 11 }}>
                          #{tag}
                        </Text>
                      </View>
                    ))}
                  </View>
                )}
              </View>
            ))}
          </View>
        )}
      </View>

      {/* Formulaire de tâche */}
      
    </ScrollView>
  );
} 