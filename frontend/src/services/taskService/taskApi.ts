import { apiHelpers } from "../api/client"
import { Task, TaskStats } from "@/src/types"

const taskApi = {
  getTasks: async (): Promise<Task[]> => {
    try {
      const response = await apiHelpers.get<{data: Task[]}>('/tasks')
      return response.data
    } catch (error) {
      console.error('Get tasks error:', error)
      throw error
    }
  },

  getTask: async (id: string): Promise<Task> => {
    try {
      const response = await apiHelpers.get<{data: Task}>(`/tasks/${id}`)
      return response.data
    } catch (error) {
      console.error('Get task error:', error)
      throw error
    }
  },

  createTask: async (task: {
    title: string
    description?: string
    priority: 'low' | 'medium' | 'high'
    dueDate?: Date
    categoryId?: string
    durationPlanned?: number
  }): Promise<Task> => {
    try {
      const response = await apiHelpers.post<{data: Task}>('/tasks', task)
      return response.data
    } catch (error) {
      console.error('Create task error:', error)
      throw error
    }
  },

  updateTask: async (id: string, task: Partial<Task>): Promise<Task> => {
    try {
      const response = await apiHelpers.put<{data: Task}>(`/tasks/${id}`, task)
      return response.data
    } catch (error) {
      console.error('Update task error:', error)
      throw error
    }
  },

  deleteTask: async (id: string): Promise<void> => {
    try {
      await apiHelpers.delete(`/tasks/${id}`)
    } catch (error) {
      console.error('Delete task error:', error)
      throw error
    }
  },

  completeTask: async (id: string): Promise<Task> => {
    try {
      const response = await apiHelpers.post<{data: Task}>(`/tasks/${id}/complete`, {})
      return response.data
    } catch (error) {
      console.error('Complete task error:', error)
      throw error
    }
  },

  getTaskStats: async (): Promise<TaskStats> => {
    try {
      const response = await apiHelpers.get<{data: TaskStats}>('/tasks/stats')
      return response.data
    } catch (error) {
      console.error('Get task stats error:', error)
      throw error
    }
  }
}

export default taskApi 