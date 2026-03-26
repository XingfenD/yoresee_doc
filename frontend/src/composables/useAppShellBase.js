import { computed, ref } from 'vue';

const DEFAULT_AVATAR = 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png';

export function useAppShellBase({ locale, router, userStore, defaultActiveMenu, menuItems = [] }) {
  const systemName = ref('Yoresee');
  const activeMenu = ref(defaultActiveMenu);
  const isDarkMode = computed(() => userStore.darkMode);
  const userInfo = computed(() => userStore.userInfo);
  const userAvatar = computed(() => userInfo.value?.avatar || DEFAULT_AVATAR);

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
    menuItems,
    currentLanguage,
    initLanguage,
    handleLanguageChange,
    toggleTheme,
    handleLogout,
    handleMenuSelect
  };
}
