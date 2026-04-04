import { computed, ref } from 'vue';

export function useAsyncAction() {
  const pendingCount = ref(0);
  const pending = computed(() => pendingCount.value > 0);

  const runAsync = async (action, options = {}) => {
    const {
      onSuccess,
      onError,
      onFinally,
      rethrow = false,
      fallback = null
    } = options;

    if (typeof action !== 'function') {
      return fallback;
    }

    pendingCount.value += 1;
    try {
      const result = await action();
      if (typeof onSuccess === 'function') {
        await onSuccess(result);
      }
      return result;
    } catch (error) {
      if (typeof onError === 'function') {
        await onError(error);
      }
      if (rethrow) {
        throw error;
      }
      return fallback;
    } finally {
      pendingCount.value = Math.max(0, pendingCount.value - 1);
      if (typeof onFinally === 'function') {
        await onFinally();
      }
    }
  };

  return {
    pending,
    runAsync
  };
}
