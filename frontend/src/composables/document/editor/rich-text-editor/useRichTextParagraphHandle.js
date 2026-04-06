import { computed, onBeforeUnmount, ref, watch } from 'vue';

const BLOCK_SELECTOR = [
  'p',
  'h1',
  'h2',
  'h3',
  'h4',
  'h5',
  'h6',
  'blockquote',
  'pre',
  'li',
  'table',
  'hr',
  '.mindmap-node',
  '.drawio-node',
  'yoresee-mindmap',
  'yoresee-drawio'
].join(',');

const PREFERRED_BLOCK_NODE_NAMES = new Set([
  'listItem',
  'paragraph',
  'heading',
  'blockquote',
  'codeBlock',
  'horizontalRule',
  'table',
  'mindmapBlock',
  'drawioBlock'
]);

const NON_EMPTY_BLOCK_TYPES = new Set([
  'horizontalRule',
  'table',
  'mindmapBlock',
  'drawioBlock'
]);
const HANDLE_SIZE = 22;

const normalizeBlockType = (value) => {
  switch (value) {
    case 'heading':
      return 'heading';
    case 'listItem':
      return 'list';
    case 'blockquote':
      return 'quote';
    case 'codeBlock':
      return 'code';
    case 'table':
      return 'table';
    case 'horizontalRule':
      return 'divider';
    case 'mindmapBlock':
      return 'mindmap';
    case 'drawioBlock':
      return 'drawio';
    default:
      return 'paragraph';
  }
};

const dedupeNumberList = (list = []) => {
  const seen = new Set();
  const result = [];
  list.forEach((item) => {
    const value = Number(item);
    if (!Number.isFinite(value) || seen.has(value)) {
      return;
    }
    seen.add(value);
    result.push(value);
  });
  return result;
};

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

  const enabled = computed(() => Boolean(editorRef.value?.view && scrollContainerRef.value));
  const paragraphType = computed(() => hoveredMeta.value?.type || 'paragraph');
  const paragraphIsEmpty = computed(() => Boolean(hoveredMeta.value?.isEmpty));
  const paragraphHandleMode = computed(() => (paragraphIsEmpty.value ? 'empty' : 'normal'));

  const buildDefaultActions = () => ([
    {
      key: 'add-above',
      label: labels.addAbove || 'Add Paragraph Above',
      iconKey: 'add-above'
    },
    {
      key: 'add-below',
      label: labels.addBelow || 'Add Paragraph Below',
      iconKey: 'add-below'
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

  const getEditorDom = () => editorRef.value?.view?.dom || null;

  const resolveHoveredDomBlock = (target) => {
    const currentEditorDom = getEditorDom();
    if (!currentEditorDom) {
      return null;
    }
    const sourceTarget = target?.nodeType === Node.TEXT_NODE ? target.parentElement : target;
    if (!(sourceTarget instanceof Element) || !currentEditorDom.contains(sourceTarget)) {
      return null;
    }

    let domBlock = sourceTarget.closest(BLOCK_SELECTOR);
    if (!domBlock || !currentEditorDom.contains(domBlock)) {
      return null;
    }

    if (domBlock.tagName === 'P' && domBlock.parentElement?.tagName === 'LI') {
      domBlock = domBlock.parentElement;
    }

    return domBlock;
  };

  const resolveRangeFromDomBlock = (domBlock) => {
    const view = editorRef.value?.view;
    const state = view?.state;
    if (!view || !state || !domBlock) {
      return null;
    }
    try {
      const pos = view.posAtDOM(domBlock, 0);
      const safePos = Math.max(0, Math.min(pos, state.doc.content.size));
      const directNode = state.doc.nodeAt(safePos);
      if (directNode && (directNode.isBlock || PREFERRED_BLOCK_NODE_NAMES.has(directNode.type?.name))) {
        return {
          from: safePos,
          to: safePos + directNode.nodeSize
        };
      }
      return {
        from: safePos,
        to: Math.max(0, Math.min(safePos + domBlock.textContent.length + 2, state.doc.content.size))
      };
    } catch (_) {
      return null;
    }
  };

  const resolveRangeByMouseCoords = (event) => {
    const view = editorRef.value?.view;
    const state = view?.state;
    if (!view || !state || typeof event?.clientX !== 'number' || typeof event?.clientY !== 'number') {
      return null;
    }

    const coords = view.posAtCoords({
      left: event.clientX,
      top: event.clientY
    });
    if (!coords || typeof coords.pos !== 'number') {
      return null;
    }

    const pos = Math.max(0, Math.min(coords.pos, state.doc.content.size));
    return { from: pos, to: pos };
  };

  const resolveMetaFromRange = (range) => {
    const view = editorRef.value?.view;
    const state = view?.state;
    if (!view || !state || !range) {
      return null;
    }

    const directNode = state.doc.nodeAt(range.from);
    if (directNode && (directNode.isBlock || PREFERRED_BLOCK_NODE_NAMES.has(directNode.type?.name))) {
      const nodeType = String(directNode.type?.name || 'paragraph');
      const normalizedType = normalizeBlockType(nodeType);
      const text = String(directNode.textContent || '').replace(/\u200b/g, '').trim();
      const isEmpty = !NON_EMPTY_BLOCK_TYPES.has(nodeType) && text.length === 0;
      return {
        from: range.from,
        to: range.from + directNode.nodeSize,
        topFrom: range.from,
        topTo: range.from + directNode.nodeSize,
        nodeType,
        type: normalizedType,
        isEmpty
      };
    }

    const probePos = Math.max(0, Math.min(range.from + 1, state.doc.content.size));
    const $probe = state.doc.resolve(probePos);

    for (let depth = $probe.depth; depth > 0; depth -= 1) {
      const node = $probe.node(depth);
      if (!node || (!node.isBlock && !PREFERRED_BLOCK_NODE_NAMES.has(node.type?.name))) {
        continue;
      }
      const from = $probe.before(depth);
      const to = $probe.after(depth);
      const topFrom = $probe.before(1);
      const topTo = $probe.after(1);
      const nodeType = String(node.type?.name || 'paragraph');
      const normalizedType = normalizeBlockType(nodeType);
      const text = String(node.textContent || '').replace(/\u200b/g, '').trim();
      const isEmpty = !NON_EMPTY_BLOCK_TYPES.has(nodeType) && text.length === 0;
      return {
        from,
        to,
        topFrom,
        topTo,
        nodeType,
        type: normalizedType,
        isEmpty
      };
    }

    return null;
  };

  const updatePosition = () => {
    const domBlock = hoveredDomBlock.value;
    const container = scrollContainerRef.value;
    if (!domBlock || !container || !domBlock.isConnected) {
      reset();
      return;
    }
    const blockRect = domBlock.getBoundingClientRect();
    const containerRect = container.getBoundingClientRect();
    const top = blockRect.top - containerRect.top + container.scrollTop - HANDLE_SIZE;
    const left = blockRect.left - containerRect.left + container.scrollLeft - HANDLE_SIZE;

    style.value = {
      top: `${top}px`,
      left: `${left}px`
    };
    visible.value = true;
  };

  const updateHoveredBlock = (domBlock, event) => {
    if (!domBlock) {
      reset();
      return;
    }

    const range = resolveRangeFromDomBlock(domBlock) || resolveRangeByMouseCoords(event);
    const meta = resolveMetaFromRange(range);
    if (!meta || meta.from >= meta.to) {
      reset();
      return;
    }

    hoveredDomBlock.value = domBlock;
    hoveredMeta.value = meta;
    updatePosition();
  };

  const handleMouseMove = (event) => {
    if (event?.target instanceof Element && event.target.closest('.rich-text-paragraph-handle')) {
      return;
    }
    updateHoveredBlock(resolveHoveredDomBlock(event?.target), event);
  };

  const handleMouseLeave = (event) => {
    const relatedTarget = event?.relatedTarget;
    if (relatedTarget instanceof Element && relatedTarget.closest('.rich-text-paragraph-handle')) {
      return;
    }
    if (hoveringHandle.value) {
      return;
    }
    reset();
  };

  const handleScrollOrResize = () => {
    if (!hoveredDomBlock.value) {
      return;
    }
    updatePosition();
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

    const candidates = dedupeNumberList(positions).map((pos) => Math.max(0, Math.min(pos, state.doc.content.size)));
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

  const runParagraphAction = (key) => {
    const action = paragraphActions.value.find((item) => item?.key === key);
    if (!action) {
      return;
    }
    const ctx = {
      blockType: paragraphType.value,
      isEmpty: paragraphIsEmpty.value,
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
    hoveringHandle.value = true;
  };

  const handleHandleLeave = () => {
    hoveringHandle.value = false;
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
