// API Response Types

export interface ApiResponse<T> {
    success: boolean
    data: T
    message?: string
    error?: ApiError
}

export interface ApiError {
    code: string
    message: string
    details?: Record<string, string[]>
}

export interface PaginatedResponse<T> {
    success: boolean
    data: T[]
    pagination: Pagination
}

export interface Pagination {
    page: number
    page_size: number
    total_items: number
    total_pages: number
    has_next: boolean
    has_previous: boolean
}

export interface PaginationParams {
    page?: number
    page_size?: number
    sort_by?: string
    sort_order?: 'asc' | 'desc'
    search?: string
}

// Common Entity Types
export interface BaseEntity {
    id: string
    created_at: string
    updated_at: string
}

export interface SoftDeleteEntity extends BaseEntity {
    deleted_at?: string
}

// Select Option for dropdowns
export interface SelectOption<T = string> {
    label: string
    value: T
    disabled?: boolean
}
