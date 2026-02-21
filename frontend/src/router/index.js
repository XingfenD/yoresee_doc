import { createRouter, createWebHistory } from 'vue-router';
import { useUserStore } from '../store/user';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/Login.vue')
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('../views/Register.vue')
    },
    {
      path: '/',
      name: 'Home',
      component: () => import('../views/Home.vue'),
      meta: { requiresAuth: true }
    }
  ]
});

router.beforeEach((to, from) => {
  const userStore = useUserStore();
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  
  if (requiresAuth && !userStore.token) {
    return '/login';
  }
});

export default router;