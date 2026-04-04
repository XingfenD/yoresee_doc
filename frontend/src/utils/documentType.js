export const DEFAULT_DOCUMENT_TYPE = '1';

const LEGACY_ALIAS_MAP = {
  markdown: '1',
  document_type_markdown: '1',
  table: '2',
  document_type_table: '2',
  markdown_slide: '3',
  document_type_markdown_slide: '3',
  slide: '3',
  slides: '3',
  document_type_slide: '3'
};

const normalizeRawType = (value) => {
  if (value === null || value === undefined) {
    return '';
  }
  if (typeof value === 'number' && Number.isFinite(value)) {
    return String(Math.trunc(value));
  }
  const text = String(value).trim().toLowerCase();
  return text;
};

export const normalizeDocumentType = (value, fallback = DEFAULT_DOCUMENT_TYPE) => {
  const raw = normalizeRawType(value);
  if (!raw) {
    return fallback;
  }
  if (/^\d+$/.test(raw)) {
    return raw;
  }
  if (LEGACY_ALIAS_MAP[raw]) {
    return LEGACY_ALIAS_MAP[raw];
  }
  return fallback;
};

export const toDocumentTypeNumber = (value, fallback = Number(DEFAULT_DOCUMENT_TYPE)) => {
  const normalized = normalizeDocumentType(value, '');
  if (!normalized) {
    return fallback;
  }
  const parsed = Number(normalized);
  if (!Number.isFinite(parsed)) {
    return fallback;
  }
  return parsed;
};

export const matchDocumentType = (candidate, expected) => {
  const normalizedExpected = normalizeDocumentType(expected, '');
  if (!normalizedExpected) {
    return true;
  }
  return normalizeDocumentType(candidate, '') === normalizedExpected;
};

export const buildDocumentTypeOptions = (t) => [
  { value: '1', label: t('document.typeMarkdown') },
  { value: '2', label: t('document.typeTable') },
  { value: '3', label: t('document.typeSlide') }
];
