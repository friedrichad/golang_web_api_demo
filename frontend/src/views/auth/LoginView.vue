<template>
  <div class="login-page">
    <LoginForm @submit="handleLogin" />
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import LoginForm from '@/views/auth/components/LoginItem.vue'
import { loginApi } from '@/views/auth/auth'

const router = useRouter()

const handleLogin = async (form: {
  username: string
  password: string
}) => {
  try {
    await loginApi(form.username, form.password)

    alert('Login successful!')

    router.push('/')
  } catch (error: any) {
    alert(error.message)
  }
}
</script>