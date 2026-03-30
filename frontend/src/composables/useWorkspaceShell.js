import { Collection, Document, House, Search, Tickets } from '@element-plus/icons-vue';
import { useAppShellBase } from '@/composables/useAppShellBase';

export function useWorkspaceShell({
  locale,
  router,
  userStore,
  defaultActiveMenu = 'home'
}) {
  const workspaceMenuItems = [
    { key: 'home', labelKey: 'navigation.home', icon: House, route: '/' },
    { key: 'search', labelKey: 'navigation.search', icon: Search, route: '/search' },
    { key: 'documents', labelKey: 'navigation.myDocuments', icon: Document, route: '/mydocuments' },
    { key: 'knowledge-base', labelKey: 'navigation.knowledgeBase', icon: Collection, route: '/knowledge-base' },
    { key: 'templates', labelKey: 'navigation.templates', icon: Tickets, route: '/templates' }
  ];

  const shell = useAppShellBase({
    locale,
    router,
    userStore,
    defaultActiveMenu,
    menuItems: workspaceMenuItems
  });

  return {
    ...shell,
    workspaceMenuItems,
  };
}
