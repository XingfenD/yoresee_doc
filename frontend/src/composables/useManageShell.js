import { computed, ref } from 'vue';
import { House, Setting, Ticket, User, UserFilled, OfficeBuilding } from '@element-plus/icons-vue';

const DEFAULT_AVATAR = 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png';

export function useManageShell({ locale, router, userStore, defaultActiveMenu }) {
  const systemName = ref('Yoresee');
  const activeMenu = ref(defaultActiveMenu);
  const isDarkMode = computed(() => userStore.darkMode);

  const userInfo = computed(() => userStore.userInfo);
  const userAvatar = computed(() => userInfo.value?.avatar || DEFAULT_AVATAR);

  const manageMenuItems = [
    { key: 'home', labelKey: 'navigation.home', icon: House, route: '/' },
    { key: 'manage-user', labelKey: 'system.menu.user', icon: User, route: '/manage/user' },
    { key: 'manage-user-group', labelKey: 'system.menu.userGroup', icon: UserFilled, route: '/manage/user_group' },
    { key: 'manage-organization', labelKey: 'system.menu.organization', icon: OfficeBuilding, route: '/manage/organization' },
    { key: 'manage-invite', labelKey: 'system.menu.invite', icon: Ticket, route: '/manage/invitations' },
    { key: 'manage-security', labelKey: 'system.menu.security', icon: Setting, route: '/manage/security' }
  ];

  const currentLanguage = computed({
    get: () => locale.value,
    set: (value) => {
      locale.value = value;
      localStorage.setItem('language', value);
    }
  });

  const initLanguage = () => {
    const savedLanguage = localStorage.getItem('language');
    if (savedLanguage) {
      currentLanguage.value = savedLanguage;
    }
  };

  const handleLanguageChange = (command) => {
    currentLanguage.value = command;
  };

  const toggleTheme = () => {
    userStore.toggleDarkMode();
  };

  const handleLogout = () => {
    userStore.logout();
    router.push('/login');
  };

  const handleMenuSelect = (key) => {
    activeMenu.value = key;
  };

  return {
    systemName,
    activeMenu,
    isDarkMode,
    userInfo,
    userAvatar,
    manageMenuItems,
    currentLanguage,
    initLanguage,
    handleLanguageChange,
    toggleTheme,
    handleLogout,
    handleMenuSelect
  };
}
