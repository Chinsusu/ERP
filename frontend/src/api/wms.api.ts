import api from './axios'
import type { ApiResponse, PaginatedResponse, PaginationParams } from '@/types/api.types'
import type { Stock, Lot, GoodsReceiptNote, GRNItem, StockMovement } from '@/types/business.types'

export interface StockFilters extends PaginationParams {
    material_id?: string
    warehouse_id?: string
    below_minimum?: boolean
    search?: string
}

export interface LotFilters extends PaginationParams {
    material_id?: string
    status?: string
    expiring_within_days?: number
    search?: string
}

export interface GRNFilters extends PaginationParams {
    status?: string
    supplier_id?: string
    po_id?: string
    from_date?: string
    to_date?: string
    search?: string
}

export const wmsApi = {
    // Stock Overview
    getStockOverview: () =>
        api.get<ApiResponse<{
            total_materials: number
            total_value: number
            low_stock_count: number
            expiring_soon_count: number
        }>>('/warehouse/stock/overview').then(r => r.data),

    listStock: (params: StockFilters = {}) =>
        api.get<PaginatedResponse<Stock>>('/warehouse/stock', { params }).then(r => r.data),

    getStockByMaterial: (materialId: string) =>
        api.get<ApiResponse<Stock[]>>(`/warehouse/stock/by-material/${materialId}`).then(r => r.data),

    getExpiringStock: (days: number = 30) =>
        api.get<PaginatedResponse<Lot>>('/warehouse/stock/expiring', { params: { days } }).then(r => r.data),

    // Lots
    listLots: (params: LotFilters = {}) =>
        api.get<PaginatedResponse<Lot>>('/warehouse/lots', { params }).then(r => r.data),

    getLot: (id: string) =>
        api.get<ApiResponse<Lot>>(`/warehouse/lots/${id}`).then(r => r.data),

    getLotMovements: (lotId: string) =>
        api.get<ApiResponse<StockMovement[]>>(`/warehouse/lots/${lotId}/movements`).then(r => r.data),

    updateLotStatus: (id: string, status: string, reason?: string) =>
        api.patch<ApiResponse<Lot>>(`/warehouse/lots/${id}/status`, { status, reason }).then(r => r.data),

    // Goods Receipt Notes
    listGRN: (params: GRNFilters = {}) =>
        api.get<PaginatedResponse<GoodsReceiptNote>>('/warehouse/grn', { params }).then(r => r.data),

    getGRN: (id: string) =>
        api.get<ApiResponse<GoodsReceiptNote>>(`/warehouse/grn/${id}`).then(r => r.data),

    createGRN: (data: Partial<GoodsReceiptNote>) =>
        api.post<ApiResponse<GoodsReceiptNote>>('/warehouse/grn', data).then(r => r.data),

    updateGRN: (id: string, data: Partial<GoodsReceiptNote>) =>
        api.put<ApiResponse<GoodsReceiptNote>>(`/warehouse/grn/${id}`, data).then(r => r.data),

    completeGRN: (id: string) =>
        api.patch<ApiResponse<GoodsReceiptNote>>(`/warehouse/grn/${id}/complete`).then(r => r.data),

    cancelGRN: (id: string, reason: string) =>
        api.patch<ApiResponse<GoodsReceiptNote>>(`/warehouse/grn/${id}/cancel`, { reason }).then(r => r.data),

    // GRN Items
    addGRNItem: (grnId: string, item: Partial<GRNItem>) =>
        api.post<ApiResponse<GRNItem>>(`/warehouse/grn/${grnId}/items`, item).then(r => r.data),

    updateGRNItem: (grnId: string, itemId: string, item: Partial<GRNItem>) =>
        api.put<ApiResponse<GRNItem>>(`/warehouse/grn/${grnId}/items/${itemId}`, item).then(r => r.data),

    // Stock Movements
    listMovements: (params: PaginationParams & { material_id?: string; lot_id?: string; movement_type?: string }) =>
        api.get<PaginatedResponse<StockMovement>>('/warehouse/movements', { params }).then(r => r.data),

    // Goods Issue
    createGoodsIssue: (data: {
        reference_type: string
        reference_id: string
        items: { material_id: string; lot_id: string; quantity: number }[]
    }) =>
        api.post<ApiResponse<StockMovement[]>>('/warehouse/goods-issue', data).then(r => r.data),
}
