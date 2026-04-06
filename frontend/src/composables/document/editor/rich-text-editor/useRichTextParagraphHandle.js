import { computed, onBeforeUnmount, ref, watch } from 'vue';
import {
  dedupeNumberList,
  resolveDomBlockFromMeta,
  resolveHandleStyle,
  resolveHoveredDomBlock,
  resolveMetaFromRange,
  resolvePositionFromMeta,
  resolveRangeByMouseCoords,
  resolveRangeFromDomBlock
} from './paragraph-handle/paragraphHandleUtils';

export function useRichTextParagraphHandle({
  editorRef,
  scrollContainerRef,
  labels = {},
  resolveActions = null,
  onMutated
}) {
  const visible = ref(false);
  const style = ref({ top: '0px', left: '0px' });
  const hoveredDomBlock = ref(null);
  const hoveredMeta = ref(null);
  const hoveringHandle = ref(false);

  let editorDom = null;
  let mouseHost = null;
  let scrollContainer = null;
  let hideTimer = 0;
  const HIDE_DELAY_MS = 240;
  const TRANSITION_HIDE_DELAY_MS = 320;

  const enabled = computed(() => Boolean(editorRef.value?.view && scrollContainerRef.value));
  const paragraphType = computed(() => hoveredMeta.value?.type || 'paragraph');
  const paragraphIsEmpty = computed(() => Boolean(hoveredMeta.value?.isEmpty));
  const paragraphHandleMode = computed(() => (paragraphIsEmpty.value ? 'empty' : 'normal'));

  const buildDefaultActions = () => ([
    {
      key: 'insert-empty-paragraph',
      label: labels.insertEmptyParagraph || '插入空段落',
      iconKey: 'insert',
      children: [
        {
          key: 'add-above',
          label: labels.addAbove || '在上方',
          iconKey: 'add-above'
        },
        {
          key: 'add-below',
          label: labels.addBelow || '在下方',
          iconKey: 'add-below'
        }
      ]
    },
    {
      key: 'delete',
      label: labels.delete || 'Delete Paragraph',
      iconKey: 'delete',
      danger: true
    }
  ]);

  const paragraphActions = computed(() => {
    const defaults = buildDefaultActions();
    if (typeof resolveActions !== 'function') {
      return defaults;
    }
    const resolved = resolveActions({
      blockType: paragraphType.value,
      isEmpty: paragraphIsEmpty.value,
      defaults: defaults.slice()
    });
    if (!Array.isArray(resolved) || resolved.length === 0) {
      return defaults;
    }
    return resolved;
  });

  const reset = () => {
    visible.value = false;
    hoveredDomBlock.value = null;
    hoveredMeta.value = null;
  };

  const clearHideTimer = () => {
    if (hideTimer) {
      window.clearTimeout(hideTimer);
      hideTimer = 0;
    }
  };

  const scheduleReset = (delay = HIDE_DELAY_MS) => {
    clearHideTimer();
    hideTimer = window.setTimeout(() => {
      hideTimer = 0;
      if (!hoveringHandle.value) {
        reset();
      }
    }, delay);
  };

  const getEditorDom = () => editorRef.value?.view?.dom || null;
  const getEditorView = () => editorRef.value?.view || null;
  const getEditorState = () => getEditorView()?.state || null;

  const updatePosition = (nextMeta = hoveredMeta.value, nextDomBlock = hoveredDomBlock.value) => {
    const container = scrollContainerRef.value;
    if (!container || !nextMeta) {
      reset();
      return;
    }

    let blockRect = null;
    if (nextDomBlock?.isConnected) {
      blockRect = nextDomBlock.getBoundingClientRect();
    }

    if (!blockRect) {
      const coords = resolvePositionFromMeta({
        view: getEditorView(),
        meta: nextMeta
      });
      if (coords) {
        blockRect = {
          top: coords.top,
          left: coords.left
        };
      }
    }

    if (!blockRect) {
      if (!hoveringHandle.value) {
        reset();
      }
      return;
    }

    const nextStyle = resolveHandleStyle({
      container,
      blockRect
    });
    if (!nextStyle) {
      if (!hoveringHandle.value) {
        reset();
      }
      return;
    }

    style.value = nextStyle;
    visible.value = true;
  };

  const updateHoveredBlock = (domBlock, event) => {
    if (!domBlock) {
      reset();
      return;
    }

    const view = getEditorView();
    const state = getEditorState();
    const range = resolveRangeFromDomBlock({ view, state, domBlock })
      || resolveRangeByMouseCoords({ view, state, event });
    const meta = resolveMetaFromRange({ state, range });

    if (!meta || meta.from >= meta.to) {
      reset();
      return;
    }

    const resolvedDomBlock = domBlock || resolveDomBlockFromMeta({
      view,
      editorRoot: getEditorDom(),
      meta
    });

    hoveredDomBlock.value = resolvedDomBlock;
    hoveredMeta.value = meta;
    updatePosition(meta, resolvedDomBlock);
  };

  const handleMouseMove = (event) => {
    clearHideTimer();
    if (event?.target instanceof Element && event.target.closest('.rich-text-paragraph-handle')) {
      return;
    }

    const domBlock = resolveHoveredDomBlock({
      target: event?.target,
      editorDom: getEditorDom()
    });

    if (!domBlock) {
      if (hoveringHandle.value) {
        return;
      }
      if (hoveredMeta.value) {
        // Keep a short grace period when cursor moves from tiny/empty paragraphs to the handle.
        scheduleReset(TRANSITION_HIDE_DELAY_MS);
      }
      return;
    }

    updateHoveredBlock(domBlock, event);
  };

  const handleMouseLeave = (event) => {
    const relatedTarget = event?.relatedTarget;
    if (relatedTarget instanceof Element && relatedTarget.closest('.rich-text-paragraph-handle')) {
      return;
    }
    if (hoveringHandle.value) {
      return;
    }
    scheduleReset(TRANSITION_HIDE_DELAY_MS);
  };

  const handleScrollOrResize = () => {
    if (!hoveredMeta.value) {
      return;
    }
    updatePosition(hoveredMeta.value, hoveredDomBlock.value);
  };

  const bindEvents = () => {
    editorDom = getEditorDom();
    scrollContainer = scrollContainerRef.value;
    mouseHost = scrollContainerRef.value;

    if (!editorDom || !scrollContainer || !mouseHost) {
      return;
    }

    mouseHost.addEventListener('mousemove', handleMouseMove);
    mouseHost.addEventListener('mouseleave', handleMouseLeave);
    if (editorDom !== mouseHost) {
      editorDom.addEventListener('mousemove', handleMouseMove);
      editorDom.addEventListener('mouseleave', handleMouseLeave);
    }
    scrollContainer.addEventListener('scroll', handleScrollOrResize, { passive: true });
    window.addEventListener('resize', handleScrollOrResize);
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
    if (scrollContainer) {
      scrollContainer.removeEventListener('scroll', handleScrollOrResize);
      scrollContainer = null;
    }
    window.removeEventListener('resize', handleScrollOrResize);
    reset();
  };

  const tryInsertParagraphAt = (positions = []) => {
    const editor = editorRef.value;
    const view = editor?.view;
    const state = view?.state;
    const paragraphNodeType = state?.schema?.nodes?.paragraph;
    if (!editor || !view || !state || !paragraphNodeType) {
      return false;
    }

    const candidates = dedupeNumberList(positions)
      .map((pos) => Math.max(0, Math.min(pos, state.doc.content.size)));

    for (const pos of candidates) {
      try {
        const tr = view.state.tr.insert(pos, paragraphNodeType.create());
        view.dispatch(tr.scrollIntoView());
        editor.commands.focus(Math.max(0, Math.min(pos + 1, view.state.doc.content.size)));
        onMutated?.({ type: 'insert', position: pos, blockType: paragraphType.value });
        return true;
      } catch (_) {
        // fallback to next candidate
      }
    }

    return false;
  };

  const addParagraphAbove = () => {
    const meta = hoveredMeta.value;
    if (!meta) {
      return false;
    }
    const inserted = tryInsertParagraphAt([meta.from, meta.topFrom]);
    if (inserted) {
      reset();
    }
    return inserted;
  };

  const addParagraphBelow = () => {
    const meta = hoveredMeta.value;
    if (!meta) {
      return false;
    }
    const inserted = tryInsertParagraphAt([meta.to, meta.topTo]);
    if (inserted) {
      reset();
    }
    return inserted;
  };

  const deleteHoveredParagraph = () => {
    const editor = editorRef.value;
    const view = editor?.view;
    const meta = hoveredMeta.value;
    if (!editor || !view || !meta || meta.from >= meta.to) {
      return false;
    }

    const state = view.state;
    if (state.doc.childCount <= 1) {
      editor.commands.clearContent(true);
      onMutated?.({ type: 'delete', blockType: paragraphType.value });
      reset();
      return true;
    }

    view.dispatch(state.tr.deleteRange(meta.from, meta.to).scrollIntoView());
    editor.commands.focus();
    onMutated?.({ type: 'delete', blockType: paragraphType.value });
    reset();
    return true;
  };

  const focusHoveredBlock = (position = 'start') => {
    const editor = editorRef.value;
    const meta = hoveredMeta.value;
    const size = editor?.state?.doc?.content?.size;
    if (!editor || !meta || !Number.isFinite(size)) {
      return false;
    }
    const targetPos = position === 'end'
      ? Math.max(meta.from + 1, meta.to - 1)
      : meta.from + 1;
    const safePos = Math.max(1, Math.min(targetPos, size));
    editor.commands.focus(safePos);
    return true;
  };

  const runParagraphAction = (key) => {
    const findActionByKey = (actions, targetKey) => {
      if (!Array.isArray(actions) || !targetKey) {
        return null;
      }
      for (const item of actions) {
        if (item?.key === targetKey) {
          return item;
        }
        if (Array.isArray(item?.children) && item.children.length) {
          const childMatched = findActionByKey(item.children, targetKey);
          if (childMatched) {
            return childMatched;
          }
        }
      }
      return null;
    };

    const action = findActionByKey(paragraphActions.value, key);
    if (!action) {
      return;
    }
    const ctx = {
      blockType: paragraphType.value,
      isEmpty: paragraphIsEmpty.value,
      editor: editorRef.value,
      meta: hoveredMeta.value,
      focusBlock: focusHoveredBlock,
      addParagraphAbove,
      addParagraphBelow,
      deleteParagraph: deleteHoveredParagraph
    };
    if (typeof action.handler === 'function') {
      action.handler(ctx);
      return;
    }
    if (key === 'add-above') {
      addParagraphAbove();
      return;
    }
    if (key === 'add-below') {
      addParagraphBelow();
      return;
    }
    if (key === 'delete') {
      deleteHoveredParagraph();
    }
  };

  const handleHandleEnter = () => {
    clearHideTimer();
    hoveringHandle.value = true;
  };

  const handleHandleLeave = () => {
    hoveringHandle.value = false;
    scheduleReset(HIDE_DELAY_MS);
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

  onBeforeUnmount(() => {
    clearHideTimer();
    unbindEvents();
  });

  return {
    paragraphHandleVisible: visible,
    paragraphHandleStyle: style,
    paragraphHandleMode,
    paragraphType,
    paragraphIsEmpty,
    paragraphActions,
    runParagraphAction,
    addParagraphAbove,
    addParagraphBelow,
    deleteHoveredParagraph,
    handleHandleEnter,
    handleHandleLeave
  };
}
