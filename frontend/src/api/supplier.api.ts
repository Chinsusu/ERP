import api from './axios'
import type { ApiResponse, PaginatedResponse, PaginationParams } from '@/types/api.types'
import type { Supplier, SupplierContact, SupplierCertification, SupplierEvaluation } from '@/types/business.types'

export interface SupplierFilters extends PaginationParams {
    supplier_type?: string
    status?: string
    is_approved?: boolean
    search?: string
}

export const supplierApi = {
    // Suppliers
    list: (params: SupplierFilters = {}) =>
        api.get<PaginatedResponse<Supplier>>('/suppliers', { params }).then(r => r.data),

    getById: (id: string) =>
        api.get<ApiResponse<Supplier>>(`/suppliers/${id}`).then(r => r.data),

    create: (data: Partial<Supplier>) =>
        api.post<ApiResponse<Supplier>>('/suppliers', data).then(r => r.data),

    update: (id: string, data: Partial<Supplier>) =>
        api.put<ApiResponse<Supplier>>(`/suppliers/${id}`, data).then(r => r.data),

    delete: (id: string) =>
        api.delete<ApiResponse<void>>(`/suppliers/${id}`).then(r => r.data),

    approve: (id: string) =>
        api.patch<ApiResponse<Supplier>>(`/suppliers/${id}/approve`).then(r => r.data),

    block: (id: string, reason: string) =>
        api.patch<ApiResponse<Supplier>>(`/suppliers/${id}/block`, { reason }).then(r => r.data),

    // Contacts
    getContacts: (supplierId: string) =>
        api.get<ApiResponse<SupplierContact[]>>(`/suppliers/${supplierId}/contacts`).then(r => r.data),

    addContact: (supplierId: string, data: Partial<SupplierContact>) =>
        api.post<ApiResponse<SupplierContact>>(`/suppliers/${supplierId}/contacts`, data).then(r => r.data),

    updateContact: (supplierId: string, contactId: string, data: Partial<SupplierContact>) =>
        api.put<ApiResponse<SupplierContact>>(`/suppliers/${supplierId}/contacts/${contactId}`, data).then(r => r.data),

    deleteContact: (supplierId: string, contactId: string) =>
        api.delete<ApiResponse<void>>(`/suppliers/${supplierId}/contacts/${contactId}`).then(r => r.data),

    // Certifications
    getCertifications: (supplierId: string) =>
        api.get<ApiResponse<SupplierCertification[]>>(`/suppliers/${supplierId}/certifications`).then(r => r.data),

    addCertification: (supplierId: string, data: Partial<SupplierCertification>) =>
        api.post<ApiResponse<SupplierCertification>>(`/suppliers/${supplierId}/certifications`, data).then(r => r.data),

    updateCertification: (supplierId: string, certId: string, data: Partial<SupplierCertification>) =>
        api.put<ApiResponse<SupplierCertification>>(`/suppliers/${supplierId}/certifications/${certId}`, data).then(r => r.data),

    deleteCertification: (supplierId: string, certId: string) =>
        api.delete<ApiResponse<void>>(`/suppliers/${supplierId}/certifications/${certId}`).then(r => r.data),

    // Evaluations
    getEvaluations: (supplierId: string) =>
        api.get<ApiResponse<SupplierEvaluation[]>>(`/suppliers/${supplierId}/evaluations`).then(r => r.data),

    addEvaluation: (supplierId: string, data: Partial<SupplierEvaluation>) =>
        api.post<ApiResponse<SupplierEvaluation>>(`/suppliers/${supplierId}/evaluations`, data).then(r => r.data),
}
