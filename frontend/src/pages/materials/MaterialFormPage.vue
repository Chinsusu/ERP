<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { useToast } from 'primevue/usetoast'
import { useAppStore } from '@/stores/app.store'
import { materialApi } from '@/api/material.api'
import { useApiMutation } from '@/composables/useApi'
import type { Material, Category, UnitOfMeasure } from '@/types/business.types'

import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'
import InputNumber from 'primevue/inputnumber'
import Checkbox from 'primevue/checkbox'
import Divider from 'primevue/divider'

const route = useRoute()
const router = useRouter()
const toast = useToast()
const appStore = useAppStore()

const materialId = computed(() => route.params.id as string)
const isEditMode = computed(() => !!materialId.value && materialId.value !== 'new')

// Form state
const form = ref<Partial<Material>>({
  material_code: '',
  name: '',
  description: '',
  material_type: 'RAW_MATERIAL',
  category_id: undefined,
  inci_name: '',
  cas_number: '',
  origin_country: '',
  is_organic: false,
  is_natural: false,
  storage_conditions: '',
  requires_cold_storage: false,
  shelf_life_days: undefined,
  is_hazardous: false,
  allergen_info: '',
  base_uom_id: undefined,
  min_stock_quantity: undefined,
  max_stock_quantity: undefined,
  reorder_point: undefined,
  standard_cost: undefined,
  currency: 'VND',
  is_active: true
})

const errors = ref<Record<string, string>>({})
const isSubmitting = ref(false)

// Load existing material for edit
const { data: materialData, isLoading } = useQuery({
  queryKey: ['material', materialId],
  queryFn: () => materialApi.getById(materialId.value),
  enabled: isEditMode
})

watch(materialData, (data) => {
  if (data?.data) {
    form.value = { ...data.data }
  }
})

// Load categories and UoMs
const { data: categoriesData } = useQuery({
  queryKey: ['categories', 'MATERIAL'],
  queryFn: () => materialApi.listCategories('MATERIAL')
})

const { data: uomData } = useQuery({
  queryKey: ['uom'],
  queryFn: () => materialApi.listUom()
})

const categories = computed(() => categoriesData.value?.data || [])
const units = computed(() => uomData.value?.data || [])

const materialTypes = [
  { label: 'Raw Material', value: 'RAW_MATERIAL' },
  { label: 'Packaging', value: 'PACKAGING' },
  { label: 'Semi-Finished', value: 'SEMI_FINISHED' }
]

onMounted(() => {
  appStore.setPageTitle(isEditMode.value ? 'Edit Material' : 'New Material')
})

// Mutations
const createMutation = useApiMutation(
  (data: Partial<Material>) => materialApi.create(data),
  {
    successMessage: 'Material created successfully',
    invalidateKeys: [['materials']],
    onSuccess: () => router.push('/materials')
  }
)

const updateMutation = useApiMutation(
  (data: Partial<Material>) => materialApi.update(materialId.value, data),
  {
    successMessage: 'Material updated successfully',
    invalidateKeys: [['materials'], ['material', materialId.value]],
    onSuccess: () => router.push(`/materials/${materialId.value}`)
  }
)

function validate(): boolean {
  errors.value = {}
  
  if (!form.value.material_code?.trim()) {
    errors.value.material_code = 'Material code is required'
  }
  if (!form.value.name?.trim()) {
    errors.value.name = 'Name is required'
  }
  if (!form.value.material_type) {
    errors.value.material_type = 'Material type is required'
  }
  if (!form.value.base_uom_id) {
    errors.value.base_uom_id = 'Base UoM is required'
  }
  
  return Object.keys(errors.value).length === 0
}

async function onSubmit() {
  if (!validate()) return
  
  isSubmitting.value = true
  try {
    if (isEditMode.value) {
      await updateMutation.mutateAsync(form.value)
    } else {
      await createMutation.mutateAsync(form.value)
    }
  } finally {
    isSubmitting.value = false
  }
}

function goBack() {
  router.back()
}
</script>

<template>
  <div class="material-form-page">
    <!-- Header -->
    <div class="page-header">
      <div class="header-left">
        <Button icon="pi pi-arrow-left" text @click="goBack" />
        <h1 class="page-title">{{ isEditMode ? 'Edit Material' : 'New Material' }}</h1>
      </div>
    </div>

    <div v-if="isLoading" class="loading-state">
      <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
    </div>

    <form v-else @submit.prevent="onSubmit" class="form-container">
      <!-- Basic Info -->
      <Card class="form-card">
        <template #title>Basic Information</template>
        <template #content>
          <div class="form-grid">
            <div class="form-field">
              <label for="material_code">Material Code *</label>
              <InputText 
                id="material_code" 
                v-model="form.material_code" 
                :class="{ 'p-invalid': errors.material_code }"
                :disabled="isEditMode"
              />
              <small class="p-error">{{ errors.material_code }}</small>
            </div>

            <div class="form-field">
              <label for="name">Name *</label>
              <InputText 
                id="name" 
                v-model="form.name" 
                :class="{ 'p-invalid': errors.name }"
              />
              <small class="p-error">{{ errors.name }}</small>
            </div>

            <div class="form-field">
              <label for="material_type">Material Type *</label>
              <Dropdown 
                id="material_type"
                v-model="form.material_type" 
                :options="materialTypes"
                optionLabel="label"
                optionValue="value"
                :class="{ 'p-invalid': errors.material_type }"
              />
              <small class="p-error">{{ errors.material_type }}</small>
            </div>

            <div class="form-field">
              <label for="category_id">Category</label>
              <Dropdown 
                id="category_id"
                v-model="form.category_id" 
                :options="categories"
                optionLabel="name"
                optionValue="id"
                placeholder="Select category"
                showClear
              />
            </div>

            <div class="form-field full-width">
              <label for="description">Description</label>
              <Textarea id="description" v-model="form.description" rows="3" />
            </div>
          </div>
        </template>
      </Card>

      <!-- Specifications -->
      <Card class="form-card">
        <template #title>Specifications</template>
        <template #content>
          <div class="form-grid">
            <div class="form-field">
              <label for="inci_name">INCI Name</label>
              <InputText id="inci_name" v-model="form.inci_name" />
            </div>

            <div class="form-field">
              <label for="cas_number">CAS Number</label>
              <InputText id="cas_number" v-model="form.cas_number" />
            </div>

            <div class="form-field">
              <label for="origin_country">Origin Country</label>
              <InputText id="origin_country" v-model="form.origin_country" />
            </div>

            <div class="form-field checkbox-row">
              <Checkbox id="is_organic" v-model="form.is_organic" :binary="true" />
              <label for="is_organic">Organic</label>
            </div>

            <div class="form-field checkbox-row">
              <Checkbox id="is_natural" v-model="form.is_natural" :binary="true" />
              <label for="is_natural">Natural</label>
            </div>
          </div>
        </template>
      </Card>

      <!-- Storage & Safety -->
      <Card class="form-card">
        <template #title>Storage & Safety</template>
        <template #content>
          <div class="form-grid">
            <div class="form-field">
              <label for="storage_conditions">Storage Conditions</label>
              <InputText id="storage_conditions" v-model="form.storage_conditions" />
            </div>

            <div class="form-field">
              <label for="shelf_life_days">Shelf Life (days)</label>
              <InputNumber id="shelf_life_days" v-model="form.shelf_life_days" :min="0" />
            </div>

            <div class="form-field checkbox-row">
              <Checkbox id="requires_cold_storage" v-model="form.requires_cold_storage" :binary="true" />
              <label for="requires_cold_storage">Requires Cold Storage</label>
            </div>

            <div class="form-field checkbox-row">
              <Checkbox id="is_hazardous" v-model="form.is_hazardous" :binary="true" />
              <label for="is_hazardous">Hazardous Material</label>
            </div>

            <div class="form-field full-width">
              <label for="allergen_info">Allergen Information</label>
              <Textarea id="allergen_info" v-model="form.allergen_info" rows="2" />
            </div>
          </div>
        </template>
      </Card>

      <!-- Inventory Settings -->
      <Card class="form-card">
        <template #title>Inventory Settings</template>
        <template #content>
          <div class="form-grid">
            <div class="form-field">
              <label for="base_uom_id">Base UoM *</label>
              <Dropdown 
                id="base_uom_id"
                v-model="form.base_uom_id" 
                :options="units"
                optionLabel="name"
                optionValue="id"
                placeholder="Select unit"
                :class="{ 'p-invalid': errors.base_uom_id }"
              />
              <small class="p-error">{{ errors.base_uom_id }}</small>
            </div>

            <div class="form-field">
              <label for="standard_cost">Standard Cost</label>
              <InputNumber id="standard_cost" v-model="form.standard_cost" :minFractionDigits="0" :maxFractionDigits="2" />
            </div>

            <div class="form-field">
              <label for="min_stock_quantity">Min Stock Quantity</label>
              <InputNumber id="min_stock_quantity" v-model="form.min_stock_quantity" :minFractionDigits="0" :maxFractionDigits="3" />
            </div>

            <div class="form-field">
              <label for="max_stock_quantity">Max Stock Quantity</label>
              <InputNumber id="max_stock_quantity" v-model="form.max_stock_quantity" :minFractionDigits="0" :maxFractionDigits="3" />
            </div>

            <div class="form-field">
              <label for="reorder_point">Reorder Point</label>
              <InputNumber id="reorder_point" v-model="form.reorder_point" :minFractionDigits="0" :maxFractionDigits="3" />
            </div>

            <div class="form-field checkbox-row">
              <Checkbox id="is_active" v-model="form.is_active" :binary="true" />
              <label for="is_active">Active</label>
            </div>
          </div>
        </template>
      </Card>

      <!-- Actions -->
      <div class="form-actions">
        <Button label="Cancel" severity="secondary" @click="goBack" />
        <Button 
          type="submit" 
          :label="isEditMode ? 'Update' : 'Create'" 
          :loading="isSubmitting"
        />
      </div>
    </form>
  </div>
</template>

<style scoped>
.material-form-page {
  padding: 1.5rem;
  max-width: 1000px;
  margin: 0 auto;
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

.page-title {
  font-size: 1.5rem;
  font-weight: 600;
}

.loading-state {
  display: flex;
  justify-content: center;
  padding: 3rem;
}

.form-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-card {
  background: var(--surface-card);
  border-radius: 12px;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-field.full-width {
  grid-column: 1 / -1;
}

.form-field label {
  font-weight: 500;
  font-size: 0.875rem;
}

.form-field.checkbox-row {
  flex-direction: row;
  align-items: center;
  gap: 0.5rem;
}

.form-field.checkbox-row label {
  font-weight: 400;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding-top: 1rem;
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
