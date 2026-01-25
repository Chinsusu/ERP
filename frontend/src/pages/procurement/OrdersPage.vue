<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { useAppStore } from '@/stores/app.store'
import { procurementApi, type POFilters } from '@/api/procurement.api'
import type { PurchaseOrder } from '@/types/business.types'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import Tag from 'primevue/tag'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'

const router = useRouter()
const appStore = useAppStore()

onMounted(() => appStore.setPageTitle('Purchase Orders'))

const filters = ref<POFilters>({ page: 1, page_size: 10, search: '' })
const searchValue = ref('')
const selectedStatus = ref<string | null>(null)

const statusOptions = [
  { label: 'All Status', value: null },
  { label: 'Draft', value: 'DRAFT' },
  { label: 'Pending Approval', value: 'PENDING_APPROVAL' },
  { label: 'Approved', value: 'APPROVED' },
  { label: 'Confirmed', value: 'CONFIRMED' },
  { label: 'Partially Received', value: 'PARTIALLY_RECEIVED' },
  { label: 'Received', value: 'RECEIVED' },
  { label: 'Cancelled', value: 'CANCELLED' }
]

const { data, isLoading, refetch } = useQuery({
  queryKey: ['purchase-orders', filters],
  queryFn: () => procurementApi.listPO(filters.value)
})

const items = computed(() => data.value?.data || [])
const totalRecords = computed(() => data.value?.pagination?.total_items || 0)

function onPage(event: any) {
  filters.value.page = event.page + 1
  filters.value.page_size = event.rows
}

function onSearch() {
  filters.value.search = searchValue.value
  filters.value.page = 1
}

watch(selectedStatus, (val) => { filters.value.status = val || undefined; filters.value.page = 1 })

function viewPO(po: PurchaseOrder) { router.push(`/procurement/orders/${po.id}`) }
function editPO(po: PurchaseOrder) { router.push(`/procurement/orders/${po.id}/edit`) }
function createPO() { router.push('/procurement/orders/new') }

function getStatusSeverity(status: string) {
  switch (status) {
    case 'DRAFT': return 'secondary'
    case 'PENDING_APPROVAL': return 'warning'
    case 'APPROVED': return 'info'
    case 'CONFIRMED': return 'success'
    case 'PARTIALLY_RECEIVED': return 'warning'
    case 'RECEIVED': return 'success'
    case 'CANCELLED': return 'danger'
    default: return 'info'
  }
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function formatCurrency(amount: number, currency: string) {
  return amount?.toLocaleString() + ' ' + currency
}
</script>

<template>
  <div class="po-list-page">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">Purchase Orders</h1>
        <p class="page-subtitle">Manage orders to suppliers</p>
      </div>
      <div class="header-right">
        <Button label="Create PO" icon="pi pi-plus" @click="createPO" />
      </div>
    </div>

    <div class="filters-bar">
      <IconField>
        <InputIcon class="pi pi-search" />
        <InputText v-model="searchValue" placeholder="Search PO..." @keyup.enter="onSearch" />
      </IconField>
      <Dropdown v-model="selectedStatus" :options="statusOptions" optionLabel="label" optionValue="value" placeholder="Status" class="filter-dropdown" />
      <Button icon="pi pi-refresh" text rounded @click="() => refetch()" />
    </div>

    <div class="card">
      <DataTable :value="items" :loading="isLoading" :paginator="true" :rows="filters.page_size" :totalRecords="totalRecords" :lazy="true" :rowsPerPageOptions="[10, 25, 50]" @page="onPage" dataKey="id" stripedRows>
        <template #empty><div class="empty-message"><i class="pi pi-shopping-cart"></i><p>No orders found</p></div></template>
        
        <Column field="po_number" header="PO #" sortable style="width: 120px">
          <template #body="{ data }"><span class="font-mono">{{ data.po_number }}</span></template>
        </Column>
        <Column field="supplier" header="Supplier" style="min-width: 180px">
          <template #body="{ data }">
            <div class="font-medium">{{ data.supplier?.name || '-' }}</div>
          </template>
        </Column>
        <Column field="order_date" header="Order Date" style="width: 120px">
          <template #body="{ data }">{{ formatDate(data.order_date) }}</template>
        </Column>
        <Column field="expected_delivery_date" header="Expected" style="width: 120px">
          <template #body="{ data }">{{ data.expected_delivery_date ? formatDate(data.expected_delivery_date) : '-' }}</template>
        </Column>
        <Column field="status" header="Status" style="width: 140px">
          <template #body="{ data }"><Tag :value="data.status.replace('_', ' ')" :severity="getStatusSeverity(data.status)" /></template>
        </Column>
        <Column field="total_amount" header="Total" style="width: 140px">
          <template #body="{ data }">{{ formatCurrency(data.total_amount, data.currency) }}</template>
        </Column>
        <Column header="Actions" style="width: 100px">
          <template #body="{ data }">
            <Button icon="pi pi-eye" text rounded size="small" @click="viewPO(data)" />
            <Button icon="pi pi-pencil" text rounded size="small" @click="editPO(data)" :disabled="!['DRAFT', 'PENDING_APPROVAL'].includes(data.status)" />
          </template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>

<style scoped>
.po-list-page { padding: 1.5rem; }
.page-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 1.5rem; }
.page-title { font-size: 1.5rem; font-weight: 600; margin-bottom: 0.25rem; }
.page-subtitle { color: var(--text-color-secondary); font-size: 0.875rem; }
.filters-bar { display: flex; gap: 1rem; margin-bottom: 1rem; flex-wrap: wrap; align-items: center; }
.filter-dropdown { min-width: 160px; }
.card { background: var(--surface-card); border-radius: 12px; padding: 1rem; box-shadow: 0 2px 8px rgba(0,0,0,0.08); }
.empty-message { text-align: center; padding: 3rem; color: var(--text-color-secondary); }
.empty-message i { font-size: 3rem; margin-bottom: 1rem; }
.font-mono { font-family: 'JetBrains Mono', monospace; }
.font-medium { font-weight: 500; }
</style>
