<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { useAppStore } from '@/stores/app.store'
import { procurementApi, type PRFilters } from '@/api/procurement.api'
import type { PurchaseRequisition } from '@/types/business.types'

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

onMounted(() => appStore.setPageTitle('Purchase Requisitions'))

const filters = ref<PRFilters>({ page: 1, page_size: 10, search: '' })
const searchValue = ref('')
const selectedStatus = ref<string | null>(null)
const selectedPriority = ref<string | null>(null)

const statusOptions = [
  { label: 'All Status', value: null },
  { label: 'Draft', value: 'DRAFT' },
  { label: 'Submitted', value: 'SUBMITTED' },
  { label: 'Approved', value: 'APPROVED' },
  { label: 'Rejected', value: 'REJECTED' },
  { label: 'Converted', value: 'CONVERTED' }
]

const priorityOptions = [
  { label: 'All Priorities', value: null },
  { label: 'Low', value: 'LOW' },
  { label: 'Normal', value: 'NORMAL' },
  { label: 'High', value: 'HIGH' },
  { label: 'Urgent', value: 'URGENT' }
]

const { data, isLoading, refetch } = useQuery({
  queryKey: ['purchase-requisitions', filters],
  queryFn: () => procurementApi.listPR(filters.value)
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

function viewPR(pr: PurchaseRequisition) { router.push(`/procurement/requisitions/${pr.id}`) }
function editPR(pr: PurchaseRequisition) { router.push(`/procurement/requisitions/${pr.id}/edit`) }
function createPR() { router.push('/procurement/requisitions/new') }

function getStatusSeverity(status: string) {
  switch (status) {
    case 'DRAFT': return 'secondary'
    case 'SUBMITTED': return 'info'
    case 'APPROVED': return 'success'
    case 'REJECTED': return 'danger'
    case 'CONVERTED': return 'warning'
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

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function formatCurrency(amount: number, currency: string) {
  return amount?.toLocaleString() + ' ' + currency
}
</script>

<template>
  <div class="pr-list-page">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">Purchase Requisitions</h1>
        <p class="page-subtitle">Manage purchase requests from departments</p>
      </div>
      <div class="header-right">
        <Button label="Create PR" icon="pi pi-plus" @click="createPR" />
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
        <template #empty><div class="empty-message"><i class="pi pi-file-edit"></i><p>No requisitions found</p></div></template>
        
        <Column field="pr_number" header="PR #" sortable style="width: 120px">
          <template #body="{ data }"><span class="font-mono">{{ data.pr_number }}</span></template>
        </Column>
        <Column field="title" header="Title" style="min-width: 200px">
          <template #body="{ data }">
            <div class="font-medium">{{ data.title }}</div>
            <div class="text-sm text-muted">{{ data.requester_name }}</div>
          </template>
        </Column>
        <Column field="priority" header="Priority" style="width: 100px">
          <template #body="{ data }"><Tag :value="data.priority" :severity="getPrioritySeverity(data.priority)" /></template>
        </Column>
        <Column field="status" header="Status" style="width: 110px">
          <template #body="{ data }"><Tag :value="data.status" :severity="getStatusSeverity(data.status)" /></template>
        </Column>
        <Column field="total_amount" header="Amount" style="width: 140px">
          <template #body="{ data }">{{ formatCurrency(data.total_amount, data.currency) }}</template>
        </Column>
        <Column field="required_date" header="Required" style="width: 120px">
          <template #body="{ data }">{{ data.required_date ? formatDate(data.required_date) : '-' }}</template>
        </Column>
        <Column header="Actions" style="width: 100px">
          <template #body="{ data }">
            <Button icon="pi pi-eye" text rounded size="small" @click="viewPR(data)" />
            <Button icon="pi pi-pencil" text rounded size="small" @click="editPR(data)" :disabled="data.status !== 'DRAFT'" />
          </template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>

<style scoped>
.pr-list-page { padding: 1.5rem; }
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
.text-muted { color: var(--text-color-secondary); }
</style>
