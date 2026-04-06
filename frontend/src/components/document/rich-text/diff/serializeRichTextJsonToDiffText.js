import { resolveRichTextPreviewDiffAdapterRegistry } from '@/components/document/rich-text/components/registry';

const defaultAdapterRegistry = resolveRichTextPreviewDiffAdapterRegistry();

const collectText = (node) => {
  if (!node || typeof node !== 'object') {
    return '';
  }
  if (node.type === 'text') {
    return String(node.text || '');
  }
  if (!Array.isArray(node.content)) {
    return '';
  }
  return node.content.map((child) => collectText(child)).join('');
};

const serializeNode = (node, output, mode, adapterRegistry, orderedIndexRef = { value: 1 }) => {
  if (!node || typeof node !== 'object') {
    return;
  }

  const adapter = adapterRegistry[node.type];
  const adapterSerializer = mode === 'preview' ? adapter?.toPreview : adapter?.toDiff;
  if (typeof adapterSerializer === 'function') {
    const result = adapterSerializer(node);
    if (result) {
      output.push(result);
    }
    return;
  }

  if (node.type === 'heading') {
    const level = Number(node?.attrs?.level || 1);
    output.push(`${'#'.repeat(Math.max(1, Math.min(6, level)))} ${collectText(node)}`.trim());
    return;
  }

  if (node.type === 'paragraph') {
    const text = collectText(node).trim();
    if (text) {
      output.push(text);
    }
    return;
  }

  if (node.type === 'blockquote') {
    const text = collectText(node).trim();
    if (text) {
      output.push(`> ${text}`);
    }
    return;
  }

  if (node.type === 'codeBlock') {
    const lang = String(node?.attrs?.language || '').trim();
    const text = collectText(node);
    output.push(`\`\`\`${lang}\n${text}\n\`\`\``.trim());
    return;
  }

  if (node.type === 'bulletList' && Array.isArray(node.content)) {
    node.content.forEach((item) => {
      const text = collectText(item).trim();
      if (text) {
        output.push(`- ${text}`);
      }
    });
    return;
  }

  if (node.type === 'orderedList' && Array.isArray(node.content)) {
    const start = Number(node?.attrs?.start || 1);
    orderedIndexRef.value = Number.isFinite(start) ? start : 1;
    node.content.forEach((item) => {
      const text = collectText(item).trim();
      if (text) {
        output.push(`${orderedIndexRef.value}. ${text}`);
        orderedIndexRef.value += 1;
      }
    });
    return;
  }

  if (Array.isArray(node.content)) {
    node.content.forEach((child) => serializeNode(child, output, mode, adapterRegistry, orderedIndexRef));
  }
};

const serializeRichTextJson = (doc, mode = 'diff', adapterRegistry = defaultAdapterRegistry) => {
  const output = [];
  serializeNode(doc, output, mode, adapterRegistry);
  return output.join('\n\n');
};

export const serializeRichTextJsonToDiffText = (
  doc,
  adapterRegistry = defaultAdapterRegistry
) => serializeRichTextJson(doc, 'diff', adapterRegistry);

export const serializeRichTextJsonToPreviewText = (
  doc,
  adapterRegistry = defaultAdapterRegistry
) => serializeRichTextJson(doc, 'preview', adapterRegistry);
