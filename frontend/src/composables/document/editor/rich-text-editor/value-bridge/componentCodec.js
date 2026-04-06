export const encodeDrawioSourceToBase64 = (value) => {
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

export const normalizeDrawioSource = (value) => {
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

export const resolveMindmapSource = (node) => {
  if (!node || typeof node.getAttribute !== 'function') {
    return '';
  }
  const rawSource = node.getAttribute('data-source') || node.getAttribute('source') || '';
  if (rawSource) {
    try {
      return decodeURIComponent(String(rawSource));
    } catch (_) {
      return String(rawSource);
    }
  }
  if (typeof node.querySelector === 'function') {
    const textarea = node.querySelector('textarea');
    if (textarea?.value) {
      return String(textarea.value);
    }
  }
  return '';
};

export const resolveDrawioSource = (node) => {
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

export const extractDrawioBlocksFromHtml = (sourceHtml) => {
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

export const extractLanguageFromCodeClass = (className = '') => {
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
