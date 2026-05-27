<script setup lang="ts">
import {ref} from 'vue'
import { loginApi } from '../service/auth'
const username = ref('')
const password = ref('')

const login = async () => {
    try {
        await loginApi(username.value, password.value)
    } catch (error) {
        console.error('Login failed:', error)
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
            <button type="submit" >Login</button>
        </form>
    </div>
</template>