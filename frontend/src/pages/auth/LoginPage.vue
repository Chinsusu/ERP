<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Checkbox from 'primevue/checkbox'
import Message from 'primevue/message'

const router = useRouter()
const route = useRoute()
const { login, loading } = useAuth()

const email = ref('')
const password = ref('')
const rememberMe = ref(false)
const error = ref('')

const isFormValid = computed(() => {
  return email.value.trim() !== '' && password.value.trim() !== ''
})

async function handleLogin() {
  error.value = ''
  
  const result = await login({
    email: email.value,
    password: password.value,
    remember_me: rememberMe.value
  })
  
  if (!result.success) {
    error.value = result.error || 'Login failed. Please try again.'
  } else {
    // Redirect to original page or dashboard
    const redirect = route.query.redirect as string
    router.push(redirect || '/dashboard')
  }
}
</script>

<template>
  <div class="login-page">
    <div class="login-container">
      <!-- Brand Section -->
      <div class="brand-section">
        <div class="brand-content">
          <i class="pi pi-sparkles brand-icon"></i>
          <h1 class="brand-title">ERP Cosmetics</h1>
          <p class="brand-subtitle">Comprehensive ERP Solution for Natural Cosmetics Manufacturing</p>
        </div>
        <div class="brand-features">
          <div class="feature">
            <i class="pi pi-check-circle"></i>
            <span>Lot Traceability</span>
          </div>
          <div class="feature">
            <i class="pi pi-check-circle"></i>
            <span>FEFO Inventory</span>
          </div>
          <div class="feature">
            <i class="pi pi-check-circle"></i>
            <span>GMP Compliance</span>
          </div>
          <div class="feature">
            <i class="pi pi-check-circle"></i>
            <span>Formula Protection</span>
          </div>
        </div>
      </div>

      <!-- Login Form Section -->
      <div class="form-section">
        <div class="login-form">
          <div class="form-header">
            <h2>Welcome Back</h2>
            <p>Sign in to continue to your dashboard</p>
          </div>

          <Message v-if="error" severity="error" :closable="false" class="mb-3">
            {{ error }}
          </Message>

          <form @submit.prevent="handleLogin">
            <div class="form-group">
              <label for="email">Email</label>
              <InputText
                id="email"
                v-model="email"
                type="email"
                placeholder="Enter your email"
                class="w-full"
                :disabled="loading"
              />
            </div>

            <div class="form-group">
              <label for="password">Password</label>
              <Password
                id="password"
                v-model="password"
                placeholder="Enter your password"
                :feedback="false"
                toggleMask
                class="w-full"
                inputClass="w-full"
                :disabled="loading"
              />
            </div>

            <div class="form-options">
              <div class="remember-me">
                <Checkbox v-model="rememberMe" :binary="true" inputId="rememberMe" />
                <label for="rememberMe">Remember me</label>
              </div>
              <RouterLink to="/forgot-password" class="forgot-link">
                Forgot password?
              </RouterLink>
            </div>

            <Button
              type="submit"
              label="Sign In"
              class="w-full login-btn"
              :loading="loading"
              :disabled="!isFormValid || loading"
            />
          </form>

          <div class="form-footer">
            <p>Demo credentials: <strong>admin@company.vn</strong> / <strong>admin123</strong></p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  padding: 1rem;
}

.login-container {
  display: flex;
  width: 100%;
  max-width: 1000px;
  background: var(--surface-card);
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
}

/* Brand Section */
.brand-section {
  flex: 1;
  padding: 3rem;
  background: linear-gradient(135deg, #e91e63 0%, #9c27b0 100%);
  color: white;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.brand-content {
  margin-bottom: 2rem;
}

.brand-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.brand-title {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.brand-subtitle {
  font-size: 1rem;
  opacity: 0.9;
  line-height: 1.6;
}

.brand-features {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.feature {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 0.95rem;
}

.feature i {
  font-size: 1rem;
}

/* Form Section */
.form-section {
  flex: 1;
  padding: 3rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-form {
  width: 100%;
  max-width: 360px;
}

.form-header {
  text-align: center;
  margin-bottom: 2rem;
}

.form-header h2 {
  font-size: 1.75rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: var(--text-color);
}

.form-header p {
  color: var(--text-color-secondary);
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: var(--text-color);
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.remember-me {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.remember-me label {
  font-size: 0.875rem;
  color: var(--text-color-secondary);
  cursor: pointer;
}

.forgot-link {
  font-size: 0.875rem;
  color: var(--primary-color);
}

.login-btn {
  background: linear-gradient(90deg, #e91e63, #9c27b0);
  border: none;
  padding: 0.875rem;
  font-weight: 600;
}

.login-btn:hover {
  background: linear-gradient(90deg, #c2185b, #7b1fa2);
}

.form-footer {
  margin-top: 1.5rem;
  text-align: center;
  font-size: 0.8125rem;
  color: var(--text-color-secondary);
}

.form-footer strong {
  color: var(--text-color);
}

.w-full {
  width: 100%;
}

.mb-3 {
  margin-bottom: 1rem;
}

/* Responsive */
@media (max-width: 768px) {
  .login-container {
    flex-direction: column;
  }

  .brand-section {
    padding: 2rem;
  }

  .brand-features {
    flex-direction: row;
    flex-wrap: wrap;
    gap: 0.5rem 1rem;
  }

  .form-section {
    padding: 2rem;
  }
}
</style>
