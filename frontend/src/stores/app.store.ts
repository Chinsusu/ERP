import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useAppStore = defineStore('app', () => {
    // State
    const sidebarCollapsed = ref(localStorage.getItem('sidebar_collapsed') === 'true')
    const darkMode = ref(localStorage.getItem('dark_mode') === 'true')
    const pageTitle = ref('Dashboard')
    const breadcrumbs = ref<{ label: string; to?: string }[]>([])

    // Watch and persist to localStorage
    watch(sidebarCollapsed, (value) => {
        localStorage.setItem('sidebar_collapsed', String(value))
    })

    watch(darkMode, (value) => {
        localStorage.setItem('dark_mode', String(value))
        updateDarkModeClass(value)
    })

    // Apply dark mode class to document
    function updateDarkModeClass(isDark: boolean) {
        if (isDark) {
            document.documentElement.classList.add('p-dark')
        } else {
            document.documentElement.classList.remove('p-dark')
        }
    }

    // Actions
    function toggleSidebar() {
        sidebarCollapsed.value = !sidebarCollapsed.value
    }

    function toggleDarkMode() {
        darkMode.value = !darkMode.value
    }

    function setDarkMode(value: boolean) {
        darkMode.value = value
    }

    function setPageTitle(title: string) {
        pageTitle.value = title
        document.title = `${title} | ERP Cosmetics`
    }

    function setBreadcrumbs(items: { label: string; to?: string }[]) {
        breadcrumbs.value = items
    }

    // Initialize dark mode on startup
    function initDarkMode() {
        updateDarkModeClass(darkMode.value)
    }

    return {
        // State
        sidebarCollapsed,
        darkMode,
        pageTitle,
        breadcrumbs,
        // Actions
        toggleSidebar,
        toggleDarkMode,
        setDarkMode,
        setPageTitle,
        setBreadcrumbs,
        initDarkMode
    }
})
