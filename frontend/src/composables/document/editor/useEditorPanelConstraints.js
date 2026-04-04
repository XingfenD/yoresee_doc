import { onBeforeUnmount, onMounted, watch } from 'vue';
import { usePanelSidebar } from '@/composables/layout/usePanelSidebar';

export function useEditorPanelConstraints(options = {}) {
  const {
    editorLayoutRef,
    commentSidebarRef,
    isCommentCollapsed,
    sidebarWidthStorageKey = 'docSidebarWidth',
    sidebarCollapsedStorageKey = 'sidebarCollapsed',
    minEditorMainWidth = 420
  } = options;

  const getVisibleCommentWidth = () => {
    if (isCommentCollapsed.value) {
      return 0;
    }
    const exposedWidth = Number(commentSidebarRef.value?.getCurrentWidth?.());
    if (Number.isFinite(exposedWidth) && exposedWidth > 0) {
      return exposedWidth;
    }
    const storedWidth = Number(localStorage.getItem('commentSidebarWidth'));
    if (Number.isFinite(storedWidth) && storedWidth > 0) {
      return Math.min(Math.max(storedWidth, 280), 560);
    }
    return 320;
  };

  const computeSidebarMaxWidth = () => {
    const layoutRect = editorLayoutRef.value?.getBoundingClientRect();
    if (!layoutRect) {
      return 520;
    }
    const available = layoutRect.width - getVisibleCommentWidth() - minEditorMainWidth - 8;
    return Math.min(520, Math.max(220, available));
  };

  const {
    width: sidebarWidth,
    collapsed: isSidebarCollapsed,
    resizing: isSidebarResizing,
    toggleCollapsed: toggleSidebar,
    startResize
  } = usePanelSidebar({
    defaultWidth: 280,
    minWidth: 220,
    maxWidth: 520,
    resizeEdge: 'right',
    collapsedStorageKey: sidebarCollapsedStorageKey,
    widthStorageKey: sidebarWidthStorageKey,
    getMaxWidth: computeSidebarMaxWidth,
    onWidthChange: (value) => {
      document.documentElement.style.setProperty('--sidebar-width', `${value}px`);
    }
  });

  const clampSidebarWidth = () => {
    const maxWidth = computeSidebarMaxWidth();
    if (sidebarWidth.value > maxWidth) {
      sidebarWidth.value = maxWidth;
    }
  };

  const handleCommentWidthChange = () => {
    clampSidebarWidth();
  };

  watch(
    () => isCommentCollapsed.value,
    () => {
      requestAnimationFrame(() => {
        clampSidebarWidth();
      });
    }
  );

  onMounted(() => {
    window.addEventListener('resize', clampSidebarWidth);
    clampSidebarWidth();
  });

  onBeforeUnmount(() => {
    window.removeEventListener('resize', clampSidebarWidth);
  });

  return {
    isSidebarCollapsed,
    isSidebarResizing,
    toggleSidebar,
    startResize,
    handleCommentWidthChange,
    clampSidebarWidth
  };
}
