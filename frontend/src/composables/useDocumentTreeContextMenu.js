import { computed, onBeforeUnmount, onMounted, ref } from 'vue';

const resolveBool = (valueOrGetterOrRef, fallback = false) => {
  if (typeof valueOrGetterOrRef === 'function') {
    return valueOrGetterOrRef();
  }
  if (valueOrGetterOrRef && typeof valueOrGetterOrRef === 'object' && 'value' in valueOrGetterOrRef) {
    return valueOrGetterOrRef.value;
  }
  if (valueOrGetterOrRef === undefined) {
    return fallback;
  }
  return Boolean(valueOrGetterOrRef);
};

export function useDocumentTreeContextMenu(options = {}) {
  const { enabled = true } = options;

  const contextMenu = ref({
    visible: false,
    x: 0,
    y: 0,
    data: null
  });

  const closeContextMenu = () => {
    if (!contextMenu.value.visible) {
      return;
    }
    contextMenu.value.visible = false;
  };

  const openContextMenu = (event, data) => {
    if (!resolveBool(enabled, true)) {
      return false;
    }

    event.preventDefault();
    event.stopPropagation();

    const menuWidth = 150;
    const menuHeight = 120;
    let x = event.clientX;
    let y = event.clientY;
    if (x + menuWidth > window.innerWidth) {
      x = window.innerWidth - menuWidth - 8;
    }
    if (y + menuHeight > window.innerHeight) {
      y = window.innerHeight - menuHeight - 8;
    }

    contextMenu.value = {
      visible: true,
      x,
      y,
      data
    };

    return true;
  };

  const contextMenuStyle = computed(() => ({
    left: `${contextMenu.value.x}px`,
    top: `${contextMenu.value.y}px`
  }));

  onMounted(() => {
    window.addEventListener('click', closeContextMenu);
    window.addEventListener('scroll', closeContextMenu, true);
  });

  onBeforeUnmount(() => {
    window.removeEventListener('click', closeContextMenu);
    window.removeEventListener('scroll', closeContextMenu, true);
  });

  return {
    contextMenu,
    contextMenuStyle,
    openContextMenu,
    closeContextMenu
  };
}
