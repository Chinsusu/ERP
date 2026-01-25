import { useQuery, useMutation, useQueryClient, type UseQueryOptions, type UseMutationOptions } from '@tanstack/vue-query'
import { useToast } from 'primevue/usetoast'
import type { ApiResponse, ApiError, PaginatedResponse, PaginationParams } from '@/types'
import { computed, type Ref, type ComputedRef } from 'vue'

// Generic API query hook
export function useApiQuery<T>(
    queryKey: string[],
    queryFn: () => Promise<ApiResponse<T>>,
    options?: Omit<UseQueryOptions<ApiResponse<T>>, 'queryKey' | 'queryFn'>
) {
    const query = useQuery({
        queryKey,
        queryFn,
        staleTime: 5 * 60 * 1000, // 5 minutes
        ...options
    })

    const data = computed(() => query.data.value?.data)
    const error = computed(() => query.error.value as ApiError | null)

    return {
        ...query,
        data,
        error
    }
}

// Paginated query hook
export function usePaginatedQuery<T>(
    queryKey: (string | Ref<any> | ComputedRef<any>)[],
    queryFn: (params: PaginationParams) => Promise<PaginatedResponse<T>>,
    params: Ref<PaginationParams>,
    options?: Omit<UseQueryOptions<PaginatedResponse<T>>, 'queryKey' | 'queryFn'>
) {
    const query = useQuery({
        queryKey: [...queryKey, params],
        queryFn: () => queryFn(params.value),
        staleTime: 2 * 60 * 1000, // 2 minutes for lists
        ...options
    })

    const items = computed(() => query.data.value?.data || [])
    const pagination = computed(() => query.data.value?.pagination)
    const total = computed(() => pagination.value?.total_items || 0)

    return {
        ...query,
        items,
        pagination,
        total
    }
}

// Mutation hook with toast notifications
export function useApiMutation<TData, TVariables>(
    mutationFn: (variables: TVariables) => Promise<ApiResponse<TData>>,
    options?: {
        successMessage?: string
        errorMessage?: string
        invalidateKeys?: string[][]
        onSuccess?: (data: TData) => void
        onError?: (error: ApiError) => void
    }
) {
    const toast = useToast()
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn,
        onSuccess: (response) => {
            if (options?.successMessage) {
                toast.add({
                    severity: 'success',
                    summary: 'Success',
                    detail: options.successMessage,
                    life: 3000
                })
            }

            // Invalidate related queries
            if (options?.invalidateKeys) {
                options.invalidateKeys.forEach(key => {
                    queryClient.invalidateQueries({ queryKey: key })
                })
            }

            if (options?.onSuccess && response.data) {
                options.onSuccess(response.data)
            }
        },
        onError: (error: ApiError) => {
            toast.add({
                severity: 'error',
                summary: 'Error',
                detail: options?.errorMessage || error.message || 'An error occurred',
                life: 5000
            })

            if (options?.onError) {
                options.onError(error)
            }
        }
    })
}

// Delete mutation with confirmation
export function useDeleteMutation<T>(
    deleteFn: (id: string) => Promise<ApiResponse<T>>,
    options?: {
        resourceName?: string
        invalidateKeys?: string[][]
        onSuccess?: () => void
    }
) {
    return useApiMutation(deleteFn, {
        successMessage: `${options?.resourceName || 'Item'} deleted successfully`,
        errorMessage: `Failed to delete ${options?.resourceName?.toLowerCase() || 'item'}`,
        invalidateKeys: options?.invalidateKeys,
        onSuccess: options?.onSuccess
    })
}
