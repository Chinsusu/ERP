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
        api.get<PaginatedResponse<PurchaseRequisition>>('/procurement/requisitions', { params }).then(r => r.data),

    getPR: (id: string) =>
        api.get<ApiResponse<PurchaseRequisition>>(`/procurement/requisitions/${id}`).then(r => r.data),

    createPR: (data: Partial<PurchaseRequisition>) =>
        api.post<ApiResponse<PurchaseRequisition>>('/procurement/requisitions', data).then(r => r.data),

    updatePR: (id: string, data: Partial<PurchaseRequisition>) =>
        api.put<ApiResponse<PurchaseRequisition>>(`/procurement/requisitions/${id}`, data).then(r => r.data),

    deletePR: (id: string) =>
        api.delete<ApiResponse<void>>(`/procurement/requisitions/${id}`).then(r => r.data),

    submitPR: (id: string) =>
        api.patch<ApiResponse<PurchaseRequisition>>(`/procurement/requisitions/${id}/submit`).then(r => r.data),

    approvePR: (id: string, comments?: string) =>
        api.patch<ApiResponse<PurchaseRequisition>>(`/procurement/requisitions/${id}/approve`, { comments }).then(r => r.data),

    rejectPR: (id: string, reason: string) =>
        api.patch<ApiResponse<PurchaseRequisition>>(`/procurement/requisitions/${id}/reject`, { reason }).then(r => r.data),

    // PR Items
    addPRItem: (prId: string, item: Partial<PRItem>) =>
        api.post<ApiResponse<PRItem>>(`/procurement/requisitions/${prId}/items`, item).then(r => r.data),

    updatePRItem: (prId: string, itemId: string, item: Partial<PRItem>) =>
        api.put<ApiResponse<PRItem>>(`/procurement/requisitions/${prId}/items/${itemId}`, item).then(r => r.data),

    deletePRItem: (prId: string, itemId: string) =>
        api.delete<ApiResponse<void>>(`/procurement/requisitions/${prId}/items/${itemId}`).then(r => r.data),

    // Purchase Orders
    listPO: (params: POFilters = {}) =>
        api.get<PaginatedResponse<PurchaseOrder>>('/procurement/orders', { params }).then(r => r.data),

    getPO: (id: string) =>
        api.get<ApiResponse<PurchaseOrder>>(`/procurement/orders/${id}`).then(r => r.data),

    createPO: (data: Partial<PurchaseOrder>) =>
        api.post<ApiResponse<PurchaseOrder>>('/procurement/orders', data).then(r => r.data),

    updatePO: (id: string, data: Partial<PurchaseOrder>) =>
        api.put<ApiResponse<PurchaseOrder>>(`/procurement/orders/${id}`, data).then(r => r.data),

    deletePO: (id: string) =>
        api.delete<ApiResponse<void>>(`/procurement/orders/${id}`).then(r => r.data),

    submitPO: (id: string) =>
        api.patch<ApiResponse<PurchaseOrder>>(`/procurement/orders/${id}/submit`).then(r => r.data),

    approvePO: (id: string, comments?: string) =>
        api.patch<ApiResponse<PurchaseOrder>>(`/procurement/orders/${id}/approve`, { comments }).then(r => r.data),

    rejectPO: (id: string, reason: string) =>
        api.patch<ApiResponse<PurchaseOrder>>(`/procurement/orders/${id}/reject`, { reason }).then(r => r.data),

    confirmPO: (id: string) =>
        api.patch<ApiResponse<PurchaseOrder>>(`/procurement/orders/${id}/confirm`).then(r => r.data),

    generatePOFromPR: (prId: string, supplierId: string) =>
        api.post<ApiResponse<PurchaseOrder>>('/procurement/orders/from-pr', { pr_id: prId, supplier_id: supplierId }).then(r => r.data),

    // PO Items
    addPOItem: (poId: string, item: Partial<POItem>) =>
        api.post<ApiResponse<POItem>>(`/procurement/orders/${poId}/items`, item).then(r => r.data),

    updatePOItem: (poId: string, itemId: string, item: Partial<POItem>) =>
        api.put<ApiResponse<POItem>>(`/procurement/orders/${poId}/items/${itemId}`, item).then(r => r.data),

    deletePOItem: (poId: string, itemId: string) =>
        api.delete<ApiResponse<void>>(`/procurement/orders/${poId}/items/${itemId}`).then(r => r.data),
}
