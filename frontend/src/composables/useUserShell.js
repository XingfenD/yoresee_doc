import { computed, ref } from 'vue';
import { House, User, Ticket, Setting, Bell } from '@element-plus/icons-vue';

const DEFAULT_AVATAR = 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png';

export function useUserShell({ locale, router, userStore, defaultActiveMenu = 'user-center' }) {
  const systemName = ref('Yoresee');
  const activeMenu = ref(defaultActiveMenu);

  const isDarkMode = computed(() => userStore.darkMode);
  const userInfo = computed(() => userStore.userInfo);
  const userAvatar = computed(() => userInfo.value?.avatar || DEFAULT_AVATAR);

  const userMenuItems = [
    { key: 'home', labelKey: 'navigation.home', icon: House, route: '/' },
    { key: 'user-center', labelKey: 'user.menu.center', icon: User, route: '/user_info/example' },
    { key: 'user-notifications', labelKey: 'user.menu.notifications', icon: Bell, route: '/user_info/notifications' },
    { key: 'user-invite', labelKey: 'user.menu.invite', icon: Ticket, route: '/user_info/invatations' },
    { key: 'user-security', labelKey: 'user.menu.security', icon: Setting, route: '/user_info/example' }
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
    userMenuItems,
    currentLanguage,
    initLanguage,
    handleLanguageChange,
    toggleTheme,
    handleLogout,
    handleMenuSelect
  };
}
