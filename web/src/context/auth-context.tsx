import {
  createContext,
  useContext,
  useState,
  useEffect,
  type ReactNode,
} from 'react'
import { useNavigate } from 'react-router-dom'
import {
  authApi,
  getAccessToken,
  setTokens,
  clearTokens,
  type User,
  type LoginRequest,
  type RegisterRequest,
} from '@/lib/api'

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
}

interface AuthContextValue extends AuthState {
  login: (data: LoginRequest) => Promise<void>
  register: (data: RegisterRequest) => Promise<void>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)

interface AuthProviderProps {
  children: ReactNode
}

export function AuthProvider({ children }: AuthProviderProps) {
  const navigate = useNavigate()
  const [state, setState] = useState<AuthState>({
    user: null,
    isAuthenticated: false,
    isLoading: true,
  })

  useEffect(() => {
    const initAuth = async () => {
      // Dev mode: bypass authentication with mock user
      if (import.meta.env.VITE_AUTH_BYPASS === 'true') {
        const mockUser: User = {
          id: 'dev-user',
          email: 'dev@memwright.local',
          name: 'Dev User',
        }
        setState({ user: mockUser, isAuthenticated: true, isLoading: false })
        return
      }

      const token = getAccessToken()

      if (!token) {
        setState({ user: null, isAuthenticated: false, isLoading: false })
        return
      }

      try {
        const response = await authApi.me()
        setState({
          user: response.data,
          isAuthenticated: true,
          isLoading: false,
        })
      } catch {
        clearTokens()
        setState({ user: null, isAuthenticated: false, isLoading: false })
      }
    }

    initAuth()
  }, [])

  const login = async (data: LoginRequest): Promise<void> => {
    const response = await authApi.login(data)
    const { user, accessToken, refreshToken } = response.data

    setTokens(accessToken, refreshToken)
    setState({ user, isAuthenticated: true, isLoading: false })
    navigate('/')
  }

  const register = async (data: RegisterRequest): Promise<void> => {
    const response = await authApi.register(data)
    const { user, accessToken, refreshToken } = response.data

    setTokens(accessToken, refreshToken)
    setState({ user, isAuthenticated: true, isLoading: false })
    navigate('/')
  }

  const logout = async (): Promise<void> => {
    try {
      await authApi.logout()
    } catch {
      // Ignore logout errors - we'll clear tokens anyway
    } finally {
      clearTokens()
      setState({ user: null, isAuthenticated: false, isLoading: false })
      navigate('/login')
    }
  }

  return (
    <AuthContext.Provider value={{ ...state, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext)

  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }

  return context
}
