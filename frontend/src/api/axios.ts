import axios, { AxiosError, type InternalAxiosRequestConfig } from 'axios'
import { useAuthStore } from '@/stores/auth.store'

// Create axios instance
const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL || '/api/v1',
    timeout: 30000,
    headers: {
        'Content-Type': 'application/json'
    }
})

// Enable mock handlers in development if no real API
if (import.meta.env.DEV && import.meta.env.VITE_USE_MOCKS !== 'false') {
    import('@/mocks/handlers').then(({ setupMockHandlers }) => {
        setupMockHandlers(api)
    })
}

// Request interceptor - attach JWT token
api.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
        const authStore = useAuthStore()
        const token = authStore.accessToken

        if (token && config.headers) {
            config.headers.Authorization = `Bearer ${token}`
        }

        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// Response interceptor - handle 401 and refresh token
api.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
        const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean }
        const authStore = useAuthStore()

        // If 401 and not already retrying
        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true

            try {
                // Attempt to refresh token
                const refreshToken = authStore.refreshToken
                if (refreshToken) {
                    const response = await axios.post('/api/v1/auth/refresh', {
                        refresh_token: refreshToken
                    })

                    const { access_token, refresh_token } = response.data.data
                    authStore.setTokens(access_token, refresh_token)

                    // Retry original request with new token
                    if (originalRequest.headers) {
                        originalRequest.headers.Authorization = `Bearer ${access_token}`
                    }
                    return api(originalRequest)
                }
            } catch (refreshError) {
                // Refresh failed, logout user
                authStore.logout()
                window.location.href = '/login'
                return Promise.reject(refreshError)
            }
        }

        return Promise.reject(error)
    }
)

export default api
