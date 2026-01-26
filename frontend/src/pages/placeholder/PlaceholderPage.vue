<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import Button from 'primevue/button'

const route = useRoute()
const title = computed(() => (route.meta.title as string) || 'Module')
</script>

<template>
  <div class="placeholder-page">
    <div class="content card">
      <div class="illustration">
        <i class="pi pi-construct" style="font-size: 5rem; color: #7367f0"></i>
      </div>
      <h1 class="title">{{ title }} is under development</h1>
      <p class="subtitle">
        We are working hard to bring you the best experience. 
        This module will be available in the next update.
      </p>
      <div class="actions">
        <Button label="Back to Dashboard" icon="pi pi-home" @click="$router.push('/dashboard')" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.placeholder-page {
  padding: 2rem;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
}

.card {
  background: white;
  padding: 4rem 2rem;
  border-radius: 12px;
  text-align: center;
  max-width: 600px;
  width: 100%;
  box-shadow: 0 4px 24px 0 rgba(34, 41, 47, 0.1);
}

.illustration {
  margin-bottom: 2rem;
}

.title {
  font-size: 1.75rem;
  font-weight: 600;
  color: #5e5873;
  margin-bottom: 1rem;
}

.subtitle {
  color: #6e6b7b;
  margin-bottom: 2.5rem;
  line-height: 1.6;
}
</style>
