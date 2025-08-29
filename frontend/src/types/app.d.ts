
interface Category {
  id: string
  userId: string
  name: string
  type: 'income' | 'expense' | 'task'
  icon: string
  color: string
  parentId?: string
  createdAt: Date
  updatedAt: Date
}
 

interface Account {
  id: string
  userId: string
  name: string
  type: 'checking' | 'savings' | 'mobile_money'
  balance: number
  icon: string
  color: string
  currency: string
  createdAt: Date
  updatedAt: Date
  account_number: string
  created_at: Date
}

// Types pour les transactions
interface Transaction {
  id: string
  userId: string
  accountId: string
  account_id: string
  categoryId: string
  category_id: string
  type: 'income' | 'expense'
  amount: number
  description: string
  date: Date
  recurring: boolean
  createdAt: Date
  updatedAt: Date
  created_at: Date
  updated_at: Date
}

interface AccountDetails {
  account: Account
  transactions: AccountTransaction[]
}

interface AccountTransaction extends Transaction {
  category?: Category
  account?: Account
}

// Types pour les budgets
interface Budget {
  id: string
  user_id: string
  category_id: string
  name: string
  amount_planned: number
  amount_spent: number
  period: string
  created_at: Date
  updated_at: Date
  status: string
  percentage_used: number
  remaining_amount: number
  days_remaining: number
  category?: Category
}