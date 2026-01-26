<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { useAppStore } from '@/stores/app.store'
import { materialApi } from '@/api/material.api'
import type { Product } from '@/types/business.types'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Tag from 'primevue/tag'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'

const router = useRouter()
const appStore = useAppStore()

onMounted(() => {
  appStore.setPageTitle('Products')
})

const filters = ref({
  page: 1,
  page_size: 10,
  search: ''
})

const searchValue = ref('')

// Query using materialApi.listProducts
const { data, isLoading, refetch } = useQuery({
  queryKey: ['products', filters],
  queryFn: () => materialApi.listProducts(filters.value)
})

const products = computed(() => data.value?.data || [])
const totalRecords = computed(() => data.value?.pagination?.total_items || 0)

function onPage(event: any) {
  filters.value.page = event.page + 1
  filters.value.page_size = event.rows
}

function onSearch() {
  filters.value.search = searchValue.value
  filters.value.page = 1
}

function viewProduct(product: Product) {
  router.push(`/products/${product.id}`)
}

function createProduct() {
  router.push('/products/new')
}
</script>

<template>
  <div class="product-list-page">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">Products</h1>
        <p class="page-subtitle">Manage finished goods and sales items</p>
      </div>
      <div class="header-right">
        <Button label="Add Product" icon="pi pi-plus" @click="createProduct" />
      </div>
    </div>

    <div class="filters-bar">
      <IconField>
        <InputIcon class="pi pi-search" />
        <InputText 
          v-model="searchValue" 
          placeholder="Search products..." 
          @keyup.enter="onSearch"
        />
      </IconField>
      <Button icon="pi pi-refresh" text rounded @click="() => refetch()" v-tooltip="'Refresh'" />
    </div>

    <div class="card">
      <DataTable
        :value="products"
        :loading="isLoading"
        :paginator="true"
        :rows="filters.page_size"
        :totalRecords="totalRecords"
        :lazy="true"
        @page="onPage"
        stripedRows
        showGridlines
      >
        <template #empty>
          <div class="empty-message">
            <i class="pi pi-inbox"></i>
            <p>No products found</p>
          </div>
        </template>

        <Column field="product_code" header="Code" style="width: 150px">
          <template #body="{ data }">
            <span class="font-mono">{{ data.product_code }}</span>
          </template>
        </Column>

        <Column field="name" header="Name" style="min-width: 200px">
          <template #body="{ data }">
            <div class="font-medium">{{ data.name }}</div>
          </template>
        </Column>

        <Column field="category" header="Category" style="width: 200px">
          <template #body="{ data }">
            {{ data.category?.name || '-' }}
          </template>
        </Column>

        <Column field="standard_price" header="Price" style="width: 120px">
          <template #body="{ data }">
            {{ data.standard_price ? new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(data.standard_price) : '-' }}
          </template>
        </Column>

        <Column field="is_active" header="Status" style="width: 100px">
          <template #body="{ data }">
            <Tag :value="data.is_active ? 'Active' : 'Inactive'" :severity="data.is_active ? 'success' : 'danger'" />
          </template>
        </Column>

        <Column header="Actions" style="width: 100px">
          <template #body="{ data }">
            <Button icon="pi pi-eye" text rounded size="small" @click="viewProduct(data)" />
          </template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>

<style scoped>
.product-list-page { padding: 1.5rem; }
.page-header { display: flex; justify-content: space-between; margin-bottom: 1.5rem; }
.page-title { font-size: 1.5rem; font-weight: 600; margin-bottom: 0.25rem; }
.page-subtitle { color: #6e6b7b; font-size: 0.875rem; }
.filters-bar { display: flex; gap: 1rem; margin-bottom: 1rem; }
.card { background: white; border-radius: 12px; padding: 1rem; box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08); }
.empty-message { text-align: center; padding: 3rem; color: #6e6b7b; }
.empty-message i { font-size: 3rem; margin-bottom: 1rem; }
.font-mono { font-family: monospace; }
.font-medium { font-weight: 500; }
</style>
