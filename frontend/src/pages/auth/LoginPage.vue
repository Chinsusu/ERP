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
    <!-- Branding Logo (Top Left) -->
    <div class="brand-logo">
      <img src="@/assets/logo.svg" alt="Logo" class="logo-img" />
      <span class="logo-text">VyVy's ERP</span>
    </div>

    <div class="login-wrapper">
      <!-- Left Section: Illustration -->
      <div class="illustration-section">
        <div class="illustration-content">
          <img src="@/assets/images/login-illustration.png" alt="Login Illustration" class="main-img" />
        </div>
      </div>

      <!-- Right Section: Login Form -->
      <div class="form-section">
        <div class="form-container">
          <div class="form-header">
            <h2 class="welcome-title">Welcome to VyVy's ERP! 👋</h2>
            <p class="welcome-subtitle">Please sign-in to your account and start the adventure</p>
          </div>

          <Message v-if="error" severity="error" :closable="false" class="mb-4">
            {{ error }}
          </Message>

          <form @submit.prevent="handleLogin" class="auth-form">
            <div class="form-group">
              <label for="email" class="form-label">Email</label>
              <InputText
                id="email"
                v-model="email"
                type="email"
                placeholder="admin@company.vn"
                class="w-full custom-input"
                :disabled="loading"
              />
            </div>

            <div class="form-group">
              <div class="label-wrapper">
                <label for="password" class="form-label">Password</label>
                <RouterLink to="/forgot-password" class="forgot-link">
                  Forgot Password?
                </RouterLink>
              </div>
              <Password
                id="password"
                v-model="password"
                placeholder="12345678"
                :feedback="false"
                toggleMask
                class="w-full custom-password"
                inputClass="w-full custom-input"
                :disabled="loading"
              />
            </div>

            <div class="form-options">
              <div class="remember-me">
                <Checkbox v-model="rememberMe" :binary="true" inputId="rememberMe" />
                <label for="rememberMe" class="remember-label">Remember Me</label>
              </div>
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
            <p>New on our platform? <RouterLink to="/register" class="create-account">Create an account</RouterLink></p>
          </div>
          
          <div class="demo-info">
            <p>Demo: <strong>admin@company.vn</strong> / <strong>12345678</strong></p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  background-color: #f8f7fa;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Brand Logo (Top Left Overlay) */
.brand-logo {
  position: absolute;
  top: 1.5rem;
  left: 2rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  z-index: 10;
}

.logo-img {
  width: 32px;
  height: 32px;
}

.logo-text {
  font-size: 1.5rem;
  font-weight: 700;
  color: #5e5873;
  letter-spacing: -0.5px;
}

.login-wrapper {
  display: flex;
  width: 100%;
  height: 100vh;
  box-shadow: none;
}

/* Illustration Section (Left 70%) */
.illustration-section {
  flex: 2.2;
  background-color: #f8f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  position: relative;
}

.illustration-content {
  width: 100%;
  max-width: 700px;
  text-align: center;
}

.main-img {
  width: 100%;
  max-width: 600px;
  height: auto;
  filter: drop-shadow(0 20px 50px rgba(0,0,0,0.1));
}

/* Form Section (Right 30%) */
.form-section {
  flex: 1;
  background-color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  box-shadow: -10px 0 30px rgba(0, 0, 0, 0.03);
}

.form-container {
  width: 100%;
  max-width: 400px;
}

.form-header {
  margin-bottom: 2rem;
}

.welcome-title {
  font-size: 1.625rem;
  font-weight: 600;
  color: #5e5873;
  margin-bottom: 0.5rem;
}

.welcome-subtitle {
  color: #6e6b7b;
  font-size: 0.9375rem;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-label {
  display: block;
  font-size: 0.8125rem;
  font-weight: 500;
  color: #5e5873;
  margin-bottom: 0.375rem;
  text-transform: uppercase;
}

.label-wrapper {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.375rem;
}

.forgot-link {
  font-size: 0.8125rem;
  color: #7367f0;
  text-decoration: none;
}

:deep(.custom-input) {
  border-radius: 0.375rem !important;
  border: 1px solid #d8d6de !important;
  padding: 0.572rem 1rem !important;
  font-size: 0.9375rem !important;
  transition: border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out !important;
}

:deep(.custom-input:focus) {
  border-color: #7367f0 !important;
  box-shadow: 0 3px 10px 0 rgba(34, 41, 47, 0.1) !important;
}

.form-options {
  margin-bottom: 1.5rem;
}

.remember-me {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.remember-label {
  font-size: 0.9375rem;
  color: #6e6b7b;
  cursor: pointer;
}

.login-btn {
  background: #7367f0 !important;
  border: none !important;
  padding: 0.75rem !important;
  font-weight: 500 !important;
  font-size: 1rem !important;
  box-shadow: 0 8px 25px -8px rgba(115, 103, 240, 0.48) !important;
  border-radius: 0.375rem !important;
}

.login-btn:hover {
  background: #6558d3 !important;
  box-shadow: 0 8px 25px -8px rgba(115, 103, 240, 0.6) !important;
}

.form-footer {
  margin-top: 1.5rem;
  text-align: center;
  font-size: 0.9375rem;
  color: #6e6b7b;
}

.create-account {
  color: #7367f0;
  text-decoration: none;
  font-weight: 500;
}

.demo-info {
  margin-top: 2rem;
  padding: 1rem;
  background-color: #f8f7fa;
  border-radius: 0.375rem;
  text-align: center;
  font-size: 0.8125rem;
  color: #a5a3ae;
}

.demo-info strong {
  color: #6e6b7b;
}

/* Responsive */
@media (max-width: 1200px) {
  .illustration-section {
    flex: 1.5;
  }
}

@media (max-width: 992px) {
  .illustration-section {
    display: none;
  }
  
  .form-section {
    flex: 1;
    box-shadow: none;
  }

  .brand-logo {
    position: static;
    margin-bottom: 2rem;
    justify-content: center;
  }

  .form-container {
    display: flex;
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
