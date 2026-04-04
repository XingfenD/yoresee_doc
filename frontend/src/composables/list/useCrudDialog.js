import { ref } from 'vue';

const toFormObject = (value) => ({ ...(value || {}) });

const resolveFactory = (valueOrFactory) => {
  if (typeof valueOrFactory === 'function') {
    return valueOrFactory();
  }
  return valueOrFactory || {};
};

export function useCrudDialog(options = {}) {
  const {
    initialForm = () => ({}),
    mapOpenForm = null,
    validate = null,
    submitRequest = null,
    closeOnSuccess = true,
    onSuccess = null,
    onError = null
  } = options;

  const visible = ref(false);
  const submitting = ref(false);
  const form = ref(toFormObject(resolveFactory(initialForm)));

  const reset = () => {
    form.value = toFormObject(resolveFactory(initialForm));
  };

  const open = (payload = undefined) => {
    if (typeof mapOpenForm === 'function') {
      form.value = toFormObject(mapOpenForm(payload));
    } else if (payload && typeof payload === 'object') {
      form.value = toFormObject(payload);
    } else {
      reset();
    }
    visible.value = true;
  };

  const close = () => {
    visible.value = false;
  };

  const submit = async () => {
    if (submitting.value) {
      return false;
    }
    if (typeof validate === 'function' && validate(form.value) === false) {
      return false;
    }
    if (typeof submitRequest !== 'function') {
      return false;
    }

    submitting.value = true;
    try {
      const result = await submitRequest(form.value);
      if (closeOnSuccess) {
        visible.value = false;
      }
      if (typeof onSuccess === 'function') {
        await onSuccess(result, form.value);
      }
      return true;
    } catch (error) {
      if (typeof onError === 'function') {
        onError(error, form.value);
      } else {
        console.error('[useCrudDialog] submit failed', error);
      }
      return false;
    } finally {
      submitting.value = false;
    }
  };

  return {
    visible,
    submitting,
    form,
    open,
    close,
    reset,
    submit
  };
}
