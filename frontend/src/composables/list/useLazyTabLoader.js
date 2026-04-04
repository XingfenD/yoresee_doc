import { ref, watch } from 'vue';

export function useLazyTabLoader(options = {}) {
  const { initialTab = '', tabs = {} } = options;
  const activeTab = ref(initialTab);

  const ensureTabLoaded = async (tab = activeTab.value) => {
    const target = tabs[tab];
    if (!target || typeof target.load !== 'function') return;
    if (typeof target.loaded === 'function' && target.loaded()) return;
    await target.load();
  };

  watch(activeTab, async (tab) => {
    await ensureTabLoaded(tab);
  });

  return {
    activeTab,
    ensureTabLoaded
  };
}
