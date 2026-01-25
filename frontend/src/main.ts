import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin } from '@tanstack/vue-query'
import PrimeVue from 'primevue/config'
import Aura from '@primeuix/themes/aura'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'
import Tooltip from 'primevue/tooltip'

import App from './App.vue'
import router from './router'
import { useAppStore } from './stores/app.store'

// PrimeIcons
import 'primeicons/primeicons.css'

// Global styles
import './assets/styles/main.css'

// Create Vue app
const app = createApp(App)

// Pinia store
const pinia = createPinia()
app.use(pinia)

// Initialize app store (dark mode)
const appStore = useAppStore()
appStore.initDarkMode()

// PrimeVue with Aura theme
app.use(PrimeVue, {
    theme: {
        preset: Aura,
        options: {
            prefix: 'p',
            darkModeSelector: '.p-dark',
            cssLayer: false
        }
    },
    ripple: true
})

// PrimeVue services
app.use(ToastService)
app.use(ConfirmationService)

// PrimeVue directives
app.directive('tooltip', Tooltip)

// Vue Query
app.use(VueQueryPlugin)

// Router
app.use(router)

// Global error handler
app.config.errorHandler = (err, _instance, info) => {
    console.error('Global error:', err)
    console.error('Error info:', info)
}

// Mount app
app.mount('#app')
