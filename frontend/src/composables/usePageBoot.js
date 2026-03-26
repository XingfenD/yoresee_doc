import { ref } from 'vue';

export function usePageBoot(options = {}) {
  const { initLanguage, fetchSystemInfo } = options;
  const booting = ref(false);
  const booted = ref(false);

  const boot = async (...tasks) => {
    if (booting.value) {
      return;
    }
    booting.value = true;
    try {
      if (typeof initLanguage === 'function') {
        initLanguage();
      }
      if (typeof fetchSystemInfo === 'function') {
        await fetchSystemInfo();
      }
      for (const task of tasks) {
        if (typeof task === 'function') {
          await task();
        }
      }
      booted.value = true;
    } finally {
      booting.value = false;
    }
  };

  return {
    boot,
    booting,
    booted
  };
}
