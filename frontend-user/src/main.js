import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import './style.css'
import { useAuth } from './stores/auth'
import Login from './views/Login.vue'
import Home from './views/Home.vue'
import Learn from './views/Learn.vue'
import ExamHall from './views/ExamHall.vue'
import Result from './views/Result.vue'
import Certs from './views/Certs.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: Login },
    { path: '/', component: Home },
    { path: '/learn/:id', component: Learn },
    { path: '/exam/:id', component: ExamHall },
    { path: '/result/:id', component: Result },
    { path: '/certs', component: Certs },
  ],
})
router.beforeEach((to) => {
  const auth = useAuth()
  if (to.path !== '/login' && !auth.token) return '/login'
  if (to.path === '/login' && auth.token) return '/'
})
createApp(App).use(createPinia()).use(router).mount('#app')
