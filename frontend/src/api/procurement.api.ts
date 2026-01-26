import api from './axios'
import type { ApiResponse, PaginatedResponse, PaginationParams } from '@/types/api.types'
import type { PurchaseRequisition, PurchaseOrder, PRItem, POItem } from '@/types/business.types'

export interface PRFilters extends PaginationParams {
    status?: string
    priority?: string
    requester_id?: string
    from_date?: string
    to_date?: string
    search?: string
}

export interface POFilters extends PaginationParams {
    status?: string
    supplier_id?: string
    from_date?: string
    to_date?: string
    search?: string
}

export const procurementApi = {
    // Purchase Requisitions
    listPR: (params: PRFilters = {}) =>
        api.get<PaginatedResponse<PurchaseRequisition>>('/purchase-requisitions', { params }).then(r => r.data),

    getPR: (id: string) =>
        api.get<ApiResponse<PurchaseRequisition>>(`/purchase-requisitions/${id}`).then(r => r.data),

    createPR: (data: Partial<PurchaseRequisition>) =>
        api.post<ApiResponse<PurchaseRequisition>>('/purchase-requisitions', data).then(r => r.data),

    updatePR: (id: string, data: Partial<PurchaseRequisition>) =>
        api.put<ApiResponse<PurchaseRequisition>>(`/purchase-requisitions/${id}`, data).then(r => r.data),

    deletePR: (id: string) =>
        api.delete<ApiResponse<void>>(`/purchase-requisitions/${id}`).then(r => r.data),

    submitPR: (id: string) =>
        api.post<ApiResponse<PurchaseRequisition>>(`/purchase-requisitions/${id}/submit`).then(r => r.data),

    approvePR: (id: string, comments?: string) =>
        api.post<ApiResponse<PurchaseRequisition>>(`/purchase-requisitions/${id}/approve`, { comments }).then(r => r.data),

    rejectPR: (id: string, reason: string) =>
        api.post<ApiResponse<PurchaseRequisition>>(`/purchase-requisitions/${id}/reject`, { reason }).then(r => r.data),

    // PR Items
    addPRItem: (prId: string, item: Partial<PRItem>) =>
        api.post<ApiResponse<PRItem>>(`/purchase-requisitions/${prId}/items`, item).then(r => r.data),

    updatePRItem: (prId: string, itemId: string, item: Partial<PRItem>) =>
        api.put<ApiResponse<PRItem>>(`/purchase-requisitions/${prId}/items/${itemId}`, item).then(r => r.data),

    deletePRItem: (prId: string, itemId: string) =>
        api.delete<ApiResponse<void>>(`/purchase-requisitions/${prId}/items/${itemId}`).then(r => r.data),

    // Purchase Orders
    listPO: (params: POFilters = {}) =>
        api.get<PaginatedResponse<PurchaseOrder>>('/purchase-orders', { params }).then(r => r.data),

    getPO: (id: string) =>
        api.get<ApiResponse<PurchaseOrder>>(`/purchase-orders/${id}`).then(r => r.data),

    createPO: (data: Partial<PurchaseOrder>) =>
        api.post<ApiResponse<PurchaseOrder>>('/purchase-orders', data).then(r => r.data),

    updatePO: (id: string, data: Partial<PurchaseOrder>) =>
        api.put<ApiResponse<PurchaseOrder>>(`/purchase-orders/${id}`, data).then(r => r.data),

    deletePO: (id: string) =>
        api.delete<ApiResponse<void>>(`/purchase-orders/${id}`).then(r => r.data),

    submitPO: (id: string) =>
        api.post<ApiResponse<PurchaseOrder>>(`/purchase-orders/${id}/submit`).then(r => r.data),

    approvePO: (id: string, comments?: string) =>
        api.post<ApiResponse<PurchaseOrder>>(`/purchase-orders/${id}/approve`, { comments }).then(r => r.data),

    rejectPO: (id: string, reason: string) =>
        api.post<ApiResponse<PurchaseOrder>>(`/purchase-orders/${id}/reject`, { reason }).then(r => r.data),

    confirmPO: (id: string) =>
        api.post<ApiResponse<PurchaseOrder>>(`/purchase-orders/${id}/confirm`).then(r => r.data),

    generatePOFromPR: (prId: string, supplierId: string) =>
        api.post<ApiResponse<PurchaseOrder>>(`/purchase-requisitions/${prId}/convert-to-po`, { supplier_id: supplierId }).then(r => r.data),

    // PO Items
    addPOItem: (poId: string, item: Partial<POItem>) =>
        api.post<ApiResponse<POItem>>(`/purchase-orders/${poId}/items`, item).then(r => r.data),

    updatePOItem: (poId: string, itemId: string, item: Partial<POItem>) =>
        api.put<ApiResponse<POItem>>(`/purchase-orders/${poId}/items/${itemId}`, item).then(r => r.data),

    deletePOItem: (poId: string, itemId: string) =>
        api.delete<ApiResponse<void>>(`/purchase-orders/${poId}/items/${itemId}`).then(r => r.data),
}
