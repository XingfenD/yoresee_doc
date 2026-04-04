import { createRouter, createWebHistory } from 'vue-router';
import { useUserStore } from '../store/user';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/auth/Login.vue')
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('../views/auth/Register.vue')
    },
    {
      path: '/',
      name: 'Home',
      component: () => import('../views/workspace/Home.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/mydocuments',
      name: 'MyDocuments',
      component: () => import('../views/workspace/MyDocuments.vue'),
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
      component: () => import('../views/document/DocumentEditor.vue'),
      meta: { requiresAuth: true },
      props: (route) => ({ kbId: 'personal', docId: route.params.docId })
    },
    {
      path: '/mydocument/:docId/attachments',
      name: 'MyDocumentAttachments',
      component: () => import('../views/document/DocumentAttachments.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/mydocument/:docId/history',
      name: 'MyDocumentHistory',
      component: () => import('../views/document/DocumentHistory.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/mydocument/:docId/setting',
      name: 'MyDocumentSettings',
      component: () => import('../views/document/DocumentSettings.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/mydocument/:docId/attachment/:attachmentId',
      name: 'MyDocumentAttachmentOpen',
      component: () => import('../views/document/DocumentAttachmentOpen.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/knowledge-base',
      name: 'KnowledgeBase',
      component: () => import('../views/knowledge-base/KnowledgeBase.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/templates',
      name: 'Templates',
      component: () => import('../views/template/Templates.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/search',
      name: 'Search',
      component: () => import('../views/workspace/Search.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/template/:templateId',
      name: 'TemplatePreview',
      component: () => import('../views/template/TemplatePreview.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/knowledge-base/:id',
      name: 'KnowledgeBaseDetail',
      component: () => import('../views/knowledge-base/KnowledgeBaseDetail.vue'),
      meta: { requiresAuth: true },
      props: true
    },
    {
      path: '/knowledge-base/:kbId/document/:docId',
      name: 'KnowledgeBaseDocumentEditor',
      component: () => import('../views/document/DocumentEditor.vue'),
      meta: { requiresAuth: true },
      props: true
    },
    {
      path: '/knowledge-base/:kbId/document/:docId/attachments',
      name: 'KnowledgeBaseDocumentAttachments',
      component: () => import('../views/document/DocumentAttachments.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/knowledge-base/:kbId/document/:docId/history',
      name: 'KnowledgeBaseDocumentHistory',
      component: () => import('../views/document/DocumentHistory.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/knowledge-base/:kbId/document/:docId/setting',
      name: 'KnowledgeBaseDocumentSettings',
      component: () => import('../views/document/DocumentSettings.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/knowledge-base/:kbId/document/:docId/attachment/:attachmentId',
      name: 'KnowledgeBaseDocumentAttachmentOpen',
      component: () => import('../views/document/DocumentAttachmentOpen.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/user_info/profile',
      name: 'UserProfile',
      component: () => import('../views/user/UserProfile.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/user_info/setting',
      name: 'UserSetting',
      component: () => import('../views/user/UserSetting.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/user_info/invitations',
      name: 'UserInvitations',
      component: () => import('../views/user/UserInvitations.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/user_info/notifications',
      name: 'UserNotifications',
      component: () => import('../views/user/UserNotifications.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/manage/security',
      name: 'SystemManageSecurity',
      component: () => import('../views/manage/SystemManageSecurity.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/manage/user',
      name: 'SystemManageUser',
      component: () => import('../views/manage/SystemManageUser.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/manage/user_group',
      name: 'SystemManageUserGroup',
      component: () => import('../views/manage/SystemManageUserGroup.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/manage/user_group/:externalID',
      name: 'SystemManageUserGroupDetail',
      component: () => import('../views/manage/SystemManageUserGroupDetail.vue'),
      meta: { requiresAuth: true },
      props: true
    },
    {
      path: '/manage/organization',
      name: 'SystemManageOrganization',
      component: () => import('../views/manage/SystemManageOrganization.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/manage/organization/:externalID',
      name: 'SystemManageOrganizationDetail',
      component: () => import('../views/manage/SystemManageOrganizationDetail.vue'),
      meta: { requiresAuth: true },
      props: true
    },
    {
      path: '/manage/invitations',
      name: 'SystemManageInvitations',
      component: () => import('../views/manage/SystemManageInvitations.vue'),
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
