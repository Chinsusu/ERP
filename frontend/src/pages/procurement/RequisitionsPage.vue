<script setup lang="ts">
import { onMounted } from 'vue'
import { useAppStore } from '@/stores/app.store'

const appStore = useAppStore()

onMounted(() => {
  appStore.setPageTitle('Purchase Requisitions')
})
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">Purchase Requisitions</h1>
    </div>
    <div class="card">
      <p>Purchase requisitions page - Coming soon</p>
    </div>
  </div>
</template>
