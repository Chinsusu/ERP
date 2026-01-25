import { ref, computed } from 'vue'
import type { PaginationParams, Pagination } from '@/types'

export function usePagination(defaultPageSize = 10) {
    const page = ref(1)
    const pageSize = ref(defaultPageSize)
    const totalItems = ref(0)
    const totalPages = ref(0)
    const sortBy = ref<string>('')
    const sortOrder = ref<'asc' | 'desc'>('asc')
    const search = ref('')

    const hasNext = computed(() => page.value < totalPages.value)
    const hasPrevious = computed(() => page.value > 1)

    const paginationParams = computed<PaginationParams>(() => ({
        page: page.value,
        page_size: pageSize.value,
        sort_by: sortBy.value || undefined,
        sort_order: sortBy.value ? sortOrder.value : undefined,
        search: search.value || undefined
    }))

    function setFromResponse(pagination: Pagination) {
        page.value = pagination.page
        pageSize.value = pagination.page_size
        totalItems.value = pagination.total_items
        totalPages.value = pagination.total_pages
    }

    function nextPage() {
        if (hasNext.value) {
            page.value++
        }
    }

    function previousPage() {
        if (hasPrevious.value) {
            page.value--
        }
    }

    function goToPage(pageNum: number) {
        if (pageNum >= 1 && pageNum <= totalPages.value) {
            page.value = pageNum
        }
    }

    function setSort(field: string, order: 'asc' | 'desc' = 'asc') {
        sortBy.value = field
        sortOrder.value = order
        page.value = 1 // Reset to first page on sort change
    }

    function setSearch(value: string) {
        search.value = value
        page.value = 1 // Reset to first page on search
    }

    function reset() {
        page.value = 1
        pageSize.value = defaultPageSize
        sortBy.value = ''
        sortOrder.value = 'asc'
        search.value = ''
        totalItems.value = 0
        totalPages.value = 0
    }

    return {
        // State
        page,
        pageSize,
        totalItems,
        totalPages,
        sortBy,
        sortOrder,
        search,
        // Computed
        hasNext,
        hasPrevious,
        paginationParams,
        // Actions
        setFromResponse,
        nextPage,
        previousPage,
        goToPage,
        setSort,
        setSearch,
        reset
    }
}
