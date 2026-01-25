<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { useAppStore } from '@/stores/app.store'
import { supplierApi } from '@/api/supplier.api'

import Card from 'primevue/card'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Rating from 'primevue/rating'
import Dialog from 'primevue/dialog'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()

const supplierId = computed(() => route.params.id as string)
const showCertDialog = ref(false)

const { data, isLoading } = useQuery({
  queryKey: ['supplier', supplierId],
  queryFn: () => supplierApi.getById(supplierId.value),
  enabled: computed(() => !!supplierId.value)
})

const supplier = computed(() => data.value?.data)

const { data: certificationsData } = useQuery({
  queryKey: ['supplier-certifications', supplierId],
  queryFn: () => supplierApi.getCertifications(supplierId.value),
  enabled: computed(() => !!supplierId.value)
})

const { data: evaluationsData } = useQuery({
  queryKey: ['supplier-evaluations', supplierId],
  queryFn: () => supplierApi.getEvaluations(supplierId.value),
  enabled: computed(() => !!supplierId.value)
})

const certifications = computed(() => certificationsData.value?.data || [])
const evaluations = computed(() => evaluationsData.value?.data || [])

onMounted(() => {
  appStore.setPageTitle('Supplier Details')
})

function goBack() {
  router.push('/suppliers')
}

function editSupplier() {
  router.push(`/suppliers/${supplierId.value}/edit`)
}

function getStatusSeverity(status: string) {
  switch (status) {
    case 'ACTIVE': case 'VALID': return 'success'
    case 'INACTIVE': case 'EXPIRED': return 'danger'
    case 'PENDING': return 'warning'
    case 'BLOCKED': return 'danger'
    default: return 'info'
  }
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}
</script>

<template>
  <div class="supplier-detail-page">
    <!-- Header -->
    <div class="page-header">
      <div class="header-left">
        <Button icon="pi pi-arrow-left" text @click="goBack" class="back-btn" />
        <div>
          <h1 class="page-title">{{ supplier?.name || 'Loading...' }}</h1>
          <p class="page-subtitle">{{ supplier?.supplier_code }}</p>
        </div>
      </div>
      <div class="header-right">
        <Button label="Edit" icon="pi pi-pencil" @click="editSupplier" />
      </div>
    </div>

    <div v-if="isLoading" class="loading-state">
      <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
    </div>

    <template v-else-if="supplier">
      <!-- Info Cards -->
      <div class="info-grid">
        <Card class="info-card">
          <template #title>Basic Information</template>
          <template #content>
            <div class="info-row">
              <span class="info-label">Supplier Code</span>
              <span class="info-value font-mono">{{ supplier.supplier_code }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Legal Name</span>
              <span class="info-value">{{ supplier.legal_name || supplier.name }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Type</span>
              <Tag :value="supplier.supplier_type" />
            </div>
            <div class="info-row">
              <span class="info-label">Tax ID</span>
              <span class="info-value font-mono">{{ supplier.tax_id || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Status</span>
              <Tag :value="supplier.status" :severity="getStatusSeverity(supplier.status)" />
            </div>
          </template>
        </Card>

        <Card class="info-card">
          <template #title>Contact Information</template>
          <template #content>
            <div class="info-row">
              <span class="info-label">Address</span>
              <span class="info-value">{{ supplier.address || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">City</span>
              <span class="info-value">{{ supplier.city || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Country</span>
              <span class="info-value">{{ supplier.country || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Phone</span>
              <span class="info-value">{{ supplier.phone || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Email</span>
              <span class="info-value">{{ supplier.email || '-' }}</span>
            </div>
          </template>
        </Card>

        <Card class="info-card">
          <template #title>Business Terms</template>
          <template #content>
            <div class="info-row">
              <span class="info-label">Payment Terms</span>
              <span class="info-value">{{ supplier.payment_terms ? `${supplier.payment_terms} days` : '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Credit Limit</span>
              <span class="info-value">{{ supplier.credit_limit ? `${supplier.credit_limit.toLocaleString()} ${supplier.currency}` : '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Currency</span>
              <span class="info-value">{{ supplier.currency }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">Approved</span>
              <Tag :value="supplier.is_approved ? 'Yes' : 'No'" :severity="supplier.is_approved ? 'success' : 'warning'" />
            </div>
          </template>
        </Card>

        <Card class="info-card">
          <template #title>Performance Rating</template>
          <template #content>
            <div class="info-row">
              <span class="info-label">Quality</span>
              <Rating :modelValue="supplier.quality_rating || 0" readonly :cancel="false" />
            </div>
            <div class="info-row">
              <span class="info-label">Delivery</span>
              <Rating :modelValue="supplier.delivery_rating || 0" readonly :cancel="false" />
            </div>
            <div class="info-row">
              <span class="info-label">Overall</span>
              <Rating :modelValue="supplier.overall_rating || 0" readonly :cancel="false" />
            </div>
          </template>
        </Card>
      </div>

      <!-- Tabs -->
      <TabView class="mt-4">
        <TabPanel value="0" header="Contacts">
          <DataTable :value="supplier.contacts || []" stripedRows>
            <template #empty>No contacts found</template>
            <Column field="name" header="Name" />
            <Column field="position" header="Position" />
            <Column field="department" header="Department" />
            <Column field="phone" header="Phone" />
            <Column field="email" header="Email" />
            <Column field="is_primary" header="Primary">
              <template #body="{ data }">
                <Tag v-if="data.is_primary" value="Primary" severity="success" />
              </template>
            </Column>
          </DataTable>
        </TabPanel>

        <TabPanel value="0" header="Certifications">
          <div class="tab-header">
            <Button label="Add Certification" icon="pi pi-plus" size="small" @click="showCertDialog = true" />
          </div>
          <DataTable :value="certifications" stripedRows>
            <template #empty>No certifications found</template>
            <Column field="certification_type" header="Type" />
            <Column field="certificate_number" header="Certificate #" />
            <Column field="issuing_body" header="Issuing Body" />
            <Column field="issue_date" header="Issue Date">
              <template #body="{ data }">{{ formatDate(data.issue_date) }}</template>
            </Column>
            <Column field="expiry_date" header="Expiry Date">
              <template #body="{ data }">{{ formatDate(data.expiry_date) }}</template>
            </Column>
            <Column field="status" header="Status">
              <template #body="{ data }">
                <Tag :value="data.status" :severity="getStatusSeverity(data.status)" />
              </template>
            </Column>
          </DataTable>
        </TabPanel>

        <TabPanel value="0" header="Evaluations">
          <DataTable :value="evaluations" stripedRows>
            <template #empty>No evaluations found</template>
            <Column field="evaluation_date" header="Date">
              <template #body="{ data }">{{ formatDate(data.evaluation_date) }}</template>
            </Column>
            <Column field="quality_score" header="Quality">
              <template #body="{ data }">
                <Rating :modelValue="data.quality_score" readonly :cancel="false" />
              </template>
            </Column>
            <Column field="delivery_score" header="Delivery">
              <template #body="{ data }">
                <Rating :modelValue="data.delivery_score" readonly :cancel="false" />
              </template>
            </Column>
            <Column field="price_score" header="Price">
              <template #body="{ data }">
                <Rating :modelValue="data.price_score" readonly :cancel="false" />
              </template>
            </Column>
            <Column field="overall_score" header="Overall">
              <template #body="{ data }">
                <Rating :modelValue="data.overall_score" readonly :cancel="false" />
              </template>
            </Column>
            <Column field="comments" header="Comments" />
          </DataTable>
        </TabPanel>
      </TabView>
    </template>

    <!-- Certification Dialog -->
    <Dialog v-model:visible="showCertDialog" header="Add Certification" :modal="true" style="width: 500px">
      <p>Certification form coming soon...</p>
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showCertDialog = false" />
        <Button label="Save" @click="showCertDialog = false" />
      </template>
    </Dialog>
  </div>
</template>

<style scoped>
.supplier-detail-page {
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
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
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

.tab-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 1rem;
}
</style>
