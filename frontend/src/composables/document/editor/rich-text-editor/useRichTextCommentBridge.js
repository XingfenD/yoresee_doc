import { nextTick, ref } from 'vue';
import { COMMENT_ANCHOR_ATTR } from '@/components/document/rich-text/extensions/commentAnchorExtension';

const createCommentAnchorId = () => `rt_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`;

export function useRichTextCommentBridge(options = {}) {
  const {
    editorRef,
    scrollContainerRef,
    commentEnabledRef,
    onCommentAdd,
    onCommentRemove,
    onCommentChanged
  } = options;

  const selectionCommentVisible = ref(false);
  const selectionCommentStyle = ref({ top: '0px', left: '0px' });
  const selectionCommentHovering = ref(false);
  const selectionUpdateRaf = ref(0);

  const isSelectionInsideEditor = () => {
    const root = editorRef.value?.view?.dom;
    const selection = window.getSelection?.();
    if (!root || !selection || selection.rangeCount === 0 || selection.isCollapsed) {
      return false;
    }
    const range = selection.getRangeAt(0);
    const common = range.commonAncestorContainer;
    const commonElement = common?.nodeType === Node.ELEMENT_NODE ? common : common?.parentElement;
    return Boolean(commonElement && root.contains(commonElement));
  };

  const addInlineComment = () => {
    if (!commentEnabledRef.value || !editorRef.value) {
      return;
    }
    const { from, to, empty } = editorRef.value.state.selection;
    if (empty || from === to) {
      return;
    }
    const selectedText = editorRef.value.state.doc.textBetween(from, to, ' ');
    const anchorId = createCommentAnchorId();
    const chain = editorRef.value.chain().focus();
    if (typeof chain.setCommentAnchor === 'function') {
      chain.setCommentAnchor(anchorId).run();
    } else {
      chain.setMark('commentAnchor', { anchorId }).run();
    }
    onCommentAdd?.({
      id: anchorId,
      text: selectedText || ''
    });
    onCommentChanged?.();
    selectionCommentVisible.value = false;
  };

  const updateSelectionCommentTrigger = () => {
    if (!editorRef.value || !scrollContainerRef.value) {
      selectionCommentVisible.value = false;
      return;
    }
    if (!isSelectionInsideEditor()) {
      selectionCommentVisible.value = false;
      return;
    }
    const instance = editorRef.value;
    const selection = instance.state.selection;
    const { from, to, empty } = selection;
    if (empty || from === to) {
      selectionCommentVisible.value = false;
      return;
    }
    const selectedText = instance.state.doc.textBetween(from, to, ' ').trim();
    if (!selectedText) {
      selectionCommentVisible.value = false;
      return;
    }
    let fromCoords;
    let toCoords;
    try {
      fromCoords = instance.view.coordsAtPos(from);
      toCoords = instance.view.coordsAtPos(to);
    } catch (_) {
      selectionCommentVisible.value = false;
      return;
    }
    const container = scrollContainerRef.value;
    const containerRect = container.getBoundingClientRect();
    const triggerWidth = commentEnabledRef.value ? 308 : 278;
    const triggerHeight = 36;
    const rawLeft = toCoords.right - containerRect.left + container.scrollLeft + 8;
    const rawTop = Math.min(fromCoords.top, toCoords.top) - containerRect.top + container.scrollTop - triggerHeight - 6;
    const minLeft = container.scrollLeft + 8;
    const maxLeft = container.scrollLeft + container.clientWidth - triggerWidth - 8;
    const minTop = container.scrollTop + 8;
    const maxTop = container.scrollTop + container.clientHeight - triggerHeight - 8;
    selectionCommentStyle.value = {
      left: `${Math.min(maxLeft, Math.max(minLeft, rawLeft))}px`,
      top: `${Math.min(maxTop, Math.max(minTop, rawTop))}px`
    };
    selectionCommentVisible.value = true;
  };

  const requestSelectionCommentTriggerUpdate = () => {
    if (selectionUpdateRaf.value) {
      cancelAnimationFrame(selectionUpdateRaf.value);
    }
    selectionUpdateRaf.value = requestAnimationFrame(() => {
      updateSelectionCommentTrigger();
      selectionUpdateRaf.value = 0;
    });
  };

  const getCommentAnchorElements = (ids = null) => {
    const root = editorRef.value?.view?.dom;
    if (!root) {
      return [];
    }
    const all = Array.from(root.querySelectorAll(`span[${COMMENT_ANCHOR_ATTR}]`));
    if (!Array.isArray(ids) || ids.length === 0) {
      return all;
    }
    const idSet = new Set(ids.map((id) => String(id || '').trim()).filter(Boolean));
    return all.filter((element) => idSet.has(String(element.getAttribute(COMMENT_ANCHOR_ATTR) || '').trim()));
  };

  const getCommentIds = () => {
    const elements = getCommentAnchorElements();
    const container = scrollContainerRef.value;
    const containerRect = container?.getBoundingClientRect?.();
    const resultMap = new Map();

    elements.forEach((element) => {
      const anchorId = String(element.getAttribute(COMMENT_ANCHOR_ATTR) || '').trim();
      if (!anchorId) {
        return;
      }
      let top = 0;
      if (container && containerRect) {
        const rect = element.getBoundingClientRect();
        top = rect.top - containerRect.top + container.scrollTop;
      }
      if (!resultMap.has(anchorId) || top < resultMap.get(anchorId)) {
        resultMap.set(anchorId, top);
      }
    });

    return Array.from(resultMap.entries()).map(([id, top]) => ({ id, top }));
  };

  const hlCommentIds = (ids = []) => {
    getCommentAnchorElements(ids).forEach((element) => {
      element.classList.add('comment-anchor-highlight');
    });
  };

  const unHlCommentIds = (ids = []) => {
    getCommentAnchorElements(ids).forEach((element) => {
      element.classList.remove('comment-anchor-highlight');
    });
  };

  const removeCommentIds = (ids = []) => {
    if (!editorRef.value || !Array.isArray(ids) || ids.length === 0) {
      return;
    }
    const idSet = new Set(ids.map((id) => String(id || '').trim()).filter(Boolean));
    if (idSet.size === 0) {
      return;
    }
    let transaction = editorRef.value.state.tr;
    editorRef.value.state.doc.descendants((node, pos) => {
      if (!node.isText || !Array.isArray(node.marks) || node.marks.length === 0) {
        return;
      }
      node.marks.forEach((mark) => {
        if (mark.type.name !== 'commentAnchor') {
          return;
        }
        const anchorId = String(mark.attrs?.anchorId || '').trim();
        if (!idSet.has(anchorId)) {
          return;
        }
        transaction = transaction.removeMark(pos, pos + node.nodeSize, mark);
      });
    });
    if (!transaction.docChanged) {
      return;
    }
    editorRef.value.view.dispatch(transaction);
    onCommentRemove?.(Array.from(idSet));
    onCommentChanged?.();
  };

  const scrollToCommentId = async (id) => {
    const targetId = String(id || '').trim();
    if (!targetId) {
      return;
    }
    await nextTick();
    const container = scrollContainerRef.value;
    if (!container) {
      return;
    }
    const target = getCommentAnchorElements([targetId])[0];
    if (!target) {
      return;
    }
    const containerRect = container.getBoundingClientRect();
    const targetRect = target.getBoundingClientRect();
    const top = Math.max(targetRect.top - containerRect.top + container.scrollTop - 24, 0);
    container.scrollTo({ top, behavior: 'smooth' });
    hlCommentIds([targetId]);
  };

  const clearSelectionRaf = () => {
    if (selectionUpdateRaf.value) {
      cancelAnimationFrame(selectionUpdateRaf.value);
      selectionUpdateRaf.value = 0;
    }
  };

  const handleEditorBlur = () => {
    if (!selectionCommentHovering.value) {
      selectionCommentVisible.value = false;
    }
  };

  return {
    selectionCommentVisible,
    selectionCommentStyle,
    selectionCommentHovering,
    addInlineComment,
    updateSelectionCommentTrigger,
    requestSelectionCommentTriggerUpdate,
    handleEditorBlur,
    clearSelectionRaf,
    getCommentIds,
    hlCommentIds,
    unHlCommentIds,
    removeCommentIds,
    scrollToCommentId
  };
}
