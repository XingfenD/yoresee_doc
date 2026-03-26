import { House, User, Ticket, Setting, Bell } from '@element-plus/icons-vue';
import { useAppShellBase } from '@/composables/useAppShellBase';

export function useUserShell({ locale, router, userStore, defaultActiveMenu = 'user-center' }) {
  const userMenuItems = [
    { key: 'home', labelKey: 'navigation.home', icon: House, route: '/' },
    { key: 'user-center', labelKey: 'user.menu.center', icon: User, route: '/user_info/example' },
    { key: 'user-notifications', labelKey: 'user.menu.notifications', icon: Bell, route: '/user_info/notifications' },
    { key: 'user-invite', labelKey: 'user.menu.invite', icon: Ticket, route: '/user_info/invatations' },
    { key: 'user-security', labelKey: 'user.menu.security', icon: Setting, route: '/user_info/example' }
  ];

  const shell = useAppShellBase({
    locale,
    router,
    userStore,
    defaultActiveMenu,
    menuItems: userMenuItems
  });

  return {
    ...shell,
    userMenuItems,
  };
}
