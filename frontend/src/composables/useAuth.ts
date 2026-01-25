import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'
import { authApi } from '@/api/auth.api'
import { userApi } from '@/api/user.api'
import type { LoginRequest } from '@/types'

export function useAuth() {
    const router = useRouter()
    const authStore = useAuthStore()

    // Computed properties
    const isAuthenticated = computed(() => authStore.isAuthenticated)
    const currentUser = computed(() => authStore.user)
    const loading = computed(() => authStore.loading)

    /**
     * Login with email and password
     */
    async function login(data: LoginRequest) {
        try {
            authStore.setLoading(true)

            const response = await authApi.login(data)
            const { access_token, refresh_token, user } = response.data

            authStore.setTokens(access_token, refresh_token)
            authStore.setUser(user)

            // Redirect to dashboard
            router.push('/dashboard')

            return { success: true }
        } catch (error: any) {
            const message = error.response?.data?.message || 'Login failed'
            return { success: false, error: message }
        } finally {
            authStore.setLoading(false)
        }
    }

    /**
     * Logout current user
     */
    async function logout() {
        try {
            await authApi.logout()
        } catch (error) {
            // Ignore logout API errors
        } finally {
            authStore.logout()
            router.push('/login')
        }
    }

    /**
     * Fetch current user profile
     */
    async function fetchCurrentUser() {
        try {
            const response = await userApi.getCurrentUser()
            authStore.setUser(response.data)
            return response.data
        } catch (error) {
            authStore.logout()
            throw error
        }
    }

    /**
     * Check if user has a specific permission
     */
    function hasPermission(permission: string): boolean {
        return authStore.hasPermission(permission)
    }

    /**
     * Check if user has any of the given permissions
     */
    function hasAnyPermission(permissions: string[]): boolean {
        return authStore.hasAnyPermission(permissions)
    }

    /**
     * Check if user has a specific role
     */
    function hasRole(role: string): boolean {
        return authStore.userRoles.includes(role)
    }

    return {
        // State
        isAuthenticated,
        currentUser,
        loading,
        // Actions
        login,
        logout,
        fetchCurrentUser,
        hasPermission,
        hasAnyPermission,
        hasRole
    }
}
