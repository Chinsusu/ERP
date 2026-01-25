// Mock API handlers - intercepts axios requests and returns mock data
import type { AxiosInstance } from 'axios'
import type { PaginatedResponse, ApiResponse } from '@/types/api.types'
import {
    mockMaterials, mockCategories, mockUoMs,
    mockSuppliers,
    mockPRs, mockPOs,
    mockStock, mockGRNs,
    mockBOMs, mockWorkOrders
} from './data'

// Helper to paginate array
function paginate<T>(items: T[], page: number, pageSize: number): PaginatedResponse<T> {
    const start = (page - 1) * pageSize
    const end = start + pageSize
    const data = items.slice(start, end)
    const totalPages = Math.ceil(items.length / pageSize)

    return {
        success: true,
        data,
        pagination: {
            page,
            page_size: pageSize,
            total_items: items.length,
            total_pages: totalPages,
            has_next: page < totalPages,
            has_previous: page > 1
        }
    }
}

// Helper to wrap in ApiResponse
function wrapResponse<T>(data: T): ApiResponse<T> {
    return {
        success: true,
        data,
        message: 'OK'
    }
}

// Search/filter helper
function filterBySearch<T extends Record<string, any>>(items: T[], search: string, fields: string[]): T[] {
    if (!search) return items
    const lowerSearch = search.toLowerCase()
    return items.filter(item =>
        fields.some(field => {
            const value = item[field]
            return value && String(value).toLowerCase().includes(lowerSearch)
        })
    )
}

// Setup mock interceptors
export function setupMockHandlers(axios: AxiosInstance): void {
    console.log('🎭 Mock API handlers enabled')

    axios.interceptors.request.use(async (config) => {
        const url = config.url || ''
        const method = config.method?.toUpperCase()
        const params = config.params || {}

        // Extract query params
        const page = parseInt(params.page) || 1
        const pageSize = parseInt(params.page_size) || 10
        const search = params.search || ''

        let mockResponse: any = null

        // ============================================
        // MATERIALS
        // ============================================
        if (url.includes('/master-data/materials') && method === 'GET') {
            if (url.match(/\/materials\/[^/]+$/)) {
                // Get single material
                const id = url.split('/').pop()
                const material = mockMaterials.find(m => m.id === id)
                mockResponse = wrapResponse(material)
            } else {
                // List materials
                let filtered = filterBySearch(mockMaterials, search, ['material_code', 'name', 'inci_name'])
                if (params.material_type) {
                    filtered = filtered.filter(m => m.material_type === params.material_type)
                }
                if (params.is_active !== undefined) {
                    filtered = filtered.filter(m => m.is_active === (params.is_active === 'true'))
                }
                mockResponse = paginate(filtered, page, pageSize)
            }
        }

        // Categories
        if (url.includes('/master-data/categories') && method === 'GET') {
            mockResponse = wrapResponse(mockCategories)
        }

        // UoM
        if (url.includes('/master-data/uom') && method === 'GET') {
            mockResponse = wrapResponse(mockUoMs)
        }

        // ============================================
        // SUPPLIERS
        // ============================================
        if (url.includes('/suppliers') && !url.includes('/contacts') && !url.includes('/certifications') && method === 'GET') {
            if (url.match(/\/suppliers\/[^/]+$/)) {
                const id = url.split('/').pop()
                const supplier = mockSuppliers.find(s => s.id === id)
                mockResponse = wrapResponse(supplier)
            } else {
                let filtered = filterBySearch(mockSuppliers, search, ['supplier_code', 'name', 'legal_name'])
                if (params.supplier_type) {
                    filtered = filtered.filter(s => s.supplier_type === params.supplier_type)
                }
                if (params.status) {
                    filtered = filtered.filter(s => s.status === params.status)
                }
                mockResponse = paginate(filtered, page, pageSize)
            }
        }

        // Supplier contacts/certifications/evaluations
        if (url.includes('/suppliers/') && url.includes('/contacts') && method === 'GET') {
            mockResponse = wrapResponse([])
        }
        if (url.includes('/suppliers/') && url.includes('/certifications') && method === 'GET') {
            mockResponse = wrapResponse([])
        }
        if (url.includes('/suppliers/') && url.includes('/evaluations') && method === 'GET') {
            mockResponse = wrapResponse([])
        }

        // ============================================
        // PROCUREMENT - PR
        // ============================================
        if (url.includes('/procurement/requisitions') && method === 'GET') {
            if (url.match(/\/requisitions\/[^/]+$/)) {
                const id = url.split('/').pop()
                const pr = mockPRs.find(p => p.id === id)
                mockResponse = wrapResponse(pr)
            } else {
                let filtered = filterBySearch(mockPRs, search, ['pr_number', 'title'])
                if (params.status) {
                    filtered = filtered.filter(p => p.status === params.status)
                }
                if (params.priority) {
                    filtered = filtered.filter(p => p.priority === params.priority)
                }
                mockResponse = paginate(filtered, page, pageSize)
            }
        }

        // PROCUREMENT - PO
        if (url.includes('/procurement/orders') && method === 'GET') {
            if (url.match(/\/orders\/[^/]+$/)) {
                const id = url.split('/').pop()
                const po = mockPOs.find(p => p.id === id)
                mockResponse = wrapResponse(po)
            } else {
                let filtered = filterBySearch(mockPOs, search, ['po_number'])
                if (params.status) {
                    filtered = filtered.filter(p => p.status === params.status)
                }
                mockResponse = paginate(filtered, page, pageSize)
            }
        }

        // ============================================
        // WMS - STOCK
        // ============================================
        if (url.includes('/warehouse/stock/overview') && method === 'GET') {
            mockResponse = wrapResponse({
                total_materials: mockStock.length,
                total_value: 125000000,
                low_stock_count: 1,
                expiring_soon_count: 2
            })
        }

        if (url.includes('/warehouse/stock/expiring') && method === 'GET') {
            mockResponse = paginate([], page, pageSize)
        }

        if (url.includes('/warehouse/stock') && !url.includes('overview') && !url.includes('expiring') && method === 'GET') {
            let filtered = mockStock
            if (params.below_minimum) {
                filtered = filtered.filter(s => s.available_quantity < (s.material?.min_stock_quantity || 0))
            }
            mockResponse = paginate(filtered, page, pageSize)
        }

        // WMS - LOTS
        if (url.includes('/warehouse/lots') && method === 'GET') {
            mockResponse = paginate([], page, pageSize)
        }

        // WMS - GRN
        if (url.includes('/warehouse/grn') && method === 'GET') {
            if (url.match(/\/grn\/[^/]+$/)) {
                const id = url.split('/').pop()
                const grn = mockGRNs.find(g => g.id === id)
                mockResponse = wrapResponse(grn)
            } else {
                let filtered = filterBySearch(mockGRNs, search, ['grn_number', 'po_number'])
                if (params.status) {
                    filtered = filtered.filter(g => g.status === params.status)
                }
                mockResponse = paginate(filtered, page, pageSize)
            }
        }

        // ============================================
        // MANUFACTURING - BOM
        // ============================================
        if (url.includes('/manufacturing/bom') && method === 'GET') {
            if (url.match(/\/bom\/[^/]+$/)) {
                const id = url.split('/').pop()
                const bom = mockBOMs.find(b => b.id === id)
                mockResponse = wrapResponse(bom)
            } else {
                let filtered = filterBySearch(mockBOMs, search, ['bom_code'])
                if (params.status) {
                    filtered = filtered.filter(b => b.status === params.status)
                }
                mockResponse = paginate(filtered, page, pageSize)
            }
        }

        // MANUFACTURING - WORK ORDERS
        if (url.includes('/manufacturing/work-orders') && method === 'GET') {
            if (url.match(/\/work-orders\/[^/]+$/)) {
                const id = url.split('/').pop()
                const wo = mockWorkOrders.find(w => w.id === id)
                mockResponse = wrapResponse(wo)
            } else {
                let filtered = filterBySearch(mockWorkOrders, search, ['wo_number'])
                if (params.status) {
                    filtered = filtered.filter(w => w.status === params.status)
                }
                if (params.priority) {
                    filtered = filtered.filter(w => w.priority === params.priority)
                }
                mockResponse = paginate(filtered, page, pageSize)
            }
        }

        // If we have a mock response, return it as a resolved adapter
        if (mockResponse !== null) {
            console.log(`🎭 Mock: ${method} ${url}`, mockResponse)

            // Create a mock adapter response
            config.adapter = () => {
                return Promise.resolve({
                    data: mockResponse,
                    status: 200,
                    statusText: 'OK',
                    headers: {},
                    config
                })
            }
        }

        return config
    })
}

export default setupMockHandlers
