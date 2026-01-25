<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { useAppStore } from '@/stores/app.store'
import { wmsApi, type StockFilters } from '@/api/wms.api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Tag from 'primevue/tag'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import Card from 'primevue/card'
import ProgressBar from 'primevue/progressbar'

const router = useRouter()
const appStore = useAppStore()

onMounted(() => appStore.setPageTitle('Stock Overview'))

// Stock overview stats
const { data: overviewData, isLoading: overviewLoading } = useQuery({
  queryKey: ['stock-overview'],
  queryFn: () => wmsApi.getStockOverview()
})

const overview = computed(() => overviewData.value?.data)

// Stock list
const filters = ref<StockFilters>({ page: 1, page_size: 10, search: '' })
const searchValue = ref('')
const showLowStock = ref(false)

const { data: stockData, isLoading, refetch } = useQuery({
  queryKey: ['stock', filters],
  queryFn: () => wmsApi.listStock(filters.value)
})

const stocks = computed(() => stockData.value?.data || [])
const totalRecords = computed(() => stockData.value?.pagination?.total_items || 0)

// Expiring stock
const { data: expiringData } = useQuery({
  queryKey: ['expiring-stock'],
  queryFn: () => wmsApi.getExpiringStock(30)
})

const expiringCount = computed(() => expiringData.value?.pagination?.total_items || 0)

function onPage(event: any) {
  filters.value.page = event.page + 1
  filters.value.page_size = event.rows
}

function onSearch() {
  filters.value.search = searchValue.value
  filters.value.page = 1
}

watch(showLowStock, (val) => {
  filters.value.below_minimum = val || undefined
  filters.value.page = 1
})

function viewLots(materialId: string) {
  router.push({ path: '/warehouse/lots', query: { material_id: materialId } })
}

function goToExpiring() {
  router.push('/warehouse/expiring')
}

function goToGRN() {
  router.push('/warehouse/grn')
}

function getStockLevel(stock: any) {
  if (!stock.material?.min_stock_quantity) return 100
  return Math.min((stock.available_quantity / stock.material.min_stock_quantity) * 100, 100)
}

function getStockSeverity(stock: any) {
  const level = getStockLevel(stock)
  if (level < 30) return 'danger'
  if (level < 60) return 'warning'
  return 'success'
}
</script>

<template>
  <div class="stock-overview-page">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">Stock Overview</h1>
        <p class="page-subtitle">Warehouse inventory and stock levels</p>
      </div>
      <div class="header-right">
        <Button label="Create GRN" icon="pi pi-plus" @click="goToGRN" />
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <Card class="stat-card">
        <template #content>
          <div class="stat-content">
            <i class="pi pi-box stat-icon"></i>
            <div>
              <div class="stat-value">{{ overview?.total_materials || 0 }}</div>
              <div class="stat-label">Total Materials</div>
            </div>
          </div>
        </template>
      </Card>

      <Card class="stat-card warning" @click="showLowStock = !showLowStock" style="cursor: pointer">
        <template #content>
          <div class="stat-content">
            <i class="pi pi-exclamation-triangle stat-icon"></i>
            <div>
              <div class="stat-value">{{ overview?.low_stock_count || 0 }}</div>
              <div class="stat-label">Low Stock Items</div>
            </div>
          </div>
        </template>
      </Card>

      <Card class="stat-card danger" @click="goToExpiring" style="cursor: pointer">
        <template #content>
          <div class="stat-content">
            <i class="pi pi-clock stat-icon"></i>
            <div>
              <div class="stat-value">{{ expiringCount }}</div>
              <div class="stat-label">Expiring Soon (30d)</div>
            </div>
          </div>
        </template>
      </Card>

      <Card class="stat-card">
        <template #content>
          <div class="stat-content">
            <i class="pi pi-dollar stat-icon"></i>
            <div>
              <div class="stat-value">{{ (overview?.total_value || 0).toLocaleString() }}</div>
              <div class="stat-label">Total Value (VND)</div>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Filters -->
    <div class="filters-bar">
      <IconField>
        <InputIcon class="pi pi-search" />
        <InputText v-model="searchValue" placeholder="Search materials..." @keyup.enter="onSearch" />
      </IconField>
      <Button 
        :label="showLowStock ? 'Show All' : 'Low Stock Only'" 
        :severity="showLowStock ? 'warning' : 'secondary'"
        @click="showLowStock = !showLowStock" 
      />
      <Button icon="pi pi-refresh" text rounded @click="() => refetch()" />
    </div>

    <!-- Stock Table -->
    <div class="card">
      <DataTable :value="stocks" :loading="isLoading" :paginator="true" :rows="filters.page_size" :totalRecords="totalRecords" :lazy="true" :rowsPerPageOptions="[10, 25, 50]" @page="onPage" dataKey="id" stripedRows>
        <template #empty><div class="empty-message"><i class="pi pi-inbox"></i><p>No stock data found</p></div></template>
        
        <Column field="material.material_code" header="Material Code" style="width: 120px">
          <template #body="{ data }"><span class="font-mono">{{ data.material?.material_code }}</span></template>
        </Column>
        <Column field="material.name" header="Material Name" style="min-width: 200px">
          <template #body="{ data }"><div class="font-medium">{{ data.material?.name }}</div></template>
        </Column>
        <Column field="quantity" header="On Hand" style="width: 100px">
          <template #body="{ data }">{{ data.quantity?.toLocaleString() }} {{ data.uom?.code }}</template>
        </Column>
        <Column field="reserved_quantity" header="Reserved" style="width: 100px">
          <template #body="{ data }">{{ data.reserved_quantity?.toLocaleString() }}</template>
        </Column>
        <Column field="available_quantity" header="Available" style="width: 100px">
          <template #body="{ data }">
            <span :class="{ 'text-danger': data.available_quantity < (data.material?.min_stock_quantity || 0) }">
              {{ data.available_quantity?.toLocaleString() }}
            </span>
          </template>
        </Column>
        <Column header="Stock Level" style="width: 150px">
          <template #body="{ data }">
            <ProgressBar :value="getStockLevel(data)" :showValue="false" style="height: 8px" :class="getStockSeverity(data)" />
          </template>
        </Column>
        <Column header="Actions" style="width: 80px">
          <template #body="{ data }">
            <Button icon="pi pi-list" text rounded size="small" @click="viewLots(data.material_id)" v-tooltip="'View Lots'" />
          </template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>

<style scoped>
.stock-overview-page { padding: 1.5rem; }
.page-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 1.5rem; }
.page-title { font-size: 1.5rem; font-weight: 600; margin-bottom: 0.25rem; }
.page-subtitle { color: var(--text-color-secondary); font-size: 0.875rem; }

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.stat-card { border-radius: 12px; }
.stat-card.warning { border-left: 4px solid var(--yellow-500); }
.stat-card.danger { border-left: 4px solid var(--red-500); }

.stat-content { display: flex; align-items: center; gap: 1rem; }
.stat-icon { font-size: 2rem; color: var(--primary-color); }
.stat-value { font-size: 1.5rem; font-weight: 600; }
.stat-label { color: var(--text-color-secondary); font-size: 0.875rem; }

.filters-bar { display: flex; gap: 1rem; margin-bottom: 1rem; flex-wrap: wrap; align-items: center; }
.card { background: var(--surface-card); border-radius: 12px; padding: 1rem; box-shadow: 0 2px 8px rgba(0,0,0,0.08); }
.empty-message { text-align: center; padding: 3rem; color: var(--text-color-secondary); }
.empty-message i { font-size: 3rem; margin-bottom: 1rem; }
.font-mono { font-family: 'JetBrains Mono', monospace; }
.font-medium { font-weight: 500; }
.text-danger { color: var(--red-500); font-weight: 600; }
</style>
