const formatMindmap = (node) => {
  const source = String(node?.attrs?.source || '').trim();
  if (!source) {
    return '[Mindmap]';
  }
  const lines = source.split('\n').slice(0, 8).join('\n');
  return `[Mindmap]\n${lines}`;
};

export const mindmapPreviewDiffAdapter = {
  toPreview: formatMindmap,
  toDiff: formatMindmap
};
