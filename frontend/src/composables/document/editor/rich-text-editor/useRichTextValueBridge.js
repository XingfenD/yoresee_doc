import { computed, ref } from 'vue';
import TurndownService from 'turndown';
import { marked } from 'marked';
import { COMMENT_ANCHOR_ATTR } from '@/components/document/rich-text/extensions/commentAnchorExtension';

const EMPTY_DOC = Object.freeze({ type: 'doc', content: [] });
const EMPTY_DOC_JSON = '{"type":"doc","content":[]}';

const decodeMindmapSource = (value) => {
  if (!value) {
    return '';
  }
  try {
    return decodeURIComponent(String(value));
  } catch (_) {
    return String(value);
  }
};

const encodeDrawioSourceToBase64 = (value) => {
  const source = String(value || '');
  if (!source) {
    return '';
  }
  try {
    if (typeof TextEncoder !== 'undefined') {
      const bytes = new TextEncoder().encode(source);
      let binary = '';
      const chunkSize = 0x8000;
      for (let index = 0; index < bytes.length; index += chunkSize) {
        const chunk = bytes.subarray(index, index + chunkSize);
        binary += String.fromCharCode(...chunk);
      }
      return btoa(binary);
    }
    return btoa(unescape(encodeURIComponent(source)));
  } catch (_) {
    return '';
  }
};

const decodeDrawioSourceFromBase64 = (value) => {
  const source = String(value || '').trim();
  if (!source) {
    return '';
  }
  try {
    const binary = atob(source);
    if (typeof TextDecoder !== 'undefined') {
      const bytes = Uint8Array.from(binary, (char) => char.charCodeAt(0));
      return new TextDecoder().decode(bytes);
    }
    return decodeURIComponent(escape(binary));
  } catch (_) {
    return '';
  }
};

const normalizeDrawioSource = (value) => {
  const source = String(value || '').trim();
  if (!source) {
    return '';
  }
  if (source.startsWith('base64:')) {
    const decoded = decodeDrawioSourceFromBase64(source.slice('base64:'.length));
    return decoded || '';
  }
  return source;
};

const resolveMindmapSourceFromNode = (node) => {
  if (!node || typeof node.getAttribute !== 'function') {
    return '';
  }
  const rawSource = node.getAttribute('data-source') || node.getAttribute('source') || '';
  if (rawSource) {
    return decodeMindmapSource(rawSource);
  }
  if (typeof node.querySelector === 'function') {
    const textarea = node.querySelector('textarea');
    if (textarea?.value) {
      return String(textarea.value);
    }
  }
  return '';
};

const resolveDrawioSourceFromNode = (node) => {
  if (!node || typeof node.getAttribute !== 'function') {
    return '';
  }
  const source = node.getAttribute('data-diagram') || node.getAttribute('diagram') || '';
  if (!source) {
    return '';
  }
  try {
    return decodeURIComponent(String(source));
  } catch (_) {
    return String(source);
  }
};

const extractDrawioBlocksFromHtml = (sourceHtml) => {
  const source = String(sourceHtml || '');
  if (!source) {
    return [];
  }

  const blocks = [];
  const pattern = /<([a-z0-9-]+)\b([^>]*)>/gi;
  let matched = pattern.exec(source);
  while (matched) {
    const tagName = String(matched[1] || '').toLowerCase();
    const attrsSource = String(matched[2] || '');
    const dataTypeMatch = attrsSource.match(/\bdata-type=["']([^"']+)["']/i);
    const dataType = String(dataTypeMatch?.[1] || '').toLowerCase();
    const hasDrawioType = tagName === 'yoresee-drawio' || dataType.includes('drawio');
    if (!hasDrawioType) {
      matched = pattern.exec(source);
      continue;
    }
    const dataDiagramMatch = attrsSource.match(/\b(?:data-diagram|diagram)=["']([^"']+)["']/i);
    if (!dataDiagramMatch?.[1]) {
      matched = pattern.exec(source);
      continue;
    }

    let decodedSource = '';
    try {
      decodedSource = decodeURIComponent(String(dataDiagramMatch[1]));
    } catch (_) {
      decodedSource = String(dataDiagramMatch[1]);
    }
    decodedSource = decodedSource.trim();
    if (!decodedSource) {
      matched = pattern.exec(source);
      continue;
    }

    const encoded = encodeDrawioSourceToBase64(decodedSource);
    if (encoded) {
      blocks.push(`\`\`\`drawio\nbase64:${encoded}\n\`\`\``);
    }

    matched = pattern.exec(source);
  }
  return blocks;
};

const normalizeLanguageToken = (value = '') =>
  String(value || '')
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9_+-]/g, '');

const extractLanguageFromCodeClass = (className = '') => {
  const source = String(className || '').trim();
  if (!source) {
    return '';
  }
  const matched = source.match(/(?:^|\s)language-([a-z0-9_+-]+)(?:\s|$)/i);
  if (!matched?.[1]) {
    return '';
  }
  return normalizeLanguageToken(matched[1]);
};

const createRichTextTurndown = () => {
  const turndown = new TurndownService({
    headingStyle: 'atx',
    bulletListMarker: '-',
    codeBlockStyle: 'fenced',
    emDelimiter: '*'
  });

  turndown.addRule('fencedCodeBlockWithLanguage', {
    filter: (node) => {
      if (!node || node.nodeName !== 'PRE' || !node.firstChild) {
        return false;
      }
      return node.firstChild.nodeName === 'CODE';
    },
    replacement: (_, node) => {
      const codeNode = node.firstChild;
      const className = typeof codeNode.getAttribute === 'function' ? codeNode.getAttribute('class') || '' : '';
      const language = extractLanguageFromCodeClass(className);
      const rawText = String(codeNode.textContent || '').replace(/\n+$/, '');
      const languageSuffix = language ? language : '';
      return `\n\n\`\`\`${languageSuffix}\n${rawText}\n\`\`\`\n\n`;
    }
  });

  turndown.addRule('yoreseeMindmap', {
    filter: (node) => {
      if (!node) {
        return false;
      }
      const nodeName = String(node.nodeName || '').toLowerCase();
      if (nodeName === 'yoresee-mindmap') {
        return true;
      }
      if (typeof node.getAttribute !== 'function') {
        return false;
      }
      const dataType = String(node.getAttribute('data-type') || '').toLowerCase();
      const hasDataSource = Boolean(node.getAttribute('data-source') || node.getAttribute('source'));
      return hasDataSource && dataType.includes('mindmap');
    },
    replacement: (_, node) => {
      const source = resolveMindmapSourceFromNode(node).trim();
      if (!source) {
        return '';
      }
      return `\n\n\`\`\`mindmap\n${source}\n\`\`\`\n\n`;
    }
  });

  turndown.addRule('yoreseeDrawio', {
    filter: (node) => {
      if (!node) {
        return false;
      }
      const nodeName = String(node.nodeName || '').toLowerCase();
      if (nodeName === 'yoresee-drawio') {
        return true;
      }
      if (typeof node.getAttribute !== 'function') {
        return false;
      }
      const dataType = String(node.getAttribute('data-type') || '').toLowerCase();
      const hasData = Boolean(node.getAttribute('data-diagram') || node.getAttribute('diagram'));
      return hasData && dataType.includes('drawio');
    },
    replacement: (_, node) => {
      const source = resolveDrawioSourceFromNode(node).trim();
      if (!source) {
        return '';
      }
      const encoded = encodeDrawioSourceToBase64(source);
      if (!encoded) {
        return '';
      }
      return `\n\n\`\`\`drawio\nbase64:${encoded}\n\`\`\`\n\n`;
    }
  });

  turndown.addRule('yoreseeCommentAnchor', {
    filter: (node) => node.nodeName === 'SPAN' && node.getAttribute(COMMENT_ANCHOR_ATTR),
    replacement: (content, node) => {
      const anchorId = String(node.getAttribute(COMMENT_ANCHOR_ATTR) || '').trim();
      if (!anchorId) {
        return content || '';
      }
      return `<span ${COMMENT_ANCHOR_ATTR}="${anchorId}">${content || ''}</span>`;
    }
  });

  return turndown;
};

marked.setOptions({
  gfm: true,
  breaks: true
});

const markdownToHtml = (value) => {
  const markdownSource = String(value || '');
  if (!markdownSource.trim()) {
    return '<p></p>';
  }

  const sourceWithComponentBlocks = markdownSource
    .replace(/```drawio\s*\n([\s\S]*?)```/gi, (_, drawioSource) => {
      const normalized = normalizeDrawioSource(drawioSource);
      if (!normalized) {
        return '';
      }
      const encoded = encodeURIComponent(normalized);
      return `\n<yoresee-drawio data-diagram="${encoded}" data-type="drawio"></yoresee-drawio>\n`;
    })
    .replace(/```mindmap\s*\n([\s\S]*?)```/gi, (_, mindmapSource) => {
      const encoded = encodeURIComponent(String(mindmapSource || '').trim());
      return `\n<yoresee-mindmap data-source="${encoded}"></yoresee-mindmap>\n`;
    });

  const parsed = marked.parse(sourceWithComponentBlocks, { async: false });
  return typeof parsed === 'string' ? parsed : markdownSource;
};

const htmlToMarkdown = (value, turndown) => {
  const source = String(value || '');
  if (!source.trim()) {
    return '';
  }
  const drawioFallbackBlocks = extractDrawioBlocksFromHtml(source);
  let markdown = turndown.turndown(source);
  if (drawioFallbackBlocks.length > 0 && !/```drawio[\s\S]*?```/i.test(markdown)) {
    markdown = `${markdown.trimEnd()}\n\n${drawioFallbackBlocks.join('\n\n')}`.trim();
  }
  return markdown.replace(/\n{3,}/g, '\n\n').trimEnd();
};

const safeClone = (value) => {
  try {
    return JSON.parse(JSON.stringify(value));
  } catch (_) {
    return { ...EMPTY_DOC };
  }
};

const normalizeJsonDoc = (value) => {
  if (value && typeof value === 'object' && value.type === 'doc') {
    return safeClone(value);
  }
  return { ...EMPTY_DOC };
};

const serializeJsonDoc = (value) => {
  try {
    return JSON.stringify(normalizeJsonDoc(value));
  } catch (_) {
    return EMPTY_DOC_JSON;
  }
};

export function useRichTextValueBridge(options = {}) {
  const { editorRef, valueFormatRef } = options;
  const applyingModelValue = ref(false);
  const lastEmittedValue = ref('');
  const isJsonValueMode = computed(() => valueFormatRef.value === 'json');
  const turndown = createRichTextTurndown();

  const serializeModelValue = (value) => {
    if (isJsonValueMode.value) {
      return serializeJsonDoc(value);
    }
    return String(value || '');
  };

  const modelValueFromEditor = (instance) => {
    if (isJsonValueMode.value) {
      return safeClone(instance.getJSON());
    }
    return htmlToMarkdown(instance.getHTML(), turndown);
  };

  const resolveInitialEditorContent = (modelValue) => {
    if (isJsonValueMode.value) {
      return normalizeJsonDoc(modelValue);
    }
    return markdownToHtml(String(modelValue || ''));
  };

  const applyModelValueToEditor = (modelValue) => {
    if (!editorRef.value) {
      return;
    }
    const serialized = serializeModelValue(modelValue);
    if (serialized === lastEmittedValue.value) {
      return;
    }
    applyingModelValue.value = true;
    if (isJsonValueMode.value) {
      editorRef.value.commands.setContent(normalizeJsonDoc(modelValue), false);
    } else {
      editorRef.value.commands.setContent(markdownToHtml(String(modelValue || '')), false, {
        preserveWhitespace: true
      });
    }
    lastEmittedValue.value = serialized;
    applyingModelValue.value = false;
  };

  const syncLastEmittedValue = (value) => {
    lastEmittedValue.value = serializeModelValue(value);
  };

  return {
    applyingModelValue,
    isJsonValueMode,
    lastEmittedValue,
    serializeModelValue,
    modelValueFromEditor,
    resolveInitialEditorContent,
    applyModelValueToEditor,
    syncLastEmittedValue
  };
}
