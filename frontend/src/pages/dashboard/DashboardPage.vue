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

// Mock data for dashboard stats
const stats = [
  { label: 'Active Orders', value: '127', icon: 'pi pi-shopping-cart', color: '#e91e63', trend: '+12%' },
  { label: 'Low Stock Items', value: '23', icon: 'pi pi-exclamation-triangle', color: '#ff9800', trend: '-5%' },
  { label: 'Work Orders', value: '45', icon: 'pi pi-cog', color: '#9c27b0', trend: '+8%' },
  { label: 'Pending Approvals', value: '8', icon: 'pi pi-clock', color: '#2196f3', trend: '0%' }
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
.dashboard-page {
  padding: 1.5rem;
  max-width: 1400px;
  margin: 0 auto;
}

/* Welcome Section */
.welcome-section {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
}

.welcome-text h1 {
  font-size: 1.75rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  background: linear-gradient(90deg, var(--text-color), var(--primary-color));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.welcome-text p {
  color: var(--text-color-secondary);
}

.date-display {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--text-color-secondary);
  font-size: 0.9375rem;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: var(--surface-card);
  border-radius: 16px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--text-color);
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-color-secondary);
}

.stat-trend {
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.25rem 0.5rem;
  border-radius: 20px;
}

.stat-trend.positive {
  background: rgba(76, 175, 80, 0.15);
  color: #4caf50;
}

.stat-trend.negative {
  background: rgba(244, 67, 54, 0.15);
  color: #f44336;
}

/* Content Grid */
.content-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 1.5rem;
}

@media (max-width: 1024px) {
  .content-grid {
    grid-template-columns: 1fr;
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.view-all {
  font-size: 0.875rem;
  font-weight: 500;
}

/* Activities */
.activities-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 0.75rem;
  border-radius: 8px;
  transition: background 0.2s;
}

.activity-item:hover {
  background: rgba(0, 0, 0, 0.02);
}

.activity-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-top: 0.5rem;
}

.activity-indicator.success { background: #4caf50; }
.activity-indicator.warning { background: #ff9800; }
.activity-indicator.info { background: #2196f3; }

.activity-action {
  font-weight: 500;
  margin-bottom: 0.25rem;
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

/* Quick Actions */
.actions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  padding: 1.25rem;
  background: var(--surface-ground);
  border: 1px solid var(--surface-border);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.action-btn i {
  font-size: 1.25rem;
}

.action-btn span {
  font-size: 0.875rem;
  font-weight: 500;
}
</style>
