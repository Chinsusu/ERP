import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'

export async function authGuard(
    to: RouteLocationNormalized,
    _from: RouteLocationNormalized,
    next: NavigationGuardNext
) {
    const authStore = useAuthStore()
    const requiresAuth = to.meta.requiresAuth !== false
    const requiredPermission = to.meta.permission as string | undefined

    // Check if route requires authentication
    if (requiresAuth) {
        // Not authenticated, redirect to login
        if (!authStore.accessToken) {
            return next({
                path: '/login',
                query: { redirect: to.fullPath }
            })
        }

        // Check specific permission if required
        if (requiredPermission && !authStore.hasPermission(requiredPermission)) {
            // User doesn't have required permission
            return next({ name: 'Dashboard' }) // Redirect to dashboard instead of 403
        }
    }

    // Already authenticated, trying to access auth pages
    if (to.path === '/login' && authStore.accessToken) {
        return next({ path: '/dashboard' })
    }

    next()
}
