<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits(['submit'])

const username = ref('')
const password = ref('')
const message = ref('')
const isError = ref(false)

const login = async () => {
  try {
    emit('submit', {
      username: username.value,
      password: password.value
    })

    message.value = ''
  } catch (error: any) {
    message.value = error.message
    isError.value = true
  }
}
</script>

<template>
  <div class="login-page">
    <h1>Login</h1>

    <form @submit.prevent="login" class="login-form">
      <div>
        <label for="username">Username:</label>

        <input
          id="username"
          v-model="username"
          type="text"
          required
        />
      </div>

      <div>
        <label for="password">Password:</label>

        <input
          id="password"
          v-model="password"
          type="password"
          required
        />
      </div>

      <button type="submit">
        Login
      </button>

      <p
        v-if="message"
        :style="{
          color: isError ? 'red' : 'green'
        }"
      >
        {{ message }}
      </p>
    </form>
  </div>
</template>


<style scoped>
.login-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
h1 {
  text-align: center;
}

input {
  padding: 10px;
}

button {
    background-color: rgb(38, 243, 38);
  padding: 10px;
  cursor: pointer;
}
.login-page {
  display: flexbox;
  align-items: center;
  justify-content: center;
  max-width: 400px;
  margin: 0 auto;
  padding: 20px;
}
</style>
