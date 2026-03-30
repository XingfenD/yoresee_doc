import { House, Search, Setting, Ticket, User, UserFilled, OfficeBuilding } from '@element-plus/icons-vue';
import { useAppShellBase } from '@/composables/useAppShellBase';

export function useManageShell({ locale, router, userStore, defaultActiveMenu }) {
  const manageMenuItems = [
    { key: 'home', labelKey: 'navigation.home', icon: House, route: '/' },
    { key: 'search', labelKey: 'navigation.search', icon: Search, route: '/search' },
    { key: 'manage-user', labelKey: 'system.menu.user', icon: User, route: '/manage/user' },
    { key: 'manage-user-group', labelKey: 'system.menu.userGroup', icon: UserFilled, route: '/manage/user_group' },
    { key: 'manage-organization', labelKey: 'system.menu.organization', icon: OfficeBuilding, route: '/manage/organization' },
    { key: 'manage-invite', labelKey: 'system.menu.invite', icon: Ticket, route: '/manage/invitations' },
    { key: 'manage-security', labelKey: 'system.menu.security', icon: Setting, route: '/manage/security' }
  ];

  const shell = useAppShellBase({
    locale,
    router,
    userStore,
    defaultActiveMenu,
    menuItems: manageMenuItems
  });

  return {
    ...shell,
    manageMenuItems,
  };
}
