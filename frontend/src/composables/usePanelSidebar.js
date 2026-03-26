import { computed, onBeforeUnmount, onMounted, ref, unref, watch } from 'vue';

export function usePanelSidebar(options = {}) {
  const {
    defaultWidth = 320,
    minWidth = 240,
    maxWidth = 560,
    resizeEdge = 'right',
    widthStorageKey = '',
    collapsedStorageKey = '',
    defaultCollapsed = false,
    externalCollapsed = null,
    getMaxWidth = null,
    onWidthChange = null
  } = options;

  const toValidWidth = (value) => {
    const dynamicMax = typeof getMaxWidth === 'function' ? getMaxWidth() : maxWidth;
    const hardMax = Math.max(minWidth, dynamicMax || maxWidth);
    return Math.min(Math.max(value, minWidth), hardMax);
  };

  const parseStoredBool = (raw, fallback) => {
    if (raw === 'true') return true;
    if (raw === 'false') return false;
    return fallback;
  };

  const storedWidth = Number(widthStorageKey ? localStorage.getItem(widthStorageKey) : NaN);
  const initialWidth = Number.isFinite(storedWidth) ? storedWidth : defaultWidth;
  const width = ref(toValidWidth(initialWidth));

  const innerCollapsed = ref(
    parseStoredBool(collapsedStorageKey ? localStorage.getItem(collapsedStorageKey) : null, defaultCollapsed)
  );
  const collapsed = computed(() => {
    if (externalCollapsed) {
      return !!unref(externalCollapsed);
    }
    return innerCollapsed.value;
  });

  const resizing = ref(false);
  const startX = ref(0);
  const startWidth = ref(width.value);

  const stopResize = () => {
    if (!resizing.value) return;
    resizing.value = false;
    document.body.style.cursor = '';
    document.body.style.userSelect = '';
    window.removeEventListener('mousemove', onResizeMove);
    window.removeEventListener('mouseup', stopResize);
  };

  const onResizeMove = (event) => {
    if (!resizing.value) return;
    const delta = resizeEdge === 'left'
      ? startX.value - event.clientX
      : event.clientX - startX.value;
    width.value = toValidWidth(startWidth.value + delta);
  };

  const startResize = (event) => {
    if (collapsed.value) return;
    event.preventDefault();
    resizing.value = true;
    startX.value = event.clientX;
    startWidth.value = width.value;
    document.body.style.cursor = 'col-resize';
    document.body.style.userSelect = 'none';
    window.addEventListener('mousemove', onResizeMove);
    window.addEventListener('mouseup', stopResize);
  };

  const setCollapsed = (value) => {
    if (externalCollapsed) {
      return;
    }
    innerCollapsed.value = !!value;
  };

  const toggleCollapsed = () => {
    if (externalCollapsed) {
      return;
    }
    innerCollapsed.value = !innerCollapsed.value;
  };

  watch(width, (value) => {
    if (widthStorageKey) {
      localStorage.setItem(widthStorageKey, `${value}`);
    }
    if (typeof onWidthChange === 'function') {
      onWidthChange(value);
    }
  });

  watch(collapsed, () => {
    stopResize();
  });

  if (!externalCollapsed && collapsedStorageKey) {
    watch(innerCollapsed, (value) => {
      localStorage.setItem(collapsedStorageKey, String(value));
    });
  }

  onMounted(() => {
    if (typeof onWidthChange === 'function') {
      onWidthChange(width.value);
    }
  });

  onBeforeUnmount(() => {
    stopResize();
  });

  return {
    width,
    collapsed,
    resizing,
    setCollapsed,
    toggleCollapsed,
    startResize,
    stopResize
  };
}
