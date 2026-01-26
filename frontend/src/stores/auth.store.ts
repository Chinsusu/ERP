import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'

export const useAuthStore = defineStore('auth', () => {
    // State
    const user = ref<User | null>(null)
    const accessToken = ref<string | null>(localStorage.getItem('access_token'))
    const refreshToken = ref<string | null>(localStorage.getItem('refresh_token'))
    const loading = ref(false)
    const initialized = ref(false)

    // Getters
    const isAuthenticated = computed(() => !!accessToken.value && !!user.value)

    const userPermissions = computed(() => user.value?.permissions || [])

    const userRoles = computed(() => user.value?.roles?.map(r => r.name) || [])

    // Check if user has a specific permission
    // Permission format: service:resource:action (e.g., procurement:po:create)
    function hasPermission(permission: string): boolean {
        if (!user.value) return false

        // Super admin has all permissions
        if (userRoles.value.includes('super_admin')) return true

        // Check exact match
        if (userPermissions.value.includes(permission)) return true

        // Check wildcard permissions
        const [service, resource, action] = permission.split(':')
        const wildcardPatterns = [
            `${service}:*:*`,
            `${service}:${resource}:*`,
            '*:*:*'
        ]

        return wildcardPatterns.some(pattern => userPermissions.value.includes(pattern))
    }

    // Check if user has any of the given permissions
    function hasAnyPermission(permissions: string[]): boolean {
        return permissions.some(p => hasPermission(p))
    }

    // Check if user has all of the given permissions
    function hasAllPermissions(permissions: string[]): boolean {
        return permissions.every(p => hasPermission(p))
    }

    // Actions
    function setUser(userData: User) {
        user.value = userData
    }

    function setTokens(access: string, refresh: string) {
        accessToken.value = access
        refreshToken.value = refresh
        localStorage.setItem('access_token', access)
        localStorage.setItem('refresh_token', refresh)
    }

    function logout() {
        user.value = null
        accessToken.value = null
        refreshToken.value = null
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
    }

    function setLoading(value: boolean) {
        loading.value = value
    }

    async function initialize() {
        if (initialized.value || !accessToken.value) {
            initialized.value = true
            return
        }

        try {
            const { userApi } = await import('@/api/user.api')
            const response = await userApi.getCurrentUser()
            if (response.success) {
                user.value = response.data
            }
        } catch (error) {
            console.error('Failed to initialize auth store:', error)
            // If token is invalid, logout
            logout()
        } finally {
            initialized.value = true
        }
    }

    return {
        // State
        user,
        accessToken,
        refreshToken,
        loading,
        // Getters
        isAuthenticated,
        userPermissions,
        userRoles,
        // Actions
        hasPermission,
        hasAnyPermission,
        hasAllPermissions,
        setUser,
        setTokens,
        logout,
        setLoading,
        initialize,
        initialized
    }
})
