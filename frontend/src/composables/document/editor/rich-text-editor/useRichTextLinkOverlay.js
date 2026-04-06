import { computed, onBeforeUnmount, ref, watch } from 'vue';

const HOVER_HIDE_DELAY = 120;

const normalizeLinkUrl = (value = '') => {
  const input = String(value || '').trim();
  if (!input) {
    return '';
  }
  if (/^(https?:\/\/|mailto:|tel:)/i.test(input)) {
    return input;
  }
  return `https://${input}`;
};

export function useRichTextLinkOverlay(options = {}) {
  const {
    editorRef,
    scrollContainerRef,
    setLink,
    unsetLink
  } = options;

  const linkDialogVisible = ref(false);
  const linkDialogValue = ref('');
  const linkHoverVisible = ref(false);
  const linkHoverStyle = ref({ top: '0px', left: '0px' });
  const linkHoverHref = ref('');
  const linkHoveringCard = ref(false);

  const hoverAnchorPos = ref(null);
  const dialogSource = ref({
    type: 'selection',
    pos: null
  });

  let mouseHost = null;
  let editorDom = null;
  let hideTimer = 0;

  const enabled = computed(() => Boolean(editorRef.value?.view && scrollContainerRef.value));

  const clearHideTimer = () => {
    if (hideTimer) {
      window.clearTimeout(hideTimer);
      hideTimer = 0;
    }
  };

  const hideHoverCard = () => {
    linkHoverVisible.value = false;
    linkHoverHref.value = '';
    hoverAnchorPos.value = null;
  };

  const scheduleHideHoverCard = () => {
    clearHideTimer();
    hideTimer = window.setTimeout(() => {
      hideTimer = 0;
      if (linkHoveringCard.value || linkDialogVisible.value) {
        return;
      }
      hideHoverCard();
    }, HOVER_HIDE_DELAY);
  };

  const getEditorDom = () => editorRef.value?.view?.dom || null;

  const resolveAnchorPos = (anchorEl) => {
    const view = editorRef.value?.view;
    const size = editorRef.value?.state?.doc?.content?.size;
    if (!view || !anchorEl || !Number.isFinite(size)) {
      return null;
    }
    try {
      const pos = view.posAtDOM(anchorEl, 0);
      const safePos = Math.max(1, Math.min(pos + 1, size));
      return safePos;
    } catch (_) {
      return null;
    }
  };

  const updateHoverCardFromAnchor = (anchorEl) => {
    const container = scrollContainerRef.value;
    if (!anchorEl || !container) {
      hideHoverCard();
      return;
    }
    const href = String(anchorEl.getAttribute('href') || '').trim();
    if (!href) {
      hideHoverCard();
      return;
    }
    const rect = anchorEl.getBoundingClientRect();
    const containerRect = container.getBoundingClientRect();

    const cardWidth = 280;
    const cardHeight = 108;
    const rawLeft = rect.left - containerRect.left + container.scrollLeft;
    const rawTop = rect.bottom - containerRect.top + container.scrollTop + 8;
    const minLeft = container.scrollLeft + 8;
    const maxLeft = container.scrollLeft + container.clientWidth - cardWidth - 8;
    const minTop = container.scrollTop + 8;
    const maxTop = container.scrollTop + container.clientHeight - cardHeight - 8;

    linkHoverStyle.value = {
      left: `${Math.max(minLeft, Math.min(rawLeft, Math.max(minLeft, maxLeft)))}px`,
      top: `${Math.max(minTop, Math.min(rawTop, Math.max(minTop, maxTop)))}px`
    };
    linkHoverHref.value = href;
    hoverAnchorPos.value = resolveAnchorPos(anchorEl);
    linkHoverVisible.value = true;
  };

  const handleMouseMove = (event) => {
    if (linkDialogVisible.value) {
      return;
    }
    clearHideTimer();
    const root = getEditorDom();
    const target = event?.target;
    const source = target?.nodeType === Node.TEXT_NODE ? target.parentElement : target;
    if (!(source instanceof Element) || !root?.contains(source)) {
      if (!linkHoveringCard.value) {
        scheduleHideHoverCard();
      }
      return;
    }
    const anchor = source.closest('a[href]');
    if (!anchor || !root.contains(anchor)) {
      if (!linkHoveringCard.value) {
        scheduleHideHoverCard();
      }
      return;
    }
    updateHoverCardFromAnchor(anchor);
  };

  const handleMouseLeave = () => {
    if (linkHoveringCard.value || linkDialogVisible.value) {
      return;
    }
    scheduleHideHoverCard();
  };

  const bindEvents = () => {
    editorDom = getEditorDom();
    mouseHost = scrollContainerRef.value;
    if (!editorDom || !mouseHost) {
      return;
    }
    mouseHost.addEventListener('mousemove', handleMouseMove);
    mouseHost.addEventListener('mouseleave', handleMouseLeave);
    if (editorDom !== mouseHost) {
      editorDom.addEventListener('mousemove', handleMouseMove);
      editorDom.addEventListener('mouseleave', handleMouseLeave);
    }
  };

  const unbindEvents = () => {
    clearHideTimer();
    if (mouseHost) {
      mouseHost.removeEventListener('mousemove', handleMouseMove);
      mouseHost.removeEventListener('mouseleave', handleMouseLeave);
      mouseHost = null;
    }
    if (editorDom) {
      editorDom.removeEventListener('mousemove', handleMouseMove);
      editorDom.removeEventListener('mouseleave', handleMouseLeave);
      editorDom = null;
    }
    hideHoverCard();
  };

  const openSelectionLinkDialog = () => {
    const currentHref = String(editorRef.value?.getAttributes?.('link')?.href || '').trim();
    linkDialogValue.value = currentHref;
    dialogSource.value = { type: 'selection', pos: null };
    linkDialogVisible.value = true;
  };

  const openHoverLinkDialog = () => {
    linkDialogValue.value = linkHoverHref.value;
    dialogSource.value = { type: 'hover', pos: hoverAnchorPos.value };
    linkDialogVisible.value = true;
  };

  const closeLinkDialog = () => {
    linkDialogVisible.value = false;
  };

  const applyLinkDialog = () => {
    const source = dialogSource.value;
    const normalized = normalizeLinkUrl(linkDialogValue.value);
    if (!normalized) {
      unsetLink?.({ pos: source.pos });
      linkDialogVisible.value = false;
      if (source.type === 'hover') {
        hideHoverCard();
      }
      return;
    }
    setLink?.(normalized, { pos: source.pos });
    linkDialogVisible.value = false;
    if (source.type === 'hover') {
      linkHoverHref.value = normalized;
    }
  };

  const removeHoverLink = () => {
    unsetLink?.({ pos: hoverAnchorPos.value });
    hideHoverCard();
  };

  const handleLinkCardEnter = () => {
    clearHideTimer();
    linkHoveringCard.value = true;
  };

  const handleLinkCardLeave = () => {
    linkHoveringCard.value = false;
    if (!linkDialogVisible.value) {
      scheduleHideHoverCard();
    }
  };

  watch(
    enabled,
    (next) => {
      unbindEvents();
      if (next) {
        bindEvents();
      }
    },
    { immediate: true }
  );

  watch(
    () => editorRef.value,
    () => {
      if (!enabled.value) {
        return;
      }
      unbindEvents();
      bindEvents();
    }
  );

  watch(linkDialogVisible, (next) => {
    if (!next && !linkHoveringCard.value) {
      scheduleHideHoverCard();
    }
  });

  onBeforeUnmount(() => {
    unbindEvents();
  });

  return {
    linkDialogVisible,
    linkDialogValue,
    linkHoverVisible,
    linkHoverStyle,
    linkHoverHref,
    openSelectionLinkDialog,
    openHoverLinkDialog,
    closeLinkDialog,
    applyLinkDialog,
    removeHoverLink,
    handleLinkCardEnter,
    handleLinkCardLeave
  };
}
