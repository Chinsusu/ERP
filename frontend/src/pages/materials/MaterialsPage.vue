<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { useAppStore } from '@/stores/app.store'
import { materialApi, type MaterialFilters } from '@/api/material.api'
import type { Material } from '@/types/business.types'

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

onMounted(() => {
  appStore.setPageTitle('Materials')
})

// Filters
const filters = ref<MaterialFilters>({
  page: 1,
  page_size: 10,
  search: '',
  material_type: undefined,
  is_active: undefined
})

const searchValue = ref('')
const selectedType = ref<string | null>(null)
const selectedStatus = ref<boolean | null>(null)

const materialTypes = [
  { label: 'All Types', value: null },
  { label: 'Raw Material', value: 'RAW_MATERIAL' },
  { label: 'Packaging', value: 'PACKAGING' },
  { label: 'Semi-Finished', value: 'SEMI_FINISHED' }
]

const statusOptions = [
  { label: 'All Status', value: null },
  { label: 'Active', value: true },
  { label: 'Inactive', value: false }
]

// Query
const { data, isLoading, refetch } = useQuery({
  queryKey: ['materials', filters],
  queryFn: () => materialApi.list(filters.value)
})

const materials = computed(() => data.value?.data || [])
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
  filters.value.material_type = val || undefined
  filters.value.page = 1
})

watch(selectedStatus, (val) => {
  filters.value.is_active = val ?? undefined
  filters.value.page = 1
})

function viewMaterial(material: Material) {
  router.push(`/materials/${material.id}`)
}

function editMaterial(material: Material) {
  router.push(`/materials/${material.id}/edit`)
}

function createMaterial() {
  router.push('/materials/new')
}

function getTypeSeverity(type: string) {
  switch (type) {
    case 'RAW_MATERIAL': return 'info'
    case 'PACKAGING': return 'warning'
    case 'SEMI_FINISHED': return 'secondary'
    default: return 'info'
  }
}

function getTypeLabel(type: string) {
  switch (type) {
    case 'RAW_MATERIAL': return 'Raw Material'
    case 'PACKAGING': return 'Packaging'
    case 'SEMI_FINISHED': return 'Semi-Finished'
    default: return type
  }
}
</script>

<template>
  <div class="material-list-page">
    <!-- Header -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">Materials</h1>
        <p class="page-subtitle">Manage raw materials, packaging, and semi-finished goods</p>
      </div>
      <div class="header-right">
        <Button label="Add Material" icon="pi pi-plus" @click="createMaterial" />
      </div>
    </div>

    <!-- Filters -->
    <div class="filters-bar">
      <IconField>
        <InputIcon class="pi pi-search" />
        <InputText 
          v-model="searchValue" 
          placeholder="Search materials..." 
          @keyup.enter="onSearch"
        />
      </IconField>

      <Dropdown 
        v-model="selectedType" 
        :options="materialTypes" 
        optionLabel="label" 
        optionValue="value"
        placeholder="Material Type"
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
        :value="materials"
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
        tableStyle="min-width: 50rem"
      >
        <template #empty>
          <div class="empty-message">
            <i class="pi pi-inbox"></i>
            <p>No materials found</p>
          </div>
        </template>

        <Column field="material_code" header="Code" sortable style="width: 120px">
          <template #body="{ data }">
            <span class="font-mono">{{ data.material_code }}</span>
          </template>
        </Column>

        <Column field="name" header="Name" sortable style="min-width: 200px">
          <template #body="{ data }">
            <div>
              <div class="font-medium">{{ data.name }}</div>
              <div class="text-sm text-muted" v-if="data.inci_name">{{ data.inci_name }}</div>
            </div>
          </template>
        </Column>

        <Column field="material_type" header="Type" sortable style="width: 140px">
          <template #body="{ data }">
            <Tag :value="getTypeLabel(data.material_type)" :severity="getTypeSeverity(data.material_type)" />
          </template>
        </Column>

        <Column field="category" header="Category" style="width: 150px">
          <template #body="{ data }">
            {{ data.category?.name || '-' }}
          </template>
        </Column>

        <Column field="base_uom" header="UoM" style="width: 80px">
          <template #body="{ data }">
            {{ data.base_uom?.code || '-' }}
          </template>
        </Column>

        <Column field="is_active" header="Status" style="width: 100px">
          <template #body="{ data }">
            <Tag 
              :value="data.is_active ? 'Active' : 'Inactive'" 
              :severity="data.is_active ? 'success' : 'danger'" 
            />
          </template>
        </Column>

        <Column header="Actions" style="width: 120px">
          <template #body="{ data }">
            <div class="action-buttons">
              <Button icon="pi pi-eye" text rounded size="small" @click="viewMaterial(data)" v-tooltip="'View'" />
              <Button icon="pi pi-pencil" text rounded size="small" @click="editMaterial(data)" v-tooltip="'Edit'" />
            </div>
          </template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>

<style scoped>
.material-list-page {
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
