import api from './axios'
import type { ApiResponse, LoginRequest, LoginResponse, RefreshTokenResponse, ForgotPasswordRequest } from '@/types'

export const authApi = {
    /**
     * Login with email and password
     */
    async login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
        const response = await api.post<ApiResponse<LoginResponse>>('/auth/login', data)
        return response.data
    },

    /**
     * Logout current user
     */
    async logout(): Promise<ApiResponse<null>> {
        const response = await api.post<ApiResponse<null>>('/auth/logout')
        return response.data
    },

    /**
     * Refresh access token using refresh token
     */
    async refreshToken(refreshToken: string): Promise<ApiResponse<RefreshTokenResponse>> {
        const response = await api.post<ApiResponse<RefreshTokenResponse>>('/auth/refresh', {
            refresh_token: refreshToken
        })
        return response.data
    },

    /**
     * Request password reset email
     */
    async forgotPassword(data: ForgotPasswordRequest): Promise<ApiResponse<null>> {
        const response = await api.post<ApiResponse<null>>('/auth/forgot-password', data)
        return response.data
    },

    /**
     * Verify current token is valid
     */
    async verifyToken(): Promise<ApiResponse<{ valid: boolean }>> {
        const response = await api.get<ApiResponse<{ valid: boolean }>>('/auth/verify')
        return response.data
    }
}
