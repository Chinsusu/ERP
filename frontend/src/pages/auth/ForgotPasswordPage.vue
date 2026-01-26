<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '@/api/auth.api'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import Message from 'primevue/message'

const router = useRouter()

const email = ref('')
const loading = ref(false)
const error = ref('')
const success = ref(false)

const isFormValid = computed(() => {
  return email.value.trim() !== '' && email.value.includes('@')
})

async function handleSubmit() {
  error.value = ''
  loading.value = true
  
  try {
    await authApi.forgotPassword({ email: email.value })
    success.value = true
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to send reset email. Please try again.'
  } finally {
    loading.value = false
  }
}

function goToLogin() {
  router.push('/login')
}
</script>

<template>
  <div class="forgot-password-page">
    <div class="forgot-container">
      <div class="forgot-form">
        <!-- Success State -->
        <template v-if="success">
          <div class="success-content">
            <div class="success-icon">
              <i class="pi pi-check-circle"></i>
            </div>
            <h2>Check Your Email</h2>
            <p>We've sent a password reset link to <strong>{{ email }}</strong></p>
            <p class="text-muted">Didn't receive the email? Check your spam folder.</p>
            <Button
              label="Back to Login"
              class="w-full mt-4"
              @click="goToLogin"
            />
          </div>
        </template>

        <!-- Form State -->
        <template v-else>
          <div class="form-header">
            <div class="header-icon">
              <i class="pi pi-lock"></i>
            </div>
            <h2>Forgot Password?</h2>
            <p>Enter your email address and we'll send you a link to reset your password.</p>
          </div>

          <Message v-if="error" severity="error" :closable="false" class="mb-3">
            {{ error }}
          </Message>

          <form @submit.prevent="handleSubmit">
            <div class="form-group">
              <label for="email">Email Address</label>
              <InputText
                id="email"
                v-model="email"
                type="email"
                placeholder="Enter your email"
                class="w-full"
                :disabled="loading"
              />
            </div>

            <Button
              type="submit"
              label="Send Reset Link"
              class="w-full submit-btn"
              :loading="loading"
              :disabled="!isFormValid || loading"
            />
          </form>

          <div class="form-footer">
            <RouterLink to="/login" class="back-link">
              <i class="pi pi-arrow-left"></i>
              Back to Login
            </RouterLink>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
.forgot-password-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #25293c 0%, #2f3349 50%, #25293c 100%);
  padding: 1rem;
}

.forgot-container {
  width: 100%;
  max-width: 420px;
}

.forgot-form {
  background: var(--surface-card);
  padding: 2.5rem;
  border-radius: 12px;
  box-shadow: 0 4px 24px 0 rgba(47, 43, 61, 0.18);
}

.form-header {
  text-align: center;
  margin-bottom: 2rem;
}

.header-icon {
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #7367f0, #9e95f5);
  border-radius: 50%;
  margin: 0 auto 1rem;
}

.header-icon i {
  font-size: 1.75rem;
  color: white;
}

.form-header h2 {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: var(--text-color);
}

.form-header p {
  color: var(--text-color-secondary);
  font-size: 0.9375rem;
  line-height: 1.5;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: var(--text-color);
}

.submit-btn {
  background: var(--primary-color) !important;
  border: none !important;
  padding: 0.875rem;
  font-weight: 600;
  box-shadow: 0 2px 6px 0 rgba(115, 103, 240, 0.5);
}

.submit-btn:hover {
  background: var(--primary-hover) !important;
  transform: translateY(-1px);
}

.form-footer {
  margin-top: 1.5rem;
  text-align: center;
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--text-color-secondary);
  font-size: 0.875rem;
  transition: color 0.2s;
}

.back-link:hover {
  color: var(--primary-color);
}

/* Success State */
.success-content {
  text-align: center;
}

.success-icon {
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(40, 199, 111, 0.15);
  border-radius: 50%;
  margin: 0 auto 1.5rem;
}

.success-icon i {
  font-size: 2.5rem;
  color: #28c76f;
}

.success-content h2 {
  font-size: 1.5rem;
  margin-bottom: 0.75rem;
  color: var(--text-color);
}

.success-content p {
  color: var(--text-color-secondary);
  margin-bottom: 0.5rem;
}

.text-muted {
  font-size: 0.875rem;
  opacity: 0.7;
}

.w-full {
  width: 100%;
}

.mt-4 {
  margin-top: 1.5rem;
}

.mb-3 {
  margin-bottom: 1rem;
}
</style>
