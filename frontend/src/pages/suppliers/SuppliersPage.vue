<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { useAppStore } from '@/stores/app.store'
import { supplierApi, type SupplierFilters } from '@/api/supplier.api'
import type { Supplier } from '@/types/business.types'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import Tag from 'primevue/tag'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import Rating from 'primevue/rating'

const router = useRouter()
const appStore = useAppStore()

onMounted(() => {
  appStore.setPageTitle('Suppliers')
})

// Filters
const filters = ref<SupplierFilters>({
  page: 1,
  page_size: 10,
  search: '',
  supplier_type: undefined,
  status: undefined
})

const searchValue = ref('')
const selectedType = ref<string | null>(null)
const selectedStatus = ref<string | null>(null)

const supplierTypes = [
  { label: 'All Types', value: null },
  { label: 'Manufacturer', value: 'MANUFACTURER' },
  { label: 'Distributor', value: 'DISTRIBUTOR' },
  { label: 'Trader', value: 'TRADER' }
]

const statusOptions = [
  { label: 'All Status', value: null },
  { label: 'Active', value: 'ACTIVE' },
  { label: 'Inactive', value: 'INACTIVE' },
  { label: 'Pending', value: 'PENDING' },
  { label: 'Blocked', value: 'BLOCKED' }
]

// Query
const { data, isLoading, refetch } = useQuery({
  queryKey: ['suppliers', filters],
  queryFn: () => supplierApi.list(filters.value)
})

const suppliers = computed(() => data.value?.data || [])
const totalRecords = computed(() => data.value?.pagination?.total_items || 0)

// Handlers
function onPage(event: any) {
  filters.value.page = event.page + 1
  filters.value.page_size = event.rows
}

function onSort(event: any) {
  filters.value.sort_by = event.sortField
  filters.value.sort_order = event.sortOrder === 1 ? 'asc' : 'desc'
}

function onSearch() {
  filters.value.search = searchValue.value
  filters.value.page = 1
}

watch(selectedType, (val) => {
  filters.value.supplier_type = val || undefined
  filters.value.page = 1
})

watch(selectedStatus, (val) => {
  filters.value.status = val || undefined
  filters.value.page = 1
})

function viewSupplier(supplier: Supplier) {
  router.push(`/suppliers/${supplier.id}`)
}

function editSupplier(supplier: Supplier) {
  router.push(`/suppliers/${supplier.id}/edit`)
}

function createSupplier() {
  router.push('/suppliers/new')
}

function getStatusSeverity(status: string) {
  switch (status) {
    case 'ACTIVE': return 'success'
    case 'INACTIVE': return 'secondary'
    case 'PENDING': return 'warning'
    case 'BLOCKED': return 'danger'
    default: return 'info'
  }
}

function getTypeSeverity(type: string) {
  switch (type) {
    case 'MANUFACTURER': return 'info'
    case 'DISTRIBUTOR': return 'warning'
    case 'TRADER': return 'secondary'
    default: return 'info'
  }
}
</script>

<template>
  <div class="supplier-list-page">
    <!-- Header -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">Suppliers</h1>
        <p class="page-subtitle">Manage suppliers and vendor relationships</p>
      </div>
      <div class="header-right">
        <Button label="Add Supplier" icon="pi pi-plus" @click="createSupplier" />
      </div>
    </div>

    <!-- Filters -->
    <div class="filters-bar">
      <IconField>
        <InputIcon class="pi pi-search" />
        <InputText 
          v-model="searchValue" 
          placeholder="Search suppliers..." 
          @keyup.enter="onSearch"
        />
      </IconField>

      <Dropdown 
        v-model="selectedType" 
        :options="supplierTypes" 
        optionLabel="label" 
        optionValue="value"
        placeholder="Supplier Type"
        class="filter-dropdown"
      />

      <Dropdown 
        v-model="selectedStatus" 
        :options="statusOptions" 
        optionLabel="label" 
        optionValue="value"
        placeholder="Status"
        class="filter-dropdown"
      />

      <Button icon="pi pi-refresh" text rounded @click="() => refetch()" v-tooltip="'Refresh'" />
    </div>

    <!-- DataTable -->
    <div class="card">
      <DataTable
        :value="suppliers"
        :loading="isLoading"
        :paginator="true"
        :rows="filters.page_size"
        :totalRecords="totalRecords"
        :lazy="true"
        :rowsPerPageOptions="[10, 25, 50]"
        @page="onPage"
        @sort="onSort"
        dataKey="id"
        stripedRows
        showGridlines
      >
        <template #empty>
          <div class="empty-message">
            <i class="pi pi-building"></i>
            <p>No suppliers found</p>
          </div>
        </template>

        <Column field="supplier_code" header="Code" sortable style="width: 100px">
          <template #body="{ data }">
            <span class="font-mono">{{ data.supplier_code }}</span>
          </template>
        </Column>

        <Column field="name" header="Supplier Name" sortable style="min-width: 200px">
          <template #body="{ data }">
            <div>
              <div class="font-medium">{{ data.name }}</div>
              <div class="text-sm text-muted" v-if="data.legal_name && data.legal_name !== data.name">{{ data.legal_name }}</div>
            </div>
          </template>
        </Column>

        <Column field="supplier_type" header="Type" sortable style="width: 130px">
          <template #body="{ data }">
            <Tag :value="data.supplier_type" :severity="getTypeSeverity(data.supplier_type)" />
          </template>
        </Column>

        <Column field="city" header="Location" style="width: 150px">
          <template #body="{ data }">
            {{ [data.city, data.country].filter(Boolean).join(', ') || '-' }}
          </template>
        </Column>

        <Column field="overall_rating" header="Rating" style="width: 140px">
          <template #body="{ data }">
            <Rating v-if="data.overall_rating" :modelValue="data.overall_rating" readonly :cancel="false" />
            <span v-else class="text-muted">Not rated</span>
          </template>
        </Column>

        <Column field="status" header="Status" sortable style="width: 100px">
          <template #body="{ data }">
            <Tag :value="data.status" :severity="getStatusSeverity(data.status)" />
          </template>
        </Column>

        <Column header="Actions" style="width: 120px">
          <template #body="{ data }">
            <div class="action-buttons">
              <Button icon="pi pi-eye" text rounded size="small" @click="viewSupplier(data)" v-tooltip="'View'" />
              <Button icon="pi pi-pencil" text rounded size="small" @click="editSupplier(data)" v-tooltip="'Edit'" />
            </div>
          </template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>

<style scoped>
.supplier-list-page {
  padding: 1.5rem;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.5rem;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
}

.page-subtitle {
  color: var(--text-color-secondary);
  font-size: 0.875rem;
}

.filters-bar {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
  align-items: center;
}

.filter-dropdown {
  min-width: 160px;
}

.card {
  background: var(--surface-card);
  border-radius: 12px;
  padding: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.empty-message {
  text-align: center;
  padding: 3rem;
  color: var(--text-color-secondary);
}

.empty-message i {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.font-mono {
  font-family: 'JetBrains Mono', monospace;
}

.font-medium {
  font-weight: 500;
}

.text-sm {
  font-size: 0.8125rem;
}

.text-muted {
  color: var(--text-color-secondary);
}

.action-buttons {
  display: flex;
  gap: 0.25rem;
}
</style>
