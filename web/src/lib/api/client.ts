import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL
if (!API_BASE_URL) {
  throw new Error('VITE_API_URL environment variable is required')
}

const STORAGE_KEYS = {
  ACCESS_TOKEN: 'memwright_access_token',
  REFRESH_TOKEN: 'memwright_refresh_token',
} as const

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export function getAccessToken(): string | null {
  return localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN)
}

export function getRefreshToken(): string | null {
  return localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN)
}

export function setTokens(accessToken: string, refreshToken: string): void {
  localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, accessToken)
  localStorage.setItem(STORAGE_KEYS.REFRESH_TOKEN, refreshToken)
}

export function clearTokens(): void {
  localStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN)
  localStorage.removeItem(STORAGE_KEYS.REFRESH_TOKEN)
}

let isRefreshing = false
let failedQueue: Array<{
  resolve: (token: string) => void
  reject: (error: unknown) => void
}> = []

function processQueue(error: unknown, token: string | null = null): void {
  failedQueue.forEach((promise) => {
    if (error) {
      promise.reject(error)
    } else if (token) {
      promise.resolve(token)
    }
  })
  failedQueue = []
}

apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getAccessToken()
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean
    }

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        })
          .then((token) => {
            if (originalRequest.headers) {
              originalRequest.headers.Authorization = `Bearer ${token}`
            }
            return apiClient(originalRequest)
          })
          .catch((err) => Promise.reject(err))
      }

      originalRequest._retry = true
      isRefreshing = true

      const refreshToken = getRefreshToken()

      if (!refreshToken) {
        clearTokens()
        window.location.href = '/login'
        return Promise.reject(error)
      }

      try {
        const response = await axios.post<{
          accessToken: string
          refreshToken: string
        }>(`${API_BASE_URL}/auth/refresh`, { refreshToken })

        const { accessToken, refreshToken: newRefreshToken } = response.data
        setTokens(accessToken, newRefreshToken)

        processQueue(null, accessToken)

        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${accessToken}`
        }

        return apiClient(originalRequest)
      } catch (refreshError) {
        processQueue(refreshError, null)
        clearTokens()
        window.location.href = '/login'
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }

    return Promise.reject(error)
  }
)

export interface AuthTokens {
  accessToken: string
  refreshToken: string
}

export interface User {
  id: string
  email: string
  name: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  name: string
}

export interface AuthResponse {
  user: User
  accessToken: string
  refreshToken: string
}

export const authApi = {
  login: (data: LoginRequest) =>
    apiClient.post<AuthResponse>('/auth/login', data),

  register: (data: RegisterRequest) =>
    apiClient.post<AuthResponse>('/auth/register', data),

  logout: () => apiClient.post('/auth/logout'),

  refresh: (refreshToken: string) =>
    apiClient.post<AuthTokens>('/auth/refresh', { refreshToken }),

  me: () => apiClient.get<User>('/auth/me'),
}

export interface Deck {
  id: string
  userId: string
  parentId: string | null
  name: string
  description: string
  cardCount: number
  newCount: number
  dueCount: number
  children?: Deck[]
  createdAt: string
  updatedAt: string
}

export interface CreateDeckRequest {
  name: string
  description?: string
  parentId?: string | null
}

export interface UpdateDeckRequest {
  name?: string
  description?: string
  parentId?: string | null
}

export interface DeckStats {
  totalCards: number
  newCards: number
  dueCards: number
  learnedCards: number
  masteredCards: number
}

export interface ReviewActivity {
  date: string
  count: number
}

export interface GlobalStats {
  totalDecks: number
  totalCards: number
  newCards: number
  cardsLearned: number
  cardsToReview: number
  masteredCards: number
  currentStreak: number
  longestStreak: number
  dailyAverage: number
  daysLearned: number
  totalDays: number
  retentionRate: number
  totalReviews: number
  reviewTime: number
}

export const deckApi = {
  getAll: () => apiClient.get<Deck[]>('/decks'),

  getById: (id: string) => apiClient.get<Deck>(`/decks/${id}`),

  create: (data: CreateDeckRequest) => apiClient.post<Deck>('/decks', data),

  update: (id: string, data: UpdateDeckRequest) =>
    apiClient.patch<Deck>(`/decks/${id}`, data),

  delete: (id: string) => apiClient.delete(`/decks/${id}`),

  getSubdecks: (id: string) => apiClient.get<Deck[]>(`/decks/${id}/subdecks`),

  getStats: (id: string) => apiClient.get<DeckStats>(`/stats/decks/${id}`),
}

export const statsApi = {
  getGlobal: () => apiClient.get<GlobalStats>('/stats/global'),

  getHeatmap: (year?: number) =>
    apiClient.get<ReviewActivity[]>('/stats/heatmap', {
      params: year ? { year } : undefined,
    }),
}
