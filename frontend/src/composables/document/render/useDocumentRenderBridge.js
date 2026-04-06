import { normalizeDocumentType } from '@/utils/documentType';

const MARKDOWN_DOC_TYPE = '1';
const TABLE_DOC_TYPE = '2';
const MARKDOWN_SLIDE_DOC_TYPE = '3';
const YORESEE_RICH_TEXT_DOC_TYPE = '4';

const toText = (value) => {
  if (value === null || value === undefined) {
    return '';
  }
  return String(value);
};

const tryParseJson = (value) => {
  const text = toText(value).trim();
  if (!text) {
    return null;
  }
  try {
    return JSON.parse(text);
  } catch (_) {
    return null;
  }
};

const unwrapTemplateContent = (value) => {
  const rawText = toText(value);
  const parsed = tryParseJson(rawText);
  if (!parsed || typeof parsed !== 'object') {
    return rawText;
  }
  if (typeof parsed.content === 'string') {
    return parsed.content;
  }
  return rawText;
};

const normalizeSource = (value, { isTemplate = false } = {}) =>
  isTemplate ? unwrapTemplateContent(value) : toText(value);

const stripCommentAnchorNoise = (value) =>
  toText(value)
    .replace(/<span\b[^>]*\bdata-comment-anchor-id=(["'])[^"']*\1[^>]*>([\s\S]*?)<\/span>/gi, '$2')
    .replace(/<\/?span\b[^>]*\bdata-comment-anchor-id=(["'])[^"']*\1[^>]*\/?>/gi, '');

const normalizeRows = (rows) => {
  if (!Array.isArray(rows) || rows.length === 0) {
    return [];
  }
  const width = Math.max(1, ...rows.map((row) => (Array.isArray(row) ? row.length : 0)));
  return rows.map((row) => {
    const values = Array.isArray(row) ? row : [];
    return Array.from({ length: width }, (_, index) => {
      const cell = values[index];
      return cell === null || cell === undefined ? '' : String(cell);
    });
  });
};

const parseRowsFromParsedPayload = (parsed) => {
  if (Array.isArray(parsed)) {
    return normalizeRows(parsed);
  }
  if (parsed && Array.isArray(parsed.rows)) {
    return normalizeRows(parsed.rows);
  }
  if (parsed && typeof parsed.content === 'string') {
    return parseTableRows(parsed.content, { isTemplate: false });
  }
  return [];
};

const parseTableRows = (value, options = {}) => {
  const source = normalizeSource(value, options);
  const parsed = tryParseJson(source);
  return parseRowsFromParsedPayload(parsed);
};

const escapeMarkdownCell = (value) =>
  String(value || '')
    .replace(/\|/g, '\\|')
    .replace(/\n/g, '<br />')
    .trim();

const rowsToMarkdownTable = (rows) => {
  const normalizedRows = normalizeRows(rows);
  if (normalizedRows.length === 0) {
    return '';
  }

  const columnCount = Math.max(1, normalizedRows[0]?.length || 0);
  const firstRow = normalizedRows[0] || [];
  const hasHeader = firstRow.some((item) => String(item || '').trim() !== '');
  const header = hasHeader
    ? firstRow
    : Array.from({ length: columnCount }, (_, index) => `Col ${index + 1}`);
  const bodyRows = hasHeader ? normalizedRows.slice(1) : normalizedRows;

  const headerLine = `| ${header.map(escapeMarkdownCell).join(' | ')} |`;
  const separatorLine = `| ${Array.from({ length: columnCount }, () => '---').join(' | ')} |`;
  const bodyLines = bodyRows.map((row) => `| ${row.map(escapeMarkdownCell).join(' | ')} |`);
  return [headerLine, separatorLine, ...bodyLines].join('\n');
};

const rowsToDiffText = (rows) =>
  normalizeRows(rows)
    .map((row, index) => `[R${index + 1}] ${row.map((cell) => toText(cell).replace(/\t/g, ' ')).join('\t')}`)
    .join('\n');

const parseSlideSections = (value, options = {}) => {
  const source = normalizeSource(value, options);
  const parsed = tryParseJson(source);
  if (parsed && Array.isArray(parsed.slides)) {
    return parsed.slides.map((slide, index) => {
      if (typeof slide === 'string') {
        return slide.trim();
      }
      if (slide && typeof slide.content === 'string') {
        return slide.content.trim();
      }
      if (slide && typeof slide.markdown === 'string') {
        return slide.markdown.trim();
      }
      return `# Slide ${index + 1}`;
    }).filter(Boolean);
  }

  return source
    .replace(/\r\n/g, '\n')
    .split(/\n-{3,}\n/g)
    .map((item) => item.trim())
    .filter(Boolean);
};

const slidesToPreviewMarkdown = (sections) => {
  if (!Array.isArray(sections) || sections.length === 0) {
    return '';
  }
  return sections
    .map((section, index) => `## Slide ${index + 1}\n\n${section}`)
    .join('\n\n---\n\n');
};

const slidesToDiffText = (sections) => {
  if (!Array.isArray(sections) || sections.length === 0) {
    return '';
  }
  return sections
    .map((section, index) => `[Slide ${index + 1}]\n${section}`)
    .join('\n\n---\n\n');
};

const buildRendererRegistry = () => ({
  [MARKDOWN_DOC_TYPE]: {
    toPreviewMarkdown: (value, options) => normalizeSource(value, options),
    toDiffText: (value) => toText(value)
  },
  [YORESEE_RICH_TEXT_DOC_TYPE]: {
    toPreviewMarkdown: (value, options) => stripCommentAnchorNoise(normalizeSource(value, options)),
    toDiffText: (value) => stripCommentAnchorNoise(toText(value))
  },
  [TABLE_DOC_TYPE]: {
    toPreviewMarkdown: (value, options) => rowsToMarkdownTable(parseTableRows(value, options)),
    toDiffText: (value) => rowsToDiffText(parseTableRows(value, { isTemplate: false }))
  },
  [MARKDOWN_SLIDE_DOC_TYPE]: {
    toPreviewMarkdown: (value, options) => slidesToPreviewMarkdown(parseSlideSections(value, options)),
    toDiffText: (value) => slidesToDiffText(parseSlideSections(value, { isTemplate: false }))
  }
});

const documentRenderRegistry = buildRendererRegistry();

const resolveRenderer = (documentType) => {
  const normalized = normalizeDocumentType(documentType, MARKDOWN_DOC_TYPE);
  return documentRenderRegistry[normalized] || documentRenderRegistry[MARKDOWN_DOC_TYPE];
};

export const resolveDocumentPreviewKind = (documentType) => {
  const normalized = normalizeDocumentType(documentType, MARKDOWN_DOC_TYPE);
  if (normalized === TABLE_DOC_TYPE) {
    return 'table';
  }
  if (normalized === MARKDOWN_SLIDE_DOC_TYPE) {
    return 'slide';
  }
  return 'markdown';
};

export const resolveTablePreviewRows = ({
  content = '',
  isTemplate = false
} = {}) => parseTableRows(content, { isTemplate });

export const resolveSlidePreviewSections = ({
  content = '',
  isTemplate = false
} = {}) => parseSlideSections(content, { isTemplate });

export const resolveDocumentPreviewContent = ({
  content = '',
  documentType = MARKDOWN_DOC_TYPE,
  isTemplate = false
} = {}) => {
  const renderer = resolveRenderer(documentType);
  return renderer.toPreviewMarkdown(content, { isTemplate });
};

export const resolveDocumentDiffContentPair = ({
  leftContent = '',
  rightContent = '',
  documentType = MARKDOWN_DOC_TYPE
} = {}) => {
  const renderer = resolveRenderer(documentType);
  return {
    // Final guard: strip comment-anchor wrappers for every doc type before diff rendering.
    // This prevents noisy diffs when content accidentally contains inline anchor markup.
    leftText: stripCommentAnchorNoise(renderer.toDiffText(leftContent)),
    rightText: stripCommentAnchorNoise(renderer.toDiffText(rightContent))
  };
};

export const resolveTableDiffRowsPair = ({
  leftContent = '',
  rightContent = ''
} = {}) => ({
  leftRows: parseTableRows(leftContent, { isTemplate: false }),
  rightRows: parseTableRows(rightContent, { isTemplate: false })
});
