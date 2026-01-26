<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAppStore } from '@/stores/app.store'
import { useAuthStore } from '@/stores/auth.store'
import Card from 'primevue/card'

const appStore = useAppStore()
const authStore = useAuthStore()

const user = computed(() => authStore.user)

onMounted(() => {
  appStore.setPageTitle('Dashboard')
})

// Mock data for dashboard stats - Vuexy color palette
const stats = [
  { label: 'Active Orders', value: '127', icon: 'pi pi-shopping-cart', color: '#7367f0', trend: '+12%' },
  { label: 'Low Stock Items', value: '23', icon: 'pi pi-exclamation-triangle', color: '#ff9f43', trend: '-5%' },
  { label: 'Work Orders', value: '45', icon: 'pi pi-cog', color: '#00cfe8', trend: '+8%' },
  { label: 'Pending Approvals', value: '8', icon: 'pi pi-clock', color: '#ea5455', trend: '0%' }
]

const recentActivities = [
  { action: 'PO-2026-0089 confirmed', user: 'Nguyen Van A', time: '5 minutes ago', type: 'success' },
  { action: 'GRN-2026-0156 completed', user: 'Tran Thi B', time: '15 minutes ago', type: 'success' },
  { action: 'Low stock alert: RM-0025', user: 'System', time: '30 minutes ago', type: 'warning' },
  { action: 'Work Order WO-2026-0034 started', user: 'Le Van C', time: '1 hour ago', type: 'info' },
  { action: 'New supplier approved: ABC Corp', user: 'Admin', time: '2 hours ago', type: 'success' }
]
</script>

<template>
  <div class="dashboard-page">
    <!-- Welcome Section -->
    <div class="welcome-section">
      <div class="welcome-text">
        <h1>Welcome back, {{ user?.first_name || 'User' }}!</h1>
        <p>Here's what's happening with your business today.</p>
      </div>
      <div class="date-display">
        <i class="pi pi-calendar"></i>
        {{ new Date().toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="stats-grid">
      <div v-for="stat in stats" :key="stat.label" class="stat-card">
        <div class="stat-icon" :style="{ background: stat.color + '20', color: stat.color }">
          <i :class="stat.icon"></i>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stat.value }}</div>
          <div class="stat-label">{{ stat.label }}</div>
        </div>
        <div class="stat-trend" :class="{ positive: stat.trend.startsWith('+'), negative: stat.trend.startsWith('-') }">
          {{ stat.trend }}
        </div>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="content-grid">
      <!-- Recent Activities -->
      <Card class="activities-card">
        <template #title>
          <div class="card-header">
            <span>Recent Activities</span>
            <a href="#" class="view-all">View All</a>
          </div>
        </template>
        <template #content>
          <div class="activities-list">
            <div v-for="(activity, index) in recentActivities" :key="index" class="activity-item">
              <div class="activity-indicator" :class="activity.type"></div>
              <div class="activity-content">
                <div class="activity-action">{{ activity.action }}</div>
                <div class="activity-meta">
                  <span class="activity-user">{{ activity.user }}</span>
                  <span class="activity-time">{{ activity.time }}</span>
                </div>
              </div>
            </div>
          </div>
        </template>
      </Card>

      <!-- Quick Actions -->
      <Card class="quick-actions-card">
        <template #title>Quick Actions</template>
        <template #content>
          <div class="actions-grid">
            <button class="action-btn">
              <i class="pi pi-plus"></i>
              <span>New PO</span>
            </button>
            <button class="action-btn">
              <i class="pi pi-inbox"></i>
              <span>Receive Goods</span>
            </button>
            <button class="action-btn">
              <i class="pi pi-file"></i>
              <span>Create BOM</span>
            </button>
            <button class="action-btn">
              <i class="pi pi-users"></i>
              <span>Add Customer</span>
            </button>
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>

<style scoped>
/* Vuexy Dashboard */
.dashboard-page {
  padding: 1.5rem;
  max-width: 1400px;
  margin: 0 auto;
}

/* Welcome Section - Vuexy Style */
.welcome-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding: 1.5rem;
  background: var(--surface-card);
  border-radius: 10px;
  box-shadow: 0 4px 24px 0 rgba(34, 41, 47, 0.1);
}

.welcome-text h1 {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
  color: var(--text-color);
}

.welcome-text p {
  color: var(--text-color-secondary);
  font-size: 0.9375rem;
  margin: 0;
}

.date-display {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--text-color-secondary);
  font-size: 0.875rem;
  background: var(--surface-ground);
  padding: 0.5rem 1rem;
  border-radius: 6px;
}

/* Stats Grid - Vuexy Cards */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 576px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}

/* Stat Card - Vuexy Style */
.stat-card {
  background: var(--surface-card);
  border-radius: 10px;
  padding: 1.25rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: 0 4px 24px 0 rgba(34, 41, 47, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  position: relative;
  overflow: hidden;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
  background: currentColor;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 6px 18px -8px rgba(34, 41, 47, 0.56);
}

.stat-card:hover::before {
  opacity: 0.5;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.375rem;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-value {
  font-size: 1.625rem;
  font-weight: 600;
  color: var(--text-color);
  line-height: 1.2;
}

.stat-label {
  font-size: 0.8125rem;
  color: var(--text-color-secondary);
  margin-top: 0.125rem;
}

.stat-trend {
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.25rem 0.625rem;
  border-radius: 4px;
  align-self: flex-start;
}

.stat-trend.positive {
  background: rgba(40, 199, 111, 0.12);
  color: #28c76f;
}

.stat-trend.negative {
  background: rgba(234, 84, 85, 0.12);
  color: #ea5455;
}

/* Content Grid */
.content-grid {
  display: grid;
  grid-template-columns: 1.5fr 1fr;
  gap: 1.5rem;
}

@media (max-width: 1024px) {
  .content-grid {
    grid-template-columns: 1fr;
  }
}

/* Card Styles - Vuexy */
:deep(.p-card) {
  border-radius: 10px;
  box-shadow: 0 4px 24px 0 rgba(34, 41, 47, 0.1);
  border: none;
}

:deep(.p-card-body) {
  padding: 1.25rem;
}

:deep(.p-card-title) {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-color);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.view-all {
  font-size: 0.8125rem;
  font-weight: 500;
  color: #7367f0;
}

.view-all:hover {
  text-decoration: underline;
}

/* Activities - Vuexy Style */
.activities-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 0.875rem 0;
  border-bottom: 1px solid var(--surface-border);
  transition: background 0.2s ease;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-item:hover {
  background: transparent;
}

.activity-indicator {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  margin-top: 0.375rem;
  flex-shrink: 0;
}

.activity-indicator.success { background: #28c76f; }
.activity-indicator.warning { background: #ff9f43; }
.activity-indicator.info { background: #00cfe8; }

.activity-action {
  font-weight: 500;
  font-size: 0.9375rem;
  margin-bottom: 0.25rem;
  color: var(--text-color);
}

.activity-meta {
  display: flex;
  gap: 0.5rem;
  font-size: 0.8125rem;
  color: var(--text-color-secondary);
}

.activity-user::after {
  content: '•';
  margin-left: 0.5rem;
}

/* Quick Actions - Vuexy Style */
.actions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.875rem;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.625rem;
  padding: 1.25rem 0.75rem;
  background: var(--surface-ground);
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: var(--text-color);
}

.action-btn:hover {
  background: linear-gradient(72.47deg, #7367f0 22.16%, rgba(115, 103, 240, 0.7) 76.47%);
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px 0 rgba(115, 103, 240, 0.4);
}

.action-btn i {
  font-size: 1.375rem;
}

.action-btn span {
  font-size: 0.8125rem;
  font-weight: 500;
}
</style>
