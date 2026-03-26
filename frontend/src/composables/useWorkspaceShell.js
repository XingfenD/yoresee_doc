import { Collection, Document, House, Tickets } from '@element-plus/icons-vue';
import { useAppShellBase } from '@/composables/useAppShellBase';

export function useWorkspaceShell({
  locale,
  router,
  userStore,
  defaultActiveMenu = 'home'
}) {
  const workspaceMenuItems = [
    { key: 'home', labelKey: 'navigation.home', icon: House, route: '/' },
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

  const fetchSystemInfo = async () => {
    try {
      const info = await userStore.fetchSystemInfo();
      shell.systemName.value = info.system_name;
    } catch (err) {
      console.error('获取系统信息失败:', err);
    }
  };

  return {
    ...shell,
    workspaceMenuItems,
    fetchSystemInfo
  };
}
