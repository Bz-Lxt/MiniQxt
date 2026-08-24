import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import './style.css'
import { useAuth } from './stores/auth'
import Login from './views/Login.vue'
import Shell from './layouts/Shell.vue'
import Dashboard from './views/Dashboard.vue'
import Org from './views/Org.vue'
import Questions from './views/Questions.vue'
import PaperBoard from './views/PaperBoard.vue'
import Exams from './views/Exams.vue'
import ManualGrade from './views/ManualGrade.vue'
import Audit from './views/Audit.vue'
import Courses from './views/Courses.vue'
import Certs from './views/Certs.vue'
import Queue from './views/Queue.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: Login },
    {
      path: '/',
      component: Shell,
      children: [
        { path: '', component: Dashboard },
        { path: 'org', component: Org },
        { path: 'questions', component: Questions },
        { path: 'papers', component: PaperBoard },
        { path: 'exams', component: Exams },
        { path: 'grade', component: ManualGrade },
        { path: 'audit', component: Audit },
        { path: 'audit/:id', component: Audit },
        { path: 'courses', component: Courses },
        { path: 'certs', component: Certs },
        { path: 'queue', component: Queue },
      ],
    },
  ],
})

router.beforeEach((to) => {
  const auth = useAuth()
  if (to.path !== '/login' && !auth.token) return '/login'
  if (to.path === '/login' && auth.token) return '/'
})

createApp(App).use(createPinia()).use(router).mount('#app')
