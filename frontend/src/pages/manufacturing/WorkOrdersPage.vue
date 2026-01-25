<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { useAppStore } from '@/stores/app.store'
import { manufacturingApi, type WorkOrderFilters } from '@/api/manufacturing.api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import Tag from 'primevue/tag'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import ProgressBar from 'primevue/progressbar'

const router = useRouter()
const appStore = useAppStore()

onMounted(() => appStore.setPageTitle('Work Orders'))

const filters = ref<WorkOrderFilters>({ page: 1, page_size: 10, search: '' })
const searchValue = ref('')
const selectedStatus = ref<string | null>(null)
const selectedPriority = ref<string | null>(null)

const statusOptions = [
  { label: 'All Status', value: null },
  { label: 'Draft', value: 'DRAFT' },
  { label: 'Planned', value: 'PLANNED' },
  { label: 'Released', value: 'RELEASED' },
  { label: 'In Progress', value: 'IN_PROGRESS' },
  { label: 'Completed', value: 'COMPLETED' },
  { label: 'Cancelled', value: 'CANCELLED' }
]

const priorityOptions = [
  { label: 'All Priorities', value: null },
  { label: 'Low', value: 'LOW' },
  { label: 'Normal', value: 'NORMAL' },
  { label: 'High', value: 'HIGH' },
  { label: 'Urgent', value: 'URGENT' }
]

const { data, isLoading, refetch } = useQuery({
  queryKey: ['work-orders', filters],
  queryFn: () => manufacturingApi.listWorkOrders(filters.value)
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
watch(selectedPriority, (val) => { filters.value.priority = val || undefined; filters.value.page = 1 })

function viewWO(wo: any) { router.push(`/manufacturing/work-orders/${wo.id}`) }
function createWO() { router.push('/manufacturing/work-orders/new') }

function getStatusSeverity(status: string) {
  switch (status) {
    case 'DRAFT': return 'secondary'
    case 'PLANNED': return 'info'
    case 'RELEASED': return 'info'
    case 'IN_PROGRESS': return 'warning'
    case 'COMPLETED': return 'success'
    case 'CANCELLED': return 'danger'
    default: return 'info'
  }
}

function getPrioritySeverity(priority: string) {
  switch (priority) {
    case 'LOW': return 'secondary'
    case 'NORMAL': return 'info'
    case 'HIGH': return 'warning'
    case 'URGENT': return 'danger'
    default: return 'info'
  }
}

function getProgress(wo: any) {
  if (!wo.planned_quantity) return 0
  return Math.round((wo.completed_quantity / wo.planned_quantity) * 100)
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}
</script>

<template>
  <div class="wo-list-page">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">Work Orders</h1>
        <p class="page-subtitle">Manufacturing production orders</p>
      </div>
      <div class="header-right">
        <Button label="Create Work Order" icon="pi pi-plus" @click="createWO" />
      </div>
    </div>

    <div class="filters-bar">
      <IconField>
        <InputIcon class="pi pi-search" />
        <InputText v-model="searchValue" placeholder="Search..." @keyup.enter="onSearch" />
      </IconField>
      <Dropdown v-model="selectedStatus" :options="statusOptions" optionLabel="label" optionValue="value" placeholder="Status" class="filter-dropdown" />
      <Dropdown v-model="selectedPriority" :options="priorityOptions" optionLabel="label" optionValue="value" placeholder="Priority" class="filter-dropdown" />
      <Button icon="pi pi-refresh" text rounded @click="() => refetch()" />
    </div>

    <div class="card">
      <DataTable :value="items" :loading="isLoading" :paginator="true" :rows="filters.page_size" :totalRecords="totalRecords" :lazy="true" :rowsPerPageOptions="[10, 25, 50]" @page="onPage" dataKey="id" stripedRows>
        <template #empty><div class="empty-message"><i class="pi pi-cog"></i><p>No work orders found</p></div></template>
        
        <Column field="wo_number" header="WO #" sortable style="width: 120px">
          <template #body="{ data }"><span class="font-mono">{{ data.wo_number }}</span></template>
        </Column>
        <Column field="product" header="Product" style="min-width: 180px">
          <template #body="{ data }">
            <div class="font-medium">{{ data.product?.name || '-' }}</div>
          </template>
        </Column>
        <Column field="planned_quantity" header="Quantity" style="width: 120px">
          <template #body="{ data }">
            {{ data.completed_quantity || 0 }} / {{ data.planned_quantity }} {{ data.uom?.code }}
          </template>
        </Column>
        <Column header="Progress" style="width: 120px">
          <template #body="{ data }">
            <ProgressBar :value="getProgress(data)" :showValue="true" style="height: 12px" />
          </template>
        </Column>
        <Column field="priority" header="Priority" style="width: 90px">
          <template #body="{ data }"><Tag :value="data.priority" :severity="getPrioritySeverity(data.priority)" /></template>
        </Column>
        <Column field="status" header="Status" style="width: 110px">
          <template #body="{ data }"><Tag :value="data.status.replace('_', ' ')" :severity="getStatusSeverity(data.status)" /></template>
        </Column>
        <Column field="planned_start_date" header="Schedule" style="width: 140px">
          <template #body="{ data }">
            <div class="text-sm">{{ formatDate(data.planned_start_date) }} - {{ formatDate(data.planned_end_date) }}</div>
          </template>
        </Column>
        <Column header="Actions" style="width: 80px">
          <template #body="{ data }">
            <Button icon="pi pi-eye" text rounded size="small" @click="viewWO(data)" />
          </template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>

<style scoped>
.wo-list-page { padding: 1.5rem; }
.page-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 1.5rem; }
.page-title { font-size: 1.5rem; font-weight: 600; margin-bottom: 0.25rem; }
.page-subtitle { color: var(--text-color-secondary); font-size: 0.875rem; }
.filters-bar { display: flex; gap: 1rem; margin-bottom: 1rem; flex-wrap: wrap; align-items: center; }
.filter-dropdown { min-width: 140px; }
.card { background: var(--surface-card); border-radius: 12px; padding: 1rem; box-shadow: 0 2px 8px rgba(0,0,0,0.08); }
.empty-message { text-align: center; padding: 3rem; color: var(--text-color-secondary); }
.empty-message i { font-size: 3rem; margin-bottom: 1rem; }
.font-mono { font-family: 'JetBrains Mono', monospace; }
.font-medium { font-weight: 500; }
.text-sm { font-size: 0.8125rem; }
</style>
