<script setup lang="ts">
import { onMounted } from 'vue'
import { useAppStore } from '@/stores/app.store'
const appStore = useAppStore()
onMounted(() => appStore.setPageTitle('Work Orders'))
</script>
<template>
  <div class="page-container">
    <div class="page-header"><h1 class="page-title">Work Orders</h1></div>
    <div class="card"><p>Work orders management page - Coming soon</p></div>
  </div>
</template>
