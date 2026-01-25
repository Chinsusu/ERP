<script setup lang="ts">
import { onMounted } from 'vue'
import { useAppStore } from '@/stores/app.store'
const appStore = useAppStore()
onMounted(() => appStore.setPageTitle('Goods Receipt'))
</script>
<template>
  <div class="page-container">
    <div class="page-header"><h1 class="page-title">Goods Receipt Notes</h1></div>
    <div class="card"><p>GRN management page - Coming soon</p></div>
  </div>
</template>
