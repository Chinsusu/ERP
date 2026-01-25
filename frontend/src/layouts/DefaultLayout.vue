<script setup lang="ts">
import { ref, computed } from 'vue'
import { RouterView } from 'vue-router'
import { useAppStore } from '@/stores/app.store'
import AppHeader from '@/components/layout/AppHeader.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'

const appStore = useAppStore()

const sidebarCollapsed = computed(() => appStore.sidebarCollapsed)
const mainContentStyle = computed(() => ({
  marginLeft: sidebarCollapsed.value ? 'var(--sidebar-collapsed-width)' : 'var(--sidebar-width)',
  transition: 'margin-left var(--transition-duration) ease'
}))
</script>

<template>
  <div class="default-layout">
    <!-- Sidebar -->
    <AppSidebar />
    
    <!-- Main Content Area -->
    <div class="main-content" :style="mainContentStyle">
      <!-- Header -->
      <AppHeader />
      
      <!-- Page Content -->
      <main class="page-content">
        <RouterView v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </RouterView>
      </main>
      
      <!-- Footer -->
      <footer class="app-footer">
        <span>&copy; 2026 ERP Cosmetics. All rights reserved.</span>
      </footer>
    </div>
  </div>
</template>

<style scoped>
.default-layout {
  min-height: 100vh;
  background: var(--surface-ground);
}

.main-content {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.page-content {
  flex: 1;
  padding: 1.5rem;
  margin-top: var(--header-height);
  overflow-y: auto;
}

.app-footer {
  padding: 1rem 1.5rem;
  text-align: center;
  color: var(--text-color-secondary);
  font-size: 0.875rem;
  border-top: 1px solid var(--surface-border);
  background: var(--surface-card);
}

/* Page transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
