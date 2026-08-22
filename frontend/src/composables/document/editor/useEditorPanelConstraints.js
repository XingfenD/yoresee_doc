import { onBeforeUnmount, onMounted, watch } from 'vue';
import { usePanelSidebar } from '@/composables/layout/usePanelSidebar';

export function useEditorPanelConstraints(options = {}) {
  const {
    editorLayoutRef,
    commentSidebarRef,
    isCommentCollapsed,
    onLayoutChange = null,
    sidebarWidthStorageKey = 'docSidebarWidth',
    sidebarCollapsedStorageKey = 'sidebarCollapsed',
    minEditorMainWidth = 520,
    pageWidth = 720
  } = options;
  const COLLAPSE_ANIMATION_MS = 320;
  const COMMENT_RESIZE_DEBOUNCE_MS = 140;
  const LAYOUT_GUTTER = 8;
  const PAGE_MIN_MARGIN = 16;
  let layoutSettleTimer = null;
  let resizeObserver = null;

  const emitLayoutChange = (delay = 0) => {
    if (layoutSettleTimer) {
      clearTimeout(layoutSettleTimer);
      layoutSettleTimer = null;
    }
    layoutSettleTimer = window.setTimeout(() => {
      layoutSettleTimer = null;
      if (typeof onLayoutChange === 'function') {
        onLayoutChange();
      }
    }, Math.max(0, delay));
  };

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
    setCollapsed: setSidebarCollapsed,
    toggleCollapsed: toggleSidebar,
    startResize
  } = usePanelSidebar({
    defaultWidth: 280,
    minWidth: 180,
    maxWidth: 520,
    resizeEdge: 'right',
    collapsedStorageKey: sidebarCollapsedStorageKey,
    widthStorageKey: sidebarWidthStorageKey,
    getMaxWidth: computeSidebarMaxWidth,
    onWidthChange: (value) => {
      document.documentElement.style.setProperty('--sidebar-width', `${value}px`);
      if (!isSidebarResizing.value) {
        emitLayoutChange(0);
      }
    }
  });

  const clampSidebarWidth = () => {
    const maxWidth = computeSidebarMaxWidth();
    if (sidebarWidth.value > maxWidth) {
      sidebarWidth.value = maxWidth;
    }
  };

  // When the editing area is too narrow to fit the fixed-width document page,
  // reclaim space by collapsing the directory and comment sidebars first.
  // The document then falls back to a horizontal scrollbar if still narrower.
  // Sidebars we auto-collapsed are restored when the area widens again, while
  // sidebars the user collapsed manually are left alone.
  let responsiveDirForced = false;
  let responsiveCommentForced = false;
  const applyResponsiveSidebars = () => {
    const layoutRect = editorLayoutRef.value?.getBoundingClientRect();
    if (!layoutRect || layoutRect.width <= 0) {
      return;
    }
    const dirVisible = isSidebarCollapsed.value ? 0 : sidebarWidth.value;
    const commentVisible = getVisibleCommentWidth();
    const editorArea = layoutRect.width - commentVisible - dirVisible - LAYOUT_GUTTER;
    const narrow = editorArea > 0 && editorArea < pageWidth + PAGE_MIN_MARGIN;
    if (narrow) {
      if (!isSidebarCollapsed.value) {
        setSidebarCollapsed(true);
        responsiveDirForced = true;
      }
      if (!isCommentCollapsed.value) {
        isCommentCollapsed.value = true;
        responsiveCommentForced = true;
      }
    } else {
      if (responsiveDirForced) {
        setSidebarCollapsed(false);
        responsiveDirForced = false;
      }
      if (responsiveCommentForced) {
        isCommentCollapsed.value = false;
        responsiveCommentForced = false;
      }
    }
  };

  const handleCommentWidthChange = () => {
    clampSidebarWidth();
    emitLayoutChange(COMMENT_RESIZE_DEBOUNCE_MS);
  };

  watch(
    () => isSidebarCollapsed.value,
    () => {
      requestAnimationFrame(() => {
        clampSidebarWidth();
        emitLayoutChange(COLLAPSE_ANIMATION_MS);
      });
    }
  );

  watch(
    () => isSidebarResizing.value,
    (resizing) => {
      if (!resizing) {
        emitLayoutChange(0);
      }
    }
  );

  watch(
    () => isCommentCollapsed.value,
    () => {
      requestAnimationFrame(() => {
        clampSidebarWidth();
        emitLayoutChange(COLLAPSE_ANIMATION_MS);
      });
    }
  );

  onMounted(() => {
    // Always start with directory sidebar expanded when entering editor.
    setSidebarCollapsed(false);
    window.addEventListener('resize', handleWindowResize);
    if (typeof ResizeObserver !== 'undefined' && editorLayoutRef.value) {
      resizeObserver = new ResizeObserver(handleLayoutResize);
      resizeObserver.observe(editorLayoutRef.value);
    }
    clampSidebarWidth();
    applyResponsiveSidebars();
    emitLayoutChange(0);
  });

  onBeforeUnmount(() => {
    window.removeEventListener('resize', handleWindowResize);
    if (resizeObserver) {
      resizeObserver.disconnect();
      resizeObserver = null;
    }
    if (layoutSettleTimer) {
      clearTimeout(layoutSettleTimer);
      layoutSettleTimer = null;
    }
  });

  const handleWindowResize = () => {
    clampSidebarWidth();
    applyResponsiveSidebars();
  };

  const handleLayoutResize = () => {
    applyResponsiveSidebars();
  };

  return {
    isSidebarCollapsed,
    isSidebarResizing,
    toggleSidebar,
    startResize,
    handleCommentWidthChange,
    clampSidebarWidth
  };
}
