import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { authGuard } from './guards'

const routes: RouteRecordRaw[] = [
    // Auth routes (no layout)
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/pages/auth/LoginPage.vue'),
        meta: { requiresAuth: false, title: 'Login' }
    },
    {
        path: '/forgot-password',
        name: 'ForgotPassword',
        component: () => import('@/pages/auth/ForgotPasswordPage.vue'),
        meta: { requiresAuth: false, title: 'Forgot Password' }
    },

    // Main app routes (with DefaultLayout)
    {
        path: '/',
        component: () => import('@/layouts/DefaultLayout.vue'),
        meta: { requiresAuth: true },
        children: [
            {
                path: '',
                redirect: '/dashboard'
            },
            {
                path: 'dashboard',
                name: 'Dashboard',
                component: () => import('@/pages/dashboard/DashboardPage.vue'),
                meta: { title: 'Dashboard' }
            },

            // User Management
            {
                path: 'users',
                name: 'Users',
                component: () => import('@/pages/users/UsersPage.vue'),
                meta: { title: 'Users', permission: 'user:user:read' }
            },

            // Master Data - Materials
            {
                path: 'materials',
                name: 'Materials',
                component: () => import('@/pages/materials/MaterialsPage.vue'),
                meta: { title: 'Materials', permission: 'master_data:material:read' }
            },
            {
                path: 'materials/new',
                name: 'MaterialCreate',
                component: () => import('@/pages/materials/MaterialFormPage.vue'),
                meta: { title: 'New Material', permission: 'master_data:material:create' }
            },
            {
                path: 'materials/:id',
                name: 'MaterialDetail',
                component: () => import('@/pages/materials/MaterialDetailPage.vue'),
                meta: { title: 'Material Details', permission: 'master_data:material:read' }
            },
            {
                path: 'materials/:id/edit',
                name: 'MaterialEdit',
                component: () => import('@/pages/materials/MaterialFormPage.vue'),
                meta: { title: 'Edit Material', permission: 'master_data:material:update' }
            },

            // Suppliers
            {
                path: 'suppliers',
                name: 'Suppliers',
                component: () => import('@/pages/suppliers/SuppliersPage.vue'),
                meta: { title: 'Suppliers', permission: 'supplier:supplier:read' }
            },
            {
                path: 'suppliers/new',
                name: 'SupplierCreate',
                component: () => import('@/pages/suppliers/SupplierFormPage.vue'),
                meta: { title: 'New Supplier', permission: 'supplier:supplier:create' }
            },
            {
                path: 'suppliers/:id',
                name: 'SupplierDetail',
                component: () => import('@/pages/suppliers/SupplierDetailPage.vue'),
                meta: { title: 'Supplier Details', permission: 'supplier:supplier:read' }
            },
            {
                path: 'suppliers/:id/edit',
                name: 'SupplierEdit',
                component: () => import('@/pages/suppliers/SupplierFormPage.vue'),
                meta: { title: 'Edit Supplier', permission: 'supplier:supplier:update' }
            },

            // Procurement
            {
                path: 'procurement',
                children: [
                    {
                        path: 'requisitions',
                        name: 'PurchaseRequisitions',
                        component: () => import('@/pages/procurement/RequisitionsPage.vue'),
                        meta: { title: 'Purchase Requisitions', permission: 'procurement:pr:read' }
                    },
                    {
                        path: 'orders',
                        name: 'PurchaseOrders',
                        component: () => import('@/pages/procurement/OrdersPage.vue'),
                        meta: { title: 'Purchase Orders', permission: 'procurement:po:read' }
                    }
                ]
            },

            // Warehouse
            {
                path: 'warehouse',
                children: [
                    {
                        path: 'stock',
                        name: 'Stock',
                        component: () => import('@/pages/warehouse/StockPage.vue'),
                        meta: { title: 'Stock', permission: 'wms:stock:read' }
                    },
                    {
                        path: 'grn',
                        name: 'GRN',
                        component: () => import('@/pages/warehouse/GRNPage.vue'),
                        meta: { title: 'Goods Receipt', permission: 'wms:grn:read' }
                    }
                ]
            },

            // Manufacturing
            {
                path: 'manufacturing',
                children: [
                    {
                        path: 'bom',
                        name: 'BOM',
                        component: () => import('@/pages/manufacturing/BOMPage.vue'),
                        meta: { title: 'Bill of Materials', permission: 'manufacturing:bom:read' }
                    },
                    {
                        path: 'work-orders',
                        name: 'WorkOrders',
                        component: () => import('@/pages/manufacturing/WorkOrdersPage.vue'),
                        meta: { title: 'Work Orders', permission: 'manufacturing:wo:read' }
                    }
                ]
            },

            // Sales
            {
                path: 'sales',
                children: [
                    {
                        path: 'customers',
                        name: 'Customers',
                        component: () => import('@/pages/sales/CustomersPage.vue'),
                        meta: { title: 'Customers', permission: 'sales:customer:read' }
                    },
                    {
                        path: 'orders',
                        name: 'SalesOrders',
                        component: () => import('@/pages/sales/OrdersPage.vue'),
                        meta: { title: 'Sales Orders', permission: 'sales:order:read' }
                    }
                ]
            },

            // Marketing
            {
                path: 'marketing',
                children: [
                    {
                        path: 'campaigns',
                        name: 'Campaigns',
                        component: () => import('@/pages/marketing/CampaignsPage.vue'),
                        meta: { title: 'Campaigns', permission: 'marketing:campaign:read' }
                    },
                    {
                        path: 'kols',
                        name: 'KOLs',
                        component: () => import('@/pages/marketing/KOLsPage.vue'),
                        meta: { title: 'KOL Management', permission: 'marketing:kol:read' }
                    }
                ]
            },

            // Settings
            {
                path: 'settings',
                children: [
                    {
                        path: 'roles',
                        name: 'Roles',
                        component: () => import('@/pages/settings/RolesPage.vue'),
                        meta: { title: 'Roles', permission: 'auth:role:read' }
                    }
                ]
            }
        ]
    },

    // 404 Not Found
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('@/pages/NotFoundPage.vue'),
        meta: { title: 'Not Found' }
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// Apply navigation guards
router.beforeEach(authGuard)

// Set page title
router.afterEach((to) => {
    const title = to.meta.title as string || "VyVy's ERP"
    document.title = `${title} | VyVy's ERP`
})

export default router
