// Auth Types
export interface LoginRequest {
    email: string
    password: string
    remember_me?: boolean
}

export interface LoginResponse {
    access_token: string
    refresh_token: string
    expires_in: number
    token_type: string
    user: User
}

export interface RefreshTokenRequest {
    refresh_token: string
}

export interface RefreshTokenResponse {
    access_token: string
    refresh_token: string
    expires_in: number
}

export interface ForgotPasswordRequest {
    email: string
}

// User Types
export interface User {
    id: string
    email: string
    first_name: string
    last_name: string
    full_name: string
    avatar_url?: string
    department_id?: string
    department_name?: string
    employee_code?: string
    is_active: boolean
    roles: Role[]
    permissions: string[]
    created_at: string
    updated_at: string
}

export interface Role {
    id: string
    name: string
    display_name: string
    description?: string
    permissions: Permission[]
}

export interface Permission {
    id: string
    name: string
    display_name: string
    description?: string
    service: string
    resource: string
    action: string
}

// Token Payload (decoded JWT)
export interface TokenPayload {
    sub: string        // user_id
    email: string
    exp: number        // expiration timestamp
    iat: number        // issued at timestamp
    roles: string[]
    permissions: string[]
}
