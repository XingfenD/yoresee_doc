import { ElMessage } from 'element-plus';
import { useAsyncAction } from '@/composables/useAsyncAction';

export const isActionCancelled = (error) => {
  if (error === null || error === undefined) {
    return false;
  }
  if (typeof error === 'string') {
    return error === 'cancel' || error === 'close';
  }
  const message = error?.message || '';
  return message === 'cancel' || message === 'close';
};

export function useApiAction(options = {}) {
  const {
    t,
    defaultErrorKey = 'common.requestFailed',
    defaultErrorMessage = 'Request failed'
  } = options;

  const { pending, runAsync } = useAsyncAction();

  const resolveErrorMessage = (message) => {
    if (message) {
      return message;
    }
    if (typeof t === 'function' && defaultErrorKey) {
      return t(defaultErrorKey);
    }
    return defaultErrorMessage;
  };

  const shouldIgnore = (error, ignoreError) => {
    if (typeof ignoreError === 'function') {
      return ignoreError(error);
    }
    return Boolean(ignoreError);
  };

  const handleApiError = (error, options = {}) => {
    const {
      context = 'api request',
      errorMessage,
      showErrorMessage = true,
      logError = true,
      ignoreError = false
    } = options;

    if (shouldIgnore(error, ignoreError)) {
      return;
    }

    if (logError) {
      console.error(`${context} failed`, error);
    }
    if (showErrorMessage) {
      ElMessage.error(resolveErrorMessage(errorMessage));
    }
  };

  const createApiErrorHandler = (options = {}) => (error) => {
    handleApiError(error, options);
  };

  const runApi = async (action, options = {}) => {
    const {
      successMessage,
      onSuccess,
      onError,
      onFinally,
      rethrow = false,
      fallback = null
    } = options;

    return runAsync(action, {
      rethrow,
      fallback,
      onSuccess: async (result) => {
        if (successMessage) {
          ElMessage.success(successMessage);
        }
        if (typeof onSuccess === 'function') {
          await onSuccess(result);
        }
      },
      onError: async (error) => {
        handleApiError(error, options);
        if (typeof onError === 'function') {
          await onError(error);
        }
      },
      onFinally
    });
  };

  const runSilent = async (action, options = {}) =>
    runApi(action, {
      showErrorMessage: false,
      logError: false,
      ...options
    });

  const runWithLoading = async (loadingRef, action, options = {}) => {
    if (!loadingRef || typeof loadingRef.value !== 'boolean') {
      return runApi(action, options);
    }
    if (loadingRef.value) {
      return options.fallback ?? null;
    }
    loadingRef.value = true;
    const originalOnFinally = options.onFinally;
    try {
      return await runApi(action, {
        ...options,
        onFinally: async () => {
          if (typeof originalOnFinally === 'function') {
            await originalOnFinally();
          }
        }
      });
    } finally {
      loadingRef.value = false;
    }
  };

  return {
    pending,
    runApi,
    runSilent,
    runWithLoading,
    handleApiError,
    createApiErrorHandler
  };
}
