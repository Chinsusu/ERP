<script setup lang="ts">
import { onMounted } from 'vue'
import { useAppStore } from '@/stores/app.store'
const appStore = useAppStore()
onMounted(() => appStore.setPageTitle('KOL Management'))
</script>
<template>
  <div class="page-container">
    <div class="page-header"><h1 class="page-title">KOL Management</h1></div>
    <div class="card"><p>KOL/Influencer management page - Coming soon</p></div>
  </div>
</template>
