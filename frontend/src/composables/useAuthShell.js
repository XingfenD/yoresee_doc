import { computed, ref } from 'vue';

export function useAuthShell({ locale, userStore }) {
  const systemName = ref('Yoresee');
  const isDarkMode = computed(() => userStore.darkMode);
  const currentLanguage = ref(localStorage.getItem('language') || 'en');

  const handleLanguageChange = (command) => {
    currentLanguage.value = command;
    locale.value = command;
    localStorage.setItem('language', command);
  };

  const toggleTheme = () => {
    userStore.toggleDarkMode();
  };

  const initLanguage = () => {
    const savedLanguage = localStorage.getItem('language');
    if (savedLanguage) {
      currentLanguage.value = savedLanguage;
      locale.value = savedLanguage;
    }
  };

  const fetchSystemInfo = async (onLoaded) => {
    try {
      const info = await userStore.fetchSystemInfo();
      systemName.value = info.system_name;
      onLoaded?.(info);
    } catch (err) {
      console.error('获取系统信息失败:', err);
    }
  };

  return {
    systemName,
    isDarkMode,
    currentLanguage,
    handleLanguageChange,
    toggleTheme,
    initLanguage,
    fetchSystemInfo
  };
}
