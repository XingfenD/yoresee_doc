export const BLOCK_SELECTOR = [
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
  '.rich-table-node',
  'yoresee-mindmap',
  'yoresee-drawio',
  'yoresee-table'
].join(',');

const PREFERRED_BLOCK_NODE_NAMES = new Set([
  'listItem',
  'paragraph',
  'heading',
  'blockquote',
  'codeBlock',
  'horizontalRule',
  'table',
  'tableBlock',
  'mindmapBlock',
  'drawioBlock'
]);

const NON_EMPTY_BLOCK_TYPES = new Set([
  'horizontalRule',
  'table',
  'tableBlock',
  'mindmapBlock',
  'drawioBlock'
]);

export const HANDLE_SIZE = 22;

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
    case 'tableBlock':
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

export const dedupeNumberList = (list = []) => {
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

export const resolveHoveredDomBlock = ({ target, editorDom }) => {
  if (!editorDom) {
    return null;
  }
  const sourceTarget = target?.nodeType === Node.TEXT_NODE ? target.parentElement : target;
  if (!(sourceTarget instanceof Element) || !editorDom.contains(sourceTarget)) {
    return null;
  }

  let domBlock = sourceTarget.closest(BLOCK_SELECTOR);
  if (!domBlock || !editorDom.contains(domBlock)) {
    return null;
  }

  if (domBlock.tagName === 'P' && domBlock.parentElement?.tagName === 'LI') {
    domBlock = domBlock.parentElement;
  }

  return domBlock;
};

export const resolveRangeFromDomBlock = ({ view, state, domBlock }) => {
  if (!view || !state || !domBlock) {
    return null;
  }
  try {
    const rawPos = view.posAtDOM(domBlock, 0);
    const size = state.doc.content.size;
    const candidates = dedupeNumberList([
      rawPos - 1,
      rawPos,
      rawPos + 1
    ]).map((pos) => Math.max(0, Math.min(pos, size)));

    for (const safePos of candidates) {
      const directNode = state.doc.nodeAt(safePos);
      if (directNode && (directNode.isBlock || PREFERRED_BLOCK_NODE_NAMES.has(directNode.type?.name))) {
        return {
          from: safePos,
          to: safePos + directNode.nodeSize
        };
      }
    }

    const safePos = Math.max(0, Math.min(rawPos, size));
    return {
      from: safePos,
      to: Math.max(0, Math.min(safePos + domBlock.textContent.length + 2, state.doc.content.size))
    };
  } catch (_) {
    return null;
  }
};

export const resolveRangeByMouseCoords = ({ view, state, event }) => {
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

export const resolveMetaFromRange = ({ state, range }) => {
  if (!state || !range) {
    return null;
  }

  const resolveDirectMetaByPos = (pos) => {
    const directNode = state.doc.nodeAt(pos);
    if (!directNode || (!directNode.isBlock && !PREFERRED_BLOCK_NODE_NAMES.has(directNode.type?.name))) {
      return null;
    }
    const nodeType = String(directNode.type?.name || 'paragraph');
    const normalizedType = normalizeBlockType(nodeType);
    const text = String(directNode.textContent || '').replace(/\u200b/g, '').trim();
    const isEmpty = !NON_EMPTY_BLOCK_TYPES.has(nodeType)
      && text.length === 0
      && Number(directNode.content?.size || 0) === 0;
    return {
      from: pos,
      to: pos + directNode.nodeSize,
      topFrom: pos,
      topTo: pos + directNode.nodeSize,
      nodeType,
      type: normalizedType,
      isEmpty
    };
  };

  const directCandidates = dedupeNumberList([
    range.from,
    range.from - 1,
    range.from + 1
  ]).map((pos) => Math.max(0, Math.min(pos, state.doc.content.size)));

  for (const pos of directCandidates) {
    const directMeta = resolveDirectMetaByPos(pos);
    if (directMeta) {
      return directMeta;
    }
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
    const isEmpty = !NON_EMPTY_BLOCK_TYPES.has(nodeType)
      && text.length === 0
      && Number(node.content?.size || 0) === 0;
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

export const resolvePositionFromMeta = ({ view, meta }) => {
  if (!view || !meta) {
    return null;
  }
  const upperBound = Math.max(1, view.state.doc.content.size);
  const candidatePositions = dedupeNumberList([
    meta.from + 1,
    meta.from,
    Math.max(meta.to - 1, 1)
  ]).map((pos) => Math.max(1, Math.min(pos, upperBound)));

  for (const pos of candidatePositions) {
    try {
      const coords = view.coordsAtPos(pos);
      if (coords?.left != null && coords?.top != null) {
        return coords;
      }
    } catch (_) {
      // continue fallback candidates
    }
  }
  return null;
};

export const resolveDomBlockFromMeta = ({ view, editorRoot, meta }) => {
  if (!view || !editorRoot || !meta) {
    return null;
  }
  try {
    const dom = view.nodeDOM(meta.from);
    if (!dom) {
      return null;
    }
    const base = dom.nodeType === Node.TEXT_NODE ? dom.parentElement : dom;
    if (!(base instanceof Element)) {
      return null;
    }
    if (base.matches(BLOCK_SELECTOR)) {
      return base;
    }
    const closest = base.closest(BLOCK_SELECTOR);
    if (closest && editorRoot.contains(closest)) {
      return closest;
    }
    return null;
  } catch (_) {
    return null;
  }
};

export const resolveHandleStyle = ({ container, blockRect, handleSize = HANDLE_SIZE }) => {
  if (!container || !blockRect) {
    return null;
  }
  const containerRect = container.getBoundingClientRect();
  const rawTop = blockRect.top - containerRect.top + container.scrollTop - handleSize;
  const rawLeft = blockRect.left - containerRect.left + container.scrollLeft - handleSize;
  const minTop = container.scrollTop + 2;
  const minLeft = container.scrollLeft + 2;
  const maxTop = container.scrollTop + container.clientHeight - handleSize - 2;
  const maxLeft = container.scrollLeft + container.clientWidth - handleSize - 2;
  const top = Math.max(minTop, Math.min(rawTop, Math.max(minTop, maxTop)));
  const left = Math.max(minLeft, Math.min(rawLeft, Math.max(minLeft, maxLeft)));
  return {
    top: `${top}px`,
    left: `${left}px`
  };
};
