<script setup lang="ts">
import Button from 'primevue/button'
import { useRouter } from 'vue-router'

const router = useRouter()

function goHome() {
  router.push('/dashboard')
}
</script>

<template>
  <div class="not-found-page">
    <div class="not-found-content">
      <div class="error-code">404</div>
      <h1>Page Not Found</h1>
      <p>The page you're looking for doesn't exist or has been moved.</p>
      <Button label="Go to Dashboard" icon="pi pi-home" @click="goHome" />
    </div>
  </div>
</template>

<style scoped>
.not-found-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--surface-ground);
  text-align: center;
  padding: 2rem;
}

.error-code {
  font-size: 8rem;
  font-weight: 800;
  background: linear-gradient(135deg, #7367f0, #9e95f5);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  line-height: 1;
  margin-bottom: 1rem;
}

h1 {
  font-size: 2rem;
  margin-bottom: 0.75rem;
}

p {
  color: var(--text-color-secondary);
  margin-bottom: 2rem;
}
</style>
