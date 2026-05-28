import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../modules/HomeView.vue'
import LoginView from '@/views/auth/LoginView.vue'
import DefaultLayout from '@/layouts/DefaultLayout.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },

    {
      path: '/',
      component: DefaultLayout,
      children: [
        {
          path: '',
          name: 'home',
          component: () => import('../modules/HomeView.vue'),
        },
        {
          path: 'about',
          component: () => import('../modules/AboutView.vue'),
        },
      ],
    }
  ],
})

export default router
