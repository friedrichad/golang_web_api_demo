<script setup lang="ts">
import { ref } from 'vue'
import { loginApi } from '../services/auth'
import { useRouter } from 'vue-router'

const router = useRouter()
const username = ref('')
const password = ref('')
const message = ref('')
const isError = ref(false)


const login = async () => {
    try {
        await loginApi(username.value, password.value)
        message.value = 'Login successful!'
        isError.value = false
        alert(message.value)
        router.push('/')
    } catch (error: any) {
        message.value = error.message
        isError.value = true
        alert(message.value)
        console.error(error)
    }
}
</script>


<template>
    <div>
        <h1>Login</h1>
        <form @submit.prevent="login">
            <div>
                <label for="username">Username:</label>
                <input id="username" v-model="username" type="text" required />
            </div>
            <div>
                <label for="password">Password:</label>
                <input id="password" v-model="password" type="password" required />
            </div>
            <button type="submit">Login</button>
            <p v-if="message" :style="{
                color: isError ? 'red' : 'green'
            }">
                {{ message }}
            </p>
        </form>
    </div>
</template>