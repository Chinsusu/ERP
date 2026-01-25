import api from './axios'
import type { ApiResponse, PaginatedResponse, User, PaginationParams } from '@/types'

export const userApi = {
    /**
     * Get current logged in user profile
     */
    async getCurrentUser(): Promise<ApiResponse<User>> {
        const response = await api.get<ApiResponse<User>>('/v1/users/me')
        return response.data
    },

    /**
     * Get paginated list of users
     */
    async getUsers(params?: PaginationParams): Promise<PaginatedResponse<User>> {
        const response = await api.get<PaginatedResponse<User>>('/v1/users', { params })
        return response.data
    },

    /**
     * Get user by ID
     */
    async getUserById(id: string): Promise<ApiResponse<User>> {
        const response = await api.get<ApiResponse<User>>(`/v1/users/${id}`)
        return response.data
    },

    /**
     * Update user
     */
    async updateUser(id: string, data: Partial<User>): Promise<ApiResponse<User>> {
        const response = await api.put<ApiResponse<User>>(`/v1/users/${id}`, data)
        return response.data
    },

    /**
     * Delete user
     */
    async deleteUser(id: string): Promise<ApiResponse<null>> {
        const response = await api.delete<ApiResponse<null>>(`/v1/users/${id}`)
        return response.data
    },

    /**
     * Update current user profile
     */
    async updateProfile(data: Partial<User>): Promise<ApiResponse<User>> {
        const response = await api.put<ApiResponse<User>>('/v1/users/me', data)
        return response.data
    }
}
