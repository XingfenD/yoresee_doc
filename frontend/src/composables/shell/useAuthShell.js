import { useAppShellBase } from '@/composables/shell/useAppShellBase';

export function useAuthShell({ locale, userStore }) {
  const shell = useAppShellBase({
    locale,
    // auth 场景没有侧边栏路由跳转需求，这里提供一个空 router 适配基座 API
    router: { push: () => {} },
    userStore,
    defaultActiveMenu: 'home',
    menuItems: []
  });

  return {
    systemName: shell.systemName,
    isDarkMode: shell.isDarkMode,
    currentLanguage: shell.currentLanguage,
    handleLanguageChange: shell.handleLanguageChange,
    toggleTheme: shell.toggleTheme,
    initLanguage: shell.initLanguage,
    fetchSystemInfo: shell.fetchSystemInfo
  };
}
