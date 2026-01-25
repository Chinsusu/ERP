import api from './axios'
import type { ApiResponse, PaginatedResponse, PaginationParams } from '@/types/api.types'
import type { BillOfMaterials, BOMItem, WorkOrder, MaterialIssue, QCInspection, TraceabilityRecord } from '@/types/business.types'

export interface BOMFilters extends PaginationParams {
    product_id?: string
    status?: string
    search?: string
}

export interface WorkOrderFilters extends PaginationParams {
    status?: string
    product_id?: string
    priority?: string
    from_date?: string
    to_date?: string
    search?: string
}

export const manufacturingApi = {
    // Bill of Materials
    listBOM: (params: BOMFilters = {}) =>
        api.get<PaginatedResponse<BillOfMaterials>>('/manufacturing/bom', { params }).then(r => r.data),

    getBOM: (id: string) =>
        api.get<ApiResponse<BillOfMaterials>>(`/manufacturing/bom/${id}`).then(r => r.data),

    createBOM: (data: Partial<BillOfMaterials>) =>
        api.post<ApiResponse<BillOfMaterials>>('/manufacturing/bom', data).then(r => r.data),

    updateBOM: (id: string, data: Partial<BillOfMaterials>) =>
        api.put<ApiResponse<BillOfMaterials>>(`/manufacturing/bom/${id}`, data).then(r => r.data),

    deleteBOM: (id: string) =>
        api.delete<ApiResponse<void>>(`/manufacturing/bom/${id}`).then(r => r.data),

    activateBOM: (id: string) =>
        api.patch<ApiResponse<BillOfMaterials>>(`/manufacturing/bom/${id}/activate`).then(r => r.data),

    // BOM Items
    addBOMItem: (bomId: string, item: Partial<BOMItem>) =>
        api.post<ApiResponse<BOMItem>>(`/manufacturing/bom/${bomId}/items`, item).then(r => r.data),

    updateBOMItem: (bomId: string, itemId: string, item: Partial<BOMItem>) =>
        api.put<ApiResponse<BOMItem>>(`/manufacturing/bom/${bomId}/items/${itemId}`, item).then(r => r.data),

    deleteBOMItem: (bomId: string, itemId: string) =>
        api.delete<ApiResponse<void>>(`/manufacturing/bom/${bomId}/items/${itemId}`).then(r => r.data),

    // Work Orders
    listWorkOrders: (params: WorkOrderFilters = {}) =>
        api.get<PaginatedResponse<WorkOrder>>('/manufacturing/work-orders', { params }).then(r => r.data),

    getWorkOrder: (id: string) =>
        api.get<ApiResponse<WorkOrder>>(`/manufacturing/work-orders/${id}`).then(r => r.data),

    createWorkOrder: (data: Partial<WorkOrder>) =>
        api.post<ApiResponse<WorkOrder>>('/manufacturing/work-orders', data).then(r => r.data),

    updateWorkOrder: (id: string, data: Partial<WorkOrder>) =>
        api.put<ApiResponse<WorkOrder>>(`/manufacturing/work-orders/${id}`, data).then(r => r.data),

    deleteWorkOrder: (id: string) =>
        api.delete<ApiResponse<void>>(`/manufacturing/work-orders/${id}`).then(r => r.data),

    releaseWorkOrder: (id: string) =>
        api.patch<ApiResponse<WorkOrder>>(`/manufacturing/work-orders/${id}/release`).then(r => r.data),

    startWorkOrder: (id: string) =>
        api.patch<ApiResponse<WorkOrder>>(`/manufacturing/work-orders/${id}/start`).then(r => r.data),

    completeWorkOrder: (id: string, data: { completed_quantity: number; rejected_quantity: number }) =>
        api.patch<ApiResponse<WorkOrder>>(`/manufacturing/work-orders/${id}/complete`, data).then(r => r.data),

    cancelWorkOrder: (id: string, reason: string) =>
        api.patch<ApiResponse<WorkOrder>>(`/manufacturing/work-orders/${id}/cancel`, { reason }).then(r => r.data),

    // Material Issues
    issueMaterials: (workOrderId: string, items: { material_id: string; lot_id: string; quantity: number }[]) =>
        api.post<ApiResponse<MaterialIssue[]>>(`/manufacturing/work-orders/${workOrderId}/issue-materials`, { items }).then(r => r.data),

    // QC Inspections
    listInspections: (params: PaginationParams & { reference_type?: string; status?: string }) =>
        api.get<PaginatedResponse<QCInspection>>('/manufacturing/qc/inspections', { params }).then(r => r.data),

    getInspection: (id: string) =>
        api.get<ApiResponse<QCInspection>>(`/manufacturing/qc/inspections/${id}`).then(r => r.data),

    createInspection: (data: Partial<QCInspection>) =>
        api.post<ApiResponse<QCInspection>>('/manufacturing/qc/inspections', data).then(r => r.data),

    updateInspection: (id: string, data: Partial<QCInspection>) =>
        api.put<ApiResponse<QCInspection>>(`/manufacturing/qc/inspections/${id}`, data).then(r => r.data),

    completeInspection: (id: string, result: 'PASS' | 'FAIL' | 'CONDITIONAL', comments?: string) =>
        api.patch<ApiResponse<QCInspection>>(`/manufacturing/qc/inspections/${id}/complete`, { result, comments }).then(r => r.data),

    // Traceability
    traceForward: (lotId: string) =>
        api.get<ApiResponse<TraceabilityRecord[]>>(`/manufacturing/traceability/forward/${lotId}`).then(r => r.data),

    traceBackward: (lotId: string) =>
        api.get<ApiResponse<TraceabilityRecord[]>>(`/manufacturing/traceability/backward/${lotId}`).then(r => r.data),

    getTraceabilityTree: (lotId: string) =>
        api.get<ApiResponse<TraceabilityRecord>>(`/manufacturing/traceability/tree/${lotId}`).then(r => r.data),
}
