<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { useAppStore } from '@/stores/app.store'
import { materialApi } from '@/api/material.api'

import Card from 'primevue/card'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()

const materialId = computed(() => route.params.id as string)

const { data, isLoading, error } = useQuery({
  queryKey: ['material', materialId],
  queryFn: () => materialApi.getById(materialId.value),
  enabled: computed(() => !!materialId.value)
})

const material = computed(() => data.value?.data)

const { data: suppliersData } = useQuery({
  queryKey: ['material-suppliers', materialId],
  queryFn: () => materialApi.getSuppliers(materialId.value),
  enabled: computed(() => !!materialId.value)
})

const suppliers = computed(() => suppliersData.value?.data || [])

onMounted(() => {
  appStore.setPageTitle('Material Details')
})

function goBack() {
  router.push('/materials')
}

function editMaterial() {
  router.push(`/materials/${materialId.value}/edit`)
}

function getTypeSeverity(type: string) {
  switch (type) {
    case 'RAW_MATERIAL': return 'info'
    case 'PACKAGING': return 'warning'
    case 'SEMI_FINISHED': return 'secondary'
    default: return 'info'
  }
}
</script>

<template>
  <div class="material-detail-page">
    <!-- Header -->
    <div class="page-header">
      <div class="header-left">
        <Button icon="pi pi-arrow-left" text @click="goBack" class="back-btn" />
        <div>
          <h1 class="page-title">{{ material?.name || 'Loading...' }}</h1>
          <p class="page-subtitle">{{ material?.material_code }}</p>
        </div>
      </div>
      <div class="header-right">
        <Button label="Edit" icon="pi pi-pencil" @click="editMaterial" />
      </div>
    </div>

    <div v-if="isLoading" class="loading-state">
      <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
    </div>

    <template v-else-if="material">
      <!-- Info Cards -->
      <div class="info-grid">
        <Card class="info-card">
          <template #title>Basic Information</template>
          <template #content>
            <div class="info-row">
              <span class="info-label">Material Code</span>
              <span class="info-value font-mono">{{ material.material_code }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Name</span>
              <span class="info-value">{{ material.name }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Type</span>
              <Tag :value="material.material_type" :severity="getTypeSeverity(material.material_type)" />
            </div>
            <div class="info-row">
              <span class="info-label">Category</span>
              <span class="info-value">{{ material.category?.name || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Status</span>
              <Tag :value="material.is_active ? 'Active' : 'Inactive'" :severity="material.is_active ? 'success' : 'danger'" />
            </div>
          </template>
        </Card>

        <Card class="info-card">
          <template #title>Specifications</template>
          <template #content>
            <div class="info-row">
              <span class="info-label">INCI Name</span>
              <span class="info-value">{{ material.inci_name || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">CAS Number</span>
              <span class="info-value font-mono">{{ material.cas_number || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Origin Country</span>
              <span class="info-value">{{ material.origin_country || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Organic</span>
              <Tag :value="material.is_organic ? 'Yes' : 'No'" :severity="material.is_organic ? 'success' : 'secondary'" />
            </div>
            <div class="info-row">
              <span class="info-label">Natural</span>
              <Tag :value="material.is_natural ? 'Yes' : 'No'" :severity="material.is_natural ? 'success' : 'secondary'" />
            </div>
          </template>
        </Card>

        <Card class="info-card">
          <template #title>Storage & Safety</template>
          <template #content>
            <div class="info-row">
              <span class="info-label">Storage Conditions</span>
              <span class="info-value">{{ material.storage_conditions || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Cold Storage</span>
              <Tag :value="material.requires_cold_storage ? 'Required' : 'Not Required'" :severity="material.requires_cold_storage ? 'warning' : 'secondary'" />
            </div>
            <div class="info-row">
              <span class="info-label">Shelf Life</span>
              <span class="info-value">{{ material.shelf_life_days ? `${material.shelf_life_days} days` : '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Hazardous</span>
              <Tag :value="material.is_hazardous ? 'Yes' : 'No'" :severity="material.is_hazardous ? 'danger' : 'success'" />
            </div>
          </template>
        </Card>

        <Card class="info-card">
          <template #title>Inventory Settings</template>
          <template #content>
            <div class="info-row">
              <span class="info-label">Base UoM</span>
              <span class="info-value">{{ material.base_uom?.name }} ({{ material.base_uom?.code }})</span>
            </div>
            <div class="info-row">
              <span class="info-label">Min Stock</span>
              <span class="info-value">{{ material.min_stock_quantity || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Max Stock</span>
              <span class="info-value">{{ material.max_stock_quantity || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Reorder Point</span>
              <span class="info-value">{{ material.reorder_point || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Standard Cost</span>
              <span class="info-value">{{ material.standard_cost ? `${material.standard_cost.toLocaleString()} ${material.currency}` : '-' }}</span>
            </div>
          </template>
        </Card>
      </div>

      <!-- Tabs -->
      <TabView class="mt-4">
        <TabPanel value="0" header="Approved Suppliers">
          <DataTable :value="suppliers" stripedRows>
            <template #empty>No approved suppliers</template>
            <Column field="supplier.name" header="Supplier" />
            <Column field="supplier_material_code" header="Supplier Code" />
            <Column field="unit_price" header="Unit Price">
              <template #body="{ data }">
                {{ data.unit_price?.toLocaleString() }} {{ data.currency }}
              </template>
            </Column>
            <Column field="lead_time_days" header="Lead Time">
              <template #body="{ data }">
                {{ data.lead_time_days ? `${data.lead_time_days} days` : '-' }}
              </template>
            </Column>
            <Column field="is_preferred" header="Preferred">
              <template #body="{ data }">
                <Tag :value="data.is_preferred ? 'Yes' : 'No'" :severity="data.is_preferred ? 'success' : 'secondary'" />
              </template>
            </Column>
          </DataTable>
        </TabPanel>
        <TabPanel value="0" header="Description">
          <p>{{ material.description || 'No description available' }}</p>
        </TabPanel>
        <TabPanel value="0" header="Allergen Info">
          <p>{{ material.allergen_info || 'No allergen information available' }}</p>
        </TabPanel>
      </TabView>
    </template>
  </div>
</template>

<style scoped>
.material-detail-page {
  padding: 1.5rem;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.back-btn {
  margin-right: 0.5rem;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
}

.page-subtitle {
  color: var(--text-color-secondary);
  font-family: 'JetBrains Mono', monospace;
}

.loading-state {
  display: flex;
  justify-content: center;
  padding: 3rem;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1rem;
}

.info-card {
  background: var(--surface-card);
  border-radius: 12px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0;
  border-bottom: 1px solid var(--surface-border);
}

.info-row:last-child {
  border-bottom: none;
}

.info-label {
  color: var(--text-color-secondary);
  font-size: 0.875rem;
}

.info-value {
  font-weight: 500;
}

.font-mono {
  font-family: 'JetBrains Mono', monospace;
}

.mt-4 {
  margin-top: 1.5rem;
}
</style>
