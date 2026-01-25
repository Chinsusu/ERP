<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { useAppStore } from '@/stores/app.store'
import { manufacturingApi, type BOMFilters } from '@/api/manufacturing.api'

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

onMounted(() => appStore.setPageTitle('Bill of Materials'))

const filters = ref<BOMFilters>({ page: 1, page_size: 10, search: '' })
const searchValue = ref('')
const selectedStatus = ref<string | null>(null)

const statusOptions = [
  { label: 'All Status', value: null },
  { label: 'Draft', value: 'DRAFT' },
  { label: 'Active', value: 'ACTIVE' },
  { label: 'Obsolete', value: 'OBSOLETE' }
]

const { data, isLoading, refetch } = useQuery({
  queryKey: ['bom-list', filters],
  queryFn: () => manufacturingApi.listBOM(filters.value)
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

function viewBOM(bom: any) { router.push(`/manufacturing/bom/${bom.id}`) }
function editBOM(bom: any) { router.push(`/manufacturing/bom/${bom.id}/edit`) }
function createBOM() { router.push('/manufacturing/bom/new') }

function getStatusSeverity(status: string) {
  switch (status) {
    case 'DRAFT': return 'secondary'
    case 'ACTIVE': return 'success'
    case 'OBSOLETE': return 'danger'
    default: return 'info'
  }
}
</script>

<template>
  <div class="bom-list-page">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">Bill of Materials</h1>
        <p class="page-subtitle">Product formulations and recipes</p>
      </div>
      <div class="header-right">
        <Button label="Create BOM" icon="pi pi-plus" @click="createBOM" />
      </div>
    </div>

    <div class="filters-bar">
      <IconField>
        <InputIcon class="pi pi-search" />
        <InputText v-model="searchValue" placeholder="Search BOM..." @keyup.enter="onSearch" />
      </IconField>
      <Dropdown v-model="selectedStatus" :options="statusOptions" optionLabel="label" optionValue="value" placeholder="Status" class="filter-dropdown" />
      <Button icon="pi pi-refresh" text rounded @click="() => refetch()" />
    </div>

    <div class="card">
      <DataTable :value="items" :loading="isLoading" :paginator="true" :rows="filters.page_size" :totalRecords="totalRecords" :lazy="true" :rowsPerPageOptions="[10, 25, 50]" @page="onPage" dataKey="id" stripedRows>
        <template #empty><div class="empty-message"><i class="pi pi-file"></i><p>No BOMs found</p></div></template>
        
        <Column field="bom_code" header="BOM Code" sortable style="width: 120px">
          <template #body="{ data }"><span class="font-mono">{{ data.bom_code }}</span></template>
        </Column>
        <Column field="product" header="Product" style="min-width: 200px">
          <template #body="{ data }">
            <div class="font-medium">{{ data.product?.name || '-' }}</div>
            <div class="text-sm text-muted">{{ data.product?.product_code }}</div>
          </template>
        </Column>
        <Column field="version" header="Version" style="width: 80px">
          <template #body="{ data }"><Tag :value="data.version" severity="secondary" /></template>
        </Column>
        <Column field="yield_quantity" header="Yield" style="width: 100px">
          <template #body="{ data }">{{ data.yield_quantity }} {{ data.yield_uom?.code }}</template>
        </Column>
        <Column field="items" header="Components" style="width: 100px">
          <template #body="{ data }">{{ data.items?.length || 0 }}</template>
        </Column>
        <Column field="status" header="Status" style="width: 100px">
          <template #body="{ data }"><Tag :value="data.status" :severity="getStatusSeverity(data.status)" /></template>
        </Column>
        <Column header="Actions" style="width: 100px">
          <template #body="{ data }">
            <Button icon="pi pi-eye" text rounded size="small" @click="viewBOM(data)" />
            <Button icon="pi pi-pencil" text rounded size="small" @click="editBOM(data)" :disabled="data.status === 'OBSOLETE'" />
          </template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>

<style scoped>
.bom-list-page { padding: 1.5rem; }
.page-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 1.5rem; }
.page-title { font-size: 1.5rem; font-weight: 600; margin-bottom: 0.25rem; }
.page-subtitle { color: var(--text-color-secondary); font-size: 0.875rem; }
.filters-bar { display: flex; gap: 1rem; margin-bottom: 1rem; flex-wrap: wrap; align-items: center; }
.filter-dropdown { min-width: 120px; }
.card { background: var(--surface-card); border-radius: 12px; padding: 1rem; box-shadow: 0 2px 8px rgba(0,0,0,0.08); }
.empty-message { text-align: center; padding: 3rem; color: var(--text-color-secondary); }
.empty-message i { font-size: 3rem; margin-bottom: 1rem; }
.font-mono { font-family: 'JetBrains Mono', monospace; }
.font-medium { font-weight: 500; }
.text-sm { font-size: 0.8125rem; }
.text-muted { color: var(--text-color-secondary); }
</style>
