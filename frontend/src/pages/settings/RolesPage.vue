<script setup lang="ts">
import { onMounted } from 'vue'
import { useAppStore } from '@/stores/app.store'
const appStore = useAppStore()
onMounted(() => appStore.setPageTitle('Roles'))
</script>
<template>
  <div class="page-container">
    <div class="page-header"><h1 class="page-title">Roles & Permissions</h1></div>
    <div class="card"><p>Roles management page - Coming soon</p></div>
  </div>
</template>
