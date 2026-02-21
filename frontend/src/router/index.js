import { createRouter, createWebHistory } from 'vue-router';
import Login from '@/views/Login.vue';
import Home from '@/views/Home.vue';

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const token = localStorage.getItem('token');
  
  if (requiresAuth && !token) {
    next('/login');
  } else if (!requiresAuth && token && to.path === '/login') {
    next('/');
  } else {
    next();
  }
});

export default router;