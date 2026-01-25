<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app.store'
import { useAuthStore } from '@/stores/auth.store'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import Menu from 'primevue/menu'
import Badge from 'primevue/badge'
import { ref } from 'vue'

const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()

const userMenu = ref()
const notificationMenu = ref()

const user = computed(() => authStore.user)
const darkMode = computed(() => appStore.darkMode)
const sidebarCollapsed = computed(() => appStore.sidebarCollapsed)

const userMenuItems = [
  {
    label: 'Profile',
    icon: 'pi pi-user',
    command: () => router.push('/profile')
  },
  {
    label: 'Settings',
    icon: 'pi pi-cog',
    command: () => router.push('/settings')
  },
  { separator: true },
  {
    label: 'Logout',
    icon: 'pi pi-sign-out',
    command: () => logout()
  }
]

const notifications = ref([
  { id: 1, title: 'Low stock alert', message: 'Material RM-0025 is below minimum', time: '5 min ago', type: 'warning' },
  { id: 2, title: 'PO approved', message: 'PO-2026-0045 has been approved', time: '1 hour ago', type: 'success' }
])

function toggleUserMenu(event: Event) {
  userMenu.value.toggle(event)
}

function toggleNotifications(event: Event) {
  notificationMenu.value.toggle(event)
}

function toggleDarkMode() {
  appStore.toggleDarkMode()
}

function toggleSidebar() {
  appStore.toggleSidebar()
}

async function logout() {
  authStore.logout()
  router.push('/login')
}

const userInitials = computed(() => {
  if (!user.value) return 'U'
  const first = user.value.first_name?.[0] || ''
  const last = user.value.last_name?.[0] || ''
  return (first + last).toUpperCase() || 'U'
})
</script>

<template>
  <header class="app-header" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
    <div class="header-left">
      <Button
        icon="pi pi-bars"
        text
        rounded
        class="sidebar-toggle"
        @click="toggleSidebar"
        v-tooltip.bottom="'Toggle Sidebar'"
      />
      
      <div class="breadcrumb">
        <span class="breadcrumb-item">{{ appStore.pageTitle }}</span>
      </div>
    </div>

    <div class="header-right">
      <!-- Dark Mode Toggle -->
      <Button
        :icon="darkMode ? 'pi pi-sun' : 'pi pi-moon'"
        text
        rounded
        class="header-btn"
        @click="toggleDarkMode"
        v-tooltip.bottom="darkMode ? 'Light Mode' : 'Dark Mode'"
      />

      <!-- Notifications -->
      <div class="notification-wrapper">
        <Button
          icon="pi pi-bell"
          text
          rounded
          class="header-btn"
          @click="toggleNotifications"
          v-tooltip.bottom="'Notifications'"
        >
          <Badge 
            v-if="notifications.length > 0" 
            :value="notifications.length" 
            severity="danger" 
            class="notification-badge"
          />
        </Button>
        <Menu ref="notificationMenu" :popup="true" class="notification-menu">
          <template #start>
            <div class="notification-header">
              <span>Notifications</span>
              <Button label="Mark all read" text size="small" />
            </div>
          </template>
          <template #item>
            <div v-for="notification in notifications" :key="notification.id" class="notification-item">
              <i :class="['pi', notification.type === 'warning' ? 'pi-exclamation-triangle text-warning' : 'pi-check-circle text-success']"></i>
              <div class="notification-content">
                <div class="notification-title">{{ notification.title }}</div>
                <div class="notification-message">{{ notification.message }}</div>
                <div class="notification-time">{{ notification.time }}</div>
              </div>
            </div>
          </template>
        </Menu>
      </div>

      <!-- User Menu -->
      <div class="user-wrapper">
        <Button
          text
          rounded
          class="user-btn"
          @click="toggleUserMenu"
        >
          <Avatar 
            :label="userInitials" 
            shape="circle" 
            class="user-avatar"
            style="background: linear-gradient(135deg, #e91e63, #9c27b0); color: white;"
          />
          <span class="user-name">{{ user?.full_name || 'User' }}</span>
          <i class="pi pi-chevron-down"></i>
        </Button>
        <Menu ref="userMenu" :model="userMenuItems" :popup="true" />
      </div>
    </div>
  </header>
</template>

<style scoped>
.app-header {
  position: fixed;
  top: 0;
  right: 0;
  left: var(--sidebar-width);
  height: var(--header-height);
  background: var(--surface-card);
  border-bottom: 1px solid var(--surface-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1.5rem;
  z-index: 999;
  transition: left var(--transition-duration) ease;
}

.app-header.sidebar-collapsed {
  left: var(--sidebar-collapsed-width);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.sidebar-toggle {
  display: none;
}

@media (max-width: 992px) {
  .sidebar-toggle {
    display: flex;
  }
}

.breadcrumb {
  font-size: 1rem;
  font-weight: 500;
  color: var(--text-color);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.header-btn {
  color: var(--text-color-secondary);
}

.header-btn:hover {
  color: var(--text-color);
}

.notification-wrapper {
  position: relative;
}

.notification-badge {
  position: absolute;
  top: 4px;
  right: 4px;
  pointer-events: none;
}

.user-wrapper {
  margin-left: 0.5rem;
}

.user-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem;
  color: var(--text-color);
}

.user-avatar {
  width: 32px;
  height: 32px;
  font-size: 0.875rem;
}

.user-name {
  font-weight: 500;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

@media (max-width: 576px) {
  .user-name {
    display: none;
  }
}

/* Notification Menu */
.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--surface-border);
  font-weight: 600;
}

.notification-item {
  display: flex;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--surface-border);
}

.notification-item:last-child {
  border-bottom: none;
}

.notification-content {
  flex: 1;
}

.notification-title {
  font-weight: 500;
  font-size: 0.875rem;
}

.notification-message {
  font-size: 0.8125rem;
  color: var(--text-color-secondary);
}

.notification-time {
  font-size: 0.75rem;
  color: var(--text-color-secondary);
  margin-top: 0.25rem;
}

.text-warning {
  color: #ff9800;
}

.text-success {
  color: #4caf50;
}
</style>
