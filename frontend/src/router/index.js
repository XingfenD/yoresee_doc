import { createRouter, createWebHistory } from 'vue-router';
import { useUserStore } from '../store/user';
import i18n from '../i18n';

const { t } = i18n.global;

const APP_NAME = 'Yoresee';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/auth/Login.vue'),
      meta: { titleKey: 'login.title' }
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('../views/auth/Register.vue'),
      meta: { titleKey: 'register.title' }
    },
    {
      path: '/',
      name: 'Home',
      component: () => import('../views/workspace/Home.vue'),
      meta: { requiresAuth: true, titleKey: 'navigation.home' }
    },
    {
      path: '/mydocuments',
      name: 'MyDocuments',
      component: () => import('../views/workspace/MyDocuments.vue'),
      meta: { requiresAuth: true, titleKey: 'navigation.myDocuments' }
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
      meta: { requiresAuth: true, dynamicTitle: true },
      props: (route) => ({ kbId: 'personal', docId: route.params.docId })
    },
    {
      path: '/mydocument/:docId/attachments',
      name: 'MyDocumentAttachments',
      component: () => import('../views/document/DocumentAttachments.vue'),
      meta: { requiresAuth: true, titleKey: 'document.attachments.title' }
    },
    {
      path: '/mydocument/:docId/history',
      name: 'MyDocumentHistory',
      component: () => import('../views/document/DocumentHistory.vue'),
      meta: { requiresAuth: true, titleKey: 'history.title' }
    },
    {
      path: '/mydocument/:docId/setting',
      name: 'MyDocumentSettings',
      component: () => import('../views/document/DocumentSettings.vue'),
      meta: { requiresAuth: true, titleKey: 'document.settings.title' }
    },
    {
      path: '/mydocument/:docId/attachment/:attachmentId',
      name: 'MyDocumentAttachmentOpen',
      component: () => import('../views/document/DocumentAttachmentOpen.vue'),
      meta: { requiresAuth: true, titleKey: 'document.attachments.previewTitle' }
    },
    {
      path: '/knowledge-base',
      name: 'KnowledgeBase',
      component: () => import('../views/knowledge-base/KnowledgeBase.vue'),
      meta: { requiresAuth: true, titleKey: 'knowledgeBase.title' }
    },
    {
      path: '/templates',
      name: 'Templates',
      component: () => import('../views/template/Templates.vue'),
      meta: { requiresAuth: true, titleKey: 'templates.title' }
    },
    {
      path: '/search',
      name: 'Search',
      component: () => import('../views/workspace/Search.vue'),
      meta: { requiresAuth: true, titleKey: 'search.title' }
    },
    {
      path: '/template/:templateId',
      name: 'TemplatePreview',
      component: () => import('../views/template/TemplatePreview.vue'),
      meta: { requiresAuth: true, dynamicTitle: true }
    },
    {
      path: '/knowledge-base/:id',
      name: 'KnowledgeBaseDetail',
      component: () => import('../views/knowledge-base/KnowledgeBaseDetail.vue'),
      meta: { requiresAuth: true, dynamicTitle: true },
      props: true
    },
    {
      path: '/knowledge-base/:kbId/document/:docId',
      name: 'KnowledgeBaseDocumentEditor',
      component: () => import('../views/document/DocumentEditor.vue'),
      meta: { requiresAuth: true, dynamicTitle: true },
      props: true
    },
    {
      path: '/knowledge-base/:kbId/document/:docId/attachments',
      name: 'KnowledgeBaseDocumentAttachments',
      component: () => import('../views/document/DocumentAttachments.vue'),
      meta: { requiresAuth: true, titleKey: 'document.attachments.title' }
    },
    {
      path: '/knowledge-base/:kbId/document/:docId/history',
      name: 'KnowledgeBaseDocumentHistory',
      component: () => import('../views/document/DocumentHistory.vue'),
      meta: { requiresAuth: true, titleKey: 'history.title' }
    },
    {
      path: '/knowledge-base/:kbId/document/:docId/setting',
      name: 'KnowledgeBaseDocumentSettings',
      component: () => import('../views/document/DocumentSettings.vue'),
      meta: { requiresAuth: true, titleKey: 'document.settings.title' }
    },
    {
      path: '/knowledge-base/:kbId/document/:docId/attachment/:attachmentId',
      name: 'KnowledgeBaseDocumentAttachmentOpen',
      component: () => import('../views/document/DocumentAttachmentOpen.vue'),
      meta: { requiresAuth: true, titleKey: 'document.attachments.previewTitle' }
    },
    {
      path: '/user_info/profile',
      name: 'UserProfile',
      component: () => import('../views/user/UserProfile.vue'),
      meta: { requiresAuth: true, titleKey: 'user.profileTitle' }
    },
    {
      path: '/user_info/setting',
      name: 'UserSetting',
      component: () => import('../views/user/UserSetting.vue'),
      meta: { requiresAuth: true, titleKey: 'user.settingTitle' }
    },
    {
      path: '/user_info/invitations',
      name: 'UserInvitations',
      component: () => import('../views/user/UserInvitations.vue'),
      meta: { requiresAuth: true, titleKey: 'user.invite.title' }
    },
    {
      path: '/user_info/notifications',
      name: 'UserNotifications',
      component: () => import('../views/user/UserNotifications.vue'),
      meta: { requiresAuth: true, titleKey: 'user.notifications.title' }
    },
    {
      path: '/manage/security',
      name: 'SystemManageSecurity',
      component: () => import('../views/manage/SystemManageSecurity.vue'),
      meta: { requiresAuth: true, titleKey: 'system.security.title' }
    },
    {
      path: '/manage/user',
      name: 'SystemManageUser',
      component: () => import('../views/manage/SystemManageUser.vue'),
      meta: { requiresAuth: true, titleKey: 'system.user.title' }
    },
    {
      path: '/manage/user_group',
      name: 'SystemManageUserGroup',
      component: () => import('../views/manage/SystemManageUserGroup.vue'),
      meta: { requiresAuth: true, titleKey: 'system.userGroup.title' }
    },
    {
      path: '/manage/user_group/:externalID',
      name: 'SystemManageUserGroupDetail',
      component: () => import('../views/manage/SystemManageUserGroupDetail.vue'),
      meta: { requiresAuth: true, dynamicTitle: true },
      props: true
    },
    {
      path: '/manage/organization',
      name: 'SystemManageOrganization',
      component: () => import('../views/manage/SystemManageOrganization.vue'),
      meta: { requiresAuth: true, titleKey: 'system.organization.title' }
    },
    {
      path: '/manage/organization/:externalID',
      name: 'SystemManageOrganizationDetail',
      component: () => import('../views/manage/SystemManageOrganizationDetail.vue'),
      meta: { requiresAuth: true, dynamicTitle: true },
      props: true
    },
    {
      path: '/manage/invitations',
      name: 'SystemManageInvitations',
      component: () => import('../views/manage/SystemManageInvitations.vue'),
      meta: { requiresAuth: true, titleKey: 'system.invite.title' }
    },
    {
      path: '/manage',
      name: 'SystemManage',
      redirect: '/manage/user'
    },
    {
      path: '/404',
      name: 'NotFound',
      component: () => import('../views/error/NotFound.vue'),
      meta: { titleKey: '404' }
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: (to) => ({
        path: '/404',
        query: {
          from: to.fullPath
        }
      })
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

router.afterEach((to) => {
  if (to.meta.dynamicTitle) return;
  const titleKey = to.meta.titleKey;
  if (titleKey) {
    const pageTitle = t(titleKey);
    document.title = `${pageTitle} - ${APP_NAME}`;
  } else {
    document.title = APP_NAME;
  }
});

export default router;
