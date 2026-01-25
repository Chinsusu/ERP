import api from './axios'
import type { ApiResponse, PaginatedResponse, PaginationParams } from '@/types/api.types'
import type { Material, Category, UnitOfMeasure, Product } from '@/types/business.types'

export interface MaterialFilters extends PaginationParams {
    material_type?: string
    category_id?: string
    is_active?: boolean
    search?: string
}

export const materialApi = {
    // Materials
    list: (params: MaterialFilters = {}) =>
        api.get<PaginatedResponse<Material>>('/master-data/materials', { params }).then(r => r.data),

    getById: (id: string) =>
        api.get<ApiResponse<Material>>(`/master-data/materials/${id}`).then(r => r.data),

    create: (data: Partial<Material>) =>
        api.post<ApiResponse<Material>>('/master-data/materials', data).then(r => r.data),

    update: (id: string, data: Partial<Material>) =>
        api.put<ApiResponse<Material>>(`/master-data/materials/${id}`, data).then(r => r.data),

    delete: (id: string) =>
        api.delete<ApiResponse<void>>(`/master-data/materials/${id}`).then(r => r.data),

    getSuppliers: (id: string) =>
        api.get<ApiResponse<any[]>>(`/master-data/materials/${id}/suppliers`).then(r => r.data),

    // Categories
    listCategories: (type?: 'MATERIAL' | 'PRODUCT') =>
        api.get<ApiResponse<Category[]>>('/master-data/categories', { params: { category_type: type } }).then(r => r.data),

    getCategory: (id: string) =>
        api.get<ApiResponse<Category>>(`/master-data/categories/${id}`).then(r => r.data),

    createCategory: (data: Partial<Category>) =>
        api.post<ApiResponse<Category>>('/master-data/categories', data).then(r => r.data),

    updateCategory: (id: string, data: Partial<Category>) =>
        api.put<ApiResponse<Category>>(`/master-data/categories/${id}`, data).then(r => r.data),

    // Units of Measure
    listUom: () =>
        api.get<ApiResponse<UnitOfMeasure[]>>('/master-data/uom').then(r => r.data),

    getUom: (id: string) =>
        api.get<ApiResponse<UnitOfMeasure>>(`/master-data/uom/${id}`).then(r => r.data),

    // Products
    listProducts: (params: PaginationParams = {}) =>
        api.get<PaginatedResponse<Product>>('/master-data/products', { params }).then(r => r.data),

    getProduct: (id: string) =>
        api.get<ApiResponse<Product>>(`/master-data/products/${id}`).then(r => r.data),

    createProduct: (data: Partial<Product>) =>
        api.post<ApiResponse<Product>>('/master-data/products', data).then(r => r.data),

    updateProduct: (id: string, data: Partial<Product>) =>
        api.put<ApiResponse<Product>>(`/master-data/products/${id}`, data).then(r => r.data),
}
