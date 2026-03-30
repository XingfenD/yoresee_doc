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
    },
    {
      path: '/mydocuments',
      name: 'MyDocuments',
      component: () => import('../views/MyDocuments.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/documents',
      name: 'MyDocumentsAlias',
      redirect: '/mydocuments'
    },
    {
      path: '/mydocument/:docId',
      name: 'MyDocumentEditor',
      component: () => import('../views/DocumentEditor.vue'),
      meta: { requiresAuth: true },
      props: (route) => ({ kbId: 'personal', docId: route.params.docId })
    },
    {
      path: '/knowledge-base',
      name: 'KnowledgeBase',
      component: () => import('../views/KnowledgeBase.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/templates',
      name: 'Templates',
      component: () => import('../views/Templates.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/search',
      name: 'Search',
      component: () => import('../views/Search.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/template/:templateId',
      name: 'TemplatePreview',
      component: () => import('../views/TemplatePreview.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/knowledge-base/:id',
      name: 'KnowledgeBaseDetail',
      component: () => import('../views/KnowledgeBaseDetail.vue'),
      meta: { requiresAuth: true },
      props: true
    },
    {
      path: '/knowledge-base/:kbId/document/:docId',
      name: 'KnowledgeBaseDocumentEditor',
      component: () => import('../views/DocumentEditor.vue'),
      meta: { requiresAuth: true },
      props: true
    },
    {
      path: '/user_info/example',
      name: 'UserInfoExample',
      component: () => import('../views/UserInfoExample.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/user_info/invitations',
      name: 'UserInvitations',
      component: () => import('../views/UserInvitations.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/user_info/notifications',
      name: 'UserNotifications',
      component: () => import('../views/UserNotifications.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/manage/security',
      name: 'SystemManageSecurity',
      component: () => import('../views/SystemManageSecurity.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/manage/user',
      name: 'SystemManageUser',
      component: () => import('../views/SystemManageUser.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/manage/user_group',
      name: 'SystemManageUserGroup',
      component: () => import('../views/SystemManageUserGroup.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/manage/user_group/:externalID',
      name: 'SystemManageUserGroupDetail',
      component: () => import('../views/SystemManageUserGroupDetail.vue'),
      meta: { requiresAuth: true },
      props: true
    },
    {
      path: '/manage/organization',
      name: 'SystemManageOrganization',
      component: () => import('../views/SystemManageOrganization.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/manage/organization/:externalID',
      name: 'SystemManageOrganizationDetail',
      component: () => import('../views/SystemManageOrganizationDetail.vue'),
      meta: { requiresAuth: true },
      props: true
    },
    {
      path: '/manage/invitations',
      name: 'SystemManageInvitations',
      component: () => import('../views/SystemManageInvitations.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/manage',
      name: 'SystemManage',
      redirect: '/manage/user'
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
