<script setup lang="ts">
import { onMounted } from 'vue'
import { useAppStore } from '@/stores/app.store'
const appStore = useAppStore()
onMounted(() => appStore.setPageTitle('Sales Orders'))
</script>
<template>
  <div class="page-container">
    <div class="page-header"><h1 class="page-title">Sales Orders</h1></div>
    <div class="card"><p>Sales orders management page - Coming soon</p></div>
  </div>
</template>
