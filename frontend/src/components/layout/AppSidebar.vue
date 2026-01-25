<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app.store'
import { useAuthStore } from '@/stores/auth.store'

const route = useRoute()
const appStore = useAppStore()
const authStore = useAuthStore()

const collapsed = computed(() => appStore.sidebarCollapsed)

// Navigation menu items with permissions
const menuItems = [
  {
    label: 'Dashboard',
    icon: 'pi pi-home',
    to: '/dashboard',
    permission: null
  },
  {
    label: 'Master Data',
    icon: 'pi pi-database',
    children: [
      { label: 'Materials', to: '/materials', permission: 'master_data:material:read' },
      { label: 'Products', to: '/products', permission: 'master_data:product:read' },
      { label: 'Categories', to: '/categories', permission: 'master_data:category:read' }
    ]
  },
  {
    label: 'Supply Chain',
    icon: 'pi pi-truck',
    children: [
      { label: 'Suppliers', to: '/suppliers', permission: 'supplier:supplier:read' },
      { label: 'Purchase Requisitions', to: '/procurement/requisitions', permission: 'procurement:pr:read' },
      { label: 'Purchase Orders', to: '/procurement/orders', permission: 'procurement:po:read' }
    ]
  },
  {
    label: 'Warehouse',
    icon: 'pi pi-box',
    children: [
      { label: 'Stock', to: '/warehouse/stock', permission: 'wms:stock:read' },
      { label: 'Goods Receipt', to: '/warehouse/grn', permission: 'wms:grn:read' },
      { label: 'Goods Issue', to: '/warehouse/goods-issue', permission: 'wms:gi:read' }
    ]
  },
  {
    label: 'Manufacturing',
    icon: 'pi pi-cog',
    children: [
      { label: 'Bill of Materials', to: '/manufacturing/bom', permission: 'manufacturing:bom:read' },
      { label: 'Work Orders', to: '/manufacturing/work-orders', permission: 'manufacturing:wo:read' },
      { label: 'Quality Control', to: '/manufacturing/qc', permission: 'manufacturing:qc:read' }
    ]
  },
  {
    label: 'Sales',
    icon: 'pi pi-shopping-cart',
    children: [
      { label: 'Customers', to: '/sales/customers', permission: 'sales:customer:read' },
      { label: 'Sales Orders', to: '/sales/orders', permission: 'sales:order:read' },
      { label: 'Quotations', to: '/sales/quotations', permission: 'sales:quotation:read' }
    ]
  },
  {
    label: 'Marketing',
    icon: 'pi pi-megaphone',
    children: [
      { label: 'Campaigns', to: '/marketing/campaigns', permission: 'marketing:campaign:read' },
      { label: 'KOL Management', to: '/marketing/kols', permission: 'marketing:kol:read' }
    ]
  },
  {
    label: 'Reports',
    icon: 'pi pi-chart-bar',
    to: '/reports',
    permission: 'report:report:read'
  },
  {
    label: 'Settings',
    icon: 'pi pi-sliders-h',
    children: [
      { label: 'Users', to: '/users', permission: 'user:user:read' },
      { label: 'Roles', to: '/settings/roles', permission: 'auth:role:read' }
    ]
  }
]

// Filter menu items based on permissions
function hasAccess(item: any): boolean {
  if (!item.permission) return true
  return authStore.hasPermission(item.permission)
}

function hasAnyChildAccess(item: any): boolean {
  if (!item.children) return hasAccess(item)
  return item.children.some((child: any) => hasAccess(child))
}

function isActive(to: string): boolean {
  return route.path === to || route.path.startsWith(to + '/')
}

function toggleSidebar() {
  appStore.toggleSidebar()
}
</script>

<template>
  <aside class="app-sidebar" :class="{ collapsed }">
    <!-- Logo -->
    <div class="sidebar-header">
      <div class="logo">
        <img src="@/assets/logo.svg" alt="Logo" class="logo-icon" v-if="!collapsed" />
        <span v-if="!collapsed" class="logo-text">ERP Cosmetics</span>
        <i v-else class="pi pi-sparkles logo-icon-small"></i>
      </div>
    </div>

    <!-- Navigation Menu -->
    <nav class="sidebar-nav">
      <ul class="menu-list">
        <template v-for="item in menuItems" :key="item.label">
          <li v-if="hasAnyChildAccess(item)" class="menu-item">
            <!-- Simple link (no children) -->
            <RouterLink
              v-if="!item.children"
              :to="item.to"
              class="menu-link"
              :class="{ active: isActive(item.to) }"
              v-tooltip.right="collapsed ? item.label : null"
            >
              <i :class="item.icon"></i>
              <span v-if="!collapsed" class="menu-label">{{ item.label }}</span>
            </RouterLink>

            <!-- Parent with children -->
            <template v-else>
              <div class="menu-parent" v-tooltip.right="collapsed ? item.label : null">
                <i :class="item.icon"></i>
                <span v-if="!collapsed" class="menu-label">{{ item.label }}</span>
              </div>
              <ul v-if="!collapsed" class="submenu-list">
                <li v-for="child in item.children" :key="child.label">
                  <RouterLink
                    v-if="hasAccess(child)"
                    :to="child.to"
                    class="submenu-link"
                    :class="{ active: isActive(child.to) }"
                  >
                    {{ child.label }}
                  </RouterLink>
                </li>
              </ul>
            </template>
          </li>
        </template>
      </ul>
    </nav>

    <!-- Collapse Toggle -->
    <div class="sidebar-footer">
      <button class="collapse-btn" @click="toggleSidebar" v-tooltip.right="collapsed ? 'Expand' : 'Collapse'">
        <i :class="collapsed ? 'pi pi-angle-right' : 'pi pi-angle-left'"></i>
      </button>
    </div>
  </aside>
</template>

<style scoped>
.app-sidebar {
  position: fixed;
  top: 0;
  left: 0;
  height: 100vh;
  width: var(--sidebar-width);
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
  color: #fff;
  display: flex;
  flex-direction: column;
  transition: width var(--transition-duration) ease;
  z-index: 1000;
  overflow: hidden;
}

.app-sidebar.collapsed {
  width: var(--sidebar-collapsed-width);
}

.sidebar-header {
  padding: 1.25rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.logo-icon {
  width: 36px;
  height: 36px;
}

.logo-icon-small {
  font-size: 1.5rem;
  color: var(--primary-color);
}

.logo-text {
  font-size: 1.125rem;
  font-weight: 700;
  background: linear-gradient(90deg, #e91e63, #9c27b0);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  white-space: nowrap;
}

.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  padding: 1rem 0;
}

.menu-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.menu-item {
  margin-bottom: 0.25rem;
}

.menu-link,
.menu-parent {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1.25rem;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all var(--transition-duration);
  border-left: 3px solid transparent;
}

.menu-link:hover,
.menu-parent:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
}

.menu-link.active {
  background: rgba(233, 30, 99, 0.15);
  color: var(--primary-color);
  border-left-color: var(--primary-color);
}

.menu-link i,
.menu-parent i {
  font-size: 1.125rem;
  width: 1.5rem;
  text-align: center;
}

.menu-label {
  white-space: nowrap;
  overflow: hidden;
}

.submenu-list {
  list-style: none;
  padding: 0 0 0 2.5rem;
  margin: 0.25rem 0;
}

.submenu-link {
  display: block;
  padding: 0.5rem 1rem;
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.875rem;
  transition: all var(--transition-duration);
  border-radius: 4px;
}

.submenu-link:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.05);
}

.submenu-link.active {
  color: var(--primary-color);
  background: rgba(233, 30, 99, 0.1);
}

.sidebar-footer {
  padding: 1rem;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.collapse-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  padding: 0.5rem;
  background: rgba(255, 255, 255, 0.05);
  border: none;
  border-radius: 8px;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all var(--transition-duration);
}

.collapse-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

/* Collapsed state styles */
.collapsed .sidebar-header {
  padding: 1rem;
  display: flex;
  justify-content: center;
}

.collapsed .menu-link,
.collapsed .menu-parent {
  justify-content: center;
  padding: 0.75rem;
}

.collapsed .menu-link i,
.collapsed .menu-parent i {
  margin: 0;
}
</style>
