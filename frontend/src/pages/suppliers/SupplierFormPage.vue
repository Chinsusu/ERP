<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { useAppStore } from '@/stores/app.store'
import { supplierApi } from '@/api/supplier.api'
import { useApiMutation } from '@/composables/useApi'
import type { Supplier } from '@/types/business.types'

import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'
import InputNumber from 'primevue/inputnumber'
import Checkbox from 'primevue/checkbox'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()

const supplierId = computed(() => route.params.id as string)
const isEditMode = computed(() => !!supplierId.value && supplierId.value !== 'new')

// Form state
const form = ref<Partial<Supplier>>({
  supplier_code: '',
  name: '',
  legal_name: '',
  supplier_type: 'MANUFACTURER',
  tax_id: '',
  address: '',
  city: '',
  country: '',
  phone: '',
  email: '',
  website: '',
  payment_terms: 30,
  credit_limit: undefined,
  currency: 'VND',
  status: 'PENDING',
  is_approved: false
})

const errors = ref<Record<string, string>>({})
const isSubmitting = ref(false)

// Load existing supplier
const { data: supplierData, isLoading } = useQuery({
  queryKey: ['supplier', supplierId],
  queryFn: () => supplierApi.getById(supplierId.value),
  enabled: isEditMode
})

watch(supplierData, (data) => {
  if (data?.data) {
    form.value = { ...data.data }
  }
})

const supplierTypes = [
  { label: 'Manufacturer', value: 'MANUFACTURER' },
  { label: 'Distributor', value: 'DISTRIBUTOR' },
  { label: 'Trader', value: 'TRADER' }
]

const statusOptions = [
  { label: 'Active', value: 'ACTIVE' },
  { label: 'Inactive', value: 'INACTIVE' },
  { label: 'Pending', value: 'PENDING' }
]

onMounted(() => {
  appStore.setPageTitle(isEditMode.value ? 'Edit Supplier' : 'New Supplier')
})

// Mutations
const createMutation = useApiMutation(
  (data: Partial<Supplier>) => supplierApi.create(data),
  {
    successMessage: 'Supplier created successfully',
    invalidateKeys: [['suppliers']],
    onSuccess: () => router.push('/suppliers')
  }
)

const updateMutation = useApiMutation(
  (data: Partial<Supplier>) => supplierApi.update(supplierId.value, data),
  {
    successMessage: 'Supplier updated successfully',
    invalidateKeys: [['suppliers'], ['supplier', supplierId.value]],
    onSuccess: () => router.push(`/suppliers/${supplierId.value}`)
  }
)

function validate(): boolean {
  errors.value = {}
  
  if (!form.value.supplier_code?.trim()) {
    errors.value.supplier_code = 'Supplier code is required'
  }
  if (!form.value.name?.trim()) {
    errors.value.name = 'Name is required'
  }
  if (!form.value.supplier_type) {
    errors.value.supplier_type = 'Supplier type is required'
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
  <div class="supplier-form-page">
    <!-- Header -->
    <div class="page-header">
      <div class="header-left">
        <Button icon="pi pi-arrow-left" text @click="goBack" />
        <h1 class="page-title">{{ isEditMode ? 'Edit Supplier' : 'New Supplier' }}</h1>
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
              <label for="supplier_code">Supplier Code *</label>
              <InputText 
                id="supplier_code" 
                v-model="form.supplier_code" 
                :class="{ 'p-invalid': errors.supplier_code }"
                :disabled="isEditMode"
              />
              <small class="p-error">{{ errors.supplier_code }}</small>
            </div>

            <div class="form-field">
              <label for="name">Company Name *</label>
              <InputText 
                id="name" 
                v-model="form.name" 
                :class="{ 'p-invalid': errors.name }"
              />
              <small class="p-error">{{ errors.name }}</small>
            </div>

            <div class="form-field">
              <label for="legal_name">Legal Name</label>
              <InputText id="legal_name" v-model="form.legal_name" />
            </div>

            <div class="form-field">
              <label for="supplier_type">Supplier Type *</label>
              <Dropdown 
                id="supplier_type"
                v-model="form.supplier_type" 
                :options="supplierTypes"
                optionLabel="label"
                optionValue="value"
                :class="{ 'p-invalid': errors.supplier_type }"
              />
              <small class="p-error">{{ errors.supplier_type }}</small>
            </div>

            <div class="form-field">
              <label for="tax_id">Tax ID</label>
              <InputText id="tax_id" v-model="form.tax_id" />
            </div>

            <div class="form-field">
              <label for="status">Status</label>
              <Dropdown 
                id="status"
                v-model="form.status" 
                :options="statusOptions"
                optionLabel="label"
                optionValue="value"
              />
            </div>
          </div>
        </template>
      </Card>

      <!-- Contact Info -->
      <Card class="form-card">
        <template #title>Contact Information</template>
        <template #content>
          <div class="form-grid">
            <div class="form-field full-width">
              <label for="address">Address</label>
              <Textarea id="address" v-model="form.address" rows="2" />
            </div>

            <div class="form-field">
              <label for="city">City</label>
              <InputText id="city" v-model="form.city" />
            </div>

            <div class="form-field">
              <label for="country">Country</label>
              <InputText id="country" v-model="form.country" />
            </div>

            <div class="form-field">
              <label for="phone">Phone</label>
              <InputText id="phone" v-model="form.phone" />
            </div>

            <div class="form-field">
              <label for="email">Email</label>
              <InputText id="email" v-model="form.email" type="email" />
            </div>

            <div class="form-field">
              <label for="website">Website</label>
              <InputText id="website" v-model="form.website" />
            </div>
          </div>
        </template>
      </Card>

      <!-- Business Terms -->
      <Card class="form-card">
        <template #title>Business Terms</template>
        <template #content>
          <div class="form-grid">
            <div class="form-field">
              <label for="payment_terms">Payment Terms (days)</label>
              <InputNumber id="payment_terms" v-model="form.payment_terms" :min="0" />
            </div>

            <div class="form-field">
              <label for="credit_limit">Credit Limit</label>
              <InputNumber id="credit_limit" v-model="form.credit_limit" :min="0" />
            </div>

            <div class="form-field">
              <label for="currency">Currency</label>
              <InputText id="currency" v-model="form.currency" />
            </div>

            <div class="form-field checkbox-row">
              <Checkbox id="is_approved" v-model="form.is_approved" :binary="true" />
              <label for="is_approved">Approved Supplier</label>
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
.supplier-form-page {
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
