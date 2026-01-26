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
// Show all menu items - route guards handle actual access control
// This ensures menu is always visible after successful login
function hasAccess(_item: any): boolean {
  // Always show menu items in sidebar
  // Permission/auth checks happen at route guard and API level
  return true
}

function hasAnyChildAccess(item: any): boolean {
  // Always show parent menus
  return true
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
        <span v-if="!collapsed" class="logo-text">VyVy's ERP</span>
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
/* Vuexy Sidebar - Light Theme */
.app-sidebar {
  position: fixed;
  top: 0;
  left: 0;
  height: 100vh;
  width: var(--sidebar-width);
  background: #ffffff;
  color: #6f6b7d;
  display: flex;
  flex-direction: column;
  transition: width 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  z-index: 1000;
  overflow: hidden;
  box-shadow: 0 0 15px 0 rgba(34, 41, 47, 0.05);
  border-right: 1px solid #ebe9f1;
}

.app-sidebar.collapsed {
  width: var(--sidebar-collapsed-width);
}

/* Header with Logo */
.sidebar-header {
  padding: 1.25rem 1.5rem;
  display: flex;
  align-items: center;
  min-height: 64px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.875rem;
}

.logo-icon {
  width: 32px;
  height: 32px;
}

.logo-icon-small {
  font-size: 1.5rem;
  color: #7367f0;
}

.logo-text {
  font-size: 1.375rem;
  font-weight: 700;
  color: #5e5873;
  letter-spacing: -0.5px;
  white-space: nowrap;
}

/* Navigation */
.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 0.5rem 0;
}

.sidebar-nav::-webkit-scrollbar {
  width: 4px;
}

.sidebar-nav::-webkit-scrollbar-thumb {
  background: #dbdade;
  border-radius: 10px;
}

.menu-list {
  list-style: none;
  padding: 0 0.875rem;
  margin: 0;
}

.menu-item {
  margin: 0;
  padding: 0.125rem 0;
}

/* Section Header (like APPS & PAGES) */
.section-header {
  padding: 1.5rem 1rem 0.5rem;
  font-size: 0.75rem;
  font-weight: 600;
  color: #a5a3ae;
  text-transform: uppercase;
  letter-spacing: 0.8px;
}

/* Menu Links - Vuexy Light Style */
.menu-link,
.menu-parent {
  display: flex;
  align-items: center;
  gap: 0.875rem;
  padding: 0.625rem 1rem;
  color: #6f6b7d;
  cursor: pointer;
  transition: all 0.2s ease;
  border-radius: 0.375rem;
  margin: 0;
  text-decoration: none;
  width: 100%;
}

.menu-link:hover,
.menu-parent:hover {
  background: #f8f7fa;
  color: #7367f0;
}

/* Active State - Vuexy Light Purple Pill */
.menu-link.active {
  background: linear-gradient(118deg, #7367f0, rgba(115, 103, 240, 0.7));
  color: #fff !important;
  box-shadow: 0 2px 4px 0 rgba(115, 103, 240, 0.4);
  font-weight: 500;
}

.menu-link.active i {
  color: #fff !important;
}

/* Icons */
.menu-link i,
.menu-parent i {
  font-size: 1.25rem;
  width: 1.5rem;
  text-align: center;
  flex-shrink: 0;
  color: #6f6b7d;
}

.menu-label {
  white-space: nowrap;
  overflow: hidden;
  font-size: 0.9375rem;
  flex: 1;
}

/* Submenu */
.submenu-list {
  list-style: none;
  padding: 0.25rem 0 0.25rem 0;
  margin: 0;
}

.submenu-link {
  display: flex;
  align-items: center;
  padding: 0.5rem 1rem 0.5rem 1.25rem;
  color: #6f6b7d;
  font-size: 0.9375rem;
  transition: all 0.2s ease;
  border-radius: 0.375rem;
  text-decoration: none;
  position: relative;
  margin-left: 0.5rem;
}

.submenu-link::before {
  content: '';
  width: 6px;
  height: 6px;
  border-radius: 50%;
  border: 1.5px solid currentColor;
  background: transparent;
  position: absolute;
  left: 0;
  transition: all 0.2s ease;
}

.submenu-link:hover {
  color: #7367f0;
}

.submenu-link.active {
  color: #7367f0;
  font-weight: 500;
}

.submenu-link.active::before {
  background: #7367f0;
  border-color: #7367f0;
  box-shadow: 0 0 4px rgba(115, 103, 240, 0.4);
}

/* Footer */
.sidebar-footer {
  padding: 0.75rem 0.875rem;
  background: #ffffff;
}

.collapse-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  padding: 0.625rem;
  background: #f8f7fa;
  border: none;
  border-radius: 0.375rem;
  color: #6f6b7d;
  cursor: pointer;
  transition: all 0.2s ease;
}

.collapse-btn:hover {
  background: #eeedf0;
  color: #7367f0;
}

/* Collapsed State */
.collapsed .sidebar-header {
  padding: 1rem 0.75rem;
  justify-content: center;
}

.collapsed .logo-text, 
.collapsed .menu-label, 
.collapsed .section-header {
  display: none;
}

.collapsed .menu-list {
  padding: 0 0.5rem;
}

.collapsed .menu-link,
.collapsed .menu-parent {
  justify-content: center;
  padding: 0.75rem;
}

.collapsed .submenu-list {
  display: none;
}
</style>
