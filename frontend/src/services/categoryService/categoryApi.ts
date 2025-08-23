import  { apiHelpers } from "../api/client"
import { Category } from "@/src/types"

const categoryApi = {
  getCategories: async (): Promise<Category[]> => {
    try {
      const response = await apiHelpers.get<{data:  Category[]}>('/categories')
      console.log("response getCategories", response)
      
      return response.data
    } catch (error) {
      console.error('Get categories error:', error)
      throw error
    }
  },

  getCategoriesByType: async (type: 'income' | 'expense' | 'task'): Promise<Category[]> => {
    try {
      const response = await apiHelpers.get<{data: {categories: Category[]}}>(`/categories?type=${type}`)
      return response.data.categories
    } catch (error) {
      console.error('Get categories by type error:', error)
      throw error
    }
  },

  createCategory: async (category: { name: string; type: string; parentId?: string }): Promise<Category> => {
    try {
      const response = await apiHelpers.post<{data: {category: Category}}>('/categories', category)
      return response.data.category
    } catch (error) {
      console.error('Create category error:', error)
      throw error
    }
  }
}

export default categoryApi 