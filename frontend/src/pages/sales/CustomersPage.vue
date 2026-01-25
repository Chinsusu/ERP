<script setup lang="ts">
import { onMounted } from 'vue'
import { useAppStore } from '@/stores/app.store'
const appStore = useAppStore()
onMounted(() => appStore.setPageTitle('Customers'))
</script>
<template>
  <div class="page-container">
    <div class="page-header"><h1 class="page-title">Customers</h1></div>
    <div class="card"><p>Customer management page - Coming soon</p></div>
  </div>
</template>
