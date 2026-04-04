export const flattenCommentTree = (items) => {
  if (!Array.isArray(items) || items.length === 0) {
    return [];
  }

  const childrenMap = new Map();
  const idSet = new Set();

  items.forEach((item) => {
    if (item?.external_id) {
      idSet.add(item.external_id);
    }
  });

  items.forEach((item) => {
    const parentId = item.parent_external_id || '';
    if (!childrenMap.has(parentId)) {
      childrenMap.set(parentId, []);
    }
    childrenMap.get(parentId).push(item);
  });

  const result = [];
  const walk = (node, level) => {
    node.level = level;
    result.push(node);
    if (!node.external_id) {
      return;
    }
    const children = childrenMap.get(node.external_id) || [];
    children.forEach((child) => walk(child, level + 1));
  };

  const roots = [];
  items.forEach((item) => {
    const parentId = item.parent_external_id || '';
    if (!parentId || !idSet.has(parentId)) {
      roots.push(item);
    }
  });

  roots.forEach((root) => walk(root, 0));
  return result;
};

export const getIndentStyle = (item) => ({
  paddingLeft: `${Math.min(item?.level || 0, 3) * 16}px`
});

export const buildReplyLabel = (item, comments, t) => {
  if (!item?.parent_external_id) {
    return '';
  }
  const parent = comments.find((entry) => entry.external_id === item.parent_external_id);
  return t('document.commentReplyTo', { name: parent?.creator_name || t('document.commentUnknown') });
};
